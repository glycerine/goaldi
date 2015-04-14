#SRC: icon/queens.icn

############################################################################
#
#	File:     queens.icn
#
#	Subject:  Program to generate solutions to the n-queens problem
#
#	Author:   Stephen B. Wampler
#
#	Date:     June 10, 1988
#
############################################################################
#
#   This file is in the public domain.
#
############################################################################
#
#     This program displays the solutions to the non-attacking n-
#  queens problem: the ways in which n queens can be placed on an
#  n-by-n chessboard so that no queen can attack another. A positive
#  integer can be given as a command line argument to specify the
#  number of queens. For example,
#
#          iconx queens -n8
#
#  displays the solutions for 8 queens on an 8-by-8 chessboard.  The
#  default value in the absence of an argument is 6.  One solution
#  for six queens is:
#
#         -------------------------
#         |   | Q |   |   |   |   |
#         -------------------------
#         |   |   |   | Q |   |   |
#         -------------------------
#         |   |   |   |   |   | Q |
#         -------------------------
#         | Q |   |   |   |   |   |
#         -------------------------
#         |   |   | Q |   |   |   |
#         -------------------------
#         |   |   |   |   | Q |   |
#         -------------------------
#
#  Comments: There are many approaches to programming solutions to
#  the n-queens problem.  This program is worth reading for
#  its programming techniques.
#
############################################################################

global n
global solution

procedure main(arg) {
	# local i, opts
	local i
	local opts

	n := integer(\arg) | 6
	if n <= 0 then stop("-n needs a positive numeric parameter")

	solution := list(n)		# ... and a list of column solutions
	write(n,"-Queens:")
	every q(1)			# start by placing queen in first column
}

# q(c) - place a queen in column c.
#
procedure q(c) {
	local r
	# static up, down, rows
	static up
	static down
	static rows

	if /rows := list(n,0) then {
		/up := list(2*n-1,0)
		/down := list(2*n-1,0)
	}
	every 0 = rows[r := 1 to n] = up[n+r-c] = down[r+c-1] &
		rows[r] <- up[n+r-c] <- down[r+c-1] <- 1        do {
			solution[c] := r	# record placement.
			if c = n then show() else q(c + 1)		# try to place next queen.
		}
}

# show the solution on a chess board.
#
procedure show() {
	# static count, line, border
	static count
	static line
	static border

	if /count := 0 then {
		line := repl("|   ",n) || "|"
		border := repl("----",n) || "-"
	}
	write("solution: ", count+:=1)
	write("  ", border)
	every line[4*(!solution - 1) + 3] <- "Q" do {
		write("  ", line)
		write("  ", border)
	}
	write()
}