#!/bin/bash

7z a -tzip jhbuild.zip ~/.local ~/gtk
./tools/dropbox/dropbox_uploader.sh upload jhbuild.zip Public/jhbuild.zip
