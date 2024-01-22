package structs

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description:
		Chronos is a simple, fast and lightweight documentation server written in Go.
*/

// ArticlesResponse is the response struct for the /articles endpoint.
type ArticlesResponse struct {
	Title         string    `json:"title"`
	SupportedLang []string  `json:"SupportedLang"`
	Articles      []Article `json:"articles"`
}
