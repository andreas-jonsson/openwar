// +build tui

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
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/andreas-jonsson/openwar/game"
	"github.com/andreas-jonsson/openwar/resource"

	wmenu "gopkg.in/dixonwille/wmenu.v2"
)

var (
	arch = "[Not Installed]"
	war  *resource.Archive

	cfg = &game.Config{
		Fullscreen: false,
		Widescreen: true,
		WC2Input:   true,
	}
)

func Start() {
	cfg.Debug.Map = "HUMAN01"
	cfg.Debug.Race = "Human"

	mainMenu()
}

func optionsToString() string {
	var s string
	if cfg.Fullscreen {
		s += "Fullscreen "
	}
	if cfg.Widescreen {
		s += "Widescreen "
	}
	if cfg.WC2Input {
		s += "WC2Input "
	}

	if s == "" {
		return "[]"
	}
	return "[" + s + "\x08]"
}

func mainMenu() {
	banner()

	for i := 0; true; i++ {
		menu := wmenu.NewMenu("\nSelect an option.")
		menu.Option("Start Game", "0", war != nil, func() error {
			if war == nil {
				installArchiveMenu()
			} else {
				banner()
				game.Start(cfg, war)
				os.Exit(0)
			}
			return nil
		})

		menu.Option("Archive: "+arch, "1", war == nil, installArchiveMenu)
		menu.Option("Options: "+optionsToString(), "2", false, optionsMenu)
		menu.Option("Race: "+cfg.Debug.Race, "3", false, raceMenu)
		menu.Option("Map: "+cfg.Debug.Map, "4", false, mapMenu)

		if i > 0 {
			menu.ClearOnMenuRun()
		}

		menu.Run()
	}
}

func optionsMenu() error {
	f := func(opts []wmenu.Opt) error {
		cfg.Fullscreen = false
		cfg.Widescreen = false
		cfg.WC2Input = false

		for _, opt := range opts {
			switch opt.Text {
			case "Fullscreen":
				cfg.Fullscreen = true
			case "Widescreen":
				cfg.Widescreen = true
			case "WC2Input":
				cfg.WC2Input = true
			}
		}
		return nil
	}

	menu := wmenu.NewMenu("\nSelect options.")
	menu.Action(func(opt wmenu.Opt) error { return f([]wmenu.Opt{opt}) })
	menu.MultipleAction(f)

	menu.Option("Fullscreen", "0", cfg.Fullscreen, nil)
	menu.Option("Widescreen", "1", cfg.Widescreen, nil)
	menu.Option("WC2Input", "2", cfg.WC2Input, nil)

	menu.LoopOnInvalid()
	menu.ClearOnMenuRun()

	if err := menu.Run(); err != nil {
		log.Panicln(err)
	}
	return nil
}

func raceMenu() error {
	menu := wmenu.NewMenu("\nSelect a race.")
	menu.Action(func(opt wmenu.Opt) error {
		cfg.Debug.Race = opt.Text
		return nil
	})

	menu.Option("Human", "0", true, nil)
	menu.Option("Orc", "1", false, nil)

	menu.ClearOnMenuRun()
	menu.LoopOnInvalid()

	if err := menu.Run(); err != nil {
		log.Panicln(err)
	}
	return nil
}

func mapMenu() error {
	if war == nil {
		installArchiveMenu()
		if war == nil {
			return nil
		}
	}

	menu := wmenu.NewMenu("\nSelect a map.")
	menu.Action(func(opt wmenu.Opt) error {
		cfg.Debug.Map = opt.Text
		return nil
	})

	maps := []string{}
	for file, _ := range war.Files {
		if path.Ext(file) == ".TER" {
			maps = append(maps, strings.TrimSuffix(file, ".TER"))
		}
	}
	sort.Strings(maps)
	for i, m := range maps {
		menu.Option(m, fmt.Sprint(i), i == 0, nil)
	}

	menu.ClearOnMenuRun()
	menu.LoopOnInvalid()

	if err := menu.Run(); err != nil {
		log.Panicln(err)
	}
	return nil
}

func installArchiveMenu() error {
	usr, _ := user.Current()
	menu := wmenu.NewMenu("Do you want to search for it in your home directory?")
	menu.Action(func(opt wmenu.Opt) error {
		if opt.ID == 0 {
			searchForArchive(usr.HomeDir)
		}
		return nil
	})

	menu.IsYesNo(0)
	menu.LoopOnInvalid()
	menu.ClearOnMenuRun()

	if err := menu.Run(); err != nil {
		log.Panicln(err)
	}
	return nil
}

func searchForArchive(searchPath string) bool {
	war = nil
	arch = ""

	filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if info.IsDir() {
			fmt.Printf("\033[H\033[2J" + path)
		} else if strings.ToUpper(filepath.Base(path)) == "DATA.WAR" {
			if archive, err := resource.OpenArchive(path); err == nil {
				war = archive
				arch = path
				return errors.New(path)
			}
		}
		return nil
	})

	return war != nil
}
