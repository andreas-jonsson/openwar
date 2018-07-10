// +build mobile

/*
Copyright (C) 2016-2018 Andreas T Jonsson

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
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/andreas-jonsson/openwar/game"
	"github.com/andreas-jonsson/openwar/platform"
	"github.com/andreas-jonsson/openwar/resource"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

func Start() {
	banner()
	app.Main(func(a app.App) {
		var (
			glctx  gl.Context
			sz     size.Event
			images *glutil.Images
			glimg  *glutil.Image
		)

		paintDoneChan := make(chan struct{})

		cfg := createNewConfig(false)
		cfg.Widescreen = false

		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					images = glutil.NewImages(glctx)
					glimg = images.NewImage(512, 512)
					platform.ExternalBackBuffer = glimg.RGBA

					file, size := openAssets()
					defer file.Close()

					if war, err := resource.OpenArchiveFrom(file, size); err == nil {
						go game.Start(cfg, war)
					} else {
						log.Panicln(err)
					}

					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					close(platform.PaintEventChan)
					close(platform.InputEventChan)

					glimg.Release()
					images.Release()
					glctx = nil
				}
			case size.Event:
				sz = e
				select {
				case platform.InputEventChan <- e:
				default:
				}
			case paint.Event:
				if glctx == nil || e.External {
					continue
				}

				glctx.ClearColor(0, 0, 0, 1)
				glctx.Clear(gl.COLOR_BUFFER_BIT)

				platform.PaintEventChan <- paintDoneChan
				glimg.Upload()
				<-paintDoneChan

				glimg.Draw(sz, geom.Point{0, 0}, geom.Point{512, 0}, geom.Point{0, 512}, glimg.RGBA.Bounds())
				a.Publish()
				a.Send(paint.Event{})
			case touch.Event, key.Event:
				select {
				case platform.InputEventChan <- e:
				default:
				}
			}
		}
	})
}

func openAssets() (asset.File, int64) {
	file, err := asset.Open("DATA.WAR")
	if err != nil {
		err = filepath.Walk("/storage/emulated/0/Download", func(path string, info os.FileInfo, err error) error {
			if info != nil && !info.IsDir() && strings.ToUpper(filepath.Base(path)) == "DATA.WAR" {
				log.Println("Found resources: " + path)
				return errors.New(path)
			}
			return nil
		})

		if err == nil {
			log.Panicln("Could not find DATA.WAR in your Download folder!")
		}

		file, err = os.Open(err.Error())
		if err != nil {
			log.Panicln(err)
		}
	}

	size, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		log.Panicln(err)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		log.Panicln(err)
	}

	return file, size
}
