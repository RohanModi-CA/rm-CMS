package main

import "fmt"
import "io"
import "net/http"
import "cms/parsers"

func main() {
	
	http.Handle("/", http.FileServer(http.Dir("../site")))

	http.HandleFunc("/upload_markdown", process_markdown_file)

	fmt.Println("Started listening on the Port 8080.")
	http.ListenAndServe(":8080", nil)


}


// This takes the HTTP request that sends the markdown file for processing.
func process_markdown_file(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Call to /upload_markdown received.")

	// The client JavaScript has sent us a "multipart/form-data" object, which the http
	// library can parse for us. We'll "allocate" 32 megabytes for it.
	r.ParseMultipartForm(30 << 20)

	// In Javascript, we labeled the file as "mdfile", so let's grab that now.
	file, handler, err := r.FormFile("mdfile")

	if err != nil {
		fmt.Println("Error retrieving the md file. ")
		fmt.Println(err)
		return
	}

	defer file.Close()

	// The +v format specifier prints values as well as struct field names.
	fmt.Printf("Uploaded file: %+v\n", handler.Filename)
	fmt.Printf("File size: %+v\n", handler.Size)

	filebytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading the markdown file")
		fmt.Println(err)
		return
	}

	// Now, we'll process the file text.
	html_out := parsers.MainCall(string(filebytes))

	// Let's tell the client we're sending it HTML.
	w.Header().Set("content-type", "text/html")


	// And send the content back to the server, alongside a success code.
	w.WriteHeader(200);
	fmt.Fprintf(w, html_out)
}


