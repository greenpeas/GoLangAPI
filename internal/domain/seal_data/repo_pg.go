package sealData

import (
	"context"
	app_interface "seal/internal/app/interface"
	"seal/internal/repository/pg"
	"seal/internal/repository/pg/query"
)

type repo struct {
	ctx    context.Context
	db     pg.DbClient
	logger app_interface.Logger
}

func NewRepo(ctx context.Context, db pg.DbClient, logger app_interface.Logger) Repo {
	return &repo{ctx, db, logger}
}

func (r *repo) List(params ListParams) ([]SealData, error) {
	order := "sd.dev_time"
	if params.OrderDesc {
		order += " DESC"
	}

	q := query.New[SealData](r.ctx, r.db).
		Select("sd.dev_time", "").
		AddSelect("sd.status_raw", "").
		AddSelect("sd.status", "").
		AddSelect("sd.errors_raw", "").
		AddSelect("sd.errors", "").
		AddSelect("sd.sensitivity_range", "").
		AddSelect("sd.battery_level", "").
		AddSelect("sd.rssi", "").
		AddSelect("sd.temperature", "").
		AddSelect("sd.sensitivity_cable", "").
		AddSelect("sd.build_version", "").
		AddSelect("sd.count_commands_in_queue", "").
		AddSelect("jsonb_build_object('id', m.id, 'serial', m.serial)", "modem").
		//AddSelect("row_to_json(md.*)", "modem_data").
		From("seals_data", "sd").
		LeftJoin("m", "modems", "m.id = sd.modem").
		//LeftJoin("md", "modems_data", "md.dev_time = sd.modem_time and md.modem = sd.modem").
		Where(query.GREAT, "sd.dev_time", params.TimeFrom).
		AndFilterWhere(query.LITTLE_OR_EQ, "sd.dev_time", params.TimeTo).
		AndWhere(query.EQUEL, "seal", params.SealId).
		OrderBy(order).
		Limit(params.Limit)

	data, err := q.All()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}
