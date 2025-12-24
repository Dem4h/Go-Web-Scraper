package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Yellow = "\033[33m"
var Cyan = "\033[36m"
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

func checkAtrr(dom, l string, doc *html.Node, links *[]string, v *map[string]bool, wg *sync.WaitGroup) {
	defer wg.Done()
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
						wg.Add(1)
						go Scrape(dom, nUrl, v, wg)
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
func Scrape(dom, l string, v *map[string]bool, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	res, err := ParsePage(dom, l)
	if err != nil {
		fmt.Println(err)
		return
	}
	doc, err := html.Parse(res.Body)
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}
	links := make([]string, 0)
	checkAtrr(dom, l, doc, &links, v, wg)
}

func main() {
	arg := os.Args[1]
	v := make(map[string]bool)
	if arg[len(arg)-1] == "/"[0] {
		arg = arg[:len(arg)-1]
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go Scrape(arg, arg, &v, &wg)
	wg.Wait()
	fmt.Println(Cyan + "\nDead link:")
	for k := range v {
		if v[k] == false {
			fmt.Println(k)
		}

	}
}
