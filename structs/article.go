package structs

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2023
	Description:
		Chronos is a simple, fast and lightweight documentation server written in Go.
*/

import "github.com/russross/blackfriday/v2"

type Article struct {
	Title           string
	Description     string
	PublicationDate string
	Authors         []string
	Body            string
	Language        string
	Path            string
	Url             string
}

// ParseBody parses the body of an article and converts it from Markdown to HTML.
func (a *Article) ParseBody() {
	parsedBody := blackfriday.Run([]byte(a.Body))
	a.Body = string(parsedBody)
}
