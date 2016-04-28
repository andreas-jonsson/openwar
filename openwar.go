/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/andreas-jonsson/openwar/platform"
	"github.com/andreas-jonsson/openwar/resource"
	"github.com/andreas-jonsson/openwar/resource/debug"
	"github.com/andreas-jonsson/openwar/xmi"
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
	_, err = resource.LoadPalettes(arch)
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

	debug.DumpImg(images, "")
	debug.DumpArchive(arch, "")

	arch.Close()

	if err = platform.Init(); err != nil {
		panic(err)
	}
	defer platform.Shutdown()

	rend, err := platform.NewRenderer(640, 360, "OpenWar")
	if err != nil {
		panic(err)
	}
	defer rend.Shutdown()

	for {
		for event := platform.PollEvent(); event != nil; event = platform.PollEvent() {
			switch event.(type) {
			case *platform.QuitEvent:
				return
			}
		}

		rend.Clear()

		time.Sleep(16 * time.Millisecond)

		rend.Present()
	}
}

func convertXMI(arch *resource.Archive) error {
	for file, data := range arch.Files {
		if path.Ext(file) == ".XMI" {
			mid, err := xmi.ToMid(data, xmi.NoConversion)
			if err != nil {
				return err
			}
			arch.Files[file] = mid
		}
	}
	return nil
}
