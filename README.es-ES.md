

# OpenWar

[![OSX & Linux Build](https://travis-ci.org/andreas-jonsson/openwar.svg?branch=master)](https://travis-ci.org/andreas-jonsson/openwar) [![Windows Build](https://ci.appveyor.com/api/projects/status/erhgfi08p3amtaec?svg=true)](https://ci.appveyor.com/project/andreas-t-jonsson/openwar)
[![Donate](https://img.shields.io/badge/donate-itch.io-red.svg)](https://phix.itch.io/openwar/purchase)

## Acerca de

OpenWar es un motor de juego alternativo para Warcraft: Orcs & Humans. Por lo tanto, a menos que tengas una *copia legal* de Warcraft: Orcs & Humans (se requiere la versión original para MS-DOS, no funcionará con las versiones para Mac o Demo), OpenWar te será bastante inútil, ya que no incluye gráficos ni sonidos por sí mismo.

Dado que OpenWar es un motor de juego diferente, no todo funcionará exactamente igual que en el Warcraft original. Si deseas la experiencia inalterada del Warcraft original, todavía deberás usar el juego original.

<img src="https://raw.githubusercontent.com/andreas-jonsson/openwar/master/doc/screenshot1.png" width="250"><img src="https://raw.githubusercontent.com/andreas-jonsson/openwar/master/doc/screenshot4.png" width="300"><img src="https://raw.githubusercontent.com/andreas-jonsson/openwar/master/doc/screenshot2.png" width="250">

## Aviso legal

OpenWar *no* es un producto oficial de Blizzard; es una modificación de Warcraft: Orcs & Humans, hecha por fans de Warcraft para fans de Warcraft. Necesitas una copia de la versión original de Warcraft: Orcs & Humans para MS-DOS para poder utilizar OpenWar. *Warcraft* es una marca registrada de [Blizzard Entertainment](https://www.blizzard.com).

## Compilación

Utiliza [Homebrew](http://brew.sh) o [Linuxbrew](http://linuxbrew.sh) para la compilación.

```bash
brew tap andreas-jonsson/tap
brew install openwar
```

o la forma clásica y tradicional.

```bash
# Start by installing all external dependencies:
# GCC/LLVM/Mingw, Go1.8, Git, SDL2, SDL2_mixer (with Timidity support).

export GOPATH=$HOME                               # Make sure you have a GOPATH set to your Go workspace.
go get github.com/andreas-jonsson/openwar         # Download the project using Go.
cd $GOPATH/src/github.com/andreas-jonsson/openwar
go build openwar.go                               # Build or run OpenWar.
```

## Desarrollo

El juego aún no es jugable, pero la mayoría de los recursos ya están decodificados y pueden cargarse.

* Imágenes y Sprites
* Paletas
* Tilesets
* Mapas
* Sonidos
* Música
* Diálogos, texto de misiones, etc.

## Otros proyectos

Puedes encontrar más clones y remakes de juegos de código abierto en [osgameclones.com](http://osgameclones.com/).

## Licencia
```
Este programa es software libre: puede redistribuirlo y/o modificarlo
bajo los términos de la Licencia Pública General de GNU publicada por
la Free Software Foundation, ya sea la versión 3 de la Licencia, o
(a su elección) cualquier versión posterior.

Este programa se distribuye con la esperanza de que sea útil,
pero SIN NINGUNA GARANTÍA; sin incluso la garantía implícita de
COMERCIABILIDAD o IDONEIDAD PARA UN PROPÓSITO PARTICULAR. Consulte la
Licencia Pública General de GNU para más detalles.

Debería haber recibido una copia de la Licencia Pública General de GNU
junto con este programa. Si no, consulte <http://www.gnu.org/licenses/>.
```
