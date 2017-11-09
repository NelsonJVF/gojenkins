package gojenkins

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"encoding/json"
	"strconv"
	"bytes"
	"time"
)

// Configuration is struct for Jenkins access information
type Configuration struct {
	User string `yaml:"user"` // Username for Jenkins
	Pass string `yaml:"pass"` // Password from Jenkins Username
	Url string `yaml:"url"`	// Jenkins URL
	UrlExtraPath string `yaml:"urlExtraPath"`	// Jenkins configuration path
	Port string `yaml:"port"` // Jenkins Port
	Project string `yaml:"project"` // Some projects have more than one Jenkins instance, so just lable as you wish
	Crumb string `yaml:"crumb"` // Jenkins Crumb
	Timeout int `yaml:"timeout"` // Jenkins request timeout
}

// getCrumbResponse is a struct to get the user crumb response
type getCrumbResponse struct {
	Class             string `json:"_class"`
	Crumb             string `json:"crumb"`
	CrumbRequestField string `json:"crumbRequestField"`
}

// JenkinsJobsResponse is a struct for Jenkins Jobs Response
type JenkinsJobsResponse struct {
	Class          string `json:"_class"`
	AssignedLabels []struct {
	} `json:"assignedLabels"`
	Mode            string      `json:"mode"`
	NodeDescription string      `json:"nodeDescription"`
	NodeName        string      `json:"nodeName"`
	NumExecutors    int         `json:"numExecutors"`
	Description     interface{} `json:"description"`
	Jobs            []struct {
		Class string `json:"_class"`
		Name  string `json:"name"`
		URL   string `json:"url"`
		Color string `json:"color"`
	} `json:"jobs"`
	OverallLoad struct {
	} `json:"overallLoad"`
	PrimaryView struct {
		Class string `json:"_class"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"primaryView"`
	QuietingDown   bool `json:"quietingDown"`
	SlaveAgentPort int  `json:"slaveAgentPort"`
	UnlabeledLoad  struct {
		Class string `json:"_class"`
	} `json:"unlabeledLoad"`
	UseCrumbs   bool `json:"useCrumbs"`
	UseSecurity bool `json:"useSecurity"`
	Views       []struct {
		Class string `json:"_class"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"views"`
}

// JenkinsJobsLastBuildResponse is a struct for Jenkins Jobs Last Build esponse
type JenkinsJobsLastBuildResponse struct {
	Class   string `json:"_class"`
	Actions []struct {
		Class  string `json:"_class"`
		Causes []struct {
			Class            string `json:"_class"`
			ShortDescription string `json:"shortDescription"`
			UserID           string `json:"userId"`
			UserName         string `json:"userName"`
		} `json:"causes"`
	} `json:"actions"`
	Artifacts         []interface{} `json:"artifacts"`
	Building          bool          `json:"building"`
	Description       interface{}   `json:"description"`
	DisplayName       string        `json:"displayName"`
	Duration          int           `json:"duration"`
	EstimatedDuration int           `json:"estimatedDuration"`
	Executor          interface{}   `json:"executor"`
	FullDisplayName   string        `json:"fullDisplayName"`
	ID                string        `json:"id"`
	KeepLog           bool          `json:"keepLog"`
	Number            int           `json:"number"`
	QueueID           int           `json:"queueId"`
	Result            string        `json:"result"`
	Timestamp         int64         `json:"timestamp"`
	URL               string        `json:"url"`
	BuiltOn           string        `json:"builtOn"`
	ChangeSet         struct {
		Class string        `json:"_class"`
		Items []interface{} `json:"items"`
		Kind  interface{}   `json:"kind"`
	} `json:"changeSet"`
}

// JenkinsBuildsJobResponse is a struct for Jenkins Builds Job esponse
type JenkinsBuildsJobResponse struct {
	Class  string `json:"_class"`
	Builds []struct {
		Class     string `json:"_class"`
		ID        string `json:"id"`
		Number    int    `json:"number"`
		QueueID   int    `json:"queueId"`
		Result    string `json:"result"`
		Timestamp int64  `json:"timestamp"`
	} `json:"builds"`
}

// JenkinsJobDetailsResponse is a struct for Jenkins Job Details esponse
type JenkinsJobDetailsResponse struct {
	Class   string `json:"_class"`
	Actions []struct {
		Class      string `json:"_class,omitempty"`
		Parameters []struct {
			Class string `json:"_class"`
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"parameters,omitempty"`
		Causes []struct {
			Class            string `json:"_class"`
			ShortDescription string `json:"shortDescription"`
			UserID           string `json:"userId"`
			UserName         string `json:"userName"`
		} `json:"causes,omitempty"`
	} `json:"actions"`
	Artifacts         []interface{} `json:"artifacts"`
	Building          bool          `json:"building"`
	Description       interface{}   `json:"description"`
	DisplayName       string        `json:"displayName"`
	Duration          int           `json:"duration"`
	EstimatedDuration int           `json:"estimatedDuration"`
	Executor          interface{}   `json:"executor"`
	FullDisplayName   string        `json:"fullDisplayName"`
	ID                string        `json:"id"`
	KeepLog           bool          `json:"keepLog"`
	Number            int           `json:"number"`
	QueueID           int           `json:"queueId"`
	Result            string        `json:"result"`
	Timestamp         int64         `json:"timestamp"`
	URL               string        `json:"url"`
	BuiltOn           string        `json:"builtOn"`
	ChangeSet         struct {
		Class string        `json:"_class"`
		Items []interface{} `json:"items"`
		Kind  interface{}   `json:"kind"`
	} `json:"changeSet"`
	Culprits []interface{} `json:"culprits"`
}

// hTTPResponse is a struct for the http response
type hTTPResponse struct {
	Header 	http.Header
	Body 		[]byte
}

// Config is a variable to store the Jenkins server information
var Config []Configuration

// hTTPRequest is generic HTTP caller
func hTTPRequest(url string, method string, user string, pass string, crumb string, timeout int, parameters url.Values) (hTTPResponse, error) {
	var hTTPResp hTTPResponse

	timeoutVal := time.Duration(time.Duration(timeout) * time.Second)
	client := &http.Client{
		Timeout: timeoutVal,
	}
	r, _ := http.NewRequest(method, url, bytes.NewBufferString(parameters.Encode()))

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

// prepareJenkinsCall is a function to call jenkins rest api
func prepareJenkinsCall(project string, urlPath string, method string, parameters url.Values) (hTTPResponse, error) {
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
			URL = c.Url
			URLExtraPath = c.UrlExtraPath
			port = c.Port
			crumb = c.Crumb
			timeout = c.Timeout

			break
		}
	}

	if(len(url) == 0) {
		err := fmt.Errorf(" ---------- Jenkins configuration is missing  ---------- ")
		return resp, err
	}

	if(len(crumb) == 0) {
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

// RunJenkinsJob is a function to run a jenkins job
func RunJenkinsJob(project string, job string, parameters url.Values) (string, error) {
	var returnJobId string
	var path string

	path = fmt.Sprintf("/job/%s/buildWithParameters?delay=0sec", job)

	jobsResp, err := prepareJenkinsCall(project, path, "POST", parameters)
	if err != nil {
		return "", err
	}

	returnJobId = jobsResp.Header.Get("Location")

	return returnJobId, nil
}

// GetJenkinsJobs is a function to get jenkins logs
func GetJenkinsJobs(project string) (JenkinsJobsResponse, error) {
	var jenkinsJobs JenkinsJobsResponse

	jobsResp, err := prepareJenkinsCall(project, "/api/json?pretty=true", "GET", nil)
	if err != nil {
		return jenkinsJobs, err
	}

	if err := json.Unmarshal(jobsResp.Body, &jenkinsJobs); err != nil {
		return jenkinsJobs, err
	}

	return jenkinsJobs, nil
}

// GetLastBuild is a function to get the last build
func GetLastBuild(project string, job string) (JenkinsJobsLastBuildResponse, error) {
	var lastBuildResp JenkinsJobsLastBuildResponse
	var URL string

	URL = fmt.Sprintf("/job/%s/lastBuild/api/json?pretty=true", job)

	lastBuildResponse, errCall := prepareJenkinsCall(project, URL, "GET", nil)
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
	var jobDetailsResp JenkinsJobDetailsResponse
	var URL string

	URL = fmt.Sprintf("/job/%s/%s/api/json", job, strconv.Itoa(number))

	jobDetailsResponse, errCall := prepareJenkinsCall(project, URL, "GET", nil)
	if errCall != nil {
		return jobDetailsResp, errCall
	}

	if err := json.Unmarshal(jobDetailsResponse.Body, &jobDetailsResp); err != nil {
		return jobDetailsResp, err
	}

	return jobDetailsResp, nil
}

// GetJobLogs is a function to get job logs
func GetJobLogs(project string, job string, number int) (string, error) {
	var URL string

	URL = fmt.Sprintf("/job/%s/%s/consoleText", job, strconv.Itoa(number))
	tempResp, errCall := prepareJenkinsCall(project, URL, "GET", nil)
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
	tempResp, errCall := prepareJenkinsCall(project, URL, "GET", nil)
	if errCall != nil {
		return buildsJob, errCall
	}

	if err := json.Unmarshal(tempResp.Body, &buildsJob); err != nil {
		return buildsJob, err
	}

	return buildsJob, nil
}

// GetJobIdFromBuild is a function to get the job id from a given build
func GetJobIdFromBuild(project string, job string, buildId int) int {
	var jobId int

	resp, err := GetBuildsJob(project, job)
	if err != nil {
		log.Println(err)
	}

	for _, build := range resp.Builds {
		if build.QueueID == buildId {
			jobId = build.Number
		}
	}

	return jobId
}
