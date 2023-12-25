package controllers

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	tModels "github.com/akashc777/csvToPdf/models/postgres/Templates"
	"github.com/go-chi/chi/v5"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	pdfCPUModel "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	business "github.com/akashc777/csvToPdf/business/Templates"
	"github.com/akashc777/csvToPdf/helpers"
)

func Login(w http.ResponseWriter, r *http.Request) {

	var userInfo tModels.Template

	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {

		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/user Failed to parse the body of the req err: %+v",
			err)

		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envolope{
			"msg": "Failed to parse the body of the req",
			"err": err,
		})
		return

	}

	tokenString, err := business.LoginUserByEmailID(userInfo.CreatedBy)

	if err != nil {

		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/user Failed to login err: %+v",
			err)

		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envolope{
			"msg": "Failed to parse the body of the req",
			"err": err,
		})
		return

	}
	http.SetCookie(w, &http.Cookie{
		Name:    "Authorization",
		Value:   tokenString,
		Expires: time.Now().Add(time.Hour),
	})

	helpers.WriteJSON(w, http.StatusOK, helpers.Envolope{
		"msg":         "successfully loggedIn !",
		"tokenString": tokenString,
	})
}

func CreateTemplate(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userEmailID := ctx.Value(helpers.UserEmailContextKey).(string)

	var templateInfo tModels.Template
	err := json.NewDecoder(r.Body).Decode(&templateInfo)
	if err != nil {

		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/CreateTemplate Failed to parse the body of the req err: %+v",
			err)

		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envolope{
			"msg": "Failed to parse the body of the req",
			"err": err,
		})
		return

	}

	templateInfo.CreatedBy = userEmailID

	err = business.CreateTemplate(templateInfo.CreatedBy, templateInfo.HtmlTemplate, templateInfo.TemplateName)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/CreateTemplate Failed to create template err: %+v",
			err)

		er := errors.New("Failed to create template")
		helpers.ErrorJSON(w, er, http.StatusBadRequest)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envolope{
		"msg": "successfully created template!",
	})
}

func GetTemplateByName(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userEmailID := ctx.Value(helpers.UserEmailContextKey).(string)

	var inputTemplateInfo tModels.Template

	err := json.NewDecoder(r.Body).Decode(&inputTemplateInfo)
	if err != nil {

		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/CreateTemplate Failed to parse the body of the req err: %+v",
			err)

		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envolope{
			"msg": "Failed to parse the body of the req",
			"err": err,
		})
		return

	}
	inputTemplateInfo.CreatedBy = userEmailID

	templateInfo, err := business.GetTemplateByName(&inputTemplateInfo.CreatedBy, &inputTemplateInfo.TemplateName)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/CreateTemplate Failed to create template err: %+v",
			err)

		er := errors.New("Failed to create template")
		helpers.ErrorJSON(w, er, http.StatusBadRequest)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envolope{
		"msg":      "successfully fetched template!",
		"template": templateInfo,
	})
}

func GetTemplateNames(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userEmailID := ctx.Value(helpers.UserEmailContextKey).(string)

	var inputTemplateInfo tModels.Template

	err := json.NewDecoder(r.Body).Decode(&inputTemplateInfo)
	if err != nil {

		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/CreateTemplate Failed to parse the body of the req err: %+v",
			err)

		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envolope{
			"msg": "Failed to parse the body of the req",
			"err": err,
		})
		return

	}
	inputTemplateInfo.CreatedBy = userEmailID

	templateNames, err := business.GetTemplateNamesByEmailID(&inputTemplateInfo.CreatedBy)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/CreateTemplate Failed to create template err: %+v",
			err)

		er := errors.New("Failed to create template")
		helpers.ErrorJSON(w, er, http.StatusBadRequest)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envolope{
		"msg":           "successfully fetched template!",
		"templateNames": templateNames,
	})
}

func UpdateTemplate(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userEmailID := ctx.Value(helpers.UserEmailContextKey).(string)

	var inputTemplateInfo tModels.Template

	err := json.NewDecoder(r.Body).Decode(&inputTemplateInfo)
	if err != nil {

		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/CreateTemplate Failed to parse the body of the req err: %+v",
			err)

		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envolope{
			"msg": "Failed to parse the body of the req",
			"err": err,
		})
		return

	}
	inputTemplateInfo.CreatedBy = userEmailID

	err = business.UpdateTemplateByName(&inputTemplateInfo)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/CreateTemplate Failed to create template err: %+v",
			err)

		er := errors.New("Failed to create template")
		helpers.ErrorJSON(w, er, http.StatusBadRequest)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envolope{
		"msg": "successfully updated template!",
	})
}

func DeleteTemplate(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userEmailID := ctx.Value(helpers.UserEmailContextKey).(string)

	var inputTemplateInfo tModels.Template

	err := json.NewDecoder(r.Body).Decode(&inputTemplateInfo)
	if err != nil {

		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/CreateTemplate Failed to parse the body of the req err: %+v",
			err)

		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envolope{
			"msg": "Failed to parse the body of the req",
			"err": err,
		})
		return

	}
	inputTemplateInfo.CreatedBy = userEmailID

	err = business.DeleteTemplateByName(&inputTemplateInfo.CreatedBy, &inputTemplateInfo.TemplateName)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/CreateTemplate Failed to create template err: %+v",
			err)

		er := errors.New("Failed to create template")
		helpers.ErrorJSON(w, er, http.StatusBadRequest)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envolope{
		"msg": "successfully deleted template!",
	})
}

func GetPDF(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userEmailID := ctx.Value(helpers.UserEmailContextKey).(string)

	templateName := r.FormValue("templateName")

	templateInfo, err := business.GetTemplateByName(&userEmailID, &templateName)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/GetPDF Failed to get template err: %+v",
			err)

		er := errors.New("Failed to get template")
		helpers.ErrorJSON(w, er, http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("csv_file")
	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/GetPDF Failed to get csv file err: %+v",
			err)

		er := errors.New("Failed to get csv file")
		helpers.ErrorJSON(w, er, http.StatusBadRequest)

	}
	fp := csv.NewReader(file)
	records, err := fp.ReadAll()

	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/GetPDF Failed to read from csv file err: %+v",
			err)

		er := errors.New("Failed to read from csv")
		helpers.ErrorJSON(w, er, http.StatusBadRequest)
		return
	}

	var zipFileBytes bytes.Buffer
	zipWriter := zip.NewWriter(&zipFileBytes)

	funcMap := template.FuncMap{
		"splitLines": func(s string) []string {
			return strings.Split(s, "\n")
		},
	}

	tmpl, err := template.New("invoice").Funcs(funcMap).Parse(templateInfo.HtmlTemplate)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/GetPDF Failed to parse HTML template err: %+v",
			err)

		er := errors.New("Failed to parse HTML template")
		helpers.ErrorJSON(w, er, http.StatusBadRequest)
		return
	}

	header := records[0]
	dataList := make([]map[string]string, 0)

	// Parse CSV records
	for _, record := range records[1:] {
		data := make(map[string]string)
		for i, colValue := range record {
			data[header[i]] = colValue
		}
		dataList = append(dataList, data)
	}

	for i, data := range dataList {
		pdfg, err := wkhtmltopdf.NewPDFGenerator()
		if err != nil {
			helpers.MessageLogs.ErrorLog.Printf(
				"controllers/GetPDF failed to initialise new PDF generator err: %+v",
				err)

			er := errors.New("Error failed to initialise new PDF generator")
			helpers.ErrorJSON(w, er, http.StatusBadRequest)
			return
		}

		// Add data to PDF using the HTML template
		var tplOutput strings.Builder
		err = tmpl.Execute(&tplOutput, data)
		if err != nil {
			helpers.MessageLogs.ErrorLog.Printf(
				"controllers/GetPDF Error executing HTML template err: %+v",
				err)

			er := errors.New("Error executing HTML template")
			helpers.ErrorJSON(w, er, http.StatusBadRequest)
			return
		}
		page := wkhtmltopdf.NewPageReader(strings.NewReader(tplOutput.String()))
		pdfg.AddPage(page)
		err = pdfg.Create()
		if err != nil {
			helpers.MessageLogs.ErrorLog.Printf(
				"controllers/GetPDF Error in generating PDF err: %+v",
				err)

			er := errors.New("Error in generating PDF")
			helpers.ErrorJSON(w, er, http.StatusBadRequest)
			return
		}

		// Add PDF file to the zip archive in memory
		zipFile, err := zipWriter.Create(fmt.Sprintf("invoice_%d.pdf", i+1))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating PDF file in zip: %s", err), http.StatusInternalServerError)

		}

		pdfcpuconnfig := pdfCPUModel.NewDefaultConfiguration()
		pdfcpuconnfig.OwnerPW = data["pdf_password"]
		pdfcpuconnfig.UserPW = data["pdf_password"]
		pdfcpuconnfig.EncryptUsingAES = true
		err = api.Encrypt(bytes.NewReader(pdfg.Bytes()), zipFile, pdfcpuconnfig)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing encrypted PDF to zip: %s", err), http.StatusInternalServerError)
			return
		}
	}

	err = zipWriter.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error closing zip file: %s", err), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=invoices.zip")

	// Write the zip file to the response
	_, err = w.Write(zipFileBytes.Bytes())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing response: %s", err), http.StatusInternalServerError)
		return
	}
}

func OAuthLogin(w http.ResponseWriter, r *http.Request) {
	// TODO: add empty check for this
	redirectURI := strings.TrimSpace(r.URL.Query().Get("redirect_uri"))

	provider := chi.URLParam(r, "oauth_provider")
	ctx := r.Context()
	url, err := business.OAuthLogin(ctx, redirectURI, provider)

	if err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envolope{
			"msg": "Failed to OAuth login URL",
			"err": err,
		})
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func OAuthCallback(w http.ResponseWriter, r *http.Request) {
	// TODO: add empty check for this

	redirectURL := r.FormValue("redirect_uri")
	ctx := r.Context()
	oauthCode := r.FormValue("code")
	provider := chi.URLParam(r, "oauth_provider")

	emailID, err := business.OAuthCallback(ctx, oauthCode, provider)
	if err != nil {

		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/OAuthCallback Failed to get user email err: %+v",
			err)

		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envolope{
			"msg": "Failed to get user email",
			"err": err,
		})
		return

	}
	// set email
	tokenString, err := business.LoginUserByEmailID(emailID)

	if err != nil {

		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/OAuthCallback Failed to login err: %+v",
			err)

		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envolope{
			"msg": "Failed to parse the body of the req",
			"err": err,
		})
		return

	}
	http.SetCookie(w, &http.Cookie{
		Name:    "Authorization",
		Value:   tokenString,
		Expires: time.Now().Add(time.Hour),
	})

	http.Redirect(w, r, redirectURL, http.StatusFound)

}
