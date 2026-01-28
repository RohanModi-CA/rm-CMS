package versioning

import (
	"cms/misc"
	"os"
	"time"
)

func GetWebsiteTree(CS *misc.ConversionState) {
	/*
		We're going to try to clone just the structure to get it.

	*/

	const repositoryURL string = "git@github.com:RohanModi-CA/rohanmodi.ca.git"

	var localRepoParentDir string
	var cloneCommandList []string
	var lsTreeCommandList []string
	var tree_command_list []string

	localRepoParentDir = "/tmp/RM-CMS-TMPDIR-" + string(time.Now().UnixNano()) + "/"
	cloneCommandList = []string{"git", "clone", "--filter=blob:none", "--no-checkout", repositoryURL, localRepoParentDir}
	cloneCommandOptions := misc.ExecOptions{
		CommandList: cloneCommandList,
		LogCommand:  false,
		Dir:         "",
		DontPanic:   false,
		CS:          CS,
		Stdin:       "",
	}

	misc.RunExecCommand(cloneCommandOptions)

	// Clean up
	defer os.RemoveAll(localRepoParentDir)

	lsTreeCommandList = []string{"git", "ls-tree", "-r", "HEAD", "--name-only"}
	ls_tree_command_options := misc.ExecOptions{
		CommandList: lsTreeCommandList,
		LogCommand:  false,
		Dir:         localRepoParentDir,
		DontPanic:   false,
		CS:          CS,
		Stdin:       "",
	}
	ls_out := misc.RunExecCommandOut(ls_tree_command_options)

	tree_command_list = []string{"tree", "--fromfile", "-J"}
	tree_command_options := misc.ExecOptions{
		CommandList: tree_command_list,
		LogCommand:  false,
		Dir:         localRepoParentDir,
		DontPanic:   false,
		CS:          CS,
		Stdin:       ls_out,
	}

	tree_json_out := misc.RunExecCommandOut(tree_command_options)
	_ = tree_json_out
	//print(tree_json_out)
}
