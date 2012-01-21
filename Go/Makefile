include $(GOROOT)/src/Make.inc

TARG=gorman

CGOFILES=\
	gorman.go

CGO_CFLAGS=-I$(RMANTREE)/include
CGO_LDFLAGS+=-L$(RMANTREE)/lib
CGO_LDFLAGS+=-lprman
CGO_LDFLAGS+=-Wl,-rpath $(RMANTREE)/lib

CLEANFILES+=sample
CLEANFILES+=debug.rib

include $(GOROOT)/src/Make.pkg

sample: install sample.go
	$(GC) sample.go
	$(LD) -r $(RMANTREE)/lib -o $@ sample.$O

