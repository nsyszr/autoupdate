package postgres

import (
	"database/sql"
	"fmt"

	"github.com/insys-icom/autoupdate/pkg/repository"
)

// DAO represents a PostgreSQL implementation of DAO.
type DAO struct {
	DB *sql.DB
}

// NewDAO creates a PostgreSQL DAO object
func NewDAO(db *sql.DB) *DAO {
	return &DAO{
		DB: db,
	}
}

// GetConnectionByID returns a connection for a given schema and ID
func (dao *DAO) GetConnectionByID(schema string, id int64) (*repository.Connection, error) {
	sqlStmt := fmt.Sprintf(
		`SELECT id, title, api_url, api_key, api_secret, api_use_ssl, region 
         FROM %s.connections 
         WHERE id = $1`, schema)

	entity := &repository.Connection{}
	if err := dao.DB.QueryRow(sqlStmt, id).Scan(&entity.ID,
		&entity.Title, &entity.APIURL, &entity.APIKey, &entity.APISecret,
		&entity.APIUseSSL, &entity.Region); err != nil {
		return nil, err
	}

	return entity, nil
}

// CreateConnection creates a new connection in the given schema
func (dao *DAO) CreateConnection(schema string, entity repository.Connection) (*repository.Connection, error) {
	sqlStmt := fmt.Sprintf(
		`INSERT INTO %s.connections (title, api_url, api_key, api_secret, 
            api_use_ssl, region) 
         VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, schema)

	var id int64
	if err := dao.DB.QueryRow(sqlStmt, entity.Title, entity.APIURL,
		entity.APIKey, entity.APISecret, entity.APIUseSSL, entity.Region).Scan(&id); err != nil {
		return nil, err
	}

	if id == 0 {
		return nil, fmt.Errorf("invalid ID returned by SQL INSERT statement")
	}

	return dao.GetConnectionByID(schema, id)
}

// GetRepositoryByID returns a repository for a given schema and ID
func (dao *DAO) GetRepositoryByID(schema string, id int64) (*repository.Repository, error) {
	sqlStmt := fmt.Sprintf(
		`SELECT id, connection_id, title, bucket_name, object_prefix, base_url 
         FROM %s.package_repositories 
         WHERE id = $1`, schema)

	entity := &repository.Repository{}
	if err := dao.DB.QueryRow(sqlStmt, id).Scan(&entity.ID,
		&entity.ConnectionID, &entity.Title, &entity.BucketName,
		&entity.ObjectPrefix, &entity.BaseURL); err != nil {
		return nil, err
	}

	return entity, nil
}

// CreateRepository creates a new repository in the given schema
func (dao *DAO) CreateRepository(schema string, entity repository.Repository) (*repository.Repository, error) {
	sqlStmt := fmt.Sprintf(
		`INSERT INTO %s.package_repositories (connection_id, title, bucket_name,
            object_prefix, base_url) 
         VALUES ($1, $2, $3, $4, $5) RETURNING id`, schema)

	var id int64
	if err := dao.DB.QueryRow(sqlStmt, entity.ConnectionID, entity.Title,
		entity.BucketName, entity.ObjectPrefix, entity.BaseURL).Scan(&id); err != nil {
		return nil, err
	}

	if id == 0 {
		return nil, fmt.Errorf("invalid ID returned by SQL INSERT statement")
	}

	return dao.GetRepositoryByID(schema, id)
}

// GetUpdatePacketByID returns a update packet for a given schema and ID
func (dao *DAO) GetUpdatePacketByID(schema string, id int64) (*repository.UpdatePacket, error) {
	sqlStmt := fmt.Sprintf(
		`SELECT id, repository_id, title, object_prefix, object_name 
         FROM %s.packages 
         WHERE id = $1`, schema)

	entity := &repository.UpdatePacket{}
	if err := dao.DB.QueryRow(sqlStmt, id).Scan(&entity.ID,
		&entity.RepositoryID, &entity.Title, &entity.ObjectPrefix,
		&entity.ObjectName); err != nil {
		return nil, err
	}

	return entity, nil
}

// CreateUpdatePacket creates a new update packet in the given schema
func (dao *DAO) CreateUpdatePacket(schema string, entity repository.UpdatePacket) (*repository.UpdatePacket, error) {
	sqlStmt := fmt.Sprintf(
		`INSERT INTO %s.packages (repository_id, title, object_prefix,
            object_name) 
         VALUES ($1, $2, $3, $4) RETURNING id`, schema)

	var id int64
	if err := dao.DB.QueryRow(sqlStmt, entity.RepositoryID, entity.Title,
		entity.ObjectPrefix, entity.ObjectPrefix).Scan(&id); err != nil {
		return nil, err
	}

	if id == 0 {
		return nil, fmt.Errorf("invalid ID returned by SQL INSERT statement")
	}

	return dao.GetUpdatePacketByID(schema, id)
}

// GetUpdatePacketArtifactByID returns a update packet artifact for a given schema and ID
func (dao *DAO) GetUpdatePacketArtifactByID(schema string, id int64) (*repository.UpdatePacketArtifact, error) {
	sqlStmt := fmt.Sprintf(
		`SELECT id, package_id, object_name, object_exists, order_index,
            file_name, file_type_id, descr, required_software, version 
         FROM %s.package_artifacts 
         WHERE id = $1`, schema)

	entity := &repository.UpdatePacketArtifact{}
	if err := dao.DB.QueryRow(sqlStmt, id).Scan(&entity.ID,
		&entity.UpdatePacketID, &entity.ObjectName, &entity.ObjectExists,
		&entity.OrderIndex, &entity.FileName, &entity.FileTypeID,
		&entity.Description, &entity.RequiredSoftware, &entity.Version); err != nil {
		return nil, err
	}

	return entity, nil
}

// CreateUpdatePacketArtifact creates a new update packet artifact in the given schema
func (dao *DAO) CreateUpdatePacketArtifact(schema string, entity repository.UpdatePacketArtifact) (*repository.UpdatePacketArtifact, error) {
	sqlStmt := fmt.Sprintf(
		`INSERT INTO %s.packages (package_id, object_name, object_exits,
            order_index, file_name, file_type_id, descr, required_software,
            version) 
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`, schema)

	var id int64
	if err := dao.DB.QueryRow(sqlStmt, entity.UpdatePacketID, entity.ObjectName,
		entity.ObjectExists, entity.OrderIndex, entity.FileName,
		entity.FileTypeID, entity.Description, entity.RequiredSoftware,
		entity.Version).Scan(&id); err != nil {
		return nil, err
	}

	if id == 0 {
		return nil, fmt.Errorf("invalid ID returned by SQL INSERT statement")
	}

	return dao.GetUpdatePacketArtifactByID(schema, id)
}

// ExistsTargetBySlug checks if a target for given schema and slug exists
func (dao *DAO) ExistsTargetBySlug(schema, slug string) (bool, error) {
	sqlStmt := fmt.Sprintf(
		`SELECT count(id)
         FROM %s.update_targets 
         WHERE slug = $1`, schema)

	var count int64
	if err := dao.DB.QueryRow(sqlStmt, slug).Scan(&count); err != nil {
		return false, err
	}

	return (count == 1), nil
}

// GetUpdateTargetBySlug returns a target for a given schema and slug.
func (dao *DAO) GetUpdateTargetBySlug(schema, slug string) (*repository.UpdateTarget, error) {
	sqlStmt := fmt.Sprintf(
		`SELECT id, title, slug, download_list_id, is_enabled
         FROM %s.update_targets 
         WHERE slug = $1`, schema)

	entity := &repository.UpdateTarget{}
	if err := dao.DB.QueryRow(sqlStmt, slug).Scan(&entity.ID,
		&entity.Title, &entity.Slug, &entity.DownloadListID, &entity.Enabled); err != nil {
		return nil, err
	}

	return entity, nil
}

// GetDownloadListByID returns a download list for a given schema and slug.
func (dao *DAO) GetDownloadListByID(schema string, id int64) (*repository.DownloadList, error) {
	sqlStmt := fmt.Sprintf(
		`SELECT id, title, default_update_packet_id
         FROM %s.download_lists 
         WHERE id = $1`, schema)

	entity := &repository.DownloadList{}
	if err := dao.DB.QueryRow(sqlStmt, id).Scan(&entity.ID,
		&entity.Title, &entity.DefaultUpdatePacketID); err != nil {
		return nil, err
	}

	return entity, nil
}

// FindDownloadListLineItemsByDownloadListID returns all download list items for a given schema and download list ID.
func (dao *DAO) FindDownloadListLineItemsByDownloadListID(schema string, downloadListID int64) ([]*repository.DownloadListLineItem, error) {
	sqlStmt := fmt.Sprintf(
		`SELECT id, download_list_id, serial_number, update_packet_id
         FROM %s.download_list_line_items
         WHERE download_list_id = $1`, schema)

	var entites []*repository.DownloadListLineItem
	rows, err := dao.DB.Query(sqlStmt, downloadListID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		entity := &repository.DownloadListLineItem{}
		if err := rows.Scan(&entity.ID, &entity.DownloadListID,
			&entity.SerialNumber, &entity.UpdatePacketID); err != nil {
			return nil, err
		}
		entites = append(entites, entity)
	}

	return entites, nil
}
