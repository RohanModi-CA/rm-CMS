package parsers

import "log"
import "regexp"
import "net/url"
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

func images(line string, cs *ConversionState) string {
	/* Identifies images embedded within the input string, converts them to webP with a hashed filename. Adds them  
	returns the processed line. Appends to cs.images_relative_paths the paths of the image it finds.
	// TODO: FIX MULTIPLE IMAGES ON THE SAME LINE.
	*/


	var processed_line string
	var image_path string
	var new_webp_paths [2]string

	processed_line = line

	image_regex := regexp.MustCompile(`!\[([^\]]*)\]\(([^\)]+)\)`)

	matches := image_regex.FindStringSubmatch(processed_line)
	if (matches != nil) {
		image_path = matches[2]

		decoded_path, err := url.PathUnescape(image_path)
		if (err != nil) {
			log.Fatalf("Unable to decode path %v", image_path)
		}


		
		// Now, we're going to attempt to run the image conversion.
		new_webp_paths = to_webp("uploads/"+decoded_path)


		replacement := fmt.Sprintf(`<div class="picture"><img alt="$1" loading="lazy" decoding="async" src="https://static.rohanmodi.ca/images/%s" onerror="%s=null; if(%s){this.src='images/%s';}" ></div>`, new_webp_paths[1], "this.onerror", "window.IN_DEVELOPMENT", new_webp_paths[1])
		processed_line = image_regex.ReplaceAllString(processed_line, replacement)

		cs.ImagesRelativePaths = append(cs.ImagesRelativePaths, new_webp_paths[0])
				
	}

	return processed_line
}

func hyperlinks(line string) string {
	var processed_line string

	processed_line = line

	hyperlink_regex := regexp.MustCompile(`(^|[^!])\[(\S.*)\]\((\S.*)\)`)
	processed_line = hyperlink_regex.ReplaceAllString(processed_line, `$1<a href="$3">$2</a>`)

	return processed_line
}


func process_line(line string, cs *ConversionState) string {

	var processed_line string

	processed_line = line

	// Note that by this point, we should have already handled foldable headers. 
	// If a --- exists as its own line at this point, it will be made into an hrule.
	processed_line = hrule(processed_line)
	processed_line = in_line_code(processed_line)
	processed_line = bold_and_italicize(processed_line)
	processed_line = images(processed_line, cs)
	processed_line = hyperlinks(processed_line)


	return processed_line
}


