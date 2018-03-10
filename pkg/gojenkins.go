package gojenkins

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// DefaultTimeout is the default timeout for jenkins call
var DefaultTimeout = 30

// hTTPRequest is generic HTTP caller
func hTTPRequest(URL string, method string, user string, pass string, crumb string, parameters url.Values) (hTTPResponse, error) {
	var hTTPResp hTTPResponse

	timeoutVal := time.Duration(time.Duration(DefaultTimeout) * time.Second)
	client := &http.Client{
		Timeout: timeoutVal,
	}
	r, _ := http.NewRequest(method, URL, bytes.NewBufferString(parameters.Encode()))

	r.SetBasicAuth(user, pass)
	r.Header.Set("Jenkins-Crumb", crumb)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, errDo := client.Do(r)
	if errDo != nil {
		return hTTPResp, errDo
	}

	body, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		log.Printf("ioutil.ReadAll err   #%v ", errRead)
		return hTTPResp, errRead
	}

	hTTPResp.Header = resp.Header
	hTTPResp.Body = body

	return hTTPResp, nil
}

// getCrumb is a function to get the jenkins user crumb
func getCrumb(url string, user string, pass string) (string, error) {
	var crumbResp getCrumbResponse

	urlCrumb := fmt.Sprintf("%s%s", url, "/crumbIssuer/api/json")
	resp, err := hTTPRequest(urlCrumb, "GET", user, pass, "", nil)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(resp.Body, &crumbResp); err != nil {
		return "", err
	}

	return crumbResp.Crumb, nil
}

// prepareCall is a function to call jenkins rest api
func prepareCall(URL string, URLPath string, username string, password string, method string, parameters url.Values) (hTTPResponse, error) {

	var resp hTTPResponse

	if len(URL) == 0 {
		return resp, errors.New("Please, provide the Jenkins URL")
	}

	// Get Jenkins Crumb for configured user
	crumb, err := getCrumb(URL, username, password)
	if err != nil {
		return resp, err
	}

	urlCall := fmt.Sprintf("%s%s", URL, URLPath)

	resp, err = hTTPRequest(urlCall, method, username, password, crumb, parameters)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RunJob is a function to run a jenkins job
func RunJob(URL string, username string, password string, job string, parameters url.Values) (string, error) {
	var jobID string
	var URLPath string

	URLPath = fmt.Sprintf("/job/%s/buildWithParameters?delay=0sec", job)

	resp, err := prepareCall(URL, URLPath, username, password, "POST", parameters)
	if err != nil {
		return "", err
	}

	jobID = resp.Header.Get("Location")

	return jobID, nil
}

// GetJobs is a function to get jenkins logs
func GetJobs(URL string, username string, password string) (JenkinsJobsResponse, error) {
	var jobs JenkinsJobsResponse

	resp, err := prepareCall(URL, "/api/json?pretty=true", username, password, "GET", nil)
	if err != nil {
		return jobs, err
	}

	if err := json.Unmarshal(resp.Body, &jobs); err != nil {
		return jobs, err
	}

	return jobs, nil
}

// GetLastBuild is a function to get the last build
func GetLastBuild(URL string, username string, password string, job string) (JenkinsJobsLastBuildResponse, error) {
	var lastBuildResp JenkinsJobsLastBuildResponse
	var URLPath string

	URLPath = fmt.Sprintf("/job/%s/lastBuild/api/json?pretty=true", job)

	resp, errCall := prepareCall(URL, URLPath, username, password, "GET", nil)
	if errCall != nil {
		return lastBuildResp, errCall
	}

	if err := json.Unmarshal(resp.Body, &lastBuildResp); err != nil {
		return lastBuildResp, err
	}

	return lastBuildResp, nil
}

// GetJobDetails is a function to get job details
func GetJobDetails(URL string, username string, password string, job string, number int) (JenkinsJobDetailsResponse, error) {
	var jobDetails JenkinsJobDetailsResponse
	var URLPath string

	URLPath = fmt.Sprintf("/job/%s/%s/api/json", job, strconv.Itoa(number))

	resp, errCall := prepareCall(URL, URLPath, username, password, "GET", nil)
	if errCall != nil {
		return jobDetails, errCall
	}

	if err := json.Unmarshal(resp.Body, &jobDetails); err != nil {
		return jobDetails, err
	}

	return jobDetails, nil
}

// GetJobLogs is a function to get job logs
func GetJobLogs(URL string, username string, password string, job string, number int) (string, error) {
	var URLPath string

	URLPath = fmt.Sprintf("/job/%s/%s/consoleText", job, strconv.Itoa(number))

	resp, errCall := prepareCall(URL, URLPath, username, password, "GET", nil)
	if errCall != nil {
		return "", errCall
	}

	return string(resp.Body), nil
}

// GetBuildsJob is a function to get builds from a given job
func GetBuildsJob(URL string, username string, password string, job string) (JenkinsBuildsJobResponse, error) {
	var buildsJob JenkinsBuildsJobResponse
	var URLPath string

	URLPath = fmt.Sprintf("/job/%s/api/json?tree=builds[number,status,timestamp,id,queueId,result]", job)

	resp, errCall := prepareCall(URL, URLPath, username, password, "GET", nil)
	if errCall != nil {
		return buildsJob, errCall
	}

	if err := json.Unmarshal(resp.Body, &buildsJob); err != nil {
		return buildsJob, err
	}

	return buildsJob, nil
}

// GetJobIDFromBuild is a function to get the job id from a given build
func GetJobIDFromBuild(URL string, username string, password string, job string, buildID int) (int, error) {
	var jobID int

	resp, err := GetBuildsJob(URL, username, password, job)
	if err != nil {
		return jobID, err
	}

	for _, build := range resp.Builds {
		if build.QueueID == buildID {
			jobID = build.Number
		}
	}

	return jobID, nil
}
