package terraformrefs

import (
	"github.com/champ-oss/gemini/mocks/mock_adapter"
	"github.com/champ-oss/gemini/mocks/mock_repository"
	"github.com/champ-oss/gemini/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test_ProcessRepo_Success checks that ProcessRepo() is successful
func Test_ProcessRepo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockRepository(ctrl)
	gitClient := mock_adapter.NewMockGitClient(ctrl)

	owner := "owner1"
	name := "repo1"
	branch := "branch1"
	tfFile := "main.tf"

	runs := []*model.WorkflowRun{
		{
			Owner:        "owner1",
			Repo:         "repo1",
			Branch:       "branch1",
			Sha:          "03a631016daebdbd49205c7efdcbe2bc050aeaa3",
			RunUpdatedAt: 1643753231,
		},
	}

	contents := `
	module "module1" {
      source = "git::git@github.com:champ-oss/terraform-foo1.git?ref=242f9bc82dd93fe95a700d8d5f04063129917d9c"
      git    = var.git
    }
	`

	expected := []*model.TerraformRef{
		{
			Owner:        "owner1",
			Repo:         "repo1",
			Branch:       "branch1",
			Sha:          "03a631016daebdbd49205c7efdcbe2bc050aeaa3",
			FileName:     tfFile,
			RunUpdatedAt: 1643753231,
			ModuleName:   "module1",
			SourceOwner:  "champ-oss",
			SourceRepo:   "terraform-foo1",
			SourceRef:    "242f9bc82dd93fe95a700d8d5f04063129917d9c",
		},
	}

	repo.EXPECT().GetWorkflowRunsByName(owner, name, branch, "apply").Return(runs)
	gitClient.EXPECT().GetContents(owner, name, tfFile, "03a631016daebdbd49205c7efdcbe2bc050aeaa3").Return(contents, nil)

	repo.EXPECT().AddTerraformRefs(expected).Return(int64(1), nil)
	svc := service{repo: repo, gitClient: gitClient}
	assert.Nil(t, svc.ProcessRepo(owner, name, branch))
}
