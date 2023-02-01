package src

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Rayato159/rayato-go-facebook-oauth/models"
	"github.com/Rayato159/rayato-go-facebook-oauth/utils"
)

type errRes struct {
	Error *errDatum `json:"error"`
}

type errDatum struct {
	Message   string `json:"message"`
	Type      string `json:"type"`
	Code      int    `json:"code"`
	FbtraceID string `json:"fbtrace_id"`
}

type facebookEnv struct {
	GraphUrl     string
	Version      string
	CallbackUrl  string
	ClientId     string
	ClientSecret string
}

// Fields of public_profile -> https://developers.facebook.com/docs/graph-api/reference/user
type UserFields string

const (
	Id      UserFields = "string"
	Name    UserFields = "name"
	Email   UserFields = "email" // email permssion required
	Picture UserFields = "picture.width(240).height(240){url}"
)

// Permission
type UserPermission string

const (
	Birthday UserPermission = "user_birthday"
	Gender   UserPermission = "user_gender"
)

func NewGoFacebookOauth(version, callbackUrl, clientId, clientSecret string) *facebookEnv {
	// Set default
	if version == "" {
		version = "15.0"
	}

	return &facebookEnv{
		GraphUrl:     "https://graph.facebook.com",
		Version:      fmt.Sprintf("v%s", strings.TrimPrefix(version, "v")),
		CallbackUrl:  callbackUrl,
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
}

func fieldsConcator(fields []UserFields) (string, error) {
	var result string
	if len(fields) == 0 {
		fields = []UserFields{
			"id",
			"name",
			"email",
			"picture",
		}
	}
	for i := range fields {
		if string(fields[i]) == "" {
			return "", fmt.Errorf("fields is invalid")
		}

		if i != len(fields)-1 {
			result += fmt.Sprintf("%s,", string(fields[i]))
		} else {
			result += string(fields[i])
		}
	}
	return result, nil
}

func permissionConcator(str *string, permissions []UserPermission) error {
	var result string
	for i := range permissions {
		if string(permissions[i]) == "" {
			return fmt.Errorf("permissions is invalid")
		}

		if i != len(permissions) {
			result += fmt.Sprintf("%s,", string(permissions[i]))
		} else {
			result += string(permissions[i])
		}
	}
	return nil
}

func (f *facebookEnv) GetCallbackUrl(state string, permissions ...UserPermission) (string, error) {
	if state == "" {
		state = "none"
	}

	permissionQuery := "public_profile,email"
	if len(permissions) != 0 {
		if err := permissionConcator(&permissionQuery, permissions); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf(
		"https://www.facebook.com/%s/dialog/oauth?client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		f.Version,
		f.ClientId,
		f.CallbackUrl,
		permissionQuery,
		state,
	), nil
}

func (f *facebookEnv) GetAccessToken(code string) (*models.UserAccessToken, error) {
	// Config a request before fire an API
	url := fmt.Sprintf(
		"%s/%s/oauth/access_token?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s",
		f.GraphUrl,
		f.Version,
		f.ClientId,
		f.CallbackUrl,
		f.ClientSecret,
		code,
	)
	resJson, err := utils.FireHttpRequest("GET", url)
	if err != nil {
		return nil, fmt.Errorf("http get error: %v", err)
	}

	// Success response
	resSuccess := new(models.UserAccessToken)
	if err := json.Unmarshal(resJson, &resSuccess); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	// Error response
	resErr := &errRes{
		Error: &errDatum{},
	}
	if err := json.Unmarshal(resJson, &resErr); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}
	if resErr.Error.FbtraceID != "" {
		return nil, fmt.Errorf(resErr.Error.Message)
	}

	return resSuccess, nil
}

func (f *facebookEnv) GetUserData(userAccessToken string, fields ...UserFields) (*models.UserProfile, error) {
	// Url builder
	fieldsQuery, err := fieldsConcator(fields)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(
		"%s/%s/me?fields=%s&access_token=%s",
		f.GraphUrl,
		f.Version,
		fieldsQuery,
		userAccessToken,
	)

	// Config a request before fire an API
	resJson, err := utils.FireHttpRequest("GET", url)
	if err != nil {
		return nil, fmt.Errorf("http get error: %v", err)
	}

	// Success response
	resSuccess := new(models.UserProfile)
	if err := json.Unmarshal(resJson, &resSuccess); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	// Error response
	resErr := &errRes{
		Error: &errDatum{},
	}
	if err := json.Unmarshal(resJson, &resErr); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}
	if resErr.Error.FbtraceID != "" {
		return nil, fmt.Errorf(resErr.Error.Message)
	}

	return resSuccess, nil
}

func (f *facebookEnv) Logout(userAccessToken string) (*models.UserLogoutRes, error) {
	// Config a request before fire an API
	url := fmt.Sprintf(
		"%s/%s/me/permissions?&access_token=%s",
		f.GraphUrl,
		f.Version,
		userAccessToken,
	)
	resJson, err := utils.FireHttpRequest("DELETE", url)
	if err != nil {
		return nil, fmt.Errorf("http get error: %v", err)
	}

	// Success response
	resSuccess := new(models.UserLogoutRes)
	if err := json.Unmarshal(resJson, &resSuccess); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	// Error response
	resErr := &errRes{
		Error: &errDatum{},
	}
	if err := json.Unmarshal(resJson, &resErr); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}
	if resErr.Error.FbtraceID != "" {
		return nil, fmt.Errorf(resErr.Error.Message)
	}
	if !resSuccess.Success {
		return nil, fmt.Errorf("logout error")
	}

	return resSuccess, nil
}
