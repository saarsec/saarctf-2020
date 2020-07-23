from __future__ import print_function, unicode_literals
import os
import shutil
import sys
import subprocess

os.chdir(os.path.dirname(os.path.abspath(__file__)))

# TEST:
# docker run --name schlossberg -it -v "`pwd`/dist:/opt/dist" debian:buster bash

PREPARATION = '''
apt-get update && apt-get install build-essential g++ cmake libmicrohttpd-dev llvm-5.0 llvm-5.0-dev
useradd -m schlossberg
'''

INSTALLATION = '''
# todo extract
# compile backend
pushd . ; ( cd backend && mkdir build && cd build && cmake .. && make -j 4 && find -iname '*.o' -delete ) ; popd
# set permissions
chown schlossberg:schlossberg data
# todo configure nginx
'''

RUN = '''
cd backend/build && ./SchlossbergCaveServer
'''


def filter(text):
	lines = text.split('\n')
	output = []
	for l in lines:
		if '--- DIST_REMOVE ---' in l:
			break
		output.append(l)
	return '\n'.join(output)


def copy_backend_recursive(src, dest):
	blacklist = ['python', 'pybind', 'debug', 'build', '.idea', '.git', 'todo', 'codesamples_old', 'tests']

	for fname in os.listdir(src):
		if any([b in fname for b in blacklist]):
			print('(excluded)', src + fname)
			continue

		if os.path.isdir(src + fname):
			os.mkdir(dest + fname)
			copy_backend_recursive(src + fname + '/', dest + fname + '/')
		else:
			if fname == 'CMakeLists.txt':
				with open(src + fname, 'r') as f1:
					with open(dest + fname, 'w') as f2:
						f2.write(filter(f1.read()))
			else:
				shutil.copy(src + fname, dest + fname)


if __name__ == '__main__':
	# Create folder structure
	if os.path.exists('dist'):
		shutil.rmtree('dist')
	os.mkdir('dist')
	os.mkdir('dist/backend')
	os.mkdir('dist/data')

	# Create frontend
	subprocess.check_call("cd frontend && npm run build", shell=True)
	subprocess.check_call("mv frontend/dist/SchlossbergCaves dist/frontend", shell=True)
	subprocess.check_call("cd dist/frontend && find -iname '*.map' -delete", shell=True)

	# Create backend
	copy_backend_recursive('backend/', 'dist/backend/')

	# Create default data
	shutil.copytree('data/cave-templates', 'dist/data/cave-templates')

	# Test mode?
	if len(sys.argv) > 1 and sys.argv[1] == '--test':
		os.mkdir('dist/testserver')
		os.mkdir('dist/scripts')
		shutil.copy('testserver/server.js', 'dist/testserver/')
		shutil.copy('testserver/package.json', 'dist/testserver/')
		shutil.copytree('scripts/gameserver', 'dist/scripts/gameserver')
		shutil.copytree('scripts/gamelib', 'dist/scripts/gameserver/gamelib')
		shutil.copy('scripts/generate-paths.py', 'dist/scripts/')
		with open('dist/README.txt', 'w') as f:
			f.write('''
Backend port:  9081
Frontend port: 9080

## Preparation:
{}

## Installation:
{}

## Run C++ backend server:
{}

## Run test server (will be nginx in the final ctf):
cd ../testserver
npm install
node server.js 

## Run gameserver testscript
cd ../scripts
python3 generate-paths.py   # one-time initialization, takes time
cd gameserver
python3 schlossberg_cave_interface.py [roundnumber]   # store and retrieve flags for round (roundnumber is optional)

'''.format(PREPARATION, INSTALLATION, RUN))

	# Compress
	if '--compress' in sys.argv:
		os.system("cd dist && tar -zcvf schlossberg.tar.gz *")

	print('')
	print('## Preparation: ')
	print(PREPARATION)
	print('')
	print('## Installation: ')
	print(INSTALLATION)
	print('')
	print('## Run:')
	print(RUN)
