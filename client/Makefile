include $(GOROOT)/src/Make.inc

TARG=autohttperf
GOFILES=\
		client.go \
		parse.go \
		types.go \
		utils.go \

include $(GOROOT)/src/Make.cmd

## This exists solely to allow 'gotest' to work properly with this command,
## since by default it only works with packages

testpackage: _test/$(TARG).a
testpackage-clean:
		rm -f _test/$(TARG).a _gotest_.$O

_gotest_.$O: $(GOFILES) $(GOTESTFILES) $(PREREQ)
		@mkdir -p _test/ofiles/$(dir)
		@for FILE in $(GOFILES) $(GOTESTFILES); do \
                cat $$FILE | sed 's/^package .*/package $(TARG)/' > \
                        _test/ofiles/$$FILE; \
        done
		cd _test/ofiles && $(GC) -o ../../$@ $(GOFILES) $(GOTESTFILES)

_test/$(TARG).a: _gotest_.$O $(OFILES)
		@mkdir -p _test/$(dir)
		rm -f _test/$(TARG).a
		gopack grc $@ _gotest_.$O $(OFILES)

importpath:
		@echo $(TARG)

clean: testclean

testclean:
		rm -rf _test*

