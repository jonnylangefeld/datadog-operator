package controllers

import (
	"errors"

	"github.com/golang/mock/gomock"
	dd "github.com/zorkian/go-datadog-api"

	"github.com/jonnylangefeld/datadog-operator/pkg/datadog"
	"github.com/jonnylangefeld/datadog-operator/pkg/datadog/mocks"
)

func getDatadogClientMock() datadog.ClientInterface {
	datadogClient := mocks.NewMockClientInterface(mockController)

	datadogClient.
		EXPECT().
		UpdateMonitor(gomock.Any()).
		Return(errors.New("Not Found"))

	m, _ := testMonitor.Transform()
	m.Id = dd.Int(1)
	create := datadogClient.
		EXPECT().
		CreateMonitor(gomock.Any()).
		Return(m, nil)

	datadogClient.
		EXPECT().
		UpdateMonitor(gomock.Any()).
		Return(nil).
		After(create)

	datadogClient.
		EXPECT().
		DeleteMonitor(gomock.Any()).
		Return(nil)

	return datadogClient
}
