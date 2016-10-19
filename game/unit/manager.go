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

package unit

import (
	"encoding/json"
	"image"
	"image/color"
	"io/ioutil"
	"log"

	"github.com/andreas-jonsson/openwar/data"
	"github.com/andreas-jonsson/openwar/platform"
	"github.com/andreas-jonsson/openwar/resource"
)

type (
	UnitConfigs struct {
		Buildings map[string]BuildingConfig
	}

	Manager struct {
		terrainPal  color.Palette
		resources   *resource.Resources
		units       []Unit
		unitConfigs UnitConfigs
	}
)

func NewManager(res *resource.Resources, terrainPal color.Palette) *Manager {
	fp, err := data.FS.Open("units.json")
	if err != nil {
		log.Panicln(err)
	}
	defer fp.Close()

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		log.Panicln(err)
	}

	mgr := &Manager{terrainPal: terrainPal, resources: res}
	if err := json.Unmarshal(data, &mgr.unitConfigs); err != nil {
		log.Panicln(err)
	}
	return mgr
}

func (mgr *Manager) SpawnBuilding(name string, pos image.Point) Unit {
	pal := resource.BlackToAlpha(resource.CombinePalettes(mgr.terrainPal, mgr.resources.Palettes["SPRITE0.PAL"]))
	cfg := mgr.unitConfigs.Buildings[name]
	unit := NewBuilding(name, &cfg, mgr.resources, pal)
	unit.SetTilePosition(pos)

	mgr.units = append(mgr.units, unit)
	return unit
}

func (mgr *Manager) Render(renderer platform.Renderer, vp image.Rectangle, cameraPos image.Point) error {
	origo := vp.Min.Sub(cameraPos)
	for _, u := range mgr.units {
		if err := u.Render(renderer, origo); err != nil {
			return err
		}
	}
	return nil
}

func (mgr *Manager) AllUnits() []Unit {
	return mgr.units
}
