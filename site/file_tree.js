/**
 * @param {HTMLDivElement} parent_node - This is generally the "directory_name" div of where you want to place the new dir. Unless its the first directory.
 * @param {string} dir_name - The name of the directory to create. Will not be checked for validity.
 * @param {boolean} open - Whether the details of this directory should be open by default.
 * @returns {{directory_details: HTMLDetailsElement, directory_name_div: HTMLDivElement}} - The directory_details for the new div, as well as the new directory_name_div, where you would put nested files, etc.
 */
export function create_new_dir_div(parent_node, dir_name, open=false)
{
	let dir_details = document.createElement("details")
	dir_details.classList.add("directory_details")
	parent_node.appendChild(dir_details)

	dir_details.open = open;

	let dir_summary_name = document.createElement("summary")
	dir_summary_name.classList.add("directory_summary_name")
	let dir_summary_name_span = document.createElement("span")
	dir_summary_name_span.textContent = dir_name
	dir_summary_name_span.classList.add("dir_summary_name_span")
	dir_details.appendChild(dir_summary_name)
	dir_summary_name.appendChild(dir_summary_name_span)

	let dir_div = document.createElement("div")
	dir_div.classList.add("directory_name")
	dir_details.appendChild(dir_div)	

	let out = {directory_details: dir_details, directory_name_div: dir_div}
	return out
}


export async function initialize_file_tree()
{
	// Add the stylesheet
	const link = document.createElement("link")
	link.rel = "stylesheet"
	link.href = "file_tree.css"
	document.head.append(link)

	const left_tree_panel_entries = document.getElementById("left-tree-panel-entries")

	// This takes the entire site_map.
	function renderTree(site_map)
	{
		let actual_tree = site_map[0]
		
		// We need to recursively dive in.
		function recursive_create(entry, count=0, root_div=left_tree_panel_entries)
		{
			// We will do something a bit hacky and just say that if
			// count is zero, then we are on the root '.', which we want to
			// autoexpand. Sorry!

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

			/** @type {boolean} */ 
			let first_entry = (count == 0)

			let {directory_name_div} = create_new_dir_div(root_div, dir_name, first_entry)

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
					let file_name_span = document.createElement("span")
					file_name_span.textContent = file_name
					file_name_span.classList.add("file_name_span")
					file_div.classList.add("file_name")
					directory_name_div.appendChild(file_div)
					file_div.appendChild(file_name_span)
				}
				else if (dir_contents[i]["type"] == "directory")
				{
					recursive_create(dir_contents[i], ++count, directory_name_div)
				}
			}
		}

		recursive_create(actual_tree)
	}

	// Now, here we are going to start building it, but we need our JSON. Let's assume we have that in site_map.json.
	
	return fetch('./site_map.json')
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
}
