package models

import (
	"context"
	"database/sql"
	"github.com/akashc777/csvToPdf/initializers/postgresInit"
	"time"

	"github.com/akashc777/csvToPdf/helpers"
	"github.com/google/uuid"
)

type Template struct {
	Id              uuid.UUID
	HtmlTemplate    string `json:"html_template,omitempty"`
	CreatedBy       string `json:"created_by,omitempty"`
	TemplateName    string `json:"template_name,omitempty"`
	NewTemplateName string `json:"new_template_name,omitempty"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}

// create template
func CreateTemplate(query string, createdByEmailId, htmlTemplate, templateName string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), postgresInit.DBConn.DBTimeout)
	defer cancel()
	_, err = postgresInit.DBConn.SqlDB.ExecContext(
		ctx,
		query,
		createdByEmailId,
		htmlTemplate,
		templateName,
	)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"domain/CreateTemplate Failed to create new user err: %+v",
			err)
		return
	}

	return
}

// get template
func GetTemplateByName(query string, emailID *string, template_name *string) (pdfTemplateInfo *Template, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), postgresInit.DBConn.DBTimeout)
	defer cancel()

	var templateInfo Template
	var createdAt sql.NullTime
	var updatedAt sql.NullTime
	var deletedAt sql.NullTime

	row := postgresInit.DBConn.SqlDB.QueryRowContext(ctx, query, emailID, template_name)
	err = row.Scan(
		&templateInfo.Id,
		&templateInfo.CreatedBy,
		&templateInfo.HtmlTemplate,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"models/GetTemplateByName Failed to get template by name err: %+v",
			err)
		return
	}

	if createdAt.Valid {
		templateInfo.CreatedAt = createdAt.Time
	}

	if updatedAt.Valid {
		templateInfo.UpdatedAt = updatedAt.Time
	}

	if deletedAt.Valid {
		templateInfo.DeletedAt = deletedAt.Time
	}

	return &templateInfo, nil
}

// get all template name by email
func GetTemplateNamesByEmailID(query string, emailID *string) (templateNames []*string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), postgresInit.DBConn.DBTimeout)
	defer cancel()
	rows, err := postgresInit.DBConn.SqlDB.QueryContext(ctx, query, emailID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var templateName string
		err = rows.Scan(
			&templateName,
		)

		if err != nil {
			helpers.MessageLogs.ErrorLog.Printf(
				"models/GetTemplateNamesByEmailID Failed to scan template name err: %+v",
				err)
			return
		}

		templateNames = append(templateNames, &templateName)
	}

	return
}

// update teemplate
func UpdateTemplateByName(query string, TemplateName string, emailID string, htmlTemplate string, newTemplateName string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), postgresInit.DBConn.DBTimeout)
	defer cancel()
	_, err = postgresInit.DBConn.SqlDB.ExecContext(
		ctx,
		query,
		htmlTemplate,
		newTemplateName,
		emailID,
		TemplateName,
	)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"models/UpdateTemplateByName Failed to update user name err: %+v",
			err)
		return
	}

	return
}

// delete template
func DeleteTemplateByName(query string, emailID *string, templateName *string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), postgresInit.DBConn.DBTimeout)
	defer cancel()

	_, err = postgresInit.DBConn.SqlDB.ExecContext(
		ctx,
		query,
		emailID,
		templateName,
	)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"models/DeleteTemplateByName Failed to delete template err: %+v",
			err)
		return
	}
	return
}
