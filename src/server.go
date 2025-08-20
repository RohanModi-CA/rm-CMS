package main

import "fmt"
import "io"
import "path/filepath"
import "net/http"
import "os"
import "cms/parsers"
import "cms/versioning"

// Our global conversion state. Uppercase means it is exported throughout main package.
var GlobalConversionState parsers.ConversionState


func main() {
	http.Handle("/", http.FileServer(http.Dir("../site")))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("cms-resources/images"))))

	http.HandleFunc("/upload_markdown", process_markdown_file)
	http.HandleFunc("/resources_dump", process_resource_dump)
	http.HandleFunc("/push-static-images", push_static_images)

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
	html_out := parsers.MainCall(string(filebytes), &GlobalConversionState)

	// Let's tell the client we're sending it HTML.
	w.Header().Set("content-type", "text/html")


	// And send the content back to the server, alongside a success code.
	w.WriteHeader(200);
	fmt.Fprintf(w, html_out)
}


func process_resource_dump(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Call to /resources_dump received.")

	// This time, our object will be much larger. We'll allocate 100mb.
	err := r.ParseMultipartForm(100 << 20)
	if (err != nil) {
		fmt.Printf("Could not parse multipart form: %v", err)
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	// r.MultipartForm.File is a map of key []*FileHeader
	files := r.MultipartForm.File["files[]"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded under the key 'files[]'", http.StatusBadRequest)
		return
	}

	// Now, we're going to clear the uploaded files directory to make room for the new files. 
	
	if err := os.RemoveAll("uploads/resources"); err != nil {
		fmt.Printf("Error: %v\n", err)
		panic("error")
	}

	if err := os.Mkdir("uploads/resources", 0755); err != nil {
		fmt.Printf("Error recreating directory: %v\n", err)
		panic("error")
	}




	fmt.Printf("Received %d files to upload", len(files))

	var uploaded_file_names []string

	// Now we're going to loop through this and save them. We want to preserve the filenames.
	for _, file_header := range files {
		//fmt.Printf("Processing file: %s (Size: %d bytes)", file_header.Filename, file_header.Size)

		// Open it.
		file, err := file_header.Open()
		if err != nil {
			fmt.Printf("Error opening file %s: %v", file_header.Filename, err)
			http.Error(w, "Error processing file", http.StatusInternalServerError)
			return
		}
		defer file.Close() // Important to close the file

		// Create the destination file
		destination_path := filepath.Join("uploads","resources", file_header.Filename)
		dst, err := os.Create(destination_path)
		if err != nil {
			fmt.Printf("Error creating destination file %s: %v", destination_path, err)
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Now we copy the uploaded file's content to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			fmt.Printf("Error copying file content for %s: %v", file_header.Filename, err)
			http.Error(w, "Error saving file content", http.StatusInternalServerError)
			return
		}

		// Add the filename to our list of successful uploads
		uploaded_file_names = append(uploaded_file_names, file_header.Filename)
		w.WriteHeader(200)
		fmt.Fprintf(w, "Successfully uploaded the files")
	} // end for loop
}

func push_static_images(w http.ResponseWriter, r *http.Request) {

	versioning.PublishStatics(&GlobalConversionState)
	w.WriteHeader(204)
}





