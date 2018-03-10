package gojenkins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Config is a variable to store the Jenkins server information
var Config []Configuration

// hTTPRequest is generic HTTP caller
func hTTPRequest(URL string, method string, user string, pass string, crumb string, timeout int, parameters url.Values) (hTTPResponse, error) {
	var hTTPResp hTTPResponse

	timeoutVal := time.Duration(time.Duration(timeout) * time.Second)
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
func getCrumb(user string, pass string, url string, port string, timeout int, urlExtraPath string) (string, error) {
	var crumbResp getCrumbResponse

	urlCrumb := fmt.Sprintf("%s:%s%s%s", url, port, urlExtraPath, "/crumbIssuer/api/json")
	resp, err := hTTPRequest(urlCrumb, "GET", user, pass, "", timeout, nil)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(resp.Body, &crumbResp); err != nil {
		log.Println("Error getJenkinsCrumb()")
		return "", err
	}

	crumb := crumbResp.Crumb

	return crumb, nil
}

// prepareCall is a function to call jenkins rest api
func prepareCall(project string, urlPath string, method string, parameters url.Values) (hTTPResponse, error) {
	var username string
	var password string
	var URL string
	var URLExtraPath string
	var port string
	var crumb string
	var timeout int
	var err error
	var resp hTTPResponse

	for _, c := range Config {
		if c.Project == project {
			username = c.User
			password = c.Pass
			URL = c.URL
			URLExtraPath = c.URLExtraPath
			port = c.Port
			crumb = c.Crumb
			timeout = c.Timeout

			break
		}
	}

	if len(URL) == 0 {
		err := fmt.Errorf(" ---------- Jenkins configuration is missing  ---------- ")
		return resp, err
	}

	if len(crumb) == 0 {
		// Get Jenkins Crumb for configured user
		crumb, err = getCrumb(username, password, URL, port, timeout, URLExtraPath)
		if err != nil {
			return resp, err
		}

		for _, c := range Config {
			if c.Project == project {
				c.Crumb = crumb
			}
		}
	}

	urlCall := fmt.Sprintf("%s:%s%s%s", URL, port, URLExtraPath, urlPath)

	resp, err = hTTPRequest(urlCall, method, username, password, crumb, timeout, parameters)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RunJob is a function to run a jenkins job
func RunJob(project string, job string, parameters url.Values) (string, error) {
	var jobID string
	var path string

	path = fmt.Sprintf("/job/%s/buildWithParameters?delay=0sec", job)

	jobsResp, err := prepareCall(project, path, "POST", parameters)
	if err != nil {
		return "", err
	}

	jobID = jobsResp.Header.Get("Location")

	return jobID, nil
}

// GetJobs is a function to get jenkins logs
func GetJobs(project string) (JenkinsJobsResponse, error) {
	var jobs JenkinsJobsResponse

	jobsResp, err := prepareCall(project, "/api/json?pretty=true", "GET", nil)
	if err != nil {
		return jobs, err
	}

	if err := json.Unmarshal(jobsResp.Body, &jobs); err != nil {
		return jobs, err
	}

	return jobs, nil
}

// GetLastBuild is a function to get the last build
func GetLastBuild(project string, job string) (JenkinsJobsLastBuildResponse, error) {
	var lastBuildResp JenkinsJobsLastBuildResponse
	var URL string

	URL = fmt.Sprintf("/job/%s/lastBuild/api/json?pretty=true", job)

	lastBuildResponse, errCall := prepareCall(project, URL, "GET", nil)
	if errCall != nil {
		return lastBuildResp, errCall
	}

	if err := json.Unmarshal(lastBuildResponse.Body, &lastBuildResp); err != nil {
		return lastBuildResp, err
	}

	return lastBuildResp, nil
}

// GetJobDetails is a function to get job details
func GetJobDetails(project string, job string, number int) (JenkinsJobDetailsResponse, error) {
	var jobDetails JenkinsJobDetailsResponse
	var URL string

	URL = fmt.Sprintf("/job/%s/%s/api/json", job, strconv.Itoa(number))

	response, errCall := prepareCall(project, URL, "GET", nil)
	if errCall != nil {
		return jobDetails, errCall
	}

	if err := json.Unmarshal(response.Body, &jobDetails); err != nil {
		return jobDetails, err
	}

	return jobDetails, nil
}

// GetJobLogs is a function to get job logs
func GetJobLogs(project string, job string, number int) (string, error) {
	var URL string

	URL = fmt.Sprintf("/job/%s/%s/consoleText", job, strconv.Itoa(number))
	tempResp, errCall := prepareCall(project, URL, "GET", nil)
	if errCall != nil {
		return "", errCall
	}

	return string(tempResp.Body), nil
}

// GetBuildsJob is a function to get builds from a given job
func GetBuildsJob(project string, job string) (JenkinsBuildsJobResponse, error) {
	var buildsJob JenkinsBuildsJobResponse
	var URL string

	URL = fmt.Sprintf("/job/%s/api/json?tree=builds[number,status,timestamp,id,queueId,result]", job)
	tempResp, errCall := prepareCall(project, URL, "GET", nil)
	if errCall != nil {
		return buildsJob, errCall
	}

	if err := json.Unmarshal(tempResp.Body, &buildsJob); err != nil {
		return buildsJob, err
	}

	return buildsJob, nil
}

// GetJobIDFromBuild is a function to get the job id from a given build
func GetJobIDFromBuild(project string, job string, buildID int) int {
	var jobID int

	resp, err := GetBuildsJob(project, job)
	if err != nil {
		log.Println(err)
	}

	for _, build := range resp.Builds {
		if build.QueueID == buildID {
			jobID = build.Number
		}
	}

	return jobID
}
