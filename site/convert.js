document.addEventListener("DOMContentLoaded", () =>
{
	const file_manager_directory_input = document.getElementById("file-manager-directory-input");
	const file_manager_files_buttons_div = document.getElementById("file-manager-files-buttons");
	const markdown_editor_textarea = document.getElementById("markdown-editor-textarea");
	const preview_html_button = document.getElementById("preview-html-button");

	

	file_manager_directory_input.addEventListener("change", () => 
	{
		let files = file_manager_directory_input.files
		parse_file_change(files)
	});


	function parse_file_change(files)
	{
		// First, we need to find out the name of the directory.
		let sample_filepath = files[0]["webkitRelativePath"];
		let directory_name = sample_filepath.split("/")[0];

		let filetree = {directory_name: {}};

		for (let i=0; i<files.length; ++i)
		{
			let file = files[i];

			// Only add a button for markdown files.
			if (!(file.name.substr(-3) === ".md"))
			{
				continue;
			}

			console.log(file);
			let file_button = document.createElement("button");
			file_button.textContent = file.webkitRelativePath;
			file_button.addEventListener("click", async () => 
			{
				let filetext = await file.text();
				launch_markdown_editor(filetext);
			});

			file_manager_files_buttons_div.appendChild(file_button);
		}
	}


	function launch_markdown_editor(relativepath)
	{
		markdown_editor_textarea.classList.remove("hidden");
		preview_html_button.classList.remove("hidden");
		markdown_editor_textarea.value = relativepath;
	}


});


