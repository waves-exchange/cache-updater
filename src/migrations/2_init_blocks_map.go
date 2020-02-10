package main;

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
	"github.com/ventuary-lab/cache-updater/src/entities"
)


func init () {
	const TABLE_NAME = entities.BLOCKS_MAP_NAME

	migrations.MustRegisterTx(
		func(db migrations.DB) error {
			fmt.Printf("creating %v table...\n", TABLE_NAME)
			_, err := db.Exec(fmt.Sprintf(
				`CREATE TABLE %[1]v (
					height bigint PRIMARY KEY,
					timestamp bigint
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