import sys

from gamelib import *

from utils import *

def store_flag_errorlog(target, flag, gameserver_key, timeout=5, verbose=False):
    conn, gameserver_id, _ = connectAndHandshake(target, gameserver_key, timeout=timeout)
    if verbose:
        print("Finished handshake. Our ID: %s" % gameserver_id)
    message = dict()
    message['message'] = flag
    message['id'] = gameserver_id
    message_json = json.dumps(message)
    conn.send(message_json)
    conn.close()


def store_flag_datalog(target, flag, gameserver_key, timeout=5, verbose=False):
    conn, gameserver_id, _ = connectAndHandshake(target, gameserver_key, timeout=timeout)
    if verbose:
        print("Finished handshake. Our ID: %s" % gameserver_id)
    message = dict()
    message['message'] = flag
    message['id'] = gameserver_id
    message['amount'] = str(random.randint(3, 20))
    message['type'] = "GetMaterialsCommand"
    message_json = json.dumps(message)
    conn.send(message_json)
    conn.close()


def retrieve_errorlog(target, gameserver_key, timeout=5):
    conn, gameserver_id, _ = connectAndHandshake(target, gameserver_key, timeout=timeout)

    message = dict()
    message['type'] = 'GetErrorCommand'
    message['asked_id'] = gameserver_id
    message_json = json.dumps(message)
    conn.send(message_json)
    resp = conn.recvuntil("}", timeout=timeout)
    resp = json.loads(resp)
    assert resp["type"] == "GetEvent"
    key, content = resp["message"].split("-")
    assert len(key) == len(content)
    content = base64.b64decode(content)
    key = base64.b64decode(key)
    content_decoded = bytes([k ^ c for k, c in zip(key, content)])
    content_decoded = content_decoded.decode()
    conn.close()
    return content_decoded


def retrieve_datalog(target, gameserver_key, timeout=5):
    conn, gameserver_id, _ = connectAndHandshake(target, gameserver_key, timeout=timeout)

    message = dict()
    message['type'] = 'GetDataCommand'
    message['asked_id'] = gameserver_id
    message_json = json.dumps(message)
    conn.send(message_json)
    resp = conn.recvuntil("}", timeout=timeout)
    resp = json.loads(resp)
    assert resp["type"] == "GetEvent"
    key, content = resp["message"].split("-")
    assert len(key) == len(content)
    content = base64.b64decode(content)
    key = base64.b64decode(key)
    content_decoded = bytes([k ^ c for k, c in zip(key, content)])
    content_decoded = content_decoded.decode()
    conn.close()
    return content_decoded


def send_GetMaterials(target, gameserver_key, timeout=5):
    """
    returns "NOT-LOGISTIC" in case we are talking to a production server
    """
    conn, gameserver_id, server_type = connectAndHandshake(target, gameserver_key)
    if server_type != ServerType.LOGISTIC and server_type is not None:
        # this only makes sense for logistic servers
        return server_type, "NOT-LOGISTIC"
    message = dict()
    message['type'] = 'GetMaterialsCommand'
    message['asked_id'] = gameserver_id
    message['message'] = random_string(random.randint(8, 16))
    # get random amount of materials
    message['amount'] = str(random.randint(3, 20))

    message_json = json.dumps(message)
    conn.send(message_json)
    resp = conn.recvuntil("}", timeout=timeout)
    resp = json.loads(resp)
    conn.close()
    return server_type, resp


def send_GetProducts(target, gameserver_key, timeout=5):
    """
    returns "NOT-LOGISTIC" in case we are talking to a production server
    """
    conn, gameserver_id, server_type = connectAndHandshake(target, gameserver_key)
    if server_type != ServerType.LOGISTIC and server_type is not None:
        # this only makes sense for logistic servers
        return server_type, "NOT-LOGISTIC"
    message = dict()
    message['type'] = 'GetProductsCommand'
    message['message'] = random_string(random.randint(8, 16))
    # get random amount of materials
    message['amount'] = str(random.randint(3, 20))

    message_json = json.dumps(message)
    conn.send(message_json)
    resp = conn.recvuntil("}", timeout=timeout)
    resp = json.loads(resp)
    conn.close()
    return server_type, resp


class Stahl4Interface(ServiceInterface):
    # service name
    name = 'Stahl4'
    # connection timeout in seconds
    timeout = TIMEOUT
    # disable pwntools logging
    context.log_level = 'error'

    gameserver_key = '''-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAps3Y/k5VLmIzhkzW2O7EnL+SEWVJG/FOzVkZTQG2vaABQfnd
EwqDIbphsO9Lu6AOh+Y+rFl1W4Vk0KXH/Clex078MGFGMOPoQqTV8d7EwyhLad0s
dEIHalDSwVs2dSvPhgIM6F2Mo2Rj9xk/h3tDT6xVU9ZF9hh58IG9+ICTSCJAxHP6
6RFMUkp8bepD6nUyVb10GRnty9Cv1oEKLk4hk2ADtAS+1UCiWZlzCoFTUx9WcUai
iUfs7o0tWH8u7arfP4tR8DArSCTR8O2aKQfZNjEkpa93FadG5ptCt8J1HtBEg/Xg
p5UmGdp6m7C/hqp12OFmWUV1HeqCqtmDDJ1OzQIDAQABAoIBAB2xMtlB9GPHU0iq
0OvYxOLqLwQttW/l6pMfMyPEFFV48ABgi/vhuSn91Xn7fs2tVijW0X04h219N070
076NVrO/2aB7EFWPnD6QvQzLDNd4isQTfOBNCKjiLkIBDJaa79v47XdHf8tuCWVo
9ygUxwydrqq3z8hY3wvy7HAZ4x3JwXK9ZWRGB1uW0Mw4xlgQgztGMLx+2GI2fqGO
uW4Zn3umRMfDpNWKTvSamoJbhJqlypdUbcMpbkiQVZjgQzXgXMmsW0BmVsCx3Lki
zHDs8UGo0ZPhBE54ah+F1jsrtT7OPxGBpsYOL1efy4xEeep6GEUtOkXsr/KUWsxm
B25dB/kCgYEAw2/MEqmYhqI3+HCVBOFGQJ89Slzx+nNUrG94BHbGl/GU8WIhkTYE
9Qla0SgcwieKtmk0PeowFTp1/srasRSwkJaTOYjYojlc9BfgoDm67rAO1k7OY8XH
pyh44QqaBxE1fJCrcMaCP3YpYfwF2u629vMsm0gjHsdgZIv3uh5wplMCgYEA2n6Y
MT3F2RhCWGLEn0Oq5mT4EL9ee/2HXKBHldPm6ncjrCW5/3gxAFR3tN3qZ7E2JoY6
ZSg3t7WGepXD020yme10Gmfr5Wt1i4kYyqEa7GXf4ECdZU50CnrxO216OnkkVFqs
9b6k+WJF+7MygOcJmVWAI7IjkaWp0ZAjqjMzUl8CgYA9CYl+uqD8S2WXnfGsmH7M
Z+9IbkbUpXhoclfVbRMRGU4fJTq+k1RxAv7uG8z/hkH4PlsmiyGBP8TiUpCChaev
QJeyfF7MK4YwZdSttdn/+dRocixbVMXKGwXFov3//wvpX3Vrv1OmZkz+YSui+LMD
5WalCJ6PWk1smZpA8ojUKQKBgEW0lYFAH9p+rsvb1raos+EE3U8afl44J/MY/z2B
eO3cTHkjIA+snJVqXTZKhfnGw2vO7tpO1le5hcmd9feBot8QrjWuacerXLjDaDFc
7GX2qlG0y4ICYWrmhgdbid8Vvs1akEtmIuOcwo7mQHp3Osy8RkEdF9PjciX1QiuO
YhUpAoGAOGbFRv1/tyB0U/6mekBy8ccWQ7dldGvENKjOKtKLHT1VnTHwFp3986FI
czQlwezGGENnGlcXRctgffY6cEJpgFAq5p4gnPcGvfqQBrxF4jpX1iXY7CkfOpE1
yOh7HNOBFi4ApS4bATKGQzZK5F0y0Z+TpvDHC9YHhJcJv6Tq2So=
-----END RSA PRIVATE KEY-----'''

    # This is necessary for local testing
    def __init__(self, id):
        self.gameserver_key = RSA.import_key(self.gameserver_key)
        super(Stahl4Interface, self).__init__(id)

    def check_integrity(self, team, round, verbose=False):
        # for integrity check we use randomized gameserver keys to prevent fingerprinting
        gameserver_key = getRSAKey()
        try:
            # check GetMaterials endpoint
            server_type, resp = send_GetMaterials(team.ip, gameserver_key, self.timeout)
            # note that server type can not be None as this is the first connection to the server (with this key)
            assert server_type is not None
            if resp != "NOT-LOGISTIC" and resp['type'] != "DeliveredMaterialsEvent" \
                    and resp['type'] != "NoMaterialsEvent":
                raise MumbleException(f"GetMaterials returned {resp}")

            if verbose:
                if resp == "NOT-LOGISTIC":
                    print("Server is not a logistic server")
                elif resp == "NoMaterialsEvent":
                    print("Server got not enough materials")
                else:
                    print(f"Retrieved {resp['amount']} materials")

            # check GetProducts endpoint
            if server_type == ServerType.LOGISTIC:
                # it is important that we only check this for logistic servers as we can only
                # catch the wrong server-type in the FIRST connection NOT IN THE SECOND
                # (due to refactoring this is not strictly needed anymore, one could refactor this)
                server_type2, resp = send_GetProducts(team.ip, gameserver_key, self.timeout)
                if resp != "NOT-LOGISTIC" and resp['type'] != "FetchedProductsEvent" \
                        and resp['type'] != "NoProductsEvent":
                    raise MumbleException(f"GetProducts returned {resp}")
                if verbose:
                    if resp == "NoProductsEvent":
                        print("Server got not enough products")
                    else:
                        print(f"Retrieved {resp['amount']} products")

            # store and retrieve errorlog
            value_to_store = random_string(random.randint(8, 16))
            store_flag_errorlog(team.ip, value_to_store, gameserver_key, self.timeout)
            resp = retrieve_errorlog(team.ip, gameserver_key, self.timeout)
            if value_to_store not in resp:
                raise MumbleException(f"Failed errorlog check: got response {resp} - missing: {value_to_store}")

            # store and retrieve datalog
            value_to_store = random_string(random.randint(8, 16))
            store_flag_datalog(team.ip, value_to_store, gameserver_key, self.timeout)
            resp = retrieve_datalog(team.ip, gameserver_key, self.timeout)
            if value_to_store not in resp:
                raise MumbleException(f"Failed datalog check: got response {resp} - missing: {value_to_store}")

        except (IOError, EOFError, pwnlib.exception.PwnlibException):
            raise OfflineException('Could not check integrity')
        except (json.JSONDecodeError, binascii.Error):
            raise MumbleException('Could not check integrity')

    def store_flags(self, team, round):
        gameserver_key = self.gameserver_key
        flag1 = self.get_flag(team, round, 1)
        flag2 = self.get_flag(team, round, 2)
        try:
            # flag store 1
            store_flag_errorlog(team.ip, flag1, gameserver_key, self.timeout)
            # flag store 2
            store_flag_datalog(team.ip, flag2, gameserver_key, self.timeout)
        except (IOError, EOFError, pwnlib.exception.PwnlibException):
            raise OfflineException('Could not store flag')
        except (json.JSONDecodeError, binascii.Error):
            raise MumbleException('Could not store flag')
        return 2

    def retrieve_flags(self, team, round):
        gameserver_key = self.gameserver_key
        flag1 = self.get_flag(team, round, 1)
        flag2 = self.get_flag(team, round, 2)

        try:
            # flag store 1
            errorlog_content = retrieve_errorlog(team.ip, gameserver_key, self.timeout)
            for flag in self.search_flags(errorlog_content):
                if flag == flag1:
                    break
            else:
                raise FlagMissingException("Flag not found in errorlog")

            # flag store 2
            datalog_content = retrieve_datalog(team.ip, gameserver_key, self.timeout)
            for flag in self.search_flags(datalog_content):
                if flag == flag2:
                    break
            else:
                raise FlagMissingException("Flag not found in datalog")
            return 2
        except (IOError, EOFError, pwnlib.exception.PwnlibException):
            raise OfflineException('Could not connect to peer')
        except (json.JSONDecodeError, binascii.Error):
            raise MumbleException('Could not connect to peer')


if __name__ == '__main__':
    # TEST CODE
    team = Team(12, 'n00bs', '10.32.1.2')
    round = 2
    service = Stahl4Interface(0x1337)
    print('[1] Integrity check...')
    service.check_integrity(team, round, verbose=False)
    print('Passed.')

    print('[2] Store flags...')
    flags = service.store_flags(team, round)
    print('Done ({} flags).'.format(flags))

    print('[3] Retrieve the flags in the next round')
    flags = service.retrieve_flags(team, round)
    print('Done ({} flags).'.format(flags))

    print('[4] Integrity check...')
    service.check_integrity(team, round, verbose=False)
    print('Passed.')
