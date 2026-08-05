package model

type MetaData struct {
	ID            string
	Name          string
	Size          int64
	MIMEType      string
	UploadTime    date
	LastDownload  date
	DownloadCount int
}
