package video

import (
	"context"

	"encore.dev/storage/sqldb"
)

type MetadataParams struct {
	Title string
	URL   string
}

type Metadata struct {
	Id    string
	Title string
	URL   string
}

type MetadataCollection struct {
	Metadatas []Metadata
}

var db = sqldb.NewDatabase("video", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

//encore:api public path=/video/:title
func GetVideo(ctx context.Context, title string) (*Metadata, error) {
	response, err := selectMetadata(ctx, title)

	return &response, err
}

//encore:api public method=GET path=/video
func GetAllVideos(ctx context.Context) (*MetadataCollection, error) {
	m, err := selectAllMetadata(ctx)

	return &m, err
}

//encore:api public method=POST path=/video
func PostVideo(ctx context.Context, m *MetadataParams) error {
	err := insertMetadata(ctx, m)
	return err
}

func insertMetadata(ctx context.Context, p *MetadataParams) error {
	_, err := db.Exec(ctx, `
        INSERT INTO metadata
        VALUES (DEFAULT, $1, $2)
    `, p.Title, p.URL)
	return err
}

func selectMetadata(ctx context.Context, title string) (Metadata, error) {
	m := Metadata{}

	err := db.QueryRow(ctx, `
		SELECT id, title, url
		FROM metadata
		WHERE title = $1
	`, title).Scan(&m.Id, &m.Title, &m.URL)

	return m, err
}

func selectAllMetadata(ctx context.Context) (MetadataCollection, error) {
	m := MetadataCollection{}

	rows, err := db.Query(ctx, `
		SELECT * FROM metadata
	`)

	for rows.Next() {
		var metadata Metadata
		if err := rows.Scan(&metadata.Id, &metadata.Title, &metadata.URL); err != nil {
			return m, err
		}
		m.Metadatas = append(m.Metadatas, metadata)
	}

	return m, err
}
