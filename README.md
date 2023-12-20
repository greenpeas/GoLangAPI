# Main API

## Migration
`go run cmd/migrate/main.go create name`
`go run cmd/migrate/main.go down`

## Swagger install
`go install github.com/swaggo/swag/cmd/swag@latest`


### Generate swagger docs
`swag init -g .\cmd\main.go --outputTypes json`

### GOPRIVATE

`go env -w GOPRIVATE=gitlab.kvant.online/*`

and

`git config --global url."ssh://git@192.168.158.80".insteadOf "https://gitlab.kvant.online"`

## DB extentions
```
CREATE EXTENSION IF NOT EXISTS "btree_gist";
CREATE EXTENSION cube;
CREATE EXTENSION earthdistance;
```

### Update contracts
`go get -u gitlab.kvant.online/seal/grpc-contracts`

### Порядок обновления
`driver -> gateway -> api (если есть миграции) -> sync -> syncRaw`