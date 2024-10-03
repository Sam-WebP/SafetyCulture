package folder

import (
	"strings"

	"github.com/gofrs/uuid"
)

func GetAllFolders() []*Folder {
	data := GetSampleData()
	result := make([]*Folder, len(data))
	for i := range data {
		result[i] = &data[i]
	}
	return result
}

func (d *driver) GetFoldersByOrgID(orgID uuid.UUID) []*Folder {
	if folders, exists := d.foldersByOrgID[orgID]; exists {
		return folders
	}
	return []*Folder{}
}

func (d *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]*Folder, error) {
	parentFolder, err := d.getFolderByNameAndOrgID(name, orgID)
	if err != nil {
		return nil, err
	}

	childFolders := []*Folder{}
	d.traverseChildren(parentFolder, &childFolders)

	return childFolders, nil
}

func (d *driver) traverseChildren(folder *Folder, collection *[]*Folder) {
	for _, child := range folder.Children {
		*collection = append(*collection, child)
		d.traverseChildren(child, collection)
	}
}

func (d *driver) getFolderByNameAndOrgID(name string, orgID uuid.UUID) (*Folder, error) {
	lowerName := strings.ToLower(name)
	if orgMap, exists := d.nameIndex[lowerName]; exists {
		if folder, exists := orgMap[orgID]; exists {
			return folder, nil
		} else {
			// Folder exists with this name but not in the specified orgID
			return nil, ErrFolderNotInOrganization
		}
	}
	return nil, ErrFolderNotFound
}
