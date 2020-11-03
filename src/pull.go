package main

import (
	"github.com/spf13/cobra"
)

func CommandPull(cmd *cobra.Command, args []string) {
	mustLoadConfig(cmd)

	component := "all"
	componentId := "all"
	if len(args) > 1 {
		component = args[0]
		if len(args) == 2 {
			componentId = args[1]
		}
	}

	if component == "all" {
		pullAll(componentId)
	} else if component == "assignments" {
		pullAssignments(componentId)
	} else if component == "modules" {
		pullModules(componentId)
	} else if component == "pages" {
		pullPages(componentId)
	} else if component == "quizzes" {
		pullQuizzes(componentId)
	}
}

func pullAll(componentId string) {
	pullAssignments(componentId)
	pullModules(componentId)
	pullPages(componentId)
	pullQuizzes(componentId)
}
