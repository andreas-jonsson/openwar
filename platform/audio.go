/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package platform

type AudioPlayer interface {
	LoadMusic(name string, data []byte) error
	LoadSound(name string, data []byte) error
	Shutdown()
}
