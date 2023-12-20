package service

import (
	"context"
	app_interface "seal/internal/app/interface"
	"seal/internal/repository/pg/query"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Params struct {
	Ctx      context.Context
	Db       *pgxpool.Pool
	Logger   app_interface.Logger
	StopChan chan bool
}

func RunCheckTelemetry(params Params) {
	go func() {
		waitLoopNumber := 0

		for {
			select {
			case <-params.StopChan:
				params.Logger.Debug("Stop RunCheckTelemetry")
				return
			case <-time.After(100 * time.Millisecond):
				if waitLoopNumber == 30 {
					doCheckTelemetry(params)

					waitLoopNumber = 0
					continue
				}

				waitLoopNumber++
			}
		}
	}()
}

func doCheckTelemetry(params Params) {
	q := `with t as (
	select t.dev_time, sh.modem , sh.route, longitude, latitude 
	from coordinates t
	inner join shipping sh on sh.modem = t.modem and sh.time_start is not null 
		and t.dev_time >= sh.time_start and (t.dev_time <= sh.time_end or sh.time_end is null)
	where 
		dev_time > NOW() - INTERVAL '1 DAY' 
		and min_distance_to_route is null
		and (select sa.id  
			from coordinates t1
			left join secret_areas sa on sa.area@>point(t1.latitude, t1.longitude)
			where dev_time >= sh.time_start and (dev_time <= sh.time_end or sh.time_end is null) and modem = t.modem
			order by dev_time desc limit 1) is null
	order by dev_time
	limit 5000
)
update coordinates 
	set min_distance_to_route =  floor(sub.min_distance_to_route)
	from (
		select 
			t.*, r.*, 
			(point(t.longitude, t.latitude)<@>r.point)*1609.34 min_distance_to_route 
		from t
		inner join lateral 
			(select point(rp.longitude, rp.latitude) point 
				from route_points rp 
				where route = t.route 
				order by point(rp.longitude, rp.latitude) <-> point(t.longitude, t.latitude) 
				limit 1) r ON true
	) sub
	where 
		coordinates.dev_time = sub.dev_time
		and coordinates.modem = sub.modem 
		and  coordinates.longitude = sub.longitude 
		and coordinates.latitude = sub.latitude`

	commandTag, err := params.Db.Exec(params.Ctx, q)
	params.Logger.DebugOrError(err, query.NewLogSql(q).SetResult(commandTag.RowsAffected() > 0).SetError(err).GetMsg())
}
