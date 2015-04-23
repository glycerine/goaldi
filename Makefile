#  Goaldi Makefile
#
#  Assumptions:
#	$GOPATH is set per Go documentation
#	$GOPATH/bin (first GOPATH component) is destination for built programs
#	$GOPATH/bin are in search path, as is the Go compiler
#
#	Goaldi builds itself by a bootstrapping process.
#	If you break it so badly that this doesn't work:
#		run "make clean uninstall"  (probably optional, but a good idea)
#		restore the original source code
#		run "make full"

PKG = goaldi
PROGS = $(PKG)/goaldi
# GOBIN expands in shell to {first component of $GOPATH}/bin
GOBIN = $${GOPATH%%:*}/bin

#  default action: set up, build all, run test suite, run expt.gd if present
default:	setup build doc test expt

#  quick rebuild and test
quick:		build qktest expt

#  configure Git pre-commit hook
HOOKMASTER = ./pre-commit.hook
HOOKFILE = .git/hooks/pre-commit
setup:	$(HOOKFILE)
$(HOOKFILE):	$(HOOKMASTER)
	cp $(HOOKMASTER) $(HOOKFILE)

#  build using existing translator if available
build:
	$(GOBIN)/goaldi -x -l /dev/null || $(MAKE) boot
	cd tran; $(MAKE)
	go install $(PROGS)

#  bootstrap build goaldi using stable translator binary
boot:
	cd tran; $(MAKE) boot
	go install $(PROGS)

#  full three-pass rebuild using bootstrapping from old stable front end
full:
	#
	# make an executable that embeds an old version of the front end
	#
	cd tran; $(MAKE) boot
	cd runtime; go test
	go install $(PROGS)
	cd tests; $(MAKE) quick
	#
	# make an executable embedding the latest front end, built by the old one
	#
	cd tran; $(MAKE) clean; $(MAKE) GEN=1 gtran.go
	go install $(PROGS)
	cd tests; $(MAKE) quick
	#
	# make an executable embedding the latest front end as built by itself
	#
	cd tran; $(MAKE) clean; $(MAKE) GEN=2 gtran.go
	go install $(PROGS)
	$(MAKE) doc
	cd tests; $(MAKE) quick
	#
	# looks good in quick tests; now run full suite
	#
	$(MAKE) test

#  extract stdlib documentation from the Goaldi binary
doc:	.FORCE
	cd doc; $(MAKE)

#  run Go unit tests; build and link demos; run Goaldi test suite
test:
	cd runtime; go test
	cd demos; $(MAKE) link
	cd tests; $(MAKE)

#  run a single quick test
qktest:
	cd tests; $(MAKE) quick

#  run expt.gd (presumably the test of the moment) if present
#  passes $GXOPTS to interpreter if set in environment
expt:
	test -f expt.gd && $(GOBIN)/goaldi $$GXOPTS expt.gd || :

#  install the newly built translator as the stable version for future builds
#  (be sure this is a good one or you'll lose the ability to bootstrap)
accept:
	cd tran; $(MAKE) accept

#  prepare Go source for check-in by running standard Go reformatter
format:
	go fmt *.go
	go fmt ir/*.go
	go fmt goaldi/*.go
	go fmt runtime/*.go
	go fmt extensions/*.go

#  gather together source for single-file editing; requires "bundle" util
bundle:
	@bundle `find * -name '*.go' ! -name gtran.go`

#  remove temporary and built files from source tree
#  and also subpackages built and saved in $GOPATH
clean:
	go clean $(PKG) $(PROGS)
	cd tran; $(MAKE) clean
	cd tests; $(MAKE) clean
	cd doc; $(MAKE) clean
	rm -rf $(GOBIN)/../pkg/*/goaldi

#  remove files placed elsewhere in $GOPATH
uninstall:
	go clean -i $(PKG) $(PROGS)
	rm -rf $(GOBIN)/goaldi $(GOBIN)/../pkg/*/goaldi


.FORCE:
