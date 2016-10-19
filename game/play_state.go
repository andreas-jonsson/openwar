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

package game

import (
	"image"
	"log"
	"time"

	"github.com/andreas-jonsson/openwar/game/unit"
	"github.com/andreas-jonsson/openwar/platform"
	"github.com/andreas-jonsson/openwar/resource"
)

const (
	scrollLeft = 1 << iota
	scrollRight
	scrollUp
	scrollDown
)

const scrollSpeed = 0.1

type playState struct {
	g *Game
	p *player

	scrollDirection  int
	cameraX, cameraY float64

	ter   terrain
	units *unit.Manager
	res   resource.Resources
}

func NewPlayState(g *Game) GameState {
	ter, err := newTerrain(g, g.config.Debug.Map)
	if err != nil {
		log.Panicln(err)
	}

	race := humanRace
	if g.config.Debug.Race == "Orc" {
		race = orcRace
	}

	unitManager := unit.NewManager(&g.resources, ter.terrainPalette())

	return &playState{
		g:     g,
		p:     newPlay(g, unitManager, race, ter.terrainPalette()),
		res:   g.resources,
		ter:   ter,
		units: unitManager,
	}
}

func (s *playState) Name() string {
	return "play"
}

func (s *playState) Enter(from GameState, args ...interface{}) error {
	s.g.musicPlayer.random(10 * time.Second)

	snd, _ := s.g.soundPlayer.Sound("OREADY.VOC")
	snd.Play(-1, 0, 0)

	s.units.SpawnBuilding("HumanTownhall", image.Point{10, 4})
	s.units.SpawnBuilding("HumanFarm", image.Point{10, 8})
	s.units.SpawnBuilding("HumanBarrack", image.Point{6, 7})
	s.units.SpawnBuilding("HumanTower", image.Point{13, 8})

	return nil
}

func (s *playState) Exit(to GameState) error {
	return nil
}

func (s *playState) Update() error {
	g := s.g
	g.PollAll()

	if pos, updateCamera := s.p.hud.handleMouse(platform.Mouse()); updateCamera {
		s.cameraX = float64(pos.X)
		s.cameraY = float64(pos.Y)
	}

	s.updateScroll(g.dt)
	return nil
}

func (s *playState) updateScroll(dt float64) {
	g := s.g
	pos := g.cursorPos
	max := g.renderer.BackBuffer().Bounds().Max
	s.scrollDirection = 0

	if pos.X <= 0 {
		s.scrollDirection |= scrollLeft
		s.cameraX -= dt * scrollSpeed
	} else if pos.X >= max.X-1 {
		s.scrollDirection |= scrollRight
		s.cameraX += dt * scrollSpeed
	}

	if pos.Y <= 0 {
		s.scrollDirection |= scrollUp
		s.cameraY -= dt * scrollSpeed
	} else if pos.Y >= max.Y-1 {
		s.scrollDirection |= scrollDown
		s.cameraY += dt * scrollSpeed
	}

	mapSize := s.ter.size()
	vp := s.p.hud.viewport()
	cameraPos := image.Point{int(s.cameraX), int(s.cameraY)}
	cameraMax := image.Point{mapSize*16 - (vp.Max.X - vp.Min.X), mapSize*16 - (vp.Max.Y - vp.Min.Y)}

	if cameraPos.X < 0 {
		cameraPos.X = 0
		s.cameraX = 0
	} else if cameraPos.X > cameraMax.X {
		cameraPos.X = cameraMax.X
		s.cameraX = float64(cameraPos.X)
	}

	if cameraPos.Y < 0 {
		cameraPos.Y = 0
		s.cameraY = 0
	} else if cameraPos.Y > cameraMax.Y {
		cameraPos.Y = cameraMax.Y
		s.cameraY = float64(cameraPos.Y)
	}

	s.setCursor()
}

func (s *playState) setCursor() {
	g := s.g
	g.currentCursor = cursorNormal

	switch {
	case s.scrollDirection == scrollUp|scrollRight:
		g.currentCursor = cursorScrollTopRight
	case s.scrollDirection == scrollDown|scrollRight:
		g.currentCursor = cursorScrollBottomRight
	case s.scrollDirection == scrollDown|scrollLeft:
		g.currentCursor = cursorScrollBottomLeft
	case s.scrollDirection == scrollUp|scrollLeft:
		g.currentCursor = cursorScrollTopLeft
	case s.scrollDirection == scrollUp:
		g.currentCursor = cursorScrollTop
	case s.scrollDirection == scrollRight:
		g.currentCursor = cursorScrollRight
	case s.scrollDirection == scrollDown:
		g.currentCursor = cursorScrollBottom
	case s.scrollDirection == scrollLeft:
		g.currentCursor = cursorScrollLeft
	}
}

func (s *playState) Render() error {
	mapSize := s.ter.size()
	vp := s.p.hud.viewport()
	cameraPos := image.Point{int(s.cameraX), int(s.cameraY)}
	cameraMax := image.Point{mapSize*16 - (vp.Max.X - vp.Min.X), mapSize*16 - (vp.Max.Y - vp.Min.Y)}

	if cameraPos.X < 0 {
		cameraPos.X = 0
		s.cameraX = 0
	} else if cameraPos.X > cameraMax.X {
		cameraPos.X = cameraMax.X
		s.cameraX = float64(cameraPos.X)
	}

	if cameraPos.Y < 0 {
		cameraPos.Y = 0
		s.cameraY = 0
	} else if cameraPos.Y > cameraMax.Y {
		cameraPos.Y = cameraMax.Y
		s.cameraY = float64(cameraPos.Y)
	}

	s.ter.render(vp, cameraPos)
	s.units.Render(s.g.renderer, vp, cameraPos)
	s.p.render(s.ter.miniMapImage(), cameraPos)
	return nil
}
