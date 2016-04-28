/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andreas-jonsson/openwar/resource"
	"github.com/andreas-jonsson/openwar/resource/debug"
)

const logo = `________                       __      __
\_____  \ ______   ____   ____/  \    /  \_____ _______
 /   |   \\____ \_/ __ \ /    \   \/\/   /\__  \\_  __ \
/    |    \  |_> >  ___/|   |  \        /  / __ \|  | \/
\_______  /   __/ \___  >___|  /\__/\  /  (____  /__|
        \/|__|        \/     \/      \/        \/
`

func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Println(logo)

	resourcePath := "./"
	if len(args) > 0 {
		resourcePath = args[0]
	}

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
	_, err = resource.LoadPal(arch)
	if err != nil {
		panic(err)
	}

	fmt.Println("Loading images...")
	images, err := resource.LoadImg(arch)
	if err != nil {
		panic(err)
	}

	fmt.Println("Loading cursors...")
	_, err = resource.LoadCur(arch, images)
	if err != nil {
		panic(err)
	}

	debug.DumpImg(images, "")
	debug.DumpArchive(arch, "")
}
