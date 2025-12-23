package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ParsePage(dom, l string) (*http.Response, error) {

	client := &http.Client{}
	u, err := url.Parse(l)
	if err != nil {
		panic(err)
	}
	udom, err := url.Parse(dom)
	if err != nil {
		panic(err)
	}
	if udom.Hostname() != u.Hostname() {
		errorDomain := errors.New("Be careful! not same domain: %s " + u.Hostname())
		return nil, errorDomain
	}
	res, err := client.Get(l)
	if err != nil {
		panic(err)
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	return res, nil
}
func formatURL(dom, l string) string {
	if l[0] == "/"[0] {
		return dom + l
	}
	return l
}

func checkAtrr(dom, l string, doc *html.Node, links *[]string, v *map[string]bool) {
	ls := links
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if _, prs := (*v)[a.Val]; !prs && a.Key == "href" && a.Val != "#" {
					*ls = append(*ls, a.Val)
					(*v)[a.Val] = true
					nUrl := formatURL(dom, a.Val)
					fmt.Println(nUrl)
					*ls = append(*ls, nUrl)
					Scrape(dom, nUrl, v)
				}

			}

		}
	}
	if len(*ls) == 0 {
		(*v)[l] = false
	}
}
func Scrape(dom, l string, v *map[string]bool) {
	res, err := ParsePage(dom, l)
	if err != nil {
		fmt.Println(err)
		return
	}
	doc, err := html.Parse(res.Body)

	if err != nil {
		panic(err)
	}
	links := make([]string, 0)
	checkAtrr(dom, l, doc, &links, v)
}

func main() {
	arg := os.Args[1]
	v := make(map[string]bool)
	Scrape(arg, arg, &v)
	fmt.Println("\nDead link:")
	for k := range v {
		if v[k] == false {
			fmt.Println(k)
		}

	}
}
