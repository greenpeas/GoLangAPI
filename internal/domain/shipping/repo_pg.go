package shipping

import (
	"context"
	"errors"
	app_interface "seal/internal/app/interface"
	"seal/internal/repository/pg"
	"seal/internal/repository/pg/query"
	"seal/pkg/app_error"

	"github.com/jackc/pgx/v5"
)

type repo struct {
	db     pg.DbClient
	logger app_interface.Logger
	ctx    context.Context
}

func NewRepo(ctx context.Context, db pg.DbClient, logger app_interface.Logger) Repo {
	return &repo{db, logger, ctx}
}

func (r *repo) Create(sh Db) (Shipping, error) {
	q := `INSERT INTO shipping
		(author, custom_number, create_date, number, transport, route)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	qp := []any{sh.Author, sh.CustomNumber, sh.CreateDate, sh.Number, sh.Transport, sh.Route}

	err := r.db.QueryRow(r.ctx, q, qp...).Scan(&sh.Id)

	r.logger.DebugOrError(err, query.NewLogSql(q, qp...).SetResult(sh.Id).SetError(err).GetMsg())

	if err != nil {
		return Shipping{}, err
	}

	return r.GetById(sh.Id)
}

func (r *repo) Update(sh Db) (Shipping, error) {
	q := `UPDATE shipping 
		set (author, custom_number, create_date, number, transport, route, status, time_start, time_end, files, modem) = 
		    ($2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		    where id = $1
		RETURNING id
	`

	qp := []any{sh.Id, sh.Author, sh.CustomNumber, sh.CreateDate, sh.Number, sh.Transport, sh.Route, sh.Status,
		sh.TimeStart, sh.TimeEnd, sh.Files, sh.Modem}

	err := r.db.QueryRow(r.ctx, q, qp...).Scan(&sh.Id)

	r.logger.DebugOrError(err, query.NewLogSql(q, qp...).SetResult(sh.Id).SetError(err).GetMsg())

	if err != nil {
		return Shipping{}, err
	}

	return r.GetById(sh.Id)
}

func (r *repo) GetById(id int) (Shipping, error) {
	q := query.New[Shipping](r.ctx, r.db).
		Select("s.id", "").
		AddSelect("s.created_at", "").
		AddSelect("s.custom_number", "").
		AddSelect("s.create_date", "").
		AddSelect("s.number", "").
		AddSelect("s.status", "").
		AddSelect("s.time_start", "").
		AddSelect("s.time_end", "").
		AddSelect("s.files", "").
		AddSelect("coalesce(s.time_start, now()) + (r.travel_time * interval '1 minute')", "estimated_arrival_time").
		AddSelect("jsonb_build_object('id', u.id, 'login', u.login)", "author").
		AddSelect("(to_jsonb(t.*) || jsonb_build_object('type', to_jsonb(tt.*)))", "transport").
		AddSelect("to_jsonb(r.*)", "route").
		AddSelect("to_jsonb(m.*) || jsonb_build_object('last', to_jsonb(ml.*))", "modem").
		AddSelect("(with r as (select distinct seal as seal_id from seals_data sd "+
			//"where sd.dev_time >= coalesce(s.time_start, s.created_at) and sd.dev_time < coalesce(s.time_end, now()) "+
			"where sd.dev_time >= s.created_at and sd.dev_time < coalesce(s.time_end, now()) "+
			"and modem = s.modem) "+
			"select coalesce(jsonb_agg(jsonb_build_object('id', seals.id, 'serial', seals.serial, 'last', to_jsonb(sd.*)) "+
			"order by seals.serial), '[]') "+
			"from r inner join seals on seals.id = r.seal_id "+
			"left join lateral (select * from seals_data where seal = seals.id and modem = m.id order by dev_time desc limit 1) sd ON true)", "seals").
		From("shipping", "s").
		LeftJoin("u", "users", "u.id=s.author").
		LeftJoin("t", "transports", "t.id=s.transport").
		LeftJoin("tt", "transport_types", `tt.id=t.type`).
		LeftJoin("r", "routes", "r.id=s.route").
		LeftJoin("m", "modems", "m.id=s.modem").
		LeftJoin("ml", "modems_data", "ml.dev_time = m.last_dev_time and ml.modem = m.id").
		Where(query.EQUEL, "s.id", id)

	data, err := q.One()

	if errors.Is(err, pgx.ErrNoRows) {
		return data, app_error.ErrNotFound
	}

	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) GetActiveByModemId(id int) (Shipping, error) {
	q := query.New[struct {
		Id int `json:"id"`
	}](r.ctx, r.db).
		Select("s.id", "").
		From("shipping", "s").
		Where(query.EQUEL, "s.modem", id).
		AndWhere(query.EQUEL, "s.status", 0)

	data, err := q.One()

	if errors.Is(err, pgx.ErrNoRows) {
		return Shipping{}, app_error.ErrNotFound
	} else if err != nil {
		return Shipping{}, err
	}

	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return r.GetById(data.Id)
}

func (r *repo) GetDbById(id int) (Db, error) {
	q := query.New[Db](r.ctx, r.db).
		Select("*", "").
		From("shipping", "s").
		Where(query.EQUEL, "s.id", id)

	data, err := q.One()

	if errors.Is(err, pgx.ErrNoRows) {
		return data, app_error.ErrNotFound
	}

	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) List(params QueryParams) (query.List[ShippingForList], error) {
	q := query.New[ShippingForList](r.ctx, r.db).
		Select("s.id", "").
		//AddSelect("s.created_at", "").
		AddSelect("s.custom_number", "").
		AddSelect("s.create_date", "").
		AddSelect("s.number", "").
		AddSelect("s.status", "").
		//AddSelect("s.time_start", "").
		//AddSelect("s.time_end", "").
		AddSelect("coalesce(s.time_start, now()) + (r.travel_time * interval '1 minute')", "estimated_arrival_time").
		AddSelect("s.files", "").
		//AddSelect("jsonb_build_object('id', u.id, 'login', u.login)", "author").
		AddSelect("(to_jsonb(t.*) || jsonb_build_object('type', to_jsonb(tt.*)))", "transport").
		AddSelect("to_jsonb(r.*)", "route").
		From("shipping", "s").
		LeftJoin("u", "users", "u.id=s.author").
		LeftJoin("t", "transports", "t.id=s.transport").
		LeftJoin("tt", "transport_types", `tt.id=t.type`).
		LeftJoin("r", "routes", "r.id=s.route").
		FilterWhere(params.FindType, "custom_number", params.Find).
		OrFilterWhere(params.FindType, "create_date", params.Find).
		OrFilterWhere(params.FindType, "number", params.Find).
		OrFilterWhere(query.IN, "status", params.Status).
		OrderBy("s.status, s.custom_number, s.create_date, s.number").
		Limit(params.Limit).
		Offset(params.Offset)

	data, err := q.GetList()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) Exists(id int) (bool, error) {
	q := query.New[Shipping](r.ctx, r.db).
		Select("id", "").
		From("shipping", "").
		Where(query.EQUEL, "id", id)

	data, err := q.Exists()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) ExistsByUnique(customNumber, createDate string, id, number int) (bool, error) {
	q := query.New[Shipping](r.ctx, r.db).
		Select("id", "").
		From("shipping", "").
		Where(query.EQUEL, "custom_number", customNumber).
		AndWhere(query.EQUEL, "create_date", createDate).
		AndWhere(query.EQUEL, "number", number).
		AndWhere(query.NOT_EQUEL, "id", id)

	data, err := q.Exists()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) DeleteById(id int) (bool, error) {
	q := `DELETE FROM shipping where id = $1`

	commandTag, err := r.db.Exec(r.ctx, q, id)
	r.logger.DebugOrError(err, query.NewLogSql(q, id).SetResult(commandTag.RowsAffected() > 0).SetError(err).GetMsg())
	return commandTag.RowsAffected() > 0, err
}
