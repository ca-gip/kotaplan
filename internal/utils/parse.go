package utils

import "regexp"

func LabelsFromString(rawLabels string) (labels map[string]string) {
	labelsPattern := regexp.MustCompile(`(?P<key>\w+)=(?P<value>[^,]+)`)

	if !labelsPattern.MatchString(rawLabels) {
		return map[string]string{}
	}

	matches := labelsPattern.FindAllStringSubmatch(rawLabels, -1)
	labels = make(map[string]string, len(matches))
	for _, match := range matches {
		labels[match[1]] = match[2]
	}

	return
}
