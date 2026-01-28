package versioning

import (
	//	"os/exec"
	//	"fmt"

	"os"
	"path/filepath"
	"time"

	//"cms/parsers"

	"cms/misc"
)

type file_to_sparse_push struct {
	relative_path   string
	sparse_set_path string
}

func sparse_clone(files_datas []file_to_sparse_push, repository_url string, CS *misc.ConversionState) {
	/* This function sparse clones a git repository into a /tmp/ directory,
	   avoiding the downloading of all files except the ones references as
	   relative_path's in the files_datas struct. The files are copied to the
	   git repository from their relative_path's, committed, and then pushed.
	   The temporary directory is then deleted.
	*/

	var local_repo_parent_dir string
	var clone_command_list []string
	var repo_paths []string
	var sparse_checkout_init_list []string
	var sparse_checkout_set_list []string
	var checkout_master_list []string
	var create_subdirectory_list []string
	var move_files_to_tmp_list []string
	var add_files_to_git_list []string
	var commit_files_list []string
	var push_files_list []string
	//var error_message string

	local_repo_parent_dir = "/tmp/RM-CMS-TMPDIR-" + string(time.Now().UnixNano()) + "/"

	// Now we clone our repo within our temporary thing.
	clone_command_list = []string{"git", "clone", "--filter=blob:none", "--depth=1", "--no-checkout", repository_url, local_repo_parent_dir}
	clone_command_options := misc.ExecOptions{
		CommandList: clone_command_list,
		LogCommand:  false,
		Dir:         "",
		DontPanic:   false,
		CS:          CS,
		Stdin:       "",
	} // Command, Print it, Dir, CS.
	misc.RunExecCommand(clone_command_options)

	// Make sure we clean up.
	defer os.RemoveAll(local_repo_parent_dir)

	// Now, we're going to sparse checkout. This is a two step process. First, we initiate the
	// sparse checkout with sparse-checkout init --cone, and then we do sparse-checkout set
	// <repo-path to all the files>. We'll need to first get the raw filenames for each. We
	// need to ensure we also call the git commands from within the repo directory.

	// We'll do this for an easier call of sparse-checkout.
	for i := 0; i < len(files_datas); i++ {
		repo_paths = append(repo_paths, files_datas[i].sparse_set_path)
	}

	sparse_checkout_init_list = []string{"git", "sparse-checkout", "init", "--cone"}
	sparse_checkout_init_options := misc.ExecOptions{
		CommandList: sparse_checkout_init_list,
		LogCommand:  false,
		Dir:         local_repo_parent_dir,
		DontPanic:   false,
		CS:          CS,
		Stdin:       "",
	}
	misc.RunExecCommand(sparse_checkout_init_options)

	sparse_checkout_set_list = append([]string{"git", "sparse-checkout", "set"}, repo_paths...)
	sparse_checkout_set_options := misc.ExecOptions{
		CommandList: sparse_checkout_set_list,
		LogCommand:  false,
		Dir:         local_repo_parent_dir,
		DontPanic:   false,
		CS:          CS,
		Stdin:       "",
	}
	misc.RunExecCommand(sparse_checkout_set_options)

	// Now we checkout the master branch.
	checkout_master_list = []string{"git", "checkout", "master"}
	checkout_master_options := misc.ExecOptions{
		CommandList: checkout_master_list,
		LogCommand:  false,
		Dir:         local_repo_parent_dir,
		DontPanic:   false,
		CS:          CS,
		Stdin:       "",
	}
	misc.RunExecCommand(checkout_master_options)

	/* Now, we need to ensure the directory structures exist so that we can copy our files to where they
	   need to go. */

	for i := 0; i < len(files_datas); i++ {
		i_par_dir := local_repo_parent_dir + filepath.Dir(files_datas[i].sparse_set_path)
		create_subdirectory_list = []string{"mkdir", "-p", i_par_dir}
		create_subdirectory_options := misc.ExecOptions{
			CommandList: create_subdirectory_list,
			LogCommand:  false,
			Dir:         "",
			DontPanic:   false,
			CS:          CS,
			Stdin:       "",
		}
		misc.RunExecCommand(create_subdirectory_options)
	}

	// Now, we need to move the files into the directory, and add them to git.
	for i := 0; i < len(files_datas); i++ {
		// Move them

		move_files_to_tmp_list = []string{"cp", files_datas[i].relative_path, local_repo_parent_dir + files_datas[i].sparse_set_path}
		move_files_to_tmp_options := misc.ExecOptions{
			CommandList: move_files_to_tmp_list,
			LogCommand:  false,
			Dir:         "",
			DontPanic:   false,
			CS:          CS,
			Stdin:       "",
		}
		misc.RunExecCommand(move_files_to_tmp_options)

		// Add them
		add_files_to_git_list = []string{"git", "add", files_datas[i].sparse_set_path}
		add_files_to_git_options := misc.ExecOptions{
			CommandList: add_files_to_git_list,
			LogCommand:  false,
			Dir:         local_repo_parent_dir,
			DontPanic:   false,
			CS:          CS,
			Stdin:       "",
		}
		misc.RunExecCommand(add_files_to_git_options)

	}

	// We commit them all and push them all.
	commit_files_list = []string{"git", "commit", "-am", "File Upload, RM-CMS"}
	commit_files_options := misc.ExecOptions{
		CommandList: commit_files_list,
		LogCommand:  false,
		Dir:         local_repo_parent_dir,
		DontPanic:   true,
		CS:          CS,
		Stdin:       "",
	}
	misc.RunExecCommand(commit_files_options)

	push_files_list = []string{"git", "push", "origin", "master"}
	push_files_options := misc.ExecOptions{
		CommandList: push_files_list,
		LogCommand:  false,
		Dir:         local_repo_parent_dir,
		DontPanic:   false,
		CS:          CS,
		Stdin:       "",
	}

	misc.RunExecCommand(push_files_options)
}

func PublishStatics(CS *misc.ConversionState) {
	relative_paths := CS.ImagesRelativePaths

	const repository_url string = "git@github.com:RohanModi-CA/static.rohanmodi.ca.git"
	var files_datas []file_to_sparse_push

	for i := 0; i < len(relative_paths); i++ {
		i_file_to_sparse_push := file_to_sparse_push{relative_paths[i], "images/" + filepath.Base(relative_paths[i])}
		files_datas = append(files_datas, i_file_to_sparse_push)
	}

	sparse_clone(files_datas, repository_url, CS)
}

func PublishWebsite(CS *misc.ConversionState) {
	/* We will publish both the HTML file, and the associated .md file.
	   We have them both as strings, so we first need to convert them to files.
	*/

	const repository_url string = "git@github.com:RohanModi-CA/rohanmodi.ca.git"

	temp_html, temp_html_err := os.CreateTemp("", "temphtml-")
	temp_md, temp_md_err := os.CreateTemp("", "tempmd-")

	misc.ErrorHandlePanic(temp_html_err)
	misc.ErrorHandlePanic(temp_md_err)

	print(temp_html.Name())

	// Populate the files
	write_html_err := os.WriteFile(temp_html.Name(), []byte(CS.HtmlFileContents), 0644)
	write_md_err := os.WriteFile(temp_md.Name(), []byte(CS.MdFileContents), 0644)

	misc.ErrorHandlePanic(write_html_err)
	misc.ErrorHandlePanic(write_md_err)

	defer temp_html.Close()
	defer temp_md.Close()
	defer os.Remove(temp_html.Name())
	defer os.Remove(temp_md.Name())

	var html_and_md [2]file_to_sparse_push
	html_and_md[0].relative_path = temp_html.Name()
	html_and_md[1].relative_path = temp_md.Name()
	html_and_md[0].sparse_set_path = CS.WebsiteRelativePath + ".html"
	html_and_md[1].sparse_set_path = CS.WebsiteRelativePath + ".md"

	sparse_clone(html_and_md[:], repository_url, CS)
}
