#SRC: icon/sieve.icn
#
#          S I E V E   O F   E R A T O S T H E N E S
#

#  This program illustrates the use of tables as sets in implementing the
#  classical sieve algorithm for computing prime numbers.

procedure main(limit) {
	# local s, i
	local s
	local i

	/limit := 100
	s := table()
	every s[2 to limit] := 1
	every s.member(i := 2 to limit) do
		every s.delete(i + i to limit by i)
	write("In the first ", limit, " integers there are ", *s, " primes:")
	every write((!s.sort()).key)
}
