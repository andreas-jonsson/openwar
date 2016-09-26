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

package debug

import (
	"log"

	"github.com/andreas-jonsson/openwar/game"
	"github.com/andreas-jonsson/openwar/platform"
	"github.com/mattn/go-gtk/gtk"
)

func Start(cfg *game.Config) {
	builder := gtk.NewBuilder()
	if _, err := builder.AddFromFile(platform.RootJoin("debug.glade")); err != nil {
		log.Fatalln("could not load interface description:", err)
	}

	cfg.Debug.Race = "Orcs"

	setupDebugWindow(builder)
}

func setupDebugWindow(builder *gtk.Builder) {
	debugWindow := gtk.WidgetFromObject(builder.GetObject("debug_window"))
	debugWindow.ShowAll()
}
