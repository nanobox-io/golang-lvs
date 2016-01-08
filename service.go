package lvs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type (
	Service struct {
		Host        string   `json:"host"`
		Port        int      `json:"port"`
		Type        string   `json:"type"`
		Scheduler   string   `json:"scheduler"`
		Persistance int      `json:"persistance"`
		Netmask     string   `json:"netmask"`
		Servers     []Server `json:"servers"`
	}
)

var (
	ServiceTypeFlag = map[string]string{
		"tcp":    "-t",
		"udp":    "-u",
		"fwmark": "-f",
		"":       "-t", // default
	}

	ServiceSchedulerFlag = map[string]string{
		"rr":    "rr",
		"wrr":   "wrr",
		"lc":    "lc",
		"wlc":   "wlc",
		"lblc":  "lblc",
		"lblcr": "lblcr",
		"dh":    "dh",
		"sh":    "sh",
		"sed":   "sed",
		"nq":    "nq",
		"":      "wlc", // default
	}
)

func (s Service) FindServer(server Server) *Server {
	for i := range s.Servers {
		if s.Servers[i].Host == server.Host && s.Servers[i].Port == server.Port {
			return &s.Servers[i]
		}
	}
	return nil
}

func (s *Service) AddServer(server Server) error {
	s.Servers = append(s.Servers, server)
	return backend("ipvsadm", append([]string{"-a", ServiceTypeFlag[s.Type], s.getHostPort(), "-r"}, strings.Split(server.String(), " ")...)...)
}

func (s *Service) EditServer(server Server) error {
	for i := range s.Servers {
		if s.Servers[i].Host == server.Host && s.Servers[i].Port == server.Port {
			s.Servers = append(s.Servers[:i], append([]Server{server}, s.Servers[i+1:]...)...)
			break
		}
	}
	return backend("ipvsadm", append([]string{"-e", ServiceTypeFlag[s.Type], s.getHostPort(), "-r"}, strings.Split(server.String(), " ")...)...)
}

func (s *Service) RemoveServer(server Server) error {
	for i := range s.Servers {
		if s.Servers[i].Host == server.Host && s.Servers[i].Port == server.Port {
			s.Servers = append(s.Servers[:i], s.Servers[i+1:]...)
			break
		}
	}
	return backend("ipvsadm", "-d", ServiceTypeFlag[s.Type], s.getHostPort(), "-r", server.getHostPort())
}

func (s *Service) FromJson(bytes []byte) error {
	return json.Unmarshal(bytes, s)
}

func (s Service) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s Service) getId() string {
	return fmt.Sprintf("%v:%v", s.Host, s.Port)
}

func (s Service) getNetmask() string {
	if s.Netmask != "" {
		return fmt.Sprintf("-M %s", s.Netmask)
	} else {
		return ""
	}
}

func (s Service) getHostPort() string {
	if s.Port == 0 {
		return s.Host
	}
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s Service) String() string {
	a := make([]string, 0, 0)
	a = append(a, fmt.Sprintf("-A %s %s -s %s -p %d %s\n",
		ServiceTypeFlag[s.Type], s.getHostPort(),
		ServiceSchedulerFlag[s.Scheduler], s.Persistance, s.getNetmask()))
	for i := range s.Servers {
		a = append(a, fmt.Sprintf("-a %s %s:%d -r %s\n",
			ServiceTypeFlag[s.Type], s.Host, s.Port,
			s.Servers[i].String()))
	}
	return strings.Join(a, "")
}

func parseService(serviceString string) Service {
	service := Service{
		Scheduler:   "wlc",
		Type:        "tcp",
		Persistance: 300,
	}
	var err error
	exploded := strings.Split(serviceString, " ")
	for i := range exploded {
		switch exploded[i] {
		case "-t", "--tcp-service":
			service.Type = "tcp"
			service.Host, service.Port = parseHostPort(exploded[i+1])
		case "-u", "--udp-service":
			service.Type = "udp"
			service.Host, service.Port = parseHostPort(exploded[i+1])
		case "-f", "--fwmark-service":
			service.Type = "fwmark"
			service.Host, service.Port = parseHostPort(exploded[i+1])
		case "-s", "--scheduler":
			service.Scheduler = exploded[i+1]
		case "-p", "--persistent":
			service.Persistance, err = strconv.Atoi(exploded[i+1])
			if err != nil {
				service.Persistance = 300
			}
		case "-M", "--netmask":
			service.Netmask = exploded[i+1]
		}
	}
	return service
}
