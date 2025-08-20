package parsers

import "regexp"
import "unicode"
import "strings"
import "fmt"
import "log"

// ====== helper ======

func find_first_code_block(body [2]string) ([2]string, error) {
	/* 	
		Now, we're going to write a function that iterates through our text, and finds the first complete code block. It will wrap that within a `<pre><code class="language-html"></code></pre>`, which is what triggers [highlight.js](https://highlightjs.org/) highlighting. We omit the class if no language is specified.
		It will stop after the first complete code block, and return a 2-item string array:
		```go
		[`lorem ipsum <pre><code class="language-python)">dolor sit</code></pre>`, `amet, consectetur adipiscing elit ...`]
		```
		The first item is the current processed content. The second element is the rest, which is left to process.
		In fact, this 2-item array is what our function will take as input as well. For the first iteration, the first string will be empty, and the second element will be full. We will repeatedly call our function, which will act on the second element, until the second element is empty.
		It will error if it cannot close a code block. 
	*/
	var whole_rune []rune
	var in_markdown bool 
	var start_ticks int
	var end_ticks int
	var language string
	var content string
	var index int
	var char rune
	var code_class string
	var index_of_start int
	var out [2]string
	var matched bool

	start_ticks = 0
	end_ticks = 0
	index_of_start = 0
	in_markdown = false
	matched = false
	language = ""
	content = ""
	code_class = ""
	// We are going to process the second entry.
	whole_rune = []rune(body[1])
	for index=0; index<len(whole_rune); index++ {
		char = whole_rune[index]
		
		/* First, if we're not in markdown, we'll check to see if we have enough ticks to 
		 * be in markdown. */
		if (!in_markdown) {
			if(start_ticks >= 3) {
				// We have enough, do we have more?
				if (char == '`') {
					start_ticks ++  
					continue
				} else {
					// First, we'll move the start index forward by one, so that it is the first backtick. Unless its the first char.
					if(index_of_start != 0) {
						index_of_start ++
					}
					
					// Okay, we're in markdown. Let's check to see if we have a language.
					// If we don't, the next character will be a whitespace char.
					in_markdown = true
					
					for ; index<len(whole_rune); index++ {
						char = whole_rune[index]
						// unicode.IsSpace detects on '\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP).
						if (unicode.IsSpace(char)) {
							// Thing is over.
							break
						} else {
							language += string(char)
						}
					} // ending iterate over next word
				} // ending entering markdown
			} else { 
				// We don't have enough ticks. 
				if (start_ticks >= 0) {
					if (char == '`') {
						start_ticks ++
					} else {
						start_ticks = 0
						index_of_start = index
					} 
				} 
			}
		} else if (in_markdown) {
			if (end_ticks == start_ticks) {
				// Time to end the loop. 
				break
			} else if (end_ticks >= 0) {
				if (char == '`') {
					end_ticks ++
					if (end_ticks == start_ticks) {
						matched = true
					}
				} else {
					// Now, we have to give back any ticks that we skipped:
					content += strings.Repeat("`", end_ticks)

					end_ticks = 0
					content += string(char)
				}
			}
		} // ending being in markdown.
	} // end reading through.
	if (content != "") {
		
		// If it didn't match, we error.
		if(!matched) {
		    return [2]string{}, fmt.Errorf("unclosed code block starting with %d backticks", start_ticks)
		}

		// Now, we need to wrap the content properly.
		if (language != "") {
			code_class += "language-" + language
		}
		
		out[0] = body[0] + string(whole_rune[0:index_of_start]) + `<pre><code class="` + code_class + `">` + content + `</code></pre>`
		out[1] = string(whole_rune[index:])
	} else {
		out[0] = body[0] + body[1]
		out[1] = ""
	}
	return out, nil
}

func parse_table_row(line_rune []rune, columns int) []string{
	/* This takes one line, and attempts to parse it as though it had columns columns.
	   if it succeeds, it will return a list of strings that is columns long. If it fails, 
	   it will return an empty list. We assume the table is of format |col1|col2|etc|.
	*/

	var last_char rune
	var current_entry []rune
	var counted_cols int
	var char_i int
	var c rune
	var out []string

	// Initialize to -1 to treat the required leading | as a boundary, not a column separator.
	counted_cols = -1

	if (columns==0){
		return []string{}
	}

	for char_i=0; char_i<len(line_rune); char_i ++ {
		c = line_rune[char_i]

		if(c == '|' && last_char != '\\') {
			counted_cols ++
			if (counted_cols > 0) {
				out = append(out, string(current_entry))
				current_entry = []rune{}
			}
		} else if (c == '\\' && last_char != '\\') {
			last_char = c
		} else if (c == '\\' && last_char == '\\') {
			// Spaces are neutral characters that reset the escape state for the next character.
			last_char = ' ' 
			current_entry = append(current_entry, c)
		} else {
			last_char = c
			current_entry = append(current_entry, c)
		}

	}

	if (counted_cols == columns) {
		return out
	} else {
		return []string{}
	}
}


func build_html_table(table_data[][]string) string {
	/* This takes a 2d table_data array and builds a HTML table from it. 
	   table_data[0] is the first row of the table. */
	
	var out string
	var entry_start_tag string
	var entry_end_tag string

	entry_start_tag = "<th>"
	entry_end_tag = "</th>"

	out = "<table>"
	for row:=0; row<len(table_data); row++ {
		out += "<tr>"

		for col:=0; col<len(table_data[row]); col ++ {
			out += entry_start_tag + table_data[row][col] + entry_end_tag
		}

		out += "</tr>"
		// Everything after the first row should be a data row, not a header.
		entry_start_tag = "<td>"
		entry_end_tag = "</td>"
	} 

	out += "</table>"

	return out
}

func process_first_table(whole [2]string) [2]string {
	/* This takes a 2-long string array. The first entry is the part that has already been processed. 
	   The second entry has not been processed. This finds the first table in the second part, appends
	   everything before the table into the first entry, adds the HTML around the table and appends
	   that into the first entry, and leaves everything after the first table in the second entry.
	   This only accepts tables of the form |col1|col2|colETC|, with a separator |--|-:|---|
	*/

	var out [2]string
	var lines []string
	var line_rune []rune
	var index int
	var char_i int
	var c rune
	var is_separator bool
	var columns_count int
	var last_char rune
	var table_data [][]string
	var in_a_table bool
	var table_row int
	var row_data []string

	out = whole

	// First things first we will split our whole string into a list of strings
	lines = body_lines_to_slice(whole[1])
	in_a_table = false

	// The separator is on the second line of the table
	for index=1; index<len(lines); index ++ {
		line_rune = []rune(lines[index])

		if (!in_a_table) {
			is_separator = false
			last_char = '-'
			table_data = [][]string{}
			row_data = []string{}
			in_a_table = false
			table_row = 0

			// 2 | corresponds to 1 column
			columns_count = -1

			for char_i=0; char_i<len(line_rune); char_i ++ {
				c = line_rune[char_i]

				// Last character must be a closing pipe.
				if(char_i == len(line_rune) - 1) {
					if(c==' '){
						if(last_char=='|' && columns_count >= 0){
							// Valid.
							columns_count ++ 
							is_separator = true
							break
						} else {
							// Invalid
							break
						}
					}
					if(c=='|') {
						if(columns_count >= 0) {
							// Valid.
							columns_count ++
							is_separator = true
							break
						}
					}
				}

				if (c==' ') {
					continue
				} else if (c=='-' || c==':') {
					last_char = c
					continue
				} else if (c=='|'){
					if (last_char=='-' || last_char==':') {
						columns_count ++
						last_char = '|'
						continue
					} else {
						// Not a valid separator.
						break
					}
				} else {
					// This is not a separator
					break
				}
			} // within line

			// Now, if this line was a separator, let's check the previous one.	
			if (is_separator){
				line_rune = []rune(lines[index-1])
				row_data = parse_table_row(line_rune, columns_count)

				if(len(row_data) == columns_count) {
					table_data = append(table_data, row_data)
					in_a_table = true
					table_row ++

					is_separator = false
					// We will now append everything before this table to the processed entry of our output.
					out[0] = out[0] + strings.Join(lines[:index-1], "\n")

					continue
				} 
				is_separator = false
			}
			
		} else {
			// We're in a table.
			row_data = parse_table_row(line_rune, columns_count)
			if (len(row_data) == columns_count) {
				table_data = append(table_data,row_data)
				table_row ++


			} else {
				// Table is over
				in_a_table = false
				
				out[0] += build_html_table(table_data)
				// Dump the rest of the file into the nonprocessed entry of the output.
				out[1] = strings.Join(lines[index:], "\n")

				return out

			}
		}
	} // For over the lines

	// Handle tables that end on the last line of the file.
	if ((index == len(lines)) && in_a_table){
		// Table is over
		in_a_table = false
		
		out[0] += build_html_table(table_data)
		// Dump the rest of the file into the nonprocessed entry of the output.
		out[1] = ""
	}


	// If we reach here, that means we haven't had any tables. We've processed it all.
	out[0] = out[0] + out[1]
	out[1] = ""

	return out
}

// ====== helper ======





func foldable_header(whole string) string {
	var processed_whole string

	processed_whole = whole

	small_foldable_header_regex := regexp.MustCompile(`(?m)^[ \t]*([^\s#].*$)\r*\n[ \t]*---[ \t]*$`)
	processed_whole = small_foldable_header_regex.ReplaceAllString(processed_whole, "<div class='small-foldable-header hidden'>$1</div>")

	return processed_whole
}

func code_block(whole string) string {
	// We are going to repeatedly call find_first_code_block until it has processed the entire body.

	var whole_pair [2]string
	var err error

	whole_pair[0] = ""
	whole_pair[1] = whole

	// Apparently there is no "while" keyword in go
	for whole_pair[1] != "" {
		whole_pair, err = find_first_code_block(whole_pair)

		if(err != nil) {
			panic(err);
		}
	}

	return whole_pair[0]
}


func quote_blocks(whole string) string {
	var processed_whole string
	
	processed_whole = whole

	quote_block_regex := regexp.MustCompile(`(?m)(?:^&gt; ?.*\n?)+`)
	quote_cleaner_regex := regexp.MustCompile(`(?m)^&gt; ?`)

	processed_whole = quote_block_regex.ReplaceAllStringFunc(processed_whole, func (quote_block string) string {
		content := quote_cleaner_regex.ReplaceAllString(quote_block, "")
		content = strings.TrimSpace(content)
		
		return (fmt.Sprintf("<div class='quote-block'><blockquote>%s</blockquote></div>",content))
	})

	return processed_whole
}


func tables(whole string) string {
	// We are going to repeatedly call process_first_table until it has processed the entire body.

	var whole_pair [2]string

	whole_pair[0] = ""
	whole_pair[1] = whole

	for whole_pair[1] != "" {
		whole_pair = process_first_table(whole_pair)
	}

	return whole_pair[0]
}


func postprocess_section_header(whole string) string {
	/* First we need to find the div. Note that when we created the divs, we did it like this:
	 processed_whole = small_foldable_header_regex.ReplaceAllString(processed_whole, "<div class='small-foldable-header hidden'>$1</div>")
	 This means we just have to search for instances of this. We'll create a search that searches for this pattern, adds another div header before,
	 and then searches for the end of the section. If it is found, we add in the closing div, and if it is not, we put it in at the end, before the
	 end of the body tag.
	 */

	 /* We're going to store all starting indices of regexp matches in a list. */
	 
	 var matches_indices [][]int
	 var index_after_which_to_insert_closing int
	 var index_after_which_to_start_opening int
	 var current_div_str string
	 var header_text string
	 var processed_whole string
	 var past_length_whole int
	 var offset_length_whole int

	 processed_whole = whole
	 small_foldable_header_div_regex := regexp.MustCompile(`<div class='small-foldable-header hidden'>(.*?)</div>`)
	 matches_indices = small_foldable_header_div_regex.FindAllStringIndex(whole, -1)

	 // Now that we have this list, we're going to loop through, and check whether it goes to the end or not, and handle each case.
	 // We will first add the closing, since that comes after the opening, so adding the closing will not change the index of the start.
	 // We'll then calculate the amount of new characters we've added, and then offset all future indices by that amount. Since we are going
	 // in order, and there is no nesting or anything, this should suffice.

	 for i:=0; i<len(matches_indices); i++ {

		 past_length_whole = len(processed_whole)

		 // First, we'll extract the substring corresponding to this match to extract its header text.
		 current_div_str = processed_whole[matches_indices[i][0]:matches_indices[i][1]]

		 matches := small_foldable_header_div_regex.FindStringSubmatch(current_div_str)
		 if (matches == nil || len(matches) == 0) {
			 // Don't think this is possible but let's just handle this.
			 log.Fatal("Logic error in processing section header divs")
		 }

		 header_text = matches[1]

		 index_after_which_to_start_opening = matches_indices[i][1]
		 
		 // If it is the last match, that means it goes to the end.
		 if (i == len(matches_indices) - 1) {
			 index_after_which_to_insert_closing = len(processed_whole) - 1
		 } else {
			 index_after_which_to_insert_closing = matches_indices[i+1][0] - 1
		 } // else case: there is another foldable_header after this one.

		 // We set the closing </details> in the appropriate place.
		 processed_whole = processed_whole[:index_after_which_to_insert_closing + 1] + 
		 		"</details>" + processed_whole[index_after_which_to_insert_closing + 1 :]

		 // Now, we set the opening details and the summary header.
		 processed_whole = processed_whole[:index_after_which_to_start_opening + 1] + "<details open>\n\t<summary>" + 
		 					header_text + "</summary>\n" + processed_whole[index_after_which_to_start_opening + 1 :]
		 
		
		// Now, let's see how much we need to offset this by.
		offset_length_whole = len(processed_whole) - past_length_whole

		for j:=i+1; j<len(matches_indices); j++ {
			matches_indices[j][0] += offset_length_whole
			matches_indices[j][1] += offset_length_whole
		} 

	 } // looping over all matches
	 
	 return processed_whole
}


func process_whole(whole string) string {
	var processed_whole string
	

	processed_whole = whole
	processed_whole = code_block(processed_whole)
	processed_whole = foldable_header(processed_whole)
	processed_whole = quote_blocks(processed_whole)
	processed_whole = tables(processed_whole)


	// Post Processing
	processed_whole = postprocess_section_header(processed_whole)

	return processed_whole
}




