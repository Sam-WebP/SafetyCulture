package main

import (
	"fmt"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
)

func main() {
	orgID := uuid.FromStringOrNil(folder.DefaultOrgID)

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
	folder.PrettyPrint(folders)

	// Print folders for the specific orgID
	fmt.Printf("\nFolders for orgID: %s\n", orgID)
	folder.PrettyPrint(orgFolder)
}
