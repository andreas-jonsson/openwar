/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package game

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/openwar-hq/openwar/platform"
	"github.com/openwar-hq/openwar/resource"
)

type GameState interface {
	Name() string
	Enter(from GameState, args ...interface{}) error
	Exit(to GameState) error
	Update() error
	Render() error
}

type Game struct {
	currentState GameState
	running      bool
	states       map[string]GameState

	currentCursor int
	cursorPos     image.Point
	cursors       []resource.Image
	cursorPal     color.Palette

	renderer    platform.Renderer
	resources   resource.Resources
	musicPlayer *musicPlayer
	soundPlayer platform.AudioPlayer

	t, ft time.Time
	dt    float64

	numFrames, fps int
}

func NewGame(rend platform.Renderer, player platform.AudioPlayer, res resource.Resources) (*Game, error) {
	var err error

	g := new(Game)
	g.running = true

	g.renderer = rend
	g.resources = res
	g.soundPlayer = player
	g.musicPlayer, err = newMusicPlayer(res.Archive, player)

	g.cursorPal = resource.ClonePalette(res.Palettes["CURSORS.PAL"])
	for i := 0; i < len(g.cursorPal); i++ {
		if cr, cg, cb, _ := g.cursorPal[i].RGBA(); cr == 0 && cg == 0 && cb == 0 {
			g.cursorPal[i] = color.RGBA{}
		}
	}

	g.cursors = []resource.Image{
		res.Images["NORMAL.CUR"],
		res.Images["NOPE.CUR"],
		res.Images["CROSHAIR.CUR"],
		res.Images["TARGET01.CUR"],
		res.Images["TARGET02.CUR"],
		res.Images["TARGET03.CUR"],
		res.Images["INSPECT.CUR"],
		res.Images["TIME.CUR"],

		res.Images["SCROLLT.CUR"],
		res.Images["SCROLLTR.CUR"],
		res.Images["SCROLLR.CUR"],
		res.Images["SCROLLBR.CUR"],
		res.Images["SCROLLB.CUR"],
		res.Images["SCROLLBL.CUR"],
		res.Images["SCROLLL.CUR"],
		res.Images["SCROLLTL.CUR"],
	}

	s := map[string]GameState{
		"menu": NewMenuState(g),
		"play": NewPlayState(g),
	}

	g.states = s
	return g, err
}

func (g *Game) PollEvent() platform.Event {
	for {
		event := platform.PollEvent()
		if event == nil {
			return nil
		}

		switch t := event.(type) {
		case *platform.QuitEvent:
			g.running = false
		case *platform.KeyDownEvent:
			switch t.Key {
			case platform.KeyEsc:
				g.running = false
				continue
			case platform.KeySpace:
				g.renderer.ToggleFullscreen()
			}
			return event
		case *platform.MouseMotionEvent:
			g.cursorPos = image.Point{t.X, t.Y}
			return event
		default:
			return event
		}
	}
}

func (g *Game) SwitchState(to string, args ...interface{}) error {
	newState, ok := g.states[to]
	if !ok {
		return fmt.Errorf("invalid state: %s", to)
	}

	if g.currentState != nil {
		log.Printf("Exiting state: %v", g.currentState.Name())
		if err := g.currentState.Exit(newState); err != nil {
			return nil
		}
	}

	log.Printf("Enter state: %v", to)
	if err := newState.Enter(g.currentState, args...); err != nil {
		return err
	}

	g.currentState = newState
	return nil
}

func (g *Game) Running() bool {
	return g.running
}

func (g *Game) Update() error {
	now := time.Now()
	g.dt = float64(now.Sub(g.t).Nanoseconds() / int64(time.Millisecond))
	g.t = now

	err := g.currentState.Update()

	g.numFrames++
	if time.Since(g.ft).Nanoseconds()/int64(time.Millisecond) >= 1000 {
		g.fps = g.numFrames
		g.ft = now
		g.numFrames = 0

		g.renderer.SetWindowTitle(fmt.Sprintf("OpenWar - %d fps", g.fps))
	}

	return err
}

func (g *Game) Render() error {
	if err := g.currentState.Render(); err != nil {
		return err
	}

	cur := g.cursors[g.currentCursor]
	g.renderer.BlitImage(g.cursorPos.Add(cur.Point()), cur.Data, g.cursorPal)
	return nil
}

func (g *Game) Shutdown() {

}
