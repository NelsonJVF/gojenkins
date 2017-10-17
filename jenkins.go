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

/*
	Struct for Jenkins access information
 */
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

/*
	Get User Crumb Response
 */
type getCrumbResponse struct {
	Class             string `json:"_class"`
	Crumb             string `json:"crumb"`
	CrumbRequestField string `json:"crumbRequestField"`
}

/*
	Get Jenkings Jobs Response
 */
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

type hTTPResponse struct {
	Header 	http.Header
	Body 		[]byte
}

var Config []Configuration

/*
	Generic HTTP caller
 */
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

/*
	Get Jenkins Crumb
 */
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

/*
	Prepare Jenkins Call
		- Get and Set Jenkins information
		- Get and Set Crumb information
 */
func prepareJenkinsCall(project string, urlPath string, method string, parameters url.Values) (hTTPResponse, error) {
	var user string
	var pass string
	var url string
	var urlExtraPath string
	var port string
	var crumb string
	var timeout int
	var err error
	var resp hTTPResponse

	for _, c := range Config {
		if c.Project == project {
			user = c.User
			pass = c.Pass
			url = c.Url
			urlExtraPath = c.UrlExtraPath
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
		crumb, err = getCrumb(user, pass, url, port, timeout, urlExtraPath)
		if err != nil {
			return resp, err
		}

		for _, c := range Config {
			if c.Project == project {
				c.Crumb = crumb
			}
		}
	}

	urlCall := fmt.Sprintf("%s:%s%s%s", url, port, urlExtraPath, urlPath)

	resp, err = hTTPRequest(urlCall, method, user, pass, crumb, timeout, parameters)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func RunJenkinsJob(project string, job string, parameters string) (string, error) {
	var returnJobId string
	var path string

	if project == "Production" {
		data := url.Values{}
		data.Set("multiline", parameters)

		path = fmt.Sprintf("/job/%s/buildWithParameters?delay=0sec", job)

		jobsResp, err := prepareJenkinsCall(project, path, "POST", data)
		if err != nil {
			return "", err
		}

		returnJobId = jobsResp.Header.Get("Location")
	} else {
		if len(parameters) == 0 {
			path = fmt.Sprintf("/job/%s/buildWithParameters?delay=0sec", job)
		} else {
			path = fmt.Sprintf("/job/%s/buildWithParameters?delay=0sec&%s", job, parameters)
		}

		jobsResp, err := prepareJenkinsCall(project, path, "POST", nil)
		if err != nil {
			return "", err
		}

		fmt.Println(jobsResp.Header)

		returnJobId = jobsResp.Header.Get("Location")
	}

	return returnJobId, nil
}

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

func GetLastBuild(project string, job string) (JenkinsJobsLastBuildResponse, error) {
	var lastBuildResp JenkinsJobsLastBuildResponse
	var url string

	url = fmt.Sprintf("/job/%s/lastBuild/api/json?pretty=true", job, "GET")

	lastBuildResponse, errCall := prepareJenkinsCall(project, url, "GET", nil)
	if errCall != nil {
		return lastBuildResp, errCall
	}

	if err := json.Unmarshal(lastBuildResponse.Body, &lastBuildResp); err != nil {
		return lastBuildResp, err
	}

	return lastBuildResp, nil
}

func GetJobDetails(project string, job string, number int) (JenkinsJobDetailsResponse, error) {
	var jobDetailsResp JenkinsJobDetailsResponse
	var url string

	url = fmt.Sprintf("/job/%s/%s/api/json", job, strconv.Itoa(number))

	jobDetailsResponse, errCall := prepareJenkinsCall(project, url, "GET", nil)
	if errCall != nil {
		return jobDetailsResp, errCall
	}

	if err := json.Unmarshal(jobDetailsResponse.Body, &jobDetailsResp); err != nil {
		return jobDetailsResp, err
	}

	return jobDetailsResp, nil
}

func GetJobLogs(project string, job string, number int) (string, error) {
	var url string

	url = fmt.Sprintf("/job/%s/%s/consoleText", job, strconv.Itoa(number))
	tempResp, errCall := prepareJenkinsCall(project, url, "GET", nil)
	if errCall != nil {
		return "", errCall
	}

	return string(tempResp.Body), nil
}

func GetBuildsJob(project string, job string) (JenkinsBuildsJobResponse, error) {
	var buildsJob JenkinsBuildsJobResponse
	var url string

	url = fmt.Sprintf("/job/%s/api/json?tree=builds[number,status,timestamp,id,queueId,result]", job)
	tempResp, errCall := prepareJenkinsCall(project, url, "GET", nil)
	if errCall != nil {
		return buildsJob, errCall
	}

	if err := json.Unmarshal(tempResp.Body, &buildsJob); err != nil {
		return buildsJob, err
	}

	return buildsJob, nil
}

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
