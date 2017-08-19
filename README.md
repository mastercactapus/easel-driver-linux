# EaselDriver for linux!

This repo contains scripts for re-packaging/compiling the easel driver for use on linux systems.
The released *AppImage* files are self-contained executables with all dependencies bundled.

An exception is **libfuse2** which is preinstalled already on most systems, but some Raspberry Pi distros may require
you to install it (e.g. raspbian-lite). This is because AppImages use fuse to mount without requiring
root.

More info here: [AppImage](http://appimage.org/)

You can download pre-built releases here: [releases](https://github.com/mastercactapus/easel-driver-linux/releases)

# Building

To build the AppImage for the current platform, simply run `make` from the root of this repo.

## Requirements

To build you must have node v4.5.0 installed and upgrade npm to v3.10.7. 
A quick way to get node installed is to use [nvm](https://github.com/creationix/nvm).
Once installed, you can run `npm install -g npm@3.10.7` to set npm to the correct version.

You must also have basic build tools (make, cmake, etc..), git, and python2 available for
node to build native extensions.

The following files must be present in the repo root:

- appimagetool-<..your-arch-here..>.AppImage
- EaselDriver-0.3.3.pkg

appimagetool is part of [AppImageKit](https://github.com/probonopd/AppImageKit). If you are targeting a
platform that doesn't have pre-built versions, you will have to clone the *AppImageKit* repo and run `build.sh`
to build the tools.

EaselDriver-0.3.3.pkg can be found via the **mac** download link from the easel setup page.


# Disclaimer

This linux easel driver is not supported by Easel/Inventables and not associated in any way. The driver itself is
just recompiled from the mac version for those in the community that prefer to use linux.

