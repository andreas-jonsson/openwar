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

package platform

import (
	"os"
	"path"
	"runtime"
)

var DataPath string

func RootJoin(p ...string) string {
	return path.Join(DataPath, path.Join(p...))
}

func init() {
	wd, _ := os.Getwd()
	DataPath = path.Join(wd, "data")

	if _, err := os.Stat(DataPath); os.IsNotExist(err) {
		switch runtime.GOOS {
		case "linux", "darwin":
			DataPath = "/usr/local/share/openwar"
		}
	}
}
