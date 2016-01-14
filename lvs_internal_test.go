// Copyright (c) 2016 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package lvs

var (
	fakeRunOutput       []byte
	fakeRunErr          error
	fakeExecuteErr      error
	fakeExecuteStdinErr error
)

func fakeRun(args []string) ([]byte, error) {
	// cmd := exec.Command(args[0], args[1:]...)
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	return nil, errors.New(err.Error() + " output: " + string(output))
	// }
	return fakeRunOutput, fakeRunErr
}

func fakeExecute(exe string, args ...string) error {
	// // fmt.Printf("%s\n", strings.Join(append([]string{exe}, args...), " "))
	// cmd := exec.Command(exe, args...)
	return fakeExecuteErr
}

func fakeExecuteStdin(in, exe string, args ...string) error {
	// var err error
	// var total, part, segment int
	// var stdin io.WriteCloser

	// cmd := exec.Command(exe, args...)
	// stdin, err = cmd.StdinPipe()
	// defer stdin.Close()
	// if err = cmd.Start(); err != nil {
	// 	return err
	// }

	// total = len(in)
	// for part = 0; part != total; part += segment {
	// 	segment, err = stdin.Write([]byte(in[part:total]))
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return fakeExecuteStdinErr
}
