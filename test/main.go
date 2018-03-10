package main

import (
	"log"

	"github.com/nelsonjvf/gojenkins/pkg"
)

func main() {

	var jenkinsURL = "http://jira-server.com/"
	var jenkinsUsername = "my-jira-user"
	var jenkinsPassword = "my-jira-password"

	log.Println("Staring Testing..")

	var jenkinsJobName = "MyJenkinsJob"

	log.Println("Staring Jenkins Job - ", jenkinsJobName)
	gojenkins.RunJob(jenkinsURL, jenkinsUsername, jenkinsPassword, jenkinsJobName, nil)

	var buildNumber = 1

	log.Println("Staring Jenkins Job - ", jenkinsJobName)
	gojenkins.GetJobLogs(jenkinsURL, jenkinsUsername, jenkinsPassword, jenkinsJobName, buildNumber)

}
