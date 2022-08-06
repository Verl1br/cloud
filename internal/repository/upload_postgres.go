package repository

import (
	"fmt"

	"github.com/dhevve/uploadImage/internal/models"
	"github.com/jmoiron/sqlx"
)

type UploadPostgres struct {
	db *sqlx.DB
}

func NewUploadPostgres(db *sqlx.DB) *UploadPostgres {
	return &UploadPostgres{db: db}
}

func (r *UploadPostgres) Upload(userId int, fullFileName string) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int

	createQueryImage := fmt.Sprintf("INSERT INTO %s (name) values ($1) RETURNING id", "images")

	row := tx.QueryRow(createQueryImage, fullFileName)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createQueryUsersImages := fmt.Sprintf("INSERT INTO %s (user_id, image_id) values ($1, $2) RETURNING id", "users_images")
	_, err = tx.Exec(createQueryUsersImages, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *UploadPostgres) GetAll(userId int) ([]models.Image, error) {
	var images []models.Image

	query := fmt.Sprintf("SELECT images.id, images.name FROM %s INNER JOIN %s ON images.id = users_images.image_id WHERE users_images.user_id = $1", "images", "users_images")
	err := r.db.Select(&images, query, userId)

	return images, err
}

func (r *UploadPostgres) GetById(userId, imageId int) (models.Image, error) {
	var image models.Image

	query := fmt.Sprintf("SELECT images.id, images.name FROM %s INNER JOIN %s ON images.id = users_images.image_id WHERE users_images.user_id = $1 AND users_images.image_id = $2", "images", "users_images")
	err := r.db.Get(&image, query, userId, imageId)

	return image, err
}

func (r *UploadPostgres) Delete(userId, itemId int) error {
	query := fmt.Sprintf(`DELETE FROM images USING users_images
									WHERE images.id = users_images.image_id AND users_images.user_id = $1 AND images.id = $2`)
	_, err := r.db.Exec(query, userId, itemId)
	return err
}
