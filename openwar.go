/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package main

import (
	"fmt"
	"os"
	"path"

	"github.com/andreas-jonsson/openwar/resource"
)

func main() {
	arch, err := resource.OpenArchive("DATA.WAR")
	if err != nil {
		panic(err)
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
