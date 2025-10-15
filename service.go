package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

// --- Core Logic Functions ---

// addTask creates a new task and adds it to the list.
func addTask(tasks *[]Task, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: Missing task description for 'add' command.")
		printUsage()
		os.Exit(1)
	}
	description := args[0]

	// Generate a new unique ID.
	newID := 1
	if len(*tasks) > 0 {
		newID = (*tasks)[len(*tasks)-1].ID + 1
	}

	now := time.Now()
	task := Task{
		ID:          newID,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	*tasks = append(*tasks, task)
	fmt.Printf("✅ Added task %d: \"%s\"\n", newID, description)
}

// updateTask modifies the description of an existing task.
func updateTask(tasks *[]Task, args []string) {
	if len(args) < 2 {
		fmt.Println("Error: Missing ID and new description for 'update' command.")
		printUsage()
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Error: Invalid ID format '%s'. ID must be a number.\n", args[0])
		os.Exit(1)
	}

	newDescription := args[1]
	task, index := findTaskByID(*tasks, id)

	if task == nil {
		fmt.Printf("Error: Task with ID %d not found.\n", id)
		os.Exit(1)
	}

	(*tasks)[index].Description = newDescription
	(*tasks)[index].UpdatedAt = time.Now()
	fmt.Printf("🔄 Updated task %d.\n", id)
}

// deleteTask removes a task from the list.
func deleteTask(tasks *[]Task, args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Missing ID for 'delete' command.")
		printUsage()
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Error: Invalid ID format '%s'. ID must be a number.\n", args[0])
		os.Exit(1)
	}

	task, index := findTaskByID(*tasks, id)
	if task == nil {
		fmt.Printf("Error: Task with ID %d not found.\n", id)
		os.Exit(1)
	}

	// Remove the task from the slice.
	*tasks = append((*tasks)[:index], (*tasks)[index+1:]...)
	fmt.Printf("❌ Deleted task %d.\n", id)
}

// updateTaskStatus changes the status of an existing task.
func updateTaskStatus(tasks *[]Task, args []string, newStatus string) {
	if len(args) < 1 {
		fmt.Println("Error: Missing ID for status update command.")
		printUsage()
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Error: Invalid ID format '%s'. ID must be a number.\n", args[0])
		os.Exit(1)
	}

	task, index := findTaskByID(*tasks, id)
	if task == nil {
		fmt.Printf("Error: Task with ID %d not found.\n", id)
		os.Exit(1)
	}

	(*tasks)[index].Status = newStatus
	(*tasks)[index].UpdatedAt = time.Now()
	fmt.Printf("Updated status for task %d to '%s'.\n", id, newStatus)
}

// listTasks prints tasks based on an optional filter.
func listTasks(tasks []Task, args []string) {
	filter := "all"
	if len(args) > 0 {
		filter = args[0]
	}

	fmt.Println("--- Task List ---")
	foundTasks := false
	for _, task := range tasks {
		printTask := false
		switch filter {
		case "all":
			printTask = true
		case "done":
			if task.Status == StatusDone {
				printTask = true
			}
		case "todo":
			if task.Status == StatusTodo {
				printTask = true
			}
		case "in-progress":
			if task.Status == StatusInProgress {
				printTask = true
			}
		}

		if printTask {
			foundTasks = true
			statusIcon := "📝" // todo
			if task.Status == StatusInProgress {
				statusIcon = "⏳"
			} else if task.Status == StatusDone {
				statusIcon = "✅"
			}
			fmt.Printf("%s [%s] %d: %s\n", statusIcon, task.Status, task.ID, task.Description)
		}
	}

	if !foundTasks {
		fmt.Printf("No tasks found for filter '%s'.\n", filter)
	}
	fmt.Println("-----------------")
}

// --- File I/O and Helpers ---

// loadTasks reads the task list from the JSON file.
// If the file does not exist, it returns an empty list.
func loadTasks() ([]Task, error) {
	data, err := os.ReadFile(taskFile)
	if err != nil {
		// If the file doesn't exist, that's okay. Return an empty slice.
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// saveTasks writes the current task list to the JSON file.
func saveTasks(tasks []Task) error {
	// MarshalIndent provides nicely formatted JSON.
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(taskFile, data, 0644)
}

// findTaskByID searches for a task by its ID and returns it and its index.
func findTaskByID(tasks []Task, id int) (*Task, int) {
	for i, task := range tasks {
		if task.ID == id {
			return &tasks[i], i
		}
	}
	return nil, -1
}

// printUsage displays the help message with all available commands.
func printUsage() {
	fmt.Println("\n--- Go Task Tracker ---")
	fmt.Println("A simple command-line task management tool.")
	fmt.Println("\nUsage:")
	fmt.Println("  go run main.go <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  add \"<description>\"    Add a new task.")
	fmt.Println("  list [filter]          List tasks. Filters: all (default), todo, in-progress, done.")
	fmt.Println("  update <id> \"<desc>\"   Update a task's description.")
	fmt.Println("  delete <id>            Delete a task.")
	fmt.Println("  done <id>              Mark a task as done.")
	fmt.Println("  progress <id>          Mark a task as in-progress.")
	fmt.Println()
}
