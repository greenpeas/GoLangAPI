package transport

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

func (r *repo) Create(transport Db) (Transport, error) {
	q := `INSERT INTO transports
		(author, title, type, registration_number)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	qp := []any{transport.Author, transport.Title, transport.Type, transport.RegistrationNumber}

	err := r.db.QueryRow(r.ctx, q, qp...).Scan(&transport.Id)
	r.logger.DebugOrError(err, query.NewLogSql(q, qp...).SetResult(transport.Id).SetError(err).GetMsg())

	if err != nil {
		return Transport{}, err
	}

	return r.GetById(transport.Id)
}

func (r *repo) Update(transport Db) (Transport, error) {
	q := `UPDATE transports 
		set (author, title, type, registration_number) = ($2, $3, $4, $5)
		where id = $1
		RETURNING id
	`

	qp := []any{transport.Id, transport.Author, transport.Title, transport.Type, transport.RegistrationNumber}

	err := r.db.QueryRow(r.ctx, q, qp...).Scan(&transport.Id)
	r.logger.DebugOrError(err, query.NewLogSql(q, qp...).SetResult(transport.Id).SetError(err).GetMsg())

	if err != nil {
		return Transport{}, err
	}

	return r.GetById(transport.Id)
}

func (r *repo) GetById(id int) (Transport, error) {
	q := query.New[Transport](r.ctx, r.db).
		Select("t.id", "").
		AddSelect("t.created_at", "").
		AddSelect("t.title", "").
		AddSelect("t.registration_number", "").
		AddSelect("jsonb_build_object('id', u.id, 'login', u.login)", "author").
		AddSelect("jsonb_build_object('id', tt.id, 'title', t.title)", "type").
		From("transports", "t").
		LeftJoin("u", "users", "u.id=t.author").
		LeftJoin("tt", "transport_types", `tt.id=t.type`).
		Where(query.EQUEL, "t.id", id)

	data, err := q.One()

	if errors.Is(err, pgx.ErrNoRows) {
		return data, app_error.ErrNotFound
	}

	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) GetDbById(id int) (Db, error) {
	q := query.New[Db](r.ctx, r.db).
		Select("*", "").
		From("transports", "s").
		Where(query.EQUEL, "s.id", id)

	data, err := q.One()

	if errors.Is(err, pgx.ErrNoRows) {
		return data, app_error.ErrNotFound
	}

	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) List(params transport.QueryParams) (query.List[Transport], error) {
	q := query.New[Transport](r.ctx, r.db).
		Select("t.id", "").
		AddSelect("t.created_at", "").
		AddSelect("t.title", "").
		AddSelect("t.registration_number", "").
		AddSelect("jsonb_build_object('id', u.id, 'login', u.login)", "author").
		AddSelect("jsonb_build_object('id', tt.id, 'title', t.title)", "type").
		From("transports", "t").
		LeftJoin("u", "users", "u.id=t.author").
		LeftJoin("tt", "transport_types", `tt.id=t.type`).
		FilterWhere(params.FindType, "t.title", params.Find).
		OrderBy("t.title").
		Limit(params.Limit)

	data, err := q.GetList()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err

}

func (r *repo) Exists(id int) (bool, error) {
	q := query.New[Transport](r.ctx, r.db).
		Select("id", "").
		From("transports", "").
		Where(query.EQUEL, "id", id).
		Limit(1)

	data, err := q.Exists()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) ExistsByUnique(id int, title string) (bool, error) {
	q := query.New[Transport](r.ctx, r.db).
		Select("id", "").
		From("transports", "").
		Where(query.EQUEL, "title", title).
		AndWhere(query.NOT_EQUEL, "id", id)

	data, err := q.Exists()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) DeleteById(id int) (bool, error) {
	q := `DELETE FROM transports where id = $1`

	commandTag, err := r.db.Exec(r.ctx, q, id)
	r.logger.DebugOrError(err, query.NewLogSql(q, id).SetResult(commandTag.RowsAffected() > 0).SetError(err).GetMsg())
	return commandTag.RowsAffected() > 0, err
}
