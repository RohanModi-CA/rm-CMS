package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestFind_first_code_block(t *testing.T) {
	// We define a struct to hold our test cases
	testCases := []struct {
		name        string     // A descriptive name for the test case
		input       [2]string  // The input to the function
		want        [2]string  // The expected output array on success
		wantErr     bool       // Whether we expect an error
		errContains string     // A substring the error message should contain
	}{
		{
			name:  "Simple case with language",
			input: [2]string{"", "Some text before ```go\nfmt.Println(\"Hello\")\n``` and some text after."},
			want:  [2]string{`Some text before <pre><code class="language-go">fmt.Println("Hello")
</code></pre>`, ` and some text after.`},
			wantErr: false,
		},
		{
			name:  "No language specifier",
			input: [2]string{"", "```\nJust some code\n```"},
			want:  [2]string{`<pre><code class="">Just some code
</code></pre>`, ""},
			wantErr: false,
		},
		{
			name:  "No code block present",
			input: [2]string{"", "This is just some plain text with no code blocks."},
			want:  [2]string{"This is just some plain text with no code blocks.", ""},
			wantErr: false,
		},
		{
			name:    "Unclosed code block",
			input:   [2]string{"", "Here is some code ```python\nprint('oh no')"},
			wantErr: true,
			errContains: "unclosed code block",
		},
		{
			name:  "Code block at the very beginning",
			input: [2]string{"", "```js\nconsole.log('start');\n```The rest of the text."},
			want:  [2]string{`<pre><code class="language-js">console.log('start');
</code></pre>`, `The rest of the text.`},
			wantErr: false,
		},
		{
			name:  "Code block at the very end",
			input: [2]string{"", "The text ends with code.```sh\necho 'done'\n```"},
			want:  [2]string{`The text ends with code.<pre><code class="language-sh">echo 'done'
</code></pre>`, ""},
			wantErr: false,
		},
		{
			name:  "More than three backticks",
			input: [2]string{"", "````rust\nfn main() {}\n````"},
			want:  [2]string{`<pre><code class="language-rust">fn main() {}
</code></pre>`, ""},
			wantErr: false,
		},
		{
			name:    "Mismatched number of backticks (unclosed)",
			input:   [2]string{"", "````ruby\nputs 'hi'\n```"},
			wantErr: true,
			errContains: "unclosed code block starting with 4 backticks",
		},
		{
			name:  "Code block contains fewer backticks",
			input: [2]string{"", "````go\nvar s = `a raw string`\n````"},
			want:  [2]string{`<pre><code class="language-go">var s = ` + "`a raw string`" + `
</code></pre>`, ""},
			wantErr: false,
		},
		{
			name:  "Only one code block processed at a time",
			input: [2]string{"", "First block: ```c\nint i=0;\n```\nSecond block: ```c\nint j=0;\n```"},
			want:  [2]string{"First block: <pre><code class=\"language-c\">int i=0;\n</code></pre>", "\nSecond block: ```c\nint j=0;\n```"},
			wantErr: false,
		},
		{
			name:  "Empty input string",
			input: [2]string{"", ""},
			want:  [2]string{"", ""},
			wantErr: false,
		},
		{
			name:  "Input with only non-code-block backticks",
			input: [2]string{"", "This is `inline` code, not a ``block``."},
			want:  [2]string{"This is `inline` code, not a ``block``.", ""},
			wantErr: false,
		},
		{
			name:  "Iterative call with pre-existing content",
			input: [2]string{"<p>Previous content.</p>", "Now find this: ```\ncode\n```"},
			want:  [2]string{"<p>Previous content.</p>Now find this: <pre><code class=\"\">code\n</code></pre>", ""},
			wantErr: false,
		},
	}

	// Iterate over the test cases
	for _, tc := range testCases {
		// t.Run creates a sub-test, which gives clearer output on failure
		t.Run(tc.name, func(t *testing.T) {
			// Run the function we're testing
			got, err := find_first_code_block(tc.input)

			// Check for an expected error
			if tc.wantErr {
				if err == nil {
					t.Errorf("find_first_code_block() was expected to return an error, but did not")
					return // Stop further checks if we expected an error and got none
				}
				if !strings.Contains(err.Error(), tc.errContains) {
					t.Errorf("find_first_code_block() returned error '%v', but expected it to contain '%v'", err, tc.errContains)
				}
				return // Test is successful if the correct error was returned
			}

			// Check for an unexpected error
			if !tc.wantErr && err != nil {
				t.Errorf("find_first_code_block() returned an unexpected error: %v", err)
				return
			}

			// Check if the output matches the expected value
			if !reflect.DeepEqual(got, tc.want) {
				// Use %#v for detailed output that shows the type and quotes strings
				t.Errorf("find_first_code_block() = %#v, want %#v", got, tc.want)
			}
		})
	}
}