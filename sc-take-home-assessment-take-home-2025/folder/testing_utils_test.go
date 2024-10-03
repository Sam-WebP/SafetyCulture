package folder_test

import (
	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
)

// Shared orgID variables for tests
var (
	orgID1           = uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a")
	orgID2           = uuid.FromStringOrNil("f8a982ed-f17a-4dd9-99ca-ef05b6f5b17f")
	nonExistentOrgID = uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000")
)

// Helper type to compare folder content
type folderComparable struct {
	Name  string
	OrgId uuid.UUID
	Paths string
}

// Helper function to convert []*folder.Folder to []folderComparable
func foldersToComparable(folders []*folder.Folder) []folderComparable {
	result := make([]folderComparable, len(folders))
	for i, f := range folders {
		result[i] = folderComparable{
			Name:  f.Name,
			OrgId: f.OrgId,
			Paths: f.Paths,
		}
	}
	return result
}
