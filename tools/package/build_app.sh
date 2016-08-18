#!/bin/bash

DIR=OpenWar.app

if [ -d "$DIR" ]; then rm -rf "$DIR"; fi
cp -r "tools/package/OpenWar.app" $DIR

mkdir $DIR/Contents/MacOS
cp openwar $DIR/Contents/MacOS
cp -r data $DIR/Contents/MacOS

dylibbundler -od -b -x $DIR/Contents/MacOS/openwar -d $DIR/Contents/libs

7z a -tzip OpenWar.zip OpenWar.app
