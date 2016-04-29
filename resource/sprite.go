/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

/*
Sprite sheet files start with a 2 byte integer telling the number of frames inside the file,
followed by the sprite dimensions as 1 byte width and height. Next is a list of all frames,
starting with their y and x offset, followed by width and height, each as 1 byte value.
Last comes the offset of the frame inside the file, stored as 4 byte integer.
*/

package resource

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color/palette"
	"path"
)

type (
	Frame struct {
		X, Y          byte
		Width, Height byte
		Offset        uint32
	}

	Sprite struct {
		NumFrames     uint16
		Width, Height byte

		Frames []Frame
		Imgs   []*image.Paletted
	}

	Sprites map[string]Sprite
)

func LoadSprites(arch *Archive, images Images) (Sprites, error) {
	sprites := make(map[string]Sprite)

	for file, data := range arch.Files {
		if path.Ext(file) == ".SPR" {
			reader, _ := arch.Open(file)

			var sprite Sprite
			if err := binary.Read(reader, binary.LittleEndian, &sprite.NumFrames); err != nil {
				return sprites, err
			}

			if err := binary.Read(reader, binary.LittleEndian, &sprite.Width); err != nil {
				return sprites, err
			}

			if err := binary.Read(reader, binary.LittleEndian, &sprite.Height); err != nil {
				return sprites, err
			}

			sprite.Frames = make([]Frame, sprite.NumFrames)
			sprite.Imgs = make([]*image.Paletted, sprite.NumFrames)

			for i := uint16(0); i < sprite.NumFrames; i++ {
				if err := binary.Read(reader, binary.LittleEndian, &sprite.Frames[i]); err != nil {
					return sprites, err
				}
			}

			var (
				img        *image.Paletted
				prevOffset uint32 = 0xFFFFFFFF
			)

			for i, frame := range sprite.Frames {
				if frame.Offset == prevOffset {

					// NOTE This is strange... could it be compressed?
					fmt.Fprintln(Logger, "Frame repeat:", file, i)

					img = sprite.Imgs[i-1]
				} else {
					rect := image.Rect(0, 0, int(frame.Width), int(frame.Height))
					img = image.NewPaletted(rect, palette.Plan9)
					pixReader := bytes.NewReader(data[frame.Offset:])

					if err := binary.Read(pixReader, binary.LittleEndian, &img.Pix); err != nil {
						return sprites, err
					}
				}

				prevOffset = frame.Offset
				sprite.Imgs[i] = img
				images[fmt.Sprintf("%s,%d", file, i)] = Image{X: int(frame.X), Y: int(frame.Y), Data: img}
			}

			sprites[file] = sprite
		}
	}

	return sprites, nil
}
