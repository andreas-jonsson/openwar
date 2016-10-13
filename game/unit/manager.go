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
	"image"
	"image/color"

	"github.com/andreas-jonsson/openwar/game/unit/building"
	"github.com/andreas-jonsson/openwar/platform"
	"github.com/andreas-jonsson/openwar/resource"
)

type Unit interface {
	building.Building
}

type Manager struct {
	units []Unit
}

func NewManager(res *resource.Resources, terrainPal color.Palette) *Manager {
	pal := resource.BlackToAlpha(resource.CombinePalettes(terrainPal, res.Palettes["SPRITE0.PAL"]))
	testBuilding := building.NewHumanTownhall(res, pal)
	testBuilding.SetTilePosition(image.Point{10, 4})

	return &Manager{[]Unit{testBuilding}}
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
