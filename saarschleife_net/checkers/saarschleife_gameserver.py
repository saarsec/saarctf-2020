import html
import random
import sys
from typing import Dict, List, Tuple

import requests

try:
	from gamelib import *
	from gamelib import usernames
	from gamelib.gamelogger import GameLogger
except ImportError:
	from .gamelib import *
	from .gamelib import usernames
	from .gamelib.gamelogger import GameLogger


class SaarschleifeConnection:
	def __init__(self, ip: str):
		self.ip = ip
		self.session = requests.session()
		self.session.timeout = TIMEOUT

	def get(self, url, **kwargs):
		if url.startswith('/'):
			url = 'http://' + self.ip + ':8080' + url
		GameLogger.log(f'GET "{url}" ...')
		response = self.session.get(url, timeout=TIMEOUT, **kwargs)
		GameLogger.log(' => ', response)
		return response

	def post(self, url, **kwargs):
		GameLogger.log(f'POST "{url}" ...')
		response = self.session.post('http://' + self.ip + ':8080' + url, timeout=TIMEOUT, **kwargs)
		GameLogger.log(' => ', response)
		return response

	def login(self, username: str, password: str):
		resp = self.post('/api/login-register', json={'username': username, 'password': password})
		assert_requests_response(resp, 'text/plain; charset=utf-8')

	def logout(self):
		resp = self.post('/api/logout', json={})
		assert_requests_response(resp, 'text/plain; charset=utf-8')

	def list_visitors(self) -> List[str]:
		resp = self.get('/visitor')
		assert_requests_response(resp, 'text/html; charset=utf-8')
		return [m.group(1) for m in re.finditer(r' data-username="([^"]+)"', resp.text)]

	def me(self) -> Dict:
		resp = self.get('/api/me')
		assert_requests_response(resp)
		return resp.json()

	def get_visitor(self, username: str) -> Dict:
		resp = self.get('/api/visitor/{}'.format(username))
		assert_requests_response(resp)
		return resp.json()

	def befriend(self, username: str) -> str:
		resp = self.post('/api/visitor/friends/{}'.format(username), json={})
		assert_requests_response(resp, 'text/plain; charset=utf-8')
		return resp.text

	def statusmessage_set(self, title: str, text: str, is_public: bool):
		resp = self.post('/api/statusmessage/set', json={
			'title': title, 'text': text, 'isPublic': is_public
		})
		assert_requests_response(resp, 'text/plain; charset=utf-8')

	def statusmessage_get(self, username: str) -> Dict:
		resp = self.get('/api/statusmessage/{}'.format(username))
		assert_requests_response(resp)
		return resp.json()

	def statusmessage_remove(self):
		resp = self.post('/api/statusmessage/remove', json={})
		assert_requests_response(resp, 'text/plain; charset=utf-8')

	def statusmessage_spread(self, username: str):
		resp = self.post('/api/statusmessage/spread/{}'.format(username), json={})
		assert_requests_response(resp, 'text/plain; charset=utf-8')
		return resp.text

	def secret_found(self, name: str, solution: str) -> str:
		resp = self.post('/api/secret/found', json={'name': name, 'solution': solution})
		assert_requests_response(resp, 'text/plain; charset=utf-8')
		return resp.text


def assert_valid_statusmessage(statusmsg) -> Dict:
	assert type(statusmsg) == dict, 'Statusmessage must be an object'
	for k in ['title', 'text', 'isPublic', 'writtenBy']:
		assert k in statusmsg, 'Statusmessage is lacking property ' + k
	return statusmsg


def assert_valid_visitor(visitor) -> Dict:
	assert type(visitor) == dict, 'Visitor must be an object'
	for k in ["username", "points", "statusMessage", "isFriend", "created"]:
		assert k in visitor, 'Visitor is lacking properties'
	return visitor


def find_secret_solution(name, ip):
	ippart = '0' if ip == 'servicehost' else ip.split('.')[2]
	return 'geocache{' + hashlib.sha256((name + '|' + str(ippart)).encode('utf-8')).hexdigest()[:12] + '}'


typical_post_titles = [
	'Random thought',
	'Message to the world',
	'Hey!',
	'On the secret hunt',
	'Nice place',
	'Getting sunburned',
	'Need a break'
]
typical_posts = [
	'Hey mom, look what I just found!',
	'Wandering around forever. Lovely place!',
	'This is the most beautiful river I\'ve ever seen!',
	'Doesn\'t the sunset look beautiful on my new $5000 phone?',
	'Anyone give hints for the second secret? It\'s too hard for me!!!111!!!',
	'Soo sad I can\'t find the secrets on my own. Why hiding them so hard?',
	'Trading secret 2 hints against secret 1 hints. Anyone?',
	'Just found all the secrets. Feeling great!',
	'Hungry. Want a sandwich. Where to go?',
	'Why doesn\'t the river just flow straigt? That would be more efficient!',
	'Don\'t miss the treetop path. Fantastic view!',
	'Wow, impressive how far you can see from here!',
	'Why so much nature? I HATE NATURE!',
	'Enjoying the sun up here. And you?'
]


def get_status_message(flag=None) -> Tuple[str, str]:
	title = random.choice(typical_post_titles)
	text = random.choice(typical_posts)
	if flag:
		text += ' ' + flag
	return title, text


class SaarschleifeNetServiceInterface(ServiceInterface):
	name = 'saarschleife.net'

	secret_list = [("Cloef", 200), ("Treetop Path", 300), ("Observation Tower", 500)]

	def discover_secrets(self, team: Team, conn: SaarschleifeConnection, count=1000):
		if count == 0:
			return
		total_points = 0
		secrets = list(self.secret_list)
		random.shuffle(secrets)
		for (secret, points) in secrets[0:count]:
			solution = find_secret_solution(secret, team.ip)
			resp = conn.secret_found(secret, solution)
			total_points += points
			assert ' {} points'.format(total_points) in resp, 'Points wrong after discovering secret'
		me = assert_valid_visitor(conn.me())
		assert_equals(total_points, me['points'])

	html_link_pattern = re.compile(r'<link[^>]*\shref="([^"]+)"')
	html_script_pattern = re.compile(r'<script[^>]*\ssrc="([^"]+)"')

	def get_resources(self, text: str):
		resources = self.html_link_pattern.findall(text) + self.html_script_pattern.findall(text)
		map(html.unescape, resources)
		return resources

	def check_integrity(self, team, round_number):
		try:
			conn = SaarschleifeConnection(team.ip)
			resp = conn.get('/')
			assert_requests_response(resp, 'text/html; charset=utf-8')
			if 'ktor' not in resp.headers['Server']:
				GameLogger.log(resp.headers)
				raise MumbleException("Malicious server response")

			# Check if css / js files are being served
			resources = self.get_resources(resp.text)
			interesting_resources = [r for r in resources if 'google' not in r][::-1]
			expected_strings = {"Error loading module 'frontend'"}
			for res in interesting_resources:
				response = conn.get(res)
				print('[{}] {}'.format(response.status_code, res))
				if response.status_code == 200:
					for s in list(expected_strings):
						if s in response.text:
							expected_strings.remove(s)
					if len(expected_strings) == 0:
						break

			if len(expected_strings) > 0:
				print('Not found:')
				for s in expected_strings:
					print('-', s)
				raise MumbleException("Resources (css/js) could not be retrieved")

		except IOError:
			raise OfflineException('Server not reachable')

	def store_flags(self, team, round_number):
		"""
		create 2 visitors, A and B
		- A stores public message
		- B checks readability
		- B solves all, befriends A
		- B stores secret message
		- A checks readability and shares
		- B removes its message / replaces with public
		=> Finally: A stores the flag as non-public status message, with B as author. A and B are friends, both have access.
		:param Team team:
		:param int round_number:
		:return:
		"""
		user_a = usernames.generate_username()
		pass_a = usernames.generate_password()
		user_b = usernames.generate_username()
		pass_b = usernames.generate_password()
		print(f'User A: {user_a} : {pass_a}')
		print(f'User B: {user_b} : {pass_b}')
		self.store(team, round_number, 'credentials', [user_a, pass_a, user_b, pass_b])

		try:
			flag = self.get_flag(team, round_number)
			conn_a = SaarschleifeConnection(team.ip)
			conn_b = SaarschleifeConnection(team.ip)
			late_b = random.randint(1, 100) < 60
			if not late_b:
				conn_b.login(user_b, pass_b)

			# A: store public message
			conn_a.login(user_a, pass_a)
			title, text = get_status_message()
			conn_a.statusmessage_set(title, text, True)

			# B: read public message
			if late_b:
				conn_b.login(user_b, pass_b)
			assert_valid_visitor(conn_b.get_visitor(user_a))
			status = assert_valid_statusmessage(conn_b.statusmessage_get(user_a))
			assert status['text'] == text, 'Could not retrieve public status message'

			# B: solve all challenges, get a friend of A
			self.discover_secrets(team, conn_b)
			conn_b.befriend(user_a)

			# A: solve some challenges
			self.discover_secrets(team, conn_a, random.randint(0, 2))

			# B: store secret message
			title, text = get_status_message(flag)
			conn_b.statusmessage_set(title, text, False)

			# A: read secret message and share
			status = assert_valid_statusmessage(conn_a.statusmessage_get(user_b))
			assert status['title'] == title and status['text'] == text, 'Could not retrieve private status message'
			conn_a.statusmessage_spread(user_b)
			conn_a.logout()

			# B: remove or replace message
			if random.randint(1, 100) <= 35:
				print('(remove original message)')
				conn_b.statusmessage_remove()
			else:
				print('(replace original message)')
				title, text = get_status_message()
				conn_b.statusmessage_set(title, text, True)

			conn_b.logout()

			return 1
		except IOError:
			raise OfflineException('Server not reachable (store)')

	def retrieve_flags(self, team: Team, round_number: int):
		"""
		A stores the flag as non-public status message, with B as author. A and B are friends, both have access.
		:param Team team:
		:param int round_number:
		:return:
		"""
		try:
			user_a, pass_a, user_b, pass_b = self.load(team, round_number, 'credentials')
		except TypeError:
			raise FlagMissingException("Flag has never been stored")

		try:
			conn = SaarschleifeConnection(team.ip)
			retrieve_as_a = random.randint(1, 100) <= 40
			if retrieve_as_a:
				print('Retrieve as A')
				conn.login(user_a, pass_a)
			else:
				print('Retrieve as B')
				conn.login(user_b, pass_b)

			# check if user A is on list
			visitors = conn.list_visitors()
			if user_a not in visitors:
				raise MumbleException('Usernames missing on /visitor')

			# check if status message is present
			if retrieve_as_a:
				visitor = assert_valid_visitor(conn.me())
			else:
				visitor = assert_valid_visitor(conn.get_visitor(user_a))
			assert_equals(True, visitor['statusMessage'])

			# retrieve message
			message = assert_valid_statusmessage(conn.statusmessage_get(user_a))
			assert message['writtenBy'] == user_b, 'Author of retrieved message is wrong'
			flag = self.get_flag(team, round_number)
			if flag not in message['title'] and flag not in message['text']:
				raise FlagMissingException(flag + " not found")

			conn.logout()

			return 1
		except IOError:
			raise OfflineException('Server not reachable (retrieve)')


def main_test():
	# TEST CODE
	team = Team(12, 'n00bs', '127.0.0.1')
	round_number = int(sys.argv[1]) if len(sys.argv) > 1 else random.randint(2, 2000)
	service = SaarschleifeNetServiceInterface(7)

	print('[1] Integrity check...')
	service.check_integrity(team, round_number)
	print('Passed.')

	print('[2] Store flags...')
	flags = service.store_flags(team, round_number)
	print('Done ({} flags).'.format(flags))

	print('[3] Retrieve the flags in the next round')
	flags = service.retrieve_flags(team, round_number)
	print('Done ({} flags).'.format(flags))


if __name__ == '__main__':
	main_test()
