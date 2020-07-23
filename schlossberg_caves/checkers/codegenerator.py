import random
import string
from collections import OrderedDict
from typing import List, Union, Dict, Tuple

try:
	from . import cavelib
except ImportError:
	import cavelib


class SaarlangStatement:
	def __init__(self, text: str, subscope: List['SaarlangStatement'] = None, uses: List[str] = None, defines=False) -> None:
		self.name = None  # type: str
		self.text = text  # type: str
		self.subscope = subscope or []  # type: List[SaarlangStatement]
		self.used_symbols = uses or []  # type: List[str]
		self.is_definition = defines

	def to_list(self) -> List[str]:
		return [self.text]

	def subscope_list(self, indent: int = 0) -> List[str]:
		result = []
		for stmt in self.subscope:
			for line in stmt.to_list():
				result.append('\t' * indent + line)
		return result

	def uses(self, symbol: str) -> bool:
		if symbol in self.used_symbols:
			return True
		for stmt in self.subscope:
			if stmt.uses(symbol):
				return True
		return False

	def __str__(self):
		return '\n'.join(self.to_list())

	def __repr__(self):
		return '\n'.join(self.to_list())


class Stmt(SaarlangStatement):
	def __init__(self, text: str, *args) -> None:
		SaarlangStatement.__init__(self, text.format(*args), uses=list(args))


class WhileStmt(SaarlangStatement):
	def to_list(self):
		return ['solang ' + self.text + ': {'] + self.subscope_list(1) + ['}']


class IfStmt(SaarlangStatement):
	def to_list(self):
		return ['falls ' + self.text + ': {'] + self.subscope_list(1) + ['}']


class MultiStmt(SaarlangStatement):
	def __init__(self, *args) -> None:
		SaarlangStatement.__init__(self, '', subscope=list(args))

	def to_list(self):
		return self.subscope_list()


class SLConstant(SaarlangStatement):
	def __init__(self, value: str, name: str) -> None:
		SaarlangStatement.__init__(self, value)
		self.name = name

	def to_list(self):
		return ['const {}: int = {};'.format(self.name, self.text)]


class SLVariableDef(SaarlangStatement):
	def __init__(self, value: str, name: str, vartype: str = 'int') -> None:
		SaarlangStatement.__init__(self, value)
		self.name = name
		self.type = vartype

	def to_list(self):
		return ['var {}: {} = {};'.format(self.name, self.type, self.text)]


class SLFunction(SaarlangStatement):
	def __init__(self, text: str, name: str, subscope: List['SaarlangStatement'] = None) -> None:
		SaarlangStatement.__init__(self, text, subscope=subscope)
		self.name = name

	def to_list(self):
		return [self.text + ': {'] + self.subscope_list(1) + ['}']


class SLModule(SaarlangStatement):
	def __init__(self, name: str, subscope: List['SaarlangStatement'] = None) -> None:
		SaarlangStatement.__init__(self, '', subscope=subscope)
		self.name = name

	def to_list(self):
		result = []
		for definition in self.subscope:
			result.append('\n'.join(definition.to_list()))
		return ['\n\n'.join(result)]

	def add_import(self, filename):
		self.subscope.insert(0, Stmt('holmol "{}";', filename))


class Variable:
	def __init__(self, name: str, vartype: str, value: Union[int, List[int]]) -> None:
		self.name = name
		self.type = vartype
		if self.is_array:
			assert isinstance(value, list)
			self.value = [cast_to_type(v, self.basic_type) for v in value]  # type: Union[int, List[int]]
		else:
			assert isinstance(value, int)
			self.value = cast_to_type(value, vartype)  # type: Union[int, List[int]]

	@property
	def is_array(self):
		return self.type.startswith('lischd ')

	@property
	def basic_type(self):
		return self.type[7:]

	def __len__(self):
		return len(self.value) if self.is_array else 0

	def __getitem__(self, item):
		return self.value[item]

	def __setitem__(self, key, value):
		self.value[key] = value


RUFF = SaarlangStatement('ruff;')
RUNNER = SaarlangStatement('runner;')
DONIWWER = SaarlangStatement('doniwwer;')
RIWWER = SaarlangStatement('riwwer;')
FERDISCH = SaarlangStatement('ferdisch;')


def direction_to_saarlang(direction: int) -> SaarlangStatement:
	if direction == cavelib.UP:
		return RUFF
	elif direction == cavelib.RIGHT:
		return DONIWWER
	elif direction == cavelib.DOWN:
		return RUNNER
	elif direction == cavelib.LEFT:
		return RIWWER
	raise Exception('Invalid direction!')


def cast_to_type(value: int, vartype: str) -> int:
	if vartype == 'int':
		value &= 0xffffffffffffffff
		if value & 0x8000000000000000:
			value -= 0x10000000000000000
		return value
	elif vartype == 'byte':
		return value & 0xff
	else:
		raise Exception('Invalid type: ' + vartype)


def number_to_saarlang(num: int) -> str:
	if num < 0:
		return '0-' + str(abs(num))
	else:
		return str(num)


class CodeGenerator:
	def __init__(self):
		self.use_arithmetic_for_constants = True
		self.use_for_compression = True
		self.use_function_compression = True
		self.use_function_split = True
		self.use_multi_modules = True
		self.use_constants = True
		self.use_comments_in_expr = True
		self.use_variables = True
		self.use_stdlib = True
		self.use_additional_imports = True
		self.known_parameters = {}
		self.known_constants = {}
		self.known_variables = []  # type: List[Variable]
		self._output = []  # type: List[str]
		self.__known_symbols = set()  # type: set

	@property
	def expected_output(self) -> str:
		return ''.join(self._output)

	def random_symbol(self, charset: str, min_length: int, max_length: int) -> str:
		length = random.randint(min_length, max_length)
		while True:
			symbol = ''.join(random.choice(charset) for _ in range(length))
			if symbol not in self.__known_symbols:
				self.__known_symbols.add(symbol)
				return symbol

	def random_var(self, min_length=8, max_length=12) -> str:
		return self.random_symbol(string.ascii_lowercase, min_length, max_length)

	def random_param(self, min_length=5, max_length=8) -> str:
		return self.random_symbol(string.ascii_lowercase, min_length, max_length)

	def random_funcname(self, min_length=8, max_length=12) -> str:
		return self.random_symbol(string.ascii_lowercase, min_length, max_length)

	def random_filename(self, min_length=5, max_length=8) -> str:
		return self.random_symbol(string.ascii_lowercase, min_length, max_length) + '.sl'

	def random_string(self, min_length=4, max_length=20) -> str:
		charset = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!"#$%&\'()*+,-./:;<=>?@[]^_`{|}~'
		length = random.randint(min_length, max_length)
		return ''.join(random.choice(charset) for _ in range(length))

	def random_comment(self):
		length = random.randint(10, 300)
		if random.randint(0, 1):
			return '//' + ''.join(random.choice(string.ascii_letters + string.digits + ',.- #+/<>|~!"§$%&()=?`;:_') for _ in range(length)) + '\n'
		else:
			return '/*' + ''.join(random.choice(string.ascii_letters + string.digits + ',.- #+/<>|~!"§$%&()=?`;:_\n') for _ in range(length)) + '*/'

	def random_constant_expr(self, c: int) -> str:
		if self.use_arithmetic_for_constants and random.randint(0, 99) < 70:
			numbers = []  # type: List[Union[int, str]]
			s = 0
			for _ in range(random.randint(1, 5)):
				n = random.randint(-1000, 1000)
				if len(numbers) == 0:
					n = abs(n)
					numbers.append(n)
					s += n
				elif random.randint(0, 100) < 15:
					r = random.randint(1, 13)
					n = abs(n)
					if r == 0:
						n2 = random.randint(0, 100)
						numbers.append('{}*{}'.format(n, n2))
						s += n * n2
					elif r == 1:
						n2 = random.randint(0, 100)
						numbers.append('{}*{}'.format(n2, n))
						s += n * n2
					elif r == 2:
						n2 = random.randint(1, 20)
						numbers.append('{}/{}'.format(n, n2))
						s += n // n2
					elif r == 3:
						n2 = random.randint(1, 20)
						numbers.append('{}%{}'.format(n, n2))
						s += n % n2
					elif r == 4:
						large_nums = [random.randint(0x3fffffffffffffff, 0x7fffffffffffffff)]
						while sum(large_nums) <= 0xffffffffffffffff:
							large_nums.append(random.randint(0x3fffffffffffffff, 0x7fffffffffffffff))
						numbers += large_nums
						s += sum(large_nums)
					elif r == 5 and self.use_comments_in_expr:
						numbers.append(str(n) + self.random_comment())
						s += n
					elif r % 4 == 0 and self.known_parameters:
						param, value = random.choice(list(self.known_parameters.items()))
						numbers.append(param)
						s += value
					elif r % 4 <= 1 and self.known_constants:
						param, value = random.choice(list(self.known_constants.items()))
						numbers.append(param)
						s += value
					elif self.known_variables:
						var = random.choice(self.known_variables)
						if var.is_array:
							index = random.randint(-1, len(var) - 1)
							if index < 0:
								numbers.append('grees {}'.format(var.name))
								s += len(var)
							else:
								numbers.append('{}@{}'.format(var.name, index))
								s += var[index]
						else:
							numbers.append(var.name)
							s += var.value
					else:
						numbers.append(n)
						s += n
				else:
					numbers.append(n)
					s += n
			s = cast_to_type(s, 'int')
			numbers.append(c - s)
			return '(' + ' + '.join(map(str, numbers)).replace('+ -', '- ') + ')'
		else:
			return number_to_saarlang(c)

	def simple_commands(self, path: List[int]) -> List[SaarlangStatement]:
		return [direction_to_saarlang(direction) for direction in path]

	def path_to_saarlang_simple(self, path: List[Union[int, SaarlangStatement, None]]) -> List[SaarlangStatement]:
		commands = []
		for p in path:
			if p is None:
				pass
			elif isinstance(p, int):
				commands.append(direction_to_saarlang(p))
			else:
				commands.append(p)
		return commands

	def path_to_saarlang(self, path: List[Union[int, SaarlangStatement, None]]) -> List[SaarlangStatement]:
		if len(path) < 5 or not self.use_for_compression:
			return self.path_to_saarlang_simple(path)
		patterns = []
		for l in range(3, 21):
			i = 0
			while i <= len(path) - l:
				if path[i:i + l] == path[i + l:i + l + l]:
					c = 2
					while path[i:i + l] == path[i + c * l:i + c * l + l]:
						c += 1
					patterns.append((i, l, c))
					i += c * l
				else:
					i += 1
		patterns.sort(key=lambda p: p[1] * (p[2] - 1), reverse=True)
		if not patterns:
			return self.path_to_saarlang_simple(path)
		for i, l, c in patterns:
			if len([j for j in path[i:i + l * c] if type(j) != int]) > 0:
				continue
			varname = self.random_var()
			substmts = self.path_to_saarlang(path[i:i + l])
			offset = random.randint(0, 10)
			path[i] = MultiStmt(
				SaarlangStatement('var {}: int = {};'.format(varname, self.random_constant_expr(c + offset))),
				WhileStmt(varname + ' > ' + self.random_constant_expr(offset), substmts + [Stmt('{0} = {0} - 1;', varname)])
			)
			for j in range(i + 1, i + l * c):
				path[j] = None
		return self.path_to_saarlang_simple(path)

	def path_to_functions(self, path: List[int]) -> Tuple[List[SLFunction], List[SaarlangStatement]]:
		class MyDefaultDict(dict):
			def __missing__(self, key):
				self[key] = (len(key), [])
				return self[key]

		patterns = MyDefaultDict()
		pathstr = ''.join(map(str, path))
		for l in range(8, 20):
			for i in range(len(path) - l):
				patterns[pathstr[i:i + l]][1].append(i)
		pattern_list = [p for p in patterns.values() if len(p[1]) >= 4]
		pattern_list.sort(key=lambda p: p[0] * (len(p[1])) - 1, reverse=True)

		new_path = path  # type: List[Union[int, SaarlangStatement, None]]
		functions = []
		for length, locations in pattern_list:
			some_pos = locations[0]
			if new_path[some_pos] is None or new_path[some_pos + length] is None:
				continue

			# parameters (functions are always called with fixed parameter values)
			old_known_parameters = self.known_parameters
			self.known_parameters = OrderedDict()
			for _ in range(max(0, random.randint(-5, 5))):
				self.known_parameters[self.random_param()] = random.randint(0, 10000)
			func_sig = ', '.join('{}: int'.format(param) for param in self.known_parameters.keys())
			call_sig = ', '.join(map(str, self.known_parameters.values()))

			# create function and reset parameters afterwards
			func_name = self.random_funcname()
			func_body = self.path_to_saarlang(new_path[some_pos:some_pos + length])
			self.known_parameters = old_known_parameters
			functions.append(SLFunction('eija {}({}) gebbtserick int'.format(func_name, func_sig), func_name, subscope=func_body))

			# create calls
			for pos in locations:
				if new_path[pos] is None or new_path[pos + length] is None:
					continue
				new_path[pos] = Stmt('mach {}({});', func_name, call_sig)
				for j in range(pos + 1, pos + length):
					new_path[j] = None

		return functions, self.path_to_saarlang(new_path)

	def add_conditionals_variables_to_function(self, func: SLFunction):
		if len(func.subscope) < 10:
			return

		old_commands = list(func.subscope)
		commands = []  # type: List[SaarlangStatement]
		self.known_variables = []

		while len(old_commands) > 0:
			x = min(random.randint(1, 15), len(old_commands))
			portion = old_commands[:x]
			old_commands = old_commands[x:]
			portion_contains_definitions = any(stmt.is_definition for stmt in portion)
			rand = random.randint(0, 100)
			if rand < 20 and len(self.known_variables) > 0:
				# reassign var
				var = random.choice(self.known_variables)
				if var.is_array:
					if rand < 3:
						varname2 = self.random_var()
						commands.append(SLVariableDef(var.name, varname2, var.type))
						var.name = varname2
					else:
						for i, v in enumerate(var.value):
							if random.randint(0, 99) < 30:
								v2 = random.randint(-1024, 1024)
								commands.append(SaarlangStatement('{}@{} = {};'.format(var.name, i, self.random_constant_expr(v2))))
								var[i] = cast_to_type(v2, var.basic_type)
					commands.extend(portion)
				else:
					value2 = random.randint(-1024, 1024)
					commands.append(SaarlangStatement('{} = {};'.format(var.name, self.random_constant_expr(value2))))
					commands.extend(portion)
					var.value = cast_to_type(value2, var.type)

			elif rand < 40:
				# define var
				varname = self.random_var()
				vartype = 'int' if random.randint(0, 9) < 7 else 'byte'

				if random.randint(0, 9) < 3:
					vartype = 'lischd ' + vartype
					value = [random.randint(-1000, 1000) for _ in range(random.randint(1, 20))]
					if rand % 2 == 0:
						commands.append(SLVariableDef('neie {} ({})'.format(vartype, len(value)), varname, vartype))
					else:
						commands.append(SaarlangStatement('var {}: {}({});'.format(varname, vartype, len(value))))
					for i, v in enumerate(value):
						commands.append(SaarlangStatement('{}@{} = {};'.format(varname, i, self.random_constant_expr(v))))
				else:
					value = random.randint(-1000, 1000)
					commands.append(SLVariableDef(self.random_constant_expr(value), varname, vartype))

				commands.extend(portion)
				self.known_variables.append(Variable(varname, vartype, value))

			elif rand < 60 and not portion_contains_definitions:
				# true if (+reassign)
				v = random.randint(-1000, 1000)
				if rand % 5 == 0:
					cond = '{} == {}'.format(self.random_constant_expr(v), number_to_saarlang(v))
				elif rand % 5 == 1:
					cond = '{} > {}'.format(self.random_constant_expr(v), number_to_saarlang(v - 1))
				elif rand % 5 == 2:
					cond = '{} == {}'.format(self.random_constant_expr(v), self.random_constant_expr(v))
				elif rand % 5 == 3:
					cond = '{} != {}'.format(self.random_constant_expr(v), self.random_constant_expr(v + random.randint(1, 1000)))
				else:
					cond = '{} < {}'.format(self.random_constant_expr(v), number_to_saarlang(v + 1))
				ifstmt = IfStmt(cond, portion)
				if rand % 4 < 3 and len(self.known_variables) > 0:
					var = random.choice(self.known_variables)
					value2 = random.randint(-1024, 1024)
					if var.is_array:
						idx = random.randint(0, len(var) - 1)
						ifstmt.subscope.append(SaarlangStatement('{}@{} = {};'.format(var.name, idx, self.random_constant_expr(value2))))
						var[idx] = cast_to_type(value2, var.basic_type)
					else:
						ifstmt.subscope.append(SaarlangStatement('{} = {};'.format(var.name, self.random_constant_expr(value2))))
						var.value = cast_to_type(value2, var.type)
				commands.append(ifstmt)

			elif rand < 80 and not portion_contains_definitions:
				# false if (+reassign)
				v = random.randint(-1000, 1000)
				if rand % 5 == 0:
					cond = '{} == {}'.format(self.random_constant_expr(v), number_to_saarlang(v - 1))
				elif rand % 5 == 1:
					cond = '{} > {}'.format(self.random_constant_expr(v), number_to_saarlang(v))
				elif rand % 5 == 2:
					cond = '{} == {}'.format(self.random_constant_expr(v), self.random_constant_expr(v + random.randint(1, 1000)))
				elif rand % 5 == 3:
					cond = '{} != {}'.format(self.random_constant_expr(v), self.random_constant_expr(v))
				else:
					cond = '{} < {}'.format(self.random_constant_expr(v), number_to_saarlang(v - 1))
				ifstmt = IfStmt(cond, [random.choice((RUFF, RUNNER, RIWWER, DONIWWER)) for _ in range(random.randint(1, 12))])
				if rand % 4 < 3 and len(self.known_variables) > 0:
					var = random.choice(self.known_variables)
					value2 = random.randint(-1024, 1024)
					if var.is_array:
						idx = random.randint(0, len(var) - 1)
						ifstmt.subscope.append(SaarlangStatement('{}@{} = {};'.format(var.name, idx, self.random_constant_expr(value2))))
					else:
						ifstmt.subscope.append(SaarlangStatement('{} = {};'.format(var.name, self.random_constant_expr(value2))))
				commands.append(ifstmt)
				commands.extend(portion)
			else:
				commands.extend(portion)
		func.subscope = commands
		self.known_variables = []

	def add_output_to_path(self, commands: List[SaarlangStatement]):
		indices = [random.randint(0, len(commands) - 1) for _ in range(random.randint(1, 5))]
		indices.sort()
		count = 0
		for index in indices:
			rand = random.randint(0, 99)
			if rand < 50:
				func, sep = ('sahmol', ' ') if rand < 30 else ('sahmol_ln', '\n')
				value = random.randint(-10000, 10000)
				stmt = SaarlangStatement('mach {} ({});'.format(func, self.random_constant_expr(value)), uses=['stdlib.sl'])
				commands.insert(index + count, stmt)
				count += 1
				self._output.append(str(value) + sep)
			else:
				varname = self.random_var()
				text = self.random_string()
				stmts = [SaarlangStatement('var {}: lischd byte ({});'.format(varname, len(text) + 1), defines=True)]
				for i, c in enumerate(text):
					stmts.append(SaarlangStatement('{}@{} = {};'.format(varname, i, self.random_constant_expr(ord(c)))))
				stmts.append(SaarlangStatement('{}@{} = {};'.format(varname, len(text), 0)))
				stmts.append(SaarlangStatement('mach sahmol_as_str({});'.format(varname), uses=['stdlib.sl']))
				commands.insert(index + count, MultiStmt(*stmts))
				count += 1
				self._output.append(text + '\n')

	def split_commandlist_to_functions(self, commands: List[SaarlangStatement]) -> Tuple[List[SLFunction], List[SaarlangStatement]]:
		functions = []
		for _ in range(random.randint(5, 10)):
			if len(commands) < 30:
				continue
			pos = random.randint(0, len(commands) - 20)
			length = random.randint(15, min(50, len(commands) - pos - 1))
			func_name = self.random_funcname()
			functions.append(SLFunction('eija {}() gebbtserick int'.format(func_name), func_name, subscope=commands[pos:pos + length]))
			commands = commands[:pos] + [Stmt('mach {}();', func_name)] + commands[pos + length:]
		return functions, commands

	def split_functions_to_modules(self, functions: List[SLFunction], constants: List[SLConstant] = None) -> List[SLModule]:
		modules = [SLModule('entry.sl')]
		for _ in range(random.randint(1, 4)):
			modules.append(SLModule(self.random_filename()))

		if constants:
			random.choice(modules).subscope.extend(constants)
			for m in modules:
				m.used_symbols.extend(c.name for c in constants)

		for f in functions:
			if 'main(' in f.text:
				modules[0].subscope.append(f)
			else:
				modules[random.randint(0, len(modules) - 1)].subscope.append(f)
		# imports
		for m1 in modules:
			for m2 in modules:
				if m1.name != m2.name:
					# check if m1 has to import m2
					if any(f.name and m1.uses(f.name) for f in m2.subscope):
						m1.add_import(m2.name)
		return modules

	def add_imports(self, module: SLModule):
		for f in ['stdlib.sl', 'cavelib.sl']:
			if module.uses(f) or (self.use_additional_imports and random.randint(0, 10) == 0):
				module.add_import(f)

	def generate_modules(self, path: List[int]) -> List[SLModule]:
		constants = []  # type: List[SLConstant]
		if self.use_constants:
			for _ in range(max(0, random.randint(-2, 6))):
				cname = self.random_var()
				cvalue = random.randint(0, 20000)
				constants.append(SLConstant(self.random_constant_expr(cvalue), cname))
				self.known_constants[cname] = cvalue
		if self.use_function_compression:
			functions, commands = self.path_to_functions(path)
		else:
			functions, commands = [], self.path_to_saarlang(path)

		if self.use_stdlib:
			self.add_output_to_path(commands)

		if self.use_function_split:
			f2, commands = self.split_commandlist_to_functions(commands)
			functions += f2
		commands.append(FERDISCH)
		functions.append(SLFunction('eija main() gebbtserick int', 'main', commands))

		if self.use_variables:
			for func in functions:
				if random.randint(0, 100) < 75:
					self.add_conditionals_variables_to_function(func)

		if self.use_multi_modules:
			modules = self.split_functions_to_modules(functions, constants)
		else:
			modules = [SLModule('entry.sl', constants + functions)]
		for m in modules:
			self.add_imports(m)
		return modules

	def generate_code(self, path: List[int]) -> Dict[str, str]:
		return {m.name: str(m) for m in self.generate_modules(path)}


def program(path: List[int]) -> Tuple[Dict[str, str], str]:
	codegen = CodeGenerator()
	codegen.use_arithmetic_for_constants = random.randint(0, 99) < 90
	codegen.use_for_compression = random.randint(0, 99) < 70
	codegen.use_function_compression = random.randint(0, 99) < 50
	codegen.use_function_split = random.randint(0, 99) < 50
	codegen.use_multi_modules = random.randint(0, 99) < 50
	codegen.use_stdlib = random.randint(0, 99) < 50
	return codegen.generate_code(path), codegen.expected_output
