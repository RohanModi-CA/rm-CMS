package versioning

import (
		"os/exec"
		"fmt"
		"time"
		"os"
		"path/filepath"
		"cms/parsers"
	   )


type file_to_sparse_push struct {
	relative_path string
	sparse_set_path string
}

func sparse_clone(files_datas []file_to_sparse_push, repository_url string) {
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
	var error_message string

	local_repo_parent_dir = "/tmp/RM-CMS-TMPDIR-" + string(time.Now().UnixNano()) + "/" 

	// Now we clone our repo within our temporary thing.
	clone_command_list = []string{"git", "clone", "--filter=blob:none" , "--depth=1", "--no-checkout", repository_url, local_repo_parent_dir}
	clone_cmd := exec.Command(clone_command_list[0], clone_command_list[1:]...)
	clone_err := clone_cmd.Run()
	
	if (clone_err != nil) {
		error_message = fmt.Sprintf("Error cloning! %s", clone_err)	
		panic(error_message)
	}

	// Make sure we clean up.
	defer os.RemoveAll(local_repo_parent_dir)

	// Now, we're going to sparse checkout. This is a two step process. First, we initiate the
	// sparse checkout with sparse-checkout init --cone, and then we do sparse-checkout set 
	// <repo-path to all the files>. We'll need to first get the raw filenames for each. We 
	// need to ensure we also call the git commands from within the repo directory.

	// We'll do this for an easier call of sparse-checkout.
	for i:=0; i<len(files_datas); i++ {
		repo_paths = append(repo_paths, files_datas[i].sparse_set_path)
	}

	sparse_checkout_init_list = []string{"git", "sparse-checkout", "init", "--cone"}
	sparse_checkout_init_cmd := exec.Command(sparse_checkout_init_list[0], sparse_checkout_init_list[1:]...)
	sparse_checkout_init_cmd.Dir = local_repo_parent_dir
	sparse_checkout_init_err := sparse_checkout_init_cmd.Run()

	if (sparse_checkout_init_err != nil) {
		error_message = fmt.Sprintf("Error initing sparse checkout: %s", sparse_checkout_init_err)
		panic(error_message)
	}

	sparse_checkout_set_list = append([]string{"git", "sparse-checkout", "set"}, repo_paths...)
	sparse_checkout_set_cmd := exec.Command(sparse_checkout_set_list[0], sparse_checkout_set_list[1:]...)
	sparse_checkout_set_cmd.Dir = local_repo_parent_dir
	sparse_checkout_set_err := sparse_checkout_set_cmd.Run()

	if (sparse_checkout_set_err != nil) {
		error_message = fmt.Sprintf("Error setting the sparse checkout: %s", sparse_checkout_set_err)
		panic(error_message)
	} 

	// Now we checkout the master branch.
	checkout_master_list = []string{"git", "checkout", "master"}
	checkout_master_cmd := exec.Command(checkout_master_list[0], checkout_master_list[1:]...)
	checkout_master_cmd.Dir = local_repo_parent_dir
	checkout_master_err := checkout_master_cmd.Run()

	if (checkout_master_err != nil) {
		error_message = fmt.Sprintf("Error checking out the master branch: %s", checkout_master_err)
		panic(error_message)
	}
	

	/* Now, we need to ensure the directory structures exist so that we can copy our files to where they
	   need to go. */
	
	for i:=0; i<len(files_datas); i++ {
		i_par_dir := local_repo_parent_dir + filepath.Dir(files_datas[i].sparse_set_path)
		create_subdirectory_list = []string{"mkdir", "-p", i_par_dir}
		create_subdirectory_cmd := exec.Command(create_subdirectory_list[0], create_subdirectory_list[1:]...)
		create_subdirectory_err := create_subdirectory_cmd.Run()

		if (create_subdirectory_err != nil) {
			error_message = fmt.Sprintf("Error creating the subdirectory structure: %s", create_subdirectory_err)
			panic(error_message)
		}
	}


	// Now, we need to move the files into the directory, and add them to git.
	for i:= 0; i<len(files_datas); i++ {
		// Move them

		move_files_to_tmp_list = []string{"cp", files_datas[i].relative_path, local_repo_parent_dir + files_datas[i].sparse_set_path}
		move_files_to_tmp_cmd := exec.Command(move_files_to_tmp_list[0], move_files_to_tmp_list[1:]...)
		move_files_to_tmp_err := move_files_to_tmp_cmd.Run()

		if (move_files_to_tmp_err != nil) {
			error_message = fmt.Sprintf("Error moving the files to the temp directory: %s", move_files_to_tmp_err)
			panic(error_message)
		}

		// Add them
		add_files_to_git_list = []string{"git", "add", files_datas[i].sparse_set_path}
		add_files_to_git_cmd := exec.Command(add_files_to_git_list[0], add_files_to_git_list[1:]...)	
		add_files_to_git_cmd.Dir = local_repo_parent_dir
		add_files_to_git_err := add_files_to_git_cmd.Run()

		if (add_files_to_git_err != nil) {
			error_message = fmt.Sprintf("Error adding the files to git: %s", add_files_to_git_err)
			panic(error_message)
		}
	}

	// We commit them all and push them all.
	commit_files_list = []string{"git", "commit", "-am", "File Upload, RM-CMS"}
	commit_files_cmd := exec.Command(commit_files_list[0], commit_files_list[1:]...)
	commit_files_cmd.Dir = local_repo_parent_dir
	commit_files_err := commit_files_cmd.Run()

	if (commit_files_err != nil) {
		error_message = fmt.Sprintf("Error commiting the files: %s", commit_files_err)
		panic(error_message)
	}

	push_files_list = []string{"git", "push", "origin", "master"}
	push_files_cmd := exec.Command(push_files_list[0], push_files_list[1:]...)
	push_files_cmd.Dir = local_repo_parent_dir
	push_files_err := push_files_cmd.Run()

	if (push_files_err != nil) {
		error_message = fmt.Sprintf("Error pushing the files: %s", push_files_err)
		panic(error_message)
	}
}



func PublishStatics(cs *parsers.ConversionState) {
	relative_paths := cs.ImagesRelativePaths

	const repository_url string = "git@github.com:RohanModi-CA/static.rohanmodi.ca.git"
	var files_datas []file_to_sparse_push

	for i:=0; i<len(relative_paths); i++ {
		i_file_to_sparse_push := file_to_sparse_push{relative_paths[i], "images/"+ filepath.Base(relative_paths[i])}
		files_datas = append(files_datas, i_file_to_sparse_push)
	}

	sparse_clone(files_datas, repository_url)
}
