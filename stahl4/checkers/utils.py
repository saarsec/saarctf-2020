import json
from base64 import b64decode, b64encode
from datetime import datetime, timedelta
from enum import Enum

import pytz
from Crypto.Hash import SHA512  # pycryptodome
from Crypto.PublicKey import RSA
from Crypto.Signature import pkcs1_15
from Crypto.Util import asn1
from pwn import *


class ServerType(Enum):
    LOGISTIC = 1
    PRODUCTION = 2


def is_ip(addr):
    return socket.gethostbyname(addr) == addr


def random_string(str_length):
    letters = string.ascii_lowercase
    return ''.join(random.choice(letters) for _ in range(str_length))


def getRSAKey():
    key = RSA.generate(2048)
    return key


def generate_key_and_store_on_disk():
    gameserver_key = getRSAKey()
    with open("gameserver.key", "wb+") as fd:
        fd.write(gameserver_key.export_key('PEM'))
    return gameserver_key


def import_key():
    with open("gameserver.key", "rb+") as fd:
        gameserver_key = RSA.import_key(fd.read())
    return gameserver_key


def send_heartbeat(connection):
    connection.send(b"\xF0")
    pong = connection.recv(1, timeout=5)
    assert pong == b"\xF0"


def connectAndHandshake(ip, key, verbose=True, timeout=5):
    # timeout = 5
    SERVICE_PORT = 21485
    old_log_level = context.log_level
    if verbose:
        context.log_level = 'debug'
    else:
        context.log_level = 'error'

    binary_key = asn1.DerSequence([key.n, key.e]).encode()
    our_id = hashlib.sha256(binary_key).hexdigest()

    handshake = dict()
    pk = key.publickey().export_key(format="DER")
    c = connect(ip, SERVICE_PORT, timeout=timeout)
    nonce = c.recv(1024, timeout=timeout)

    nonce = json.loads(nonce)["nonce"].replace('"', '')
    handshake["nonce"] = nonce
    if is_ip(c.rhost):
        handshake["target"] = c.rhost
    else:
        handshake["target"] = socket.gethostbyname(c.rhost)
    handshake["pub"] = b64encode(pk).replace(b'+', b'-').replace(b'/', b'_').decode()
    t = datetime.now(pytz.timezone('Europe/Berlin')) + timedelta(minutes=1)
    handshake["expiry"] = t.isoformat('T')
    nonce = b64decode(nonce.replace('-', '+').replace('_', '/'))
    '''
        Golang notation             Java notation        C notation    Standard
    2006-01-02T15:04:05-0700     yyyy-MM-dd'T'HH:mm:ssZ        %FT%T%z        ISO 8601
    '''
    tstamp = t.strftime("%FT%T%z")
    if len(tstamp) < 24:
        tstamp = tstamp.encode() + b'\x00' * (24 - len(tstamp))
    else:
        tstamp = tstamp.encode()

    h = SHA512.new(nonce + handshake["target"].encode() + tstamp)
    handshake["hmac"] = b64encode(pkcs1_15.new(key).sign(h)).replace(b'+', b'-').replace(b'/', b'_').decode()
    c.send(json.dumps(handshake))
    event = c.recvuntil(b"}", timeout=timeout)

    server_type = ServerType(int(json.loads(event)["type"]))
    log.debug("server type %s" % server_type)
    context.log_level = old_log_level

    return c, our_id, server_type
