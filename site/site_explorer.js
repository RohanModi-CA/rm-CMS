document.addEventListener("DOMContentLoaded", ()=>
{
	const left_tree_panel_entries = document.getElementById("left-tree-panel-entries")


	// This takes the entire site_map.
	function renderTree(site_map)
	{
		let actual_tree = site_map[0]
		
		// We need to recursively dive in.
		function recursive_create(entry, count=0, root_div=left_tree_panel_entries)
		{
			if (count > 100)
			{
				throw Error("Too much recursion depth.")
			}

			// So if we're the deepest, we'll just create a div for each dir,
			// and a div for each file within the dir.
			if (entry["type"] != "directory")
			{
				throw Error("Can't handle this dir structure")
			}

			let dir_contents = entry["contents"]
			let dir_name = entry["name"]

			let dir_details = document.createElement("details")
			root_div.appendChild(dir_details)

			let dir_summary_name = document.createElement("summary")
			dir_summary_name.textContent = dir_name
			dir_details.appendChild(dir_summary_name)

			let dir_div = document.createElement("div")
			dir_div.classList.add("directory_name")
			dir_details.appendChild(dir_div)



			if (!dir_contents || !dir_name)
			{
				// It should be [{}] if empty which is truthy.
				throw Error("Bad Directory?")
			}


			for(let i=0; i<dir_contents.length; i++)
			{
				if (dir_contents[i]["type"] == "file")
				{
					let file_name = dir_contents[i]["name"]
					let file_div = document.createElement("div")
					file_div.textContent = file_name
					file_div.classList.add("file_name")
					dir_div.appendChild(file_div)
				}
				else if (dir_contents[i]["type"] == "directory")
				{
					recursive_create(dir_contents[i], ++count, dir_div)
				}
			}
		}

		recursive_create(actual_tree)
	}


	// Now, here we are going to start building it, but we need our JSON. Let's assume we have that in site_map.json.
	
	fetch('./site_map.json')
	.then(r => r.json())
	  .then(site_map => 
		{
			renderTree(site_map)
			//console.log(site_map)
	  	})
	  .catch(err => 
		{
			console.error(err)
	  	});
});
