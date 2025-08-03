package main

import "strings"
import "fmt"
import "unicode"

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
// ====== helper ======

func main() {
  
}

