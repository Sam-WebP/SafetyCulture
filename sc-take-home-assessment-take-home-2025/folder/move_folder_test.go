package folder_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func initMoveFolderDriver() folder.IDriver {
	sampleFolders := []folder.Folder{
		{Name: "alpha", Paths: "alpha", OrgId: orgID1},
		{Name: "bravo", Paths: "alpha.bravo", OrgId: orgID1},
		{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: orgID1},
		{Name: "delta", Paths: "alpha.delta", OrgId: orgID1},
		{Name: "echo", Paths: "alpha.delta.echo", OrgId: orgID1},
		{Name: "foxtrot", Paths: "foxtrot", OrgId: orgID2},
		{Name: "golf", Paths: "golf", OrgId: orgID1},
	}
	return folder.NewDriver(sampleFolders)
}

func TestMoveFolder_Success(t *testing.T) {
	driver := initMoveFolderDriver()

	testCases := []struct {
		name          string
		orgID         uuid.UUID
		source        string
		destination   string
		expectedPaths map[string]string
	}{
		{
			name:        "Move 'bravo' under 'delta'",
			orgID:       orgID1,
			source:      "bravo",
			destination: "delta",
			expectedPaths: map[string]string{
				"bravo":   "alpha.delta.bravo",
				"charlie": "alpha.delta.bravo.charlie",
			},
		},
		{
			name:        "Move 'bravo' under 'golf'",
			orgID:       orgID1,
			source:      "bravo",
			destination: "golf",
			expectedPaths: map[string]string{
				"bravo":   "golf.bravo",
				"charlie": "golf.bravo.charlie",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			updatedFolders, err := driver.MoveFolder(tc.orgID, tc.source, tc.destination)
			assert.NoError(t, err)

			// Create a map of folder names to their paths for easier lookup
			folderPaths := make(map[string]string)
			for _, folder := range updatedFolders {
				folderPaths[folder.Name] = folder.Paths
			}

			for name, expectedPath := range tc.expectedPaths {
				actualPath, exists := folderPaths[name]
				assert.True(t, exists, "Folder %s not found in updated folders", name)
				assert.Equal(t, expectedPath, actualPath, "Folder %s path mismatch", name)
			}
		})
	}
}

func TestMoveFolder_Errors(t *testing.T) {
	driver := initMoveFolderDriver()

	testCases := []struct {
		name          string
		orgID         uuid.UUID
		source        string
		destination   string
		expectedError error
	}{
		{
			name:          "Cannot move folder to itself",
			orgID:         orgID1,
			source:        "bravo",
			destination:   "bravo",
			expectedError: folder.ErrCannotMoveFolderToItself,
		},
		{
			name:          "Cannot move folder to its descendant",
			orgID:         orgID1,
			source:        "bravo",
			destination:   "charlie",
			expectedError: folder.ErrCannotMoveFolderToOwnDescendant,
		},
		{
			name:          "Cannot move folder to different organization",
			orgID:         orgID1,
			source:        "bravo",
			destination:   "foxtrot",
			expectedError: folder.ErrCannotMoveAcrossOrganizations,
		},
		{
			name:          "Source folder does not exist",
			orgID:         orgID1,
			source:        "invalid_folder",
			destination:   "delta",
			expectedError: folder.ErrFolderNotFound,
		},
		{
			name:          "Destination folder does not exist",
			orgID:         orgID1,
			source:        "bravo",
			destination:   "invalid_folder",
			expectedError: folder.ErrFolderNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := driver.MoveFolder(tc.orgID, tc.source, tc.destination)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
