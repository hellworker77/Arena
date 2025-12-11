package app

import (
	"fmt"
	"itemsParser/internal/fetch"
	"itemsParser/internal/parse"
)

/*var categories = []string{
	"spears",
}*/

var categories = []string{
	"helms",
	"armor",
	"shields",
	"weapons",
	"gloves",
	"boots",
	"belts",
	"axes",
	"bows",
	"crossbows",
	"daggers",
	"javelins",
	"maces",
	"polearms",
	"scepters",
	"spears",
	"staves",
	"spears",
	"swords",
	"throw",
	"wands",
	"barbhelms",
	"druidpelts",
	"paladinshields",
	"shrunkenheads",
	"amazonweapons",
	"katars",
	"orbs",
	"circlets",
}

/*
	var variants = []string{
		"normal",
	}
*/
var variants = []string{
	"normal",
	"exceptional",
	"elite",
}

func LoadBaseItems() {
	for _, category := range categories {
		for _, variant := range variants {
			html, err := fetch.LoadCategory(category, variant)
			if err != nil {
				panic(err)
			}

			lines, err := parse.ExtractLines(html)
			if err != nil {
				panic(err)
			}

			objects, err := parse.ExtractObjects(lines)
			if err != nil {
				panic(err)
			}

			for _, obj := range objects {
				fmt.Println(obj.IconID)
			}
			fmt.Println("")
		}
		fmt.Println("\n--------------------------------\n")
	}
}
