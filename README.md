# GoJenkins [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/NelsonJVF/gojenkins/pkg)

GoJenkins is a basic and generic package to interact with Jenkins REST API.
The idea of this package is to make your life easy, instead of learning the Jenkins REST API, you just need to set your configuration and get the information.
All the methods available return a Go object with the same structure of the response from Jenkins REST API.
Fell free to add make comments and review the code.

## Install

```bash
go get github.com/nelsonjvf/gojenkins/pkg
```

## Usage and Examples

Simply use the methods available to interact with Jenkins:

```go

// Run Jenkins Job
func RunJob(URL, username, password, job, parameters)

// Get jenkins logs
func GetJobs(URL, username, password)

// Get the last build
func GetLastBuild(URL, username, password, job)

// Get job details
func GetJobDetails(URL, username, password, job, number)

// Get job logs
func GetJobLogs(URL, username, password, job, number)

// Geet builds from a given job
func GetBuildsJob(URL, username, password, job)

// Get the job id from a given build
func GetJobIDFromBuild(URL, username, password, job, buildID)

```

Here is the full code to test it easily:

```go
package main

import (
	"log"

	"github.com/nelsonjvf/gojenkins"
)

func main() {
	
	log.Println("Run jenkins job..")

	jogToRun := "deploy-application"
	getJobsResponse := gojenkins.RunJob(jenkinsURL, jenkinsUsername, jenkinsPassword, jogToRun)
	log.Println("Build ID - ", getJobsResponse)

	log.Println("Get jenkins job log..")

	buildID := runJobResponse.id
	getLogsResponse := gojenkins.GetJobLogs(jenkinsURL, jenkinsUsername, jenkinsPassword, jogToRun, buildID)
	log.Println("Logs:")
	log.Println(getLogsResponse)

}
```

### Go Jenkins methods

We can run jenkins jobs:

```RunJob(JenkinsURL, JenkinsUsername, JenkinsPassword, JobName)```

We can also get the logs:

```GetJobLogs(JenkinsURL, JenkinsUsername, JenkinsPassword, JobName, buildID)```

Get last build from job:

```GetLastBuild(JenkinsURL, JenkinsUsername, JenkinsPassword, JobName)```

## Credits

 * [Nelson Ferreira](https://github.com/nelsonjvf)

## License

The MIT License (MIT) - see LICENSE.md for more details
