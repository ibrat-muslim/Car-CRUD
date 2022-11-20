package main

import (
	"database/sql"
	"fmt"
	"time"
)

type Car struct {
	ID        int64
	Name      string
	Price     float64
	Color     string
	Year      int64
	ImageUrl  string
	CreatedAt time.Time
	Images    []*CarImage
}

type CarImage struct {
	ID             int64
	ImageUrl       string
	SequenceNumber int32
}

type GetCarsParams struct {
	Limit  int32
	Page   int32
	Search string
}

type GetCarsResponse struct {
	Cars  []*Car
	Count int32
}

func (c *DBManager) Create(car *Car) (int64, error) {

	tx, err := c.db.Begin()
	if err != nil {
		return 0, err
	}

	query := `
		INSERT INTO cars (
			name,
			price,
			color,
			year,
			image_url
		) VALUES($1, $2, $3, $4, $5)
		RETURNING id
	`

	row := tx.QueryRow(
		query,
		car.Name,
		car.Price,
		car.Color,
		car.Year,
		car.ImageUrl,
	)

	var carID int64

	err = row.Scan(&carID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	queryInsertImage := `
		INSERT INTO car_images (
			car_id,
			image_url,
			sequence_number
		) VALUES($1, $2, $3)
	`

	for _, image := range car.Images {
		_, err := tx.Exec(
			queryInsertImage,
			carID,
			image.ImageUrl,
			image.SequenceNumber,
		)

		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return carID, nil
}

func (c *DBManager) Get(id int64) (*Car, error) {

	query := `
		SELECT
			id,
			name,
			price,
			color,
			year,
			image_url,
			created_at
		FROM cars
		WHERE id = $1
	`

	row := c.db.QueryRow(query, id)

	var car Car

	err := row.Scan(
		&car.ID,
		&car.Name,
		&car.Price,
		&car.Color,
		&car.Year,
		&car.ImageUrl,
		&car.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	queryImages := `
		SELECT 
			id,
			image_url,
			sequence_number
		FROM car_images
		WHERE car_id = $1
	`

	rows, err := c.db.Query(queryImages, id)

	car.Images = make([]*CarImage, 0)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var image CarImage

		err := rows.Scan(
			&image.ID,
			&image.ImageUrl,
			&image.SequenceNumber,
		)

		if err != nil {
			return nil, err
		}

		car.Images = append(car.Images, &image)
	}

	return &car, nil
}

func (c *DBManager) GetAll(params *GetCarsParams) (*GetCarsResponse, error) {

	var result GetCarsResponse

	result.Cars = make([]*Car, 0)

	filter := ""

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)

	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter = fmt.Sprintf(" WHERE name ILIKE '%s' OR color ILIKE '%s'", str, str)
	}

	query := `
		SELECT
			id,
			name,
			price,
			color,
			year,
			image_url,
			created_at
		FROM cars
		` + filter + `
		ORDER BY created_at DESC
		` + limit

	rows, err := c.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var car  Car

		err := rows.Scan(
			&car.ID,
			&car.Name,
			&car.Price,
			&car.Color,
			&car.Year,
			&car.ImageUrl,
			&car.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		result.Cars = append(result.Cars, &car)
	}

	queryCount := `SELECT count(1) FROM cars ` + filter

	err = c.db.QueryRow(queryCount).Scan(&result.Count)

	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

func (c *DBManager) Update(car *Car) (error) {

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	
	query := `
		UPDATE cars SET
			name = $1,
			price = $2,
			color = $3,
			year = $4,
			image_url = $5
		WHERE id = $6
	`

	result, err := tx.Exec(
		query,
		car.Name,
		car.Price,
		car.Color,
		car.Year,
		car.ImageUrl,
		car.ID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	rowsCount, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsCount == 0 {
		tx.Rollback()
		return sql.ErrNoRows
	}

	queryDeleteImages := `DELETE FROM car_images WHERE car_id = $1`

	_, err = tx.Exec(queryDeleteImages, car.ID)

	if err != nil {
		tx.Rollback()
		return err
	}

	queryInsertImage := `
		INSERT INTO car_images (
			car_id,
			image_url,
			sequence_number
		) VALUES($1, $2, $3)
	`

	for _, image := range car.Images {
		_, err := tx.Exec(
			queryInsertImage,
			car.ID,
			image.ImageUrl,
			image.SequenceNumber,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *DBManager) Delete(id int64) (error) {

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	queryDeleteImages := `DELETE FROM car_images WHERE car_id = $1`

	_, err = tx.Exec(queryDeleteImages, id)

	if err != nil {
		tx.Rollback()
		return err
	}

	quertDelete := `DELETE FROM cars WHERE id = $1`

	result, err := tx.Exec(quertDelete, id)

	if err != nil {
		tx.Rollback()
		return err
	}

	rowsCount, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsCount == 0 {
		tx.Rollback()
		return sql.ErrNoRows
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}