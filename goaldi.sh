#!/bin/sh
#
#	goaldi [options] file.gd... [--] [arg...] -- compile and run Goaldi program
#
#	To see options, run with no arguments.
#	This script assumes that gtran and gexec are in the search path.

FLAGS=acdNltvADEPTX

#  define the usage abort function
usage() {
	exec >&2
	cat <<==EOF==
Usage: $0 [-$FLAGS] file.gd... [--] [arg...]
  -N  no optimization
  -X  use experimental front end
  -c  compile only, producing IR on file.gir
  -a  compile only, producing IR on file.gir and assembly listing on file.gia
  -d  compile only, producing Dot directives on file.dot
==EOF==
	# add option descriptions from back end (gexec)
	gexec -.!! -? 2>&1 | sed -n -e '/-.!!/d' -e 's/=false: /  /p'
	exit 1
}

#  define the translation function
tran() {
	if [ -z "$XFLAG" ]; then
		# standard translator
		gtran cat $1 : yylex : parse : ast2ir $OPT : $FMT : stdout
	else
		# experimental translator
		ntran $NFLAG $1
	fi
}

#  process options
XOPTS=
WHAT=x
FMT=json_File
NFLAG=
XFLAG=
OPT=": optim -O"
while getopts $FLAGS C; do
	case $C in
	a)			WHAT=$C;;
	c)			WHAT=$C;;
	d)			WHAT=$C; FMT=dot_File;;
	N)			NFLAG=-$C; OPT="";;
	X)			XFLAG=-$C;;
	[ltvADEPT])	XOPTS="$XOPTS -$C";;
	?)			usage;;
	esac
done

if [ "$WHAT" = d -a -n "$XFLAG" ]; then
	echo 1>&2 "$0: -d and -X cannot be combined"
	exit 1
fi

#  collect source file names:
#  the first argument, always, plus any following that end in ".gd"
shift $(($OPTIND - 1))		# remove flag arguments
test $# -lt 1 && usage		# require at least one file argument
SRCS=$1						# save that argument
shift						# and remove from execution parameters

while [ "$1" != "${1%.gd}" ]; do	# while name ends in .gd
	SRCS="$SRCS $1"				# add to list
	shift						# and remove from execution parameters
done

#  remove a "--" separator argument if present
if [ "$1" = "--" ]; then
	shift
fi

#  make scratch directory for temporary files, and arrange its deletion
SCR=/tmp/goaldi.$$
trap 'X=$?; rm -rf $SCR; exit $X' 0 1 2 15
mkdir $SCR

#  compile the source files
export COEXPSIZE=300000
OBJS=
QUIT=:
for F in $SRCS; do
	B=${F%.*}
	case $WHAT in
		a)	# -a: produce file.gir and file.gia
			tran $F >$B.gir && gexec -.!! $XOPTS -l -A $B.gir >$B.gia
			QUIT="exit $?"
			;;
		c)	# -c or nothing: produce file.gir
			tran $F >$B.gir
			QUIT="exit $?"
			;;
		d)	# -d: produce file.dot
			tran $F >$B.dot
			QUIT="exit $?"
			;;
		x)	# no flag: produce temporary file.gir for later execution
			O=$SCR/${B##*/}.gir
			OBJS="$OBJS $O"
			tran $F >$O || QUIT="exit 1"
			;;
	esac
done

$QUIT	# exit if nothing more to do, or if errors in compilation

# execute compiled files
gexec -.!! $XOPTS $SCR/*.gir -- "$@"
