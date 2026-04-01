package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVideoInput(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantID       string
		wantProgress int
		wantErr      bool
	}{
		{
			name:   "plain video ID",
			input:  "dQw4w9WgXcQ",
			wantID: "dQw4w9WgXcQ",
		},
		{
			name:   "www.youtube.com watch URL",
			input:  "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			wantID: "dQw4w9WgXcQ",
		},
		{
			name:         "youtube.com watch URL with t param",
			input:        "https://youtube.com/watch?v=dQw4w9WgXcQ&t=9780",
			wantID:       "dQw4w9WgXcQ",
			wantProgress: 9780,
		},
		{
			name:   "m.youtube.com with extra params",
			input:  "https://m.youtube.com/watch?v=dQw4w9WgXcQ&pp=abc",
			wantID: "dQw4w9WgXcQ",
		},
		{
			name:   "youtube.com live URL",
			input:  "https://www.youtube.com/live/dQw4w9WgXcQ?si=abc",
			wantID: "dQw4w9WgXcQ",
		},
		{
			name:   "youtu.be short URL",
			input:  "https://youtu.be/dQw4w9WgXcQ",
			wantID: "dQw4w9WgXcQ",
		},
		{
			name:         "www.youtube.com watch URL with t param",
			input:        "https://www.youtube.com/watch?v=dQw4w9WgXcQ&t=120",
			wantID:       "dQw4w9WgXcQ",
			wantProgress: 120,
		},
		{
			name:    "URL missing video ID",
			input:   "https://www.youtube.com/watch",
			wantErr: true,
		},
		{
			name:    "invalid t param",
			input:   "https://www.youtube.com/watch?v=dQw4w9WgXcQ&t=abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseVideoInput(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantID, got.id)
			assert.Equal(t, tt.wantProgress, got.progressSeconds)
		})
	}
}
