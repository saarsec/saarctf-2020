holmol "cavelib.sl";
holmol "stdlib.sl";

eija is_free(x1: int, y1: int, direction: int) gebbtserick int: {
	var x: int = x1;
	var y: int = y1;
	falls direction == 0: y = y - 1;
	sonschd: falls direction == 1: x = x + 1;
	sonschd: falls direction == 2: y = y + 1;
	sonschd: falls direction == 3: x = x - 1;

	serick 1 - mach issdowas(x, y);
}

eija random_direction() gebbtserick int: {
	var r: int = mach ebbes() unn 63;
	falls r < 32: serick 0;
	falls r < 48: serick 1;
	falls r < 56: serick 2;
	serick 3;
}

eija main() gebbtserick int: {
	var direction: int = 0;
	var i: int = 4000;
	solang i > 0: {
		direction = mach random_direction();
		solang mach is_free(wo_x, wo_y, direction) unn i > 0: {
			falls direction == 0: ruff;
			sonschd: falls direction == 1: doniwwer;
			sonschd: falls direction == 2: runner;
			sonschd: falls direction == 3: riwwer;
			i = i - 1;
			falls (i % 5) == 0: direction = mach random_direction();
		}
	}
	ferdisch;
}
