package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	// "io/ioutil"
)

// Entry struct
type Entry struct {
	ipaddress     string
	identityCheck string
	user          string
	datetime      string
	*Request
	status    string
	size      string
	referer   string
	userAgent string
	cookie    string
}

// Request struct
type Request struct {
	method   string
	path     string
	protocol string
}

func newEntry(ipaddress string, identityCheck string, user string, datetime string, request string, status string, size string, referer string, userAgent string, cookie string) *Entry {
	e := new(Entry)
	e.ipaddress = ipaddress
	e.identityCheck = identityCheck
	e.user = user
	e.datetime = datetime
	e.Request = newRequest(strings.Split(request, " "))
	e.status = status
	e.size = size
	e.referer = referer
	e.userAgent = userAgent
	e.cookie = cookie
	return e
}

func newRequest(req []string) *Request {
	r := new(Request)
	r.method = req[0]
	r.path = req[1]
	r.protocol = req[2]
	return r
}

func main() {
	var fp *os.File
	var err error

	fmt.Printf(">> read file: %s\n", os.Args[1])
	fp, err = os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	regEx := `(\S+)\s+(\S+)\s+(\S+)\s+\[(.*?)\]\s+"(.*?)"\s+(\S+)\s+(\S+)\s+"(.*?)"\s+"(.*?)"\s+"(.*?)"`
	r := regexp.MustCompile(regEx)

	domain := "http://localhost/"
	client := &http.Client{}

	reader := bufio.NewReaderSize(fp, 4096)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		ret := r.FindStringSubmatch(string(line))
		//fmt.Println(len(ret))
		/*
		   for i, v := range ret {
		           fmt.Println(i, v)
		   }
		*/

		//var entry *Entry = newEntry(ret[1], ret[2], ret[3], ret[4], ret[5], ret[6], ret[7], ret[8], ret[9], ret[10])
		entry := newEntry(ret[1], ret[2], ret[3], ret[4], ret[5], ret[6], ret[7], ret[8], ret[9], ret[10])
		fmt.Println(entry.Request.path)

		req, _ := http.NewRequest("GET", domain+entry.Request.path, nil)
		//req.SetBasicAuth("user", `pass`)
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		println(res.Status)
		println(res.StatusCode)

	}
}

func splitline(line string) {
	params := strings.Split(line, " ")
	fmt.Println(len(params))
	for i := 0; i < len(params); i++ {
		s := params[i][0:1]
		e := params[i][len(params[i])-1 : len(params[i])]

		if s == "\"" || s == "[" {
			fmt.Println(s, e)
		}

		fmt.Println(params[i])
	}
}
