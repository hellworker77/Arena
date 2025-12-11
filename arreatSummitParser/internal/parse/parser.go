package parse

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func normalizeKey(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "  ", " ")
	return s
}

func newDocumentFromString(s string) (*goquery.Document, error) {
	reader := strings.NewReader(s)
	return goquery.NewDocumentFromReader(reader)
}

type LineEntry struct {
	Icon string
	Line string
}

func ExtractLines(html string) ([]string, error) {
	var cleanLines []string
	var icons []string

	doc, err := newDocumentFromString(html)
	if err != nil {
		return nil, err
	}

	doc.Find("center table[cellpadding='5'], center table[cellpadding='5']").Each(func(i int, table *goquery.Selection) {
		table.Find("img").Each(func(j int, img *goquery.Selection) {
			if src, ok := img.Attr("src"); ok {
				icons = append(icons, src)
			}
		})
	})
	iconIndex := 0
	doc.Find("center table[cellpadding='5'], center table[cellpadding='5']").Each(func(i int, table *goquery.Selection) {
		formated := strings.TrimSpace(table.Text())
		separated := strings.Split(formated, "\n")

		for _, line := range separated {
			line = strings.TrimSpace(line)
			if line != "" {
				if isInIgnoreList(line) {
					cleanLines = append(cleanLines, strings.ToLower(line))
				} else {
					cleanLines = append(cleanLines, line)
				}

				if isNameLine(line) && !isInIgnoreList(line) && propsHandlers[line] == nil && (!isValueLine(line) || strings.Contains(icons[iconIndex], strings.ToLower(line))) && iconIndex < len(icons) {

					cleanLines = append(cleanLines, icons[iconIndex])
					iconIndex++
				}

			}
		}

	})
	iconIndex = 0
	icons = nil
	/*for _, line := range cleanLines {
		fmt.Println(line)
	}*/
	return cleanLines, nil
}

func ExtractObjects(lines []string) ([]Base, error) {
	var list []Base
	var model []string
	var current *Base
	propIndex := 0

	for _, raw := range lines {
		line := normalizeKey(raw)

		if isInIgnoreList(line) {
			continue
		}

		if strings.Contains(line, ".gif") && current != nil {

			arr := strings.Split(line, "/")
			line = arr[len(arr)-1]

			current.IconID = strings.Split(line, ".")[0]
			continue
		}

		if isNameLine(line) && propsHandlers[line] == nil {
			if current != nil {
				list = append(list, *current)
			}

			current = &Base{ID: line}
			propIndex = 0
			continue
		}

		if isModelLine(line) {
			model = append(model, line)
			continue
		}

		if current != nil && propIndex < len(model) {
			applyStat(current, line, model[propIndex])
			propIndex++
		}
	}

	if current != nil {
		list = append(list, *current)
	}

	return list, nil
}
