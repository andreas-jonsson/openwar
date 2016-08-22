#!/bin/bash

./tools/dropbox/dropbox_uploader.sh download Public/jhbuild.zip
export PATH=~/.local/bin:$PATH

7z x jhbuild.zip -o~/ -aoa
