//
// Copyright (c) 2017 Tencent
//
// SPDX-License-Identifier: Apache-2.0
//

package scheduler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/edgexfoundry/edgex-go/core/domain/models"
	"io"
	"io/ioutil"
	"net/http"
)

// Common http const
const (
	ContentType        = "Content-Type"
	ContentTypeJsonVal = "application/json"
)

const (
	ScheduleApiPath      = "/api/v1/schedule"
	ScheduleEventApiPath = "/api/v1/scheduleevent"
	UrlPattern           = "http://%s:%d%s"
)

type SchedulerClient interface {
	QuerySchedule(id string) (models.Schedule, error)
	QueryScheduleWithName(scheduleName string) (models.Schedule, error)
	AddSchedule(schedule models.Schedule) error
	UpdateSchedule(schedule models.Schedule) error
	RemoveSchedule(id string) error
	QueryScheduleEvent(id string) (models.ScheduleEvent, error)
	AddScheduleEvent(scheduleEvent models.ScheduleEvent) error
	UpdateScheduleEvent(scheduleEvent models.ScheduleEvent) error
	RemoveScheduleEvent(id string) error
}

// Struct to represent the scheduler client
type SchedulerRestClient struct {
	SchedulerServiceHost string
	SchedulerServicePort int
	OwningService        string
}

func NewSchedulerRestClient(ssh string, ssp int, os string) SchedulerClient {
	return &SchedulerRestClient{
		SchedulerServiceHost: ssh,
		SchedulerServicePort: ssp,
		OwningService:        os,
	}
}

// Function to get a schedule from the remote scheduler server
func (src SchedulerRestClient) QuerySchedule(id string) (models.Schedule, error) {
	client := &http.Client{}

	remoteScheduleUrl := fmt.Sprintf(UrlPattern, src.SchedulerServiceHost, src.SchedulerServicePort, ScheduleApiPath)
	remoteScheduleUrl = remoteScheduleUrl + "/" + id

	jsonBytes, err := doGet(remoteScheduleUrl, client)
	if err != nil {
		return models.Schedule{}, err
	}

	schedule := models.Schedule{}

	if err := json.Unmarshal(jsonBytes, &schedule); err != nil {
		return models.Schedule{}, err
	}

	return schedule, nil
}

// Function to get a schedule with schedule name from the remote scheduler server
func (src SchedulerRestClient) QueryScheduleWithName(scheduleName string) (models.Schedule, error) {
	client := &http.Client{}

	remoteScheduleUrl := fmt.Sprintf(UrlPattern, src.SchedulerServiceHost, src.SchedulerServicePort, ScheduleApiPath)
	remoteScheduleUrl = remoteScheduleUrl + "/name/" + scheduleName

	jsonBytes, err := doGet(remoteScheduleUrl, client)
	if err != nil {
		return models.Schedule{}, err
	}

	schedule := models.Schedule{}

	if err := json.Unmarshal(jsonBytes, &schedule); err != nil {
		return models.Schedule{}, err
	}

	return schedule, nil
}

// Function to send a schedule to the remote scheduler server
func (src SchedulerRestClient) AddSchedule(schedule models.Schedule) error {
	client := &http.Client{}

	requestBody, err := schedule.MarshalJSON()
	if err != nil {
		return err
	}

	remoteScheduleUrl := fmt.Sprintf(UrlPattern, src.SchedulerServiceHost, src.SchedulerServicePort, ScheduleApiPath)

	return doPost(remoteScheduleUrl, bytes.NewBuffer(requestBody), client)
}

// Function to update a schedule to the remote scheduler server
func (src SchedulerRestClient) UpdateSchedule(schedule models.Schedule) error {
	client := &http.Client{}

	requestBody, err := schedule.MarshalJSON()
	if err != nil {
		return err
	}

	remoteScheduleUrl := fmt.Sprintf(UrlPattern, src.SchedulerServiceHost, src.SchedulerServicePort, ScheduleApiPath)

	return doPut(remoteScheduleUrl, bytes.NewBuffer(requestBody), client)
}

// Function to remove a schedule to the remote scheduler server
func (src SchedulerRestClient) RemoveSchedule(id string) error {
	client := &http.Client{}

	remoteScheduleUrl := fmt.Sprintf(UrlPattern, src.SchedulerServiceHost, src.SchedulerServicePort, ScheduleApiPath)
	remoteScheduleUrl = remoteScheduleUrl + "/" + id

	return doDelete(remoteScheduleUrl, client)
}

// Function to get a schedule event from the remote scheduler server
func (src SchedulerRestClient) QueryScheduleEvent(id string) (models.ScheduleEvent, error) {
	client := &http.Client{}

	remoteScheduleEventUrl := fmt.Sprintf(UrlPattern, src.SchedulerServiceHost, src.SchedulerServicePort, ScheduleEventApiPath)
	remoteScheduleEventUrl = remoteScheduleEventUrl + "/" + id

	jsonBytes, err := doGet(remoteScheduleEventUrl, client)
	if err != nil {
		return models.ScheduleEvent{}, err
	}

	scheduleEvent := models.ScheduleEvent{}

	if err := json.Unmarshal(jsonBytes, &scheduleEvent); err != nil {
		return models.ScheduleEvent{}, err
	}

	return scheduleEvent, nil
}

// Function to send a schedule event to the remote scheduler server
func (src SchedulerRestClient) AddScheduleEvent(scheduleEvent models.ScheduleEvent) error {
	client := &http.Client{}

	requestBody, err := scheduleEvent.MarshalJSON()
	if err != nil {
		return err
	}

	remoteScheduleEventUrl := fmt.Sprintf(UrlPattern, src.SchedulerServiceHost, src.SchedulerServicePort, ScheduleEventApiPath)

	return doPost(remoteScheduleEventUrl, bytes.NewBuffer(requestBody), client)
}

// Function to update a schedule event to the remote scheduler server
func (src SchedulerRestClient) UpdateScheduleEvent(scheduleEvent models.ScheduleEvent) error {
	client := &http.Client{}

	requestBody, err := scheduleEvent.MarshalJSON()
	if err != nil {
		return err
	}

	remoteScheduleEventUrl := fmt.Sprintf(UrlPattern, src.SchedulerServiceHost, src.SchedulerServicePort, ScheduleEventApiPath)

	return doPut(remoteScheduleEventUrl, bytes.NewBuffer(requestBody), client)
}

// Function to remove a schedule event to the remote scheduler server
func (src SchedulerRestClient) RemoveScheduleEvent(id string) error {
	client := &http.Client{}

	remoteScheduleEventUrl := fmt.Sprintf(UrlPattern, src.SchedulerServiceHost, src.SchedulerServicePort, ScheduleEventApiPath)
	remoteScheduleEventUrl = remoteScheduleEventUrl + "/" + id

	return doDelete(remoteScheduleEventUrl, client)
}

// Function to do get request
func doGet(url string, client *http.Client) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return []byte{}, err
	}

	return sendRequestAndGetResponse(client, req)
}

// Function to do post request
func doPost(url string, binaryReqBody io.Reader, client *http.Client) error {
	req, err := http.NewRequest(http.MethodPost, url, binaryReqBody)
	req.Header.Add(ContentType, ContentTypeJsonVal)

	if err != nil {
		return err
	}

	return sendRequest(client, req)
}

// Function to do put request
func doPut(url string, binaryReqBody io.Reader, client *http.Client) error {
	req, err := http.NewRequest(http.MethodPut, url, binaryReqBody)
	req.Header.Add(ContentType, ContentTypeJsonVal)

	if err != nil {
		return err
	}

	return sendRequest(client, req)
}

// Function to do delete request
func doDelete(url string, client *http.Client) error {
	req, err := http.NewRequest(http.MethodDelete, url, nil)

	if err != nil {
		return err
	}

	return sendRequest(client, req)
}

// Function to actually make the HTTP request
func sendRequest(client *http.Client, req *http.Request) error {
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		resp.Close = true
		return nil
	} else {
		return err
	}
}

// Function to actually make the HTTP request and get the response body
func sendRequestAndGetResponse(client *http.Client, req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	resp.Close = true

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return bodyBytes, nil
}
