package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type FileData interface {
	GetFileData() File
	GetFormValues() url.Values
}

type File struct {
	Reader   io.Reader
	FileName string
	FileKey  string
}

type HttpSuccessBasicResponse struct {
	Success bool `json:"success"`
}

type HttpErrorResponse struct {
	Success bool               `json:"success,omitempty"`
	Message string             `json:"message,omitempty"`
	Errors  *map[string]string `json:"errors,omitempty"`
}

type SyncHangupResponse struct {
	Success bool    `json:"success"`
	Code    int     `json:"code"`
	Message *string `json:"message"`
}

type HttpConfig struct {
	BaseUrl string
	AppId   string
}

func InitHttpConfig(baseUrl string, appId string) *HttpConfig {
	return &HttpConfig{BaseUrl: baseUrl, AppId: appId}
}

func Get[T any](config HttpConfig, path string, resultType T) (*T, *HttpErrorResponse) {
	url := getUrl(config, path)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	return HttpRequest[T](request, url, nil, resultType)
}

func Post[T any](config HttpConfig, path string, payload interface{}, resultType T) (*T, *HttpErrorResponse) {
	url := getUrl(config, path)
	jsonData, _ := json.Marshal(payload)
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	return HttpRequest[T](request, url, payload, resultType)
}

func Put[T any](config HttpConfig, path string, payload interface{}, resultType T) (*T, *HttpErrorResponse) {
	url := getUrl(config, path)
	jsonData, _ := json.Marshal(payload)
	request, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	return HttpRequest[T](request, url, payload, resultType)
}

func Patch[T any](config HttpConfig, path string, payload interface{}, resultType T) (*T, *HttpErrorResponse) {
	url := getUrl(config, path)
	jsonData, _ := json.Marshal(payload)
	request, _ := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonData))
	return HttpRequest[T](request, url, payload, resultType)
}

func Delete[T any](config HttpConfig, path string, resultType T) (*T, *HttpErrorResponse) {
	url := getUrl(config, path)
	request, _ := http.NewRequest(http.MethodDelete, url, nil)
	return HttpRequest[T](request, url, nil, resultType)
}

func Download(config HttpConfig, path string) ([]byte, *HttpErrorResponse) {
	url := getUrl(config, path)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	return downloadFile(request)
}

func Upload(config HttpConfig, path string, data FileData) (*HttpSuccessBasicResponse, *HttpErrorResponse) {
	url := getUrl(config, path)
	request, _ := http.NewRequest(http.MethodPost, url, nil)
	return uploadFile(request, data)
}

func HttpRequest[T any](request *http.Request, url string, payload interface{}, successResponse T) (*T, *HttpErrorResponse) {
	var errorResponse = HttpErrorResponse{Success: false, Message: "Internal server error"}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)

	if err != nil {
		return nil, &errorResponse
	}

	if response.StatusCode == 204 {
		return &successResponse, nil
	}

	defer response.Body.Close()

	var result interface{}
	var body, _ = io.ReadAll(response.Body)

	if len(body) == 0 && response.StatusCode == 200 {
		return &successResponse, nil
	}

	err = json.Unmarshal([]byte(body), &result)

	if object, ok := result.(map[string]interface{}); ok {
		if object["error"] == true {
			return nil, getErrorDetails(object)
		} else if object["success"] == false {
			return nil, getErrorDetails(object)
		}
	}

	if err != nil {
		return nil, &errorResponse
	}

	response.Body = io.NopCloser(bytes.NewReader(body))

	err = json.NewDecoder(response.Body).Decode(&successResponse)
	if err != nil {
		return nil, &errorResponse
	}

	return &successResponse, nil
}

func downloadFile(request *http.Request) ([]byte, *HttpErrorResponse) {
	var errorResponse = HttpErrorResponse{Success: false, Message: "No file was downloaded"}
	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)

	if err != nil {
		return nil, &errorResponse
	}

	defer response.Body.Close()

	contentDisposition := response.Header.Get("Content-Disposition")

	if !strings.Contains(contentDisposition, "attachment") {
		return nil, &errorResponse
	}

	fileData, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, &errorResponse
	}

	return fileData, nil
}

func uploadFile(request *http.Request, data FileData) (*HttpSuccessBasicResponse, *HttpErrorResponse) {
	requestBodyBuffer := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBodyBuffer)

	fileData := data.GetFileData()
	formData := data.GetFormValues()

	fileWriter, err := writer.CreateFormFile(fileData.FileKey, fileData.FileName)

	if err != nil {
		return nil, &HttpErrorResponse{Success: false, Message: "Failed to create form file"}
	}

	_, err = io.Copy(fileWriter, fileData.Reader)

	if err != nil {
		return nil, &HttpErrorResponse{Success: false, Message: "Failed to copy file data"}
	}

	for key, values := range formData {
		for _, value := range values {
			_ = writer.WriteField(key, value)
		}
	}

	err = writer.Close()

	if err != nil {
		return nil, &HttpErrorResponse{Success: false, Message: "Failed to close writer"}
	}

	bodyBytes, _ := io.ReadAll(requestBodyBuffer)

	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)

	if err != nil {
		return nil, &HttpErrorResponse{Success: false, Message: err.Error()}
	}

	if response.StatusCode != 200 {
		var responseResult map[string]interface{}

		responseBody, _ := io.ReadAll(response.Body)
		err = json.Unmarshal(responseBody, &responseResult)

		if err != nil {
			return nil, &HttpErrorResponse{Success: false, Message: err.Error()}
		}

		if responseResult["error"] == true {
			return nil, getErrorDetails(responseResult)
		}

		return nil, &HttpErrorResponse{Success: false, Message: fmt.Sprintf("Unknown error with status %v", response.Status)}
	}

	defer response.Body.Close()
	defer request.Body.Close()
	defer writer.Close()

	return &HttpSuccessBasicResponse{Success: true}, nil
}

func getUrl(config HttpConfig, url string) string {
	path := config.BaseUrl + url
	if strings.Contains(path, "?") {
		return path + "&appid=" + config.AppId
	}

	return path + "?appid=" + config.AppId
}

func getErrorDetails(obj map[string]interface{}) *HttpErrorResponse {
	if _, ok := obj["message"]; ok {
		return &HttpErrorResponse{Success: false, Message: obj["message"].(string)}
	}

	return &HttpErrorResponse{Success: false, Message: "Unknown error"}
}
