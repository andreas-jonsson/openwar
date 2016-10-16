// +build ignore

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
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

type fsWrapper struct {
	internal http.FileSystem
}

func (fs fsWrapper) Open(name string) (http.File, error) {
	return fs.internal.Open("data/src/" + name)
}

func main() {
	fs := fsWrapper{http.Dir("")}
	err := vfsgen.Generate(&fs, vfsgen.Options{
		Filename:     "data/data.go",
		PackageName:  "data",
		VariableName: "FS",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
