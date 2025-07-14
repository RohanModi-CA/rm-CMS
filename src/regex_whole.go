package regex_whole

import "fmt"
import "regexp"


func foldable_header(whole string) string {
	var processed_whole string

	processed_whole = whole

	small_foldable_header_regex := regexp.MustCompile(`(?m)^[ \t]*([^\s#].*$)\r*\n[ \t]*---[ \t]*$`)
	processed_whole = small_foldable_header_regex.ReplaceAllString(processed_whole, "<span class='small-foldable-header'>$1</span>")

	return processed_whole
	

}



func process_whole(whole string) string {
	var processed_whole string
	
	processed_whole = whole
	processed_whole = foldable_header(processed_whole)

	return processed_whole
}



