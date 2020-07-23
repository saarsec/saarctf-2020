from nacl.public import PrivateKey, SealedBox
import nacl.signing
import nacl.encoding
import nacl
from bs4 import BeautifulSoup
from gamelib import *

import traceback
import requests
import string
import random
import time
import sys
import os
import subprocess

BASE_URI = 'http://{}:20080/mensaar/'
hardcoded_sk = b'\xc5\xb2Bn\xefLN?\xe74V2\xc9im\xa3\x0e\xa3\xab\xbc\xf3I\xf7N7\xcd\x10\x8b.\x9e\xc1\x13'
sign_key = nacl.signing.SigningKey(seed=hardcoded_sk)


class MensaServiceInterface(ServiceInterface):
    name = 'Mensa'

    def check_integrity(self, team, round_number):
        try:
            target = BASE_URI.format(team.ip)
            r = requests.get(target, timeout=TIMEOUT)
            assert_requests_response(r, contenttype='text/html; charset=utf-8')
        except IOError:
            raise OfflineException('Could not login')
        except Exception:
            traceback.print_exc()
            raise MumbleException('Could not retrieve index page')

        try:
            realpath=os.path.realpath(__file__).split("/")[:-1]
            workdir="/".join(realpath)
            proc = subprocess.Popen("php test_insert.php", shell=True, stdout=subprocess.PIPE,cwd=workdir)
            php_obj = str(proc.stdout.read()[:-1],"utf-8") # the serialized obj, static atm, TODO: make dynamic
            token = sign_key.sign(bytes(php_obj,"utf-8")) # signature of our serialized obj
            token_hash = hashlib.sha1(php_obj.encode()).hexdigest()
            data = {
                'data': php_obj,
                'token': token,
                'hash': token_hash
            }
            r = requests.post(target + "next_menu.php", data=data, timeout=TIMEOUT)
            assert_requests_response(r, contenttype='text/html; charset=utf-8')
        except IOError:
            raise OfflineException('Chef could not set new menu')
        except Exception:
            traceback.print_exc()
            raise MumbleException('Chef could not reach next menu endpoint')

        try:
            r = requests.post(target + "do_reservation.php", data={'reset': 1}, timeout=TIMEOUT)
            assert_requests_response(r, contenttype='text/html; charset=utf-8')
        except IOError:
            raise OfflineException('No new reservation could be done')
        except Exception:
            traceback.print_exc()
            raise MumbleException('Something went wrong with the new reservation')

    def store_flags(self, team, round_number):
        tld = ''.join(random.choices(string.ascii_uppercase + string.digits, k=5))
        email = "{}@{}.com".format(usernames.generate_username(), tld)
        pwd = usernames.generate_password()
        self.store(team, round_number, 'credentials', [email, pwd])
        try:
            flag1 = self.get_flag(team, round_number, 1)
            flag2 = self.get_flag(team, round_number, 2)
            if None in [flag1, flag2]:
                raise MumbleException('Could not generate new flags')
            target = BASE_URI.format(team.ip)
            s = requests.Session()
            # Register new account
            data = {
                "name": flag1,
                "email": email,
                "pwd": pwd,
            }
            r = s.post(target + "registration.php", data=data, timeout=TIMEOUT)
            assert_requests_response(r, contenttype='text/html; charset=utf-8')
            # Reserve a seat such that email is semi-public
            r = s.get(target + "reservation.php", timeout=TIMEOUT)
            assert_requests_response(r, contenttype='text/html; charset=utf-8')
            soup = BeautifulSoup(r.text, 'html.parser')

            def finder(tag):
                return tag.name == 'a' and tag.has_attr('data-trigger') and not tag.has_attr('data-content')

            seat = soup.find(finder)
            if seat is not None:
                seat = seat.get('title')[5:].split('.')
            else:
                seat = [0, 0, 0, 0]

            data = {
                "user": email,
                "row": seat[0],
                "table": seat[1],
                "seatrow": seat[2],
                "seatnumber": seat[3]
            }
            r = s.post(target + "do_reservation.php", data=data, timeout=TIMEOUT)
            assert_requests_response(r, contenttype='text/html; charset=utf-8')
            # Place feedback with FLAG No.2 as comment
            data = {
                "food": random.randint(1, 15),
                "rating": random.randint(1, 5),
                "comment": flag2
            }
            r = s.post(target + "feedback.php", data=data, timeout=TIMEOUT)
            assert_requests_response(r, contenttype='text/html; charset=utf-8')
            # Logout for whatever reasons
            r = s.get(target + "logout.php", timeout=TIMEOUT)
            assert_requests_response(r, contenttype='text/html; charset=utf-8')
            return 1
        except IOError:
            raise OfflineException('Could not register')
        except (nacl.exceptions.CryptoError, ValueError):
            raise MumbleException('Crypto operation failed!')
        except Exception as e:
            print(e)
            traceback.print_exc()
            raise MumbleException('Unknown Error during flag storage!')

    def retrieve_flags(self, team, round_number):
        credentials = self.load(team, round_number, 'credentials')
        if credentials is None:
            raise FlagMissingException('Could not restore stored credentials')
        email, pwd = credentials
        try:
            s = requests.Session()
            target = BASE_URI.format(team.ip)
            # Login using the credentials form the last round_number
            data = {
                "mail": email,
                "pwd": pwd,
            }
            r = s.post(target + "login.php", data=data, timeout=TIMEOUT)
            assert_requests_response(r, contenttype='text/html; charset=utf-8')

            #find all reserved objects 
            r = s.get(target + "reservation.php", timeout=TIMEOUT)
            assert_requests_response(r, contenttype='text/html; charset=utf-8')
            soup = BeautifulSoup(r.text, 'html.parser')
            def finder(tag):
                return tag.name == 'a' and tag.has_attr('data-trigger') and tag.has_attr('data-content')

            seats = soup.find_all(finder)
            if seats == []:
                raise MumbleException('No reservations available')
            found=False
            for seat in seats:
                if email in str(seat):
                    found=True
            if not found:
                raise MumbleException('Reservation from last tick disappeared')
            # Visit the main page to get profile information
            r = s.get(target + "index.php", timeout=TIMEOUT)
            assert_requests_response(r, contenttype='text/html; charset=utf-8')
            flag = self.search_flags(r.text)
            if len(flag) == 0:
                raise FlagMissingException("No Flag No. 1 found in Mensa")
            flag = flag.pop()
            _, _, _, payload = self.check_flag(flag, team.id, round_number)
            if not flag or not payload or payload != 1:
                raise FlagMissingException("Invalid Flag No. 1 in Mensa")
            # Validate Feedback comment is still there
            feedback_id = hashlib.sha3_256(email.encode('utf-8')).hexdigest()
            r = requests.get(target + 'get_feedback.php?id=' + feedback_id, timeout=TIMEOUT)
            assert_requests_response(r, contenttype='text/html; charset=utf-8')
            sk = PrivateKey(hardcoded_sk)
            unseal_box = SealedBox(sk)
            # decrypt the received message
            plaintext = unseal_box.decrypt(r.content).decode('utf-8')
            flag = self.search_flags(plaintext)
            if len(flag) == 0:
                raise FlagMissingException("No Flag No. 2 found in Mensa")
            flag = flag.pop()
            _, _, _, payload = self.check_flag(flag, team.id, round_number)
            if not payload or payload != 2:
                raise FlagMissingException("Invalid Flag No. 2 in Mensa")
            return 1
        except FlagMissingException:
            raise
        except MumbleException:
            raise
        except (nacl.exceptions.CryptoError, ValueError):
            raise MumbleException('Crypto operation failed!')
        except IOError:
            raise OfflineException('Could not login')
        except Exception as e:
            print(e)
            traceback.print_exc()
            raise MumbleException('Unknown Error during flag retrieve!')


if __name__ == '__main__':
    # TEST CODE
    test_team = Team(12, 'n00bs', 'localhost')
    test_round = int(sys.argv[1]) if len(sys.argv) > 1 else 2
    service = MensaServiceInterface(0)

    print('[1] Integrity check...')
    service.check_integrity(test_team, test_round)
    print('Passed.')

    print('[2] Store flags...')
    flags = service.store_flags(test_team, test_round)

    print('[3] Retrieve the flags in the next round')
    flags = service.retrieve_flags(test_team, test_round)
    print('Done ({} flags).'.format(flags))
