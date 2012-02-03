include $(GOROOT)/src/Make.inc

TARG=rman

CGOFILES=\
	rman.go

CGO_CFLAGS=-I$(RMANTREE)/include
CGO_LDFLAGS+=-L$(RMANTREE)/lib
CGO_LDFLAGS+=-lprman
CGO_LDFLAGS+=-Wl,-rpath $(RMANTREE)/lib

include $(GOROOT)/src/Make.pkg

pretty:
	gofmt -tabs=false -tabwidth=4 -w=true *.go

rpath:
	readelf -d *.o
