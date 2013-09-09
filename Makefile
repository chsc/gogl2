bindings:
	go build
	./gogl2 pulldoc -ver=2
	./gogl2 pulldoc -ver=3
	./gogl2 pulldoc -ver=4
	./gogl2 pullspec
	./gogl2 generate
	make install_bindings

install_bindings:
#	go install ./gl21
#	go install ./gl30
#	go install ./gl31
#	go install ./gl31c
#	go install ./gl32
#	go install ./gl32c
#	go install ./gl33
#	go install ./gl33c
#	go install ./gl40
#	go install ./gl41c
#	go install ./gl42
#	go install ./gl42c
#	go install ./gl43
#	go install ./gl44
#	go install ./arb
#	go install ./ext
#	go install ./ati
#	go install ./amd
#	go install ./nv
#	go install ./glx
#	go install ./wgl
