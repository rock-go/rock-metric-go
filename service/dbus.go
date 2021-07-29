//+build linux

package service

import (
	"context"
	"fmt"
	"github.com/coreos/go-systemd/v22/dbus"
	"path/filepath"
)

type Properties struct {
	ExecMainCode   int32
	ExecMainStatus int32
	ExecMainPID    int32
	// accounting
	CPUAccounting    bool
	MemoryAccounting bool
	TasksAccounting  bool
	IPAccounting     bool
	// metrics
	CPUUsageNSec     int64
	MemoryCurrent    int64
	TasksCurrent     int64
	IPIngressPackets int64
	IPIngressBytes   int64
	IPEgressPackets  int64
	IPEgressBytes    int64
	// timestamps
	ActiveEnterTimestamp   uint64
	InactiveEnterTimestamp uint64
	InactiveExitTimestamp  uint64
	ActiveExitTimestamp    uint64
	// Meta
	FragmentPath string
}

type unitFetcher func(conn *dbus.Conn, states, patterns []string) ([]dbus.UnitStatus, error)

func listUnitsByPatternWrapper(conn *dbus.Conn, states, patterns []string) ([]dbus.UnitStatus, error) {
	return conn.ListUnitsByPatternsContext(context.Background(), states, patterns)
}

func listUnitsFilteredWrapper(conn *dbus.Conn, states, patterns []string) ([]dbus.UnitStatus, error) {
	units, err := conn.ListUnitsFilteredContext(context.Background(), states)
	if err != nil {
		return nil, fmt.Errorf("ListUnitsFiltered error")
	}

	return matchUnitPatterns(patterns, units)
}

func listUnitsWrapper(conn *dbus.Conn, states, patterns []string) ([]dbus.UnitStatus, error) {
	units, err := conn.ListUnitsContext(context.Background())
	if err != nil {
		//logger.Err("ListUnits error")
		return nil, fmt.Errorf("ListUnits error")
	}

	units, err = matchUnitPatterns(patterns, units)
	if err != nil {
		//logger.Err("matching unit patterns error")
		return nil, fmt.Errorf("error matching unit patterns")
	}

	finalUnits := matchUnitState(states, units)

	return finalUnits, nil
}

func matchUnitPatterns(patterns []string, units []dbus.UnitStatus) ([]dbus.UnitStatus, error) {
	var matchUnits []dbus.UnitStatus
	if len(patterns) == 0 {
		return units, nil
	}
	for _, unit := range units {
		for _, pattern := range patterns {
			match, err := filepath.Match(pattern, unit.Name)
			if err != nil {
				//logger.Err("matching with pattern %s error: ", pattern)
				return nil, err
			}
			if match {
				matchUnits = append(matchUnits, unit)
				break
			}
		}
	}
	return matchUnits, nil
}

func matchUnitState(states []string, units []dbus.UnitStatus) []dbus.UnitStatus {
	if len(states) == 0 {
		return units
	}
	var finalUnits []dbus.UnitStatus
	for _, unit := range units {
		for _, state := range states {
			if unit.LoadState == state || unit.ActiveState == state || unit.SubState == state {
				finalUnits = append(finalUnits, unit)
				break
			}
		}
	}
	return finalUnits

}

func (p *Properties) formProperties(unit dbus.UnitStatus) *Service {
	var service Service

	service.Name = unit.Name
	service.StartType = unit.JobType
	service.ExecPath = string(unit.Path)
	service.DisplayName = unit.Name
	service.Description = unit.Description
	service.State = unit.ActiveState

	if p.ExecMainPID > 0 {
		service.Pid = uint32(p.ExecMainPID)
		service.ExitCode = uint32(p.ExecMainStatus)
	}

	return &service

}
