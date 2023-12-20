package user

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
	ctx    context.Context
	db     pg.DbClient
	logger app_interface.Logger
}

func NewRepo(ctx context.Context, db pg.DbClient, logger app_interface.Logger) Repo {
	return &repo{ctx, db, logger}
}

func (r *repo) Create(user Db) (User, error) {
	q := `INSERT INTO users 
		(login, password, role, title)
		VALUES ($1, $2, $3, $4)
		RETURNING *
	`

	qp := []any{user.Login, user.Password, user.Role, user.Title}

	rows, _ := r.db.Query(r.ctx, q, qp...)

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	r.logger.DebugOrError(err, query.NewLogSql(q, qp...).SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) Update(user Db) (User, error) {
	q := `UPDATE users 
		set login = $2,
		    password = $3,
			role = $4,
			title = $5
		where id = $1
		RETURNING *
	`

	qp := []any{user.Id, user.Login, user.Password, user.Role, user.Title}

	rows, _ := r.db.Query(r.ctx, q, qp...)

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	r.logger.DebugOrError(err, query.NewLogSql(q, qp...).SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) GetByCredentials(login, password string) (User, error) {
	q := query.New[User](r.ctx, r.db).
		Select("*", "").
		From("users", "").
		Where(query.EQUEL, "login", login).
		AndWhere(query.EQUEL, "password", password)

	data, err := q.One()

	if errors.Is(err, pgx.ErrNoRows) {
		return data, app_error.ErrNotFound
	}

	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) GetById(id int) (User, error) {
	q := query.New[User](r.ctx, r.db).
		Select("*", "").
		From("users", "").
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
		From("users", "").
		Where(query.EQUEL, "id", id)

	data, err := q.One()

	if errors.Is(err, pgx.ErrNoRows) {
		return data, app_error.ErrNotFound
	}

	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) List(params transport.QueryParams) (query.List[User], error) {
	q := query.New[User](r.ctx, r.db).
		Select("*", "").
		From("users", "").
		FilterWhere(params.FindType, "login", params.Find).
		OrderBy("login").
		Limit(params.Limit)

	data, err := q.GetList()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err

}

func (r *repo) Exists(id int) (bool, error) {
	q := query.New[User](r.ctx, r.db).
		Select("id", "").
		From("users", "").
		Where(query.EQUEL, "id", id).
		Limit(1)

	data, err := q.Exists()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) ExistsByUnique(id int, title string) (bool, error) {
	q := query.New[User](r.ctx, r.db).
		Select("id", "").
		From("users", "").
		Where(query.EQUEL, "login", title).
		AndWhere(query.NOT_EQUEL, "id", id)

	data, err := q.Exists()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) DeleteById(id int) (bool, error) {
	q := `DELETE FROM users where id = $1`

	commandTag, err := r.db.Exec(r.ctx, q, id)
	r.logger.DebugOrError(err, query.NewLogSql(q, id).SetResult(commandTag.RowsAffected() > 0).SetError(err).GetMsg())
	return commandTag.RowsAffected() > 0, err
}
