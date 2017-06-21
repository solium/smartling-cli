all: clean get build

get:
	go get

build: darwin windows.exe linux
	@

%:
	GOOS=$(basename $@) go build -o bin/smartling.$@

clean:
	rm -rf ./bin
	mkdir ./bin
