/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package platform

const (
	KEY_UP = iota
	KEY_DOWN
	KEY_LEFT
	KEY_RIGHT
)

type (
	Event     interface{}
	QuitEvent struct{}

	KeyUpEvent struct {
		Key int
	}

	KeyDownEvent struct {
		Key int
	}
)
