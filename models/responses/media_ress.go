package responses

import (
	"os"
	"time"

	"github.com/notefan-golang/config"
	"github.com/notefan-golang/models/entities"
	"gopkg.in/guregu/null.v4"
)

type Media struct {
	Id             string      `json:"id"`
	ModelType      string      `json:"model_type"`
	ModelId        string      `json:"model_id"`
	CollectionName string      `json:"collection_name"`
	Name           null.String `json:"name"`
	FileName       string      `json:"file_name"`
	FileURL        string      `json:"file_url"`
	MimeType       string      `json:"mime_type"`
	Disk           string      `json:"disk"`
	ConversionDisk null.String `json:"conversion_disk"`
	Size           int64       `json:"size"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      null.Time   `json:"updated_at"`
}

func NewMediaFromEntity(mediaEntity entities.Media) Media {
	disk := config.FSDisks[mediaEntity.Disk]
	fileURL := os.Getenv("APP_URL") + disk.URL + "/medias/" + mediaEntity.Id.String() + "/" + mediaEntity.FileName

	return Media{
		Id:             mediaEntity.Id.String(),
		ModelType:      mediaEntity.ModelType,
		ModelId:        mediaEntity.ModelId.String(),
		CollectionName: mediaEntity.CollectionName,
		Name:           null.NewString(mediaEntity.Name.String, mediaEntity.Name.Valid),
		FileName:       mediaEntity.FileName,
		FileURL:        fileURL,
		MimeType:       mediaEntity.MimeType,
		Disk:           mediaEntity.Disk,
		ConversionDisk: null.NewString(mediaEntity.ConversionDisk.String, mediaEntity.ConversionDisk.Valid),
		Size:           mediaEntity.Size,
		CreatedAt:      mediaEntity.CreatedAt,
		UpdatedAt:      null.NewTime(mediaEntity.UpdatedAt.Time, mediaEntity.UpdatedAt.Time.IsZero()),
	}
}
