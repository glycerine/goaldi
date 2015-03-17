#  gtran.gd -- main program for Goaldi front-end translator
#  (quick-and-dirty experimental version)

procedure main(args[]) {
	every ^fname := !args do {
		^fbase := if fname[-3:0]==".gd" then fname[1:-3] else fname
		^pipeline := create !open(fbase || ".gd")
		pipeline := create lex(pipeline, fbase || ".gd")
		pipeline := create tee(pipeline, fbase || ".tk")
		# pipeline := create stdout(pipeline)
		pipeline := create sink(pipeline)
		@pipeline	# wait for processes to finish and close
	}
}

#  pipeline component to copy its contents to a file
procedure tee(src, fname) {
	^f := open(fname, "w")
	while ^v := @src do {
		f.write(v.image())
		suspend v
	}
	f.close()
}

#  (terminal) pipeline component to write stream stdout
procedure stdout(src) {
	while write(image(@src))
}

#  (terminal) pipeline component to toss everything into a black hole
procedure sink(src) {
	while @src
}