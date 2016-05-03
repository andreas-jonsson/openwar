/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package game

import "image/color"

type (
	playerRace      int
	environmentType int
)

const (
	humanRace playerRace = iota
	orcRace
)

const (
	environmentForest environmentType = iota
	environmentSwamp
	environmentDungeon
)

type player struct {
	g   *Game
	hud *gameHud
}

func newPlay(g *Game, race playerRace, envPal color.Palette) *player {
	return &player{g: g, hud: newGameHud(g, race, envPal)}
}

func (p *player) render() error {
	return p.hud.render()
}
