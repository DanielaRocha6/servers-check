package connections

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"time"
)

func Fetch(url string) map[string]interface{} {

	var result map[string]interface{}
	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		defer resp.Body.Close()
		res, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
		}
		json.Unmarshal([]byte(res), &result)
	}
	return result
}

func RunWhoIs(ipAddr string) string {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return ipAddr + " is not a valid address"
	}
	ipAddr = ip.String()
	whoisServer := "whois.arin.net"
	whois := whoIs("-h", whoisServer, ipAddr)

	return whois.String()
}

func whoIs(args ...string) bytes.Buffer {
	var out bytes.Buffer
	cmd := exec.Command("whois", args...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Run()
	return out
}

func IsDomainDown(url string, c chan bool) {
	host := url
	_, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		fmt.Println(url, "is down")
		c <- true
	}
	fmt.Println(url, "is up")
	c <- false
}
