#!/bin/bash

./tools/dropbox/dropbox_uploader.sh download Public/jhbuild.7z
7z x jhbuild.7z -o$HOME -aoa > /dev/null

ls -al

wget -q http://ftp.gnome.org/pub/gnome/sources/gtk-mac-bundler/0.7/gtk-mac-bundler-0.7.3.tar.xz
tar xf gtk-mac-bundler-0.7.3.tar.xz
cd gtk-mac-bundler-0.7.3
make install
cd ..
