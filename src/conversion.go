package main

import "fmt"
import "os"
import "bufio"
import "strings"
import "regexp"



// ============ Helper functions ============
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
// ============ Helper functions ============




func create_header() string {
	var header string 
	
	header = `<!DOCTYPE HTML>
<head>
	<meta charset="UTF-8">
	<link rel="stylesheet" href="fonts.css">
	<link rel="stylesheet" href="style.css">
	<script src="script.js"></script>
</head>
`
	return header
}





type TitleInfo struct {
	Title string
	Author string
	Date string
}

func create_title_html(title_info TitleInfo) string{
	var title_html string

	title_html = fmt.Sprintf("<span class='post-title'>%s</span>\n<hr>\n<div class='author-container'><span class='author'>%s</span>\n<span class='date'>%s</span></div><hr>", title_info.Title, title_info.Author, title_info.Date)

	return title_html
}


func process(body string) string {
	var header string
	var full_html string
	var title_info TitleInfo
	var title_html string

	header = create_header()

	title_info = TitleInfo{
		Title:"Test",
		Author:"Rohan Modi",
		Date:"July 13, 2025",
	}
	title_html = create_title_html(title_info)



	

	// Add the body tags around the body
	body = "<body> <div class='post'>\n" +  title_html + body + "\n</div></body>\n"

	

	full_html = header + body

	return full_html
}


func prepreprocess_md_file(md_filepath string) string {
	var lines_slice []string
	var filecontent_str string

	// Read the file
	lines_slice = readfile_lines_to_slice(md_filepath)

	// Start iterating over the lines. 
	for _, line := range lines_slice {

		var div_classes []string
		var div_classes_str string
		var trimmed_line string

		trimmed_line = strings.TrimSpace(line)

		if (trimmed_line == ""){
			line = "<br>"
		} else if (trimmed_line == "---"){
			line = "<hr>"
		}



		div_classes = append(div_classes, "post-text")

		
		// Wrap everything in a div
		div_classes_str = strings.Join(div_classes, " ")
		filecontent_str += fmt.Sprintf("<div class='%s'>", div_classes_str) + "\n\t" + line + "\n</div>\n"
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



