#SRC: icon/table.icn
#
# Table test

procedure main() {
	# local k, kv, x, y
	local k
	local kv
	local x
	local y

	x := table()
	tdump("initial", x)
	writes("should fail ", image(?x))
	# portable with 0 or 1 entries:
	every writes(" ", ">>" | (!x).value | "\n")
	x[2] := 3;
	every writes(" ", ">>" | (!x).value | image((?x).value) | "\n")
	x[4] := 7;
	x["a"] := "A";
	tdump("+2+4+a", x)

	every kv := !x do x[kv.key] := 88
	tdump("!x=88", x)

	every x[(!x).key] := 99
	tdump("[all]=99", x)

	every k := (!x).key do
		x[k] := k
	tdump("x[k]=k", x)

	/x[1] | write("/1")
	\x[2] | write("\\2")

	x := table()
	if x.member() then write("NIL IS MEMBER")
	x[nil] := nil			| write("failed 0")
	x[1] := nil				| write("failed 1")
	x[3] := nil				| write("failed 3")
	x[5] := 55				| write("failed 5")
	(x[6] := 66 & x[7] := 77) | write("failed 67")
	x[nil] := "nil"			| write("failed n")
	if not x.member() then write("NIL IS NOT MEMBER")
	tdump("insert", x)
	x.delete(nil)			| write("failed dn")
	x.delete(3)				| write("failed d3")
	x.delete(7,1)			| write("failed d71")
	tdump("delete", x)

	x := table(0)
	write(x[47])
	tdump("t0", x)
	x[nil] := nil			| write("failed 0")
	x[1] := nil				| write("failed 1")
	x[3] := nil				| write("failed 3")
	x[5] := 55				| write("failed 5")
	(x[6] := 66 & x[7] := 77) | write("failed 67")
	x[nil] := "nil"			| write("failed n")
	tdump("t0i", x)
	x.delete(nil)			| write("failed dn")
	x.delete(3)				| write("failed d3")
	x.delete(7).delete(1)	| write("failed d71")
	tdump("t0d", x)

	write()
	x := table()
	every x[3] <- 19		# should insert key but revert to default value
	every kv := !x do
		write("{",kv.key,",",kv.value,"}")

	x := table()
	every k := 0 to 4 do
		x[k] := k + 10
	y := copy(x)
	every x[(!x).key] +:= 20
	every y[(!y).key] +:= 40
	tdump("30s", x)
	tdump("50s", y)

}


#  dump a table, assuming that keys are drawn from: nil, 0 - 9, "a" - "e"
#
#  also checks member()

procedure tdump(label, T) {
	local x

	printf("%10s :%2.0f :", label, *T)
	every x := nil | (0 to 9) | !"abcde" do
		if x === ((!T).key) then {
			writes(" [", image(x), "]", image(T[x]))
			T.member(x) | writes(":NONMEMBER")
		} else {
			T.member(x) & writes(" MEMBER:", image(x))
		}
	write()
	return
}
