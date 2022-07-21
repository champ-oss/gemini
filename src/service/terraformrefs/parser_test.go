package terraformrefs

import (
	"github.com/champ-oss/gemini/model"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func Test_parseRepo(t *testing.T) {
	result := parseRepo("git::git@github.com:champ-oss/terraform-config-example.git?ref=9543d8c92e78b1ac4cac07f6ac927e37f5982a84")
	assert.Equal(t, "terraform-config-example", result)

	result = parseRepo("bitbucket.org/hashicorp/terraform-consul-aws")
	assert.Equal(t, "terraform-consul-aws", result)

	result = parseRepo("app.terraform.io/example-corp/k8s-cluster")
	assert.Equal(t, "k8s-cluster", result)

	result = parseRepo("git::ssh://username@example.com/storage.git")
	assert.Equal(t, "storage", result)

	result = parseRepo("git::ssh://username@example.com/module123/storage.git")
	assert.Equal(t, "storage", result)

	result = parseRepo("s3::https://s3-eu-west-1.amazonaws.com/examplecorp-terraform-modules/vpc.zip")
	assert.Equal(t, "vpc", result)

	result = parseRepo("https://example.com/vpc-module?archive=zip")
	assert.Equal(t, "vpc-module", result)

	result = parseRepo("https://example.com/module123/vpc-module?archive=zip")
	assert.Equal(t, "vpc-module", result)

	result = parseRepo("foo?archive=zip")
	assert.Equal(t, "foo", result)

	result = parseRepo("foo")
	assert.Equal(t, "foo", result)

	result = parseRepo("../../")
	assert.Equal(t, "", result)
}

func Test_parseOwner(t *testing.T) {
	result := parseOwner("git::git@github.com:champ-oss/terraform-config-example.git?ref=9543d8c92e78b1ac4cac07f6ac927e37f5982a84")
	assert.Equal(t, "champ-oss", result)

	result = parseOwner("bitbucket.org/hashicorp/terraform-consul-aws?ref=main")
	assert.Equal(t, "hashicorp", result)

	result = parseOwner("app.terraform.io/example-corp/k8s-cluster?ref=main")
	assert.Equal(t, "example-corp", result)

	result = parseOwner("git::ssh://username@example.com/storage.git?ref=main")
	assert.Equal(t, "username@example.com", result)

	result = parseOwner("git::ssh://username@example.com/module123/storage.git")
	assert.Equal(t, "module123", result)

	result = parseOwner("s3::https://s3-eu-west-1.amazonaws.com/examplecorp-terraform-modules/vpc.zip")
	assert.Equal(t, "examplecorp-terraform-modules", result)

	result = parseOwner("https://example.com/vpc-module?archive=zip")
	assert.Equal(t, "example.com", result)

	result = parseOwner("https://example.com/module123/vpc-module?archive=zip")
	assert.Equal(t, "module123", result)

	result = parseOwner("foo?archive=zip")
	assert.Equal(t, "", result)

	result = parseOwner("foo")
	assert.Equal(t, "", result)

	result = parseOwner("../../")
	assert.Equal(t, "..", result)
}

func Test_parseRef(t *testing.T) {
	result := parseRef("git::git@github.com:champ-oss/terraform-config-example.git?ref=9543d8c92e78b1ac4cac07f6ac927e37f5982a84")
	assert.Equal(t, "9543d8c92e78b1ac4cac07f6ac927e37f5982a84", result)

	result = parseRef("git::git@github.com:champ-oss/terraform-config-example.git?ref=main")
	assert.Equal(t, "main", result)

	result = parseRef("git::git@github.com:champ-oss/terraform-config-example.git?ref=v1.0.0")
	assert.Equal(t, "v1.0.0", result)

	result = parseRef("https://example.com/module123/vpc-module?archive=zip")
	assert.Equal(t, "", result)

	result = parseRef("git::ssh://username@example.com/storage.git")
	assert.Equal(t, "", result)

	result = parseRef("app.terraform.io/example-corp/k8s-cluster?ref=main")
	assert.Equal(t, "main", result)
}

func Test_parseStringAsHcl(t *testing.T) {
	input := `variable "foo" {}`
	result, err := parseStringAsHcl(input, "test.tf")
	assert.NoError(t, err)
	assert.Greater(t, len(result.Bytes), 1)

	input = ``
	result, err = parseStringAsHcl(input, "test.tf")
	assert.NoError(t, err)
	assert.Equal(t, len(result.Bytes), 0)

	input = `not valid`
	result, err = parseStringAsHcl(input, "test.tf")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func Test_parseModuleCalls(t *testing.T) {
	input := `
	module "module1" {
      source = "git::git@github.com:champ-oss/terraform-foo1.git?ref=242f9bc82dd93fe95a700d8d5f04063129917d9c"
      git    = var.git
    }

	module "module2" {
      source = "git::git@github.com:champ-oss/terraform-foo2.git?ref=v1.0.0"
      git    = var.git
    }
	`
	hcl, _ := parseStringAsHcl(input, "main.tf")

	run := &model.WorkflowRun{
		Owner:        "owner1",
		Repo:         "repo1",
		Branch:       "branch1",
		Sha:          "03a631016daebdbd49205c7efdcbe2bc050aeaa3",
		RunUpdatedAt: 1643753231,
	}

	expected := []*model.TerraformRef{
		{
			Owner:        run.Owner,
			Repo:         run.Repo,
			Branch:       run.Branch,
			Sha:          run.Sha,
			FileName:     "main.tf",
			RunUpdatedAt: run.RunUpdatedAt,
			ModuleName:   "module1",
			SourceOwner:  "champ-oss",
			SourceRepo:   "terraform-foo1",
			SourceRef:    "242f9bc82dd93fe95a700d8d5f04063129917d9c",
		},
		{
			Owner:        run.Owner,
			Repo:         run.Repo,
			Branch:       run.Branch,
			Sha:          run.Sha,
			FileName:     "main.tf",
			RunUpdatedAt: run.RunUpdatedAt,
			ModuleName:   "module2",
			SourceOwner:  "champ-oss",
			SourceRepo:   "terraform-foo2",
			SourceRef:    "v1.0.0",
		},
	}

	refs, err := parseModuleCalls(hcl, run)
	sort.Slice(refs, func(i, j int) bool { return refs[i].ModuleName < refs[j].ModuleName })
	assert.NoError(t, err)
	assert.Equal(t, expected, refs)
}
