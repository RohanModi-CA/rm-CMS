// We'll create a function that prepares the file tree for this purpose.
import {initialize_file_tree} from "./file_tree.js"

function add_new_dir_buttons_and_submit(function_to_pass_filename_value_to)
{
	// We add the '+' buttons to add files, as well as 'D' buttons for new directories.

	const left_tree_panel_entries = document.getElementById("left-tree-panel-entries");
	const directory_classes = left_tree_panel_entries.getElementsByClassName("directory_summary_name")
	const submit_file_name_button = document.getElementById("submit-file-name-button");
	const selected_file_name_span = document.getElementById("selected-file-name-span");
	
	
	for (let i=0; i<directory_classes.length; ++i)
	{
		const new_class_button = document.createElement("button")
		new_class_button.textContent = "+"
		new_class_button.onclick = add_new_class_button_clicked
		new_class_button.classList.add("new_class_button")
		directory_classes[i].appendChild(new_class_button)
	
		const new_directory_button = document.createElement("button")
		new_directory_button.textContent = "D"
		new_directory_button.classList.add("new_directory_button")
		directory_classes[i].appendChild(new_directory_button)


	}

	submit_file_name_button.addEventListener("click", ()=>
	{
		function_to_pass_filename_value_to(selected_file_name_span.textContent)
	});
}


function add_new_class_button_clicked(event)
{
	// This event contains the button.
	const clicked_button = event.target;
	// This button's parent is the directory summary, whose parent is the dirdiv..
	const directory_div = clicked_button.parentElement.parentElement;
	
	// Now we want to get into the directory_name subclass.
	const directory_name_subclass = directory_div.querySelector(".directory_name")

	// Now in this element we want to create a new textbox element.
	const new_file_name = document.createElement("input")
	new_file_name.oninput = new_file_name_oninput;
	directory_name_subclass.appendChild(new_file_name);
	new_file_name.focus();
	enable_disable_all_add_new_dir_buttons(true)

	const selected_file_name_div = document.getElementById("selected-file-name-div");
	selected_file_name_div.classList.remove("hidden")
}


function add_new_dir_button_clicked(event)
{
	// The parent of the dir button is the summary, whose parent is the directory_details.
	// The directory_details contains the directory_name node, which is where we want to put the new file.
	const add_new_dir_button = event.target;
	const directory_details = add_new_dir_button.parentElement.parentElement;
	const directory_name = directory_details.querySelector(".directory_name")

	const new_dir_name_div = document.createElement("div")
	new_dir_name_div.classList.add("new_dir_name_div")
	
	const new_dir_name_input = document.createElement("input")
	new_dir_name_input.classList.add("new_dir_name_input")
	new_dir_name_div.appendChild(new_dir_name_input)
	const confirm_new_dir_name_button = document.createElement("button")
	confirm_new_dir_name_button.textContent = "Y"
	confirm_new_dir_name_button.classList.add("confirm_new_dir_name_button")
	new_dir_name_div.appendChild(confirm_new_dir_name_button)

	function confirm_new_dir_name_button_on_click(event) 
	{
		// We will convert the  TODO
	}
	

}




function get_filepath_from_input_box(input_box)
{
	// This function takes the input box created and recursively goes up to find the root. From
	// it then returns the full filepath array, excluding the ., and excluding rohanmodi.ca.
	
	function recursive_get_parents_array(dir_name_node, directories_array=[])
	{
		// This takes the 'directory_name' div nodes. Finds its parent directory name. 
		// If its parent is not the root directory, it calls again. Else it ends.

		// Is the parent of the directory_name node of class 'directory_details'?
		// If not, we're at the root. ROOT '.' will be included, we remove later, for simplicity. 
		
		const directory_details = dir_name_node.parentElement
		if (directory_details.classList.contains("directory_details"))
		{
			// we're not at root
			// get this directory name.
			const directory_name = directory_details.querySelector("summary .dir_summary_name_span").textContent
			directories_array.unshift(directory_name)

			// Next, we're interested in getting the directory_name node for the next.
			return recursive_get_parents_array(directory_details.parentElement, directories_array)
		}
		else
		{
			// we're at the root. In this case, we're done, just return the directories array.
			return directories_array
		}
	}

	// Now let's get the parents_array for our input box, which is itself within a dir_name_node.
	let parents_array = recursive_get_parents_array(input_box.parentElement)

	// Let's remove the '.' root, the first entry.
	parents_array.splice(0, 1)
	
	// ensure the parents_array joined starts with a / and ends with /, and only has one, if empty.
	let joined_parents_array = "/" + parents_array.join("/")
	joined_parents_array += (parents_array.length >= 1) ? "/" : ""

	// Now all we need to do is join these together and append the filename in the input box. 
	let path = "rohanmodi.ca" + joined_parents_array + input_box.value;

	return path
}


function new_file_name_oninput(event)
{
	const input_box = event.target


	const text_input = get_filepath_from_input_box(input_box)

	const selected_file_name_div = document.getElementById("selected-file-name-div");
	const selected_file_name_span = document.getElementById("selected-file-name-span");
	
	selected_file_name_span.textContent = text_input

	let is_valid_name = is_valid_file_name(text_input)

	const submit_file_name_button = document.getElementById("submit-file-name-button");
	submit_file_name_button.disabled = is_valid_name ? false : true
}



function enable_disable_all_add_new_dir_buttons(disable=true)
{
	const left_tree_panel_entries = document.getElementById("left-tree-panel-entries");
	const new_class_buttons = left_tree_panel_entries.getElementsByClassName("new_class_button")

	for(let i=0; i<new_class_buttons.length; ++i)
	{
		new_class_buttons[i].disabled = disable ? true : false
	}

}


function is_valid_file_name(filename)
{
	/* Legal filenames can contain most things. The main thing we want to
	 * check right now is to ensure that there is no other conflicting file.
	 */
	if (filename.length >0)
	{
		return true
	}
}


export function destination_selector_on_load(div_to_put_in, function_to_pass_filename_value_to)
{
	/* div to put in is the element of the div for the file-tree.
	 * function_to_pass_return_value_to is a function that takes  
	 * one argument and passes the selected filename to it. */

	// Add the stylesheet.
	const link = document.createElement("link")
	link.rel = "stylesheet"
	link.href = "destination_selector.css"
	document.head.appendChild(link)


	let destinaton_selector_div = document.createElement('div')
	div_to_put_in.appendChild(destinaton_selector_div)

	
	// So first let's fetch our own HTML file, then fetch the file_tree base.
	
	fetch('destination_selector.html')
	.then(r=>r.text())
	.then(data=>
	{
		destinaton_selector_div.innerHTML = data
	})
	.then(()=>
	{
	return fetch('file_tree.html')
	})
	.then(r=>r.text())
	.then(data=>
	{
		const file_tree_div = document.getElementById("file-tree")
		file_tree_div.innerHTML=data
	})
	.then(()=>
	{
		return initialize_file_tree()
	})
	.then(()=>
	{
		// Now we have to add the new directory buttons.
		add_new_dir_buttons_and_submit(function_to_pass_filename_value_to)
	})
}
