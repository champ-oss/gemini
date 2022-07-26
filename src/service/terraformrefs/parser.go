package terraformrefs

import (
	"fmt"
	"github.com/champ-oss/gemini/model"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

// parseStringAsHcl parses the string as an HCL file
func parseStringAsHcl(hclString string, fileName string) (*hcl.File, error) {
	log.Debugf("Parsing %s contents as HCL file", fileName)
	parser := hclparse.NewParser()
	mainHcl, diag := parser.ParseHCL([]byte(hclString), fileName)
	if diag.HasErrors() {
		return nil, diag
	}
	return mainHcl, nil
}

// parseModuleCalls parses the Terraform file and gets a list of module call ref information
func parseModuleCalls(hclFile *hcl.File, run *model.WorkflowRun) ([]*model.TerraformRef, error) {
	log.Info("Parsing HCL file as Terraform file and getting module calls")
	var refs []*model.TerraformRef

	mod := tfconfig.NewModule(".")
	diag := tfconfig.LoadModuleFromFile(hclFile, mod)
	if diag.HasErrors() {
		return refs, diag
	}

	for _, call := range mod.ModuleCalls {
		refs = append(refs, parseModuleCall(run, terraformFile, call))
	}
	return refs, nil
}

// parseModuleCall parse a terraform module call into a TerraformRef model
func parseModuleCall(run *model.WorkflowRun, fileName string, call *tfconfig.ModuleCall) *model.TerraformRef {
	ref := &model.TerraformRef{
		Owner:        run.Owner,
		Repo:         run.Repo,
		Branch:       run.Branch,
		Sha:          run.Sha,
		FileName:     fileName,
		RunUpdatedAt: run.RunUpdatedAt,
		ModuleName:   call.Name,
		SourceOwner:  parseOwner(call.Source),
		SourceRepo:   parseRepo(call.Source),
		SourceRef:    parseRef(call.Source),
	}
	return ref
}

// parseRef parses the Terraform source string and returns the "ref" value
// Example: "...foo?ref=main" would return "main"
func parseRef(source string) string {
	parsed, err := url.Parse(source)
	if err != nil {
		fmt.Printf("Unable to URL parse string: %s\n", source)
		panic(err)
	}
	return parsed.Query().Get("ref")
}

// parseRepo parses the Terraform source string and returns the name of the module or repo
// Example: "foo.org/company1/module1" would return "module1"
func parseRepo(source string) string {
	splitQuery := strings.Split(source, "?")

	// Get the last part of the path
	splitPath := strings.Split(splitQuery[0], "/")
	repo := splitPath[len(splitPath)-1]

	repo = stripString(repo, ".git")
	repo = stripString(repo, ".zip")

	return repo
}

// parseOwner parses the Terraform source string and returns the name of the module owner
// Example: "foo.org/company1/module1" would return "company1"
func parseOwner(source string) string {
	owner := ""

	// Try to get the 2nd to last part of the source path
	splitPath := strings.Split(source, "/")
	if len(splitPath) > 1 {
		owner = splitPath[len(splitPath)-2]
	}

	// Remove the prefix for some cases like: "git::git@github.com:owner1/"
	splitColon := strings.Split(owner, ":")
	owner = splitColon[len(splitColon)-1]
	return owner
}

// stripString removes the specified characters from the string
func stripString(s string, strip string) string {
	return strings.Replace(s, strip, "", -1)
}
