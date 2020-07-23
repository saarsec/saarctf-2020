from __future__ import print_function, unicode_literals

from typing import Dict, List, Optional

import requests
import traceback
from gamelib import *
from gamelib import usernames
import random


def assert_valid_cave(cave):
	assert type(cave) == dict, 'Cave must be an object'
	for k in ['id', 'name', 'owner', 'template_id', 'created', 'treasure_count']:
		assert k in cave, 'Cave property missing'


def assert_cave_has_treasure(cave: Dict, name: str):
	assert 'treasures' in cave, '"treasures" key missing in cave'
	for treasure in cave['treasures']:
		if treasure['name'] == name:
			return treasure
	raise AssertionError('Treasure not found')


def assert_cave_has_treasure_containing(cave: Dict, name: str):
	assert 'treasures' in cave, '"treasures" key missing in cave'
	for treasure in cave['treasures']:
		if name in treasure['name']:
			return treasure
	print(json.dumps(cave, indent=4))
	raise AssertionError('Treasure "{}" not found'.format(name))


def assert_cave_in_list(cavelist: List[Dict], idname: str):
	for cave in cavelist:
		if cave['id'] == idname or cave['name'] == idname:
			return cave
	print(cavelist)
	raise AssertionError('Cave {} not found'.format(idname))


def find_cave_in_list(cavelist: List[Dict], idname: str) -> Optional[Dict]:
	for cave in cavelist:
		if cave['id'] == idname or cave['name'] == idname:
			return cave
	print(cavelist)
	return None


class Caveserver:
	def __init__(self, url: str):
		self.url = url
		self.session = requests.Session()
		self.session.timeout = TIMEOUT

	def get(self, uri, **kwargs) -> requests.Response:
		try:
			print(f'GET "{self.url + uri}" ...')
			response = self.session.get(self.url + uri, timeout=TIMEOUT, **kwargs)
			print(' => ', response)
			return response
		except (IOError, ConnectionError):
			raise OfflineException('Could not retrieve page "{}"'.format(self.url + uri))
		except ValueError:
			traceback.print_exc()
			raise MumbleException('Could not retrieve {}'.format(uri))

	def get_unchecked(self, uri, **kwargs) -> requests.Response:
		print(f'GET "{self.url + uri}" ...')
		response = self.session.get(self.url + uri, timeout=TIMEOUT, **kwargs)
		print(' => ', response)
		return response

	def post(self, uri, **kwargs) -> requests.Response:
		print(f'POST "{self.url + uri}" ...')
		response = self.session.post(self.url + uri, timeout=TIMEOUT, **kwargs)
		print(' => ', response)
		return response

	def login(self, username, password):
		try:
			resp = self.post('/api/users/login', json={'username': username, 'password': password})
			assert_requests_response(resp)
			assert resp.json()['username'] == username, 'Wrong username after login'
		except AssertionError as e:
			e.args = ('Could not login: ' + (e.args[0] if e.args else 'AssertionError'),)
			raise
		except (IOError, ConnectionError):
			raise OfflineException('Could not login')
		except ValueError:
			traceback.print_exc()
			raise MumbleException('Could not login as {}'.format(username))

	def register(self, username, password):
		try:
			resp = self.post('/api/users/register', json={"username": username, "password": password})
			assert_requests_response(resp)
			assert resp.json()['username'] == username, 'Wrong username after registration'
		except AssertionError as e:
			e.args = ('Could not register: ' + (e.args[0] if e.args else 'AssertionError'),)
			raise
		except (IOError, ConnectionError):
			raise OfflineException('Could not register')
		except (ValueError, KeyError):
			traceback.print_exc()
			raise MumbleException('Could not register as {}'.format(username))

	def caves_list(self) -> List[Dict]:
		try:
			caves = self.get_unchecked('/api/caves/list').json()
			for cave in caves:
				assert_valid_cave(cave)
			return caves
		except AssertionError as e:
			e.args = ('Could not list caves: ' + (e.args[0] if e.args else 'AssertionError'),)
			raise
		except (IOError, ConnectionError):
			raise OfflineException('Could not list caves')
		except ValueError:
			traceback.print_exc()
			raise MumbleException('Could not list caves')

	def cave_get(self, id) -> Dict:
		try:
			cave = self.get_unchecked('/api/caves/' + id).json()
			assert_valid_cave(cave)
			return cave
		except AssertionError as e:
			e.args = ('Could not retrieve cave: ' + (e.args[0] if e.args else 'AssertionError'),)
			raise
		except (IOError, ConnectionError):
			raise OfflineException('Could not retrieve cave "{}"'.format(id))
		except ValueError:
			traceback.print_exc()
			raise MumbleException('Could not retrieve cave "{}"'.format(id))

	def cave_rent(self, template_id, name) -> Dict:
		try:
			resp = self.post('/api/caves/rent', json={'template': template_id, 'name': name})
			assert_requests_response(resp)
			cave = resp.json()
			assert_valid_cave(cave)
			assert cave['template_id'] == template_id, 'Rent cave: Wrong template'
			assert cave['name'] == name, 'Rent cave: Wrong name'
			return cave
		except AssertionError as e:
			e.args = ('Could not rent cave: ' + (e.args[0] if e.args else 'AssertionError'),)
			raise
		except (IOError, ConnectionError):
			raise OfflineException('Could not rent cave')
		except ValueError:
			traceback.print_exc()
			raise MumbleException('Could not rent cave')

	def cave_hide_treasures(self, id, names) -> List[Dict]:
		"""
		:param str id:
		:param list[str] names:
		:rtype: list[dict]
		:return: Treasures in order of names
		"""
		try:
			resp = self.post('/api/caves/hide-treasures', json={'cave_id': id, 'names': names})
			assert_requests_response(resp)
			cave = resp.json()
			assert_valid_cave(cave)
			assert cave['id'] == id, 'Hide treasure: Wrong cave id'
			treasures = []
			for name in names:
				treasures.append(assert_cave_has_treasure(cave, name))
			return treasures
		except AssertionError as e:
			e.args = ('Could not hide treasures: ' + (e.args[0] if e.args else 'AssertionError'),)
			raise
		except (IOError, ConnectionError):
			raise OfflineException('Could not hide treasures in cave "{}"'.format(id))
		except ValueError:
			traceback.print_exc()
			raise MumbleException('Could not hide treasures in cave "{}"'.format(id))

	def visit(self, id, codefiles) -> str:
		try:
			resp = self.post('/api/visit', json={'cave_id': id, 'files': codefiles})
			assert_requests_response(resp, 'text/plain; charset=utf-8')
			return resp.text
		except AssertionError as e:
			e.args = ('Could not visit cave: ' + (e.args[0] if e.args else 'AssertionError'),)
			raise
		except (IOError, ConnectionError):
			raise OfflineException('Could not visit cave "{}"'.format(id))


def random_flag_prefix() -> str:
	return random.choice(['Box', 'Chest', 'Barrel', 'Hidden corner']) + ' ' + random.choice(['containing', 'with', 'of']) + ' '


if __name__ == '__main__':
	try:
		url = 'http://localhost:9080'
		username = usernames.generate_username()
		pwd = usernames.generate_password()
		print('Password =', pwd)

		server = Caveserver(url)
		server.register(username, pwd)

		server = Caveserver(url)
		server.login(username, pwd)

		cave = server.cave_rent(random.randint(1, 49), 'ScriptTestCave')
		assert_cave_in_list(server.caves_list(), cave['id'])

		server.cave_get(cave['id'])

		for t in server.cave_hide_treasures(cave['id'], [
			usernames.generate_dummy_flag(),
			usernames.generate_password(),
			random_flag_prefix() + usernames.generate_dummy_flag(),
			usernames.generate_password(),
			usernames.generate_dummy_flag(),
		]):
			print('-', t['x'], ', ', t['y'])

		for t in server.cave_get(cave['id'])['treasures']:
			print('-', t['x'], ', ', t['y'])

		print(server.visit(cave['id'], {'entry.sl': 'eija error; ajwdkljaw'}))
	except OfflineException as e:
		print('Users would see: [OFFLINE]', e.message)
	except MumbleException as e:
		traceback.print_exc()
		print("Users would see: [MUMBLE]", e.message)
	print('Done')
