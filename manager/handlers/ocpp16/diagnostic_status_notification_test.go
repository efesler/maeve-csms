// SPDX-License-Identifier: Apache-2.0

package ocpp16_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	handlers "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp16"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
)

func TestDiagnosticsStatusNotificationHandler(t *testing.T) {

	req := &types.DiagnosticsStatusNotificationJson{
		Status: types.DiagnosticsStatusNotificationJsonStatusIdle,
	}

	got, err := handlers.DiagnosticsStatusNotificationHandler(context.Background(), "cs001", req)
	assert.NoError(t, err)

	want := &types.DiagnosticsStatusNotificationResponseJson{}

	assert.Equal(t, want, got)
}
