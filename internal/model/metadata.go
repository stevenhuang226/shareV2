package model

import "time"

type MetaData struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Size          int64  `json:"size"`
	MIMEType      string `json:"mime_type"`
	UploadTime    time.Time
	LastDownload  time.Time
	DownloadCount int
}
