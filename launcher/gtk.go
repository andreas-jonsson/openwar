// -build nogui

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

package launcher

import (
	"log"
	"unsafe"

	"github.com/andreas-jonsson/openwar/editor"
	"github.com/andreas-jonsson/openwar/game"
	"github.com/andreas-jonsson/openwar/platform"
	"github.com/andreas-jonsson/openwar/resource"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gtk"
)

var war *resource.Archive

func Start() {
	log.Println("Starting launcher...")

	gtk.Init(nil)
	builder := gtk.NewBuilder()

	if _, err := builder.AddFromFile(platform.RootJoin("launcher.glade")); err != nil {
		log.Fatalln("could not load interface description:", err)
	}

	setupLauncherWindow(builder)
	gtk.Main()
}

func createConfig(builder *gtk.Builder) *game.Config {
	fullscreen := true // (*gtk.CheckButton)(unsafe.Pointer(builder.GetObject("fullscreen_checkbutton"))).GetActive()
	widescreen := true //(*gtk.CheckButton)(unsafe.Pointer(builder.GetObject("widescreen_checkbutton"))).GetActive()
	wc2Input := true   // (*gtk.CheckButton)(unsafe.Pointer(builder.GetObject("wc2_input_checkbutton"))).GetActive()

	return &game.Config{
		Fullscreen: fullscreen,
		Widescreen: widescreen,
		WC2Input:   wc2Input,
	}
}

func setSensitive(builder *gtk.Builder, sensitive bool) {
	joinButton := (*gtk.Button)(unsafe.Pointer(builder.GetObject("join_button")))
	joinButton.SetSensitive(sensitive)

	//createButton := (*gtk.Button)(unsafe.Pointer(builder.GetObject("create_button")))
	//createButton.SetSensitive(sensitive)

	editorButton := (*gtk.Button)(unsafe.Pointer(builder.GetObject("editor_button")))
	editorButton.SetSensitive(sensitive)
}

func setupLauncherWindow(builder *gtk.Builder) {
	launcherWindow := gtk.WidgetFromObject(builder.GetObject("launcher_window"))
	launcherWindow.ShowAll()

	launcherWindow.Connect("delete_event", func() {
		gtk.MainQuit()
	})

	builder.GetObject("open_button").Connect("clicked", func() {
		fileDialog := gtk.NewFileChooserDialog("Open", launcherWindow.GetTopLevelAsWindow(), gtk.FILE_CHOOSER_ACTION_OPEN, gtk.STOCK_CANCEL, gtk.RESPONSE_CANCEL, gtk.STOCK_OPEN, gtk.RESPONSE_OK)

		filter := gtk.NewFileFilter()
		filter.SetName("Warcraft Data Archive")
		filter.AddPattern("DATA.WAR")
		fileDialog.AddFilter(filter)

		if fileDialog.Run() == gtk.RESPONSE_OK {
			file := fileDialog.GetFilename()
			entry := (*gtk.Entry)(unsafe.Pointer(builder.GetObject("resource_entry")))
			entry.SetText(file)
			fileDialog.Hide()

			img := (*gtk.Image)(unsafe.Pointer(builder.GetObject("resource_image")))

			var err error
			if war, err = resource.OpenArchive(file); err == nil {
				img.SetFromStock("gtk-ok", gtk.ICON_SIZE_BUTTON)
				setSensitive(builder, true)
			} else {
				img.SetFromStock("gtk-cancel", gtk.ICON_SIZE_BUTTON)
				setSensitive(builder, false)
			}
		} else {
			fileDialog.Hide()
		}
	}, nil)

	builder.GetObject("join_button").Connect("clicked", func() {
		launcherWindow.SetSensitive(false)
		go func() {
			game.Start(createConfig(builder), war)

			gdk.ThreadsEnter()
			launcherWindow.SetSensitive(true)
			gdk.ThreadsLeave()
		}()
	}, nil)

	builder.GetObject("editor_button").Connect("clicked", func() {
		launcherWindow.SetSensitive(false)
		editor.Start(createConfig(builder), war, func() {
			launcherWindow.SetSensitive(true)
		})
	}, nil)

	builder.GetObject("about_button").Connect("clicked", func() {
		aboutDialog := (*gtk.AboutDialog)(unsafe.Pointer(builder.GetObject("about_dialog")))
		aboutDialog.Run()
		aboutDialog.Hide()
	}, nil)
}
