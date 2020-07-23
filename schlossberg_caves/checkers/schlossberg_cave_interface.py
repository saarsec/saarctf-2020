import random
import sys
import time
from typing import Dict, Tuple

from gamelib import usernames
from gamelib import *
from gamelib.gamelogger import GameLogger

CHECK_CODE_HASHES = False

try:
	from . import caveserver
	from . import cavelib
	from . import codegenerator
except (ImportError, SystemError):
	traceback.print_exc()
	from gameserver import caveserver, cavelib, codegenerator

if CHECK_CODE_HASHES:
	try:
		from . import cpp_saarlang
	except (ImportError, SystemError):
		from gameserver import cpp_saarlang


class SchlossbergCaveServiceInterface(ServiceInterface):
	name = 'SchlossbergCaves'

	def server(self, team):
		"""
		:param team:
		:rtype caveserver.Caveserver
		:return:
		"""
		return caveserver.Caveserver('http://{}:9080'.format(team.ip))
		#if team.ip.startswith('127.') and not ':' in team.ip:
		#	return caveserver.Caveserver('http://{}:8080'.format(team.ip))
		#else:
		#	return caveserver.Caveserver('http://{}/schlossberg'.format(team.ip))

	def check_integrity(self, team, round):
		"""
		Do integrity checks that are not related to flags (checking the frontpage, or if exploit-relevant functions are still available)
		:param Team team:
		:param int round:
		:raises MumbleException: Service is broken
		:raises OfflineException: Service is not reachable
		:return:
		"""
		server = self.server(team)
		caveserver.assert_requests_response(server.get('/'), 'text/html; charset=utf-8')

	def store_flags(self, team, round):
		"""
		Send one or multiple flags to a given team. You can perform additional functionality checks here.
		:param Team team:
		:param int round:
		:raises MumbleException: Service is broken
		:raises OfflineException: Service is not reachable
		:return: number of stored flags
		"""
		username = usernames.generate_username()
		password = usernames.generate_password()
		cavename = usernames.generate_name()
		self.store(team, round, 'cave', [username, password, cavename])

		server = self.server(team)
		server.register(username, password)

		server = self.server(team)
		server.login(username, password)

		cave = server.cave_rent(random.randint(1, 49), cavename)
		caveserver.assert_cave_in_list(server.caves_list(), cave['id'])

		server.cave_get(cave['id'])

		# Prepare flags that we want to store
		treasures = [self.get_flag(team, round, 0)]
		for i in range(random.randint(1, 3)):
			treasures.append(usernames.generate_password())
		treasures.append(caveserver.random_flag_prefix() + self.get_flag(team, round, 1))
		for i in range(random.randint(1, 3)):
			treasures.append(usernames.generate_password())
		treasures.append(self.get_flag(team, round, 2))

		# check treasure list on new cave
		server.cave_hide_treasures(cave['id'], treasures)
		self.store(team, round, 'cave_id', cave['id'])
		cave = server.cave_get(cave['id'])
		for t in treasures:
			caveserver.assert_cave_has_treasure(cave, t)
		return 3

	def retrieve_flags(self, team, round):
		"""
		Retrieve all flags stored in a previous round from a given team. You can perform additional functionality checks here.
		:param Team team:
		:param int round:
		:raises FlagMissingException: Flag could not be retrieved
		:raises MumbleException: Service is broken
		:raises OfflineException: Service is not reachable
		:return:
		"""
		cave_id = self.load(team, round, 'cave_id')
		if not cave_id:
			raise FlagMissingException('Flag never stored')
		username, password, cavename = self.load(team, round, 'cave')

		# check cave list
		server = self.server(team)
		caves = server.caves_list()
		if caveserver.find_cave_in_list(caves, cavename) is None:
			raise FlagMissingException('Cave {} not found'.format(cavename))

		server = self.server(team)
		server.login(username, password)
		cave = server.cave_get(cave_id)
		treasures = []
		# for i in range(3):
		#	flag = self.get_flag(team, round, i)
		#	treasures.append(caveserver.assert_cave_has_treasure_containing(cave, flag))
		flags = set()
		for treasure in cave.get('treasures', []):
			found_flags = self.search_flags(treasure['name'])
			flags |= found_flags
			if found_flags:
				treasures.append(treasure)
		payloads = set()
		for flag in flags:
			_, _, _, payload = self.check_flag(flag, team.id, round)
			if payload is not None:
				payloads.add(payload)
		missing_flags = len(set(range(3)) - payloads)
		if missing_flags > 0:
			raise FlagMissingException('{} treasures missing'.format(missing_flags))

		# check multiple visit scripts
		random.shuffle(treasures)
		for treasure in treasures:
			server = self.server(team)
			scripts, expected_output = self.generate_visitor_script(cave['template_id'], treasure['x'], treasure['y'])
			try:
				output = server.visit(cave_id, scripts)
				print(
					'------------ SAARLANG RESPONSE ------------\n' +
					output.replace('\x00', '<\\x00>') +
					'\n-------------------------------------------'
				)
				# parse output
				if output == 'Abort.':
					raise AssertionError('Abort. (compiler crash)')
				output_saarlang_start = output.find('--- Saarlang execution starts ---')
				if output_saarlang_start < 0:
					raise AssertionError('Compiler error (code not executed)')
				output_data_start = output.find('VISITED PATH:')
				if output_data_start < 0:
					match = re.search(r'ERROR: ([^\n]+)\n', output)
					if match:
						raise AssertionError('Execution interrupted, message: "{}"'.format(match.group(1)))
					raise AssertionError('Execution never finished (crash?)')
				output_data = json.loads(output[output_data_start + 13:])
				# check stdout output
				if expected_output and expected_output not in output:
					print('Expected output:', repr(expected_output))
					print('But was:        ', repr(output[output_saarlang_start:output_data_start]))
					raise AssertionError('printed output is wrong ("sahmol")')
				# check if treasure has been found
				if not any((t['name'] == treasure['name'] for t in output_data['treasures'])):
					raise AssertionError('Could not retrieve flag')
				if CHECK_CODE_HASHES:
					p1 = output.find('CODE SIGNATURES:') + len('CODE SIGNATURES:')
					server_signatures = json.loads(output[p1:output_saarlang_start])
					expected_signatures = cpp_saarlang.getSignatures(scripts)
					for fn, sig in expected_signatures.items():
						if fn not in server_signatures or server_signatures[fn] != sig:
							raise MumbleException('Invalid signature (expected {}, was {})'.format(sig, server_signatures[fn]))
				else:
					print('Signatures not checked.')
				# TODO check path
				print('\n')
			except json.JSONDecodeError:
				raise MumbleException('Signature or result is not valid json')
			except KeyError:
				raise MumbleException('Strange result json format')
			except:
				time.sleep(0.1)
				print('=== Scripts ===')
				print('Cave {}, template_id {}'.format(cave['id'], cave['template_id']))
				for name, text in scripts.items():
					print('---', name, '---')
					print(text.replace('\x00', '<\\x00>'))
					print('\n\n')
				time.sleep(0.1)
				raise
		return len(payloads)

	def generate_visitor_script(self, template_id: int, x: int, y: int) -> Tuple[Dict[str, str], str]:
		t = time.time()
		cave = cavelib.CaveCache['schlossberg_{}.cave'.format(template_id)]
		reverse_path = cave.getRandomShortestPathToStart((x, y))
		path = []
		for direction in reverse_path[::-1]:
			if direction == cavelib.UP:
				path.append(cavelib.DOWN)
			elif direction == cavelib.RIGHT:
				path.append(cavelib.LEFT)
			elif direction == cavelib.DOWN:
				path.append(cavelib.UP)
			elif direction == cavelib.LEFT:
				path.append(cavelib.RIGHT)

		code = codegenerator.program(path)
		t = time.time() - t
		GameLogger.log('Codegen time: ', t * 1000, 'ms')
		return code


if __name__ == '__main__':
	IP = '192.168.56.101'
	ITERATIONS = 1
	RETRIES = 3

	for _ in range(ITERATIONS):
		team = Team(71, 'saarsec', IP)
		round = int(sys.argv[1]) if len(sys.argv) > 1 else 2
		service = SchlossbergCaveServiceInterface(15)

		print('[1] Integrity check...')
		service.check_integrity(team, round)
		print('Passed.')

		print('[2] Store flags...')
		flags = service.store_flags(team, round)
		print('Done ({} flags).'.format(flags))

		print('[3] Retrieve the flags in the next round')
		for _ in range(RETRIES):
			flags = service.retrieve_flags(team, round)
		print('Done ({} flags).'.format(flags))
