// SPDX-License-Identifier: Apache-2.0

package ocpp16

import (
	"context"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"golang.org/x/exp/slog"
)

func DiagnosticsStatusNotificationHandler(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*types.DiagnosticsStatusNotificationJson)
	slog.Info("diagnostic status notification", slog.String("chargeStationId", chargeStationId),
		slog.Any("status", req.Status))
	return &types.DiagnosticsStatusNotificationResponseJson{}, nil
}
