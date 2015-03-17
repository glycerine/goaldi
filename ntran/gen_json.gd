#  gen_json.icn -- create json output from intermediate representation.

procedure json_File(irgen, flagList) {
    local p
    local flag
    local s

    flag := nil
    s := "[\n"
    while p := @irgen do {
        if \flag then s ||:= ",\n"
	flag := "true"

	s ||:= json(p, "")
    }
    s ||:= "\n]"
    return s
}

procedure json_list(p, indent) {
	local s
	local flag
	local i

	s := "["
	flag := nil
	every i := !p do {
		if \flag then {
			s ||:= ","
		}
		flag := "true"
		s ||:= "\n" || indent || "\t" || json(i, indent || "\t")
	}
	s ||:= "\n" || indent || "]"
	return s
}

procedure json_record(p, indent) {
	local s
	local i

	s := "{\n" || indent || "\t\"tag\" : " || image(type(p))
	every i := 1 to *p do {
		s ||:= ",\n" || indent || "\t"
		s ||:= image(p.type()[i])
		s ||:= " : "
		s ||:= json(p[i], indent || "\t")
	}
	s ||:= "\n" || indent || "}"
	return s
	
}

procedure json(p, indent) {
	if /p then {
		return "null"
	}

	case type(p) of {
	"ir_Tmp" | "ir_TmpLabel" | "ir_TmpClosure" : return image(p.name)
	"ir_Label" : return image(p.value)
	"ir_coordinate" : return image(p.file || ":" || p.line || ":" || p.column)
	}

	if match("record", image(p)) & type(p) ~== "string" then {
		return json_record(p, indent)
	} else {
		case type(p) of {
			list: return json_list(p, indent)
			set: return json_list(p, indent)
			string: return json_image(string(p))
			number: return image(string(p))
			default: throw("bad type for json", p)
		}
	}
}

procedure json_image(s) {
	local t
	static mapping
	\mapping | {
		mapping := table()
		mapping["\x00"] := "\\u0000"
		mapping["\x01"] := "\\u0001"
		mapping["\x02"] := "\\u0002"
		mapping["\x03"] := "\\u0003"
		mapping["\x04"] := "\\u0004"
		mapping["\x05"] := "\\u0005"
		mapping["\x06"] := "\\u0006"
		mapping["\x07"] := "\\u0007"
		mapping["\b"] := "\\b"
		mapping["\t"] := "\\t"
		mapping["\n"] := "\\n"
		mapping["\v"] := "\\u000b"
		mapping["\f"] := "\\f"
		mapping["\r"] := "\\r"
		mapping["\x0e"] := "\\u000e"
		mapping["\x0f"] := "\\u000f"
		mapping["\x10"] := "\\u0010"
		mapping["\x11"] := "\\u0011"
		mapping["\x12"] := "\\u0012"
		mapping["\x13"] := "\\u0013"
		mapping["\x14"] := "\\u0014"
		mapping["\x15"] := "\\u0015"
		mapping["\x16"] := "\\u0016"
		mapping["\x17"] := "\\u0017"
		mapping["\x18"] := "\\u0018"
		mapping["\x19"] := "\\u0019"
		mapping["\x1a"] := "\\u001a"
		mapping["\e"] := "\\u001b"
		mapping["\x1c"] := "\\u001c"
		mapping["\x1d"] := "\\u001d"
		mapping["\x1e"] := "\\u001e"
		mapping["\x1f"] := "\\u001f"
		mapping[" "] := " "
		mapping["!"] := "!"
		mapping["\""] := "\\\""
		mapping["#"] := "#"
		mapping["$"] := "$"
		mapping["%"] := "%"
		mapping["&"] := "&"
		mapping["'"] := "'"
		mapping["("] := "("
		mapping[")"] := ")"
		mapping["*"] := "*"
		mapping["+"] := "+"
		mapping[","] := ","
		mapping["-"] := "-"
		mapping["."] := "."
		mapping["/"] := "/"
		mapping["0"] := "0"
		mapping["1"] := "1"
		mapping["2"] := "2"
		mapping["3"] := "3"
		mapping["4"] := "4"
		mapping["5"] := "5"
		mapping["6"] := "6"
		mapping["7"] := "7"
		mapping["8"] := "8"
		mapping["9"] := "9"
		mapping[":"] := ":"
		mapping[";"] := ";"
		mapping["<"] := "<"
		mapping["="] := "="
		mapping[">"] := ">"
		mapping["?"] := "?"
		mapping["@"] := "@"
		mapping["A"] := "A"
		mapping["B"] := "B"
		mapping["C"] := "C"
		mapping["D"] := "D"
		mapping["E"] := "E"
		mapping["F"] := "F"
		mapping["G"] := "G"
		mapping["H"] := "H"
		mapping["I"] := "I"
		mapping["J"] := "J"
		mapping["K"] := "K"
		mapping["L"] := "L"
		mapping["M"] := "M"
		mapping["N"] := "N"
		mapping["O"] := "O"
		mapping["P"] := "P"
		mapping["Q"] := "Q"
		mapping["R"] := "R"
		mapping["S"] := "S"
		mapping["T"] := "T"
		mapping["U"] := "U"
		mapping["V"] := "V"
		mapping["W"] := "W"
		mapping["X"] := "X"
		mapping["Y"] := "Y"
		mapping["Z"] := "Z"
		mapping["["] := "["
		mapping["\\"] := "\\\\"
		mapping["]"] := "]"
		mapping["^"] := "^"
		mapping["_"] := "_"
		mapping["`"] := "`"
		mapping["a"] := "a"
		mapping["b"] := "b"
		mapping["c"] := "c"
		mapping["d"] := "d"
		mapping["e"] := "e"
		mapping["f"] := "f"
		mapping["g"] := "g"
		mapping["h"] := "h"
		mapping["i"] := "i"
		mapping["j"] := "j"
		mapping["k"] := "k"
		mapping["l"] := "l"
		mapping["m"] := "m"
		mapping["n"] := "n"
		mapping["o"] := "o"
		mapping["p"] := "p"
		mapping["q"] := "q"
		mapping["r"] := "r"
		mapping["s"] := "s"
		mapping["t"] := "t"
		mapping["u"] := "u"
		mapping["v"] := "v"
		mapping["w"] := "w"
		mapping["x"] := "x"
		mapping["y"] := "y"
		mapping["z"] := "z"
		mapping["{"] := "{"
		mapping["|"] := "|"
		mapping["}"] := "}"
		mapping["~"] := "~"
		mapping["\d"] := "\\u007f"
		# assumption: chars beyond here are already UTF-8 encoded
		mapping["\x80"] := "\x80"
		mapping["\x81"] := "\x81"
		mapping["\x82"] := "\x82"
		mapping["\x83"] := "\x83"
		mapping["\x84"] := "\x84"
		mapping["\x85"] := "\x85"
		mapping["\x86"] := "\x86"
		mapping["\x87"] := "\x87"
		mapping["\x88"] := "\x88"
		mapping["\x89"] := "\x89"
		mapping["\x8a"] := "\x8a"
		mapping["\x8b"] := "\x8b"
		mapping["\x8c"] := "\x8c"
		mapping["\x8d"] := "\x8d"
		mapping["\x8e"] := "\x8e"
		mapping["\x8f"] := "\x8f"
		mapping["\x90"] := "\x90"
		mapping["\x91"] := "\x91"
		mapping["\x92"] := "\x92"
		mapping["\x93"] := "\x93"
		mapping["\x94"] := "\x94"
		mapping["\x95"] := "\x95"
		mapping["\x96"] := "\x96"
		mapping["\x97"] := "\x97"
		mapping["\x98"] := "\x98"
		mapping["\x99"] := "\x99"
		mapping["\x9a"] := "\x9a"
		mapping["\x9b"] := "\x9b"
		mapping["\x9c"] := "\x9c"
		mapping["\x9d"] := "\x9d"
		mapping["\x9e"] := "\x9e"
		mapping["\x9f"] := "\x9f"
		mapping["\xa0"] := "\xa0"
		mapping["\xa1"] := "\xa1"
		mapping["\xa2"] := "\xa2"
		mapping["\xa3"] := "\xa3"
		mapping["\xa4"] := "\xa4"
		mapping["\xa5"] := "\xa5"
		mapping["\xa6"] := "\xa6"
		mapping["\xa7"] := "\xa7"
		mapping["\xa8"] := "\xa8"
		mapping["\xa9"] := "\xa9"
		mapping["\xaa"] := "\xaa"
		mapping["\xab"] := "\xab"
		mapping["\xac"] := "\xac"
		mapping["\xad"] := "\xad"
		mapping["\xae"] := "\xae"
		mapping["\xaf"] := "\xaf"
		mapping["\xb0"] := "\xb0"
		mapping["\xb1"] := "\xb1"
		mapping["\xb2"] := "\xb2"
		mapping["\xb3"] := "\xb3"
		mapping["\xb4"] := "\xb4"
		mapping["\xb5"] := "\xb5"
		mapping["\xb6"] := "\xb6"
		mapping["\xb7"] := "\xb7"
		mapping["\xb8"] := "\xb8"
		mapping["\xb9"] := "\xb9"
		mapping["\xba"] := "\xba"
		mapping["\xbb"] := "\xbb"
		mapping["\xbc"] := "\xbc"
		mapping["\xbd"] := "\xbd"
		mapping["\xbe"] := "\xbe"
		mapping["\xbf"] := "\xbf"
		mapping["\xc0"] := "\xc0"
		mapping["\xc1"] := "\xc1"
		mapping["\xc2"] := "\xc2"
		mapping["\xc3"] := "\xc3"
		mapping["\xc4"] := "\xc4"
		mapping["\xc5"] := "\xc5"
		mapping["\xc6"] := "\xc6"
		mapping["\xc7"] := "\xc7"
		mapping["\xc8"] := "\xc8"
		mapping["\xc9"] := "\xc9"
		mapping["\xca"] := "\xca"
		mapping["\xcb"] := "\xcb"
		mapping["\xcc"] := "\xcc"
		mapping["\xcd"] := "\xcd"
		mapping["\xce"] := "\xce"
		mapping["\xcf"] := "\xcf"
		mapping["\xd0"] := "\xd0"
		mapping["\xd1"] := "\xd1"
		mapping["\xd2"] := "\xd2"
		mapping["\xd3"] := "\xd3"
		mapping["\xd4"] := "\xd4"
		mapping["\xd5"] := "\xd5"
		mapping["\xd6"] := "\xd6"
		mapping["\xd7"] := "\xd7"
		mapping["\xd8"] := "\xd8"
		mapping["\xd9"] := "\xd9"
		mapping["\xda"] := "\xda"
		mapping["\xdb"] := "\xdb"
		mapping["\xdc"] := "\xdc"
		mapping["\xdd"] := "\xdd"
		mapping["\xde"] := "\xde"
		mapping["\xdf"] := "\xdf"
		mapping["\xe0"] := "\xe0"
		mapping["\xe1"] := "\xe1"
		mapping["\xe2"] := "\xe2"
		mapping["\xe3"] := "\xe3"
		mapping["\xe4"] := "\xe4"
		mapping["\xe5"] := "\xe5"
		mapping["\xe6"] := "\xe6"
		mapping["\xe7"] := "\xe7"
		mapping["\xe8"] := "\xe8"
		mapping["\xe9"] := "\xe9"
		mapping["\xea"] := "\xea"
		mapping["\xeb"] := "\xeb"
		mapping["\xec"] := "\xec"
		mapping["\xed"] := "\xed"
		mapping["\xee"] := "\xee"
		mapping["\xef"] := "\xef"
		mapping["\xf0"] := "\xf0"
		mapping["\xf1"] := "\xf1"
		mapping["\xf2"] := "\xf2"
		mapping["\xf3"] := "\xf3"
		mapping["\xf4"] := "\xf4"
		mapping["\xf5"] := "\xf5"
		mapping["\xf6"] := "\xf6"
		mapping["\xf7"] := "\xf7"
		mapping["\xf8"] := "\xf8"
		mapping["\xf9"] := "\xf9"
		mapping["\xfa"] := "\xfa"
		mapping["\xfb"] := "\xfb"
		mapping["\xfc"] := "\xfc"
		mapping["\xfd"] := "\xfd"
		mapping["\xfe"] := "\xfe"
		mapping["\xff"] := "\xff"
	}
	t := ""
	every t ||:= mapping[!s]
	t := "\"" || t || "\""
	return t
}