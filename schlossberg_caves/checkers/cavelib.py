from __future__ import print_function

import random
from builtins import bytes

import os
import struct
from collections import OrderedDict
from typing import Tuple, List, Optional

import numpy
import time
from queue import Queue

CAVEPATH = os.path.dirname(os.path.abspath(__file__)) + '/cave-templates/'
DISTPATH = os.path.dirname(os.path.abspath(__file__)) + '/cache/'
if not os.path.exists(CAVEPATH):
	CAVEPATH = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__)))) + '/data/cave-templates/'
	DISTPATH = os.path.dirname(os.path.dirname(os.path.abspath(__file__))) + '/cache/'  # TODO

DISABLE_CAVE_CACHE = True

"""

YOU WANT TO RUN THIS SCRIPT WITH pypy!

All coordinates are (x,y).
Indices are [0..x-1].
Data is accessed data[x,y].
Directions: 1,2,4,8 = up right down left
"""

UP = 1
RIGHT = 2
DOWN = 4
LEFT = 8


def neighbour_points(point: Tuple[int, int]) -> List[Tuple[int, Tuple[int, int]]]:
	x, y = point
	return [(UP, (x, y - 1)), (RIGHT, (x + 1, y)), (DOWN, (x, y + 1)), (LEFT, (x - 1, y))]


def invert(direction: int) -> int:
	if direction == UP: return DOWN
	if direction == RIGHT: return LEFT
	if direction == DOWN: return UP
	if direction == LEFT: return RIGHT


class Cave:
	def __init__(self, start, data):
		"""
		:param (int, int) start:
		:param numpy.ndarray data:
		"""
		self.start = start
		self.data = data
		# Matrix: point -> distance to start
		self.distances = None
		# Matrix: point -> optimal direction towards start
		self.shortestPaths = None

	def __str__(self):
		return 'Cave' + str(self.data.shape)

	def __getitem__(self, item):
		return self.data.__getitem__(item)

	def calculateShortestPaths(self):
		self.distances = numpy.full(self.data.shape, 0xffff, numpy.uint16)
		self.distances[self.start] = 0
		self.shortestPaths = numpy.zeros(self.data.shape, numpy.uint8)

		queue = Queue()
		queue.put(self.start)
		while not queue.empty():
			p = queue.get()
			d = self.distances[p] + 1
			for direction, p2 in neighbour_points(p):
				if not self.data[p2]: continue
				# check if we have a better distance
				if self.distances[p2] > d:
					self.distances[p2] = d
					self.shortestPaths[p2] = invert(direction)
					queue.put(p2)
				# check if we have an additional path
				elif self.distances[p2] == d:
					self.shortestPaths[p2] += invert(direction)

	def getDistanceToStart(self, pos: Tuple[int, int]) -> Optional[int]:
		"""
		:param (int,int) pos:
		:rtype: int|None
		:return:
		"""
		if not self.distances:
			self.calculateShortestPaths()
		return self.distances[pos] if self.distances[pos] < 0xffff else None

	def getRandomShortestPathToStart(self, pos: Tuple[int, int]) -> Optional[List[int]]:
		"""
		Path with directions how to navigate from a given point to the center/start
		:param (int,int) pos:
		:return:
		"""
		path = []
		if self.distances[pos] >= 0xffff:
			return None
		while pos != self.start:
			directions = self.shortestPaths[pos]
			nextPositions = []
			for direction, pos2 in neighbour_points(pos):
				if direction & directions:
					nextPositions.append((direction, pos2))
			direction, pos = random.choice(nextPositions)
			path.append(direction)
		return path

	@staticmethod
	def load(fname) -> 'Cave':
		with open(fname, 'rb') as f:
			bindata = bytes(f.read())
		dimension = struct.unpack('II', bindata[:8])  # width / height
		start_row, start_col = struct.unpack('II', bindata[8:16])
		start = (start_col, start_row)
		data = numpy.zeros(dimension[0] * dimension[1], numpy.uint8)
		i = 0
		for b in bindata[16:]:
			for s in range(8):
				data[i] = (b >> s) & 1
				i += 1
		assert i == len(data)
		data = data.reshape((dimension[1], dimension[0]))
		data = data.T
		return Cave(start, data)

	def saveDistances(self, fname: str):
		# numpy.save(fname+'.dist', self.distances, False)
		# numpy.save(fname+'.path', self.shortestPaths)
		with open(fname + '.dist', 'wb') as f:
			f.write(self.distances.tobytes())
		with open(fname + '.path', 'wb') as f:
			f.write(self.shortestPaths.tobytes())

	def loadDistances(self, fname: str) -> bool:
		if os.path.exists(fname + '.path'):
			# self.distances = numpy.load(fname+'.dist')
			# self.shortestPaths = numpy.load(fname+'.path')
			with open(fname + '.dist', 'rb') as f:
				self.distances = numpy.frombuffer(f.read(), dtype=numpy.uint16).reshape(self.data.shape)
			with open(fname + '.path', 'rb') as f:
				self.shortestPaths = numpy.frombuffer(f.read(), dtype=numpy.uint8).reshape(self.data.shape)
			return True
		else:
			return False


class CaveCacheImpl(dict):
	def __missing__(self, key) -> Cave:
		print('Loading', CAVEPATH + key)
		cave = Cave.load(CAVEPATH + key)
		cave.loadDistances(DISTPATH + key)
		if not DISABLE_CAVE_CACHE:
			self[key] = cave
		return cave


CaveCache = CaveCacheImpl()

if __name__ == '__main__':
	times = OrderedDict()

	t = time.time()
	cave = Cave.load(CAVEPATH + 'schlossberg_3.cave')
	print(cave)
	print(cave.data)
	times['load'] = time.time() - t

	t = time.time()
	cave.calculateShortestPaths()
	times['shortest paths'] = time.time() - t

	t = time.time()
	cave.saveDistances('cave3')
	times['save dist'] = time.time() - t

	t = time.time()
	print(cave.loadDistances('cave3'))
	times['load dist'] = time.time() - t

	print("\n\n=== TIMES ===")
	for name, t in times.items():
		print(name, round(t * 1000), 'ms')
