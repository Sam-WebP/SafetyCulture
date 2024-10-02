package folder

import (
	"strings"

	"github.com/gofrs/uuid"
)

func (d *driver) MoveFolder(name string, dst string) ([]Folder, error) {
	// Find source folder
	sourceFolder, sourceOrgID, err := d.findFolderByName(name)
	if err != nil {
		return nil, err
	}

	// Find destination folder
	destinationFolder, destinationOrgID, err := d.findFolderByName(dst)
	if err != nil {
		return nil, err
	}

	// Ensure both folders are in the same organization
	if sourceOrgID != destinationOrgID {
		return nil, ErrCannotMoveAcrossOrganizations
	}

	// Prevent moving a folder to itself
	if sourceFolder.Paths == destinationFolder.Paths && sourceFolder.Name == destinationFolder.Name {
		return nil, ErrCannotMoveFolderToItself
	}

	// Prevent moving a folder to its own descendant
	if strings.HasPrefix(destinationFolder.Paths, sourceFolder.Paths+".") {
		return nil, ErrCannotMoveFolderToOwnDescendant
	}

	// Prepare new path for the source folder
	oldPath := sourceFolder.Paths
	newPath := destinationFolder.Paths + "." + sourceFolder.Name

	orgID := sourceOrgID

	// Update the paths of source folder and its descendants
	for _, folder := range d.foldersByOrgID[orgID] {
		if folder.Paths == oldPath || strings.HasPrefix(folder.Paths, oldPath+".") {
			// Compute the new path
			relativePath := strings.TrimPrefix(folder.Paths, oldPath)
			relativePath = strings.TrimPrefix(relativePath, ".")
			folder.Paths = newPath
			if relativePath != "" {
				folder.Paths += "." + relativePath
			}

			// Update folderMap
			key := generateKey(folder.Name, orgID)
			d.folderMap[key] = folder
			// No need to update nameIndex as folder's name hasn't changed
		}
	}

	return d.GetFoldersByOrgID(orgID), nil
}

func (d *driver) findFolderByName(name string) (Folder, uuid.UUID, error) {
	lowerName := strings.ToLower(name)
	if folders, exists := d.nameIndex[lowerName]; exists && len(folders) > 0 {
		// Assuming one folder per name per organization
		return *folders[0], folders[0].OrgId, nil // Return first occurrence
	}
	return Folder{}, uuid.Nil, ErrFolderNotFound
}
