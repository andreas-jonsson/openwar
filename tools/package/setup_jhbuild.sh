#!/bin/bash

rm -rf $GOPATH/pkg $GOPATH/bin

wget -q https://git.gnome.org/browse/gtk-osx/plain/gtk-osx-build-setup.sh
chmod +x gtk-osx-build-setup.sh

./gtk-osx-build-setup.sh
export PATH=~/.local/bin:$PATH
jhbuild bootstrap | /dev/null
jhbuild build meta-gtk-osx-bootstrap meta-gtk-osx-core | /dev/null

wget -q http://ftp.gnome.org/pub/gnome/sources/gtk-mac-bundler/0.7/gtk-mac-bundler-0.7.3.tar.xz
tar xf gtk-mac-bundler-0.7.3.tar.xz
cd gtk-mac-bundler-0.7.3
make install
cd ..
