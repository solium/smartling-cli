MAINTAINER = Stanislav Seletskiy <s.seletskiy@gmail.com>

all: clean get build
	@

build: darwin windows.exe linux
	@

get:
	go get

clean:
	rm -rf bin pkg
	mkdir bin

_CONTROL = echo >> pkg/DEBIAN/control

deb: get linux
	rm -rf pkg
	mkdir -p pkg/{usr/bin,DEBIAN}
	cp bin/smartling.linux pkg/usr/bin/smartling
	$(eval VERSION ?= \
		$(shell git rev-parse --short HEAD).$(shell git rev-list --count HEAD))
	$(_CONTROL) "Package: smartling"
	$(_CONTROL) "Version: $(VERSION)"
	$(_CONTROL) "Architecture: all"
	$(_CONTROL) "Section: unknown"
	$(_CONTROL) "Priority: extra"
	$(_CONTROL) "Maintainer: $(MAINTAINER)"
	$(_CONTROL) "Homepage: https://github.com/Smartling/smartling-cli"
	$(_CONTROL) "Description: CLI for Smartling Platform"
	dpkg -b pkg smartling-$(VERSION)_all.deb

%:
	GOOS=$(basename $@) go build -o bin/smartling.$@
