package business

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/akashc777/csvToPdf/initializers/oauth"
	authModels "github.com/akashc777/csvToPdf/models/postgres/OAuth"
	models "github.com/akashc777/csvToPdf/models/postgres/Templates"
	"github.com/coreos/go-oidc/v3/oidc"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	domain "github.com/akashc777/csvToPdf/domain/Templates"
	"github.com/akashc777/csvToPdf/helpers"
	jwt "github.com/golang-jwt/jwt/v5"
)

const AuthRecipeMethodGoogle = "google"
const AuthRecipeMethodGithub = "github"
const GithubUserInfoURL = "https://api.github.com/user"
const GithubUserEmailsURL = "https://api.github.com/user/emails"

type RawFileInfo struct {
	FileName string `json:"file_name"`
}

func CreateTemplate(createdByEmailId, htmlTemplate, templateName string) (err error) {

	err = domain.CreateTemplate(createdByEmailId, htmlTemplate, templateName)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"business/CreateTemplate Failed to create new template err: %+v",
			err)
		return
	}

	return
}

func GetTemplateByName(emailID *string, templateName *string) (templateInfo *models.Template, err error) {
	templateInfo, err = domain.GetTemplateByName(emailID, templateName)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"business/GetTemplateByName Failed to get template err: %+v",
			err)
		return
	}

	return
}

func GetTemplateNamesByEmailID(emailID *string) (templateNames []*string, err error) {
	templateNames, err = domain.GetTemplateNamesByEmailID(emailID)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"business/GetTemplateNamesByEmailID Failed to get template err: %+v",
			err)
		return
	}
	return
}

func UpdateTemplateByName(templateInfo *models.Template) (err error) {
	err = domain.UpdateTemplateByName(templateInfo.TemplateName, templateInfo.CreatedBy, templateInfo.HtmlTemplate, templateInfo.NewTemplateName)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"business/UpdateTemplateByName Failed to update template err: %+v",
			err)
		return
	}
	return
}

func DeleteTemplateByName(emailID *string, templateName *string) (err error) {
	err = domain.DeleteTemplateByName(emailID, templateName)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"business/DeleteTemplateByName Failed to delete template err: %+v",
			err)
		return
	}
	return
}

func LoginUserByEmailID(emailID string) (tokenString string, err error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": emailID,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})

	jwtSeceret := os.Getenv("JWT_SECRET")
	tokenString, err = token.SignedString([]byte(jwtSeceret))

	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"business/LoginUserByEmailID Failed to create token err: %+v",
			err,
		)
		err = errors.New("failed to create token")
		return
	}

	return

}

func OAuthLogin(ctx context.Context, redirectURI string, provider string) (url string, err error) {
	switch provider {
	case AuthRecipeMethodGoogle:
		// during the init of OAuthProvider authorizer url might be empty
		url = oauth.OAuthProviders.GoogleConfig.AuthCodeURL(redirectURI)

	case AuthRecipeMethodGithub:
		url = oauth.OAuthProviders.GithubConfig.AuthCodeURL(redirectURI)

	default:
		helpers.MessageLogs.ErrorLog.Printf(
			"controllers/OAuthCallback Unknown provider : %+v",
			provider)
		err = errors.New("unknown provider")
	}

	return
}

func OAuthCallback(ctx context.Context, oauthCode string, provider string) (emailId string, err error) {
	var user *authModels.User
	switch provider {
	case AuthRecipeMethodGoogle:
		user, err = processGoogleUserInfo(ctx, oauthCode)
	case AuthRecipeMethodGithub:
		user, err = processGithubUserInfo(ctx, oauthCode)
	}

	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"business/OAuthCallback Failed to process user info err: %+v",
			err)

		return
	}

	emailId = *user.Email

	if len(emailId) == 0 {

		helpers.MessageLogs.ErrorLog.Printf(
			"business/OAuthCallback Failed to get user email")

		return

	}

	return
}

func processGoogleUserInfo(ctx context.Context, code string) (*authModels.User, error) {
	oauth2Token, err := oauth.OAuthProviders.GoogleConfig.Exchange(ctx, code)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"business/processGoogleUserInfo Failed to exchange code for token err: %+v",
			err)
		return nil, fmt.Errorf("invalid google exchange code: %s", err.Error())
	}
	verifier := oauth.OIDCProviders.GoogleOIDC.Verifier(&oidc.Config{ClientID: oauth.OAuthProviders.GoogleConfig.ClientID})

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		helpers.MessageLogs.ErrorLog.Printf(
			"business/processGoogleUserInfo Failed to extract ID Token from OAuth2 token err: %+v",
			err)
		return nil, fmt.Errorf("unable to extract id_token")
	}

	// Parse and verify ID Token payload.
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"business/processGoogleUserInfo Failed to verify ID Token err: %+v",
			err)
		return nil, fmt.Errorf("unable to verify id_token: %s", err.Error())
	}
	user := &authModels.User{}
	if err := idToken.Claims(&user); err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"business/processGoogleUserInfo Failed to parse ID Token claims err: %+v",
			err)
		return nil, fmt.Errorf("unable to extract claims")
	}

	return user, nil
}

func processGithubUserInfo(ctx context.Context, code string) (*authModels.User, error) {
	oauth2Token, err := oauth.OAuthProviders.GithubConfig.Exchange(ctx, code)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"business/processGithubUserInfo Failed to exchange code for token err: %+v",
			err)
		return nil, fmt.Errorf("invalid github exchange code: %s", err.Error())
	}
	client := http.Client{}
	req, err := http.NewRequest("GET", GithubUserInfoURL, nil)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"business/processGithubUserInfo Failed to create github user info request err: %+v",
			err)
		return nil, fmt.Errorf("error creating github user info request: %s", err.Error())
	}
	req.Header.Set(
		"Authorization", fmt.Sprintf("token %s", oauth2Token.AccessToken),
	)

	response, err := client.Do(req)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"business/processGithubUserInfo Failed to request github user info err: %+v",
			err)
		return nil, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"business/processGithubUserInfo Failed to read github user info response body err: %+v",
			err)
		return nil, fmt.Errorf("failed to read github response body: %s", err.Error())
	}
	if response.StatusCode >= 400 {
		helpers.MessageLogs.ErrorLog.Printf(
			"business/processGithubUserInfo Failed to request github user info err: %+v",
			err)
		return nil, fmt.Errorf("failed to request github user info: %s", string(body))
	}

	userRawData := make(map[string]string)
	err = json.Unmarshal(body, &userRawData)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Printf(
			"business/processGithubUserInfo Failed to unmarshal uaer raw data err: %+v",
			err)
		return nil, fmt.Errorf("failed to read github response body: %s", err.Error())
	}

	name := strings.Split(userRawData["name"], " ")
	firstName := ""
	lastName := ""
	if len(name) >= 1 && strings.TrimSpace(name[0]) != "" {
		firstName = name[0]
	}
	if len(name) > 1 && strings.TrimSpace(name[1]) != "" {
		lastName = name[0]
	}

	picture := userRawData["avatar_url"]
	email := userRawData["email"]

	if email == "" {
		type GithubUserEmails struct {
			Email   string `json:"email"`
			Primary bool   `json:"primary"`
		}

		// fetch using /users/email endpoint
		req, err := http.NewRequest(http.MethodGet, GithubUserEmailsURL, nil)
		if err != nil {
			helpers.MessageLogs.ErrorLog.Printf(
				"business/processGithubUserInfo Failed to create github emails request err: %+v",
				err)
			return nil, fmt.Errorf("error creating github user info request: %s", err.Error())
		}
		req.Header.Set(
			"Authorization", fmt.Sprintf("token %s", oauth2Token.AccessToken),
		)

		response, err := client.Do(req)
		if err != nil {
			helpers.MessageLogs.ErrorLog.Printf(
				"business/processGithubUserInfo Failed to request github user email err: %+v",
				err)
			return nil, err
		}

		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			helpers.MessageLogs.ErrorLog.Printf(
				"business/processGithubUserInfo Failed to read github user email response body err: %+v",
				err)
			return nil, fmt.Errorf("failed to read github response body: %s", err.Error())
		}
		if response.StatusCode >= 400 {
			helpers.MessageLogs.ErrorLog.Printf(
				"business/processGithubUserInfo UserInfo Failed to request github user email err: %+v",
				err)
			return nil, fmt.Errorf("failed to request github user info: %s", string(body))
		}

		emailData := []GithubUserEmails{}
		err = json.Unmarshal(body, &emailData)
		if err != nil {
			helpers.MessageLogs.ErrorLog.Printf(
				"business/processGithubUserInfo Failed to parse github user email err: %+v",
				err)
			return nil, fmt.Errorf("failed to parse github user email: %s", err.Error())
		}

		for _, userEmail := range emailData {
			email = userEmail.Email
			if userEmail.Primary {
				break
			}
		}
	}

	user := &authModels.User{
		GivenName:  &firstName,
		FamilyName: &lastName,
		Picture:    &picture,
		Email:      &email,
	}

	return user, nil
}
