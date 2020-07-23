import html
import random
import subprocess
import sys
import tempfile

import requests

from gamelib import *


def generate_title(words=2, camelcase=True, alphanum=True, generator=random, sep=' '):
    components = []
    for i in range(words - 1):
        components.append(generator.choice(usernames.USERNAME_ADJECTIVES))
    components.append(generator.choice(usernames.USERNAME_NOUNS))
    if camelcase:
        components = map(usernames.ucfirst, components)
    name = sep.join(components)
    if alphanum:
        name = name.replace('-', '')
    return name


class Unauthorized(ValueError):
    def __init__(self):
        super().__init__('Not authorized')


def get_token(r):
    return re.findall(r'__RequestVerificationToken.*?value="(.*?)" />', r.text, re.S)[0]


class SaarXivAPI:
    def __init__(self, URL, timeout=5):
        self.URL = URL
        self.timeout = timeout

    def register_user(self, username, firstname, lastname, password):
        s = requests.Session()
        s.verify = False
        r1 = s.get(self.URL + '/Account/Register', timeout=self.timeout)
        assert_requests_response(r1, 'text/html; charset=utf-8')
        token = get_token(r1)
        r2 = s.post(self.URL + '/Account/Register',
                    data={'Input.Username': username, 'Input.Firstname': firstname, 'Input.Lastname': lastname,
                          'Input.Password': password, 'Input.ConfirmPassword': password,
                          '__RequestVerificationToken': token})
        assert_requests_response(r2, 'text/html; charset=utf-8')
        return SaarXivLoggedInSession(self.URL, s, timeout=self.timeout)

    def login_user(self, username, password):
        s = requests.Session()
        s.verify = False
        r1 = s.get(self.URL + '/Account/Login', timeout=self.timeout)
        assert_requests_response(r1, 'text/html; charset=utf-8')
        token = get_token(r1)
        r2 = s.post(self.URL + '/Account/Login',
                    data={'Input.Username': username, 'Input.Password': password, '__RequestVerificationToken': token},
                    timeout=self.timeout)
        assert_requests_response(r2, 'text/html; charset=utf-8')
        return SaarXivLoggedInSession(self.URL, s, timeout=self.timeout)

    def anonymous_session(self):
        s = requests.Session()
        return SaarXivAnonymousSession(self.URL, s, timeout=self.timeout)


class SaarXivAnonymousSession:
    def __init__(self, URL, s, timeout):
        self.URL = URL
        self.s = s
        self.timeout = timeout

    def get_content_via_pdf(self, paper_id):
        r1 = self.s.get(self.URL + '/Paper/Download', params={'id': paper_id}, timeout=self.timeout)
        if 'Login' in r1.url or 'AccessDenied' in r1.url or r1.headers['Content-Type'] != 'application/octet-stream':
            raise Unauthorized()
        assert_requests_response(r1, 'application/octet-stream')
        with tempfile.NamedTemporaryFile() as f:
            f.file.write(r1.content)
            f.file.flush()
            return subprocess.run(['pdftotext', f.name, '-'], stdout=subprocess.PIPE).stdout.decode()

    def get_content_via_pdf_with_token(self, paper_id, download_token):
        r1 = self.s.get(self.URL + '/Paper/Download', params={'id': paper_id}, timeout=self.timeout)
        if 'Login' in r1.url or 'AccessDenied' in r1.url:
            raise Unauthorized()
        assert_requests_response(r1, 'text/html; charset=utf-8')
        token = get_token(r1)
        r2 = self.s.post(self.URL + '/Paper/Download', params={'id': paper_id},
                         data={'DownloadToken': download_token, '__RequestVerificationToken': token},
                         timeout=self.timeout)
        assert_requests_response(r2, 'application/octet-stream')
        with tempfile.NamedTemporaryFile() as f:
            f.file.write(r2.content)
            f.file.flush()
            return subprocess.run(['pdftotext', f.name, '-'], stdout=subprocess.PIPE).stdout.decode()


class SaarXivLoggedInSession(SaarXivAnonymousSession):
    def create_paper(self, title, content, under_submission=False):
        r1 = self.s.get(self.URL + '/Paper/Create', timeout=self.timeout)
        assert_requests_response(r1, 'text/html; charset=utf-8')
        token = get_token(r1)
        r2 = self.s.post(self.URL + '/Paper/Create', data={'Input.Title': title, 'Input.Content': content,
                                                           'Input.UnderSubmission': under_submission,
                                                           '__RequestVerificationToken': token}, timeout=self.timeout)
        assert_requests_response(r2, 'text/html; charset=utf-8')
        return int(re.findall(r'Successfully created paper (.*?)</div>', r2.text, re.S)[0])

    def get_content_via_edit(self, paper_id):
        r1 = self.s.get(self.URL + '/Paper/Edit', params={'id': paper_id}, timeout=self.timeout)
        assert_requests_response(r1, 'text/html; charset=utf-8')
        if 'Login' in r1.url or 'AccessDenied' in r1.url:
            raise Unauthorized()
        return html.unescape(re.findall(r'<textarea.*?>\n(.*?)</textarea>', r1.text, re.S)[0])

    def get_download_token(self, paper_id):
        r1 = self.s.get(self.URL + '/Paper/Share', params={'id': paper_id}, timeout=self.timeout)
        assert_requests_response(r1, 'text/html; charset=utf-8')
        if 'Login' in r1.url or 'AccessDenied' in r1.url or 'Error: You are not authorized to view this file.' in r1.text:
            raise Unauthorized()
        return html.unescape(re.findall(r'<textarea.*?>\n(.*?)</textarea>', r1.text, re.S)[0])


class SaarXivInterface(ServiceInterface):
    name = 'SaarXiv'

    def check_integrity(self, team, round):
        try:
            assert_requests_response(requests.get('http://{}:5000/'.format(team.ip), timeout=gamelib.TIMEOUT),
                                     'text/html; charset=utf-8')
        except (IOError, ConnectionRefusedError):
            raise OfflineException('Could not load home page')

    def store_flags(self, team, round):

        api = SaarXivAPI('http://{}:5000'.format(team.ip), timeout=gamelib.TIMEOUT)

        username = usernames.generate_username()
        firstname = usernames.generate_name(words=1)
        lastname = usernames.generate_name(words=1)
        password = usernames.generate_password()
        self.store(team, round, 'credentials', [username, password])

        try:
            session = api.register_user(username, firstname, lastname, password)
        except IOError:
            raise OfflineException('Could not register')

        try:
            flag = self.get_flag(team, round, 1)
            title = generate_title(words=3)
            content = '''Here is the truth about {0}:
            \\begin{{verbatim}}
            flag = {1}
            \\end{{verbatim}}
            '''.format(title.lower(), flag)
            paper_id = session.create_paper(title, content, under_submission=True)
            self.store(team, round, 'paper_id', [paper_id])

            return 1
        except IOError:
            raise OfflineException('Could not create paper')

    def retrieve_flags(self, team, round):
        api = SaarXivAPI('http://{}:5000'.format(team.ip), timeout=gamelib.TIMEOUT)
        up = self.load(team, round, 'credentials')
        if up is None:
            raise FlagMissingException('Missing username/password')
        username, password = up
        paper_id = self.load(team, round, 'paper_id')
        if paper_id is None:
            raise FlagMissingException('Missing paper_id')
        try:
            session = api.login_user(username, password)

            c1 = session.get_content_via_edit(paper_id)
            # 1 - check via edit functionality
            flag = self.get_flag(team, round, 1)
            if flag not in c1:
                raise FlagMissingException("Flag not found")

            # 2 - check via pdf download
            c2 = session.get_content_via_pdf(paper_id)
            if flag not in c2:
                raise FlagMissingException("Flag not found")

            # 3 - check via sharing functionality
            download_token = session.get_download_token(paper_id)
            anonymous_session = api.anonymous_session()
            c3 = anonymous_session.get_content_via_pdf_with_token(paper_id, download_token)
            if flag not in c3:
                raise FlagMissingException("Flag not found")

            return 1
        except IOError:
            raise OfflineException('Could not login')


if __name__ == '__main__':
    # TEST CODE
    team = Team(12, 'n00bs', '127.0.0.1')
    round = int(sys.argv[1]) if len(sys.argv) > 1 else 2
    service = SaarXivInterface(7)

    print('[1] Integrity check...')
    service.check_integrity(team, round)
    print('Passed.')

    print('[2] Store flags...')
    flags = service.store_flags(team, round)
    print('Done ({} flags).'.format(flags))

    print('[3] Retrieve the flags in the next round')
    flags = service.retrieve_flags(team, round)
    print('Done ({} flags).'.format(flags))
