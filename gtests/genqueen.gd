#SRC: icon/genqueen.icn

############################################################################
#
#	File:     genqueen.icn
#
#	Subject:  Program to solve arbitrary-size n-queens problem
#
#	Author:   Peter A. Bigot
#
#	Date:     October 25, 1990
#
############################################################################
#
# This program solve the non-attacking n-queens problem for (square) boards
# of arbitrary size.  The problem consists of placing chess queens on an
# n-by-n grid such that no queen is in the same row, column, or diagonal as
# any other queen.  The output is each of the solution boards; rotations
# not considered equal.  An example of the output for n:
#
#     -----------------
#     |Q| | | | | | | |
#     -----------------
#     | | | | | | |Q| |
#     -----------------
#     | | | | |Q| | | |
#     -----------------
#     | | | | | | | |Q|
#     -----------------
#     | |Q| | | | | | |
#     -----------------
#     | | | |Q| | | | |
#     -----------------
#     | | | | | |Q| | |
#     -----------------
#     | | |Q| | | | | |
#     -----------------
#
# Usage: genqueen n
# where n is the number of rows / columns in the board.  The default for n
# is 6.
#
############################################################################

global	n                           # Number of rows/columns
global	rw                          # List of queens in each row
global	dd                          # List of queens in each down diagonal
global	ud                           # List of queens in each up diagonal

procedure main (arg) {           # Program arguments
	n := integer (\arg) | 6
	rw := list (n)
	dd := list (2*n-1)
	ud := list (2*n-1)
	solvequeen (1)
} # procedure main

# placequeen(c) -- Place a queen in every permissible position in column c.
# Suspend with each result.
procedure placequeen (c) {        # Column at which to place queen
	local r                      # Possible placement row

	every r := 1 to n do
		suspend (/rw [r] <- /dd [r+c-1] <- /ud [n+r-c] <- c)
	return fail
} # procedure placequeen

# solvequeen(c) -- Place the c'th and following column queens on the board.
# Write board if have completed it.  Suspends all viable results
procedure solvequeen (c) {        # Column for next queen placement
	if (c > n) then {
		# Have placed all required queens.  Write the board, and resume search.
		writeboard ()
		return fail
	}
	suspend placequeen (c) & solvequeen (c+1)
	return fail
} # procedure solvequeen

# writeboard() -- Write an image of the board with the queen positions
# represented by Qs.
procedure writeboard () {
	local r                        # Index over rows during print
	local c                        # Column of queen in row r
	local row                       # Depiction of row as its created

	write (repl ("--", n), "-")
	every r := 1 to n do {
		c := rw [r]
		row := repl ("| ", n) || "|"
		row [2*c] := "Q"
		write (row)
		write (repl ("--", n), "-")
	}
	write ()
} # procedure writeboard
