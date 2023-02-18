package datatype

import "regexp"

type RegexpPair struct {
	Value string
	Match bool
}

// 正则表达式分割字符串
func RegexpSplit(expr string, str string) []RegexpPair {
	var data []RegexpPair

	if m, err := regexp.Compile(expr); err == nil {
		for {
			loc := m.FindStringIndex(str)

			if loc == nil {
				if str != "" {
					data = append(data, RegexpPair{Value: str})
				}

				break
			}

			if loc[0] > 0 {
				data = append(data, RegexpPair{Value: str[:loc[0]]})
			}

			match := str[loc[0]:loc[1]]
			data = append(data, RegexpPair{Value: match, Match: true})

			str = str[loc[1]:]
		}
	} else {
		data = append(data, RegexpPair{Value: str})
	}

	return data
}

// 字符串反转
func ReverseString(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}
