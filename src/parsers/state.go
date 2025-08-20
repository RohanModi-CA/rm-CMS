package parsers

type ConversionState struct {
	CurrentlyBusy bool
	ImagesRelativePaths []string
	html_file_contents string
}

