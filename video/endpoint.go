package video

import (
	"context"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MetadataParams struct {
	Title   string
	MongoId string
}

type Metadata struct {
	Id      string
	MongoId []byte
	Title   string
}

type MetadataCollection struct {
	Metadatas []Metadata
}

//encore:api public path=/video/:id
func (s *VideoService) GetVideo(ctx context.Context, id int) (*Metadata, error) {
	response, err := selectMetadata(ctx, id)
	mongoidHex := hex.EncodeToString(response.MongoId)
	s.getVideoById(mongoidHex)

	return &response, err
}

//encore:api public method=GET path=/video
func (s *VideoService) GetAllVideos(ctx context.Context) (*MetadataCollection, error) {
	m, err := selectAllMetadata(ctx)

	for _, data := range m.Metadatas {
		log.Println(hex.EncodeToString(data.MongoId))
	}
	return &m, err
}

//encore:api public method=POST path=/video
// func (s *VideoService) PostVideoMetadata(ctx context.Context, mdata *MetadataParams) error {
// 	log.Println(mdata.MongoId)
// 	return insertMetadata(ctx, mdata)
// }

//encore:api public method=DELETE path=/video
func (s *VideoService) DeleteAllMetadata(ctx context.Context) error {
	return deleteAllMetadata(ctx)
}

// NOTE: this endpoint does not work
//
//encore:api public raw method=POST path=/video
func (s *VideoService) PostVideo(w http.ResponseWriter, req *http.Request) {
	const maxUploadSize = 500 << 20 // 500MB

	req.Body = http.MaxBytesReader(w, req.Body, maxUploadSize)

	err := req.ParseMultipartForm(maxUploadSize)
	if err != nil {
		if err.Error() == "http: request body too large" {
			http.Error(w, "File too large. Maximum size is 500 MB", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
		}

		return
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		log.Println(err)
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Buffer is set to 512 bytes because that all that is required to determine the MIME type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		internalServerError(w)
		return
	}

	contentType := http.DetectContentType(buffer)
	if !isValidVideoContentType(contentType) {
		http.Error(w, "Invalid file type. Only video files are allowed", http.StatusUnsupportedMediaType)
		return
	}

	if !isValidVideoExtension(header.Filename) {
		http.Error(w, "Invalid file extension. Only video files are allowed", http.StatusUnsupportedMediaType)
		return
	}

	title := req.FormValue("title")

	tempFile, err := os.CreateTemp("", "upload-*.mp4")
	if err != nil {
		internalServerError(w)
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		internalServerError(w)
		return
	}

	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		internalServerError(w)
		return
	}

	var id primitive.ObjectID

	if id, err = s.uploadVideo(tempFile); err != nil {
		internalServerError(w)
		return
	}

	mdata := MetadataParams{MongoId: id.Hex(), Title: title}

	if err = insertMetadata(req.Context(), &mdata); err != nil {
		log.Println(err)
	}
}

func internalServerError(w http.ResponseWriter) {
	http.Error(w, "Something went wrong. Please try again later.", http.StatusInternalServerError)
}

func isValidVideoContentType(contentType string) bool {
	switch contentType {
	case "video/mp4", "video/mpeg", "video/ogg", "video/webm", "video/avi":
		return true
	default:
		return false
	}
}

func isValidVideoExtension(filename string) bool {
	ext := filepath.Ext(filename)
	switch ext {
	case ".mp4", ".mpeg", ".mpg", ".ogg", ".webm", ".avi":
		return true
	default:
		return false
	}
}
