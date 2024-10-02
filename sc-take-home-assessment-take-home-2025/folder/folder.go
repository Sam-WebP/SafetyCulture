package folder

import (
	"errors"
	"strings"

	"github.com/gofrs/uuid"
)

type IDriver interface {
	// GetFoldersByOrgID returns all folders that belong to a specific orgID.
	GetFoldersByOrgID(orgID uuid.UUID) []Folder
	// component 1
	// Implement the following methods:
	// GetAllChildFolders returns all child folders of a specific folder.
	GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error)

	// component 2
	// Implement the following methods:
	// MoveFolder moves a folder to a new destination.
	MoveFolder(name string, dst string) ([]Folder, error)
}

// Custom errors
var (
	ErrFolderNotFound                  = errors.New("folder does not exist")
	ErrFolderNotInOrganization         = errors.New("folder does not exist in the specified organization")
	ErrCannotMoveAcrossOrganizations   = errors.New("cannot move a folder to a different organization")
	ErrCannotMoveFolderToItself        = errors.New("cannot move a folder to itself")
	ErrCannotMoveFolderToOwnDescendant = errors.New("cannot move a folder to a child of itself")
)

type driver struct {
	folderMap      map[string]*Folder      // key: lower(name)+"_"+orgID.String()
	foldersByOrgID map[uuid.UUID][]*Folder // key: orgID
	nameIndex      map[string][]*Folder    // key: lower(name)
}

func NewDriver(folders []Folder) IDriver {
	d := &driver{
		folderMap:      make(map[string]*Folder),
		foldersByOrgID: make(map[uuid.UUID][]*Folder),
		nameIndex:      make(map[string][]*Folder),
	}

	for i := range folders {
		folder := &folders[i]
		key := generateKey(folder.Name, folder.OrgId)
		d.folderMap[key] = folder

		d.foldersByOrgID[folder.OrgId] = append(d.foldersByOrgID[folder.OrgId], folder)

		lowerName := strings.ToLower(folder.Name)
		d.nameIndex[lowerName] = append(d.nameIndex[lowerName], folder)
	}

	return d
}

func generateKey(name string, orgID uuid.UUID) string {
	return strings.ToLower(name) + "_" + orgID.String()
}
