/*
Copyright (C) 2016-2017 Andreas T Jonsson

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

package network

import (
	"encoding/binary"
	"errors"
	"hash/fnv"
	"io"
	"log"
	"reflect"
)

const (
	ObjectUpdateMessage byte = iota
)

type (
	objectSync struct {
		obj        Object
		syncFields []int
	}
)

type (
	Constructor func(uint64) (Object, error)

	Object interface {
		Id() uint64
	}

	Manager struct {
		objects      map[uint64]objectSync
		constructors map[uint32]Constructor
		conn         io.ReadWriter

		msgReadyChan chan byte
		nextMsgChan  chan struct{}
	}
)

func hash32(data []byte) uint32 {
	h := fnv.New32a()
	h.Write(data)
	return h.Sum32()
}

func NewNetworkManager(conn io.ReadWriter) *Manager {
	mgr := &Manager{
		conn:         conn,
		msgReadyChan: make(chan byte),
		nextMsgChan:  make(chan struct{}),
	}

	go func() {
		for {
			var msgTy byte
			if err := binary.Read(mgr.conn, binary.LittleEndian, &msgTy); err != nil {
				close(mgr.msgReadyChan)
				return
			}

			mgr.msgReadyChan <- msgTy
			if _, ok := <-mgr.nextMsgChan; !ok {
				return
			}
		}
	}()

	return mgr
}

func (mgr *Manager) RegisterInstance(obj Object) {
	id := obj.Id()
	if _, ok := mgr.objects[id]; ok {
		log.Panicln("object is already registered")
	}
	mgr.objects[id] = objectSync{obj: obj, syncFields: findSyncFields(obj)}
}

func (mgr *Manager) RegisterConstructor(className string, f Constructor) {
	class := hash32([]byte(className))
	if _, ok := mgr.constructors[class]; ok {
		log.Panicln("constructor is already registered")
	}
	mgr.constructors[class] = f
}

func (mgr *Manager) Update() error {
	for _, objSync := range mgr.objects {
		msg := ObjectUpdateMessage
		if err := binary.Write(mgr.conn, binary.LittleEndian, msg); err != nil {
			return err
		}
		if err := binary.Write(mgr.conn, binary.LittleEndian, objSync.obj.Id()); err != nil {
			return err
		}

		t := reflect.TypeOf(objSync.obj).Elem()
		class := hash32([]byte(t.Name()))

		if err := binary.Write(mgr.conn, binary.LittleEndian, class); err != nil {
			return err
		}

		v := reflect.ValueOf(objSync.obj)
		for _, id := range objSync.syncFields {
			if err := binary.Write(mgr.conn, binary.LittleEndian, v.Field(id)); err != nil {
				return err
			}
		}
	}

	for {
		select {
		case msgTy := <-mgr.msgReadyChan:
			switch msgTy {
			case ObjectUpdateMessage:
				var (
					id    uint64
					class uint32
				)

				if err := binary.Read(mgr.conn, binary.LittleEndian, &id); err != nil {
					return err
				}
				if err := binary.Read(mgr.conn, binary.LittleEndian, &class); err != nil {
					return err
				}

				objSync, ok := mgr.objects[id]
				if !ok {
					f, ok := mgr.constructors[class]
					if !ok {
						return errors.New("invalid object id or unregistered constructor")
					}

					obj, err := f(id)
					if err != nil {
						return err
					}

					mgr.RegisterInstance(obj)
					objSync = mgr.objects[id]
				}

				v := reflect.ValueOf(objSync.obj)
				for _, id := range objSync.syncFields {
					if err := binary.Read(mgr.conn, binary.LittleEndian, v.Field(id).Interface()); err != nil {
						return err
					}
				}
			default:
				return errors.New("unknown message type")
			}
			mgr.nextMsgChan <- struct{}{}
		default:
			return nil
		}
	}
}

func findSyncFields(obj Object) []int {
	var fields []int

	t := reflect.TypeOf(obj)
	for i := 0; i < t.NumField(); i++ {
		networkTag := t.Field(i).Tag.Get("network")
		if networkTag == "sync" {
			fields = append(fields, i)
		}
	}

	return fields
}
