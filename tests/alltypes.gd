#SRC: Goaldi original
#
#  Show examples of all types, presented multiple ways,
#  and test sorting and various type and value methods

record Example(	# one example for display
	value,		# example value
	type,		# value.type() (for sorting by type)
	gtype,		# corresponding global type if any
)

record Point(x,y)								# a simple illustrative record
procedure Point.dist() { return hypot(self.x, self.y) }	# and a method for it
record Circle extends Point(r)

global allvals	# set of all example values
global ttable	# table of distinct types
global tlist	# list of distinct types

procedure main() {

	# make a list of examples with associated global type values
	^E := []
	allvals := set()
	ttable := table()
	add(E, nil, niltype)
	add(E, type(), type)
	add(E, 17, number)
	add(E, %pi, number)
	add(E, 6.02214129e23, number)
	add(E, "abcd", string)
	add(E, %stdin, file)
	add(E, channel(3), channel)
	add(E, Point, constructor)
	add(E, ^P := Point(7,5), Point)
	add(E, P.dist, methodvalue)
	add(E, Circle(7,5,2), Circle)
	add(E, main, proctype)
	add(E, ^L := [2,3,5,7,11], list)
	add(E, ^S := set([4,7,1]), set)
	add(E, ^T := table(){"Fe":"Iron","Au":"Gold"}, table)
	add(E, !T.sort())	# table element
	add(E, tuple(w:6,h:4))
	add(E, duration(3600+120+3), external)
	write("Example set: ", image(allvals))

	# define alternate acceptable sprintf representations we might see
	^subst := table() {
		"map[Au:Gold Fe:Iron]" : "map[Fe:Iron Au:Gold]",
		"map[1:true 4:true 7:true]" : "map[4:true 7:true 1:true]",
		"map[1:true 7:true 4:true]" : "map[4:true 7:true 1:true]",
		"map[4:true 1:true 7:true]" : "map[4:true 7:true 1:true]",
		"map[7:true 1:true 4:true]" : "map[4:true 7:true 1:true]",
		"map[7:true 4:true 1:true]" : "map[4:true 7:true 1:true]",
	}

	# show values various ways, checking universal methods in the process
	write()
	write("Examples sorted by value, showing presentation options:")
	E := E.sort(Example["value"])
	write()
	^format := "%-4s %-15s %-30s %s\n"
	printf(format, "ch", "x.string()", "x.image()", "printf(\"%v\")")
	printf(format, "--", "----------", "---------", "------------")
	every ^x := !E do {
		^v := x.value
		^s := check(string, v, v.string())
		^i := check(image, v, v.image())
		^t := check(type, v, v.type())
		^f := sprintf("%v", v)
		if f[1+:2] == "0x" then		# if hex address
			f := "0xXXXXXX"			# hide actual value for reproducibility
		f := \subst[f]				# substitute alternate for reproducibility
		printf(format, t.char(), s, i, f)
	}

	# make list of distinct types for instanceof testing
	tlist := [: (!ttable.sort()).key :]

	write()
	write("Examples sorted by type, showing type information:")
	# n.b. stable sort keeps ordering reproducible within type
	E := E.sort(Example["type"])
	write()
	format := "%-2s %2s  %-14s %-12s %-13s  %-4s %s"
	printf(format,
		"c", "*t", "x.type()", "t.name()", "global", "t[1]", "  instanceof\n")
	printf(format,
		"-", "--", "--------", "--------", "------", "----", "  ----------\n")
	every x := !E do {
		^v := x.value
		^t := x.type
		^t1 := t[1] | "-"
		^bt := !t | "-"
		if t1 ~=== bt then
			write("MISMATCH: t[1] / !t: ", image(t1), " ~=== ", image(bt))
		^n := string(*t)
		printf(format, t.char(), n, t, t.name(), t===x.gtype | "", t1, "")
		every t := !tlist do
			if v.instanceof(t) then writes("  ", t)
		write()
	}
	write()
}

procedure add(E, v, g) {			#: add global type and sample value
	^t := type(v)
	ttable[t] := t		# register example of this type
	^n := *allvals
	allvals.put(v)		# validate usability of v as a set member
	allvals.put(v)		# make sure it only gets added once
	^i := *allvals - n
	if i ~= 1 then {
		write("ERROR: *S increased by ", i, " after twice adding ", image(v))
	}
	return E.put(Example(v, t, g))	# return list with new Example{} added
}

procedure check(p, x, s) {			#: validate p(x) === s
	^t := p(x)
	if t ~=== s then
		write("MISMATCH: ", p, "(", x, ")===", t, " ~=== ", s)
	return s
}
