package modemLogRaw

import (
	"context"
	app_interface "seal/internal/app/interface"
	"seal/internal/repository/pg"
	"seal/internal/repository/pg/query"
)

const packetFromDevice = 0
const telemetryCommandName = "get-telemetry"

type repo struct {
	ctx    context.Context
	db     pg.DbClient
	logger app_interface.Logger
}

func NewRepo(ctx context.Context, db pg.DbClient, logger app_interface.Logger) Repo {
	return &repo{ctx, db, logger}
}

func (r *repo) List(params ListParams) ([]ModemLogRaw, error) {
	order := "reg_time"
	if params.OrderDesc {
		order += " DESC"
	}

	q := query.New[ModemLogRaw](r.ctx, r.db).
		Select("*", "").
		From("modems_log_raw", "d").
		Where(query.GREAT, "reg_time", params.From).
		AndFilterWhere(query.LITTLE_OR_EQ, "reg_time", params.To).
		AndWhere(query.EQUEL, "d.imei", params.Imei).
		OrderBy(order).
		Limit(params.Limit)

	data, err := q.All()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) ListTelemetry(params ListParams) ([]Telemetry, error) {
	order := "reg_time"
	if params.OrderDesc {
		order += " DESC"
	}

	q := query.New[Telemetry](r.ctx, r.db).
		Select("d.reg_time", "").
		AddSelect("coalesce((d.payload#>>'{current_time}')::int8, 0)", "current_time").
		AddSelect("coalesce((d.payload#>>'{status}')::int8, 0)", "status").
		AddSelect("coalesce((d.payload#>>'{errors_flags}')::int8, 0)", "errors_flags").
		AddSelect("coalesce((d.payload#>>'{positioning_time}')::int8, 0)", "positioning_time").
		AddSelect("d.payload#>>'{latitude}'", "latitude").
		AddSelect("d.payload#>>'{longitude}'", "longitude").
		AddSelect("coalesce((d.payload#>>'{altitude}')::int8, 0)", "altitude").
		AddSelect("coalesce((d.payload#>>'{satellites_count}')::int8, 0)", "satellites_count").
		AddSelect("coalesce((d.payload#>>'{speed}')::int8, 0)", "speed").
		AddSelect("coalesce((d.payload#>>'{status_gps_module}')::int8, 0)", "status_gps_module").
		AddSelect("coalesce((d.payload#>>'{rssi}')::int8, 0)", "rssi").
		AddSelect("coalesce((d.payload#>>'{battery_level}')::int8, 0)", "battery_level").
		AddSelect("coalesce((d.payload#>>'{signal_gps}')::int8, 0)", "signal_gps").
		AddSelect("coalesce((d.payload#>>'{signal_glonass}')::int8, 0)", "signal_glonass").
		From("modems_log_raw", "d").
		Where(query.GREAT, "reg_time", params.From).
		AndFilterWhere(query.LITTLE_OR_EQ, "reg_time", params.To).
		AndWhere(query.EQUEL, "d.imei", params.Imei).
		AndWhere(query.EQUEL, "src", packetFromDevice).
		AndWhere(query.EQUEL, "cmd_name", telemetryCommandName).
		OrderBy(order).
		Limit(params.Limit)

	data, err := q.All()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}
