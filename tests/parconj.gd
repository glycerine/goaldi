#SRC: goaldi original
#
#  test parallel conjunction (e1 && e2)

procedure main() {
	local i
	local j
	local k
	every i := toby(1,3) && j := toby(4,5) do
		write(": ", i, j)
	every i := toby(1,3) && j := toby(4,6) do
		write(": ", i, j)
	every i := toby(1,3) && j := toby(4,7) do
		write(": ", i, j)
	every i := toby(1,3) && j := toby(4,5) && k := toby(6, 9) do
		write(": ", i, j, k)
	write()
	1 && 0	# too simple -- this used to panic
	write("done")
}

procedure toby(i, j) {
	every local v := i to j do {
		writes(v, " ")
		suspend(v)
	}
}
