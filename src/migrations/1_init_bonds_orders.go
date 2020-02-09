package main;

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
	"github.com/ventuary-lab/cache-updater/src/entities"
)


func init () {
	const TABLE_NAME = entities.BONDS_ORDERS_NAME

	migrations.MustRegisterTx(
		func(db migrations.DB) error {
			fmt.Printf("creating %v table...\n", TABLE_NAME)
			_, err := db.Exec(fmt.Sprintf(
				`CREATE TABLE %[1]v (
					height bigint,
					timestamp bigint,
					owner text,
					price int,
					total float8,
					filledtotal float8,
					resttotal float8,
					status text,
					index int,
					amount float8,
					filledamount float8,
					restamount float8,
					pairname text,
					type text,
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