package main

import "fmt"
import "os"
import "bufio"


func create_header() string {
	var header string 
	
	header = `<!DOCTYPE HTML>
<head>
	<meta charset="UTF-8">
	<link rel="stylesheet" href="style.css">
	<script src="script.js"></script>
</head>`
	return header
}

func process(body string) string {
	var header string
	var full_html string
	header = create_header()

	// Add the body tags around the body
	body = "<body>\n" + body + "\n</body>\n"
	

	full_html = header + body

	return full_html
}

func readfile_lines_to_slice(path string) []string {
	f, err := os.Open(path)
	if (err != nil) {
		fmt.Println("Error opening file, exiting.")
		os.Exit(1)
	}
	defer f.Close()

	var lines []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	
	if (s.Err() != nil) {
		fmt.Println("Error with the readfile_lines_to_slice, exiting.")
		os.Exit(1)
	}

	return lines
}


func prepreprocess_md_file(md_filepath string) string {
	var lines_slice []string
	var filecontent_str string

	// Read the file
	lines_slice = readfile_lines_to_slice(md_filepath)

	// Start iterating over the lines. 
	for _, line := range lines_slice {

		// Wrap everything in a div
		filecontent_str += "<div>\n\t" + line + "</div>"
	}


	return filecontent_str
}



func main() {
	var body string
	var full_html string
	var md_filepath string

	md_filepath = "/home/gram/Documents/FileFolder/Obsidian/CMS/1 - Intro.md"
	body = prepreprocess_md_file(md_filepath)


	full_html = process(body)
	fmt.Println(full_html)

	
}



