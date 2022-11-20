package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	PostgresUser = "postgres"
	PostfresHost = "localhost"
	PostgresPort = 5432
	PostgresPassword = 1421
	PostgresDatabase = "demo"
)

func main() {

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%d dbname=%s sslmode=disable",
		PostfresHost,
		PostgresPort,
		PostgresUser,
		PostgresPassword,
		PostgresDatabase,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open connection: %v", err)
	}

	auto := NewDBManager(db)

	id, err := auto.Create(&Car{
		Name: "Tesla Y",
		Price: 70000.0,
		Color: "White",
		Year: 2022,
		ImageUrl: "tesla_y_url",
		Images: []*CarImage{
			{
				ImageUrl: "tesla_y_url_1",
				SequenceNumber: 1,
			},
			{
				ImageUrl: "tesla_y_url_2",
				SequenceNumber: 2,
			},
			{
				ImageUrl: "tesla_y_url_3",
				SequenceNumber: 3,
			},
		},
	})

	if err != nil {
		log.Fatalf("failed to create car: %v", err)
	}

	car, err := auto.Get(id)

	if err != nil {
		log.Fatalf("failed to get a car: %v", err)
	}

	fmt.Println(car)

	resp, err := auto.GetAll(&GetCarsParams{
		Limit: 10,
		Page: 1,
		Search: "Tesla Y",
	})

	if err != nil {
		log.Fatalf("failed to get cars: %v", err)
	}

	fmt.Println(resp)

	err = auto.Update(&Car{
		ID: id,
		Name: "BMW X6",
		Price: 90000.0,
		Color: "Black",
		Year: 2023,
		ImageUrl: "bmw_url",
		Images: []*CarImage{
			{
				ImageUrl: "bmw_url_1",
				SequenceNumber: 1,
			},
			{
				ImageUrl: "bmw_url_2",
				SequenceNumber: 2,
			},
			{
				ImageUrl: "bmw_url_3",
				SequenceNumber: 3,
			},
			{
				ImageUrl: "bmw_url_4",
				SequenceNumber: 4,
			},
		},
	})

	if err != nil {
		log.Fatalf("failed to update a car: %v", err)
	}

	err = auto.Delete(id)
	if err != nil {
		log.Fatalf("failed to delete a car: %v", err)
	}
}