package feedparser

// Link struct for the link element
type Link struct {
	Href string `xml:"href,attr"`
}

// Video struct for each video in the feed
type Video struct {
	ID        string `xml:"id"`
	Title     string `xml:"title"`
	Published Date   `xml:"published"`
	Link      Link   `xml:"link"`
	IsShort   bool
}

// Channel struct for RSS
type Channel struct {
	ID     string
	Name   string   `xml:"title"`
	Videos []*Video `xml:"entry"`
}
