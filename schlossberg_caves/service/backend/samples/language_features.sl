// Saarlang sample file

holmol "stdlib.sl";
holmol "cavelib.sl";


const test1: int = 1234;


// simple function
eija add(x: int, y: int) gebbtserick int: {
	serick x + y;
}


// control flow
eija testcf(x: int) gebbtserick int: {
	var a: int = x;
	solang a > 0: {
		falls a % 123 == 0:
			serick a;
		sonschd:
			a = a - 1;
	}
	serick a;
}


// arrays
eija testarray(x: int) gebbtserick int: {
	var a: lischd int = neie lischd int(x);
	var b: lischd byte(120);

	var i: int = 0;
	solang i < grees a: {
		a@i = i + 1;
		i = i + 1;
	}

	var c: lischd int = a;
	c@1 = 3;
	serick c@0 * 100 + c@1;
}


// navigation
eija move_around_in_cave() gebbtserick int: {
	mach sahmol(wo_x);
	mach sahmol_ln(wo_y);
	mach sahmol_ln(mach issdowas(wo_x, wo_y-1) );

	runner;
	ruff;
	doniwwer;
	riwwer;
	ferdisch;
}


eija main() gebbtserick int: {
	mach sahmol_ln( mach add(1,2) );
	mach sahmol_ln( mach testcf(249) );
	mach sahmol_ln( mach testarray(15) );
	mach move_around_in_cave();
	serick 1337;
}

