package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

const (
	appName                       = "Easel"
	cmdName                       = "easel"
	perUserDotFile                = "." + cmdName + "rc"
	urlPrefix                     = "/api/v1"
	assignmentsPath               = coursePath + "/assignments"
	assignmentPath                = assignmentsPath + "/%d"
	coursesPath                   = "/courses"
	coursePath                    = coursesPath + "/%d"
	modulesPath                   = coursePath + "/modules"
	modulePath                    = modulesPath + "/%d"
	pagesPath                     = coursePath + "/pages"
	pagePath                      = pagesPath + "/%s"
	quizzesPath                   = coursePath + "/quizzes"
	quizPath                      = quizzesPath + "/%d"
	quizQuestionsPath             = quizPath + "/questions"
	quizSubmissionsPath           = quizPath + "/submissions"
	quizSubmissionPath            = quizSubmissionsPath + "/%d"
	quizSubmissionQuestionsPath   = "/quiz_submissions/%d/questions"
	quizReportsPath               = quizPath + "/reports"
	quizReportPath                = quizReportsPath + "/%d"
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

	// Login
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

	// Init
	cmdInit := &cobra.Command{
		Use:   "init",
		Short: "Initialize the db",
		Long:  "TODO instructions",
		Run:   CommandInit,
	}
	cmd.AddCommand(cmdInit)

	// Course
	cmdCourse := &cobra.Command{
		Use:   "course <command>",
		Short: "Manage courses",
		Long:  "TODO instructions",
	}
	cmdCourseAdd := &cobra.Command{
		Use:   "add <course_canvas_url>",
		Short: "Point db to the given course",
		Long:  "TODO instructions",
		Run:   CommandCourseAdd,
	}
	cmdCourse.AddCommand(cmdCourseAdd)
	cmdCourseList := &cobra.Command{
		Use:   "list",
		Short: "List current courses",
		Long:  "TODO instructions",
		Run:   CommandCourseList,
	}
	cmdCourse.AddCommand(cmdCourseList)
	cmdCourseRemove := &cobra.Command{
		Use:   "remove <section_number|course_canvas_id>",
		Short: "Remove course",
		Long:  "TODO instructions",
		Run:   CommandCourseRemove,
	}
	cmdCourse.AddCommand(cmdCourseRemove)
	cmd.AddCommand(cmdCourse)

	// Pull
	cmdPull := &cobra.Command{
		Use:   "pull [component_type] [component_id]",
		Short: "pull a single component or all of that type if blank",
		Long:  "TODO instructions",
		Run:   CommandPull,
	}
	cmd.AddCommand(cmdPull)

	// Push
	cmdPush := &cobra.Command{
		Use:   "push [component_type] [component_id]",
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
	mustLoadConfig()
	mustCreateDb()
	// TODO: create directory structure
}

func CommandCourseAdd(cmd *cobra.Command, args []string) {
	mustLoadConfig()
	db := findDb()
	defer db.Close()

	courseId, err := getCourseIdFromUrl(args[0])
	if err != nil {
		log.Fatal(err.Error())
	}

	if c, _ := findCourse(db, courseId); c.Id > 0 {
		log.Fatalf("Course exists")
	}

	course, err := pullCourse(db, courseId)
	if err != nil {
		log.Fatal(err.Error())
	}

	course.Save(db)
}

func CommandCourseList(cmd *cobra.Command, args []string) {
	db := findDb()
	defer db.Close()

	courses, err := findCourses(db)
	if err != nil {
		log.Fatal("Failed to find courses")
	}
	for _, course := range courses {
		fmt.Println(course)
	}
}

func CommandCourseRemove(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: %s remove <course_number>", os.Args[0])
	}

	db := findDb()
	defer db.Close()

	courses, err := matchCourse(db, args[0])
	if len(courses) == 0 || err != nil {
		log.Fatalf("Could not find course for %s", args[0])
	} else if len(courses) > 1 {
		for _, course := range courses {
			fmt.Println(course)
		}
		log.Printf("The search found more than one course")
		log.Printf("   pick the correct course id from the list")
		log.Fatalf("   and run '%s course remove <id>'", os.Args[0])
	}

	err = courses[0].Remove(db)
	if err != nil {
		log.Fatalf("Failed to remove course: %v", err)
	}
	fmt.Println("Removed course", courses[0].Name)
}

func CommandPull(cmd *cobra.Command, args []string) {
	mustLoadConfig()
	db := findDb()
	defer db.Close()

	switch len(args) {
	case 0:
		// TODO: pull all components of all types
		pullCourses(db)
	case 1:
		// pull all components of single type
		componentType := args[0]
		switch componentType {
		case "assignments":
			pullAssignments(db)
		case "courses":
			pullCourses(db)
		case "modules":
			pullModules(db)
		case "pages":
			pullPages(db)
		default:
			log.Fatalf("Invalid component type: %s", componentType)
		}
	case 2:
		// pull single item of single type
		componentType := args[0]
		componentFilepath := args[1]
		switch componentType {
		case "assignments", "assignment":
			assignment := new(Assignment)
			_, err := readFile(componentFilepath, assignment)
			if err != nil {
				log.Fatalf("Failed to load assignment from file %s\n", componentFilepath)
			}
			assignment.Pull(db)
		case "courses", "course":
			courses, err := matchCourse(db, componentFilepath)
			if err != nil || len(courses) != 1 {
				log.Fatalf("Failed to find a single course for %s. %v\n", componentFilepath, err)
			}
			pullCourse(db, courses[0].CanvasId)
		case "modules", "module":
			module := new(Module)
			err := readYamlFile(componentFilepath, module)
			if err != nil {
				log.Fatalf("Failed to load module from file %s\n", componentFilepath)
			}
			module.Pull(db)
		case "pages", "page":
			page := new(Page)
			page.Url = getPageUrlFromFilepath(componentFilepath)
			page.Pull(db)
		default:
			log.Fatalf("Invalid component type: %s", componentType)
		}
	default:
		log.Fatal("Too many arguments")
	}
}

func CommandPush(cmd *cobra.Command, args []string) {
	mustLoadConfig()
	db := findDb()
	defer db.Close()

	switch len(args) {
	case 0:
		// push all components of all types
		pushCourses(db)
		pushPages(db)
	case 1:
		// push all components of single type
		componentType := args[0]
		switch componentType {
		case "courses":
			pushCourses(db)
		case "pages":
			pushPages(db)
		default:
			log.Fatalf("Invalid component type: %s", componentType)
		}
	case 2:
		// push single item of single type
		componentType := args[0]
		componentFilepath := args[1]
		switch componentType {
		case "courses":
			courses, err := matchCourse(db, componentFilepath)
			if err != nil {
				log.Fatalf("Error finding course %s\n", componentFilepath)
			}
			if len(courses) != 1 {
				for _, c := range courses {
					fmt.Println(c)
				}
				log.Fatalf("Failed to find a single match for %s. Narrow search query or select course id from list.")
			}
			courses[0].Push()
		case "pages", "page":
			pageUrl := getPageUrlFromFilepath(componentFilepath)
			// TODO: flag for notifying participants of update? set field notify_of_update
			pushPage(db, pageUrl)
		default:
			log.Fatalf("Invalid component type: %s", componentType)
		}
	default:
		log.Fatal("Too many arguments")
	}
}
