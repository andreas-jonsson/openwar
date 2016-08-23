#!/bin/bash

./tools/dropbox/dropbox_uploader.sh download Public/jhbuild.zip
7z x jhbuild.zip -o~/ -aoa > /dev/null

mv -f jhbuild/.cache ~/
mv -f jhbuild/.jhbuild* ~/
mv -f jhbuild/gtk ~/
mv -f jhbuild/tmp-jhbuild-revision ~/

ls -al

wget -q http://ftp.gnome.org/pub/gnome/sources/gtk-mac-bundler/0.7/gtk-mac-bundler-0.7.3.tar.xz
tar xf gtk-mac-bundler-0.7.3.tar.xz
cd gtk-mac-bundler-0.7.3
make install
cd ..
