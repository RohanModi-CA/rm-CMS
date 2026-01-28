package misc

type ConversionState struct {
	CurrentlyBusy bool
	ImagesRelativePaths []string
	HtmlFileContents string
	MdFileContents string
	WebsiteRelativePath string
	Level LogLevel
}



