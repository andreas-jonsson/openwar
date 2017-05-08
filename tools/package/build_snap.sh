#!/bin/bash

if [ -d snap ]; then rm -rf snap; fi
mkdir snap
cp openwar snap

snapcraft
