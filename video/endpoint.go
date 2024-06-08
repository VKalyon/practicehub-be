package video

import (
	"context"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"encore.dev/rlog"
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

	return &m, err
}

//encore:api public method=DELETE path=/video
func (s *VideoService) DeleteAllMetadata(ctx context.Context) error {
	return deleteAllMetadata(ctx)
}

//encore:api public raw method=POST path=/video
func (s *VideoService) PostVideo(w http.ResponseWriter, req *http.Request) {
	const maxUploadSize = 500 << 20 // 500MB
	const maxMemory = 64 << 20      // 64MB

	// sets the max permitted size of the body of the request to 500 MB
	req.Body = http.MaxBytesReader(w, req.Body, maxUploadSize)

	// sets the max size allocated to memory, after which it starts writing to disk
	err := req.ParseMultipartForm(maxMemory)
	if err != nil {
		if err.Error() == "http: request body too large" {
			http.Error(w, "File too large. Maximum size is 500 MB", http.StatusRequestEntityTooLarge)
			rlog.Error(err.Error())
		} else {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			rlog.Error(err.Error())
		}

		return
	}
	defer req.MultipartForm.RemoveAll()

	file, header, err := req.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		rlog.Error(err.Error())
		return
	}
	defer file.Close()

	// Buffer is set to 512 bytes because that all that is required to determine the MIME type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		rlog.Error(err.Error())
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
		rlog.Error("Invalid content type. Only video files are allowed")
		return
	}

	if !isValidVideoExtension(header.Filename) {
		http.Error(w, "Invalid file extension. Only video files are allowed", http.StatusUnsupportedMediaType)
		rlog.Error("Invalid file extension. Only video files are allowed")
		return
	}

	title := req.FormValue("title")

	tempFile, err := os.CreateTemp("", "upload-*.mp4")
	if err != nil {
		rlog.Error("Error occurred while creating temporary file")
		internalServerError(w)
		return
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		rlog.Error("Error occurred while copying file")
		internalServerError(w)
		return
	}

	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		rlog.Error("Error occurred while seeking start of file")
		internalServerError(w)
		return
	}

	var id primitive.ObjectID

	if id, err = s.uploadVideo(tempFile); err != nil {
		rlog.Error(err.Error())
		internalServerError(w)
		return
	}

	mdata := MetadataParams{MongoId: id.Hex(), Title: title}

	if err = insertMetadata(req.Context(), &mdata); err != nil {
		rlog.Error(err.Error())
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
