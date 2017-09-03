# GoJenkins [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/nelsonjvf/gojenkins) [![Build Status](http://img.shields.io/travis/fatih/structs.svg?style=flat-square)]() [![Coverage Status](http://img.shields.io/coveralls/fatih/structs.svg?style=flat-square)]()

GoJenkins is a basic and generic package to interact with Jenkins REST API.
The idea of this package is to make your life easy, instead of learning the Jenkins REST API, you just need to set your configuration and get the information.
All the methos available return an Go object with the smae struter of the response from Jenkins REST API.
Fell free to add make comments and review the code.

## Install

```bash
go get github.com/nelsonjvf/gojenkins
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
// Get Jira issue information
gojenkins.RunJob()

// Search a string in Jira
gojenkins.GetLogs()
```

There is the full code to test it easly:

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
  runJobResponse := gojira.RunJob(issueToSearch)
  fmt.Println(runJobResponse)

  fmt.Println("Calling RequestSearch method:")

  jobId := runJobResponse.id
  getLogsResponse := gojira.GetLogs(jobId)
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

We can get an issue information:

```RunJob(issueId)```

We can also search in Jira:

```GetLogs(query)```

## Credits

 * [Nelson Ferreira](https://github.com/nelsonjvf)

## License

The MIT License (MIT) - see LICENSE.md for more details
