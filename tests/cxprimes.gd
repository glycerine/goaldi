#SRC: icon/cxprimes.icn
#  prime number generation using co-expressions

procedure main(limit) {
	local n := number(limit) | 100
	local s := create (2 to n)
	while (^x := @s) do {
		write(x)
		s := create sieve(x, s)
	}
}

procedure sieve(x, s) {
	local t

	while t := @s do {
		if t % x ~= 0 then suspend t
	}
}
