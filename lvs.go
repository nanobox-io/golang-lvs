// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package lvs

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
)

var (
	Conflict       = errors.New("object already exists")
	NotFound       = errors.New("object was not found")
	DeleteFailed   = errors.New("object was not deleted")
	IpvsadmMissing = errors.New("unable to find the ipvsadm command on the system")

	// these are to allow a pluggable backend for testing, ipvsadm is
	// not needed to run the tests
	backend    = execute
	backendRun = run
)

// Load verifies that lvs can be used, and populates it with values
// from the backup file
func Load() error {
	if err := check(); err != nil {
		return err
	}

	// NYI
	// populate the ipvsadm command with what was stored in the backup
	return nil
}

func check() error {
	if err := backend("which", "ipvsadm"); err != nil {
		return IpvsadmMissing
	}
	return nil
}

// Get a list of all Vips on the system
func ListVips() ([]Vip, error) {
	return parse(parseAll, "ipvsadm", "-ln")
}

//Add a Vip to the system
func AddVip(host string, port int) (*Vip, error) {
	id := fmt.Sprintf("%v:%v", host, port)
	// check if it already exists, this also validates the id
	vip, err := GetVip(id)
	if vip != nil {
		return vip, Conflict
	} else if err != NotFound {
		return nil, err
	}

	// create the vip
	if err := backend("ipvsadm", "-A", "-t", id, "-s", "wrr", "-p", "60"); err != nil {
		return nil, err // should be a custom error. this one may not make sense
	}

	backup()

	// double check that it was created
	return GetVip(id)
}

//Get a Vip on the system, or nil, NotFound if it is not found
func GetVip(id string) (*Vip, error) {
	if err := validateId(id); err != nil {
		return nil, err
	}
	vips, err := ListVips()
	if err != nil {
		return nil, err
	}

	for _, vip := range vips {
		if vip.getId() == id {
			return &vip, nil
		}
	}
	return nil, NotFound
}

//Remove a Vip from the system
func DeleteVip(id string) error {
	_, err := GetVip(id)
	if err == NotFound {
		return NotFound
	} else if err != nil {
		return err
	}

	if err := backend("ipvsadm", "-D", "-t", id); err != nil {
		return err // I should return my own error here
	}

	_, err = GetVip(id)
	if err != NotFound {
		return err
	} else if err == nil {
		return DeleteFailed
	}

	return nil
}

//Get backend destination servers for the specified Vip
func ListServers(vid string) ([]Server, error) {
	vip, err := GetVip(vid)
	if err != nil {
		return nil, err
	}
	return vip.Servers, nil
}

//Add a backend destination server for the specified Vip
func AddServer(vid, host string, port int) (*Server, error) {
	id := fmt.Sprintf("%v:%v", host, port)
	server, err := GetServer(vid, id)
	if server != nil {
		return server, Conflict
	} else if err != NotFound {
		return nil, err
	}

	if err := backend("ipvsadm", "-a", "-t", vid, "-r", id, "-w", "100", "-m"); err != nil {
		return nil, err // I should return my own error here
	}

	backup()
	return GetServer(vid, id)
}

//Get a backend server that is a memner of the specified Vip
func GetServer(vid, id string) (*Server, error) {
	if err := validateId(id); err != nil {
		return nil, err
	}
	servers, err := ListServers(vid)
	if err != nil {
		return nil, err
	}
	for _, server := range servers {
		if server.getId() == id {
			return &server, nil
		}
	}
	return nil, NotFound
}

//Disabled servers continue to serve current connections, but no new
//connections are sent to it. THis is exposed to allow draining of a
//backend server before removing it from the pool completely.
func EnableServer(vid, id string, enable bool) error {
	if _, err := GetServer(vid, id); err != nil {
		return err
	}

	var weight string
	if enable {
		weight = "100"
	} else {
		weight = "0"
	}

	if err := backend("ipvsadm", "-e", "-t", vid, "-r", id, "-w", weight); err != nil {
		return err // I should return my own error here
	}

	backup()
	return nil
}

//Remove a backend server from a Vip.
func DeleteServer(vid, id string) error {
	if _, err := GetServer(vid, id); err != nil {
		return err
	}

	if err := backend("ipvsadm", "-d", "-t", vid, "-r", id); err != nil {
		return err // I should return my own error here
	}

	backup()
	return nil
}

func parse(fun func(*bufio.Scanner) ([]Vip, error), args ...string) ([]Vip, error) {
	output, err := backendRun(args)
	if err != nil {
		return nil, err
	}
	pipe := ioutil.NopCloser(bytes.NewReader(output))
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanWords)
	vips, err := fun(scanner)
	if err != nil {
		return []Vip{}, errors.New("failed to parse: " + string(output))
	}
	return vips, err
}

func run(args []string) ([]byte, error) {
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.New(err.Error() + " output: " + string(output))
	}
	return output, err
}

func execute(exe string, args ...string) error {
	cmd := exec.Command(exe, args...)
	return cmd.Run()
}

func validateId(id string) error {
	_, _, err := net.SplitHostPort(id)
	return err
}

func backup() {
	//NYI
}
