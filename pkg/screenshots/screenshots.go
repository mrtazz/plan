package screenshots

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Importer struct {
	source      string
	destination string
	format      string
	noDryRun    bool
}

func NewImporter(source, destination, format string, noDryRun bool) *Importer {

	if strings.HasPrefix(source, "~/") {
		dirname, _ := os.UserHomeDir()
		source = filepath.Join(dirname, source[2:])
	}

	return &Importer{
		source:      source,
		destination: destination,
		format:      format,
		noDryRun:    noDryRun,
	}
}

func (i *Importer) ImportToPlanFolder() error {
	err := filepath.Walk(i.source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("error accessing screenshot source directory '%s': %s\n", i.source, err)
			return err
		}

		// we have found a screenshot matching the pattern
		if !info.IsDir() {
			basePath := filepath.Base(path)
			screenshotDate, err := time.Parse(i.format, basePath)
			if err != nil {
				// we don't return the error here since this probably just means the
				// file isn't a screenshot we are looking for
				// TODO: better surface actual errors of malformed screenshot
				// formatting
				return nil
			}

			// run destination through time formatting to make sure any desired date
			// formats are included
			formattedDestination := screenshotDate.Format(i.destination)
			// make sure the destination exists
			exists, err := dirExists(formattedDestination)
			if err != nil {
				return err
			}
			if !exists {
				if i.noDryRun {
					if err := os.MkdirAll(formattedDestination, os.ModePerm); err != nil {
						return err
					}
				} else {
					fmt.Printf("Dry run: Not creating destination directory '%s'...\n", formattedDestination)
				}
			}
			if i.noDryRun {
				if err := os.Rename(path, filepath.Join(formattedDestination, basePath)); err != nil {
					return err
				}
			} else {
				fmt.Printf("Dry run: Not moving file '%s' to '%s'...\n", path, filepath.Join(formattedDestination, basePath))
			}
		}
		return nil
	})

	return err
}

// exists returns whether the given file or directory exists
func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
