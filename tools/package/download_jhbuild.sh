#!/bin/bash

./tools/dropbox/dropbox_uploader.sh download Public/jhbuild.zip
7z x jhbuild.zip -o~/ -aoa > /dev/null

mv -f jhbuild/.cache $HOME
mv -f jhbuild/.jhbuild* $HOME
mv -f jhbuild/gtk $HOME
mv -f jhbuild/tmp-jhbuild-revision $HOME

ls -al

wget -q http://ftp.gnome.org/pub/gnome/sources/gtk-mac-bundler/0.7/gtk-mac-bundler-0.7.3.tar.xz
tar xf gtk-mac-bundler-0.7.3.tar.xz
cd gtk-mac-bundler-0.7.3
make install
cd ..
