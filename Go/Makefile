curve: curve.go lsystem.go stack.go
	make -f Makefile.curve

image: curve
	./curve Ribbon.xml single render

rib: curve
	./curve Ribbon.xml single ribgen

lsys: curve
	./curve Ribbon.xml single lsys

torus: torus.go
	make -f Makefile.torus

pretty:
	gofmt -tabs=false -tabwidth=4 -w=true *.go

clean:
	make clean -f Makefile.curve
	make clean -f Makefile.torus
