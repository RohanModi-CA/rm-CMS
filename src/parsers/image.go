package parsers

import (
	"fmt"
	"log"
	"os"
	"github.com/davidbyttow/govips/v2/vips"
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

// new variable to ensure vips is only started once per process
var init_once sync.Once

// new function to perform the one-time startup
func init_vips() {
	vips.Startup(nil)
}

func to_webp(input_filepath string) [2]string {
	// Returns an array of strings, the first is the relative path to the
	// created webP, the second is just the filename + extension.

	// ensure vips is started, but only on the first call to
	// avoid that "govips cannot be stopped and restarted"
	init_once.Do(init_vips)

	var out [2]string
	var output_filepath string


	image, err := vips.NewImageFromFile(input_filepath)
	if (err != nil) {
		log.Fatalf("Vips failed to open image file. %s", err)
	}
	defer image.Close()

	fmt.Printf("Successfully loaded image. Dimensions: %d x %d\n", image.Width(), image.Height())

	params := vips.NewDefaultWEBPExportParams()

	fmt.Printf("Compressing to WebP with Quality=%d...\n", params.Quality)

	imageBytes, _, err := image.Export(params)
	if err != nil {
		log.Fatalf("Failed to export image to webp: %s", err)
	}

	// Now, let's hash it to get the filename.
	hash := sha256.Sum256(imageBytes)
	hashStr := hex.EncodeToString(hash[:])
	fmt.Printf("Hash: %s\n", hashStr)

	output_filepath = "cms-resources/images/" + hashStr + ".webp" 
	out[0] = output_filepath
	out[1] = hashStr + ".webp"

	// nobody needs to be able to execute a .webP file.
	err = os.WriteFile(output_filepath, imageBytes, 0644)
	if err != nil {
		log.Fatalf("Failed to save output file: %s", err)
	}

	return out
}
