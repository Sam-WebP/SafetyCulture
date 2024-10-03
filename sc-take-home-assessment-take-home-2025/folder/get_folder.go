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
	stack := []*Folder{folder}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, child := range current.Children {
			*collection = append(*collection, child)
			stack = append(stack, child)
		}
	}
}

func (d *driver) getFolderByNameAndOrgID(name string, orgID uuid.UUID) (*Folder, error) {
	lowerName := strings.ToLower(name)
	if orgMap, exists := d.nameIndex[lowerName]; exists {
		if folder, exists := orgMap[orgID]; exists {
			return folder, nil
		}
		return nil, ErrFolderNotInOrganization
	}
	return nil, ErrFolderNotFound
}
