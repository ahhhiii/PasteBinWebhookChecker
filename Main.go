package main

import (
	"bufio"
	"fmt"
	"github.com/martinlindhe/notify"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var url = ""
var latestAnswer = ""

func main() {
	fmt.Println()
	fmt.Println("le funny webhook checker changer for pastebin :trol:")
	fmt.Println("~bruhitsalex")
	fmt.Println()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter pastebin URL: ")

	text, _ := reader.ReadString('\n')
	url = strings.TrimSuffix(strings.TrimSuffix(text, "\n"), "\r")
	latestAnswer = getCurrentWebhook()

	fmt.Println()
	fmt.Println("Current content is '" + latestAnswer + "'")
	fmt.Println("Will now automatically refresh...")

	for range time.Tick(time.Minute * 2) {
		go func() {
			checkAnswer()
		}()
	}
}

func checkAnswer() {
	result := getCurrentWebhook()
	if result != latestAnswer {
		latestAnswer = result
		notify.Notify("Pastebin Checker", "Content Change", latestAnswer, "")
		fmt.Println()
		fmt.Println("Answer changed: " + latestAnswer)
		if strings.Contains(latestAnswer, "discord.com/api/webhooks") {
			deleteAnswer()
		}
	}
}

func deleteAnswer() {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", latestAnswer, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Deleted webhook: " + string(respBody))
}

func getCurrentWebhook() string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	return string(body)
}
