// Copyright (c) 2016 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package lvs

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type (
	Server struct {
		Host           string `json:"host"`
		Port           int    `json:"port"`
		Forwarder      string `json:"forwarder"`
		Weight         int    `json:"weight"`
		UpperThreshold int    `json:upper_threshold`
		LowerThreshold int    `json:lower_threshold`
	}
)

var (
	ServerForwarderFlag = map[string]string{
		"g": "-g",
		"i": "-i",
		"m": "-m",
		"":  "-g", // default
	}

	InvalidServerForwarder = errors.New("Invalid Server Forwarder")
	InvalidServerPort      = errors.New("Invalid Server Port for Forwarder")
)

func (s Server) Validate() error {
	_, ok := ServerForwarderFlag[s.Forwarder]
	if !ok {
		return InvalidServerForwarder
	}
	return nil
}

func (s *Server) FromJson(bytes []byte) error {
	return json.Unmarshal(bytes, s)
}

func (s Server) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s Server) getHostPort() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s Server) String() string {
	return fmt.Sprintf("%s %s -y %d -x %d -w %d",
		s.getHostPort(), ServerForwarderFlag[s.Forwarder],
		s.LowerThreshold, s.UpperThreshold, s.Weight)
}

func parseServer(serverString string) Server {
	server := Server{
		Forwarder: "g",
		Weight:    1,
	}
	var err error
	exploded := strings.Split(serverString, " ")
	for i := range exploded {
		switch exploded[i] {
		case "-r", "--real-server":
			server.Host, server.Port = parseHostPort(exploded[i+1])
		case "-g", "--gatewaying":
			server.Forwarder = "g"
		case "-i", "--ipip":
			server.Forwarder = "i"
		case "-m", "--masquerading":
			server.Forwarder = "m"
		case "-w", "--weight":
			server.Weight, err = strconv.Atoi(exploded[i+1])
			if err != nil {
				server.Weight = 1
			}
		case "-x", "--u-threshold":
			server.UpperThreshold, err = strconv.Atoi(exploded[i+1])
			if err != nil {
				server.UpperThreshold = 0
			}
		case "-y", "--l-threshold":
			server.LowerThreshold, err = strconv.Atoi(exploded[i+1])
			if err != nil {
				server.LowerThreshold = 0
			}
		}
	}
	return server
}
