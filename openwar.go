/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/openwar-hq/openwar/platform"
	"github.com/openwar-hq/openwar/resource"
	"github.com/openwar-hq/openwar/resource/debug"
	"github.com/openwar-hq/openwar/xmi"
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

	var warFile [1]string
	filepath.Walk(resourcePath, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() && strings.ToUpper(filepath.Base(path)) == "DATA.WAR" {
			warFile[0] = path
			fmt.Println("Found resources: " + path)
			return filepath.SkipDir
		}
		return nil
	})

	if warFile[0] == "" {
		fmt.Println("Could not find all game resources.")
		return
	}

	arch, err := resource.OpenArchive(warFile[0])
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
	_, err = resource.LoadSprites(arch, images)
	if err != nil {
		panic(err)
	}

	fmt.Println("Converting XMI...")
	if err = convertXMI(arch); err != nil {
		panic(err)
	}

	if err = platform.Init(); err != nil {
		panic(err)
	}
	defer platform.Shutdown()

	rend, err := platform.NewRenderer(640, 360, "OpenWar")
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

	pal := combinePalettes(palettes["FOREST.PAL"], palettes["SPRITE0.PAL"])

	debug.DumpImg(images, pal, "")
	debug.DumpArchive(arch, "")
	arch.Close()

	//if err = player.PlayMusic("MUSIC01.XMI", 0, 0); err != nil {
	//	panic(err)
	//}

	fmt.Println("Running...")
	for {
		for event := platform.PollEvent(); event != nil; event = platform.PollEvent() {
			switch event.(type) {
			case *platform.QuitEvent:
				return
			}
		}

		rend.Clear()

		img := images["TITLE.IMG"]
		pal := palettes["TITLE.PAL"]

		rend.BlitPal(img.Data, pal, image.Point{})

		rend.Present()
	}
}

func convertXMI(arch *resource.Archive) error {
	for file, data := range arch.Files {
		if path.Ext(file) == ".XMI" {
			mid, err := xmi.ToMidi(data, xmi.NoConversion)
			if err != nil {
				return err
			}
			arch.Files[file] = mid
		}
	}
	return nil
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

func combinePalettes(low, high color.Palette) color.Palette {
	if len(low)+len(high) != 256 {
		return nil
	}

	pal := make([]color.Color, 256)
	copy(pal, low)
	copy(pal[128:], high)
	return pal
}
