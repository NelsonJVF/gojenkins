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
func hTTPRequest(url string, method string, user string, pass string, crumb string, parameters url.Values) hTTPResponse {
	var hTTPResp hTTPResponse

	client := &http.Client{}
	r, _ := http.NewRequest(method, url, bytes.NewBufferString(parameters.Encode()))

	r.SetBasicAuth(user, pass)
	r.Header.Set("Jenkins-Crumb", crumb)
	//r.Header.Set("Content-Type", "application/json")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(r)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err   #%v ", err)
	}

	hTTPResp.Header = resp.Header
	hTTPResp.Body = body

	return hTTPResp
}

/*
	Get Jenkins Crumb
 */
func getCrumb(user string, pass string, url string, port string, urlExtraPath string) string {
	var crumbResp getCrumbResponse

	urlCrumb := fmt.Sprintf("%s:%s%s%s", url, port, urlExtraPath, "/crumbIssuer/api/json")
	bodyCrumb := hTTPRequest(urlCrumb, "GET", user, pass, "", nil)

	if err := json.Unmarshal(bodyCrumb.Body, &crumbResp); err != nil {
		panic(err)
	}

	crumb := crumbResp.Crumb

	return crumb
}

/*
	Prepare Jenkins Call
		- Get and Set Jenkins information
		- Get and Set Crumb information
 */
func prepareJenkinsCall(project string, urlPath string, method string, parameters url.Values) hTTPResponse {
	var user string
	var pass string
	var url string
	var urlExtraPath string
	var port string
	var crumb string

	for _, c := range Config {
		if c.Project == project {
			user = c.User
			pass = c.Pass
			url = c.Url
			urlExtraPath = c.UrlExtraPath
			port = c.Port
			crumb = c.Crumb

			break
		}
	}

	if(len(url) == 0) {
		log.Printf(" ---------- Jenkins configuration is missing  ---------- ")
		return hTTPResponse{
			Header: nil,
			Body: nil,
		}
	}

	if(len(crumb) == 0) {
		// Get Jenkins Crumb for configured user
		crumb = getCrumb(user, pass, url, port, urlExtraPath)

		for _, c := range Config {
			if c.Project == project {
				c.Crumb = crumb
			}
		}
	}

	urlCall := fmt.Sprintf("%s:%s%s%s", url, port, urlExtraPath, urlPath)

	callResp := hTTPRequest(urlCall, method, user, pass, crumb, parameters)

	return callResp
}

func RunJenkinsJob(project string, job string, parameters string) string {
	var returnJobId string
	var path string

	if project == "SpecificJenkins" {
		data := url.Values{}
		data.Set("jenkinsVAR", parameters)

		path = fmt.Sprintf("/job/%s/buildWithParameters?delay=0sec", job)

		jobsResp := prepareJenkinsCall(project, path, "POST", data)

		returnJobId = jobsResp.Header.Get("Location")
	} else {
		if len(parameters) == 0 {
			path = fmt.Sprintf("/job/%s/buildWithParameters?delay=0sec", job)
		} else {
			path = fmt.Sprintf("/job/%s/buildWithParameters?delay=0sec&%s", job, parameters)
		}

		jobsResp := prepareJenkinsCall(project, path, "POST", nil)

		returnJobId = jobsResp.Header.Get("Location")
	}

	return returnJobId
}

func GetJenkinsJobs(project string) JenkinsJobsResponse {
	var jenkinsJobs JenkinsJobsResponse

	jobsResp := prepareJenkinsCall(project, "/api/json?pretty=true", "GET", nil)

	if err := json.Unmarshal(jobsResp.Body, &jenkinsJobs); err != nil {
		panic(err)
	}

	return jenkinsJobs
}

func GetLastBuild(project string, job string) JenkinsJobsLastBuildResponse {
	var lastBuildResp JenkinsJobsLastBuildResponse
	var url string

	url = fmt.Sprintf("/job/%s/lastBuild/api/json?pretty=true", job, "GET")

	lastBuildResponse := prepareJenkinsCall(project, url, "GET", nil)
	if err := json.Unmarshal(lastBuildResponse.Body, &lastBuildResp); err != nil {
		panic(err)
	}

	return lastBuildResp
}

func GetJobDetails(project string, job string, number int) JenkinsJobDetailsResponse {
	var jobDetailsResp JenkinsJobDetailsResponse
	var url string

	url = fmt.Sprintf("/job/%s/%s/api/json", job, strconv.Itoa(number))

	jobDetailsResponse := prepareJenkinsCall(project, url, "GET", nil)
	if err := json.Unmarshal(jobDetailsResponse.Body, &jobDetailsResp); err != nil {
		panic(err)
	}

	return jobDetailsResp
}

func GetJobLogs(project string, job string, number int) string {
	var url string

	url = fmt.Sprintf("/job/%s/%s/consoleText", job, strconv.Itoa(number))
	tempResp := prepareJenkinsCall(project, url, "GET", nil)

	return string(tempResp.Body)
}

func GetBuildsJob(project string, job string) JenkinsBuildsJobResponse {
	var buildsJob JenkinsBuildsJobResponse
	var url string

	url = fmt.Sprintf("/job/%s/api/json?tree=builds[number,status,timestamp,id,queueId,result]", job)
	tempResp := prepareJenkinsCall(project, url, "GET", nil)
	if err := json.Unmarshal(tempResp.Body, &buildsJob); err != nil {
		panic(err)
	}

	return buildsJob
}

func GetJobIdFromBuild(project string, job string, buildId int) int {
	var jobId int

	resp := GetBuildsJob(project, job)

	for _, build := range resp.Builds {
		if build.QueueID == buildId {
			jobId = build.Number
		}
	}

	return jobId
}
