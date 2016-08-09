/*
Copyright (C) 2016 Andreas T Jonsson

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
*/

package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/andreas-jonsson/openwar/game"
	"github.com/andreas-jonsson/openwar/platform"
	"github.com/andreas-jonsson/openwar/resource"
)

const versionString = "0.0.1"

const logo = `________                       __      __
\_____  \ ______   ____   ____/  \    /  \_____ _______
 /   |   \\____ \_/ __ \ /    \   \/\/   /\__  \\_  __ \
/    |    \  |_> >  ___/|   |  \        /  / __ \|  | \/
\_______  /   __/ \___  >___|  /\__/\  /  (____  /__|
        \/|__|        \/     \/      \/        \/`

func banner() {
	fmt.Print(logo)
	fmt.Println(" Ver:", versionString)

	fmt.Println("\n\tProgrammed by:")
	for _, author := range authors {
		fmt.Println("\t\t", author)
	}
	fmt.Println()
}

var (
	resourcePath string
	fullscreen   bool
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: openwar [options]\n\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&resourcePath, "war", "./", "search path for .WAR archives")
	flag.BoolVar(&fullscreen, "fs", false, "run the game in fullscreen mode")
}

func main() {
	flag.Parse()
	banner()

	path := filepath.Walk(resourcePath, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() && strings.ToUpper(filepath.Base(path)) == "DATA.WAR" {
			fmt.Println("Found resources: " + path)
			return errors.New(path)
		}
		return nil
	}).Error()

	if path == "" {
		fmt.Println("Could not find all game resources.")
		return
	}

	rand.Seed(time.Now().UnixNano())

	//resource.Logger = os.Stdout
	//resource.LoadUnsupported = true
	arch, err := resource.OpenArchive(path)
	if err != nil {
		panic(err)
	}

	fmt.Println("Loading palettes...")
	palettes, err := resource.LoadPalettes(arch)
	if err != nil {
		panic(err)
	}

	fmt.Println("Loading images...")
	images, err := resource.LoadImages(arch)
	if err != nil {
		panic(err)
	}

	fmt.Println("Loading sprites...")
	sprites, err := resource.LoadSprites(arch, images)
	if err != nil {
		panic(err)
	}

	fmt.Println("Loading tilesets...")
	tilesets, err := resource.LoadTilesets(arch, images, palettes)
	if err != nil {
		panic(err)
	}

	if err = platform.Init(); err != nil {
		panic(err)
	}
	defer platform.Shutdown()

	rend, err := platform.NewRenderer(640, 480, "OpenWar")
	if err != nil {
		panic(err)
	}
	defer rend.Shutdown()

	player, err := platform.NewAudioPlayer()
	if err != nil {
		panic(err)
	}
	defer player.Shutdown()

	fmt.Println("Loading audio...")
	if err = loadAudio(arch, player); err != nil {
		panic(err)
	}

	//debug.DumpImg(images, resource.CombinePalettes(palettes["FOREST.PAL"], palettes["SPRITE0.PAL"]), "")
	//debug.DumpArchive(arch, "")

	res := resource.Resources{Palettes: palettes, Images: images, Sprites: sprites, Tilesets: tilesets, Archive: arch}
	g, err := game.NewGame(rend, player, res)
	if err != nil {
		panic(err)
	}
	defer g.Shutdown()

	if err := g.SwitchState("menu"); err != nil {
		panic(err)
	}

	for g.Running() {
		rend.Clear()

		if err := g.Update(); err != nil {
			panic(err)
		}
		if err := g.Render(); err != nil {
			panic(err)
		}

		rend.Present()
	}
}

func loadAudio(arch *resource.Archive, player platform.AudioPlayer) error {
	for file, data := range arch.Files {
		switch path.Ext(file) {
		case ".XMI":
			if err := player.LoadMusic(file, data); err != nil {
				return err
			}
		case ".VOC", ".WAV":
			if _, err := player.LoadSound(file, data); err != nil {
				return err
			}
		}
	}
	return nil
}
