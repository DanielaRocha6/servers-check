package helpers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type ErrorStruct struct {
	Error string
}

func gradeToNum(grade string) int {
	var ans int
	switch grade {
	case "A+":
		ans = 8
	case "A":
		ans = 7
	case "A-":
		ans = 6
	case "B":
		ans = 5
	case "C":
		ans = 4
	case "D":
		ans = 3
	case "E":
		ans = 2
	case "F":
		ans = 1
	default:
		ans = 999
	}
	return ans
}

func LessThan(grade1 string, grade2 string) bool {
	a := gradeToNum(grade1)
	b := gradeToNum(grade2)
	return a < b
}

func MoreThanAnHour(altered DomainInfo) bool {
	loc, _ := time.LoadLocation("UTC")
	lastHour := altered.LastModifiedDate
	t, err := time.Parse(time.RFC3339, lastHour)
	if err != nil {
		fmt.Println("Date error", err.Error())
	}
	now := time.Now().In(loc)
	delta := now.Sub(t.In(loc)).Hours()
	if delta >= 1 {
		return true
	}
	return false
}

func GetNowString() string {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	temp := strings.Split(now.String(), " ")
	return temp[0] + "T" + temp[1] + "Z"
}

func AnswerWithErr(err string, ctx *fasthttp.RequestCtx) {
	displayedError := ErrorStruct{Error: err}
	ansJSON, _ := json.Marshal(displayedError)
	fmt.Fprintf(ctx, "%s\n", ansJSON)
}

func ValidateURL(url string) bool {
	var regexURL = regexp.MustCompile(`^((https?|ftp|smtp):\/\/)?(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,4}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	isURL := false
	for _, match := range regexURL.FindAllString(url, -1) {
		isURL = true
		fmt.Println(match, "is a valid URL")
	}
	return isURL
}

func CheckLogoURL(logo string, servers []Server, c chan string) {
	if strings.Contains(logo, "http") || strings.Contains(logo, "www") || len(servers) == 0 {
		c <- logo
	} else {
		for i := range servers {
			if strings.Contains(servers[i].Address, ".") && strings.Contains(logo, "/") {
				if strings.HasPrefix(logo, "/") {
					c <- servers[i].Address + logo
				} else {
					c <- servers[i].Address + "/" + logo
				}

				break
			} else {
				c <- logo
			}
		}
	}
}

func ServersHadChanged(past []Server, actual []Server) bool {
	if len(actual) != len(past) {
		return true
	}
	// // Couldn't find servers
	// if len(actual) == 0 {
	// 	return false
	// }
	count := 0
	for i := range actual {
		serverA := actual[i]
		for j := range past {
			serverP := past[j]
			if serverA.Address == serverP.Address && serverA.Ssl_grade == serverP.Ssl_grade {
				count++
				if count == len(actual) {
					return false
				}
			}
		}
	}
	return true
}
