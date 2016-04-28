/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

/*
Sprite sheet files start with a 2 byte integer telling the number of frames inside the file,
followed by the sprite dimensions as 1 byte width and height. Next is a list of all frames,
starting with their y and x offset, followed by width and height, each as 1 byte value.
Last comes the offset of the frame inside the file, stored as 4 byte integer.
If the width times height is greater than the difference between this and the next
offset, then the frame is compressed as specified below. Else it is to be read as a usual
indexed 256 color bitmap.

Sprites of the Mac version are often compressed with an RLE method. Only transparency
is compressed and the compression is linewise.
Lines are assembled by reading them blockwise, where each block starts with a single
byte Head, that gives further instructions:

	0x00 → EndOfLine
	0xFF → EndOfFrame
	0x80 & Head = 0 → Head-many uncompressed pixels follow
	0x80 & Head != 0 → ((0x7F & Head) + 1)-many transparent pixels

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

			for i, frame := range sprite.Frames {
				nextOffset := uint32(len(data))
				if i < len(sprite.Frames)-1 {
					nextOffset = sprite.Frames[i+1].Offset
				}

				rect := image.Rect(0, 0, int(frame.Width), int(frame.Height))

				// Check if the frame is compressed.
				if uint32(frame.Width)*uint32(frame.Height) > nextOffset-frame.Offset {
					fmt.Fprintln(Logger, "Compressed frame:", file, i)

				} else {
					img := image.NewPaletted(rect, palette.Plan9)
					if err := binary.Read(bytes.NewReader(data[frame.Offset:]), binary.LittleEndian, &img.Pix); err != nil {
						return sprites, err
					}

					sprite.Imgs[i] = img
					images[fmt.Sprintf("%s,%d", file, i)] = img
				}
			}

			sprites[file] = sprite
		}
	}

	return sprites, nil
}
