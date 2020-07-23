import json
import sqlite3
from binascii import hexlify, unhexlify

class db:
    def __init__(self):
        self.filename = "storage.db"

    def storeData(self, id : str, ciphertext : bytes, parameters : str):
        connection = sqlite3.connect(self.filename)
        c = connection.cursor()
        c.execute("INSERT INTO messages VALUES (?, ?, ?)", [id, hexlify(ciphertext).decode(), json.dumps(parameters)])
        connection.commit()
        connection.close()


    def getData(self, id : str):
        connection = sqlite3.connect(self.filename)
        c = connection.cursor()
        c.execute("SELECT msg, params FROM messages WHERE id=?", [id])
        result = c.fetchone()
        if result is None:
            return None
        return (unhexlify(result[0]), json.loads(result[1]))
