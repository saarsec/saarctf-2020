import json
import cryptostuff
import database
from binascii import unhexlify, hexlify
import sys
from sqlite3 import IntegrityError

ERROR_INVALID_CLIENT = '{"status":"error", "msg":"Untrusted client!"}'
ERROR_INVALID_COMMAND = '{"status":"error", "msg":"Unknown command or invalid format!"}'
ERROR_INVALID_ITEM = '{"status":"error", "msg":"No such item!"}'
ERROR_INTERNAL_SERVER_ERROR = '{"status":"error", "msg":"Internal server error!", "exception_type":"%s", "exception_message":"%s"}'
ERROR_ITEM_EXISTS = '{"status":"error", "msg":"Item already exists!"}'

KEYFILE = "trustedkeys"
db = None


class ServerLogic:
    def __init__(self):
        self.db = database.db()
        self.trustedClientKeys = []

        with open(KEYFILE, "r") as f:
            keyfiles = f.read()
        keyfiles = keyfiles.replace("\r", "")  # Tribute to Microsoft. Just in case...
        keyfiles = keyfiles.split("\n")
        for f in keyfiles:
            if f.endswith(".pem"):
                self.trustedClientKeys.append(cryptostuff.loadKey(f, False))

    def isMessageFromTrustedClient(self, message: str, signature: bytes):
        for key in self.trustedClientKeys:
            if cryptostuff.verify(message.encode("utf-8"), signature, key):
                return True
        return False

    def processCommand(self, message: str, signature: bytes):
        print("received message:", message, file=sys.stderr)
        msg = json.loads(message)
        if "command" not in msg:
            return ERROR_INVALID_COMMAND
        command = msg["command"]
        try:
            if command == "db_retrieve_item":
                dbId = msg["data"]
                item = self.loadItem(dbId)
                if item is None:
                    return ERROR_INVALID_ITEM
                return '{"status":"db_retrieve_success", "data":"' + hexlify(item[0]).decode(
                    "utf-8") + '","params":' + json.dumps(item[1]) + '}'

            elif command == "db_store_item":
                if not self.isMessageFromTrustedClient(message, signature):
                    return ERROR_INVALID_CLIENT
                id = msg["id"]
                params = msg["params"]
                cipher = unhexlify(msg["cipher"])
                try:
                    self.storeItem(id, cipher, params)
                except IntegrityError:
                    return ERROR_ITEM_EXISTS
                return '{"status":"db_store_success"}'

            elif command == "check_item":
                if not self.isMessageFromTrustedClient(message, signature):
                    return ERROR_INVALID_CLIENT
                dbId = msg["id"]
                item = self.loadItem((dbId))
                if item is None:
                    return ERROR_INVALID_ITEM
                c, params = item
                p = msg["msg"].encode("utf-8")
                try:
                    c1 = cryptostuff.encrypt(params["method"], p, params)
                    if c != c1:
                        return '{"status":"check_fail"}'
                    return '{"status":"check_success"}'
                except:
                    return '{"status":"check_fail"}'

            else:
                return ERROR_INVALID_COMMAND
        except Exception as e:
            return ERROR_INTERNAL_SERVER_ERROR % (type(e), e)

    def storeItem(self, id: str, cipher: bytes, params: str):
        return self.db.storeData(id, cipher, params)

    def loadItem(self, id: str):
        return self.db.getData(id)
