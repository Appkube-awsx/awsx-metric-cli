package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type AwsCredential struct {
	Region              string `json:"region,omitempty"`
	AccessKey           string `json:"accessKey,omitempty"`
	SecretKey           string `json:"secretKey,omitempty"`
	CrossAccountRoleArn string `json:"crossAccountRoleArn,omitempty"`
	ExternalId          string `json:"externalId,omitempty"`
}

type VaultResponse struct {
	RequestId     string        `json:"request_id,omitempty"`
	LeaseId       string        `json:"lease_id,omitempty"`
	Renewable     string        `json:"renewable,omitempty"`
	LeaseDuration int64         `json:"lease_duration,omitempty"`
	Data          AwsCredential `json:"data,omitempty"`
}

type ApiResponse struct {
	Status     string        `json:"status,omitempty"`
	Message    string        `json:"message,omitempty"`
	StatusCode int64         `json:"statusCode,omitempty"`
	Data       AwsCredential `json:"data,omitempty"`
}

func GetAccountDetails(vaultUrl string, vaultToken string, accountNo string) (*VaultResponse, error) {
	log.Println("Calling account details API")
	client := &http.Client{}
	req, err := http.NewRequest("GET", vaultUrl+"/"+accountNo, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Vault-Token", vaultToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	var responseObject VaultResponse
	json.Unmarshal(bodyBytes, &responseObject)
	return &responseObject, nil

}

func GetUserCredential(cloudElementId string, apiUrl string) (apiResponse *ApiResponse, statusCode int, err error) {
	log.Println("Getting user's cloud credentials by cloud-element-id")
	url := apiUrl + "?cloudElementId=" + cloudElementId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, http.StatusInternalServerError, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil && res != nil {
		fmt.Println("error in getting response from server", "url", url, "method", req.Method, "error", err.Error(), "status code", res.StatusCode)
		return nil, res.StatusCode, fmt.Errorf("error in getting response from server. api url %s", url)
	}
	if err != nil && res == nil {
		fmt.Println("error getting response from server. no response received", "url", url, "error", err.Error())
		return nil, http.StatusInternalServerError, fmt.Errorf("error getting response from server. no response received. api url %s. no response received. Error: %s", url, err.Error())
	}
	if err == nil && res == nil {
		fmt.Println("invalid response from server and also no error", "url", url, "method", req.Method)
		return nil, http.StatusInternalServerError, fmt.Errorf("invalid response from server and also no error. api url %s", url)
	}
	if res.StatusCode >= http.StatusBadRequest {
		return nil, res.StatusCode, fmt.Errorf(res.Status)
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading response body", "url", url, "error", err.Error())
		return nil, res.StatusCode, err
	}
	bodyBytes = removeBOMContent(bodyBytes)
	var responseObject ApiResponse
	jsonMarshalError := json.Unmarshal(bodyBytes, &responseObject)
	if jsonMarshalError != nil {
		fmt.Println("json marshal error", "url", url, "error", jsonMarshalError.Error())
		return nil, http.StatusInternalServerError, jsonMarshalError
	}
	return &responseObject, http.StatusOK, nil
}

func removeBOMContent(input []byte) []byte {
	return bytes.TrimPrefix(input, []byte("\xef\xbb\xbf"))
}

func GetAccountDetailsForENV(vaultUrl string, vaultToken string, accountNo string) (*VaultResponse, error) {
	log.Println("Calling account details API")
	client := &http.Client{}
	req, err := http.NewRequest("GET", vaultUrl+"/"+accountNo, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Vault-Token", vaultToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	var responseObject VaultResponse
	json.Unmarshal(bodyBytes, &responseObject)
	return &responseObject, nil

}
