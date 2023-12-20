package seal

import (
	"context"
	"errors"
	app_interface "seal/internal/app/interface"
	"seal/internal/repository/pg"
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
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

func (r *repo) Update(seal Db) (Seal, error) {
	q := `UPDATE seals
		set comment = $2
		where id = $1
		returning id
	`

	qp := []any{seal.Id, seal.Comment}

	logSql := query.NewLogSql(q, qp...)

	if err := r.db.QueryRow(r.ctx, q, qp...).Scan(&seal.Id); err != nil {
		r.logger.Error(logSql.SetError(err).GetMsg())
		return Seal{}, err
	}

	r.logger.Debug(logSql.SetResult(seal.Id).GetMsg())

	return r.GetById(seal.Id)
}

func (r *repo) GetById(id int) (Seal, error) {
	q := query.New[Seal](r.ctx, r.db).
		Select("s.id", "").
		AddSelect("s.serial", "").
		AddSelect("s.comment", "").
		AddSelect("(select to_jsonb(t) from "+
			"(select * from seals_data where seal = s.id order by dev_time desc limit 1) t)", "last").
		From("seals", "s").
		LeftJoin("l", "seals_data", "l.dev_time = s.last_dev_time and l.seal = s.id").
		Where(query.EQUEL, "s.id", id)

	data, err := q.One()

	if errors.Is(err, pgx.ErrNoRows) {
		return data, app_error.ErrNotFound
	} else if err != nil {
		r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())
		return data, err
	}

	data.Modems, _ = r.getModemsBySerial(data.Serial)

	return data, nil
}

func (r *repo) getModemsBySerial(serial uint64) ([]Modem, error) {
	if serial == 0 {
		return []Modem{}, nil
	}

	q := query.New[Modem](r.ctx, r.db).
		Select("m.id", "").
		AddSelect("m.serial", "").
		AddSelect("row_to_json(l.*)", "last").
		From("modems", "m").
		LeftJoin("l", "modems_data", "l.dev_time = m.last_dev_time and l.modem = m.id").
		Where(query.ANY, "m.serials_of_seals", serial)

	data, err := q.All()

	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) GetDbById(id int) (Db, error) {
	q := query.New[Db](r.ctx, r.db).
		Select("r.*", "").
		From("seals", "r").
		Where(query.EQUEL, "id", id)

	data, err := q.One()

	if errors.Is(err, pgx.ErrNoRows) {
		return data, app_error.ErrNotFound
	}

	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) List(params transport.QueryParams) (query.List[SealForList], error) {
	q := query.New[SealForList](r.ctx, r.db).
		Select("s.id", "").
		AddSelect("s.serial", "").
		AddSelect("(select to_jsonb(t) from "+
			"(select * from seals_data where seal = s.id order by dev_time desc limit 1) t)", "last").
		AddSelect("(select coalesce(jsonb_agg(to_jsonb(t)), '[]') from (select m.id, m.serial "+
			"from modems m where s.serial = any(m.serials_of_seals) "+
			"order by m.serial) t)", "modems").
		From("seals", "s").
		FilterWhere(params.FindType, "serial", params.Find).
		OrderBy("s.serial").
		Limit(params.Limit).
		Offset(params.Offset)

	data, err := q.GetList()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *repo) Exists(id int) (bool, error) {
	q := query.New[Seal](r.ctx, r.db).
		Select("id", "").
		From("seals", "").
		Where(query.EQUEL, "id", id).
		Limit(1)

	data, err := q.Exists()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}
