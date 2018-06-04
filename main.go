package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const challengeTime = 30

func main() {
	inputCh := input(os.Stdin)
	timeoutChan := time.NewTimer(challengeTime * time.Second)
	challenges := 0
	corrects := 0

	fmt.Printf("%d秒の間に何単語タイプできるか？\n", challengeTime)

FINISH:
	for {
		word, err := getWord()
		if err != nil {
			panic(err)
		}
		challenges++

		fmt.Println(word)
		fmt.Print("> ")
		select {
		case input := <-inputCh:
			if input == word {
				println("○ 正解！")
				corrects++
			} else {
				println("☓ 残念…")
			}
		case <-timeoutChan.C:
			break FINISH
		}
	}

	fmt.Printf("\nお疲れ様でした\n")
	fmt.Printf("出題数:%d 正答数:%d\n", challenges, corrects)
}

func getWord() (string, error) {
	resp, errHTTP := http.Get("http://api.chew.pro/trbmb")
	if errHTTP != nil {
		return "", errHTTP
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	jsonBytes := ([]byte)(byteArray)
	var sentences []string
	errJSON := json.Unmarshal(jsonBytes, &sentences)
	if errJSON != nil {
		return "", errJSON
	}

	words := strings.Split(sentences[0], " ")

	return words[2], nil
}

func input(reader io.Reader) <-chan string {
	ch1 := make(chan string)

	go func() {
		answer := bufio.NewScanner(reader)
		for answer.Scan() {
			ch1 <- answer.Text()
		}
		close(ch1)
	}()

	return ch1
}
