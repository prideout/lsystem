clean:
	make clean -f Makefile.curve
	make clean -f Makefile.torus

curve: curve.go
	make -f Makefile.curve
    
torus: torus.go
	make -f Makefile.torus
	
pretty:
	gofmt -tabs=false -tabwidth=4 -w=true *.go
