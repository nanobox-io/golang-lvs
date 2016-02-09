// Copyright (c) 2016 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package lvs

import (
	"fmt"
	"strings"
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

func (i Ipvs) FindService(netType, host string, port int) *Service {
	for j := range i.Services {
		if i.Services[j].Host == host && i.Services[j].Port == port && i.Services[j].Type == netType {
			return &i.Services[j]
		}
	}
	return nil
}

func (i *Ipvs) AddService(service Service) error {
	err := service.Validate()
	if err != nil {
		return err
	}
	if i.FindService(service.Type, service.Host, service.Port) != nil {
		return nil
	}
	err = backend("ipvsadm", append([]string{"-A", ServiceTypeFlag[service.Type], service.getHostPort(), "-s", ServiceSchedulerFlag[service.Scheduler]}, append(service.getPersistence(), service.getNetmask()...)...)...)
	if err != nil {
		return err
	}
	for i := range service.Servers {
		err := backend("ipvsadm", append([]string{"-a", ServiceTypeFlag[service.Type], service.getHostPort(), "-r"}, strings.Split(service.Servers[i].String(), " ")...)...)
		if err != nil {
			return err
		}
	}
	i.Services = append(i.Services, service)
	return nil
}

func (i *Ipvs) EditService(service Service) error {
	err := backend("ipvsadm", append([]string{"-E", ServiceTypeFlag[service.Type], service.getHostPort(), ServiceSchedulerFlag[service.Scheduler]}, append(service.getPersistence(), service.getNetmask()...)...)...)
	if err != nil {
		return err
	}

	for j := range i.Services {
		if i.Services[j].Host == service.Host && i.Services[j].Port == service.Port && i.Services[j].Type == service.Type {
			i.Services = append(i.Services[:j], append([]Service{service}, i.Services[j+1:]...)...)
			break
		}
	}
	return nil
}

func (i *Ipvs) RemoveService(netType, host string, port int) error {
	err := backend("ipvsadm", "-D", ServiceTypeFlag[netType], fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}

	for j := range i.Services {
		if i.Services[j].Host == host && i.Services[j].Port == port && i.Services[j].Type == netType {
			i.Services = append(i.Services[:j], i.Services[j+1:]...)
			break
		}
	}
	return nil
}

func (i *Ipvs) Clear() error {
	err := backend("ipvsadm", "-C")
	if err != nil {
		return err
	}

	i.Services = make([]Service, 0, 0)
	return nil
}

func (i Ipvs) SetTimeouts() error {
	if i.Tcp > 0 || i.Tcpfin > 0 || i.Udp > 0 {
		return backend("ipvsadm", "--set", string(i.Tcp), string(i.Tcpfin), string(i.Udp))
	}
	return nil
}

func (i *Ipvs) Restore(services []Service) error {
	in := make([]string, 0, 0)
	for i := range services {
		in = append(in, services[i].String())
	}
	err := backendStdin(strings.Join(in, ""), "ipvsadm", "-R")
	if err != nil {
		return err
	}

	i.Services = services
	return nil
}

// save reads the applied ipvsadm rules from the host and saves them as i.Services
func (i *Ipvs) Save() error {
	out, err := backendRun([]string{"ipvsadm", "-S", "-n"})
	if err != nil {
		return err
	}

	i.Services = make([]Service, 0, 0)
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

func (i Ipvs) Zero() error {
	return backend("ipvsadm", "-Z")
}
