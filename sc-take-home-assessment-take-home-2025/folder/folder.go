package folder

import (
	"errors"
	"strings"

	"github.com/gofrs/uuid"
)

// IDriver defines the interface for folder operations.
type IDriver interface {
	// GetFoldersByOrgID returns all folders that belong to a specific orgID.
	GetFoldersByOrgID(orgID uuid.UUID) []*Folder
	// GetAllChildFolders returns all child folders of a specific folder.
	GetAllChildFolders(orgID uuid.UUID, name string) ([]*Folder, error)
	// MoveFolder moves a folder to a new destination.
	MoveFolder(orgID uuid.UUID, name string, dst string) ([]*Folder, error)
}

// Custom errors
var (
	ErrFolderNotFound                  = errors.New("folder does not exist")
	ErrFolderNotInOrganization         = errors.New("folder does not exist in the specified organization")
	ErrCannotMoveAcrossOrganizations   = errors.New("cannot move a folder to a different organization")
	ErrCannotMoveFolderToItself        = errors.New("cannot move a folder to itself")
	ErrCannotMoveFolderToOwnDescendant = errors.New("cannot move a folder to a child of itself")
)

// driver implements the IDriver interface.
type driver struct {
	folderMap      map[string]*Folder               // key: lower(name)+"_"+orgID.String()
	foldersByOrgID map[uuid.UUID][]*Folder          // key: orgID
	nameIndex      map[string]map[uuid.UUID]*Folder // key: lower(name) -> orgID -> Folder
}

// NewDriver initializes a new driver with the given folders.
func NewDriver(folders []Folder) IDriver {
	d := &driver{
		folderMap:      make(map[string]*Folder),
		foldersByOrgID: make(map[uuid.UUID][]*Folder),
		nameIndex:      make(map[string]map[uuid.UUID]*Folder),
	}

	// Temporary map to hold folders by their path for tree building
	pathMap := make(map[string]*Folder)

	for i := range folders {
		folder := &folders[i]
		key := generateKey(folder.Name, folder.OrgId)
		d.folderMap[key] = folder

		d.foldersByOrgID[folder.OrgId] = append(d.foldersByOrgID[folder.OrgId], folder)

		lowerName := strings.ToLower(folder.Name)
		if _, exists := d.nameIndex[lowerName]; !exists {
			d.nameIndex[lowerName] = make(map[uuid.UUID]*Folder)
		}
		d.nameIndex[lowerName][folder.OrgId] = folder

		// Populate pathMap for tree building
		pathMap[folder.Paths] = folder
	}

	// Build the tree by assigning parents and children
	for _, folder := range pathMap {
		if folder.Paths == folder.Name {
			// This is a root folder; no parent
			continue
		}
		parentPath := getParentPath(folder.Paths)
		parentFolder, exists := pathMap[parentPath]
		if exists {
			folder.Parent = parentFolder
			parentFolder.Children = append(parentFolder.Children, folder)
		}
	}

	return d
}

// generateKey creates a unique key for a folder based on its name and organization ID.
func generateKey(name string, orgID uuid.UUID) string {
	return strings.ToLower(name) + "_" + orgID.String()
}

// getParentPath extracts the parent path from a given path string.
// For example, "alpha.bravo.charlie" returns "alpha.bravo"
func getParentPath(path string) string {
	lastDot := strings.LastIndex(path, ".")
	if lastDot == -1 {
		return path
	}
	return path[:lastDot]
}
