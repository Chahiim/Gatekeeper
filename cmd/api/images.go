package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Chahiim/Gatekeeper/internal/data"
	"github.com/Chahiim/Gatekeeper/internal/validator"
)

const (
	maxUploadSize = 10 << 20 // 10 MB
	storageDir    = "./uploads/originals"
)

var allowedTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

func (app *application) createImageHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize+1<<20)

	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("file is too large (max 10 MB)"))
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("image file is required"))
		return
	}
	defer file.Close()

	mediaType := handler.Header.Get("Content-Type")
	v := validator.New()
	v.Check(handler.Filename != "", "file", "must be provided")
	v.Check(allowedTypes[mediaType], "file", "must be a JPEG or PNG image")
	v.Check(handler.Size <= maxUploadSize, "file", "must not exceed 10 MB")
	v.Check(handler.Size > 0, "file", "must not be empty")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	sniffed := http.DetectContentType(buf[:n])
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if sniffed != mediaType {
		app.badRequestResponse(w, r, fmt.Errorf("file content does not match declared type"))
		return
	}

	if err := os.MkdirAll(storageDir, 0755); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	storedName, err := generateStoredName(handler.Filename)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	dstPath := filepath.Join(storageDir, storedName)
	dst, err := os.Create(dstPath)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		os.Remove(dstPath)
		app.serverErrorResponse(w, r, err)
		return
	}

	image := &data.Image{
		OriginalName: handler.Filename,
		StoredName:   storedName,
		MediaType:    mediaType,
		FileSize:     written,
	}

	err = app.models.Images.Insert(image)
	if err != nil {
		os.Remove(dstPath)
		app.serverErrorResponse(w, r, err)
		return
	}

	job := &data.ImageJob{
		ImageID: image.ID,
		Status:  data.ImageJobStatusQueued,
	}

	err = app.models.ImageJobs.Insert(job)
	if err != nil {
		os.Remove(dstPath)
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/jobs/%d", job.ID))

	response := envelope{
		"image_id":  image.ID,
		"job_id":    job.ID,
		"status":    job.Status,
		"status_url": fmt.Sprintf("/v1/jobs/%d", job.ID),
	}

	err = app.writeJSON(w, http.StatusAccepted, response, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showJobHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	var id int64
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	job, err := app.models.ImageJobs.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	result := envelope{
		"id":          job.ID,
		"image_id":    job.ImageID,
		"status":      job.Status,
		"queued_at":   job.QueuedAt,
		"started_at":  job.StartedAt,
		"completed_at": job.CompletedAt,
	}

	if job.ErrorMessage != nil {
		result["error"] = *job.ErrorMessage
	}

	if job.Status == data.ImageJobStatusCompleted {
		variants, err := app.models.Variants.GetByImageID(job.ImageID)
		if err == nil && variants != nil {
			var variantList []envelope
			for _, v := range variants {
				variantList = append(variantList, envelope{
					"name":    v.Name,
					"width":   v.Width,
					"height":  v.Height,
					"url":     fmt.Sprintf("/v1/images/%d/variants/%s", v.ImageID, v.Name),
				})
			}
			result["variants"] = variantList
		}
	}

	err = app.writeJSON(w, http.StatusOK, result, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showVariantHandler(w http.ResponseWriter, r *http.Request) {
	imageIDStr := r.PathValue("imageID")
	var imageID int64
	_, err := fmt.Sscanf(imageIDStr, "%d", &imageID)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	name := r.PathValue("name")
	if name == "" {
		app.notFoundResponse(w, r)
		return
	}

	variants, err := app.models.Variants.GetByImageID(imageID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	for _, v := range variants {
		if v.Name == name {
			http.ServeFile(w, r, filepath.Join("./uploads/variants", v.StoredName))
			return
		}
	}

	app.notFoundResponse(w, r)
}

func generateStoredName(originalName string) (string, error) {
	ext := filepath.Ext(originalName)
	if ext == "" {
		ext = ".bin"
	}
	ext = strings.ToLower(ext)

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random name: %w", err)
	}

	return hex.EncodeToString(b) + ext, nil
}


