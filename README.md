# OpenWar

[![Chat](https://badges.gitter.im/andreas-jonsson/openwar.svg)](https://gitter.im/andreas-jonsson/openwar?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Contribute](https://img.shields.io/badge/contribute-FreedomSponsors-blue.svg)](https://freedomsponsors.org/project/341)
[![Website](https://img.shields.io/badge/project-website-red.svg)](http://www.openwar.io)

## About

OpenWar is an alternative Warcraft: Orcs & Humans game engine. So unless you have a *legal copy* of Warcraft: Orcs & Humans (original MS-DOS version required, won't work with the Mac or Demo versions) OpenWar will be pretty useless to you, since it doesn't come with any graphics or sounds itself.

Since OpenWar is a different game engine, not all things will work 100% the same as they did in the original Warcraft, if you want the original unchanged Warcraft experience, you will still have to use the original game.

<img src="https://raw.githubusercontent.com/andreas-jonsson/openwar/master/doc/screenshot1.png" width="250">
<img src="https://raw.githubusercontent.com/andreas-jonsson/openwar/master/doc/screenshot3.gif" width="250">
<img src="https://raw.githubusercontent.com/andreas-jonsson/openwar/master/doc/screenshot2.png" width="250">

## Disclaimer

OpenWar is *not* an official Blizzard product, its a Warcraft: Orcs & Humans modification, by Warcraft fans for Warcraft fans. You need a copy of the original Warcraft: Orcs & Humans MS-DOS version to make use of OpenWar. *Warcraft* is a registered trademark of [Blizzard Entertainment](https://www.blizzard.com).

## Build

##### Linux & OSX: [![OSX & Linux Build](https://travis-ci.org/andreas-jonsson/openwar.svg?branch=master)](https://travis-ci.org/andreas-jonsson/openwar)

##### Windows [![Windows Build](https://ci.appveyor.com/api/projects/status/erhgfi08p3amtaec?svg=true)](https://ci.appveyor.com/project/andreas-t-jonsson/openwar)

Use [Homebrew](http://brew.sh) or [Linuxbrew](http://linuxbrew.sh) for building.

```bash
brew tap andreas-jonsson/tap
brew install openwar
```

Or the good old fashioned way.

```bash
# Start by installing all external dependencies:
# GCC/LLVM/Mingw, Go1.6, Git, GTK+2, SDL2, SDL2_mixer (with Timidity support).

export GOPATH=$HOME                                                      # Make sure you have a GOPATH set to your Go workspace.
go get github.com/andreas-jonsson/openwar                                # Download the project using Go.
cd $GOPATH/src/github.com/andreas-jonsson/openwar && go build openwar.go # Build or run OpenWar.
```

## Development

The game is not yet playable but most resources are now decoded and loadable.

* Images & Sprites
* Palettes
* Tilesets
* Maps
* Sound
* Music
* Dialog, Mission text, etc.

## Other Projects

You can find more open-source game clones and remakes on [osgameclones.com](http://osgameclones.com/).

## License
```
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
```
