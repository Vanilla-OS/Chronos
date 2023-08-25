package structs

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2023
	Description:
		Chronos is a simple, fast and lightweight documentation server written in Go.
*/

type Repo struct {
	Id              string
	Path            string
	Articles        map[string]Article
	ArticlesGrouped map[string][]Article
	Languages       []string
}
