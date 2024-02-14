package structs

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description:
		Chronos is a simple, fast and lightweight documentation server written in Go.
*/

import "github.com/russross/blackfriday/v2"

type Article struct {
	Title           string
	Description     string
	PublicationDate string
	Authors         []string
	Tags            []string
	Body            string
	Language        string
	Path            string
	Url             string
	Slug            string
}

// ParseBody parses the body of an article and converts it from Markdown to HTML.
func (a *Article) ParseBody() {
	parsedBody := blackfriday.Run([]byte(a.Body))
	a.Body = string(parsedBody)
}

type ArticleHeader struct {
	Title           string   `yaml:"Title"`
	Description     string   `yaml:"Description"`
	PublicationDate string   `yaml:"PublicationDate"`
	Authors         []string `yaml:"Authors"`
	Tags            []string `yaml:"Tags"`
}
