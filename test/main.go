package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/nelsonjvf/gojenkins/pkg"
	"gopkg.in/yaml.v2"
)

func init() {
	// Use yaml configuration file
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &gojenkins.Config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func main() {

	fmt.Println("Staring Testing..")
	fmt.Println("Setting Configuration")

	fmt.Println(gojenkins.Config)

	gojenkins.RunJob("Test Jenkins Server", "Run Application X", nil)
	gojenkins.GetJobLogs("Project", "Run Application X", 1)

}
