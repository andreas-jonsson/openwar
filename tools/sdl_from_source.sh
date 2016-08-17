#!/bin/bash

wget -q http://www.libsdl.org/release/SDL2-2.0.3.tar.gz
tar xf SDL2-*.tar.gz
cd SDL2-* && ./configure && make && sudo make install

wget -q http://www.libsdl.org/projects/SDL_mixer/release/SDL2_mixer-2.0.1.tar.gz
tar xf SDL2_mixer-*.tar.gz
cd SDL2_mixer-* && ./configure && make && sudo make install
