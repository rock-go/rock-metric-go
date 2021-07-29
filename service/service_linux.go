package service

import (
	"context"
	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/mitchellh/mapstructure"
	"github.com/rock-go/rock/logger"
	"strings"
)

var (
	Conn     *dbus.Conn
	UnitList []unitFetcher
)

func init() {
	var err error
	Conn, err = dbus.NewWithContext(context.Background())
	if err != nil {
		logger.Errorf("connect to dbus error: %v", err)
	}

	UnitList = []unitFetcher{listUnitsByPatternWrapper, listUnitsFilteredWrapper, listUnitsWrapper}
}

func GetService(pattern string) []*Service {

	var units []dbus.UnitStatus

	if Conn == nil {
		logger.Errorf("no conn to dbus")
		return nil
	}

	units = getUnits(Conn, []string{}, []string{})
	if units == nil {
		return nil
	}

	var servicesInfo []*Service
	for _, unit := range units {
		if unit.LoadState == "not-found" {
			continue
		}

		if !strings.Contains(strings.ToLower(unit.Name), pattern) {
			continue
		}

		props, err := getProps(Conn, unit.Name)
		if err != nil {
			logger.Errorf("get properties from service [%s] error: %v", unit.Name, err)
			continue
		}

		metrics := props.formProperties(unit)

		servicesInfo = append(servicesInfo, metrics)
	}

	return servicesInfo
}

func getUnits(conn *dbus.Conn, states, patterns []string) []dbus.UnitStatus {

	for _, unit := range UnitList {
		units, err := unit(conn, states, patterns)
		if err == nil {
			logger.Debugf("get dbus unit success by %v", unit)
			return units
		} else {
			logger.Errorf("get dbus unit by %v error: %v", unit, err)
		}
	}

	logger.Errorf("get dbus unit error by all methods")
	return nil
}

func getProps(conn *dbus.Conn, unit string) (Properties, error) {
	parsed := Properties{}
	rawProps, err := conn.GetAllPropertiesContext(context.Background(), unit)
	if err != nil {
		return parsed, err
	}

	err = mapstructure.Decode(rawProps, &parsed)
	return parsed, err
}
