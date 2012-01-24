include $(GOROOT)/src/Make.inc

TARG=vmath

GOFILES=\
	m4.go \
	v3.go \
	p3.go \
	
include $(GOROOT)/src/Make.pkg

pretty:
	gofmt -tabs=false -tabwidth=4 -w=true *.go
