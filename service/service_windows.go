package service

import (
	"github.com/shirou/gopsutil/winservices"
	"golang.org/x/sys/windows"
	"strings"
)

func GetService(pattern string) []*Service {
	var serviceList []*Service
	services, _ := winservices.ListServices()

	for _, service := range services {
		newService, err := winservices.NewService(service.Name)
		if err != nil {
			continue
		}

		err = newService.GetServiceDetail()
		if err != nil {
			continue
		}

		if strings.Contains(strings.ToLower(newService.Config.DisplayName), pattern) || pattern == "all" {
			serviceList = append(serviceList, getService(newService))
		}

	}

	return serviceList
}

func getService(s *winservices.Service) *Service {
	return &Service{
		Name:        s.Name,
		StartType:   getStartType(s.Config.StartType),
		ExecPath:    s.Config.BinaryPathName,
		DisplayName: s.Config.DisplayName,
		Description: s.Config.Description,
		State:       getStateType(uint32(s.Status.State)),
		Pid:         s.Status.Pid,
		ExitCode:    s.Status.Win32ExitCode,
	}
}

func getStartType(t uint32) string {
	switch t {
	case windows.SERVICE_AUTO_START:
		return "auto_start"
	case windows.SERVICE_BOOT_START:
		return "boot_start"
	case windows.SERVICE_DEMAND_START:
		return "demand_start"
	case windows.SERVICE_DISABLED:
		return "disabled"
	default:
		return "unknown"
	}
}

func getStateType(t uint32) string {
	switch t {
	case windows.SERVICE_STOPPED:
		return "stopped"
	case windows.SERVICE_START_PENDING:
		return "start_pending"
	case windows.SERVICE_STOP_PENDING:
		return "stop_pending"
	case windows.SERVICE_RUNNING:
		return "running"
	case windows.SERVICE_CONTINUE_PENDING:
		return "continue_pending"
	case windows.SERVICE_PAUSE_PENDING:
		return "pause_pending"
	case windows.SERVICE_PAUSED:
		return "paused"
	default:
		return "unknown"
	}
}
