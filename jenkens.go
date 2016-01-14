package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var (
	project_url  *string
	project_name *string
	project_id   *string
)

func init() {
	project_url = flag.String("url", "git@git.xxx.com:beatles/serv_base.git", `Please enter the ssh_url_to_repo, which such as "git@git.xxx.com:beatles/serv_base.git"`)
	flag.Parse()
	if *project_url == "" {
		flag.PrintDefaults()
	}
	re := regexp.MustCompile("^git@git.xxx.com:(.*)/(.*)\\.git$")
	ss := re.FindAllStringSubmatch(*project_url, -1)
	if len(ss) >= 1 && len(ss[0]) >= 3 {
		project_name = &ss[0][1]
		project_id = &ss[0][2]
	} else {
		log.Fatalf("ERROR: find repo or id error\n")
	}
}

func show() {
	fmt.Printf("project_url : %s\n", *project_url)
	fmt.Printf("project_name: %s\n", *project_name)
	fmt.Printf("project_id  : %s\n", *project_id)
}

type Project struct {
	Id              int    `json:"id"`
	Ssh_url_to_repo string `json:"ssh_url_to_repo"`
}
type Auth struct {
	Title string `json:"title"`
	Key   string `json:"key"`
}

func main() {
	show()
	jenkins_url := fmt.Sprintf("https://git.xxx.com/api/v3/projects/search/%s", *project_id)
	println(jenkins_url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", jenkins_url, nil)
	if err != nil {
		log.Fatalf("ERROR: Get url error")
	}
	req.Header.Add("PRIVATE-TOKEN", "XXX")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("ERROR: Get url error")
	}
	defer resp.Body.Close()

	var projects []Project
	json_body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(json_body, &projects)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatalf("ERROR: Unmarshal json error")
	}

	var project Project
	for _, p := range projects {
		fmt.Printf("p.Ssh_url_to_repo  : %s\n", p.Ssh_url_to_repo)
		fmt.Printf("project_url         : %s\n", *project_url)
		if p.Ssh_url_to_repo == *project_url {
			project = p
			break
		}
	}

	if project.Id == 0 {
		log.Fatalf("ERROR: project no find")
	} else {
		fmt.Printf("Id: %d\n", project.Id)
	}

	deploy_key := "ssh-dss XXX"

	auth := Auth{
		Title: "deploy_key_for_jenkins",
		Key:   deploy_key,
	}

	jenkins_add_key_url := fmt.Sprintf("https://git.xxx.com/api/v3/projects/%d/keys", project.Id)

	post_json_data, err := json.Marshal(auth)
	if err != nil {
		log.Fatalf(err.Error())
	}
	//println(string(post_json_data))

	req, err = http.NewRequest("POST", jenkins_add_key_url, strings.NewReader(string(post_json_data)))
	if err != nil {
		log.Fatalf("ERROR: Get url error")
	}

	println(jenkins_add_key_url)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("PRIVATE-TOKEN", "XXX")

	resp, err = client.Do(req)
	if err != nil {
		log.Fatalf("ERROR: Get url error")
	}
	if resp.StatusCode != 200 {
		fmt.Printf("============\n")
		fmt.Printf("%s\n", resp.Status)
		fmt.Printf("============\n")
	} else {
		fmt.Printf("OK!\n")
	}
}
