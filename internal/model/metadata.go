package model

import "time"

type MetaData struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Size          int64     `json:"size"`
	MIMEType      string    `json:"mime_type"`
	UploadTime    time.Time `json:"upload_time"`
	LastDownload  time.Time `json:"last_download"`
	DownloadCount int       `json:"download_count"`
}
