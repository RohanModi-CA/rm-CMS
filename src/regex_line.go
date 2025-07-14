package regex_line

import "regexp"
import "fmt"

func bold_and_italicize(line string) string {
	var processed_line string

	processed_line = line

	/* We will take care first of bold italics, which is a triple asterisk.
	 * Then we will take care of bold, a double, and italics, the single asterisk. */
	bold_italic_re := regexp.MustCompile(`\*\*\*(.*?)\*\*\*`)
	bold_re := regexp.MustCompile(`\*\*(.*?)\*\*`)
	italic_re := regexp.MustCompile(`\*(.*?)\*`)

	processed_line = bold_italic_re.ReplaceAllString(processed_line, "<b><i>$1</i></b>")
	processed_line = bold_re.ReplaceAllString(processed_line, "<b>$1</b>")
	processed_line = italic_re.ReplaceAllString(processed_line, "<i>$1</i>")



	return processed_line
}

func in_line_code(line string) string {
	var processed_line string
	
	processed_line = line
	
	in_line_code_regex := regexp.MustCompile("`(.*?)`")
	processed_line = in_line_code_regex.ReplaceAllString(processed_line, "<code>$1</code>")


	return processed_line
}

func hrule(line string) string {
	var processed_line string

	processed_line = line

	hrule_regex := regexp.MustCompile(`^[ \t]*---[ \t]*$`)
	processed_line = hrule_regex.ReplaceAllString(processed_line, "<hr>")

	return processed_line
}



func process_line(line string) string {
	var processed_line string

	processed_line = line

	// Note that by this point, we should have already handled foldable headers. 
	// If a --- exists as its own line at this point, it will be made into an hrule.
	processed_line = hrule(processed_line)
	processed_line = in_line_code(processed_line)
	processed_line = bold_and_italicize(processed_line)

	return processed_line
}


