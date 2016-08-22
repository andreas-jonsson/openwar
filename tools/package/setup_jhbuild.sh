#!/bin/bash

wget -q https://git.gnome.org/browse/gtk-osx/plain/gtk-osx-build-setup.sh
chmod +x gtk-osx-build-setup.sh

# Ensure that we get some output at least every 5min.
while :; do echo "ping!"; sleep 300; done &

./gtk-osx-build-setup.sh
export PATH=~/.local/bin:$PATH

echo "Building bootstrap..."
jhbuild bootstrap

echo "Building osx-bootstrap..."
jhbuild build meta-gtk-osx-bootstrap

wget -q http://ftp.gnome.org/pub/gnome/sources/gtk-mac-bundler/0.7/gtk-mac-bundler-0.7.3.tar.xz
tar xf gtk-mac-bundler-0.7.3.tar.xz
cd gtk-mac-bundler-0.7.3
make install
cd ..
