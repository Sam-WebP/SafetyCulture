package folder_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
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
		source        string
		destination   string
		expectedPaths map[string]string
	}{
		{
			name:        "Move 'bravo' under 'delta'",
			source:      "bravo",
			destination: "delta",
			expectedPaths: map[string]string{
				"bravo":   "alpha.delta.bravo",
				"charlie": "alpha.delta.bravo.charlie",
			},
		},
		{
			name:        "Move 'bravo' under 'golf'",
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
			updatedFolders, err := driver.MoveFolder(tc.source, tc.destination)
			assert.NoError(t, err)

			for _, folder := range updatedFolders {
				if newPath, exists := tc.expectedPaths[folder.Name]; exists {
					assert.Equal(t, newPath, folder.Paths, "Folder %s path mismatch", folder.Name)
				}
			}
		})
	}
}

func TestMoveFolder_Errors(t *testing.T) {
	driver := initMoveFolderDriver()

	testCases := []struct {
		name          string
		source        string
		destination   string
		expectedError error
	}{
		{
			name:          "Cannot move folder to itself",
			source:        "bravo",
			destination:   "bravo",
			expectedError: folder.ErrCannotMoveFolderToItself,
		},
		{
			name:          "Cannot move folder to its descendant",
			source:        "bravo",
			destination:   "charlie",
			expectedError: folder.ErrCannotMoveFolderToOwnDescendant,
		},
		{
			name:          "Cannot move folder to different organization",
			source:        "bravo",
			destination:   "foxtrot",
			expectedError: folder.ErrCannotMoveAcrossOrganizations,
		},
		{
			name:          "Source folder does not exist",
			source:        "invalid_folder",
			destination:   "delta",
			expectedError: folder.ErrFolderNotFound,
		},
		{
			name:          "Destination folder does not exist",
			source:        "bravo",
			destination:   "invalid_folder",
			expectedError: folder.ErrFolderNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := driver.MoveFolder(tc.source, tc.destination)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
