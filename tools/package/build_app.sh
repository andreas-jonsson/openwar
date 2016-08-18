#!/bin/bash

DIR=OpenWar.app

if [ -d "$DIR" ]; then rm -rf "$DIR"; fi
cp -r "tools/package/OpenWar.app" $DIR

wget -q https://www.libsdl.org/release/SDL2-2.0.4.dmg
7z x SDL2-2.0.4.dmg
cp -r SDL2/SDL2.framework $DIR/Contents/Resources

wget -q https://www.libsdl.org/projects/SDL_mixer/release/SDL2_mixer-2.0.1.dmg
7z x SDL2_mixer-2.0.1.dmg
cp -r SDL2_mixer/SDL2_mixer.framework $DIR/Contents/Resources

mkdir $DIR/Contents/MacOS
cp openwar $DIR/Contents/MacOS
cp -r data $DIR/Contents/MacOS

7z a -tzip OpenWar.zip OpenWar.app
