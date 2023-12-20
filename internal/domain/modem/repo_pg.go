package modem

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	app_interface "seal/internal/app/interface"
	"seal/internal/repository/pg"
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
	"seal/pkg/app_error"
)

const MAX_RETURNING_ROWS = 50000

type repo struct {
	db     pg.DbClient
	logger app_interface.Logger
	ctx    context.Context
}

func NewRepo(ctx context.Context, db pg.DbClient, logger app_interface.Logger) Repo {
	return &repo{db, logger, ctx}
}

func (r *repo) Update(modem Db) (Modem, error) {
	q := `UPDATE modems
		set comment = $2
		where id = $1
		returning id
	`

	qp := []any{modem.Id, modem.Comment}

	logSql := query.NewLogSql(q, qp...)

	if err := r.db.QueryRow(r.ctx, q, qp...).Scan(&modem.Id); err != nil {
		r.logger.Error(logSql.SetError(err).GetMsg())
		return Modem{}, err
	}

	r.logger.Debug(logSql.SetResult(modem.Id).GetMsg())

	return r.GetById(modem.Id)
}

func (r *repo) GetById(id int) (Modem, error) {

	qLastCharge := "(select dev_time from modems_data where dev_time > " +
		"(select dev_time from modems_data md where dev_time < now() and dev_time < " +
		"(select dev_time from modems_data where dev_time < now() and modem = m.id and battery_level = 100 order by dev_time desc limit 1) " +
		"and modem = m.id and battery_level < 100 order by dev_time desc limit 1) and modem = m.id order by dev_time limit 1)"

	q := query.New[Modem](r.ctx, r.db).
		Select("m.id", "").
		AddSelect("m.imei", "").
		AddSelect("m.serial", "").
		AddSelect("m.iccid", "").
		AddSelect("m.last_dev_time", "").
		AddSelect("m.extra", "").
		AddSelect("m.serials_of_seals", "").
		AddSelect("m.comment", "").
		AddSelect("m.msisdn", "").
		AddSelect("(to_jsonb(l.*)) || jsonb_build_object('coordinate_lbs', to_jsonb(lbs.*))", "last").
		AddSelect("(select to_jsonb(t) from (select "+
			"c.dev_time, c.latitude, c.longitude, c.altitude, c.satellites_count, speed, status_gps_module, min_distance_to_route "+
			"from coordinates c where c.dev_time < Now() and c.modem = m.id and latitude != 'NaN' and longitude != 'NaN' "+
			"order by c.dev_time desc limit 1) t)", "last_coordinate").
		AddSelect("(select jsonb_agg(t ORDER BY t.serial) from (select s.*,  to_jsonb(sd.*) as last "+
			"from seals s "+
			"left join lateral (select * from seals_data where seal = s.id and modem = m.id order by dev_time desc limit 1) sd ON true "+
			"where s.serial = any(m.serials_of_seals)) t)", "seals").
		AddSelect(qLastCharge, "last_charge_time").
		From("modems", "m").
		LeftJoin("l", "modems_data", "l.dev_time = m.last_dev_time and l.modem = m.id").
		LeftJoin("lbs", "coordinates_lbs", "lbs.dev_time = l.dev_time and lbs.modem = l.modem").
		Where(query.EQUEL, "m.id", id)

	data, err := q.One()

	if errors.Is(err, pgx.ErrNoRows) {
		return data, app_error.ErrNotFound
	} else if err != nil {
		r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())
		return data, err
	}

	//data.Seals, _ = r.getSealsBySerials(data.SerialsOfSeals)

	return data, nil
}

func (r *repo) GetDbById(id int) (Db, error) {
	q := query.New[Db](r.ctx, r.db).
		Select("r.*", "").
		From("modems", "r").
		Where(query.EQUEL, "id", id)

	data, err := q.One()

	if errors.Is(err, pgx.ErrNoRows) {
		return data, app_error.ErrNotFound
	}

	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) GetByImei(imei uint64) (Modem, error) {
	q := query.New[struct {
		Id int `json:"id"`
	}](r.ctx, r.db).
		Select("m.id", "").
		From("modems", "m").
		Where(query.EQUEL, "m.imei", imei)

	data, err := q.One()

	if errors.Is(err, pgx.ErrNoRows) {
		return Modem{}, app_error.ErrNotFound
	} else if err != nil {
		return Modem{}, err
	}

	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return r.GetById(data.Id)
}

func (r *repo) List(params transport.QueryParams) (query.List[ModemForList], error) {
	q := query.New[ModemForList](r.ctx, r.db).
		Select("m.id", "").
		AddSelect("m.imei", "").
		AddSelect("m.serial", "").
		AddSelect("m.iccid", "").
		AddSelect("m.last_dev_time", "").
		AddSelect("m.comment", "").
		AddSelect("m.extra->>'software_version'", "software_version").
		AddSelect("m.extra->>'hardware_revision'", "hardware_revision").
		AddSelect("(select CASE WHEN l.reg_time is null THEN null else to_jsonb(t) end "+
			"from (select l.status, l.reg_time, l.errors_flags, l.rssi, l.connect_period, l.battery_level) t)", "last").
		AddSelect("(select to_jsonb(t) from (select "+
			"c.dev_time, c.latitude, c.longitude, c.altitude, c.satellites_count, speed, status_gps_module, min_distance_to_route "+
			"from coordinates c where c.dev_time < Now() and c.modem = m.id and latitude != 'NaN' and longitude != 'NaN' "+
			"order by c.dev_time desc limit 1) t)", "last_coordinate").
		From("modems", "m").
		LeftJoin("l", "modems_data", "l.dev_time = m.last_dev_time and l.modem = m.id").
		FilterWhere(params.FindType, "serial", params.Find).
		OrFilterWhere(params.FindType, "imei", params.Find).
		OrderBy("serial").
		Limit(min(params.Limit, MAX_RETURNING_ROWS)).
		Offset(params.Offset)

	data, err := q.GetList()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *repo) ListShippingReady(params transport.QueryParams) (query.List[ModemForListShippingReady], error) {
	q := query.New[ModemForListShippingReady](r.ctx, r.db).
		Select("m.id", "").
		AddSelect("m.serial", "").
		AddSelect("jsonb_agg(to_jsonb(s.*) || jsonb_build_object('last', sl.*))", "seals").
		From("modems", "m").
		LeftJoin("l", "modems_data", "l.dev_time = m.last_dev_time and l.modem = m.id").
		InnerJoin("s", "seals", "s.serial = any(m.serials_of_seals)").
		InnerJoin("sl", "seals_data", "sl.dev_time = s.last_dev_time and sl.seal = s.id").
		FilterWhere(params.FindType, "m.serial", params.Find).
		OrFilterWhere(params.FindType, "m.imei", params.Find).
		AndWhereNotExists("select * from shipping where modem=m.id and status<>2").
		OrderBy("m.serial").
		GroupBy("m.id").
		Limit(min(params.Limit, MAX_RETURNING_ROWS)).
		Offset(params.Offset)

	data, err := q.GetList()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *repo) Track(params TrackQueryParams) ([]Coordinate, error) {
	order := "c.dev_time"
	if params.OrderDesc {
		order += " DESC"
	}

	q := query.New[Coordinate](r.ctx, r.db).
		Select("c.latitude", "").
		AddSelect("c.longitude", "").
		AddSelect("c.dev_time", "").
		AddSelect("c.altitude", "").
		AddSelect("c.satellites_count", "").
		AddSelect("c.speed", "").
		AddSelect("c.status_gps_module", "").
		AddSelect("c.hdop", "").
		AddSelect("c.signal_gps", "").
		AddSelect("c.signal_glonass", "").
		AddSelect("c.min_distance_to_route", "").
		From("coordinates", "c").
		Where(query.GREAT, "c.dev_time", params.From).
		AndFilterWhere(query.LITTLE_OR_EQ, "c.dev_time", params.To).
		AndWhere(query.EQUEL, "c.modem", params.Id).
		OrderBy(order).
		Limit(min(params.Limit, MAX_RETURNING_ROWS))

	data, err := q.All()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, nil
}

func (r *repo) TrackLbs(params TrackQueryParams) ([]CoordinateLbs, error) {
	order := "c.dev_time"
	if params.OrderDesc {
		order += " DESC"
	}

	q := query.New[CoordinateLbs](r.ctx, r.db).
		Select("*", "").
		From("coordinates_lbs", "c").
		Where(query.GREAT, "c.dev_time", params.From).
		AndFilterWhere(query.LITTLE_OR_EQ, "c.dev_time", params.To).
		AndWhere(query.EQUEL, "c.modem", params.Id).
		OrderBy(order).
		Limit(min(params.Limit, MAX_RETURNING_ROWS))

	data, err := q.All()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, nil
}
