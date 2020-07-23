import sys
from binascii import unhexlify
import serverlogic

if __name__ == "__main__":
    print("command:")
    cmd = sys.stdin.readline().strip()
    if not cmd:
        sys.exit(0)
    print("signature:")
    sig = b""
    try:
        sig = unhexlify(sys.stdin.readline().strip())
    except:
        print("invalid signature format", file=sys.stderr)
    sl = serverlogic.ServerLogic()
    res = sl.processCommand(cmd, sig)
    print(res, file=sys.stderr)
    print(res)