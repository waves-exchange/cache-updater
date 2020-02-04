package main

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating test_bonds_orders table...")
		_, err := db.Exec(
			`CREATE TABLE test_bonds_orders (
				height text,
				timestamp bigint,
				owner text,
				price int,
				total double,
				filledtotal double,
				resttotal double,
				status text,
				index text,
				amount double,
				filledamount double,
				restamount double,
				pairname text,
				type text,
				uuid uuid,	   
			)`
		)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping test_bonds_orders table...")
		_, err := db.Exec(`DROP TABLE test_bonds_orders`)
		return err
	})
}