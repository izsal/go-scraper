package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type item struct {
	Name   string `json:"name"`
	Price  string `json:"price"`
	ImgUrl string `json:"imgUrl"`
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("j2store.net"),
	)

	var items []item

	c.OnHTML("div[itemprop=itemListElement]", func(e *colly.HTMLElement) {
		item := item{
			Name:   e.ChildText("h2.product-title"),
			Price:  e.ChildText("div.sale-price"),
			ImgUrl: e.ChildAttr("img", "src"),
		}

		items = append(items, item)
	})

	c.OnHTML("[title=Next]", func(h *colly.HTMLElement) {
		next_page := h.Request.AbsoluteURL(h.Attr("href"))
		c.Visit(next_page)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
	})

	c.Visit("https://j2store.net/demo/index.php/shop")
	
	content, err := json.Marshal(items)

	if err != nil {
		fmt.Println(err.Error())
	}

	os.WriteFile("products.json", content, 0644)
}
