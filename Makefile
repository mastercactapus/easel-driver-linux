.PHONY: clean AppDir

VERSION=0.3.3
BUILD=1

ARCH=$(shell uname -m)
AIA=./appimagetool-$(ARCH).AppImage
NATIVE_EXT=build/iris-lib/node_modules/serialport/build/Release/serialport.node
OUT=EaselDriver-$(VERSION)-$(BUILD)-$(ARCH).AppImage



$(OUT): build/AppRun build/easel.svg build/easel-driver.desktop build/node $(NATIVE_EXT)
	-rm $@
	$(AIA) build $@

AppDir: $(NATIVE_EXT) build/AppRun build/easel.svg build/easel-driver.desktop

build/AppRun: AppRun build
	cp AppRun build/
	chmod +x build/AppRun

build/easel.svg: easel.svg build
	cp easel.svg build/

build/easel-driver.desktop: easel-driver.desktop build
	cp easel-driver.desktop build/

$(NATIVE_EXT): build/iris-lib/iris.js build
	(cd build/iris-lib && rm -rf node_modules && npm install && rm -rf node_modules/.bin node_modules/serialport/node_modules/.bin)
	touch $(NATIVE_EXT)

build/node: build
	cp $(shell which node) build/

build/iris-lib/iris.js: EaselDriver-$(VERSION).pkg build
	rm -rf tempdir
	mkdir tempdir
	7z x -otempdir EaselDriver-$(VERSION).pkg
	mkdir -p build/iris-lib
	(cd build/iris-lib && gunzip <../../tempdir/IrisLib-$(VERSION).pkg/Payload | cpio -i)
	rm -rf tempdir
build:
	mkdir -p build

clean:
	rm -rf tempdir build $(OUT)
