import sys

from gamelib import Team
from cryptoServiceInterface import CryptoServiceInterface

if __name__ == '__main__':
    # TEST CODE
    team = Team(12, 'n00bs', '127.0.0.1')
    round = int(sys.argv[1]) if len(sys.argv) > 1 else 2
    print("Round:", round)
    service = CryptoServiceInterface(0)

    print('[1] Integrity check...')
    service.check_integrity(team, round)
    print('Passed.')

    print('[2] Store flags...')
    flags = service.store_flags(team, round)
    print('Done ({} flags).'.format(flags))

    print('[3] Retrieve the flags in the next round')
    flags = service.retrieve_flags(team, round)
    print('Done ({} flags).'.format(flags))