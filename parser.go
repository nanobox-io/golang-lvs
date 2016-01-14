// Copyright (c) 2016 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package lvs

import (
	"errors"
	"net"
	"strconv"
)

var (
	EOFError       = errors.New("ipvsadm terminated prematurely")
	UnexpecedToken = errors.New("Unexpected Token")
)

func parseHostPort(hostPort string) (string, int) {
	host, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return hostPort, 0
	}
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return hostPort, 0
	}
	return host, intPort
}
