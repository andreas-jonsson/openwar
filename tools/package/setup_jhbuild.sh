#!/bin/bash

wget -q https://git.gnome.org/browse/gtk-osx/plain/gtk-osx-build-setup.sh
chmod +x gtk-osx-build-setup.sh

# Ensure that we get some output at least every 5min.
while :; do echo "ping!"; sleep 300; done &

./gtk-osx-build-setup.sh > jhbuild_log.txt 2>&1
export PATH=~/.local/bin:$PATH
jhbuild bootstrap > jhbuild_log.txt 2>&1
jhbuild build meta-gtk-osx-bootstrap > jhbuild_log.txt 2>&1

wget -q http://ftp.gnome.org/pub/gnome/sources/gtk-mac-bundler/0.7/gtk-mac-bundler-0.7.3.tar.xz
tar xf gtk-mac-bundler-0.7.3.tar.xz
cd gtk-mac-bundler-0.7.3
make install
cd ..

7z a -tzip jhbuild.zip ~/.local ~/gtk
./tools/dropbox/dropbox_uploader.sh upload jhbuild.zip Public/jhbuild.zip
