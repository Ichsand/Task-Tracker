package main

import (
	"fmt"
	"os"
)

// main is the entry point of the application.
// It parses command-line arguments and dispatches to the appropriate function.
func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	// Load existing tasks from the JSON file.
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "add":
		addTask(&tasks, args)
	case "list":
		listTasks(tasks, args)
	case "update":
		updateTask(&tasks, args)
	case "delete":
		deleteTask(&tasks, args)
	case "done":
		updateTaskStatus(&tasks, args, StatusDone)
	case "progress":
		updateTaskStatus(&tasks, args, StatusInProgress)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}

	// Save the modified task list back to the file.
	if err := saveTasks(tasks); err != nil {
		fmt.Printf("Error saving tasks: %v\n", err)
		os.Exit(1)
	}
}
