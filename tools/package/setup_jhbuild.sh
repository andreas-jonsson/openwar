#!/bin/bash

rm -rf $GOPATH/pkg $GOPATH/bin

wget -q https://git.gnome.org/browse/gtk-osx/plain/gtk-osx-build-setup.sh
chmod +x gtk-osx-build-setup.sh

./gtk-osx-build-setup.sh > /dev/null
export PATH=~/.local/bin:$PATH
jhbuild bootstrap > /dev/null
jhbuild build meta-gtk-osx-bootstrap meta-gtk-osx-core > /dev/null

wget -q http://ftp.gnome.org/pub/gnome/sources/gtk-mac-bundler/0.7/gtk-mac-bundler-0.7.3.tar.xz
tar xf gtk-mac-bundler-0.7.3.tar.xz
cd gtk-mac-bundler-0.7.3
make install
cd ..

mkdir jhbuild
cp -r .local jhbuild
cp -r gtk jhbuild
7z a -tzip jhbuild.zip jhbuild
./tools/dropbox/dropbox_uploader.sh upload jhbuild.zip jhbuild.zip
