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
	group = flag.String("group", "", "language group - choose from germanic (g), slavic (s), romance (r)")
	word  = flag.String("word", "language", "word to translate")
)

type Response struct {
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

	languageList, err := GetLanguageList(*group)
	if err != nil {
		log.Fatalln(err)
	}

	wg := sync.WaitGroup{}
	for i := 0; i < len(languageList); i++ {
		wg.Add(1)
		l := languageList[i]
		go func() {
			translation, err := Request(*word, l, token)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println(translation)
			wg.Done()
		}()
	}

	wg.Wait()
}

func usage() {
	fmt.Println("Usage:")
	flag.PrintDefaults()
}

func GetLanguageList(group string) ([]string, error) {
	if group == "" {
		usage()
		os.Exit(1)
	}

	var slavicList = []string{"ru", "be", "bg", "bs", "mk", "pl", "sr", "sk", "sl", "cs", "hr", "uk"}
	var germanicList = []string{"af", "nl", "da", "is", "de", "no", "sv"}
	var romanceList = []string{"it", "pt", "ro", "fr", "es", "ca"}

	if group == "g" || group == "germanic" {
		return germanicList, nil
	} else if group == "s" || group == "slavic" {
		return slavicList, nil
	} else if group == "r" || group == "romance" {
		return romanceList, nil
	} else {
		return nil, errors.New("language group not supported")
	}
}

// Request makes a request to the Yandex.Translate API to get the translation of a word
// from English to a given language.
func Request(word string, language string, token string) (string, error) {
	url := fmt.Sprintf("https://translate.yandex.net/api/v1.5/tr.json/translate?key=%s&text=%s&lang=en-%s",
		token, word, language)
	response, err := http.Get(url)

	if err != nil || response.StatusCode != http.StatusOK {
		log.Fatalln("could not connect to the API, make sure your token is valid")
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var res Response
	json.Unmarshal(bytes, &res)
	output := fmt.Sprintf("%s\t%s", res.Lang, res.Text[0])
	return output, nil
}
