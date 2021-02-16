package models

type UploadFinishedResponse struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Size        int    `json:"size,omitempty"`
	DownloadUrl string `json:"@content.downloadUrl,omitempty"`
}
