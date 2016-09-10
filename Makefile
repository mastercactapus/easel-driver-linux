.PHONY: clean

VERSION=0.2.6
BUILD=1

ARCH=$(shell uname -m)
AIA=./AppImageAssistant_6-$(ARCH).AppImage
NATIVE_EXT=build/iris-lib/node_modules/serialport/build/Release/serialport.node
OUT=EaselDriver-$(VERSION)-$(ARCH)-$(BUILD).AppImage

ifneq ($(shell node --version),v4.5.0)
$(error node version must be v4.5.0 but got $(shell node --version))
else ifneq ($(shell npm --version),3.10.7)
$(error npm version must be 3.10.7 but got $(shell npm --version))
endif


$(OUT): build/AppRun build/easel.svg build/easel-driver.desktop build/node $(NATIVE_EXT)
	-rm $@
	$(AIA) build $@

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

build/iris-lib/iris.js: EaselDriver-0.2.6.pkg build
	rm -rf tempdir
	mkdir tempdir
	7z x -otempdir EaselDriver-0.2.6.pkg
	mkdir -p build/iris-lib
	(cd build/iris-lib && gunzip <../../tempdir/IrisLib-0.2.6.pkg/Payload | cpio -i)
	rm -rf tempdir
build:
	mkdir -p build

clean:
	rm -rf tempdir build $(OUT)

