package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"../helpers"

	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnectDb() {
	res, err := sql.Open("postgres", "postgresql://root@localhost:26257/domains?sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to DB", err.Error())
		createTables()
	} else {
		db = res
		fmt.Println("Successfully connected to db!")
	}

}

func createTables() {
	if _, err := db.Exec(
		"CREATE TABLE Domains (" +
			"url STRING PRIMARY KEY," +
			"serversChanged BOOL," +
			"sslGrade STRING NOT NULL," +
			"previousSslGrade STRING," +
			"logo STRING NOT NULL," +
			"title STRING NOT NULL," +
			"isDown BOOL NOT NULL," +
			"checkDate TIMESTAMPTZ NOT NULL" +
			");"); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Domains table successfully created!")
	}
	if _, err := db.Exec(
		"CREATE TABLE Servers (" +
			"domainId STRING REFERENCES domains(url) ON DELETE CASCADE," +
			"address STRING PRIMARY KEY," +
			"sslGrade STRING NOT NULL," +
			"country STRING NOT NULL," +
			"owner STRING" +
			");"); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Servers table successfully created!")
	}
}

func CheckIfExisted(url string) (helpers.DomainInfo, bool) {
	query := fmt.Sprintf("select serverschanged, sslgrade,previoussslgrade,logo,title,checkdate  FROM domains WHERE url='%s';", url)
	rows, err := db.Query(query)
	var dom helpers.DomainInfo
	if err != nil {
		fmt.Println("Error in query finding domain", url, err.Error())
	} else {
		defer rows.Close()
		for rows.Next() {
			var serverschanged bool
			var sslgrade string
			var previoussslgrade string
			var logo string
			var title string
			var checkdate string
			if err := rows.Scan(&serverschanged, &sslgrade, &previoussslgrade, &logo, &title, &checkdate); err != nil {
				log.Fatal(err.Error())
			}
			dom.Servers_Changed = serverschanged
			dom.Ssl_Grade = sslgrade
			dom.Previous_ssl_grade = previoussslgrade
			dom.Logo = logo
			dom.Title = title
			dom.LastModifiedDate = checkdate
			return dom, true
		}
	}
	return dom, false
}

func FindServersByDomain(url string) []helpers.Server {
	servers := make([]helpers.Server, 0, 0)
	query := fmt.Sprintf("select address,sslgrade,country,owner from servers where domainid='%s';", url)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error in query finding servers of", url)
	} else {
		defer rows.Close()
		for rows.Next() {
			var address string
			var sslgrade string
			var country string
			var owner string
			if err := rows.Scan(&address, &sslgrade, &country, &owner); err != nil {
				log.Fatal(err)
			}
			server := helpers.Server{
				Address:   address,
				Ssl_grade: sslgrade,
				Country:   country,
				Owner:     owner,
			}
			servers = append(servers, server)
		}
	}
	return servers
}

func deleteServersByDomain(url string, c chan error) {
	query := fmt.Sprintf("DELETE FROM servers WHERE domainId='%s';", url)
	if res, err := db.Exec(
		query); err != nil {
		c <- err
		fmt.Println("Err1:", err.Error())
	} else {
		c <- nil
		rows, _ := res.RowsAffected()
		fmt.Println("DELETED", rows, "rows")
	}
}

func insertIntoServers(servers []helpers.Server, domainURL string, c chan error, withChan bool) {
	var server helpers.Server
	for i := range servers {

		server = servers[i]
		address := server.Address
		sslGrade := server.Ssl_grade
		country := server.Country
		owner := server.Owner

		query := fmt.Sprintf("INSERT INTO servers(domainId, address, sslGrade, country, owner) VALUES ('%s', '%s' ,'%s' ,'%s' , '%s');", domainURL, address, sslGrade, country, owner)
		if res, err := db.Exec(
			query); err != nil {
			fmt.Println("Err2:", err.Error())
			if withChan {
				c <- err
			}
		} else {
			rows, _ := res.RowsAffected()
			fmt.Println("INSERTED", rows, "servers rows")
			if withChan {
				c <- nil
			}
		}
	}
}

func alterDomains(url string, altered helpers.DomainInfo, original helpers.DomainInfo) helpers.DomainInfo {
	finalDomain := original
	finalDomain.Previous_ssl_grade = altered.Ssl_Grade
	c := make(chan error)

	pastServers := FindServersByDomain(url)
	if helpers.ServersHadChanged(pastServers, original.Servers) {
		fmt.Println("Servers changed")

		go deleteServersByDomain(url, c)
		deleteErr := <-c
		go insertIntoServers(original.Servers, url, c, true)

		loc, _ := time.LoadLocation("UTC")
		now := time.Now().In(loc)
		temp := strings.Split(now.String(), " ")
		checkDate := temp[0] + " " + temp[1]

		finalDomain.Servers_Changed = true
		finalDomain.LastModifiedDate = checkDate

		//Update query
		query := fmt.Sprintf("UPDATE domains SET (serverschanged, previoussslgrade, checkdate) =(%v,'%s', TIMESTAMPTZ '%s') WHERE url = '%s';", finalDomain.Servers_Changed, finalDomain.Previous_ssl_grade, finalDomain.LastModifiedDate, url)
		if res, err := db.Exec(
			query); err != nil {
			fmt.Println("Err3:", err.Error())
		} else {
			rows, _ := res.RowsAffected()
			fmt.Printf("Updated %v rows", rows)
		}

		insertErr := <-c
		if deleteErr != nil || insertErr != nil {
			fmt.Println("Error in servers Update.")
		}
	} else {
		fmt.Println("Servers didn't change")
	}
	return finalDomain
}

func insertIntoDomains(url string, dom helpers.DomainInfo) helpers.DomainInfo {
	c := make(chan error)
	sslGrade := dom.Ssl_Grade
	logo := dom.Logo
	title := dom.Title
	isDown := dom.Is_down
	serverschanged := dom.Servers_Changed
	previoussslgrade := dom.Previous_ssl_grade
	checkDate := helpers.GetNowString()

	if len(dom.Servers) > 0 {
		query := fmt.Sprintf("INSERT INTO domains(url, sslGrade, logo, title, isDown, checkDate, serverschanged, previoussslgrade) VALUES ('%s', '%s', '%s', '%s', '%v', TIMESTAMPTZ '%s','%v', '%s');", url, sslGrade, logo, title, isDown, checkDate, serverschanged, previoussslgrade)
		if res, err := db.Exec(
			query); err != nil {
			fmt.Println("Err4:", err.Error())
		} else {
			rows, _ := res.RowsAffected()
			fmt.Printf("Inserted %v domains rows", rows)
			insertIntoServers(dom.Servers, url, c, false)
		}
	}
	return dom
}

func AnswerDomain(url string, dom helpers.DomainInfo) helpers.DomainInfo {
	domi, existed := CheckIfExisted(url)
	if existed {
		return alterDomains(url, domi, dom)
	}
	return insertIntoDomains(url, dom)
}

func GetAllDomains() helpers.Domains {
	var allDomains helpers.Domains
	query := "select *  FROM domains;"
	rows, err := db.Query(query)
	var dom helpers.DomainInfo
	if err != nil {
		fmt.Println("Error in query finding domains")
	} else {
		domains := make([]helpers.Item, 0, 0)
		defer rows.Close()
		for rows.Next() {
			var url string
			var serverschanged bool
			var sslgrade string
			var previoussslgrade string
			var logo string
			var title string
			var isDown bool
			var checkdate string
			if err := rows.Scan(&url, &serverschanged, &sslgrade, &previoussslgrade, &logo, &title, &isDown, &checkdate); err != nil {
				log.Fatal(err)
			}
			servers := FindServersByDomain(url)

			dom.Servers_Changed = serverschanged
			dom.Ssl_Grade = sslgrade
			dom.Previous_ssl_grade = previoussslgrade
			dom.Title = title
			dom.LastModifiedDate = checkdate
			dom.Servers = servers
			dom.Logo = logo

			item := helpers.Item{
				Url:  url,
				Info: dom,
			}
			domains = append(domains, item)
		}
		allDomains.Items = domains
	}
	return allDomains
}
