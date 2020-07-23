from __future__ import print_function
import os
import sys
from multiprocessing import Process

sys.path.append(os.path.dirname(os.path.abspath(__file__)))

import cavelib


def generate_files(files):
	for fname in files:
		print('-', fname, '...')
		cave = cavelib.Cave.load(cavelib.CAVEPATH + fname)
		cave.calculateShortestPaths()
		cave.saveDistances(cavelib.DISTPATH + fname)


def main():
	if not os.path.exists(cavelib.DISTPATH):
		os.mkdir(cavelib.DISTPATH)
	files = []
	for fname in sorted(os.listdir(cavelib.CAVEPATH)):
		if not fname.endswith('.cave'): continue
		files.append(fname)
	# split file list
	n = len(files) // 8
	slices = [files[i:i + n] for i in range(0, len(files), n)]

	processes = [Process(target=generate_files, args=(slice,)) for slice in slices]
	for p in processes:
		p.start()
	for p in processes:
		p.join()

	print('Done.')


if __name__ == '__main__':
	main()
