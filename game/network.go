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
	"net"

	"github.com/andreas-jonsson/openwar/network"
)

func setupNetwork(addrs []net.TCPAddr, listenPort int) (*network.Manager, error) {
	/*
		l, err := net.Listen("tcp4", fmt.Sprintf(":%d", listenPort))
		if err != nil {
			return nil, err
		}

		incommingChan := make(chan net.Conn, len(addrs))
		go func() {
			for _, addr := range addrs {
				conn, err := l.Accept()
				if err != nil {
					close(incommingChan)
					return
				}
				incommingChan <- conn
			}
		}()

		for _, addr := range addrs {

		}
	*/
	return nil, nil
}
