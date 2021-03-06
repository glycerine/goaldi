#SRC: icon/arith.icn
# test arithmetic operators and numeric coercion

procedure main() {

	numtest(0, 0)
	numtest(0, 1)
	numtest(0, -1)
	numtest(1, 0)
	numtest(1, 1)
	numtest(1, 2)
	numtest(7, 3)
	numtest(3, 8)
	numtest(6.2, 4)
	numtest(8, 2.5)
	numtest(5.4, 1.2)
	numtest("1", 2.5)
	numtest("3.4", 1.7)
	numtest("5", " 5")
	numtest(0., 0.)
	numtest(0., 1.)
	numtest(0., -1.)
	numtest(1, -2)
	numtest(1., -2.)
	numtest(-3, 2)
	numtest(-3., "2.")
	numtest(-6, -3)
	numtest(-6., -3.)
	write()

	every (^i := -9 | 0 | 5 | 191) & (^j := -23 | 0 | 9 | 61) do
		bitcombo(i, j)
	write()

	shifttest()
	write()

	every pow(-3 to 3, -3 to 3)
	every pow(.5 | 1 | 1.5, (-3 to 3) / 2.0)
	every pow(-1.5 | -1.0 | -.5 | 0.0, -3 to 3)
}

procedure numtest(a, b) {
	wr4(+a)
	wr4(b)
	wr4(abs(a))
	wr5(-b)
	wr5(a + b)
	wr5(a - b)
	wr5(a * b)
	wr5(if b ~= 0 then a // b else "-/-")
	wr5(if b ~= 0 then a / b else "-/-")
	wr5(if b ~= 0 then a % b else "-%-")
	wr5(-b)
	wr5(a < b  | "---")
	wr4(a <= b | "---")
	wr4(a = b  | "---")
	wr4(a ~= b | "---")
	wr4(a >= b | "---")
	wr4(a > b  | "---")
	write()
	return
}

procedure bitcombo(i, j) {
	every wr5(i | j | icom(i) | icom(j) |
		iand(i,j) | ior(i,j) | ixor(i,j) | iclear(i,j))
	write()
	return
}

procedure wr4(n) {			# write in 4 chars
	return printf(" %3s", string(n))
}

procedure wr5(n) {			# write in 5 chars
	return printf(" %4s", string(n))
}

procedure pow(m, n) {
	if m = 0 & n <= 0 then
		return fail
	local v := m ^ n
	printf("%f ^ %f = %f\n", m, n, v)
	return
}

procedure shifttest() {
	every local n := 10 to -10 by -1 do
		printf("shift %-2.0f %5.0f %8.0f %8.0f\n",
			n, ishift(1, n), ishift(1703, n), ishift(-251, n))
}
