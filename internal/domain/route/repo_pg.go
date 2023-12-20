package route

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

func (r *repo) Create(route Route) (Route, error) {
	var err error
	var tx pgx.Tx

	if tx, err = r.db.Begin(r.ctx); err != nil {
		return Route{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(r.ctx)
		}
	}()

	q := `INSERT INTO routes 
		(title, points, length, travel_time)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	qp := []any{route.Title, route.Points, route.Length, route.TravelTime}

	logSql := query.NewLogSql(q, qp...)

	if err = tx.QueryRow(r.ctx, q, qp...).Scan(&route.Id); err != nil {
		r.logger.Error(logSql.SetError(err).GetMsg())
		return Route{}, err
	}

	if err = r.insertCoords(route, tx); err != nil {
		return Route{}, err
	}

	tx.Commit(r.ctx)

	r.logger.Debug(logSql.SetResult("ok").GetMsg())

	return r.GetById(route.Id)
}

func (r *repo) Update(route Db) (Route, error) {
	q := `UPDATE routes
		set title = $2,
		    points = $3,
		    length = $4,
			travel_time = $5
		where id = $1
		returning id
	`

	qp := []any{route.Id, route.Title, route.Points, route.Length, route.TravelTime}

	logSql := query.NewLogSql(q, qp...)

	if err := r.db.QueryRow(r.ctx, q, qp...).Scan(&route.Id); err != nil {
		r.logger.Error(logSql.SetError(err).GetMsg())
		return Route{}, err
	}

	r.logger.Debug(logSql.SetResult(route.Id).GetMsg())

	return r.GetById(route.Id)
}

func (r *repo) insertCoords(route Route, tx pgx.Tx) error {
	var coordRows [][]interface{}

	for i := 0; i < len(route.Coords); i++ {
		coord := route.Coords[i]
		coordRows = append(coordRows, []interface{}{route.Id, i, coord[0], coord[1]})
	}

	if copyCount, err := tx.CopyFrom(
		r.ctx,
		pgx.Identifier{"route_points"},
		[]string{"route", "number", "latitude", "longitude"},
		pgx.CopyFromRows(coordRows),
	); err != nil {
		r.logger.Error(err.Error())
		return err
	} else {
		r.logger.Debug("inserted coords: ", copyCount)
	}

	return nil
}

func (r *repo) GetById(id int) (Route, error) {
	q := query.New[Route](r.ctx, r.db).
		Select("r.id", "").
		AddSelect("r.title", "").
		AddSelect("r.points", "").
		AddSelect("r.length", "").
		AddSelect("r.created_at", "").
		AddSelect("r.travel_time", "").
		AddSelect("(select jsonb_agg(t.d) from (select jsonb_build_array(latitude, longitude) as d from route_points where route = r.id order by number) t)", "coords").
		From("routes", "r").
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
		Select("r.*", "").
		From("routes", "r").
		Where(query.EQUEL, "id", id)

	data, err := q.One()

	if errors.Is(err, pgx.ErrNoRows) {
		return data, app_error.ErrNotFound
	}

	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) List(params transport.QueryParams) (query.List[Route], error) {
	q := query.New[Route](r.ctx, r.db).
		Select("r.id", "").
		AddSelect("r.title", "").
		AddSelect("r.created_at", "").
		AddSelect("r.points", "").
		AddSelect("r.length", "").
		AddSelect("r.travel_time", "").
		AddSelect("null", "coords").
		From("routes", "r").
		FilterWhere(params.FindType, "title", params.Find).
		Offset(params.Offset).
		OrderBy("r.title").
		Limit(params.Limit)

	data, err := q.GetList()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) Exists(id int) (bool, error) {
	q := query.New[Route](r.ctx, r.db).
		Select("r.id", "").
		From("routes", "r").
		Where(query.EQUEL, "id", id)

	data, err := q.Exists()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) ExistsByUnique(id int, title string) (bool, error) {
	q := query.New[Route](r.ctx, r.db).
		Select("id", "").
		From("customs", "").
		Where(query.EQUEL, "title", title).
		AndWhere(query.NOT_EQUEL, "id", id)

	data, err := q.Exists()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) DeleteById(id int) (bool, error) {
	q := `DELETE FROM routes where id = $1`

	commandTag, err := r.db.Exec(r.ctx, q, id)
	r.logger.DebugOrError(err, query.NewLogSql(q, id).SetResult(commandTag.RowsAffected() > 0).SetError(err).GetMsg())
	return commandTag.RowsAffected() > 0, err
}
