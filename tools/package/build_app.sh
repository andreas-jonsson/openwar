#!/bin/bash

export PATH=$PREFIX/bin:~/.local/bin:$PATH
export PKG_CONFIG_PATH=~/.local/usr/lib/pkgconfig:$PKG_CONFIG_PATH
export SDL_PREFIX=~/.local/usr

# This step was moved to here because it takes to long to just do it all in the setup script.
jhbuild build meta-gtk-osx-core > jhbuild_log.txt 2>&1

./tools/sdl_from_source.sh > /dev/null

rm -rf $GOPATH/pkg $GOPATH/bin

go get
go build openwar.go

cd tools/package/app-bundler
gtk-mac-bundler openwar.bundle

cd ../Output
7z a -tzip ../../OpenWar.zip OpenWar.app
cd ../../
