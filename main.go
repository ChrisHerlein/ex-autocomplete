package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var tree = &node{}

func load(seed string) {
	wordsSeed, e := ioutil.ReadFile(seed)
	if e != nil {
		panic("Error reading seed: " + e.Error())
	}
	wordsString := string(wordsSeed)
	wordsString = cleanChars(wordsString)
	words := strings.Split(wordsString, ",")
	for i := 0; i < len(words); i++ {
		tree.Add(words[i])
	}
}

func main() {
	// File should be a comma-separated words list
	seed := flag.String("seed", "seed.txt", "file to look for words to work with")
	search := flag.String("search", "", "prefix to look for")
	flag.Parse()

	load(*seed)

	if len(*search) == 0 {
		fmt.Println("Nothing to look for!")
		return
	}
	found := tree.Search(*search)
	fmt.Printf("Suggestions: %+v\n", found)
}

type node struct {
	root     string
	isWord   bool
	code     rune
	children []*node
}

func (n *node) Add(word string) {
	defer func() {
		e := recover()
		if e != nil {
			fmt.Println("Can not understand:", word)
		}
	}()
	if word == "" {
		return
	}
	rn := []rune(word)
	if len(n.children) == 0 {
		n.children = make([]*node, 27)
	}
	if n.children[rn[0]-97] == nil {
		n.children[rn[0]-97] = &node{
			root:     string(word[0]),
			code:     rn[0] - 97,
			children: make([]*node, 27),
		}
	}
	n.children[rn[0]-97].add(word, rn, 1)
}

func (n *node) add(word string, rn []rune, index int) {
	if len(word) == index {
		n.isWord = true
		return
	}
	if n.children[rn[index]-97] == nil {
		n.children[rn[index]-97] = &node{
			root:     word[0 : index+1],
			code:     rn[index] - 97,
			children: make([]*node, 27),
		}
	}
	n.children[rn[index]-97].add(word, rn, index+1)
}

func (n *node) Search(prefix string) []string {
	prefix = cleanChars(prefix)
	rn := []rune(prefix)
	return n.children[rn[0]-97].search(prefix, rn, 1)
}

func (n *node) search(prefix string, rn []rune, index int) []string {
	if index == len(prefix) {
		ans := make(chan []string, 1)
		go n.getWords(ans)
		found := <-ans
		return found
	}
	if n.children[rn[index]-97] == nil {
		return []string{}
	}
	return n.children[rn[index]-97].search(prefix, rn, index+1)
}

func (n *node) getWords(upper chan []string) {
	found := make([]string, 0)
	if n.isWord {
		found = append(found, n.root)
	}
	var sent = 0
	ans := make(chan []string, 27)
	for i := 0; i < 27; i++ {
		if n.children[i] != nil {
			sent++
			go n.children[i].getWords(ans)
		}
	}
	for i := 0; i < sent; i++ {
		found = append(found, (<-ans)...)
	}
	upper <- found
}

func cleanChars(word string) string {
	word = strings.ToLower(word)
	word = strings.ReplaceAll(word, "ñ", "{") // to put ñ at end of children
	word = strings.ReplaceAll(word, "á", "a")
	word = strings.ReplaceAll(word, "é", "e")
	word = strings.ReplaceAll(word, "è", "e")
	word = strings.ReplaceAll(word, "í", "i")
	word = strings.ReplaceAll(word, "ì", "i")
	word = strings.ReplaceAll(word, "ó", "o")
	word = strings.ReplaceAll(word, "ú", "u")
	word = strings.ReplaceAll(word, ")", "")
	word = strings.ReplaceAll(word, "(", "")
	word = strings.ReplaceAll(word, ":", "")
	word = strings.ReplaceAll(word, ".", "")
	word = strings.ReplaceAll(word, "-", "")
	word = strings.ReplaceAll(word, "!", "")
	word = strings.ReplaceAll(word, "¡", "")
	word = strings.ReplaceAll(word, "?", "")
	word = strings.ReplaceAll(word, "\"", "")
	return word
}
