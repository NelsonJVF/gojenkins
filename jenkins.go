package gojenkins

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"strconv"
)

/*
	Struct for Jenkins access information
 */
type Configuration struct {
	User string `yaml:"user"` // Username for Jenkins
	Pass string `yaml:"pass"` // Password from Jenkins Username
	Url string `yaml:"url"`	// Jenkins URL
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

type hTTPResponse struct {
	Header 	http.Header
	Body 		[]byte
}

var Config []Configuration

/*
	Generic HTTP caller
 */
func hTTPRequest(url string, method string, user string, pass string, crumb string) hTTPResponse {
	var hTTPResp hTTPResponse

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Printf("http.NewRequest err   #%v ", err)
	}

	req.SetBasicAuth(user, pass)
	req.Header.Set("Content-Type", "application/json")
	if(len(crumb) != 0) {
		req.Header.Set("Jenkins-Crumb", crumb)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("http.DefaultClient.Do err   #%v ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err   #%v ", err)
	}

	hTTPResp.Header = resp.Header
	hTTPResp.Body = body

	return hTTPResp
}

func prepareJenkinsCall(project string, urlPath string, method string) hTTPResponse {
	var user string
	var pass string
	var url string
	var port string
	var crumb string

	for _, c := range Config {
		if c.Project == project {
			user = c.User
			pass = c.Pass
			url = c.Url
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
		var crumbResp getCrumbResponse

		urlCrumb := fmt.Sprintf("%s:%s%s", url, port, "/crumbIssuer/api/json")
		bodyCrumb := hTTPRequest(urlCrumb, "GET", user, pass, "")

		if err := json.Unmarshal(bodyCrumb.Body, &crumbResp); err != nil {
			panic(err)
		}

		crumb = crumbResp.Crumb

		for _, c := range Config {
			if c.Project == project {
				c.Crumb = crumb
			}
		}
	}

	urlCall := fmt.Sprintf("%s:%s%s", url, port, urlPath)
	callResp := hTTPRequest(urlCall, method, user, pass, crumb)

	return callResp
}

func RunJenkinsJob(project string, job string, parameters string) string {
	var returnJobId string
	var path string

	if len(parameters) == 0 {
		path = fmt.Sprintf("/job/%s/buildWithParameters?delay=0sec", job)
	} else {
		path = fmt.Sprintf("/job/%s/buildWithParameters?delay=0sec&%s", job, parameters)
	}

	jobsResp := prepareJenkinsCall(project, path, "POST")
	returnJobId = jobsResp.Header.Get("Location")

	return returnJobId
}

func GetJenkinsJobs(project string) JenkinsJobsResponse {
	var jenkinsJobs JenkinsJobsResponse

	jobsResp := prepareJenkinsCall(project, "/api/json?pretty=true", "GET")

	if err := json.Unmarshal(jobsResp.Body, &jenkinsJobs); err != nil {
		panic(err)
	}

	return jenkinsJobs
}

func GetLastBuild(project string, job string) JenkinsJobsLastBuildResponse {
	var lastBuildResp JenkinsJobsLastBuildResponse
	var url string

	url = fmt.Sprintf("/job/%s/lastBuild/api/json?pretty=true", job, "GET")

	tempResp := prepareJenkinsCall(project, url, "GET")
	if err := json.Unmarshal(tempResp.Body, &lastBuildResp); err != nil {
		panic(err)
	}

	return lastBuildResp
}

func GetJobLogs(project string, job string, number int) string {
	var url string

	url = fmt.Sprintf("/job/%s/%s/consoleText", job, strconv.Itoa(number))
	tempResp := prepareJenkinsCall(project, url, "GET")

	return string(tempResp.Body)
}

func GetBuildsJob(project string, job string) JenkinsBuildsJobResponse {
	var buildsJob JenkinsBuildsJobResponse
	var url string

	url = fmt.Sprintf("/job/%s/api/json?tree=builds[number,status,timestamp,id,queueId,result]", job)
	tempResp := prepareJenkinsCall(project, url, "GET")
	if err := json.Unmarshal(tempResp.Body, &buildsJob); err != nil {
		panic(err)
	}

	return buildsJob
}
