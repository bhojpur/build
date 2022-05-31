all:
	buildc2go foobar.yml

clean:
	rm -f foobar/cgo_helpers.go foobar/cgo_helpers.h foobar/cgo_helpers.c
	rm -f foobar/const.go foobar/doc.go foobar/types.go
	rm -f foobar/foobar.go

test:
	cd foobar && go build
	