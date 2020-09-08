package datadog

import "github.com/zorkian/go-datadog-api"

//go:generate $GOPATH/bin/mockgen -destination=./mocks/datadog.go -package=mocks github.com/jonnylangefeld/datadog-operator/pkg/datadog ClientInterface
type ClientInterface interface {
	UpdateMonitor(monitor *datadog.Monitor) error
	CreateMonitor(monitor *datadog.Monitor) (*datadog.Monitor, error)
	DeleteMonitor(id int) error
}
