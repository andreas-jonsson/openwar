#!/bin/bash

if [ -d OpenWar.app ]; then rm -rf OpenWar.app; fi
cp -r "tools/package/OpenWar.app" OpenWar.app

cp openwar OpenWar.app/Contents/MacOS
dylibbundler -od -b -x OpenWar.app/Contents/MacOS/openwar -d OpenWar.app/Contents/libs
dylibbundler -b -x OpenWar.app/Contents/libs/libSDL2-2.0.0.dylib -d OpenWar.app/Contents/libs
dylibbundler -b -x OpenWar.app/Contents/libs/libSDL2_mixer-2.0.0.dylib -d OpenWar.app/Contents/libs

zip -rq OpenWar.zip OpenWar.app
rm -rf OpenWar.app
