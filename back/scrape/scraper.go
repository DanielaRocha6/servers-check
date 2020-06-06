package scrape

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func FindWebsiteTitle(titleChannel chan string, url string) {
	// Default case
	pageTitle := "Title not found."

	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println("Err:", err.Error())
	} else {
		doc.Find("head title").Each(func(index int, item *goquery.Selection) {
			title, err := item.Html()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				pageTitle = title
			}
		})
	}
	titleChannel <- pageTitle
}
func FindWebsiteLogo(logoChannel chan string, url string) {
	// Default case
	logoURL := "Logo not found."

	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		possibleImg := make(map[string]string)
		doc.Find("head link").Each(func(index int, item *goquery.Selection) {
			linkTag := item
			href, _ := linkTag.Attr("href")
			rel, _ := linkTag.Attr("rel")
			if strings.Contains(rel, "icon") {
				possibleImg[rel] = href
				logoURL = href
			}
		})
		if len(possibleImg) > 1 {
			for k, v := range possibleImg {
				if !strings.Contains(v, "http") || !strings.Contains(v, "www") {
					possibleImg[k] = ""
				}
				if strings.Contains(v, ".png") {
					logoURL = v
					break
				}
			}
		} else {
			doc.Find("body img").Each(func(index int, item *goquery.Selection) {
				linkTag := item
				imgType, _ := linkTag.Attr("type")
				srcType, _ := linkTag.Attr("src")
				if strings.Contains(imgType, "icon") {
					possibleImg[imgType] = srcType
					logoURL = srcType
				}
			})
			if len(possibleImg) > 1 {
				for k, v := range possibleImg {
					if !strings.Contains(v, "http") {
						possibleImg[k] = ""
					}
					if strings.Contains(v, ".png") {
						logoURL = v
						break
					}
				}
			}
		}
	}
	logoChannel <- logoURL
}
func GetHTML(url string, htmlChannel chan []string, secondTry bool) {
	titleChannel := make(chan string)
	logoChannel := make(chan string)

	timeout := time.Duration(4 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	_, err := client.Get(url)
	if err != nil {
		if !secondTry {
			GetHTML(url, htmlChannel, true)
		}
		fmt.Println("HTML:", err.Error())
		htmlChannel <- []string{"Logo not found.", "Title not found."}
	} else {
		go FindWebsiteTitle(titleChannel, url)
		go FindWebsiteLogo(logoChannel, url)

		logoURL := <-logoChannel
		pageTitle := <-titleChannel
		htmlChannel <- []string{logoURL, pageTitle}
	}

}
