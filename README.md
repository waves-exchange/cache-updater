# Cache Updater Overview

Neutrino Cache updater is golang script that manages data parsing from waves blockchain to well-understood database models.

The idea behind this script to manage rapidly changing data (Neutrino Auction and Liquidate orders) up-to-date. 

## Scripts reference
### `process.sh`
Process file that builds systemd service and starts it. Includes various parameters. For default behaviour use `--build-n-start`
> bash ./process.sh --build-n-start
### `typegen.sh`

Generates types into swagger-types directory depending on models declared in swagger.json provided by `https://nodes.wavesplatform.com` 
> bash ./typegen.sh


# Deployment

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
