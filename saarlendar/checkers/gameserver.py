import random
import sys
import urllib

from gamelib import *
from saarlendar import Remote, randomstring, now

DEBUG = lambda *args: None
VERBOSE = lambda *args: None

import re

re_username = re.compile("^\w+$")


def generate_username():
    while True:
        user = usernames.generate_username()
        if (re_username.match(user)):
            return user


# DEBUG = print
# VERBOSE = print

def remote(ip, logging_verbose=VERBOSE, logging_debug=DEBUG):
    r = Remote('http://{}:1337/'.format(ip).encode())
    r.DEBUG = logging_debug
    r.VERBOSE = logging_verbose
    return r


class SaarlendarChecker(ServiceInterface):
    name = 'saarlendar'

    def login(self, team, round, register=False, index=None):
        r = remote(team.ip)
        if register:
            username = generate_username()
            password = usernames.generate_password(10, 20)
            r.register(username.encode(), password.encode())
            if index is not None:
                self.store(team, round, 'credentials' + str(index), [username, password])
        else:
            assert index is not None
            try:
                username, password = self.load(team, round, 'credentials' + str(index))
            except TypeError:
                # self.load returned None => no credentials stored
                raise FlagMissingException("no login credentials known")
        r.login(username.encode(), password.encode())
        return r

    def check_integrity(self, team, round):
        try:
            r = self.login(team, round, True)
            r2 = self.login(team, round, True)
        except OSError:
            raise OfflineException('Could not login')
        try:
            r.addevent(title=randomstring(10), content=randomstring(40), public=True)
            r2.addevent(title=randomstring(10), content=randomstring(40), public=False)
        except OSError:
            raise OfflineException('Could not add event')
        try:
            r.getevents(public=True)
            r2.getevents(public=False)
        except OSError:
            raise OfflineException('Could not get events')
        try:
            r.sendmessage(to=r2.username, message=randomstring(10))
            r2.sendmessage(to=r.username, message=randomstring(13))
        except OSError:
            raise OfflineException('Could not send message')
        try:
            nonce = randomstring(10)
            assert nonce in r.shell_plain(("print('%s')" % nonce.decode()).encode())
        except OSError:
            raise OfflineException('Could not use saarjs shell')
        try:
            r2.audit()
        except OSError:
            raise OfflineException('Could not audit nginx.conf')
        if random.randint(0, 10) == 0:
            try:
                # We don't want to run that too often. It takes a few requests to iterate over all config files
                r2.audit_recursively()
            except:
                raise OfflineException('Could not do a full audit')
        else:
            print("skipping full audit...")

    def store_flags(self, team, round):
        try:
            r = self.login(team, round, True, 1)
            r.getevents(public=False)
            r.getevents(public=True)
            r2 = self.login(team, round, True, 2)
            r3 = self.login(team, round, True, 3)
            timestamp = now()
            self.store(team, round, 'timestamp', timestamp.decode())
            r.addevent(title="flag party!", content="I got some spare flags to share!", public=True,
                       timestamp=timestamp)
            r2.shell_sendmessage(r.username, b"Can you invite me to your party? I can bring flags and schwenker too!")
            flag = self.get_flag(team, round, 2)
            r2.shell_sendmessage(r2.username, b"Here is my flag: " + flag.encode())
            flag = self.get_flag(team, round, 1)
            r.shell_addevent(randomstring(10), flag.encode(), public=False, timestamp=timestamp)
            flag = self.get_flag(team, round, 3)
            r3.shell_addevent(randomstring(10), flag.encode(), public=False, timestamp=timestamp)
            return 3
        except OSError:
            raise OfflineException('Could not register')

    def retrieve_flags(self, team, round):
        r = self.login(team, round, False, 1)
        r2 = self.login(team, round, False, 2)
        r3 = self.login(team, round, False, 3)
        timestamp = self.load(team, round, 'timestamp').encode()
        try:
            e1 = r.shell_getevents(shell=Remote.xored_message_shell)[::-1]
            e1p = r.getevents(timestamp=timestamp, public=True)
            m1 = r.getmessages()
            e2 = r2.shell_getevents(shell=Remote.xored_message_shell)[::-1]
            e2p = r2.getevents(timestamp=timestamp, public=True)
            m2 = r2.shell_getmessages(shell=Remote.xored_message_shell)
            e3 = r3.shell_getevents(shell=Remote.xored_message_shell)
            m3 = r3.getmessages()
            for event in e1:
                try:
                    e = json.loads(event)
                    flag = urllib.parse.unquote(e["text"])
                    _, _, _, payload = self.check_flag(flag, team.id, round)  # returns None,None,None,None if invalid
                    if not payload or payload != 1:
                        pass
                    else:
                        print("found flag 1", flag)
                        break
                except:
                    pass
            else:
                raise FlagMissingException("Flag1 not found!")
            for event in e1p:
                try:
                    e = json.loads(event)
                    if (e["creator"] == r.username.decode() and e["public"] == 1):
                        break
                except:
                    pass
            else:
                raise MumbleException("public event not found for user1")
            for event in e3:
                try:
                    e = json.loads(event)
                    flag = urllib.parse.unquote(e["text"])
                    _, _, _, payload = self.check_flag(flag, team.id, round)  # returns None,None,None,None if invalid
                    if not payload or payload != 3:
                        pass
                    else:
                        print("found flag 3", flag)
                        break
                except:
                    pass
            else:
                raise FlagMissingException("Flag3 not found")
            flag2 = False
            inviterequest = False
            inviterequest2 = False
            for m in m2:
                try:
                    m = json.loads(m)
                    text = urllib.parse.unquote_plus(m["text"])
                    if (text.startswith("Here is my flag: ")):
                        _, _, _, payload = self.check_flag(text[len("Here is my flag: "):], team.id, round)
                        if not payload or payload != 2:
                            pass
                        else:
                            print("found flag 2", text[len("Here is my flag: "):])
                            flag2 = True
                            break
                    elif (text.startswith("Can you invite me") and m["to"] == r.username.decode()):
                        inviterequest2 = True
                except:
                    pass

            for m in m1:
                try:
                    m = json.loads(m)
                    text = urllib.parse.unquote_plus(m["text"])
                    if (text.startswith("Can you invite me") and m["from"] == r2.username.decode()):
                        inviterequest = True
                except:
                    continue
            if not flag2:
                raise FlagMissingException("Flag2 not found")
            if not inviterequest:
                raise MumbleException("Invitation not arrived")
            if not inviterequest2:
                raise MumbleException("Invitation not sent")
        except OSError:
            raise OfflineException('Could not connect')
        return 3


if __name__ == '__main__':
    # TEST CODE
    team = Team(12, 'n00bs', '127.0.0.1')
    round = int(sys.argv[1]) if len(sys.argv) > 1 else 2
    print("Round:", round)
    service = SaarlendarChecker(7)

    print('[1] Integrity check...')
    service.check_integrity(team, round)
    print('Passed.')

    print('[2] Store flags...')
    flags = service.store_flags(team, round)
    print('Done ({} flags).'.format(flags))

    print('[3] Retrieve the flags in the next round')
    flags = service.retrieve_flags(team, round)
    print('Done ({} flags).'.format(flags))
