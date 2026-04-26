package storage

import (
	"context"

	models "github.com/qweq1232/dnd_form/internal/domane/models/char"
)

func (db *Storage) GetChar(ctx context.Context,
	id, userID int32,
) (*models.Char, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, "",
		`SELECT stats FROM chars WHERE id = $1 AND user_id = $2;`, id, userID,
	)

	var stats []byte

	if err = row.Scan(stats); err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	char := models.Char{ID: id, UserID: userID, Stats: stats}

	return &char, nil
}

func (db *Storage) GetAllChar(ctx context.Context,
	userID int32,
) ([]models.Char, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	stmt := `SELECT stats, id FROM chars WHERE user_id = $1`

	var chars []models.Char

	rows, err := tx.Query(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		char := models.Char{UserID: userID}

		if err := rows.Scan(&char.Stats, &char.ID); err != nil {
			return nil, err
		}

		chars = append(chars, char)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return chars, nil
}

func (db *Storage) SaveChar(ctx context.Context,
	stats []byte,
	userID int32,
) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	stmt := `INSERT INTO chars (stats, user_id) VALUES ($1, $2);`

	if err != nil {
		return err
	}

	if _, err = tx.Exec(ctx, stmt, stats, userID); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (db *Storage) DeleteChar(ctx context.Context, id int32) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	stmt := `DELETE FROM chars WHERE id = $1`

	if _, err = tx.Exec(ctx, stmt, id); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (db *Storage) UpdateChar(ctx context.Context,
	stats []byte, id int32,
) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	stmt := `Update chars SET stats = $1 WHERE id = $2;`

	if _, err = tx.Exec(ctx, stmt, stats, id); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
