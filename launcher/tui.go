// +build !nogui,!js,!mobile

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
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/andreas-jonsson/openwar/game"
	"github.com/andreas-jonsson/openwar/platform"
	"github.com/andreas-jonsson/openwar/resource"

	wmenu "gopkg.in/dixonwille/wmenu.v2"
)

const notInstalledText = "[Not Installed]"

var (
	arch = notInstalledText

	war *resource.Archive
	cfg *game.Config
)

func Start() {
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
	cfg = loadConfig()

	dataPath := "DATA.WAR"
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		dataPath = filepath.Clean(platform.CfgRootJoin(dataPath))
	}

	if archive, err := resource.OpenArchive(dataPath); err == nil {
		war = archive
		arch = dataPath
	}

	for {
		banner()

		menu := wmenu.NewMenu("\nSelect an option or press Ctrl+C to quit.")
		menu.Option("Start Game", "0", war != nil, func() error {
			if war == nil {
				installArchiveMenu()
			} else {
				banner()
				game.Start(cfg, war)
			}
			return nil
		})

		archText := arch
		if war != nil {
			archText = fmt.Sprintf("%s (%s)", arch, war.Type)
		}

		menu.Option("Archive: "+archText, "1", war == nil, installArchiveMenu)
		menu.Option("Options: "+optionsToString(), "2", false, optionsMenu)
		menu.Option("Race: "+cfg.Debug.Race, "3", false, raceMenu)
		menu.Option("Map: "+cfg.Debug.Map, "4", false, mapMenu)

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

	for {
		banner()

		menu := wmenu.NewMenu("\nSelect options.")
		menu.Action(func(opt wmenu.Opt) error { return f([]wmenu.Opt{opt}) })
		menu.MultipleAction(f)

		menu.Option("Fullscreen", "0", cfg.Fullscreen, nil)
		menu.Option("Widescreen", "1", cfg.Widescreen, nil)
		menu.Option("WC2Input", "2", cfg.WC2Input, nil)

		if menu.Run() == nil {
			saveConfig(cfg)
			return nil
		}
	}
}

func raceMenu() error {
	for {
		banner()

		menu := wmenu.NewMenu("\nSelect a race.")
		menu.Action(func(opt wmenu.Opt) error {
			cfg.Debug.Race = opt.Text
			return nil
		})

		menu.Option("Human", "0", true, nil)
		menu.Option("Orc", "1", false, nil)

		if menu.Run() == nil {
			saveConfig(cfg)
			return nil
		}
	}
}

func mapMenu() error {
	for {
		banner()

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

		if menu.Run() == nil {
			saveConfig(cfg)
			return nil
		}
	}
}

func installArchiveMenu() error {
	msg := "Do you want to download and install the content from\nWarcraft: Orcs & Humans shareware version?"
	for {
		banner()

		menu := wmenu.NewMenu(msg)
		menu.Action(func(opt wmenu.Opt) error {
			if opt.ID == 0 {
				dst := filepath.Clean(platform.CfgRootJoin("DATA.WAR"))
				if err := downloadAndExtract(dst, "https://sites.google.com/site/openwarengine/war1sw.zip?attredirects=0&d=1"); err != nil {
					clearScreen()
					msg = "Failed to download! Retry?"
					war = nil
					arch = notInstalledText
					return err
				}

				if archive, err := resource.OpenArchive(dst); err == nil {
					war = archive
					arch = dst
				}
			}
			return nil
		})

		menu.IsYesNo(0)
		if menu.Run() == nil {
			return nil
		}
	}
}

func downloadAndExtract(dst, src string) error {
	const timeout = 10 * time.Second

	clearScreen()
	fmt.Print("Contacting server...")

	downChan := make(chan byte, 4096)
	lenChan := make(chan int)

	go func() {
		defer close(downChan)
		defer close(lenChan)

		resp, err := http.Get(src)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		icl := int(resp.ContentLength)
		lenChan <- icl
		buf := bufio.NewReader(resp.Body)

		for n := 0; n < icl; n++ {
			b, err := buf.ReadByte()
			if err != nil {
				return
			}
			downChan <- b
		}
	}()

	var buf bytes.Buffer
	contentLength := -1
	downloadError := errors.New("could not download file")
	n := 0
	u := 0.0

	for {
		select {
		case contentLength = <-lenChan:
			lenChan = nil
		case b, ok := <-downChan:
			if !ok {
				downChan = nil
			}
			buf.WriteByte(b)
			n++
		case <-time.After(timeout):
			return downloadError
		}

		if n == contentLength {
			return extractZip(&buf, contentLength, dst)
		}

		p := float64(n) / float64(contentLength)
		if p > u {
			u = p + 0.01
			clearScreen()
			fmt.Printf("Downloading... %d%%", int(p*100.0))
		}
	}
}

func extractZip(buf *bytes.Buffer, size int, dst string) error {
	reader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(size))
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		if strings.ToUpper(filepath.Base(file.Name)) == "DATA.WAR" {
			in, err := file.Open()
			if err != nil {
				return err
			}
			defer in.Close()

			out, err := os.Create(dst)
			if err != nil {
				return err
			}
			defer out.Close()

			if _, err := io.Copy(out, in); err != nil {
				return err
			}

			return nil
		}
	}

	return errors.New("could not locate DATA.WAR in archive")
}
