package misc

type ConversionState struct {
	CurrentlyBusy bool
	ImagesRelativePaths []string
	HtmlFileContents string
	WebsiteRelativePath string
	Level LogLevel
}

