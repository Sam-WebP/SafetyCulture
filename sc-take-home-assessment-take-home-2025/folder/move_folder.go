package folder

import (
	"errors"
	"strings"

	"github.com/gofrs/uuid"
)

func (d *driver) MoveFolder(orgID uuid.UUID, name string, dst string) ([]*Folder, error) {
	sourceFolder, err := d.getFolderByNameAndOrgID(name, orgID)
	if err != nil {
		return nil, err
	}

	destinationFolder, err := d.getFolderByNameAndOrgID(dst, orgID)
	if err != nil {
		if errors.Is(err, ErrFolderNotInOrganization) {
			// Destination folder exists but in a different organization
			return nil, ErrCannotMoveAcrossOrganizations
		}
		return nil, err
	}

	if sourceFolder == destinationFolder {
		return nil, ErrCannotMoveFolderToItself
	}

	// Check if destination is a child of source
	if isDescendant(sourceFolder, destinationFolder) {
		return nil, ErrCannotMoveFolderToOwnDescendant
	}

	// Update parent and children relationships
	if sourceFolder.Parent != nil {
		sourceFolder.Parent.Children = removeChild(sourceFolder.Parent.Children, sourceFolder)
	}
	sourceFolder.Parent = destinationFolder
	destinationFolder.Children = append(destinationFolder.Children, sourceFolder)

	// Update paths recursively
	newPath := destinationFolder.Paths + "." + sourceFolder.Name
	updatePathsRecursive(d, sourceFolder, newPath)

	return d.GetFoldersByOrgID(orgID), nil
}

func removeChild(children []*Folder, childToRemove *Folder) []*Folder {
	for i, child := range children {
		if child == childToRemove {
			return append(children[:i], children[i+1:]...)
		}
	}
	return children
}

func updatePathsRecursive(d *driver, folder *Folder, newPath string) {
	folder.Paths = newPath
	for _, child := range folder.Children {
		childNewPath := newPath + "." + child.Name
		updatePathsRecursive(d, child, childNewPath)
	}
}

func (d *driver) findFolderByNameAndOrgID(name string, orgID uuid.UUID) (*Folder, uuid.UUID, error) {
	lowerName := strings.ToLower(name)
	if orgMap, exists := d.nameIndex[lowerName]; exists {
		if orgID == uuid.Nil {
			// Return any folder with this name
			for oid, folder := range orgMap {
				return folder, oid, nil
			}
			return nil, uuid.Nil, ErrFolderNotFound
		} else if folder, exists := orgMap[orgID]; exists {
			return folder, orgID, nil
		} else {
			// Folder exists with this name but not in the specified orgID
			return nil, uuid.Nil, ErrFolderNotInOrganization
		}
	}
	return nil, uuid.Nil, ErrFolderNotFound
}

func isDescendant(folder *Folder, potentialDescendant *Folder) bool {
	if potentialDescendant == nil {
		return false
	}
	if potentialDescendant.Parent == folder {
		return true
	}
	return isDescendant(folder, potentialDescendant.Parent)
}
