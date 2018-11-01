package repository

import "database/sql"

type FileType struct {
	ID    int64
	Title string
	Code  string
}

type Connection struct {
	ID        int64
	Title     string
	APIURL    string
	APIKey    string
	APISecret string
	APIUseSSL bool
	Region    string
}

type Repository struct {
	ID           int64
	ConnectionID int64
	Title        string
	BucketName   string
	ObjectPrefix string
	BaseURL      string
}

type UpdatePacket struct {
	ID           int64
	RepositoryID string
	Title        string
	ObjectPrefix string
	ObjectName   string
	// Tags         []string `pg:",array"`
}

type UpdatePacketArtifact struct {
	ID               int64
	UpdatePacketID   int64
	ObjectName       string
	ObjectExists     bool
	OrderIndex       int
	FileName         string
	FileTypeID       int64
	Description      string
	RequiredSoftware string
	Version          string
}

type DownloadList struct {
	ID                    int64
	Title                 string
	DefaultUpdatePacketID sql.NullInt64
}

type DownloadListLineItem struct {
	ID             int64
	DownloadListID int64
	SerialNumber   string
	UpdatePacketID int64
}

type UpdateTarget struct {
	ID             int64
	Title          string
	Slug           string
	DownloadListID sql.NullInt64
	Enabled        bool
}

type DAO interface {
	ExistsTargetBySlug(schema, slug string) (bool, error)
	GetUpdateTargetBySlug(schema, slug string) (*UpdateTarget, error)
	GetDownloadListByID(schema string, id int64) (*DownloadList, error)
	FindDownloadListLineItemsByDownloadListID(schema string, downloadListID int64) ([]*DownloadListLineItem, error)
}
