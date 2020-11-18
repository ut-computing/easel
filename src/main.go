package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

const (
	appName                       = "Easel"
	cmdName                       = "easel"
	perUserDotFile                = "." + cmdName + "rc"
	urlPrefix                     = "/api/v1"
	coursesPath                   = "/courses"
	coursePath                    = coursesPath + "/%d"
	quizzesPath                   = coursePath + "/quizzes"
	quizPath                      = quizzesPath + "/%d"
	quizQuestionsPath             = quizPath + "/questions"
	quizSubmissionsPath           = quizPath + "/submissions"
	quizSubmissionPath            = quizSubmissionsPath + "/%d"
	quizSubmissionQuestionsPath   = "/quiz_submissions/%d/questions"
	quizReportsPath               = quizPath + "/reports"
	quizReportPath                = quizReportsPath + "/%d"
	assignmentsPath               = coursePath + "/assignments"
	assignmentPath                = assignmentsPath + "/%d"
	assignmentSubmissionsPath     = assignmentPath + "/submissions"
	gradeAssignmentSubmissionPath = assignmentSubmissionsPath + "/%d" // user_id
	progressPath                  = "/progress/%d"
)

var Config struct {
	Host      string `json:"host"`
	Token     string `json:"token"`
	apiReport bool
	apiDump   bool
}

func main() {
	log.SetFlags(log.Ltime)

	cmd := &cobra.Command{
		Use:   cmdName,
		Short: "Canvas shell management tool",
	}
	cmd.PersistentFlags().BoolVarP(&Config.apiReport, "api", "", false, "report all API requests")
	cmd.PersistentFlags().BoolVarP(&Config.apiDump, "api-dump", "", false, "dump API request and response data")

	cmdLogin := &cobra.Command{
		Use:   "login <hostname> <token>",
		Short: "login to Canvas",
		Long: fmt.Sprintf("To log in, open Canvas and click on Account > " +
			"Settings. Then under Approved Integrations click " +
			"'+ New Access Token' and fill out the form as desired. " +
			"Then click 'Generate Token'. Copy the token and paste here. " +
			"You should only need to do this once per machine."),
		Run: CommandLogin,
	}
	cmd.AddCommand(cmdLogin)

	cmdInit := &cobra.Command{
		Use:   "init",
		Short: "Initialize the db and point it to the given course",
		Long:  "TODO instructions",
		Run:   CommandInit,
	}
	cmd.AddCommand(cmdInit)

	cmdPull := &cobra.Command{
		Use:   "pull [component] [component_id]",
		Short: "pull a single component or all of that type if blank",
		Long:  "TODO instructions",
		Run:   CommandPull,
	}
	cmd.AddCommand(cmdPull)

	cmdPush := &cobra.Command{
		Use:   "push [component] [component_id]",
		Short: "push a single component or all of that type if blank",
		Long:  "TODO instructions",
		Run:   CommandPush,
	}
	cmd.AddCommand(cmdPush)

	cmd.Execute()
}

type LoginSession struct {
	Token string
}

func CommandLogin(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.Fatalf("Usage: %s login <hostname> <token>", os.Args[0])
	}
	hostname, token := args[0], args[1]

	protocol := "https://"
	if strings.HasPrefix(hostname, protocol) {
		hostname = hostname[len(protocol):]
	}
	if strings.HasSuffix(hostname, "/") {
		hostname = hostname[:len(hostname)-1]
	}

	// set up config
	Config.Host = hostname
	Config.Token = token

	// save config for later use
	mustWriteConfig()

	log.Println("login successful")
}

func CommandInit(cmd *cobra.Command, args []string) {
	mustLoadConfig(cmd)
	db := mustCreateDb()
	defer db.Close()
	courseId, err := getCourseIdFromUrl(args[0])
	if err != nil {
		log.Fatal(err.Error())
	}

	course, err := createCourse(db, courseId)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = course.dump()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func mustGetObject(path string, params url.Values, download interface{}) {
	doRequest(path, params, "GET", nil, download, false)
}

func getObject(path string, params url.Values, download interface{}) bool {
	return doRequest(path, params, "GET", nil, download, true)
}

func mustPostObject(path string, params url.Values, upload interface{}, download interface{}) {
	doRequest(path, params, "POST", upload, download, false)
}

func mustPutObject(path string, params url.Values, upload interface{}, download interface{}) {
	doRequest(path, params, "PUT", upload, download, false)
}

func doRequest(path string, params url.Values, method string, upload interface{}, download interface{}, notfoundokay bool) bool {
	if !strings.HasPrefix(path, "/") {
		log.Panicf("doRequest path must start with /")
	}
	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" {
		log.Panicf("doRequest only recognizes GET, POST, PUT, and DELETE methods")
	}
	url := fmt.Sprintf("https://%s%s%s", Config.Host, urlPrefix, path)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatalf("error creating http request: %v\n", err)
	}

	// add any parameters
	if params != nil && len(params) > 0 {
		req.URL.RawQuery = params.Encode()
	}

	if Config.apiReport {
		log.Printf("%s %s", method, req.URL)
	}

	// set the headers
	req.Header.Add("Authorization", "Bearer "+Config.Token)
	if download != nil {
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Accept-Encoding", "gzip")
	}

	// upload the payload if any
	if upload != nil && (method == "POST" || method == "PUT") {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Encoding", "gzip")
		payload := new(bytes.Buffer)
		gw := gzip.NewWriter(payload)
		uncompressed := new(bytes.Buffer)
		var jsontarget io.Writer
		if Config.apiDump {
			jsontarget = io.MultiWriter(gw, uncompressed)
		} else {
			jsontarget = gw
		}
		jw := json.NewEncoder(jsontarget)
		if err := jw.Encode(upload); err != nil {
			log.Fatalf("doRequest: JSON error encoding object to upload: %v", err)
		}
		if err := gw.Close(); err != nil {
			log.Fatalf("doRequest: gzip error encoding object to upload: %v", err)
		}
		req.Body = ioutil.NopCloser(payload)

		if Config.apiDump {
			log.Printf("Request data: %s", uncompressed)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error connecting to %s: %v", Config.Host, err)
	}
	defer resp.Body.Close()
	if notfoundokay && resp.StatusCode == http.StatusNotFound {
		return false
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("unexpected status from %s: %s", url, resp.Status)
		dumpBody(resp)
		log.Fatalf("giving up")
	}

	// parse the result if any
	if download != nil {
		body := resp.Body
		if resp.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(body)
			if err != nil {
				log.Fatalf("failed to decompress gzip result: %v", err)
			}
			body = gz
			defer gz.Close()
		}
		decoder := json.NewDecoder(body)
		if err := decoder.Decode(download); err != nil {
			log.Fatalf("failed to parse result object from server: %v", err)
		}

		if Config.apiDump {
			raw, err := json.MarshalIndent(download, "", "    ")
			if err != nil {
				log.Fatalf("doRequest: JSON error encoding downloaded object: %v", err)
			}
			log.Printf("Response data: %s", raw)
		}

		return true
	}
	return false
}

func courseDirectory(label string) string {
	re := regexp.MustCompile(`^([A-Za-z]+[- ]*\d+\w*)\b`)
	groups := re.FindStringSubmatch(label)
	if len(groups) == 2 {
		return groups[1]
	}
	return label
}

func mustLoadConfig(cmd *cobra.Command) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("unable to find home directory: %v", err)
	}
	if home == "" {
		log.Fatalf("home directory is not setn")
	}
	configFile := filepath.Join(home, perUserDotFile)

	if raw, err := ioutil.ReadFile(configFile); err != nil {
		log.Fatalf("Unable to load config file; try running '%s login'\n", os.Args[0])
	} else if err := json.Unmarshal(raw, &Config); err != nil {
		log.Printf("failed to parse %s: %v", configFile, err)
		log.Fatalf("you may wish to try deleting the file and running '%s login' again\n", os.Args[0])
	}
	if Config.apiDump {
		Config.apiReport = true
	}
}

func mustWriteConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("unable to find home directory: %v", err)
	}
	if home == "" {
		log.Fatalf("home directory is not set")
	}
	configFile := filepath.Join(home, perUserDotFile)

	raw, err := json.MarshalIndent(&Config, "", "    ")
	if err != nil {
		log.Fatalf("JSON error encoding cookie file: %v", err)
	}
	raw = append(raw, '\n')

	if err = ioutil.WriteFile(configFile, raw, 0644); err != nil {
		log.Fatalf("error writing %s: %v", configFile, err)
	}
}

func dumpBody(resp *http.Response) {
	if resp.Body == nil {
		return
	}

	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			log.Fatalf("failed to decompress gzip result: %v", err)
		}
		defer gz.Close()
		io.Copy(os.Stderr, gz)
	} else {
		io.Copy(os.Stderr, resp.Body)
	}
}
