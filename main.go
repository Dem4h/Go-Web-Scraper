package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func isSameDom(dom, l string) bool {

	u, err := url.Parse(l)
	if err != nil {
		panic(err)
	}
	udom, err := url.Parse(dom)
	if err != nil {
		panic(err)
	}
	if udom.Hostname() != u.Hostname() {
		fmt.Println(Red + "link found with different domain: " + u.Hostname())
		return false
	}
	return true
}
func ParsePage(dom, l string) (*http.Response, error) {

	client := &http.Client{}
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
					(*v)[a.Val] = true
					nUrl := formatURL(dom, a.Val)
					fmt.Println(White + nUrl)
					*ls = append(*ls, nUrl)
					if isSameDom(dom, nUrl) {

						Scrape(dom, nUrl, v)
					}
				}

			}

		}
	}
	if len(*ls) == 0 {
		fmt.Println(Cyan + "Deadlink found : " + l)
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
	if arg[len(arg)-1] == "/"[0] {
		arg = arg[:len(arg)-1]
	}
	Scrape(arg, arg, &v)
	fmt.Println(Cyan + "\nDead link:")
	for k := range v {
		if v[k] == false {
			fmt.Println(k)
		}

	}
}
