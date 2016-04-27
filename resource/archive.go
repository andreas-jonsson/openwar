/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package resource

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
)

var (
	dosRetail    = [...]byte{0x18, 0x0, 0x0, 0x0}
	dosShareware = [...]byte{0x19, 0x0, 0x0, 0x0}
	macRetail    = [...]byte{0x0, 0x0, 0x0, 0x1A}
	macShareware = [...]byte{0x0, 0x0, 0x0, 0x19}
)

var Logger = os.Stdout

type Archive struct {
	Files map[string][]byte
}

func (a *Archive) Open(file string) (io.Reader, error) {
	if f, ok := a.Files[file]; ok {
		return bytes.NewReader(f), nil
	}
	return nil, os.ErrNotExist
}

func OpenArchive(file string) (*Archive, error) {
	fp, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	info, err := fp.Stat()
	if err != nil {
		return nil, err
	}

	var archiveID [4]byte
	if _, err = fp.Read(archiveID[:]); err != nil {
		return nil, err
	}

	fmt.Fprint(Logger, "Archive ID: ")
	switch archiveID {
	case dosRetail:
		fmt.Fprintln(Logger, "DOS Retail")
	default:
		switch archiveID {
		case dosShareware:
			fmt.Fprintln(Logger, "DOS Shareware")
		case macRetail:
			fmt.Fprintln(Logger, "Mac Retail")
		case macShareware:
			fmt.Fprintln(Logger, "Mac Shareware")
		default:
			return nil, errors.New("Unknown version!")
		}
		return nil, errors.New("Unsupported version!")
	}

	var numFiles uint32
	if err = binary.Read(fp, binary.LittleEndian, &numFiles); err != nil {
		return nil, err
	}
	fmt.Fprintln(Logger, "Number of files in archive: ", numFiles)

	if int(numFiles) != len(fileMap) {
		return nil, errors.New("Table mapping mismatch!")
	}

	fileTable := make([]uint32, numFiles)
	for i := range fileTable {
		if err = binary.Read(fp, binary.LittleEndian, &fileTable[i]); err != nil {
			return nil, err
		}
	}

	arch := &Archive{make(map[string][]byte)}

	for i, offset := range fileTable {
		if isPlaceHolder(fileTable, offset, i) {
			if fileMap[i] != "" {
				fmt.Fprintf(Logger, "Incomplete WAR file. Missing '%v'.\n", fileMap[i])
			}

			fmt.Fprintln(Logger, "Skipping placeholder: ", i)
			continue
		}

		if _, err = fp.Seek(int64(offset), 0); err != nil {
			return nil, err
		}

		var size uint32
		if err = binary.Read(fp, binary.LittleEndian, &size); err != nil {
			return nil, err
		}

		isCompressed := size>>24 == 0x20
		size &= 0x00FFFFFF

		var dataLength uint32
		if i == len(fileTable)-1 {
			dataLength = uint32(info.Size()) - fileTable[i]
		} else {
			dataLength = fileTable[i+1] - fileTable[i]
		}
		dataLength -= 4

		fileName := fileMap[i]
		if fileName == "" {
			fmt.Fprintf(Logger, "Warning: Filename table is incomplete! Missing file with id %v.\n", i)
			fileName = fmt.Sprintf("%s.%v", path.Base(fp.Name()), i)
		}

		var data []byte
		if isCompressed {
			if data, err = uncompressData(fp, int(size), int(dataLength)); err != nil {
				return nil, err
			}
		} else {
			data = make([]byte, size)
			if num, err := fp.Read(data); num != len(data) || err != nil {
				return nil, err
			}
		}

		arch.Files[fileName] = data
	}

	return arch, nil
}

func isPlaceHolder(tab []uint32, offset uint32, i int) bool {
	if offset == 0x0 || offset == 0xFFFFFFFF {
		return true
	}

	// Perhaps we should use the archive size?
	if i == len(tab)-1 {
		return false
	}

	if offset == (tab[i+1] - 1) {
		return true
	}
	return false
}

func readByte(reader io.Reader) (byte, error) {
	var b [1]byte
	if n, err := reader.Read(b[:]); n != 1 || err != nil {
		return 0, err
	}
	return b[0], nil
}

func readShort(reader io.Reader) (uint16, error) {
	var short uint16
	if err := binary.Read(reader, binary.LittleEndian, &short); err != nil {
		return 0, err
	}
	return short, nil
}

func uncompressData(reader io.Reader, fileSize, dataSize int) ([]byte, error) {
	const bufferSize = 4096
	var backingBuffer bytes.Buffer

	writer := bufio.NewWriter(&backingBuffer)
	buffer := make([]byte, bufferSize)

	var (
		numWrite,
		numRead int
	)

	for numRead < dataSize {
		cmask, err := readByte(reader)
		numRead++

		if err != nil {
			return buffer, err
		}

		for i := 0; i < 8 && numWrite != fileSize; i++ {
			if cmask%2 == 1 { // uncompressed
				bufByte, err := readByte(reader)
				numRead++

				if err != nil {
					return buffer, err
				}

				buffer[numWrite%bufferSize] = bufByte
				writer.WriteByte(bufByte)
				numWrite++
			} else { // compressed
				offset, err := readShort(reader)
				numRead += 2

				if err != nil {
					return buffer, err
				}

				numBytes := offset / bufferSize
				offset %= bufferSize

				for m := uint16(0); m <= numBytes+2; m++ {
					bufByte := buffer[(offset+m)%bufferSize]
					buffer[numWrite%bufferSize] = bufByte

					writer.WriteByte(bufByte)
					numWrite++
				}
			}
			cmask /= 2
		}
	}

	writer.Flush()
	return backingBuffer.Bytes(), nil
}
