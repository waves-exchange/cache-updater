# Cache Updater Overview

Neutrino Cache updater is golang script that manages data parsing from waves blockchain to well-understood database models.

The idea behind this script to manage rapidly changing and scalable data (Neutrino Auction and Liquidate orders) up-to-date.

## Scripts reference
### `process.sh`
Process file that builds systemd service and starts it. Includes various parameters. For default behaviour use `--build-n-start`. Default systemd service name is `neutrino-cache-daemon`. Can be override using `--service` param. 
> bash ./process.sh --build-n-start
### `typegen.sh`

Generates types into swagger-types directory depending on models declared in swagger.json provided by `https://nodes.wavesplatform.com` 
> bash ./typegen.sh

### `migrations.sh`

Manages migrations. Including initial table declaring and dropping.
possible params: `--init-migration`, `--run-migration`,  `--reset-migration`

> Note: go-pg library create it's own table for managing migration stage

# Deployment

1. Firstly, we run migration, if database is empty.
> bash ./migrations.sh --run-migration
2. Compile go
> go build
3. Setup and start systemd service
> bash ./process.sh --build-n-start
4. Check status
> systemctl status neutrino-cache-daemon

## Environment file ref

| Param | Desired Value | Description
|-------|-------|-----|
| `DB_USERNAME` | `string` | `postgres username`
| `DB_PASS` | `string` | `postgres pass`
| `DB_NAME` | `string` | `Database name`
| `NODE_URL` | `URL` | `Waves Full Node URL`
| `DAPP_ADDRESS` | `waves address` | `Neutrino Auction Contract Address`
| `DB_HOST` | `string` | `Database host`
| `DB_PORT` | `string` | `Database port`
| `UPDATE_FREQUENCY` | `Milliseconds` | `Table update frequency`
