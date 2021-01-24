package models

type UploadFinishedResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Size        int    `json:"size"`
	DownloadUrl string `json:"@content.downloadUrl"`
}
