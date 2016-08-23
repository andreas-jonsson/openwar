#!/bin/bash

./tools/dropbox/dropbox_uploader.sh download Public/jhbuild.zip
7z x jhbuild.zip -o~/ -aoa > /dev/null
mv -f jhbuild/* ~/

wget -q http://ftp.gnome.org/pub/gnome/sources/gtk-mac-bundler/0.7/gtk-mac-bundler-0.7.3.tar.xz
tar xf gtk-mac-bundler-0.7.3.tar.xz
cd gtk-mac-bundler-0.7.3
make install
cd ..
