package video

import (
	"context"
	"encoding/hex"

	"encore.dev/storage/sqldb"
)

var db = sqldb.NewDatabase("video", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

func insertMetadata(ctx context.Context, p *MetadataParams) error {
	bytes, err := hex.DecodeString(p.MongoId)
	if err == nil {
		_, err = db.Exec(ctx, `
        INSERT INTO metadata
        VALUES (DEFAULT, $1, $2)
    `, p.Title, bytes)
	}

	return err
}

func selectMetadata(ctx context.Context, id int) (Metadata, error) {
	m := Metadata{}

	err := db.QueryRow(ctx, `
		SELECT id, title, mongoid
		FROM metadata
		WHERE id = $1
	`, id).Scan(&m.Id, &m.Title, &m.MongoId)

	return m, err
}

func selectAllMetadata(ctx context.Context) (MetadataCollection, error) {
	m := MetadataCollection{}

	rows, err := db.Query(ctx, `
		SELECT * FROM metadata
	`)

	for rows.Next() {
		var metadata Metadata
		if err := rows.Scan(&metadata.Id, &metadata.Title, &metadata.MongoId); err != nil {
			return m, err
		}
		m.Metadatas = append(m.Metadatas, metadata)
	}

	return m, err
}

func deleteAllMetadata(ctx context.Context) error {
	_, err := db.Exec(ctx, `
        TRUNCATE metadata;
		DELETE FROM metadata;
    `)

	return err
}
