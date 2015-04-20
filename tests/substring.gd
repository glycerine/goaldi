#SRC: icon/substring.icn

# string subscripting test

#  note that Goaldi follows Jcon (not Icon) in failing on wraparound for [+:]

procedure main() {
	# local i, j, k, s, t
	local i
	local j
	local k
	local s
	local t

	s := "abcde"
	t := "ABCDE"
	write("A. ", !s)
	every write("B. ", !s)
	every i := 0 to 6 do write("C ", i, ". ", s[i] | "--")
	every i := 0 to -6 by -1 do write("D ", i, ". ", s[i] | "--")
	every i := -10 to 10 do write("E ", i, ". ", s[3:i] | "--")
	every i := -10 to 5 do write("F ", i, ". ", s[3+:i] | "--") #some SHOULD fail
	every i := -5 to 10 do write("G ", i, ". ", s[3-:i] | "--") #some SHOULD fail

	!s := "X"
	write("H. ", s)
	every !s := "Y"
	write("I. ", s)

	every i := -6 to 6 do {
		s := "abcde"
		if s[i] := t[i] then {
			write("J ", i, ". ", s)
		} else {
			write("J ", i, ". --")
		}
	}

	every i := 1 to 6 do {
		every j := 1 to 6 do {
			s := "abcde"
			writes("K ", i, " ", j, ". ")
			if s[i:j] := "(*)" then {
				write(s)
			} else {
				write(s, " [failed]")
			}
		}
	}

	every i := 1 to 6 do {
		every j := 1 to 6 do {
			every k := 1 to 6 do {
				s := "abcde"
				writes("L ", i, " ", j, " ", k, ". ")
				if s[i:j][k:2] := "(*)" then {
					write(s)
				} else {
					write(s, " [failed]")
				}
			}
		}
	}

	s := "abcde"
	every !s <- "-" do write("M ", s)
	every s [1 to 5] <- "-" do write("N ", s)
	every s [(-5 to 6) +: 0] <- "--" do write("O ", s)

	s := "abcde"
	every s[2:4] := !"123" do write("P ", s)
	s := "fghij"
	every s[2:4] := !"456" do { write("Q ", s); s := "klmno" }

	s := "3♠4♥2♦4♣"
	write("R1: ", image(s))				# ascii and non-ascii
	write("R2: ", image(s[2+:4]))		# both, in substring
	write("R3: ", image(s[1]))			# ascii only ("3")
	write("R4: ", image(number(s[1])))	# so this should work
	write("R5: ", s[1] + s[3] + s[5] + s[7])	# and this too
	t := s[1:3] || " " || s[3:5] || " " || s[5:7] || " " || s[7:9]
	write("R6: ", image(t))
	t := s[1] || t[4] || t[7] || t[10]
	write("R7: ", image(t))
	write("R8: ", image(number(t)))

}
