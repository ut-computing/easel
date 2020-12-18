package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
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

	"github.com/deiu/linkparser"
	"gopkg.in/yaml.v2"
)

const (
	easelDb = ".easeldb"
)

func writeFile(filename, metadata, html string) error {
	text := html
	if metadata != "" {
		text = fmt.Sprintf("```\n%s```\n", metadata) + text
	}
	return ioutil.WriteFile(filename, []byte(text), 0644)
}

func mustCreateDb() {
	directory, err := filepath.Abs(".")
	if err != nil {
		log.Fatalf("error finding directory: %v", err)
	}

	_, err = os.Stat(filepath.Join(directory, easelDb))
	if err == nil {
		log.Fatal("Database already exists")
	}

	db, err := sql.Open("sqlite3", easelDb)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// TODO: create all tables (execute .sql file)
	mustCreateCoursesTable(db)
}

func findDb() *sql.DB {
	directory, err := filepath.Abs(".")
	if err != nil {
		log.Fatalf("error finding directory: %v", err)
	}

	dbName := easelDb
	stepDir := directory
	for {
		path := filepath.Join(directory, dbName)
		_, err := os.Stat(path)
		if err == nil {
			break
		}
		if !os.IsNotExist(err) {
			log.Fatalf("error searching for %s in %s: %v", dbName, directory, err)
		}

		// try moving up a directory
		stepDir = directory
		directory = filepath.Dir(directory)
		if directory == stepDir {
			log.Fatal("No database found.")
		}
	}

	dbName = directory + "/" + dbName
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func addWhereLike(where string, args []interface{}, label string, value string) (string, []interface{}) {
	if where == "" {
		where = " WHERE"
	} else {
		where += " AND"
	}
	args = append(args, "%"+strings.ToLower(value)+"%")

	// sqlite is set to use case insensitive LIKEs
	where += fmt.Sprintf(" %s LIKE ?", label)
	return where, args
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

func prepareRequest(reqUrl string, params url.Values, method string, upload interface{}, download interface{}) (*http.Request, error) {
	req, err := http.NewRequest(method, reqUrl, nil)
	if err != nil {
		return req, err
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
	return req, nil
}

func doRequest(path string, params url.Values, method string, upload interface{}, download interface{}, notfoundokay bool) bool {
	if !strings.HasPrefix(path, "/") {
		log.Panicf("doRequest path must start with /")
	}
	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" {
		log.Panicf("doRequest only recognizes GET, POST, PUT, and DELETE methods")
	}

	reqUrl := fmt.Sprintf("https://%s%s%s", Config.Host, urlPrefix, path)
	req, err := prepareRequest(reqUrl, params, method, upload, download)
	if err != nil {
		log.Fatalf("error creating http request: %v\n", err)
	}

	paginated := false // assume not paginated
	allResults := make([]map[string]interface{}, 0)
	for {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalf("error connecting to %s: %v", Config.Host, err)
		}
		defer resp.Body.Close()
		if notfoundokay && resp.StatusCode == http.StatusNotFound {
			return false
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("unexpected status from %s: %s", reqUrl, resp.Status)
			dumpBody(resp)
			log.Fatalf("giving up")
		}
		gzipped := resp.Header.Get("Content-Encoding") == "gzip"

		// grab next page url if paginated
		links := lh.ParseHeader(resp.Header.Get("Link"))
		next, ok := links["next"]["href"]
		if !ok {
			if !paginated {
				// if non-paginated response, we're done
				return parseResponse(resp.Body, gzipped, download)
			} else {
				// no more paginated results, grab last results and done
				partResults := make([]map[string]interface{}, 0)
				if parseResponse(resp.Body, gzipped, &partResults) {
					allResults = append(allResults, partResults...)
				}
				break
			}
		}

		// set up request for next page
		paginated = true // at this point, we're paginated, set flag on the first time
		// grab partial results
		partResults := make([]map[string]interface{}, 0)
		if parseResponse(resp.Body, gzipped, &partResults) {
			allResults = append(allResults, partResults...)
		}
		// prepare for next request
		req, err = prepareRequest(next, url.Values{}, method, upload, download)
		if err != nil {
			log.Fatalf("error creating http request: %v\n", err)
		}
	}

	// re-encode all results
	allJson, err := json.Marshal(allResults)
	if err != nil {
		return false
	}
	return json.Unmarshal(allJson, download) == nil
}

func parseResponse(body io.ReadCloser, gzipped bool, download interface{}) bool {
	// parse the result if any
	if download != nil {
		if gzipped {
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

func mustLoadConfig() {
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

// Reads the file's metadata into the given target struct and returns the
// contents as a string
func readFile(filename string, target interface{}) (string, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	fileParts := strings.Split(string(dat), "```")
	err = yaml.Unmarshal([]byte(fileParts[1]), target)
	return fileParts[2], err
}
