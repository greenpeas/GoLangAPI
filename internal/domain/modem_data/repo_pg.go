package modemData

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

func (r *repo) List(params ListParams) ([]ModemData, error) {
	order := "md.reg_time"
	if params.OrderDesc {
		order += " DESC"
	}

	q := query.New[ModemData](r.ctx, r.db).
		Select("md.*", "").
		AddSelect("row_to_json(lbs.*)", "coordinate_lbs").
		AddSelect("coalesce(json_agg(jsonb_build_object('dev_time', sd.dev_time, 'status', sd.status, "+
			"'seal', jsonb_build_object('id', s.id, 'serial', s.serial), 'errors', sd.errors, "+
			"'sensitivity_range', sd.sensitivity_range, 'battery_level', sd.battery_level, 'rssi', sd.rssi, "+
			"'count_commands_in_queue', sd.count_commands_in_queue)) "+
			"FILTER (WHERE sd.dev_time IS NOT NULL), '[]')", "seals_data").
		From("modems_data", "md").
		LeftJoin("lbs", "coordinates_lbs", "lbs.dev_time = md.dev_time and lbs.modem = md.modem").
		LeftJoin("sd", "seals_data", "sd.dev_time between md.dev_time - interval '1 month' and md.dev_time + interval '1 month' "+
			"and sd.modem_time = md.dev_time and sd.modem = md.modem").
		LeftJoin("s", "seals", "s.id = sd.seal").
		Where(query.GREAT, "md.dev_time", params.TimeFrom).
		AndFilterWhere(query.LITTLE_OR_EQ, "md.dev_time", params.TimeTo).
		AndWhere(query.EQUEL, "md.modem", params.ModemId).
		OrderBy(order).
		GroupBy("md.dev_time, md.modem, lbs.dev_time, lbs.modem").
		Limit(params.Limit)

	data, err := q.All()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}
