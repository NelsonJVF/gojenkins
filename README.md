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

First of all we need to configure and set your Jenkins information. For that we can use our config.yaml file as example and the following init function:

```go
func init() {
	// Use yaml configuration file
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &gojira.Config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}
```

After that you can simply use the methods available to interact with Jenkins:

```go
// Run Jenkins Job
gojenkins.RunJob()

// Get Job Logs
gojenkins.GetLogs()
```

Here is the full code to test it easily:

```go
package main

import (
	"fmt"
	"github.com/nelsonjvf/gojenkins"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
)

func init() {
	// Use yaml configuration file
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &gojira.Config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func main() {
  fmt.Println("Starting test..")

  fmt.Println("Calling Run Job method:")

  jogToRun := "deploy-application"
  getJobsResponse := gojenkins.GetJenkinsJobs(jenkinsProject)
  fmt.Println(getJobsResponse)

  fmt.Println("Calling RequestSearch method:")

  jobId := runJobResponse.id
  getLogsResponse := gojenkins.RunJenkinsJob(jobId)
  fmt.Println(getLogsResponse)
}
```

Don't forget your yaml configuration file

```yaml
user: JenkinsUser
pass: JenkinsPass
url: http:/jenkins.dev.com:8080/
```

### GoJira methods

We can run jenkins jobs:

```RunJob(issueId)```

We can also get the logs:

```GetLogs(query)```

## Credits

 * [Nelson Ferreira](https://github.com/nelsonjvf)

## License

The MIT License (MIT) - see LICENSE.md for more details
