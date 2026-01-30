package parsers

import (
	"bufio"
	"cms/config"
	"cms/misc"
	"fmt"
	"os"
	"strings"
)

// ============ Helper functions ============
func readfile_lines_to_slice(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file, exiting.")
		os.Exit(1)
	}
	defer f.Close()

	var lines []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	if s.Err() != nil {
		fmt.Println("Error with the readfile_lines_to_slice, exiting.")
		os.Exit(1)
	}

	return lines
}

func body_lines_to_slice(body string) []string {
	return strings.Split(body, "\n")
}

// ============ Helper functions ============

func create_header(CS *misc.ConversionState, Conf *config.Config) string {
	var header_prefix string
	var header_additions string
	var header_postfix string
	var header string
	var rmCMSVersionMeta string

	rmCMSVersionMeta = fmt.Sprintf(`<meta name="RM_CMS_VERSION" content="%s">`, misc.GetRMCMSVersion(CS))

	header_prefix = ` <!DOCTYPE HTML>
	<head>`

	header_additions = `<script>window.IN_DEVELOPMENT=true;</script>`
	header_additions += rmCMSVersionMeta

	header_postfix = `</head>`

	header = header_prefix + Conf.HeadContents + header_additions + header_postfix

	return header
}

type TitleInfo struct {
	Title  string
	Author string
	Date   string
}

func create_title_html(title_info TitleInfo) string {
	var title_html string

	title_html = fmt.Sprintf("<span class='post-title'>%s</span>\n<hr>\n<div class='author-container'><span class='author'>%s</span>\n<span class='date'>%s</span></div><hr>", title_info.Title, title_info.Author, title_info.Date)

	return title_html
}

func replace_all_html_special_chars(body string) string {
	// We replace the ampersand first so we don't break the rest.
	body = strings.ReplaceAll(body, "&", "&amp;")
	body = strings.ReplaceAll(body, "\"", "&quot;")
	body = strings.ReplaceAll(body, "<", "&lt;")
	body = strings.ReplaceAll(body, ">", "&gt;")
	body = strings.ReplaceAll(body, "'", "&apos;")
	// body = strings.ReplaceAll(body, " ", "&nbsp;")

	return body
}

func process(body string, CS *misc.ConversionState, Conf *config.Config) string {
	var header string
	var full_html string
	var title_info TitleInfo
	var title_html string

	header = create_header(CS, Conf)

	title_info = TitleInfo{
		Title:  "Test",
		Author: "Rohan Modi",
		Date:   "July 13, 2025",
	}
	title_html = create_title_html(title_info)

	// Add the body tags around the body
	body = "<body> <div class='post'>\n" + title_html + body + "\n</div></body>\n"

	full_html = header + body

	return full_html
}

func prepreprocess_md_file(body string, CS *misc.ConversionState) string {
	var lines_slice []string
	var filecontent_str string

	// The very first thing we need to do is replace any characters that conflict with html.
	body = replace_all_html_special_chars(body)

	// First we process the whole thing.
	body = process_whole(body)

	// Then we split into lines to deal with line parsing.
	lines_slice = body_lines_to_slice(body)

	// Start iterating over the lines.
	for _, line := range lines_slice {
		//var div_classes []string
		//var div_classes_str string
		var trimmed_line string

		trimmed_line = strings.TrimSpace(line)

		if trimmed_line == "" {
			line = "<br>"
		} else {
			line = process_line(line, CS)
		}

		/*

			div_classes = append(div_classes, "post-text")


			// Wrap everything in a div
			div_classes_str = strings.Join(div_classes, " ")
			filecontent_str += fmt.Sprintf("<div class='%s'>", div_classes_str) + "\n\t" + line + "\n</div>\n"
		*/
		filecontent_str += line + "\n"
	}

	return filecontent_str
}

func MainCall(body string, CS *misc.ConversionState, Conf *config.Config) string {
	var full_html string

	body = prepreprocess_md_file(body, CS)

	full_html = process(body, CS, Conf)
	return full_html
}
