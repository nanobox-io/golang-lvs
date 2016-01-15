// Copyright (c) 2016 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package lvs

import (
	"fmt"
	// "os/exec"
	"strings"
	// "encoding/json"
)

type (
	Ipvs struct {
		MulticastInterface string    `json:mcast_interface`
		Syncid             int       `json:syncid`
		Tcp                int       `json:tcp_timeout`
		Tcpfin             int       `json:tcp_fin_timeout`
		Udp                int       `json:udp_fin_timeout`
		Services           []Service `json:services`
	}
)

var (
	DefaultIpvs = &Ipvs{}
)

func (i Ipvs) FindService(service Service) *Service {
	for j := range i.Services {
		if i.Services[j].Host == service.Host && i.Services[j].Port == service.Port && i.Services[j].Type == service.Type {
			return &i.Services[j]
		}
	}
	return nil
}

func (i *Ipvs) AddService(service Service) error {
	testService := i.FindService(service)
	if testService != nil {
		return nil
	}
	i.Services = append(i.Services, service)
	return backend("ipvsadm", append([]string{"-A", ServiceTypeFlag[service.Type], service.getHostPort(),
		"-s", ServiceSchedulerFlag[service.Scheduler],
		"-p", fmt.Sprintf("%d", service.Persistence)}, strings.Split(service.getNetmask(), "")...)...)
}

func (i *Ipvs) EditService(service Service) error {
	for j := range i.Services {
		if i.Services[j].Host == service.Host && i.Services[j].Port == service.Port && i.Services[j].Type == service.Type {
			i.Services = append(i.Services[:j], append([]Service{service}, i.Services[j+1:]...)...)
			break
		}
	}
	return backend("ipvsadm", append([]string{"-E", ServiceTypeFlag[service.Type], service.getHostPort(),
		"-s", ServiceSchedulerFlag[service.Scheduler],
		"-p", fmt.Sprintf("%d", service.Persistence)}, strings.Split(service.getNetmask(), "")...)...)
}

func (i *Ipvs) RemoveService(service Service) error {
	for j := range i.Services {
		if i.Services[j].Host == service.Host && i.Services[j].Port == service.Port && i.Services[j].Type == service.Type {
			i.Services = append(i.Services[:j], i.Services[j+1:]...)
			break
		}
	}
	return backend("ipvsadm", "-D", ServiceTypeFlag[service.Type], service.getHostPort())
}

func (i *Ipvs) Clear() error {
	i.Services = make([]Service, 0, 0)
	return backend("ipvsadm", "-C")
}

func (i Ipvs) SetTimeouts() error {
	if i.Tcp > 0 || i.Tcpfin > 0 || i.Udp > 0 {
		return backend("ipvsadm", "--set", string(i.Tcp), string(i.Tcpfin), string(i.Udp))
	}
	return nil
}

func (i *Ipvs) Restore(services []Service) error {
	i.Services = services
	var in []string
	in = make([]string, 0, 0)
	for i := range services {
		in = append(in, services[i].String())
	}
	return backendStdin(strings.Join(in, ""), "ipvsadm", "-R")
}

func (i *Ipvs) Save() error {
	i.Services = make([]Service, 0, 0)
	out, err := backendRun([]string{"ipvsadm", "-S", "-n"})
	if err != nil {
		return err
	}
	serviceStrings := strings.Split(string(out), "-A")
	for j := range serviceStrings {
		if serviceStrings[j] == "" {
			continue
		}
		serverStrings := strings.Split(serviceStrings[j], "-a")
		serviceString := serverStrings[0]
		serverStrings = serverStrings[1:]
		// fmt.Println("Service: ", serviceString)
		service := parseService(serviceString)
		for k := range serverStrings {
			// fmt.Println("Server: ", serverStrings[j])
			server := parseServer(serverStrings[k])
			service.Servers = append(service.Servers, server)
		}
		i.Services = append(i.Services, service)
	}
	return nil
}

func (i Ipvs) StartDaemon() (error, error) {
	if i.MulticastInterface != "" {
		var err1, err2 error
		if i.Syncid > 0 {
			err1 = backend("ipvsadm", "--start-daemon", "master", "--mcast-interface", i.MulticastInterface, "--syncid", string(i.Syncid))
			err2 = backend("ipvsadm", "--start-daemon", "backup", "--mcast-interface", i.MulticastInterface, "--syncid", string(i.Syncid))
		} else {
			err1 = backend("ipvsadm", "--start-daemon", "master", "--mcast-interface", i.MulticastInterface)
			err2 = backend("ipvsadm", "--start-daemon", "backup", "--mcast-interface", i.MulticastInterface)
		}
		return err1, err2
	}
	return nil, nil
}

func (i Ipvs) StopDaemon() (error, error) {
	if i.MulticastInterface != "" {
		var err1, err2 error
		err1 = backend("ipvsadm", "--stop-daemon", "master")
		err2 = backend("ipvsadm", "--stop-daemon", "backup")
		return err1, err2
	}
	return nil, nil
}
