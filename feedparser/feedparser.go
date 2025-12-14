package feedparser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/paulrosania/go-charset/charset"
)

var (
	ErrInvalidChannelID = errors.New("invalid channel ID")
	ErrParseFailed      = errors.New("failed to parse feed")
)

var urlFormat = "https://www.youtube.com/feeds/videos.xml?channel_id=%s"

type Parser interface {
	Parse(channelID string) (*Channel, error)
}

type parser struct {
	log *slog.Logger
}

func NewParser(l *slog.Logger) *parser {
	return &parser{
		log: l,
	}
}

func (p *parser) fetch(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		p.log.Error("Failed to create request", "call", "http.NewRequest", "error", err)
		return nil, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		p.log.Error("Failed to fetch feed", "call", "http.Do", "error", err)
		return nil, err
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrInvalidChannelID
	} else if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get feed with status %d", response.StatusCode)
	}

	return response.Body, nil
}

// Parse parses a YouTube channel XML feed from a channel ID
func (p *parser) Parse(channelID string) (*Channel, error) {
	url := fmt.Sprintf(urlFormat, channelID)
	reader, err := p.fetch(url)
	if err != nil {
		return nil, err
	}

	defer reader.Close()
	xmlDecoder := xml.NewDecoder(reader)
	xmlDecoder.CharsetReader = charset.NewReader

	var channel Channel
	if err := xmlDecoder.Decode(&channel); err != nil {
		p.log.Error("Failed to decode XML for RSS feed", "call", "xml.Decode", "error", err)
		return nil, fmt.Errorf("%w: %s", ErrParseFailed, err.Error())
	}
	channel.ID = channelID
	for _, video := range channel.Videos {
		video.IsShort = strings.Contains(video.Link.Href, "/shorts/")
	}

	return &channel, nil
}
