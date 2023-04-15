package vocabulary

import (
	"strings"
	"unicode"
)

//var number_word regexp.Regexp = *regexp.MustCompile(`((^\d+)|(\d+$))`)

func TokenizeText(texts ...string) []string {

	text := strings.Join(texts, " ")

	text = strings.ReplaceAll(text, " ", " ")
	text = strings.ReplaceAll(text, "(", " ")
	text = strings.ReplaceAll(text, ")", " ")

	list := strings.Split(text, " ")
	for l := range list {
		list[l] = strings.ToLower(list[l])
	}
	return list
}

func CleanWord(word string, lang Language) string {
	word = TrimWord([]rune(word))

	for s := range Suffixes[lang] {
		word = strings.TrimSuffix(word, Suffixes[lang][s])
	}
	return word
}

func TrimWord(word []rune) string {
	var prefix int = 0
	var suffix int = len(word)

	for x := 0; x < len(word); x++ {
		if !unicode.IsLetter(word[x]) {
			prefix++
		} else {
			break
		}
	}

	for x := len(word) - 1; x >= 0; x-- {
		if suffix == prefix {
			break
		}
		if !unicode.IsLetter(word[x]) {
			suffix--
		} else {
			break
		}
	}
	//fmt.Printf("changed %s to %s\n", string(word), string(word[prefix:suffix]))
	return string(word[prefix:suffix])
}

/*
EditDistance compares two strings, and calculates the edit distance (the number of needed changes
for one string to become equal to the other)
*/
func EditDistance(s, t []rune) int {

	if len(s) == 0 {
		return len(t)
	}
	if len(t) == 0 {
		return len(s)
	}

	m := make([][]int, len(s)+1)
	for i := range m {
		m[i] = make([]int, len(t)+1)
	}

	for i := 1; i < len(s)+1; i++ {
		m[i][0] = i
	}

	for j := 1; j < len(t)+1; j++ {
		m[0][j] = j
	}

	for i := 1; i < len(s)+1; i++ {
		for j := 1; j < len(t)+1; j++ {
			m[i][j] = min(
				[3]int{
					m[i-1][j-1] + boolInt(s[i-1] != t[j-1]),
					m[i-1][j] + 1,
					m[i][j-1] + 1,
				},
			)
		}
	}

	return m[len(s)][len(t)]
}

func min(numbers [3]int) int {

	if numbers[0] < numbers[1] && numbers[0] < numbers[2] {
		return numbers[0]
	}
	if numbers[1] < numbers[2] {
		return numbers[1]
	}
	return numbers[2]
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
