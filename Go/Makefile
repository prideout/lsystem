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

torus: install torus.go
	$(GC) torus.go
	$(LD) -r $(RMANTREE)/lib -o $@ torus.$O

lsystem: install lsystem.go
	$(GC) lsystem.go

curve: install curve.go lsystem
	$(GC) curve.go
	$(LD) -r $(RMANTREE)/lib -o $@ curve.$O

pretty:
	gofmt -tabindent=false -tabwidth=4 -w=true *.go
