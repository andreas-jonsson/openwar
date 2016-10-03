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

	"github.com/andreas-jonsson/openwar/game"
	"github.com/andreas-jonsson/openwar/launcher"
)

const logo = `________                       __      __
\_____  \ ______   ____   ____/  \    /  \_____ _______
 /   |   \\____ \_/ __ \ /    \   \/\/   /\__  \\_  __ \
/    |    \  |_> >  ___/|   |  \        /  / __ \|  | \/
\_______  /   __/ \___  >___|  /\__/\  /  (____  /__|
        \/|__|        \/     \/      \/        \/`

func banner() {
	log.Print(logo)
	log.Println(" Ver:", game.VersionString)

	log.Println("\n\tAuthors:")
	for _, author := range game.Authors {
		log.Println("\t\t", author)
	}
	log.Println()
}

func main() {
	banner()
	launcher.Start()
}
