fuzz:
	mkdir -p corpus
	cp ../testdata_source/* corpus/
	go run github.com/dvyukov/go-fuzz/go-fuzz-build
	go run github.com/dvyukov/go-fuzz/go-fuzz