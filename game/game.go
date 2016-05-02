/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package game

import (
	"fmt"
	"log"

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

	renderer  platform.Renderer
	player    platform.AudioPlayer
	resources resource.Resources
}

func NewGame(rend platform.Renderer, player platform.AudioPlayer, res resource.Resources) *Game {
	g := new(Game)
	g.running = true

	g.renderer = rend
	g.player = player
	g.resources = res

	s := make(map[string]GameState)
	s["menu"] = NewMenuState(g)
	s["play"] = NewPlayState(g)

	g.states = s
	return g
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
			if t.Key == platform.KEY_SPACE {
				g.renderer.ToggleFullscreen()
			}
			return event
		default:
			return event
		}
	}
}

func (g *Game) CurrentState() GameState {
	return g.currentState
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

func (g *Game) Shutdown() {

}
