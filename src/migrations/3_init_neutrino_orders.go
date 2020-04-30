package main

import (
	"fmt"
	"github.com/go-pg/migrations/v7"
	"github.com/ventuary-lab/cache-updater/src/entities"
)


func init () {
	const TABLE_NAME = entities.NEUTRINO_ORDERS_NAME

	migrations.MustRegisterTx(
		func(db migrations.DB) error {
			fmt.Printf("creating %v table...\n", TABLE_NAME)
			_, err := db.Exec(fmt.Sprintf(
				`CREATE TABLE %[1]v (
					currency text,
					owner text,
					status text,
					type text,
					height bigint,
					order_next text,
					order_prev text,
					total bigint,
					resttotal bigint,
					is_first bool,
					is_last bool,
					order_id text PRIMARY KEY
				);
				`, TABLE_NAME))
			return err
		},
		func(db migrations.DB) error {
			fmt.Printf("dropping %v table...\n", TABLE_NAME)
			_, err := db.Exec(fmt.Sprintf(`DROP TABLE %v`, TABLE_NAME))
			return err
		},
	)
}