#!/bin/bash

export PATH=$PREFIX/bin:~/.local/bin:/usr/local/bin:$PATH
export PKG_CONFIG_PATH=~/lib/pkgconfig:$PKG_CONFIG_PATH
# export SDL_PREFIX=~/.local/usr

# ./tools/sdl_from_source.sh > /dev/null

export GOPATH=$HOME

rm -rf $GOPATH/src/github.com/andreas-jonsson/openwar
go get -u github.com/andreas-jonsson/openwar

cd $GOPATH/src/github.com/andreas-jonsson/openwar
go build openwar.go

cd tools/package/app-bundler
gtk-mac-bundler openwar.bundle

cd ../Output
zip -r ../../OpenWar.zip OpenWar.app
