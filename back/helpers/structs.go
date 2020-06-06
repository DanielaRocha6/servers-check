package helpers

type DomainInfo struct {
	Servers            []Server
	Servers_Changed    bool
	Ssl_Grade          string
	Previous_ssl_grade string
	Logo               string
	Title              string
	Is_down            bool
	LastModifiedDate   string
}
type Server struct {
	Address   string
	Ssl_grade string
	Country   string
	Owner     string
}

type Item struct {
	Url  string
	Info DomainInfo
}

type Domains struct {
	Items []Item
}
