#SRC: goaldi original
#  test numeric library functions

procedure main() {
	write("\narithmetic:")
	every testarith(-%phi | "0" | 1.0 | 2 | %e | %pi)
	testmeans()
	write("\ntrigonometry:")
	every testtrig(-%phi | "0" | 1.0 | 2 | %e | %pi)
	write("\nhypertrigonometry:")
	every testhyper(-%phi | "0" | 1.0 | 2 | %e | %pi)
	testpot()
	testrad()
	write("\nbased logarithms:")
	every testlogb(0 | 1 | %phi | 2 | 8 | 32 | 100 | 1012)
	write("\natan2:")
	every testa2(0 | %phi | %pi, 0 | %e | %pi)
	write()
	write("\ngcd:")
	testgcd()
	write("\nrandgen:")
	testrandgen()
}

procedure testout(v) {
	writes(v, ":")
	every apply(type | image | number | string, v)
	write()
	return
}

procedure testarith(v) {
	writes(v, ":")
	every apply(abs | integer | ceil | floor | log | sqrt | cbrt | exp, v)
	write()
	return
}

procedure testtrig(v) {
	writes(v, ":")
	every apply(sin | cos | tan | asin | acos | atan, v)
	write()
	return
}

procedure testhyper(v) {
	writes(v, ":")
	every apply(sinh | cosh | tanh | asinh | acosh | atanh, v)
	write()
	return
}

procedure testmeans() {
	local L := [1, 1, 2, 3, 5, 8, 13, 21, 42]
	write()
	writes("amean: ", amean ! L, "   ")
	writes("gmean: ", gmean ! L, "   ")
	writes("hmean: ", hmean ! L, "   ")
	writes("qmean: ", qmean ! L)
	write()
}

procedure testpot() {
	# local x, y
	local x
	local y

	writes("hypot: ")
	every (x := 2 | 3 | 4) & (y := 3 | 5 | 8) do
		writes(" ",x,":",y,":",hypot(x,y))
	write()
	return
}

procedure testrad() {
	local d
	local r
	writes("dtortod:")
	every d := -45 | 0 | 30 | 60 | 90 | 360/%pi | 180 do {
		r := dtor(d)
		writes(" ", d, "/", r, "/", rtod(r))
	}
	write()
	return
}

procedure testlogb(v) {
	local b
	writes(v, ":")
	every b := nil | %e | 2 | 4 | 10 do {
		writes(" log(,", b, ")", log(v, b))
	}
	write()
	return
}

procedure testa2(x,y) {
	return writes(" (",x,",",y,")",atan(x,y))
}

procedure apply(p, x, y) {
	if \y then {
		writes(" ", string(p)[3:0], "()", p(x, y) | "--")
	} else {
		writes(" ", string(p)[3:0], "()", p(x) | "--")
	}
	return
}

procedure testgcd() {
	writes(" a.", gcd(0))
	writes(" b.", gcd(1))
	writes(" c.", gcd(3))
	writes(" d.", gcd(-5))
	writes(" e.", gcd(0,3))
	writes(" f.", gcd(3,0))
	writes(" g.", gcd(15,12))
	writes(" h.", gcd(15,-12))
	writes(" i.", gcd(30,42,15))
	writes(" j.", gcd(30,42,18,28))
	writes(" k.", gcd(30,42,18,-36))
	writes(" l.", gcd(30,42,18,28,19))
	writes(" m.", gcd(0,0,0))
	write()
	return
}

procedure testrandgen() {
	seed(747)
	every writes(" ", "  stdgen:" | |?999 \ 5 | "\n")
	seed(747)	# should survive through all the following
	showgen()
	showgen(0)
	showgen(1)
	showgen(1)
	showgen(314159)
	every writes(" ", "  stdgen:" | |?999 \ 5 | "\n")
}

procedure showgen(i) {
	local g := randgen(i)
	writes("   randgen(", image(i), "):")
	every writes(" ", type(g) | ":" | |g.Intn(999) \ 5 | "\n")
}
