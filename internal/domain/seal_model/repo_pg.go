package seal_model

import (
	"context"
	app_interface "seal/internal/app/interface"
	"seal/internal/repository/pg"
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
)

type repo struct {
	db     pg.DbClient
	logger app_interface.Logger
	ctx    context.Context
}

func NewRepo(ctx context.Context, db pg.DbClient, logger app_interface.Logger) Repo {
	return &repo{db, logger, ctx}
}

func (r *repo) List(params transport.QueryParams) (query.List[SealModel], error) {
	q := query.New[SealModel](r.ctx, r.db).
		Select("*", "").
		From("seal_model", "").
		FilterWhere(params.FindType, "title", params.Find).
		OrderBy("title").
		Limit(params.Limit)

	data, err := q.GetList()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}
