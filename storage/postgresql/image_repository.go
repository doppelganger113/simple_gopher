package postgresql

import (
	"api/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

type ImageRepo struct {
	database *Database
}

func NewImageRepository(db *Database) *ImageRepo {
	return &ImageRepo{database: db}
}

func (repo ImageRepo) Get(ctx context.Context, limit, offset int, order storage.Order) (storage.ImageList, error) {
	query := `SELECT
 id, name, format, original, domain, path, sizes, created_at, updated_at, author_id
 FROM images
 ORDER BY created_at ` + string(order) + `
 LIMIT $1
 OFFSET $2
`
	rows, err := repo.database.dbPool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed querying images: %w", err)
	}
	defer rows.Close()

	var imageList []storage.Image

	for rows.Next() {
		var id, name, format, original, domain, path, authorId string
		var sizes storage.ImageSizes
		var createdAt, updatedAt *time.Time

		err = rows.Scan(
			&id, &name, &format, &original, &domain, &path, &sizes, &createdAt, &updatedAt, &authorId,
		)
		if err != nil {
			return nil, fmt.Errorf("failed scaning images: %w", err)
		}
		imageList = append(imageList, storage.Image{
			Id:        id,
			Name:      name,
			Format:    storage.ImageFormat(format),
			Original:  original,
			Domain:    domain,
			Path:      path,
			Sizes:     sizes,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			AuthorId:  authorId,
		})
	}

	if imageList == nil {
		return storage.ImageList{}, nil
	}

	return imageList, nil
}

func (repo *ImageRepo) GetOne(ctx context.Context, imageId string) (storage.Image, error) {
	query := `SELECT
id, name, format, original, domain, path, sizes, created_at, updated_at, author_id
FROM images
WHERE id = $1
LIMIT 1
`
	var image storage.Image

	err := repo.database.dbPool.QueryRow(ctx, query, imageId).Scan(
		&image.Id, &image.Name, &image.Format, &image.Original, &image.Domain, &image.Path,
		&image.Sizes, &image.CreatedAt, &image.UpdatedAt, &image.AuthorId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.Image{}, storage.NotFound{}
		}
		return storage.Image{}, err
	}

	return image, nil
}

func (repo *ImageRepo) GetOneByName(ctx context.Context, name string) (storage.Image, error) {
	return storage.Image{}, nil
}

func (repo *ImageRepo) DoesImageExist(ctx context.Context, name string) (bool, error) {
	query := "SELECT name FROM images WHERE name = $1 LIMIT 1"

	var imageName string
	err := repo.database.dbPool.QueryRow(ctx, query, name).Scan(&imageName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return imageName != "", nil
}

func (repo *ImageRepo) Create(ctx context.Context, image storage.Image) (storage.Image, error) {
	query := `INSERT INTO
 images ("name", "format", "original", "domain", "path", "sizes", "author_id")
 VALUES ($1, $2, $3, $4, $5, $6, $7)
 RETURNING id, name, format, original, domain, path, sizes, created_at, updated_at, author_id
`
	data, err := json.Marshal(image.Sizes)
	if err != nil {
		return storage.Image{}, err
	}

	var id, name, format, original, domain, path, sizes, authorId string
	var createdAt, updatedAt *time.Time

	err = repo.database.dbPool.QueryRow(
		ctx,
		query,
		image.Name,
		image.Format,
		image.Original,
		image.Domain,
		image.Path,
		string(data),
		image.AuthorId,
	).Scan(
		&id, &name, &format, &original, &domain, &path, &sizes, &createdAt, &updatedAt, &authorId,
	)

	var sizesConverted storage.ImageSizes
	if len(sizes) > 0 {
		err = json.Unmarshal([]byte(sizes), &sizesConverted)
		if err != nil {
			return storage.Image{}, err
		}
	}

	createdImage := storage.Image{
		Id:        id,
		Name:      name,
		Format:    storage.ImageFormat(format),
		Original:  original,
		Domain:    domain,
		Path:      path,
		Sizes:     sizesConverted,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		AuthorId:  authorId,
	}

	return createdImage, err
}

func (repo *ImageRepo) SetNameById(ctx context.Context, imageId, newName string) (storage.Image, error) {
	return storage.Image{}, nil
}

func (repo *ImageRepo) UpdateOne(ctx context.Context, updates storage.Image) error {
	return nil
}

func (repo *ImageRepo) DeleteOne(ctx context.Context, imageId string) error {
	query := "DELETE FROM images WHERE id = $1"

	commandTag, err := repo.database.dbPool.Exec(ctx, query, imageId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.NotFound{}
		}
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return storage.NotFound{}
	}

	return nil
}

func (repo *ImageRepo) InsertMany(ctx context.Context, images storage.ImageList) (count int64, err error) {
	for _, image := range images {
		if _, err = repo.Create(ctx, image); err != nil {
			return 0, err
		}
	}
	return int64(len(images)), nil
}

func (repo *ImageRepo) DeleteAll(ctx context.Context) (rowsAffected int64, err error) {
	query := "DELETE FROM images"
	cmdTag, err := repo.database.dbPool.Exec(ctx, query)
	if err != nil {
		return 0, err
	}

	rowsAffected = cmdTag.RowsAffected()
	return
}
