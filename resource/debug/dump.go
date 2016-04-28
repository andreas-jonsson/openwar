/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package debug

import (
	"image/png"
	"os"
	"path"

	"github.com/andreas-jonsson/openwar/resource"
)

func DumpImg(images resource.Images, p string) {
	outputPath := "img"
	if p != "" {
		outputPath = p
	}

	os.MkdirAll(outputPath, 0755)

	for file, image := range images {
		outfile, err := os.Create(path.Join(outputPath, file) + ".png")
		if err != nil {
			panic(err)
		}
		if err := png.Encode(outfile, image.Data); err != nil {
			panic(err)
		}
		outfile.Close()
	}
}

func DumpArchive(arch *resource.Archive, p string) {
	outputPath := "archive"
	if p != "" {
		outputPath = p
	}

	os.MkdirAll(outputPath, 0755)

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
