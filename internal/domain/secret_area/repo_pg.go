package secret_area

import (
	"context"
	"errors"
	"fmt"
	app_interface "seal/internal/app/interface"
	"seal/internal/domain/user"
	"seal/internal/repository/pg"
	"seal/internal/repository/pg/query"
	"seal/pkg/app_error"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type repo struct {
	db     pg.DbClient
	logger app_interface.Logger
	ctx    context.Context
}

func NewRepo(ctx context.Context, db pg.DbClient, logger app_interface.Logger) Repo {
	return &repo{db, logger, ctx}
}

const STATUS_PRESENT = 2

func getAreaForQuery(area [][2]float32) string {
	var areas []string
	for _, v := range area {
		areas = append(areas, fmt.Sprintf("(%f,%f)", v[0], v[1]))
	}

	return fmt.Sprintf("(%s)", strings.Join(areas, ", "))
}

func (r *repo) Create(sa Db) (SecretArea, error) {
	var id int

	q := `INSERT INTO secret_areas
		(author, title, description, area)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	qp := []any{sa.Author, sa.Title, sa.Description, getAreaForQuery(sa.Area)}

	err := r.db.QueryRow(r.ctx, q, qp...).Scan(&id)

	r.logger.DebugOrError(err, query.NewLogSql(q, qp...).SetResult(id).SetError(err).GetMsg())

	if err != nil {
		return SecretArea{}, err
	}

	return r.GetById(id)
}

func (r *repo) Update(sa Db) (SecretArea, error) {
	q := `UPDATE secret_areas 
		set (title, description, area) = ($2, $3, $4)
		    where id = $1
		RETURNING id
	`

	qp := []any{sa.Id, sa.Title, sa.Description, getAreaForQuery(sa.Area)}

	err := r.db.QueryRow(r.ctx, q, qp...).Scan(&sa.Id)

	r.logger.DebugOrError(err, query.NewLogSql(q, qp...).SetResult(sa.Id).SetError(err).GetMsg())

	if err != nil {
		return SecretArea{}, err
	}

	return r.GetById(sa.Id)
}

type getByIdDb struct {
	Id          int
	CreatedAt   *time.Time `db:"created_at"`
	Author      *user.Author
	Title       string
	Area        pgtype.Polygon
	Description string
}

func (r *repo) GetById(id int) (SecretArea, error) {
	q := query.New[getByIdDb](r.ctx, r.db).
		Select("sa.id", "").
		AddSelect("sa.created_at", "").
		AddSelect("sa.title", "").
		AddSelect("sa.area", "").
		AddSelect("sa.description", "").
		AddSelect("jsonb_build_object('id', u.id, 'login', u.login)", "author").
		From("secret_areas", "sa").
		LeftJoin("u", "users", "u.id=sa.author").
		Where(query.EQUEL, "sa.id", id)

	data, err := q.One()

	if errors.Is(err, pgx.ErrNoRows) {
		return SecretArea{}, app_error.ErrNotFound
	}

	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	secretArea := SecretArea{
		Id:          data.Id,
		Title:       data.Title,
		Description: data.Description,
		Author:      data.Author,
	}

	for _, v := range data.Area.P {
		secretArea.Area = append(secretArea.Area, [2]float32{float32(v.X), float32(v.Y)})
	}

	return secretArea, err
}

type getDbByIdDb struct {
	Id          int
	CreatedAt   *time.Time `db:"created_at"`
	Author      int
	Title       string
	Area        pgtype.Polygon
	Description string
}

func (r *repo) GetDbById(id int) (Db, error) {
	q := query.New[getDbByIdDb](r.ctx, r.db).
		Select("*", "").
		From("secret_areas", "sa").
		Where(query.EQUEL, "sa.id", id)

	data, err := q.One()

	if errors.Is(err, pgx.ErrNoRows) {
		return Db{}, app_error.ErrNotFound
	}

	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	serviceArea := Db{
		Id:          data.Id,
		Title:       data.Title,
		Description: data.Description,
		Author:      data.Author,
	}

	for _, v := range data.Area.P {
		serviceArea.Area = append(serviceArea.Area, [2]float32{float32(v.X), float32(v.Y)})
	}

	return serviceArea, err
}

func (r *repo) List(params QueryParams) (query.List[SecretArea], error) {
	q := query.New[SecretArea](r.ctx, r.db).
		Select("sa.id", "").
		AddSelect("sa.created_at", "").
		AddSelect("sa.title", "").
		AddSelect("sa.description", "").
		AddSelect("null", "author").
		AddSelect("null", "area").
		From("secret_areas", "sa").
		FilterWhere(params.FindType, "sa.title", params.Find).
		OrderBy("sa.title").
		Limit(params.Limit).
		Offset(params.Offset)

	data, err := q.GetList()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) ExistsByUnique(id int, title string) (bool, error) {
	q := query.New[SecretArea](r.ctx, r.db).
		Select("id", "").
		From("secret_areas", "").
		Where(query.EQUEL, "title", title).
		AndWhere(query.NOT_EQUEL, "id", id)

	data, err := q.Exists()
	r.logger.DebugOrError(err, q.GetLogSql().SetResult(data).SetError(err).GetMsg())

	return data, err
}

func (r *repo) DeleteById(id int) (bool, error) {
	q := `DELETE FROM secret_areas where id = $1`

	commandTag, err := r.db.Exec(r.ctx, q, id)
	r.logger.DebugOrError(err, query.NewLogSql(q, id).SetResult(commandTag.RowsAffected() > 0).SetError(err).GetMsg())
	return commandTag.RowsAffected() > 0, err
}
