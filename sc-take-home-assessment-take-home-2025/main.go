package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Choose an option:")
		fmt.Println("1) Run all tests")
		fmt.Println("2) Run get_folder.go tests")
		fmt.Println("3) Run move_folder.go tests")
		fmt.Println("4) Print sample folder data")
		fmt.Println("5) Exit")
		fmt.Print("Enter choice: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		choice := strings.TrimSpace(input)

		var cmd *exec.Cmd

		switch choice {
		case "1":
			fmt.Println("Running all tests...")
			cmd = exec.Command("go", "test", "./...", "-v")
		case "2":
			fmt.Println("Running get_folder.go tests...")
			cmd = exec.Command("go", "test", "./folder", "-run", "^TestGet", "-v")
		case "3":
			fmt.Println("Running move_folder.go tests...")
			cmd = exec.Command("go", "test", "./folder", "-run", "^TestMove", "-v")
		case "4":
			fmt.Println("Printing sample folder data...")
			printSampleData()
			continue
		case "5":
			fmt.Println("Exiting.")
			return
		default:
			fmt.Println("Invalid choice. Please enter 1, 2, 3, 4, or 5.")
			continue
		}

		// Set the command's standard output and error to the program's output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Run the command
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error running tests: %v\n", err)
		}

		fmt.Println()
	}
}

// Function to print sample folder data
func printSampleData() {
	// Initialize sample orgID
	orgID := uuid.FromStringOrNil(folder.DefaultOrgID)

	// Fetch all folders
	res := folder.GetAllFolders() // [](*Folder)

	// Convert []*Folder to []Folder
	folders := make([]folder.Folder, len(res))
	for i, f := range res {
		folders[i] = *f
	}

	// Initialize the driver with []Folder
	folderDriver := folder.NewDriver(folders)
	orgFolder := folderDriver.GetFoldersByOrgID(orgID)

	// Print all folders
	fmt.Println("All Folders:")
	folder.PrettyPrint(folders)

	// Print folders for the specific orgID
	fmt.Printf("\nFolders for orgID: %s\n", orgID)
	folder.PrettyPrint(orgFolder)
}
