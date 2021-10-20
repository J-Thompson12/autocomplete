package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/gorilla/mux"
)

type wordList struct {
	words map[string]int
}

type wordDistance struct {
	Word     string `json:"word"`
	Count    int    `json:"count"`
	Distance int    `json:"distance"`
}

func main() {
	words, err := scanWords("./shakespeare-complete.txt")
	if err != nil {
		panic(err)
	}
	// Thing

	router := mux.NewRouter()
	router.HandleFunc("/autocomplete", words.autoComplete).Methods("Get").Queries("term", "{term}")
	err = http.ListenAndServe(":5000", router)
	if err != nil {
		fmt.Println(err)
	}

}

func (words *wordList) autoComplete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	term := params["term"]
	if term == "" {
		http.Error(w, "Url Param 'term' is missing", 404)
		return
	}

	results := []wordDistance{}
	for word, count := range words.words {

		if term != word && len(term) >= len(word) {
			continue
		}

		if strings.HasPrefix(word, term) {
			distance := distance(term, word)
			newWord := wordDistance{Word: word, Count: count, Distance: distance}
			results = append(results, newWord)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Distance < results[j].Distance {
			return true
		}
		if results[i].Distance > results[j].Distance {
			return false
		}
		return results[i].Count > results[j].Count
	})

	if len(results) > 25 {
		writeJson(w, results[:25])
		return
	}
	writeJson(w, results)

}

func writeJson(w http.ResponseWriter, words []wordDistance) {
	js, err := json.Marshal(words)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func scanWords(path string) (wordList, error) {
	wordMap := make(map[string]int)

	file, err := os.Open(path)
	if err != nil {
		return wordList{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		for _, word := range words {
			processedString := reg.ReplaceAllString(word, "")
			if processedString != "" {
				word = strings.ReplaceAll(word, ".", "")
				word = strings.ReplaceAll(word, "!", "")
				word = strings.ReplaceAll(word, "?", "")
				word = strings.ReplaceAll(word, ",", "")
				wordMap[word]++
			}
		}
	}

	if len(wordMap) == 0 {
		panic("no words were found")
	}

	wordList := wordList{words: wordMap}
	return wordList, nil
}

// https://www.golangprograms.com/golang-program-for-implementation-of-levenshtein-distance.html#:~:text=Golang%20program%20for%20implementation%20of,to%20transform%20s%20into%20t.
func distance(source, destination string) int {
	vec1 := make([]int, len(destination)+1)
	vec2 := make([]int, len(destination)+1)

	w1 := []rune(source)
	w2 := []rune(destination)

	for i := 0; i < len(vec1); i++ {
		vec1[i] = i
	}

	for i := 0; i < len(w1); i++ {
		vec2[0] = i + 1

		for j := 0; j < len(w2); j++ {
			cost := 1
			if w1[i] == w2[j] {
				cost = 0
			}
			min := min(vec2[j]+1,
				vec1[j+1]+1,
				vec1[j]+cost)
			vec2[j+1] = min
		}

		for j := 0; j < len(vec1); j++ {
			vec1[j] = vec2[j]
		}
	}

	return vec2[len(w2)]
}

func min(value0 int, values ...int) int {
	min := value0
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}
