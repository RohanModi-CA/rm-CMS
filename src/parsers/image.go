package parsers

import (
	"fmt"
	"log"
	"os"
	"github.com/davidbyttow/govips/v2/vips"
	"crypto/sha256"
	"encoding/hex"
)

func to_webp(input_filepath string) string {
	// Returns the relative path to the output file.

	var output_filepath: str

	vips.Startup(nil)
	defer vips.Shutdown()

	image, err := vips.NewImageFromFile(input_filepath)
	if (err != nil) {
		log.Fatalf("Vips failed to open image file. %s", err)
	}
	defer image.Close()

	fmt.Printf("Successfully loaded image. Dimensions: %d x %d\n", image.Width(), image.Height())

	params := vips.NewWebpExportParams()
	params.Quality = 75
	params.StripMetaData = true
	params.ReductionEffort = 6 
	fmt.Printf("Compressing to WebP with Quality=%d...\n", params.Quality)

	imageBytes, _, err := image.Export(params)
	if err != nil {
		log.Fatalf("Failed to export image to webp: %s", err)
	}

	// Now, let's hash it to get the filename.
	hash := sha256.Sum256(imageBytes)
	hashStr := hex.EncodeToString(hash[:])
	fmt.Printf("Hash: %s\n", hashStr)

	output_filepath = "../cms-resources/images/" + hashStr + ".webp" 

	// nobody needs to be able to execute a .webP file.
	err = os.WriteFile(output_filepath, imageBytes, 0644)
	if err != nil {
		log.Fatalf("Failed to save output file: %s", err)
	}


}
