#!/bin/bash

if [ -d OpenWar.app ]; then rm -rf OpenWar.app; fi
cp -r "tools/package/OpenWar.app" OpenWar.app

cp -r data OpenWar.app/Contents/Resources
cp openwar OpenWar.app/Contents/MacOS
dylibbundler -od -b -x OpenWar.app/Contents/MacOS/openwar -d OpenWar.app/Contents/libs

zip -rq OpenWar.zip OpenWar.app
rm -rf OpenWar.app
