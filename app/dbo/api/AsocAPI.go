package api

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sast-integration/app/dto/request"
	"sast-integration/app/helper"
	"sast-integration/app/shareVar"
	"strings"
)

type AsocAPI struct {
	api_asoc_service string
	super_token      string
	stdLog           *helper.StandartLog
}

func NewAsocAPI() *AsocAPI {
	if os.Getenv("APP_ENV") == "" {
		errEnv := godotenv.Load(".env")
		helper.ErrorHandler(errEnv)
	}

	asoc_url := os.Getenv("ASOC_SERVICE_URL")

	return &AsocAPI{
		api_asoc_service: asoc_url,
		super_token:      os.Getenv("SUPER_TOKEN"),
		stdLog:           helper.NewStandardLog(shareVar.AsocAPI, shareVar.Service),
	}
}

func (a AsocAPI) CreateProject(request request.AsocProjectCreateRequest) {
	a.stdLog.NameFunc = "CreateProject"
	a.stdLog.StartFunction(request)

	byteReq, err := json.Marshal(request)
	helper.ErrorHandler(err)
	reqBody := strings.NewReader(string(byteReq))

	httpReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/sast/project", a.api_asoc_service), reqBody)
	helper.ErrorHandler(err)
	httpReq.Header.Add("content-type", "application/json")
	httpReq.Header.Add("Authorization", helper.MapTokenAuthorizationHeader(a.super_token))

	client := &http.Client{}
	response, err := client.Do(httpReq)
	helper.ErrorHandler(err)

	bodyByte, err := ioutil.ReadAll(response.Body)
	helper.ErrorHandler(err)

	log.Printf("%+v", string(bodyByte))

	if response.StatusCode == http.StatusNotFound {
		a.stdLog.NameFunc = "CreateProject"
		a.stdLog.WarningFunction(response)
	} else if response.StatusCode == http.StatusBadRequest {
		a.stdLog.NameFunc = "CreateProject"
		a.stdLog.WarningFunction(response)
	} else if response.StatusCode == http.StatusConflict {
		a.stdLog.NameFunc = "CreateProject"
		a.stdLog.WarningFunction(response)
	}

	a.stdLog.NameFunc = "CreateOne"
	a.stdLog.EndFunction(nil)

}

func (a AsocAPI) SendResult(request request.AsocSendResultRequest) {
	a.stdLog.NameFunc = "SendResult"
	a.stdLog.StartFunction(request)

	byteReq, err := json.Marshal(request)
	helper.ErrorHandler(err)
	reqBody := strings.NewReader(string(byteReq))

	httpReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/sast/receive-result", a.api_asoc_service), reqBody)
	helper.ErrorHandler(err)
	httpReq.Header.Add("content-type", "application/json")
	httpReq.Header.Add("Authorization", helper.MapTokenAuthorizationHeader(a.super_token))

	client := &http.Client{}
	response, err := client.Do(httpReq)
	helper.ErrorHandler(err)

	bodyByte, err := ioutil.ReadAll(response.Body)
	helper.ErrorHandler(err)

	log.Printf("%+v", string(bodyByte))

	if response.StatusCode == http.StatusNotFound {
		a.stdLog.NameFunc = "SendResult"
		a.stdLog.WarningFunction(response)
	} else if response.StatusCode == http.StatusBadRequest {
		a.stdLog.NameFunc = "SendResult"
		a.stdLog.WarningFunction(response)
	} else if response.StatusCode == http.StatusConflict {
		a.stdLog.NameFunc = "SendResult"
		a.stdLog.WarningFunction(response)
	}

	a.stdLog.NameFunc = "SendResult"
	a.stdLog.EndFunction(nil)

}
