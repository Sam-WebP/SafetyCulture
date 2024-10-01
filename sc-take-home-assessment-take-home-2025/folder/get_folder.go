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
		return folders
	}
	return nil
}

func (d *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
	// Find the parent folder
	parentFolder, err := d.getFolderByNameAndOrgID(name, orgID)
	if err != nil {
		return nil, err
	}

	// Find and store the child folders
	prefix := parentFolder.Paths + "."
	childFolders := []Folder{}
	for _, folder := range d.foldersByOrgID[orgID] {
		if strings.HasPrefix(folder.Paths, prefix) {
			childFolders = append(childFolders, folder)
		}
	}

	return childFolders, nil
}

func (d *driver) getFolderByNameAndOrgID(name string, orgID uuid.UUID) (*Folder, error) {
	key := generateKey(name, orgID)
	if folder, exists := d.folderMap[key]; exists {
		return &folder, nil
	}

	// Check if the folder exists in any organization.
	for _, folder := range d.folderMap {
		if strings.EqualFold(folder.Name, name) && folder.OrgId != orgID {
			return nil, ErrFolderNotInOrganization
		}
	}
	return nil, ErrFolderNotFound
}
