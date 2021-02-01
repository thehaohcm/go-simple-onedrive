package models

type ItemsInfoResponse struct {
	ODataContext   string   `json:"@odata.context"`
	ODataDeltaLink string   `json:"@odata.deltaLink"`
	Value          []*Value `json:"value"`
}

type Value struct {
	ODataType            string           `json:"@odata.type"`
	CreatedDateTime      string           `json:"createdDateTime"`
	Etag                 string           `json:"eTag"`
	ID                   string           `json:"id"`
	LastModifiedDateTime string           `json:"lastModifiedDateTime"`
	Name                 string           `json:"name"`
	WebUrl               string           `json:"webUrl"`
	Size                 int64            `json:"size"`
	ParentReference      *ParentReference `json:"parentReference"`
	Folder               *Folder          `json:"folder"`
	CreatedBy            `json:"createdBy"`
	FileSystemInfo       *FileSystemInfo `json:"fileSystemInfo"`
}

type ParentReference struct {
	DriveId   string `json:"driveId"`
	DriveType string `json:"driveType"`
	ID        string `json:"id"`
	Path      string `json:"path"`
}

type Folder struct {
	ChildCount int `json:"childCount"`
}

type CreatedBy struct {
	User User `json:"user"`
}

type User struct {
	Email       string `json:"email"`
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type FileSystemInfo struct {
	CreatedDateTime      string `json:"createdDateTime"`
	LastModifiedDateTime string `json:"lastModifiedDateTime"`
}
