package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
)

const Synonyms = "Synonyms:"
const Antonyms = "Antonyms:"
const TypeOf = "Type of:"

func main() {
	http.HandleFunc("/", indexHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	key, ok := r.URL.Query()["word"]
	if !ok {
		return
	}

	word := getWord(key[0])
	data := formatToJson(word)
	json, _ := json.Marshal(data)
	_, err := fmt.Fprint(w, string(json))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getWord(word string) io.Reader {
	if word == "" {
		word = "necropolis"
	}
	var url = fmt.Sprintf("https://www.vocabulary.com/dictionary/definition.ajax?search=%s&lang=en", word)
	response, err := http.Get(url)
	if err != nil {
		return nil
	} else {
		//data, _ := ioutil.ReadAll(response.Body)
		return response.Body
	}
	return nil
}

func formatToJson(htmlString io.Reader) Vocabulary {
	vocab := Vocabulary{}

	doc, err := goquery.NewDocumentFromReader(htmlString)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".dynamictext").Each(func(i int, selection *goquery.Selection) {
		vocab.WordTitle = selection.Text()
	})

	doc.Find(".section").Each(func(i int, selection *goquery.Selection) {
		vocab.DefinitionShort = selection.Find(".short").Text()
		vocab.DefinitionLong = selection.Find(".long").Text()
	})

	doc.Find(".ordinal").Each(func(i int, selection *goquery.Selection) {
		vocab.Definition = append(vocab.Definition, Definition{
			Type:    selection.Find(".definition .anchor").Text(),
			Title:   selection.Find("h3.definition").Text(),
			Example: selection.Find(".defContent .example").Text(),
			Synonyms: DeepDefinition{
				ListWord:    findWords(selection, Synonyms),
				Description: findDescription(selection, Synonyms),
			},
			Antonyms: DeepDefinition{
				ListWord:    findWords(selection, Antonyms),
				Description: findDescription(selection, Antonyms),
			},
			Types: DeepDefinition{
				ListWord:    findWords(selection, TypeOf),
				Description: findDescription(selection, TypeOf),
			},
		})
	})

	return vocab
}

func findWords(selection *goquery.Selection, typeWord string) []string {
	var arr []string
	selection.Find("dl.instances").Map(func(i int, selection *goquery.Selection) string {
		if selection.Find("dt").Text() == typeWord {
			arr = append(arr, selection.Find(".word").Text())
		}
		return ""
	})
	return arr
}

func findDescription(selection *goquery.Selection, typeWord string) string {
	var des string
	selection.Find("dl.instances").Map(func(i int, selection *goquery.Selection) string {
		if selection.Find("dt").Text() == typeWord {
			des = selection.Find("dd div.definition").Text()
		}
		return des
	})
	return des
}

type Vocabulary struct {
	WordTitle       string
	DefinitionShort string
	DefinitionLong  string
	Definition      []Definition
}

type Definition struct {
	Title    string
	Type     string
	Example  string
	Synonyms DeepDefinition
	Antonyms DeepDefinition
	Types    DeepDefinition
}

type DeepDefinition struct {
	ListWord    []string
	Description string
}
