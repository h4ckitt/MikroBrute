package util

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var saltRegex = regexp.MustCompile(`('\\[0-9]{3}' \+ (document.login.password.value|[0-9]{4}) \+ '(?:\\[0-9]{3})+')`)

func ExtractSalt(body, num string) string {
	//fmt.Println(strings.Contains(body, "document.login.password.value"))
	if match := saltRegex.MatchString(body); !match {
		fmt.Println("No Match")
		return ""
	}

	fmt.Println("Matches")
	salt := saltRegex.FindStringSubmatch(body)

	if len(salt) < 2 {
		return ""
	}

	return strings.Replace(salt[1], salt[2], num, 1)
	/*s := strings.Replace(salt[1], salt[2], num, 1)

	fmt.Println(s)

	return strconv.Quote(strings.Join(strings.Split(strings.Trim(s, "'"), "'"), ""))*/
}

func Parse(str string) string {
	//fmt.Println(str)
	reg := regexp.MustCompile(`\\[0-9]{3}(\+([0-9]{4})\+)(?:\\[0-9]{3})+`)

	runes := make([]rune, 0)

	submatches := reg.FindStringSubmatch(str)

	newStr := strings.TrimLeft(strings.Replace(str, submatches[1], "", 1), "\\")

	octNums := strings.Split(newStr, "\\")

	for _, elem := range octNums {
		num, err := strconv.ParseInt(elem, 8, 64)

		if err != nil {
			log.Fatalln(err)
		}

		runes = append(runes, rune(num))
	}

	return fmt.Sprintf("%s%s%s", string(runes[0]), submatches[2], string(runes[1:]))
	//'\150' + 4020 + '\002\036\063\367\305\347\024\262\316\021\157\252\005\117\062\122'

}

func Saltify(salt string) string {

	replacer := strings.NewReplacer("'", "", " ", "")

	fullMatch := replacer.Replace(salt)

	raw := Parse(fullMatch)
	//fmt.Println("Raw.. ", raw)
	return raw
}

func CheckForSuccess(body string) bool {
	//fmt.Println(body)
	if strings.Contains(body, "logged") {
		return true
	}

	return false
}
