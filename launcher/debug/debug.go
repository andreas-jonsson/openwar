// +build !ndebug

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
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"unsafe"

	"github.com/andreas-jonsson/openwar/game"
	"github.com/andreas-jonsson/openwar/platform"
	"github.com/andreas-jonsson/openwar/resource"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gtk"
)

var (
	builder *gtk.Builder
	cfg     *game.Config
)

func Start(config *game.Config) {
	cfg = config
	builder = gtk.NewBuilder()
	if _, err := builder.AddFromFile(platform.RootJoin("debug.glade")); err != nil {
		log.Fatalln("could not load interface description:", err)
	}
	setupDebugWindow()
}

func ArchiveLoaded(war *resource.Archive) {
	maps := []string{}
	mapCombobox := (*gtk.ComboBoxText)(unsafe.Pointer(builder.GetObject("map_comboboxtext")))

	for file, _ := range war.Files {
		if path.Ext(file) == ".TER" {
			maps = append(maps, strings.TrimSuffix(file, ".TER"))
		}
	}

	sort.Strings(maps)
	log.Println("Available maps:")

	for _, m := range maps {
		mapCombobox.AppendText(m)
		log.Println(m)
	}
	mapCombobox.SetActive(0)

	cfg.Debug.Map = maps[0]
	mapCombobox.Connect("changed", func() {
		cfg.Debug.Map = maps[mapCombobox.GetActive()]
	})
}

func redirectLog(logTextView *gtk.TextView) {
	r, w, err := os.Pipe()
	if err != nil {
		log.Fatalln(err)
	}

	oldStdout := os.Stdout
	os.Stdout = w

	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			text := scanner.Text() + "\n"
			fmt.Fprint(oldStdout, text)

			gdk.ThreadsEnter()

			var iter gtk.TextIter
			buffer := logTextView.GetBuffer()
			buffer.GetStartIter(&iter)
			buffer.Insert(&iter, text)

			gdk.ThreadsLeave()
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(oldStdout, "pipe error:", err)
		}
	}()
}

func setupDebugWindow() {
	logTextView := (*gtk.TextView)(unsafe.Pointer(builder.GetObject("log_textview")))
	redirectLog(logTextView)

	raceCombobox := (*gtk.ComboBoxText)(unsafe.Pointer(builder.GetObject("race_comboboxtext")))
	races := []string{"Human", "Orc"}
	for _, r := range races {
		raceCombobox.AppendText(r)
	}
	raceCombobox.SetActive(0)

	cfg.Debug.Race = races[0]
	raceCombobox.Connect("changed", func() {
		cfg.Debug.Race = races[raceCombobox.GetActive()]
	})

	debugWindow := gtk.WidgetFromObject(builder.GetObject("debug_window"))
	debugWindow.ShowAll()
}
