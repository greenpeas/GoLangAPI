package custom

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

func (r *repo) Create(custom Db) (Custom, error) {
	q := `INSERT INTO customs 
		(title)
		VALUES ($1)
		RETURNING *
	`

	qp := []any{custom.Title}

	rows, _ := r.db.Query(r.ctx, q, qp...)

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Custom])
	r.logger.DebugOrError(err, query.NewLogSql(q, qp...).SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) Update(custom Db) (Custom, error) {
	q := `UPDATE customs 
			set title = $2
		where id = $1
		RETURNING *
	`

	qp := []any{custom.Id, custom.Title}

	rows, _ := r.db.Query(r.ctx, q, qp...)

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Custom])
	r.logger.DebugOrError(err, query.NewLogSql(q, qp...).SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) GetById(id int) (Custom, error) {
	q := query.New[Custom](r.ctx, r.db).
		Select("*", "").
		From("customs", "").
		Where(query.EQUEL, "id", id)

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
		From("customs", "").
		Where(query.EQUEL, "id", id)

	data, err := q.One()

	if errors.Is(err, pgx.ErrNoRows) {
		return data, app_error.ErrNotFound
	}

	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) List(params transport.QueryParams) (query.List[Custom], error) {
	q := query.New[Custom](r.ctx, r.db).
		Select("*", "").
		From("customs", "").
		FilterWhere(params.FindType, "title", params.Find).
		OrderBy("title").
		Limit(params.Limit)

	data, err := q.GetList()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) Exists(id int) (bool, error) {
	q := query.New[Custom](r.ctx, r.db).
		Select("id", "").
		From("customs", "").
		Where(query.EQUEL, "id", id)

	data, err := q.Exists()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) ExistsByUnique(id int, title string) (bool, error) {
	q := query.New[Custom](r.ctx, r.db).
		Select("id", "").
		From("customs", "").
		Where(query.EQUEL, "title", title).
		AndWhere(query.NOT_EQUEL, "id", id)

	data, err := q.Exists()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) DeleteById(id int) (bool, error) {
	q := `DELETE FROM customs where id = $1`

	commandTag, err := r.db.Exec(r.ctx, q, id)
	r.logger.DebugOrError(err, query.NewLogSql(q, id).SetResult(commandTag.RowsAffected() > 0).SetError(err).GetMsg())
	return commandTag.RowsAffected() > 0, err
}
