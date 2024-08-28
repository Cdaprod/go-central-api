package minio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Cdaprod/go-central-api/utils"
)

type MinioAPI struct {
	BaseURL string
}

func NewMinioAPI(baseURL string) *MinioAPI {
	return &MinioAPI{BaseURL: baseURL}
}

func (a *MinioAPI) GetName() string {
	return "minio"
}

func (a *MinioAPI) Handle(method, path string, payload []byte) ([]byte, error) {
	switch {
	case method == "GET" && path == "buckets":
		return a.listBuckets()
	case method == "POST" && path == "buckets":
		return a.createBucket(payload)
	case method == "GET" && strings.HasPrefix(path, "buckets/"):
		parts := strings.Split(path, "/")
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid path: %s", path)
		}
		bucketName := parts[1]
		objectName := strings.Join(parts[2:], "/")
		return a.getObject(bucketName, objectName)
	case method == "PUT" && strings.HasPrefix(path, "buckets/"):
		parts := strings.Split(path, "/")
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid path: %s", path)
		}
		bucketName := parts[1]
		objectName := strings.Join(parts[2:], "/")
		return a.putObject(bucketName, objectName, payload)
	default:
		return nil, fmt.Errorf("unknown endpoint: %s %s", method, path)
	}
}

func (a *MinioAPI) listBuckets() ([]byte, error) {
	resp, err := http.Get(a.BaseURL + "/buckets")
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %v", err)
	}
	defer resp.Body.Close()

	return utils.ReadAndValidateJSON(resp.Body)
}

func (a *MinioAPI) createBucket(payload []byte) ([]byte, error) {
	resp, err := http.Post(a.BaseURL+"/buckets", "application/json", strings.NewReader(string(payload)))
	if err != nil {
		return nil, fmt.Errorf("failed to create bucket: %v", err)
	}
	defer resp.Body.Close()

	return utils.ReadAndValidateJSON(resp.Body)
}

func (a *MinioAPI) getObject(bucketName, objectName string) ([]byte, error) {
	resp, err := http.Get(a.BaseURL + "/buckets/" + bucketName + "/" + objectName)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %v", err)
	}
	defer resp.Body.Close()

	return utils.ReadAll(resp.Body)
}

func (a *MinioAPI) putObject(bucketName, objectName string, payload []byte) ([]byte, error) {
	req, err := http.NewRequest("PUT", a.BaseURL+"/buckets/"+bucketName+"/"+objectName, strings.NewReader(string(payload)))
	if err != nil {
		return nil, fmt.Errorf("failed to create put request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to put object: %v", err)
	}
	defer resp.Body.Close()

	return utils.ReadAndValidateJSON(resp.Body)
}