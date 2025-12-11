package parse

import "unicode"

func isInIgnoreList(line string) bool {
	ignored := []string{
		"",
		"1 Handed Axe",
		"2 Handed Axe",
		"Daggers",
		"Clubs",
		"Hammers",
		"Scepters",
		"One-Handed Swords",
		"Two Handed Swords",
		"Throwing Weapons",
		"Wands",
	}

	for _, ign := range ignored {
		if line == ign {
			return true
		}
	}
	return false
}

func isNameLine(line string) bool {
	if line == "" {
		return false
	}

	r := rune(line[0])
	return unicode.IsLetter(r)
}

func isValueLine(line string) bool {
	if line == "" {
		return false
	}

	list := []string{
		"Yes",
		"No",
		"Bow",
		"Maces",
		"Spear",
		"Javelin",
		"Dagger",
		"Axe",
		"Light",
		"Medium",
		"Heavy",
	}

	for _, val := range list {
		if line == val {
			return true
		}
	}
	return false
}

func isModelLine(line string) bool {
	return propsHandlers[line] != nil
}
