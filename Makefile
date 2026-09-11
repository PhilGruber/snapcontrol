
build:
	go build -ldflags "-X main.AppVersion=$(VERSION)"


deb: build
	rm -rf deb
	mkdir -p deb/snapcontrol/usr/bin
	mkdir deb/snapcontrol/DEBIAN
	cp deb.control deb/snapcontrol/DEBIAN/control
	sed -i'' "s/__version__/$(VERSION)/" deb/snapcontrol/DEBIAN/control
	sed -i'' "s/__arch__/$(ARCH)/" deb/snapcontrol/DEBIAN/control
	cat deb/snapcontrol/DEBIAN/control
	cp snapcontrol deb/snapcontrol/usr/bin
	dpkg-deb -Zgzip --root-owner-group --build deb/snapcontrol

