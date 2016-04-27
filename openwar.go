/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package main

import (
	"flag"
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/andreas-jonsson/openwar/resource"
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

	// No logging.
	resource.Logger = ioutil.Discard

	arch, err := resource.OpenArchive(warFile[0])
	if err != nil {
		panic(err)
	}

	images, err := resource.LoadImg(arch)
	if err != nil {
		panic(err)
	}

	for file, image := range images {
		outfile, err := os.Create(path.Join("img", file) + ".png")
		if err != nil {
			panic(err)
		}
		if err := png.Encode(outfile, image); err != nil {
			panic(err)
		}
		outfile.Close()
	}

	var (
		num        int
		outputPath string
	)

	for {
		outputPath = fmt.Sprintf("DATA%v.WAR", num)
		if _, err = os.Stat(outputPath); err != nil {
			os.Mkdir(outputPath, 0755)
			break
		}
		num++
	}

	for fileName, data := range arch.Files {
		fp, err := os.Create(path.Join(outputPath, fileName))
		if err != nil {
			panic(err)
		}

		if num, err := fp.Write(data); num != len(data) || err != nil {
			panic(err)
		}
		fp.Close()
	}
}
