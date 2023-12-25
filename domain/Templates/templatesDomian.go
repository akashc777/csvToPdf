package domain

import (
	"github.com/akashc777/csvToPdf/helpers"
	"github.com/akashc777/csvToPdf/models/postgres/Templates"
)

// create template
func CreateTemplate(createdByEmailId, htmlTemplate, templateName string) (err error) {
	query := `
		INSERT INTO templates (created_by, html_template, template_name)
		VALUES ($1, $2, $3)
	`
	err = models.CreateTemplate(query, createdByEmailId, htmlTemplate, templateName)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"domain/CreateTemplate Failed to create new template err: %+v",
			err)
		return
	}

	return
}

// get template
func GetTemplateByName(emailID *string, templateName *string) (templateInfo *models.Template, err error) {
	query := `
        SELECT id, created_by, html_template, created_at, updated_at, deleted_at
        FROM templates
        WHERE created_by = $1 AND template_name = $2
    `
	templateInfo, err = models.GetTemplateByName(query, emailID, templateName)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"domain/GetTemplateByName Failed to get template err: %+v",
			err)
		return
	}

	return
}

// get all template name by email
func GetTemplateNamesByEmailID(emailID *string) (templateNames []*string, err error) {
	query := `
        SELECT template_name
        FROM templates
        WHERE created_by = $1
    `

	templateNames, err = models.GetTemplateNamesByEmailID(query, emailID)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"domain/GetTemplateNamesByEmailID Failed to get template names err: %+v",
			err)
		return
	}

	return
}

// update teemplate
func UpdateTemplateByName(TemplateName, emailID, htmlTemplate, newTemplateName string) (err error) {

	query := `
        UPDATE templates
        SET html_template = $1, template_name = $2
        WHERE created_by = $3 AND template_name = $4
    `

	err = models.UpdateTemplateByName(query, TemplateName, emailID, htmlTemplate, newTemplateName)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"domain/UpdateTemplateByName Failed to update template err: %+v",
			err)
		return
	}

	return
}

// delete template
func DeleteTemplateByName(emailID *string, templateName *string) (err error) {
	query := `
        DELETE FROM templates
        WHERE created_by = $1 AND template_name = $2
    `

	err = models.DeleteTemplateByName(query, emailID, templateName)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"domain/UpdateTemplateByName Failed to create new user err: %+v",
			err)
		return
	}

	return
}
