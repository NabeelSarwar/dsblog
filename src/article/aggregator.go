package article

import (
	"bytes"
	"encoding/json"
	// for app engine
	// "fmt"
	"html/template"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Aggregator struct {
	articles   []Article
	TitleToUrl map[string]string
	UrlToTitle map[string]string
}

// This function aggregates all the functions and creates an Aggregator object that can be
// used for templating on the main page.
func Aggregate() Aggregator {
	files, err := ioutil.ReadDir("src/static/articles/")

	var articles []Article
	var TitleToUrl = make(map[string]string)
	var UrlToTitle = make(map[string]string)

	if err != nil {
		log.Panic(err)
	}

	for _, v := range files {
		// load the json from the file
		// these things are file values
		if v.IsDir() {
			continue
		}
		var art Article
		fileContent, err := ioutil.ReadFile("src/static/articles/" + v.Name())

		if err != nil {
			log.Panic(err)
		}

		err = json.Unmarshal(fileContent, &art)
		if err != nil {
			log.Panic(err)
		}
		articles = append(articles, art)
		TitleToUrl[art.Title] = art.Url
		UrlToTitle[art.Url] = art.Title
	}
	// for go app engine
	// fmt.Println("the number of articles is:")
	// fmt.Println(len(articles))
	return Aggregator{articles, TitleToUrl, UrlToTitle}
}

/*
Need to use template.HTML to escape properly.
This returns a short version of each article ready to be displayed
*/
func (agg Aggregator) DisplayArticle(a Article) template.HTML {
	var url = agg.TitleToUrl[a.Title]
	var display = template.HTML(`<div> <article><h2> <a href="/article/` + url + `">` + a.Title + `</a> </h4>
	<div class="row">
	<div class="group1 col-sm-4 col-md-4">
	<span class="glyphicon glyphicon-bookmark"></span>` + strings.Join(a.Tags, ", ") +

	`</div>
	<div class="group2 col-sm-8 col-md-8">
	<span class="glyphicon glyphicon-time"></span>` + a.Date.Format(time.RFC822) + `
	</div>
	</div>
	<hr>

	<br />
	<p>` + a.LimitedContent + `</p>
	<p class="text-right">
	<a href="/article/` + url + `"class="text-right">
	continue reading...
	</a>
	</p>
	</hr></article> </div>`)
	return display
}

func (agg Aggregator) DisplayArticleAll() template.HTML {
	var buffer bytes.Buffer
	for _, value := range agg.articles {
		buffer.WriteString(string(agg.DisplayArticle(value)))
	}
	return template.HTML(buffer.String())
}
