package folder

import (
	"strings"

	"github.com/gofrs/uuid"
)

func GetAllFolders() []Folder {
	return GetSampleData()
}

func (d *driver) GetFoldersByOrgID(orgID uuid.UUID) []Folder {
	if folders, exists := d.foldersByOrgID[orgID]; exists {
		// Convert []*Folder to []Folder
		result := make([]Folder, len(folders))
		for i, folderPtr := range folders {
			result[i] = *folderPtr
		}
		return result
	}
	return []Folder{}
}

func (d *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
	// Find the parent folder
	parentFolder, err := d.getFolderByNameAndOrgID(name, orgID)
	if err != nil {
		return nil, err
	}

	// Build the prefix for child paths
	prefix := parentFolder.Paths + "."
	childFolders := []Folder{}
	for _, folder := range d.foldersByOrgID[orgID] {
		if strings.HasPrefix(folder.Paths, prefix) {
			childFolders = append(childFolders, *folder)
		}
	}

	return childFolders, nil
}

func (d *driver) getFolderByNameAndOrgID(name string, orgID uuid.UUID) (*Folder, error) {
	lowerName := strings.ToLower(name)
	if folders, exists := d.nameIndex[lowerName]; exists {
		for _, folder := range folders {
			if folder.OrgId == orgID {
				return folder, nil
			}
		}
	}

	// If not found in the specified orgID, check if the folder exists in other organizations
	if folders, exists := d.nameIndex[lowerName]; exists {
		for _, folder := range folders {
			if folder.OrgId != orgID {
				return nil, ErrFolderNotInOrganization
			}
		}
	}

	return nil, ErrFolderNotFound
}
