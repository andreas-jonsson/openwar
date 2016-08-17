#!/bin/bash

DIR=openwar_0.0.1-$TRAVIS_BUILD_NUMBER

if [ -d "$DIR" ]; then rm -rf "$DIR"; fi
cp -r "tools/package/openwar_0.0.1-x" $DIR

mkdir -p $DIR/opt/openwar

cp -r data $DIR/opt/openwar
cp openwar $DIR/opt/openwar

rpl e34f19fc-199d-4fb9-b334-aed07b29a173 $TRAVIS_BUILD_NUMBER $DIR/DEBIAN/control

dpkg-deb --build $DIR
