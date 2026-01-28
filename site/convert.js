document.addEventListener("DOMContentLoaded", () =>
{
	const file_manager_directory_input = document.getElementById("file-manager-directory-input");
	const file_manager_files_buttons_div = document.getElementById("file-manager-files-buttons");
	const markdown_editor_textarea = document.getElementById("markdown-editor-textarea");
	const preview_html_button = document.getElementById("preview-html-button");
	const html_preview_iframe = document.getElementById("html-preview-iframe");
	const push_statics_button = document.getElementById("push-statics-button")
	const destination_input = document.getElementById("destination-input")
	const destination_label = document.getElementById("destination-label")
	const push_html_button = document.getElementById("push-html-button")

	

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

		let files_to_upload = [];

		for (let i=0; i<files.length; ++i)
		{
			let file = files[i];

			// We will upload any PNG files.
			if (file.name.substr(-4) === ".png")
			{
				files_to_upload.push(file)
				continue;
			}
			// Only add a button for markdown files.
			else if (!(file.name.substr(-3) === ".md"))
			{
				continue;
			}

			// console.log(file);
			let file_button = document.createElement("button");
			file_button.textContent = file.webkitRelativePath;
			file_button.addEventListener("click", async () => 
			{
				let filetext = await file.text();
				launch_markdown_editor(filetext);
			});

			file_manager_files_buttons_div.appendChild(file_button);
		}

		if (files_to_upload.length > 0)
		{
			const formdata = new FormData();
			// Now, let's upload any and all.
			for (let i=0; i<files_to_upload.length; ++i)
			{
				formdata.append("files[]", files_to_upload[i]);
				// console.log("Attached " + files_to_upload[i].name);
			}

			fetch("/resources_dump", {
				method: "POST",
				body: formdata
			})
			.then(response => response.json())
			.then(data => {
				console.log("Uploaded the files, ", data);
			}).catch(error => {
				console.error(error);
			})
		
		}
	}


	function launch_markdown_editor(filetext)
	{
		markdown_editor_textarea.classList.remove("hidden");
		preview_html_button.classList.remove("hidden");
		markdown_editor_textarea.value = filetext;
	}

	preview_html_button.addEventListener("click", () => 
	{
		/* We will disable the button until the server replies to avoid spamming the server.*/

		preview_html_button.disabled = true;

		const formdata = new FormData();
		
		// File expects an array of stuff to put in the file.
		const markdown_file = new File([markdown_editor_textarea.value], "markdown.md", {type: "text/plain"});
		formdata.append("mdfile", markdown_file);

		fetch("/upload_markdown", 
		{
			method: "POST", 
			body: formdata
		})
		.then(response => response.text())
		.then(data => 
		{
			html_preview_iframe.classList.remove("hidden");

			const html_preview_iframe_doc = html_preview_iframe.contentDocument || html_preview_iframe.contentWindow.document;
			console.log("Console says + ", data);
			html_preview_iframe_doc.open();
			html_preview_iframe_doc.write(data);
			html_preview_iframe_doc.close();

			push_statics_button.classList.remove("hidden");
			destination_input.classList.remove("hidden");
			destination_label.classList.remove("hidden");
			push_html_button.classList.remove("hidden")
			push_html_button.disabled = true


			// Reenable
			preview_html_button.disabled=false;
		});
	});

	push_statics_button.addEventListener("click", async () => 
	{
		push_statics_button.disabled = true;

		const response = await fetch("/push-static-images");
		if (response.status == 204) 
		{
			// This is a success
			alert("Success uploading images.")	
			push_statics_button.disabled = false;
		}
		else
		{
			alert(`Error uploading images, reponse: ${response.status}`);
		}

		push_html_button.disabled = false;
	});

	push_html_button.addEventListener("click", async () => 
	{
		push_html_button.disabled=true

		response = await fetch("/push-html", 
		{
			method: "POST",
			body: destination_input.value
		});

		if(response.status == 204)
		{
			alert("Success!")
		}
		else 
		{
			alert(`Error uploading website, response: ${response.status}`)
		}
	});
});


