package models

type ItemsInfo struct {
	ODataContext   string      `json:"@odata.context"`
	ODataDeltaLink string      `json:"@odata.deltaLink"`
	ItemInfos      []*ItemInfo `json:"value"`
}

type ItemInfo struct {
	ODataType            string           `json:"@odata.type,omitempty"`
	CreatedDateTime      string           `json:"createdDateTime,omitempty"`
	Etag                 string           `json:"eTag,omitempty"`
	ID                   string           `json:"id,omitempty"`
	LastModifiedDateTime string           `json:"lastModifiedDateTime,omitempty"`
	Name                 string           `json:"name,omitempty"`
	WebUrl               string           `json:"webUrl,omitempty"`
	Size                 int64            `json:"size,omitempty"`
	ParentReference      *ParentReference `json:"parentReference,omitempty"`
	Folder               *Folder          `json:"folder,omitempty"`
	CreatedBy            *CreatedBy       `json:"createdBy,omitempty"`
	FileSystemInfo       *FileSystemInfo  `json:"fileSystemInfo,omitempty"`
}

type ParentReference struct {
	DriveId   string `json:"driveId,omitempty"`
	DriveType string `json:"driveType,omitempty"`
	ID        string `json:"id,omitempty"`
	Path      string `json:"path,omitempty"`
}

type Folder struct {
	ChildCount int `json:"childCount,omitempty"`
}

type CreatedBy struct {
	User User `json:"user,omitempty"`
}

type User struct {
	Email       string `json:"email,omitempty"`
	ID          string `json:"id,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type FileSystemInfo struct {
	CreatedDateTime      string `json:"createdDateTime,omitempty"`
	LastModifiedDateTime string `json:"lastModifiedDateTime,omitempty"`
}

type ErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
