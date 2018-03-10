package gojenkins

import "net/http"

// Configuration is struct for Jenkins access information
type Configuration struct {
	User         string `yaml:"user"`         // Username for Jenkins
	Pass         string `yaml:"pass"`         // Password from Jenkins Username
	URL          string `yaml:"url"`          // Jenkins URL
	URLExtraPath string `yaml:"urlExtraPath"` // Jenkins configuration path
	Port         string `yaml:"port"`         // Jenkins Port
	Project      string `yaml:"project"`      // Some projects have more than one Jenkins instance, so just lable as you wish
	Crumb        string `yaml:"crumb"`        // Jenkins Crumb
	Timeout      int    `yaml:"timeout"`      // Jenkins request timeout
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
	Header http.Header
	Body   []byte
}
