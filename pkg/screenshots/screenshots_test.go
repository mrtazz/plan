package screenshots

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScreenshotImporter(t *testing.T) {
	assert := assert.New(t)
	i := NewImporter("testdata/input", "testdata/output/20060102",
		"Screenshot 2006-01-02 at 15.04.05.png", true)

	// write a dummy file
	err := os.WriteFile("testdata/input/Screenshot 2024-02-01 at 14.17.18.png", []byte("I'm not actually a screenshot!"), 0666)
	assert.Equal(nil, err)

	err = i.ImportToPlanFolder()
	assert.Equal(nil, err)

	// make sure the file is in the new location
	_, err = os.Stat("testdata/output/20240201/Screenshot 2024-02-01 at 14.17.18.png")
	assert.Equal(nil, err)

	// make sure the file is not in the old location anymore (i.e. was moved and
	// not copied)
	_, err = os.Stat("testdata/input/Screenshot 2024-02-01 at 14.17.18.png")
	assert.True(errors.Is(err, os.ErrNotExist))

}
