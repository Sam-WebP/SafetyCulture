package folder

import (
	"fmt"

	"github.com/gofrs/uuid"
)

func (d *driver) MoveFolder(orgID uuid.UUID, name string, dst string) ([]*Folder, error) {
	sourceFolder, err := d.getFolderByNameAndOrgID(name, orgID)
	if err != nil {
		return nil, fmt.Errorf("source folder error: %w", err)
	}

	destinationFolder, err := d.getFolderByNameAndOrgID(dst, orgID)
	if err != nil {
		if err == ErrFolderNotInOrganization {
			return nil, fmt.Errorf("destination folder error: %w", ErrCannotMoveAcrossOrganizations)
		}
		return nil, fmt.Errorf("destination folder error: %w", err)
	}

	if sourceFolder == destinationFolder {
		return nil, ErrCannotMoveFolderToItself
	}

	if isDescendant(sourceFolder, destinationFolder) {
		return nil, ErrCannotMoveFolderToOwnDescendant
	}

	// Update parent and children relationships
	if sourceFolder.Parent != nil {
		sourceFolder.Parent.Children = removeChild(sourceFolder.Parent.Children, sourceFolder)
	}
	sourceFolder.Parent = destinationFolder
	destinationFolder.Children = append(destinationFolder.Children, sourceFolder)

	// Update paths iteratively
	newPath := destinationFolder.Paths + "." + sourceFolder.Name
	updatePathsIterative(sourceFolder, newPath)

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

func updatePathsIterative(folder *Folder, newPath string) {
	stack := []struct {
		f      *Folder
		newPth string
	}{{folder, newPath}}

	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		curr.f.Paths = curr.newPth
		for _, child := range curr.f.Children {
			childNewPath := curr.newPth + "." + child.Name
			stack = append(stack, struct {
				f      *Folder
				newPth string
			}{child, childNewPath})
		}
	}
}

func isDescendant(folder *Folder, potentialDescendant *Folder) bool {
	current := potentialDescendant.Parent
	for current != nil {
		if current == folder {
			return true
		}
		current = current.Parent
	}
	return false
}
