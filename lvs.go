// Copyright (c) 2016 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package lvs

import (
	"errors"
	"io"
	"os/exec"
)

var (
	Conflict       = errors.New("object already exists")
	NotFound       = errors.New("object was not found")
	DeleteFailed   = errors.New("object was not deleted")
	IpvsadmMissing = errors.New("unable to find the ipvsadm command on the system")

	// these are to allow a pluggable backend for testing, ipvsadm is
	// not needed to run the tests
	backend      = execute
	backendRun   = run
	backendStdin = executeStdin
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

func SetTimeouts() error {
	return DefaultIpvs.SetTimeouts()
}

func StartDaemon() (error, error) {
	return DefaultIpvs.StartDaemon()
}

func StopDaemon() (error, error) {
	return DefaultIpvs.StopDaemon()
}

func Clear() error {
	return DefaultIpvs.Clear()
}

func Restore(services []Service) error {
	return DefaultIpvs.Restore(services)
}

func Save() error {
	return DefaultIpvs.Save()
}

func Zero() error {
	return DefaultIpvs.Zero()
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
	// fmt.Printf("%s\n", strings.Join(append([]string{exe}, args...), " "))
	cmd := exec.Command(exe, args...)
	return cmd.Run()
}

func executeStdin(in, exe string, args ...string) error {
	var err error
	var total, part, segment int
	var stdin io.WriteCloser

	cmd := exec.Command(exe, args...)
	stdin, err = cmd.StdinPipe()
	defer stdin.Close()
	if err = cmd.Start(); err != nil {
		return err
	}

	total = len(in)
	for part = 0; part != total; part += segment {
		segment, err = stdin.Write([]byte(in[part:total]))
		if err != nil {
			return err
		}
	}
	stdin.Close()
	return cmd.Wait()
}
