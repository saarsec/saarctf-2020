import socket
import json
from binascii import unhexlify, hexlify

from gamelib import TIMEOUT, MumbleException, OfflineException, FlagMissingException

DEFAULT_TIMEOUT = TIMEOUT
SHA512_HEX_LEN = 128


class SimpleSocket(object):
    def __init__(self, socket):
        self.socket = socket
        self.rbuf = ""

    def recvline(self):
        while not "\n" in self.rbuf:
            r = self.socket.recv(4096)
            if not r:
                self.socket.shutdown(socket.SHUT_RDWR)
                self.socket.close()
                print("xxx")
            self.rbuf += r.decode("utf-8")
        i = self.rbuf.index("\n")
        line = self.rbuf[:i]
        self.rbuf = self.rbuf[i + 1:]
        return line.strip()

    def sendline(self, msg: str):
        if len(msg) == 0 or msg[-1] != "\n":
            msg += "\n"
        self.socket.sendall(msg.encode("utf-8"))

    def close(self):
        self.socket.shutdown(socket.SHUT_RDWR)
        self.socket.close()

    @staticmethod
    def connect(host, port, timeout=DEFAULT_TIMEOUT):
        socket.setdefaulttimeout(timeout)
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.connect((host, port))
        return SimpleSocket(s)


class GameServer():
    @staticmethod
    def sendCommand(host: str, port: int, cmd: str, sig: str, timeout=DEFAULT_TIMEOUT):
        toException = False
        try:
            s = SimpleSocket.connect(host, port, timeout=timeout)

            def expectLine(line):
                l = s.recvline()
                if l != line:
                    raise MumbleException(f"""Unexpected behaviour: got "{l}", expected "{line}".""")

            expectLine("command:")
            s.sendline(cmd)
            expectLine("signature:")
            s.sendline(sig)
            res = s.recvline()
            s.close()
            return res
        except (socket.timeout, ConnectionRefusedError):
            toException = True
            try:
                s.close()
            except:
                pass
        if toException:
            raise OfflineException(f"Team {host} seems to be offline")

    @staticmethod
    def retrieveCipher(host: str, port: int, flagId: str, timeout=DEFAULT_TIMEOUT) -> (bytes, dict):
        """
        :return: A tuple (ciphertext, dictionary of parameters).
        """
        cmd = """{"command":"db_retrieve_item", "data":"%s"}""" % flagId
        obj = GameServer.sendCommand(host, port, cmd, "", timeout=timeout)
        try:
            obj = json.loads(obj)
        except Exception as e:
            raise MumbleException(f"Got {type(e)}, while parsing json. Exception message: {e}")
        if not "status" in obj.keys():
            raise MumbleException("Missing key: status")
        if obj["status"] == "error":
            if not "msg" in obj.keys():
                raise MumbleException("Missing key: msg")
            if obj["msg"] == "No such item!":
                raise FlagMissingException("")
            raise MumbleException("Got answer: " + json.dumps(obj))
        if obj["status"] != "db_retrieve_success":
            raise MumbleException(
                f"""Unexpected status: got "{obj["status"]}", " expected "db_retrieve_success". Full response: {json.dumps(obj)}""")
        if not "data" in obj.keys():
            raise MumbleException("Missing key: data")
        if not "params" in obj.keys():
            raise MumbleException("Missing key: params")
        try:
            cipher = unhexlify(obj["data"])
        except:
            raise MumbleException("Invalid cipher format")
        return cipher, obj["params"]

    @staticmethod
    def storeCipher(host: str, port: int, flagId: str, cipher: bytes, params: dict, timeout=DEFAULT_TIMEOUT):
        import os
        from cryptostuff import loadKey, sign
        cmd = """{"command":"db_store_item", "id":"%s", "cipher":"%s", "params":%s}""" % (
            flagId, hexlify(cipher).decode("utf-8"), json.dumps(params))
        keyfile = os.path.dirname(__file__) + "/gs_private.pem"
        key = loadKey(keyfile, True)
        sig = hexlify(sign(cmd, key)).decode("utf-8")
        # print("signature:", repr(sig))
        obj = GameServer.sendCommand(host, port, cmd, sig, timeout=timeout)
        try:
            obj = json.loads(obj)
        except Exception as e:
            raise MumbleException(f"Got {type(e)}, while parsing json. Exception message: {e}")
        if not "status" in obj.keys():
            raise MumbleException("Missing key: status")
        if obj["status"] != "db_store_success":
            raise MumbleException(
                f"""Unexpected status: got "{obj["status"]}", " expected "db_store_success". Full response: """ + json.dumps(
                    obj))
        return True

    @staticmethod
    def encryptAndStore_DONT_USE_ME_I_PRODUCTIVE_SCRIPTS(host: str, port: int, flagId: str, msg: str, method: str,
                                                         timeout=DEFAULT_TIMEOUT):
        """
        :return: the parameter hash
        """
        # storeCipher could fail (e.g. with an OfflineException). But we still need the parameter hash to construct the full known plaintext!
        cipher, params, msg = GameServerCrypto.encrypt(msg, method)
        GameServer.storeCipher(host, port, flagId, cipher, params, timeout=timeout)
        return msg[-SHA512_HEX_LEN:]

    @staticmethod
    def checkFlag(host: str, port: int, flagId: str, msg: str, timeout=DEFAULT_TIMEOUT):
        """
        :return: Whether the flag was potentially re-encrypted
        """
        # Again, we have a similar problem as above. But the return value here is not really that important and we are allowed to have gaps in
        # the "fixing history".
        cipher, params = GameServer.retrieveCipher(host, port, flagId, timeout=timeout)
        # print("checking: ", msg, params)
        GameServerCrypto.checkEncryption(msg, cipher, params)
        origParamHash = msg[-SHA512_HEX_LEN:]
        currentParamHash = GameServerCrypto.parameterHash(params)
        return currentParamHash != origParamHash, params

    @staticmethod
    def generatePlaintext():
        from gamelib import usernames
        import random
        from units import STUPID_UNITS
        return f"Guest {usernames.generate_name()} will bring {str(random.randrange(-42,42))} {random.choice(STUPID_UNITS)} of {random.choice(usernames.USERNAME_ADJECTIVES)} %s to the party."



class GameServerCrypto:
    @staticmethod
    def parameterHash(parameters):
        import json
        from hashlib import sha512
        m = json.dumps(parameters,
                       sort_keys=True) + "42+4711+1337 blablubb!+klammerzu(die sha-klammer)und.digest()fehltauchnochundesmusseinhexdigestseindamitdiezeichenauchprintablesind!"
        return sha512(m.encode()).hexdigest()

    @staticmethod
    def encrypt(msg: str, method):
        if method not in methodsEncryptDict:
            raise MumbleException(f'Unknown method: "{method}"')
        return methodsEncryptDict[method](msg)

    @staticmethod
    def encryptPlain(msg: str):
        from cryptostuff import encrypt

        params = {"method": "plain"}
        msg += GameServerCrypto.parameterHash(params)
        return encrypt("plain", msg.encode("utf-8"), params), params, msg

    @staticmethod
    def encryptCaesar(msg: str):
        import random
        from cryptostuff import encrypt

        key = random.randint(-2 ** 31, 2 ** 31)
        params = {"key": key, "method": "caesar"}
        msg += GameServerCrypto.parameterHash(params)
        return encrypt("caesar", msg.encode("utf-8"), params), params, msg

    @staticmethod
    def encryptRsaSmallFactor(msg: str):
        # We have a non-negligible probability for the message not being coprime to n. Hence, we need to check for this case.
        from Crypto.Util.number import getPrime, GCD, bytes_to_long
        from cryptostuff import encrypt

        bits = 2048
        smallFactorBits = 10
        e = 65537
        i = 0
        while True:
            bigFactorBits = bits - smallFactorBits
            p = getPrime(smallFactorBits)
            q = getPrime(bigFactorBits)
            n = p * q
            params = {"e": e, "n": n, "method": "saar"}
            m = msg + GameServerCrypto.parameterHash(params)
            if GCD(bytes_to_long(m.encode("utf-8")), n) == 1:
                return encrypt("saar", m.encode("utf-8"), params), params, m
            i += 1
            if i > 50:
                raise Exception("Something is wrong in GameServerCrypto.encryptRsaSmallFactor()")

    @staticmethod
    def encryptRsaWiener(msg: str):
        from Crypto.Util.number import getPrime, inverse
        from cryptostuff import encrypt

        bits = 2048
        dBits = 200

        bitsDifferencePQ = 50
        i = 0
        while True:
            p = getPrime(bits // 2 + bitsDifferencePQ)
            q = getPrime(bits // 2 - bitsDifferencePQ)
            d = getPrime(dBits)
            # print d
            e = inverse(d, (p - 1) * (q - 1))
            if len(str(hex(e)[2:])) * 4 < bits - 5:
                # Very hacky and not very accurate but it serves its purpose: We want to ensure that e is not relatively small in relation to n.
                # print "e too small."
                i += 1
                if i > 100:
                    raise Exception("Something is wrong in GameServerCrypto.encryptRsaWiener()")
                continue
            params = {"e": e, "n": p * q, "method": "saar"}
            msg += GameServerCrypto.parameterHash(params)
            return encrypt("saar", msg.encode("utf-8"), params), params, msg

    @staticmethod
    def encryptOTP(msg: str):
        from cryptostuff import encrypt

        numBytes = len(msg) + 512 // 4
        with open("/dev/urandom", "rb") as f:
            randBytes = f.read(numBytes)
        params = {"encryptionkey": hexlify(randBytes).decode("utf-8"), "method": "otp"}
        msg += GameServerCrypto.parameterHash(params)
        return encrypt("otp", msg.encode("utf-8"), params), params, msg

    @staticmethod
    def schwenkNew(msg: str, skipKeyAndIV=False):
        # import Crypto.Random
        import random
        import string
        from cryptostuff import encrypt
        if skipKeyAndIV:
            k = ""
            IV = ""
        else:
            k = "".join([random.choice(string.printable) for _ in range(Schwenk.BLOCK_SIZE)])
            IV = "".join([random.choice(string.printable) for _ in range(Schwenk.BLOCK_SIZE)])
        fakeHash = "".join(random.choice("0123456789abcdef") for _ in range(SHA512_HEX_LEN))   # To be consistent with the other plaintext formats
        msg = k + IV + msg + fakeHash
        cipher, rawParams = Schwenk.encrypt(msg.encode("utf-8"))
        params = {"schwenkerid": rawParams[0], "schwenkingoptions": rawParams[1], "method": "schwenk"}
        cipher1 = encrypt("schwenk", msg.encode("utf-8"), params)
        assert cipher1 == cipher
        return cipher, params, msg

    # @staticmethod
    # def encryptElGamal(msg: str):
    #     # We have a (very small) non-negligible probability for GCD(p,s) != 1 which we need to take into account.
    #     # Hence, we also cannot sensibly use the cryptostuff.encrypt method
    #     from Crypto.Util.number import getPrime, GCD, getRandomRange, bytes_to_long
    #     targetKeyBits = 8 * len(msg) + 2 * 512
    #     bitsPerPrime = 30
    #     i = 0
    #     while True:
    #         p = 1
    #         while len(bin(p)[2:]) < targetKeyBits:
    #             p = p * getPrime(bitsPerPrime)
    #         g = getPrime(len(bin(p)[2:]) - 1)
    #         privateKey = getRandomRange(1, p)
    #         h = pow(g, privateKey, p)
    #         params = {"p": p, "g": g, "h": h, "method": "elgamal"}
    #
    #         y = getRandomRange(1, p)
    #         s = pow(h, y, p)
    #         if GCD(s, p) == 1:
    #             c1 = pow(g, y, p)
    #             m = msg + GameServerCrypto.parameterHash(params)
    #             m = m.encode("utf-8")
    #             m_int = bytes_to_long(m)
    #             # assert m < p
    #             if m_int < p and GCD(m_int, p) == 1:
    #                 c2 = (m_int * s) % p
    #                 c = {"c1": c1, "c2": c2}
    #                 return json.dumps(c).encode("utf-8"), params, m
    #         i += 1
    #         if i > 50:
    #             raise Exception("Something is wrong in GameServerCrypto.encryptElGamal()")

    @staticmethod
    def checkEncryption(msg: str, cipher: bytes, params: dict):
        from cryptostuff import encrypt
        from cryptostuff import methods as encryptMethods
        if "method" not in params:
            raise MumbleException('Missing key: "method"')

        if params["method"] not in encryptMethods:
            raise MumbleException(f'Unknown method: "{params["method"]}"')
        if params["method"] == "elgamal":
            # As the encryption is randomized, we cannot simply encrypt the known plaintext and compare the ciphertexts...
            # TODO: I don't think we can properly check this. With the known plaintext only. Otherwise, ElGamal would be pretty broken.
            pass
        else:
            try:
                cTest = encrypt(params["method"], msg.encode("utf-8"), params)
            except Exception as e:
                raise MumbleException(f"Got exception {e} of type {type(e)} during encryption.")
            if cTest != cipher:
                raise FlagMissingException(f"Incorrect cipher. Got {repr(cipher)}, expected {repr(cTest)}. Params: " + str(params))
            return


cheapMethods = [
    "plain",
    "caesar",
    "rsa_smallfactor",
    "rsa_wiener",
    "otp",
]

methods = cheapMethods + ["schwenk"]

methodsEncryptDict = {
    "plain": GameServerCrypto.encryptPlain,
    "caesar": GameServerCrypto.encryptCaesar,
    "rsa_smallfactor": GameServerCrypto.encryptRsaSmallFactor,
    "rsa_wiener": GameServerCrypto.encryptRsaWiener,
    "otp": GameServerCrypto.encryptOTP,
    # "elgamal": GameServerCrypto.encryptElGamal,
    "schwenk": GameServerCrypto.schwenkNew
}


def cheapMethodFromRoundNumber(round: int):
    # TODO: maybe adjust such that really cheap stuff appears less?
    return cheapMethods[round % len(cheapMethods)]
    # return "plain"


class Schwenk:
    BLOCK_SIZE = 32

    @staticmethod
    def toNumber(m):
        return int(m.hex(), 16)

    @staticmethod
    def F(IV, k):
        return IV + k

    @staticmethod
    def G(k, IV):
        return k * IV.inverse(), Schwenk.F(IV, k)

    @staticmethod
    def H(m, k):
        return k ** m

    @staticmethod
    def encrypt(m):
        from hashlib import sha512
        import pickle
        import random
        from os import path

        assert len(m) >= 64, "Message must have 32 byte key and IV prepended."
        assert len(m) < 13 * Schwenk.BLOCK_SIZE, "Message too long" + str(
            len(m))  # We only have enough parameters for 13+2 blocks
        k = int(m[:Schwenk.BLOCK_SIZE].hex(), 16)
        IV = int(m[Schwenk.BLOCK_SIZE:2 * Schwenk.BLOCK_SIZE].hex(), 16)
        m = m[2 * Schwenk.BLOCK_SIZE:]
        m = Schwenk.pad(m)
        assert len(m) % Schwenk.BLOCK_SIZE == 0
        NUM_PARAM_FILES = 11
        fileIndex = random.randrange(NUM_PARAM_FILES)
        with open(path.dirname(__file__) + "/SchwenkParams%02d.pkl"%fileIndex, "rb") as f:
            params = pickle.load(f)
        blockParams = params[random.randrange(0, len(params))]
        numBlocks = len(m) // Schwenk.BLOCK_SIZE
        cipher = b""
        xorParams = []
        params = []
        for i in range(numBlocks):
            mBlock = m[i * Schwenk.BLOCK_SIZE:(i + 1) * Schwenk.BLOCK_SIZE]
            cipherBlock, k, IV, blockXor, cleanParams = \
                Schwenk.encryptBlock(blockParams[i], IV, k, Schwenk.toNumber(mBlock))
            cipher += cipherBlock
            xorParams += [blockXor]
            params += [cleanParams]

        mac = sha512(cipher).digest()
        cipherBlock, k, IV, blockXor, cleanParams = \
            Schwenk.encryptBlock(blockParams[numBlocks], IV, k, Schwenk.toNumber(mac[:Schwenk.BLOCK_SIZE]),
                                 True)
        cipher += cipherBlock
        xorParams += [blockXor]
        params += [cleanParams]

        cipherBlock, k, IV, blockXor, cleanParams = \
            Schwenk.encryptBlock(blockParams[numBlocks + 1], IV, k, Schwenk.toNumber(mac[Schwenk.BLOCK_SIZE:]),
                                 True)
        cipher += cipherBlock
        xorParams += [blockXor]
        params += [cleanParams]
        return cipher, [xorParams, params]

    @staticmethod
    def encryptBlock(blockParams, IV, k, m, isHashBlock=False):
        from algebraNewUnobf import GFElement
        from Crypto.Util.number import long_to_bytes, GCD
        import Crypto.Random
        # lastBlockRealSize = max(lastBlockRealSize, len(bin(m)[2:]), len(bin(k)[2:]),
        #                         len(bin(IV)[2:]))  # Perhaps too generous but safe.

        # nextSize = lastBlockRealSize + 3
        # base, power, modCoeffs = SchwenkNew.getBlockModulusRaw(nextSize)
        base, power, modCoeffs = blockParams
        order = base ** power - 1

        # TODO: Think of a nice way to make the case distinction seem like an honest mistake.
        kElement = GFElement.fromDenseNumber(k, base, power, modCoeffs)
        IVElement = GFElement.fromDenseNumber(IV, base, power, modCoeffs)
        kNew, IVNew = Schwenk.G(kElement, IVElement)
        k = kNew.toDenseNumber()
        while True:
            # We need to xor the message block with some properly chosen random stuff to ensure that exponentiation is efficiently invertible.
            r = Schwenk.toNumber(Crypto.Random.get_random_bytes(Schwenk.BLOCK_SIZE))
            if isHashBlock:
                m1 = m ^ r
                if GCD(order, m1) == 1:
                    break
            else:
                k1 = k ^ r
                if GCD(order, k1) == 1:
                    break

        if isHashBlock:

            cipherElement = Schwenk.H(m1, kNew)
        else:
            cipherElement = Schwenk.H(k1, GFElement.fromDenseNumber(m, base, power, modCoeffs))

        return long_to_bytes(cipherElement.toDenseNumber()), kNew.toDenseNumber(), IVNew.toDenseNumber(), r, [
            base] + modCoeffs

    @staticmethod
    def pad(m):
        if len(m) % Schwenk.BLOCK_SIZE == 0:
            m += Schwenk.BLOCK_SIZE * chr(Schwenk.BLOCK_SIZE).encode("utf-8")
        else:
            m += (Schwenk.BLOCK_SIZE - len(m) % Schwenk.BLOCK_SIZE) * chr(
                Schwenk.BLOCK_SIZE - len(m) % Schwenk.BLOCK_SIZE).encode("utf-8")
        return m
