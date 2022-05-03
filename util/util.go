package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var saltRegex = regexp.MustCompile(`('\\[0-9]{3}'( \+ document.login.password.value|[0-9]{4} \+ )'(?:\\[0-9]{3})+')`)

func ExtractSalt(body, num string) string {
	//fmt.Println(strings.Contains(body, "document.login.password.value"))
	if match := saltRegex.MatchString(body); !match {
		fmt.Println("No Match")
		return ""
	}

	salt := saltRegex.FindStringSubmatch(body)

	if len(salt) < 2 {
		return ""
	}

	s := strings.Replace(salt[1], salt[2], num, 1)

	return strconv.Quote(strings.Join(strings.Split(strings.Trim(s, "'"), "'"), ""))
}

func CheckForSuccess(body string) bool {
	fmt.Println(body)
	if strings.Contains(body, "logged") {
		return true
	}

	return false
}
