#!/bin/bash

export PATH=~/.local/bin:$PATH

# Ensure that we get some output at least every 5min.
while :; do echo "ping!"; sleep 300; done &

# This step was moved to here because it takes to long to just do it all in the setup script.
jhbuild build meta-gtk-osx-core
