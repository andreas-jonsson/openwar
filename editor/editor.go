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

package editor

import (
	"log"

	"github.com/andreas-jonsson/openwar/game"
	"github.com/andreas-jonsson/openwar/platform"
	"github.com/mattn/go-gtk/gtk"
)

func Start(cfg *game.Config) {
	builder := gtk.NewBuilder()
	if _, err := builder.AddFromFile(platform.RootJoin("editor.glade")); err != nil {
		log.Fatalln("could not load interface description:", err)
	}

	setupEditorWindow(builder)
}

func setupEditorWindow(builder *gtk.Builder) {
	editorWindow := gtk.WidgetFromObject(builder.GetObject("editor_window"))
	editorWindow.ShowAll()

	//editorWindow.Connect("delete_event", func() {
	//gtk.MainQuit()
	//})
}
