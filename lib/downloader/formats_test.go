package downloader

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadFormats(t *testing.T) []byte {
	file, err := os.Open("formats.json.gz")
	require.NoError(t, err, "should open formats.json.gz")
	defer file.Close()

	reader, err := gzip.NewReader(file)
	require.NoError(t, err, "should create gzip reader")
	defer reader.Close()

	output, err := io.ReadAll(reader)
	require.NoError(t, err, "should decompress gzip")

	return output
}

func TestGetVideoFormats(t *testing.T) {
	output := loadFormats(t)

	var formatInfo map[string]interface{}
	err := json.Unmarshal(output, &formatInfo)
	require.NoError(t, err, "should parse JSON output")

	formats, ok := formatInfo["formats"].([]interface{})
	require.True(t, ok, "should have formats array")
	require.NotEmpty(t, formats, "should have at least one format")

	// Check that we got valid format data
	var hasValidFormat bool
	for _, f := range formats {
		if fMap, ok := f.(map[string]interface{}); ok {
			if height, ok := fMap["height"].(float64); ok && height > 0 {
				hasValidFormat = true
				break
			}
		}
	}
	assert.True(t, hasValidFormat, "should have at least one format with valid height")
}

func TestParseVideoFormats(t *testing.T) {
	output := loadFormats(t)

	formats := ParseFormats(output)

	assert.NotEmpty(t, formats, "should have formats")

	// Check that all formats have valid height
	for _, f := range formats {
		assert.Greater(t, f.Height, 0, "format height should be positive")
	}

	// Check that formats are sorted by height (descending)
	for i := 1; i < len(formats); i++ {
		assert.LessOrEqual(t, formats[i].Height, formats[i-1].Height,
			"formats should be sorted by height descending")
	}
}
