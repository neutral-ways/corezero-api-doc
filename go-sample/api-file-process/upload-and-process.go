package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

type Config struct {
	ApiHost  string `json:"api_host"`
	ApiPort  string `json:"api_port"`
	ApiKey   string `json:"api_key"`
	ApiProto string `json:"api_proto"`
}

type CreateUploadRequest struct {
	ContentType string `json:"content_type"`
	Filename    string `json:"filename"`
}

type CreateUploadResponse struct {
	Data struct {
		ID           string `json:"id"`
		Entity       string `json:"entity"`
		EntityID     string `json:"entity_id"`
		Filename     string `json:"filename"`
		PreSignedURL string `json:"pre_signed_url"`
		Headers      struct {
			XAmzMetaContentType string `json:"x-amz-meta-content-type"`
			XAmzMetaEntity      string `json:"x-amz-meta-entity"`
			XAmzMetaEntityID    string `json:"x-amz-meta-entity-id"`
			XAmzMetaFilename    string `json:"x-amz-meta-filename"`
			XAmzMetaPublic      string `json:"x-amz-meta-public"`
			XAmzMetaUploader    string `json:"x-amz-meta-uploader"`
		} `json:"headers"`
	} `json:"data"`
}

type ProcessFileRequest struct {
	AttachmentID string `json:"attachment_id"`
}

type ProcessFileResponse struct {
	Data struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		//DeletedAt    interface{} `json:"deleted_at"`
		CreatedBy string `json:"created_by"`
		UpdatedBy string `json:"updated_by"`
		//		AccountID    interface{} `json:"account_id"`
		//		ProjectID    interface{} `json:"project_id"`
		//		LotID        interface{} `json:"lot_id"`
		AttachmentID string `json:"attachment_id"`
		Status       string `json:"status"`
		//		ErrorMessage interface{} `json:"error_message"`
		EndedAt   *time.Time `json:"ended_at"`
		Operation string     `json:"operation"`
	} `json:"data"`
}

type GetWorkerResponse struct {
	Data struct {
		ID           string      `json:"id"`
		CreatedAt    string      `json:"created_at"`
		UpdatedAt    string      `json:"updated_at"`
		DeletedAt    interface{} `json:"deleted_at"`
		CreatedBy    string      `json:"created_by"`
		UpdatedBy    string      `json:"updated_by"`
		AccountID    string      `json:"account_id"`
		ProjectID    string      `json:"project_id"`
		LotID        string      `json:"lot_id"`
		AttachmentID string      `json:"attachment_id"`
		Status       string      `json:"status"`
		ErrorMessage string      `json:"error_message"`
		EndedAt      string      `json:"ended_at"`
		Operation    string      `json:"operation"`
		Proccesed    int         `json:"proccesed"`
		Total        int         `json:"total"`
	} `json:"data"`
}

func main() {

	fmt.Println("Upload-and-process / CoreZero (c) 2022 \U0001F9DF")

	var filename string

	monitorPtr := flag.Bool("monitor", false, "monitor worker results")
	flag.StringVar(&filename, "file", "", "file to process (mandatory)")
	flag.Parse()

	if filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	b, _ := exists(filename)
	if !b {
		fmt.Println("file not found")
		os.Exit(1)
	}

	config := loadConfig()

	if monitorPtr != nil && *monitorPtr {
		fmt.Println("Monitor is enabled!")
	}

	fmt.Println("API file upload started")
	fmt.Printf(" - API-KEY is: %s\n", config.ApiKey)
	fmt.Printf(" - API host : %s\n\n", config.ApiHost)

	fmt.Println("step 1: create upload request")
	t, err := uploadRequest(config, filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(" - attachment_id: " + t.Data.ID)
	fmt.Println(" - entity_id: " + t.Data.EntityID)

	fmt.Println("step 2: upload file")

	err = uploadFile(config, filename, t)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Print(" - upload done. waiting 5 secs")
	for i := 0; i < 5; i++ {
		fmt.Print(".")
		time.Sleep(1 * time.Second)
	}
	fmt.Print("\n")

	fmt.Println("step 3: process file")
	r, err := processFile(config, t.Data.ID)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if monitorPtr != nil && *monitorPtr {
		var bar *progressbar.ProgressBar = nil

		fmt.Println("step 4: monitor worker")
		running := true
		for running {
			tw, err := getWorkerStatus(config, r.Data.ID)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			switch tw.Data.Status {
			case "finished":
				running = false
				if bar != nil {
					bar.Set(bar.GetMax())
				}
				break
			case "running":
				if bar == nil {
					bar = progressbar.Default(int64(tw.Data.Total))
				}
				bar.Set(tw.Data.Proccesed)
				break
			case "init":
			case "created":
				fmt.Println("waiting for worker to start")
				break
			default:
				break
			}

			time.Sleep(1 * time.Second)
		}

		fmt.Println("end")

	}

}

func loadConfig() Config {

	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("cannot load config")
		os.Exit(1)
	}

	config := Config{}
	err = json.Unmarshal([]byte(file), &config)
	if err != nil {
		fmt.Println("cannot parse config")
		os.Exit(1)
	}

	return config
}

func exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func uploadRequest(config Config, filename string) (CreateUploadResponse, error) {

	data := CreateUploadRequest{
		ContentType: "text/csv",
		Filename:    filename,
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return CreateUploadResponse{}, err
	}

	body := bytes.NewReader(payloadBytes)

	url := fmt.Sprintf("%s://%s:%s/api/v1/client/file", config.ApiProto, config.ApiHost, config.ApiPort)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return CreateUploadResponse{}, err
	}

	req.Header.Set("X-API-KEY", config.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return CreateUploadResponse{}, err
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return CreateUploadResponse{}, err
	}

	t := CreateUploadResponse{}
	err = json.Unmarshal(bytes, &t)
	if err != nil {
		return t, err
	}

	return t, nil

}

func uploadFile(config Config, filename string, uploadRespose CreateUploadResponse) error {

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	payload := strings.NewReader(string(b))

	req, err := http.NewRequest("PUT", uploadRespose.Data.PreSignedURL, payload)
	if err != nil {
		return err
	}

	req.Header.Set("X-Amz-Meta-Entity", "tx-processor")
	req.Header.Set("X-Amz-Meta-Entity-Id", uploadRespose.Data.Headers.XAmzMetaEntityID)
	req.Header.Set("X-Amz-Meta-Filename", uploadRespose.Data.Headers.XAmzMetaFilename)
	req.Header.Set("X-Amz-Meta-Uploader", uploadRespose.Data.Headers.XAmzMetaUploader)
	req.Header.Set("X-Amz-Meta-Content-Type", "text/csv")
	req.Header.Set("X-Amz-Meta-Public", "false")
	req.Header.Set("Content-Type", "text/csv")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func processFile(config Config, attachmentID string) (ProcessFileResponse, error) {

	data := ProcessFileRequest{
		AttachmentID: attachmentID,
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return ProcessFileResponse{}, err
	}

	body := bytes.NewReader(payloadBytes)

	url := fmt.Sprintf("%s://%s:%s/api/v1/client/process", config.ApiProto, config.ApiHost, config.ApiPort)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return ProcessFileResponse{}, err
	}

	req.Header.Set("X-API-KEY", config.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ProcessFileResponse{}, err
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ProcessFileResponse{}, err
	}

	t := ProcessFileResponse{}
	err = json.Unmarshal(bytes, &t)
	if err != nil {
		return t, err
	}

	return t, nil

}

func getWorkerStatus(config Config, txWorkerID string) (GetWorkerResponse, error) {

	url := fmt.Sprintf("%s://%s:%s/api/v1/client/tx-worker/%s", config.ApiProto, config.ApiHost, config.ApiPort, txWorkerID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return GetWorkerResponse{}, err
	}

	req.Header.Set("X-API-KEY", config.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return GetWorkerResponse{}, err
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetWorkerResponse{}, err
	}

	t := GetWorkerResponse{}
	err = json.Unmarshal(bytes, &t)
	if err != nil {
		return t, err
	}

	return t, nil

}
