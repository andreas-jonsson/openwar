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

package unit

import (
	"image"

	"github.com/andreas-jonsson/openwar/platform"
)

type UnitType int

const (
	BuildingType UnitType = iota
	CreatureType
	GoldMineType
)

type Unit interface {
	Id() uint64
	Type() UnitType
	Name() string
	SetPosition(pt image.Point)
	Position() image.Point
	Bounds() image.Rectangle
	Update() error
	Render(renderer platform.Renderer, cameraPos image.Point) error
}
