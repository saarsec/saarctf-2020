from gamelib import *

PORT = 4711


class CryptoServiceInterface(ServiceInterface):
    name = "Let's Schwenkrypt"
    flag_id_types: List[str] = ["hex32", "hex32"]

    def check_integrity(self, team, round):
        import socket
        from gameserver import SimpleSocket
        toException = False
        try:
            s = SimpleSocket.connect(team.ip, PORT)

            def expectLine(line):
                l = s.recvline()
                if l != line:
                    raise MumbleException(f"""Unexpected behaviour: got "{l}", expected "{line}".""")

            expectLine("command:")
            s.close()
        except (socket.timeout, ConnectionRefusedError):
            toException = True
            try:
                s.close()
            except:
                pass
        if toException:
            raise OfflineException(f"Team {team.ip} seems to be offline")

    def store_flags(self, team, round):
        from gameserver import cheapMethodFromRoundNumber, cheapMethods, GameServer, GameServerCrypto, SHA512_HEX_LEN, \
            methods

        flagsStored = 0
        method = cheapMethodFromRoundNumber(round)
        flagIndex = cheapMethods.index(method)
        flagId = self.get_flag_id(team, round, 0)
        flag = self.get_flag(team, round, flagIndex)
        msg = GameServer.generatePlaintext() % flag

        # Check if the flag is already stored
        # We could in theory just try to store the flag and check for an IntegrityException from sqlite.
        # However, this approach might not work if teams tinker with their DB layout.
        flagAlreadyStored = True
        try:
            GameServer.retrieveCipher(team.ip, PORT, flagId)
        except FlagMissingException:
            flagAlreadyStored = False

        if not flagAlreadyStored:
            cipher, params, msg = GameServerCrypto.encrypt(msg, method)
            paramHash = msg[-SHA512_HEX_LEN:]
            self.store(team, round, f"paramhash_{flagIndex}", paramHash)
            self.store(team, round, f"msg_{flagIndex}", msg)
            GameServer.storeCipher(team.ip, PORT, flagId, cipher, params)
            flagsStored += 1

        flagIndex = methods.index("schwenk")
        flagId = self.get_flag_id(team, round, 1)
        flag = self.get_flag(team, round, flagIndex)
        msg = GameServer.generatePlaintext() % flag

        flagAlreadyStored = True
        try:
            GameServer.retrieveCipher(team.ip, PORT, flagId)
        except FlagMissingException:
            flagAlreadyStored = False

        if not flagAlreadyStored:
            cipher, params, msg = GameServerCrypto.schwenkNew(msg, skipKeyAndIV=False)
            paramHash = GameServerCrypto.parameterHash(params)
            self.store(team, round, f"paramhash_{flagIndex}", paramHash)
            self.store(team, round, "schwenk_msg", msg)  # Random key + IV are added.
            GameServer.storeCipher(team.ip, PORT, flagId, cipher, params)
            flagsStored += 1

        return flagsStored

    def retrieve_flags(self, team, round):
        from gameserver import cheapMethodFromRoundNumber, cheapMethods, GameServer, GameServerCrypto, methods

        method = cheapMethodFromRoundNumber(round)
        flagIndex = cheapMethods.index(method)
        flagId = self.get_flag_id(team, round, 0)
        flag = self.get_flag(team, round, flagIndex)
        msg = self.load(team, round, f"msg_{flagIndex}")
        paramHash = self.load(team, round, f"paramhash_{flagIndex}")
        if paramHash is None or msg is None:
            # We didn't event try to store this flag
            raise FlagMissingException(f"Missing flag for team {team.id} in round {round}.")

        wasReencrypted, _ = GameServer.checkFlag(team.ip, PORT, flagId, msg)
        self.store(team, round, f"fixed_{flagIndex}", wasReencrypted)

        flagIndex = methods.index("schwenk")
        flagId = self.get_flag_id(team, round, 1)
        msg = self.load(team, round, "schwenk_msg")
        paramHash = self.load(team, round, f"paramhash_{flagIndex}")
        if paramHash is None or msg is None:
            raise FlagMissingException(f"Missing flag for team {team.id} in round {round}.")

        _, storedParams = GameServer.checkFlag(team.ip, PORT, flagId, msg)
        storedParamHash = GameServerCrypto.parameterHash(storedParams)
        self.store(team, round, f"fixed_{flagIndex}", storedParamHash != paramHash)

        return 2
