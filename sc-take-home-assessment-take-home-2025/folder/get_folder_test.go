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
)

func initGetFolderDriver() folder.IDriver {

	sampleFolders := []folder.Folder{
		{Name: "alpha", Paths: "alpha", OrgId: orgID1},
		{Name: "bravo", Paths: "alpha.bravo", OrgId: orgID1},
		{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: orgID1},
		{Name: "delta", Paths: "alpha.delta", OrgId: orgID1},
		{Name: "echo", Paths: "echo", OrgId: orgID1},
		{Name: "foxtrot", Paths: "foxtrot", OrgId: orgID2},
	}
	return folder.NewDriver(sampleFolders)
}

func TestGetFoldersByOrgID(t *testing.T) {
	driver := initGetFolderDriver()

	tests := []struct {
		name  string
		orgID uuid.UUID
		want  []folder.Folder
	}{
		{
			name:  "Existing orgID with folders",
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
			name:  "Another existing orgID with folders",
			orgID: orgID2,
			want: []folder.Folder{
				{Name: "foxtrot", Paths: "foxtrot", OrgId: orgID2},
			},
		},
		{
			name:  "Non-existent orgID",
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

func TestGetAllChildFolders_Success(t *testing.T) {
	driver := initGetFolderDriver()

	tests := []struct {
		name   string
		orgID  uuid.UUID
		folder string
		want   []folder.Folder
	}{
		{
			name:   "Root folder with children",
			orgID:  orgID1,
			folder: "alpha",
			want: []folder.Folder{
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgID1},
				{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: orgID1},
				{Name: "delta", Paths: "alpha.delta", OrgId: orgID1},
			},
		},
		{
			name:   "Parent folder that is also a child",
			orgID:  orgID1,
			folder: "bravo",
			want: []folder.Folder{
				{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: orgID1},
			},
		},
		{
			name:   "Child folder with no children",
			orgID:  orgID1,
			folder: "charlie",
			want:   []folder.Folder{},
		},
		{
			name:   "Root folder with no children",
			orgID:  orgID1,
			folder: "echo",
			want:   []folder.Folder{},
		},
		{
			name:   "Root folder in second organization",
			orgID:  orgID2,
			folder: "foxtrot",
			want:   []folder.Folder{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := driver.GetAllChildFolders(tt.orgID, tt.folder)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestGetAllChildFolders_Errors(t *testing.T) {
	driver := initGetFolderDriver()

	tests := []struct {
		name    string
		orgID   uuid.UUID
		folder  string
		wantErr error
	}{
		{
			name:    "Invalid folder",
			orgID:   orgID1,
			folder:  "invalid_folder",
			wantErr: folder.ErrFolderNotFound,
		},
		{
			name:    "Folder not in specified organization",
			orgID:   orgID1,
			folder:  "foxtrot",
			wantErr: folder.ErrFolderNotInOrganization,
		},
		{
			name:    "Incorrect orgID for existing folder",
			orgID:   orgID2,
			folder:  "alpha",
			wantErr: folder.ErrFolderNotInOrganization,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := driver.GetAllChildFolders(tt.orgID, tt.folder)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Nil(t, got)
		})
	}
}
