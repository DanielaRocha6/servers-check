package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"./db"
	"./helpers"
	"./scrape"

	"./connections"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

const sslURL = "https://api.ssllabs.com/api/v3/analyze?host="

var domains = make(map[string]helpers.DomainInfo)

func saveSSLInfo(result map[string]interface{}, url string, secondTry bool) *helpers.DomainInfo {
	serversChannel := make(chan []helpers.Server)
	minSsl := make(chan string)

	go getServersInfo(result, serversChannel, minSsl)

	servers := <-serversChannel
	sslGrade := <-minSsl

	var newDomain = helpers.DomainInfo{
		Servers:   servers,
		Ssl_Grade: sslGrade,
	}
	if len(servers) < 1 {
		if !secondTry {
			saveSSLInfo(result, url, true)
		} else {
			return &newDomain
		}
	}
	return &newDomain
}

func getServersInfo(result map[string]interface{}, serversChannel chan []helpers.Server, minSsl chan string) {
	if result["endpoints"] != nil {
		servers := make([]helpers.Server, 0, 0)
		var server helpers.Server
		min := ""
		count := 0
		for _, item := range result["endpoints"].([]interface{}) {

			if str, ok := item.(map[string]interface{})["ipAddress"].(string); ok {
				whoIsServer := connections.RunWhoIs(str)
				owner := strings.TrimSpace(strings.Split(strings.Split(whoIsServer, "OrgName:")[1], "OrgId")[0])
				country := strings.TrimSpace(strings.Split(strings.Split(whoIsServer, "Country:")[1], "RegDate")[0])
				server.Address = str
				server.Owner = owner
				server.Country = country

			}
			if str, ok := item.(map[string]interface{})["grade"].(string); ok {
				if count == 0 {
					min = str
				}
				server.Ssl_grade = str
				if helpers.LessThan(str, min) {
					min = str
				}
				count++
			}
			servers = append(servers, server)
		}
		serversChannel <- servers
		minSsl <- min
	} else {
		serversChannel <- make([]helpers.Server, 0, 0)
		minSsl <- ""
	}
}

func findDomainInfo(url string) helpers.DomainInfo {

	htmlChannel := make(chan []string)
	c := make(chan bool)

	go scrape.GetHTML("http://"+url, htmlChannel, false)
	go connections.IsDomainDown(url, c)

	sslResult := connections.Fetch(sslURL + url)
	dom := *saveSSLInfo(sslResult, url, false)

	info := <-htmlChannel
	logo, title := info[0], info[1]

	logoChan := make(chan string)
	go helpers.CheckLogoURL(logo, dom.Servers, logoChan)

	isDown := <-c

	dom.LastModifiedDate = helpers.GetNowString()
	dom.Is_down = isDown

	dom.Title = title
	dom.Servers_Changed = false
	dom.Previous_ssl_grade = ""
	logoF := <-logoChan
	dom.Logo = logoF
	return db.AnswerDomain(url, dom)
}

func allDomains(ctx *fasthttp.RequestCtx) {
	serverJSON, _ := json.Marshal(db.GetAllDomains())
	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(ctx, "%s\n", serverJSON)
}

func domainEndpoint(ctx *fasthttp.RequestCtx) {

	if url, ok := ctx.UserValue("domain").(string); ok {
		if helpers.ValidateURL(url) {
			printedError := false
			foundDom, existed := db.CheckIfExisted(url)

			var domainInfo helpers.DomainInfo
			domainInfo.Servers_Changed = false
			domainInfo.Previous_ssl_grade = ""

			if existed {
				if !strings.Contains(foundDom.Logo, "http") || !strings.Contains(foundDom.Logo, "/") {
					logoChannel := make(chan string)
					go scrape.FindWebsiteLogo(logoChannel, "http://"+url)
					logoURL := <-logoChannel
					domainInfo.Logo = logoURL
				}
				if strings.Contains(foundDom.Title, "not found") {
					titleChannel := make(chan string)
					go scrape.FindWebsiteTitle(titleChannel, "http://"+url)
					titleF := <-titleChannel
					domainInfo.Title = titleF
				}
				if !helpers.MoreThanAnHour(foundDom) {
					downChannel := make(chan bool)
					go connections.IsDomainDown(url, downChannel)
					servers := db.FindServersByDomain(url)

					if len(servers) == 0 {
						sslResult := connections.Fetch(sslURL + url)
						dom := *saveSSLInfo(sslResult, url, false)
						dom.Servers_Changed = false
						if len(dom.Servers) != 0 {
							domainInfo = db.AnswerDomain(url, dom)
						} else {
							helpers.AnswerWithErr("Couldn't find any servers. Try again!", ctx)
							printedError = true
						}
					} else {
						domainInfo = foundDom
						domainInfo.Servers = servers
					}

					isDown := <-downChannel
					domainInfo.Is_down = isDown

				} else {
					domainInfo = findDomainInfo(url)
				}
			} else {
				dom := findDomainInfo(url)
				if len(dom.Servers) > 0 {
					domainInfo = dom
				} else {
					helpers.AnswerWithErr("Couldn't find servers. Try again!", ctx)
					printedError = true
				}
			}
			if !printedError {
				serverJSON, err := json.Marshal(domainInfo)
				fmt.Fprintf(ctx, "%s\n", serverJSON)

				ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
				ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

				if err != nil {
					helpers.AnswerWithErr("Couldn't parse JSON", ctx)
				}
			}
		} else {
			helpers.AnswerWithErr("URL is not valid", ctx)
		}
	} else {
		helpers.AnswerWithErr("Param is not valid", ctx)
	}
}

func main() {
	router := fasthttprouter.New()
	db.ConnectDb()

	router.GET("/checkDomain/:domain", domainEndpoint)
	router.GET("/allDomains", allDomains)
	log.Fatal(fasthttp.ListenAndServe(":8090", router.Handler))
}
