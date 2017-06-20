all: darwin windows.exe linux
	@

%:
	GOOS=$(basename $@) go build -o smartling.$@
