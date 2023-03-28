package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	group = flag.String("group", "", "Language group - choose from Germanic (g), Slavic (s), Romance (r)")
	text  = flag.String("text", "", "English text to translate")
)

func main() {
	flag.Usage = usageAndExit
	flag.Parse()

	var token = os.Getenv("DEEPL_TRANSLATE_TOKEN")
	if token == "" {
		_, _ = fmt.Fprintln(os.Stderr, "API token missing, set it as the value of the DEEPL_TRANSLATE_TOKEN environment variable.")
		os.Exit(1)
	}

	languages, err := getLanguagesFromGroup(*group)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		usageAndExit()
	}

	if *text == "" {
		usageAndExit()
	}

	var wg sync.WaitGroup
	for _, language := range languages {
		wg.Add(1)

		go func(lang string) {
			defer wg.Done()
			translation, err := translate(*text, lang, token)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return
			}
			fmt.Printf("%s\t%s\n", lang, translation)

		}(language)
	}

	wg.Wait()
}

func usageAndExit() {
	fmt.Println("Usage:")
	flag.PrintDefaults()
	os.Exit(1)
}

func getLanguagesFromGroup(group string) ([]string, error) {
	if group == "" {
		return nil, fmt.Errorf("language group is not provided")
	}
	var (
		slavic   = []string{"RU", "BG", "PL", "CS", "LT", "LV", "SL", "SK", "UK"}
		romance  = []string{"IT", "PT-PT", "RO", "FR", "ES"}
		germanic = []string{"NL", "DA", "DE", "SV"}
	)
	switch group {
	case "g", "germanic":
		return germanic, nil
	case "s", "slavic":
		return slavic, nil
	case "r", "romance":
		return romance, nil
	default:
		return nil, fmt.Errorf("language group %s is not supported", group)
	}
}

// translate sends a request to the DeepL API to get the translation of a piece of text
// from English into a given language.
func translate(text, language, token string) (string, error) {
	data := url.Values{}
	data.Set("text", text)
	data.Set("source_lang", "EN")
	data.Set("target_lang", language)
	encodedData := data.Encode()

	const baseURL = "https://api-free.deepl.com/v2/translate"
	req, err := http.NewRequest(http.MethodPost, baseURL, strings.NewReader(encodedData))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("Authorization", "DeepL-Auth-Key "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to get a response for " + language)
	}
	defer resp.Body.Close()

	var r APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}

	if len(r.Translations) == 0 || r.Translations[0].Text == "" {
		return "", errors.New("no translation available")
	}

	return r.Translations[0].Text, nil
}

type APIResponse struct {
	Translations []struct {
		Text string `json:"text"`
	} `json:"translations"`
}
