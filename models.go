// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package lvs

import (
	"encoding/json"
	"fmt"
	// "strings"
)

type (
	ToJson interface {
		ToJson() ([]byte, error)
	}

	FromJson interface {
		FromJson([]byte) error
	}

	ider interface {
		getId() string
	}

	Vip struct {
		Host        string   `json:"host"`
		Port        int      `json:"port"`
		Schedular   string   `json:"schedular"`
		Persistance int      `json:"persistance"`
		Servers     []Server `json:"servers"`
	}
)

func (v *Vip) FromJson(bytes []byte) error {
	return json.Unmarshal(bytes, v)
}
func (v Vip) ToJson() ([]byte, error) {
	return json.Marshal(v)
}
func (v Vip) getId() string {
	return fmt.Sprintf("%v:%v", v.Host, v.Port)
}
