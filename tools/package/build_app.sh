#!/bin/bash

if [ -d OpenWar.app ]; then rm -rf OpenWar.app; fi
cp -r "tools/package/OpenWar.app" OpenWar.app

mkdir OpenWar.app/Contents/Resources
mkdir OpenWar.app/Contents/Frameworks

wget -q https://dl.dropboxusercontent.com/u/1955192/SDL2.framework.zip
unzip -dq OpenWar.app/Contents/Frameworks SDL2*.zip

wget -q https://dl.dropboxusercontent.com/u/1955192/SDL2_mixer.framework.zip
unzip -dq OpenWar.app/Contents/Frameworks SDL2_mixer*.zip

cp openwar OpenWar.app/Contents/MacOS

zip -r OpenWar.zip OpenWar.app
rm -rf OpenWar.app
