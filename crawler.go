package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"strconv"

	"github.com/gocolly/colly"
)

func main() {
	visdic := map[string]bool{}
	file, err := os.OpenFile("ptbrclosure.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	buf := bufio.NewScanner(file)
	for buf.Scan() {
		visdic[strings.Split(buf.Text(), "\t")[0]] = true
	}
	defer file.Close()

	reader := bufio.NewReader(os.Stdin)
	for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
		var palavra string
		mainCollector := colly.NewCollector()
		hyperCollector := mainCollector.Clone()

		mainCollector.OnHTML(".line > a", func(e *colly.HTMLElement) {
			url := strings.Split(e.Attr("href"), "/")
			if visdic[url[len(url)-1]] {
				return
			}
			hyperCollector.Visit(e.Attr("href"))
		})
		mainCollector.OnHTML(".pagination > li > a", func(e *colly.HTMLElement) {
			if e.Text == "»" {

				e.Request.Visit(e.Attr("href"))
			}
		})
		hyperCollector.OnHTML(".line.quote > p", func(e *colly.HTMLElement) {

			fmt.Println("palavra: ", palavra)
			fmt.Println(e.Text)

			split := strings.Split(e.Text, " ")

			text, _ := reader.ReadString('\n')
			n, err := strconv.Atoi(text[:len(text)-1])
			if err != nil {
				return
			}
			_, err = file.WriteString(fmt.Sprintln(palavra + "\t" + split[n] + "\t" + e.Text))
			if err != nil {
				panic(err)
			}
		})
		hyperCollector.OnRequest(func(r *colly.Request) {
			sp := strings.Split(r.URL.String(), "/")

			palavra = sp[len(sp)-1]
			log.Println("Hyper: " + r.URL.String())
		})
		mainCollector.OnRequest(func(r *colly.Request) {
			log.Println("Main: " + r.URL.String())
		})
		mainCollector.Visit("https://dicionario.aizeta.com/verbetes/substantivo/" + string(letter))
	}
}
