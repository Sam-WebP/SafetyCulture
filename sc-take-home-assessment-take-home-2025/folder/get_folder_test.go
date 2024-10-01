package folder_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	orgID1           = uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a")
	orgID2           = uuid.FromStringOrNil("f8a982ed-f17a-4dd9-99ca-ef05b6f5b17f")
	nonExistentOrgID = uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000")

	folders = []folder.Folder{
		{Name: "alpha", Paths: "alpha", OrgId: orgID1},
		{Name: "bravo", Paths: "alpha.bravo", OrgId: orgID1},
		{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: orgID1},
		{Name: "delta", Paths: "alpha.delta", OrgId: orgID1},
		{Name: "echo", Paths: "echo", OrgId: orgID1},
		{Name: "foxtrot", Paths: "foxtrot", OrgId: orgID2},
	}
)

func TestGetFoldersByOrgID(t *testing.T) {
	driver := folder.NewDriver(folders)

	tests := []struct {
		name  string
		orgID uuid.UUID
		want  []folder.Folder
	}{
		{
			name:  "Given an existing orgID, When GetFoldersByOrgID is called, Then it should return all folders for that org",
			orgID: orgID1,
			want: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgID1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgID1},
				{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: orgID1},
				{Name: "delta", Paths: "alpha.delta", OrgId: orgID1},
				{Name: "echo", Paths: "echo", OrgId: orgID1},
			},
		},
		{
			name:  "Given another different existing orgID, When GetFoldersByOrgID is called, Then it should return all folders for that org",
			orgID: orgID2,
			want: []folder.Folder{
				{Name: "foxtrot", Paths: "foxtrot", OrgId: orgID2},
			},
		},
		{
			name:  "Given a non-existent orgID, When GetFoldersByOrgID is called, Then it should return an empty slice",
			orgID: nonExistentOrgID,
			want:  []folder.Folder{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := driver.GetFoldersByOrgID(tt.orgID)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestGetAllChildFolders(t *testing.T) {
	driver := folder.NewDriver(folders)

	tests := []struct {
		name    string
		orgID   uuid.UUID
		folder  string
		want    []folder.Folder
		wantErr error
	}{
		{
			name:   "Given a root folder with children, When GetAllChildFolders is called, Then it should return the children",
			orgID:  orgID1,
			folder: "alpha",
			want: []folder.Folder{
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgID1},
				{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: orgID1},
				{Name: "delta", Paths: "alpha.delta", OrgId: orgID1},
			},
			wantErr: nil,
		},
		{
			name:   "Given a parent folder that is also a child, When GetAllChildFolders is called, Then it should return its children",
			orgID:  orgID1,
			folder: "bravo",
			want: []folder.Folder{
				{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: orgID1},
			},
			wantErr: nil,
		},
		{
			name:    "Given a child folder with no children, When GetAllChildFolders is called, Then it should return an empty slice",
			orgID:   orgID1,
			folder:  "charlie",
			want:    []folder.Folder{},
			wantErr: nil,
		},
		{
			name:    "Given a root folder with no children, When GetAllChildFolders is called, Then it should return an empty slice",
			orgID:   orgID1,
			folder:  "echo",
			want:    []folder.Folder{},
			wantErr: nil,
		},
		{
			name:    "Given an invalid folder, When GetAllChildFolders is called, Then it should return ErrFolderNotFound",
			orgID:   orgID1,
			folder:  "invalid_folder",
			want:    nil,
			wantErr: folder.ErrFolderNotFound,
		},
		{
			name:    "Given a folder not in the specified organization, When GetAllChildFolders is called, Then it should return ErrFolderNotInOrganization",
			orgID:   orgID1,
			folder:  "foxtrot",
			want:    nil,
			wantErr: folder.ErrFolderNotInOrganization,
		},
		{
			name:    "Given a root folder with no children that is from a second organization, When GetAllChildFolders is called, Then it should return an empty slice",
			orgID:   orgID2,
			folder:  "foxtrot",
			want:    []folder.Folder{},
			wantErr: nil,
		},
		{
			name:    "Given a folder and an incorrect orgID, When GetAllChildFolders is called, Then it should return ErrFolderNotInOrganization",
			orgID:   orgID2,
			folder:  "alpha",
			want:    nil,
			wantErr: folder.ErrFolderNotInOrganization,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := driver.GetAllChildFolders(tt.orgID, tt.folder)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.want, got)
			}
		})
	}
}
