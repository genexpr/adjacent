package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	group = flag.String("group", "g", "language group - choose from germanic (g), slavic (s), romance (r)")
	word  = flag.String("word", "language", "word to translate")
)

var (
	slavic   = []string{"ru", "be", "bg", "bs", "mk", "pl", "sr", "sk", "sl", "cs", "hr", "uk"}
	germanic = []string{"af", "nl", "da", "is", "de", "no", "sv"}
	romance  = []string{"it", "pt", "ro", "fr", "es", "ca"}
)

type response struct {
	Code int
	Lang string
	Text []string
}

func main() {
	flag.Usage = usage
	flag.Parse()

	var token = os.Getenv("YANDEX_TRANSLATE_TOKEN")
	if token == "" {
		log.Fatalf("your token for the API does not exist, set it as a value of the %s environment variable\n",
			"YANDEX_TRANSLATE_TOKEN")
	}

	languages, err := getLanguages(*group)
	if err != nil {
		usage()
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	for i := 0; i < len(languages); i++ {
		wg.Add(1)

		go func(lang string) {
			translation, err := makeRequest(*word, lang, token)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println(translation)
			wg.Done()

		}(languages[i])
	}

	wg.Wait()
}

func usage() {
	fmt.Println("Usage:")
	flag.PrintDefaults()
}

func getLanguages(group string) ([]string, error) {
	switch group {
	case "g", "germanic":
		return germanic, nil
	case "s", "slavic":
		return slavic, nil
	case "r", "romance":
		return romance, nil
	default:
		return nil, errors.New("language group not supported")
	}
}

// makeRequest makes a request to the Yandex.Translate API to get the translation of a word
// from English to a given language.
func makeRequest(word string, language string, token string) (string, error) {
	url := fmt.Sprintf("https://translate.yandex.net/api/v1.5/tr.json/translate?key=%s&text=%s&lang=en-%s",
		token, word, language)
	resp, err := http.Get(url)

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalln("could not connect to the API, make sure your token is valid")
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var r response
	json.Unmarshal(bytes, &r)
	output := fmt.Sprintf("%s\t%s", r.Lang, r.Text[0])
	return output, nil
}
