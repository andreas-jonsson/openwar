#!/bin/bash

if [ -d OpenWar.app ]; then rm -rf OpenWar.app; fi
cp -r "tools/package/OpenWar.app" OpenWar.app

cp openwar OpenWar.app/Contents/MacOS
dylibbundler -od -b -x OpenWar.app/Contents/MacOS/openwar -d OpenWar.app/Contents/libs

zip -rq openwar_${OPENWAR_VERSION}_osx.zip OpenWar.app
rm -rf OpenWar.app
