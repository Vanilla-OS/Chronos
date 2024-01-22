package core

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description:
		Chronos is a simple, fast and lightweight documentation server written in Go.
*/

import (
	"strings"
)

// parseArticleHeader parses the header of an article and extracts the details
// like title, description, publication date and authors.
func parseArticleHeader(header string) (string, string, string, []string) {
	title := ""
	description := ""
	publicationDate := ""
	authors := make([]string, 0)

	lines := strings.Split(header, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "Title":
			title = value
		case "Description":
			description = value
		case "PublicationDate":
			publicationDate = value
		case "Authors":
			authors = strings.Split(value, ",")
		}
	}

	return title, description, publicationDate, authors
}
