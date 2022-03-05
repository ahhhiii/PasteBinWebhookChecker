package main

import (
	"bufio"
	"fmt"
	"github.com/martinlindhe/notify"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var url = ""
var latestAnswer = ""

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	fmt.Println()
	fmt.Println("le funny webhook checker changer for pastebin :trol:")
	fmt.Println("~bruhitsalex")
	fmt.Println()
	reader := bufio.NewReader(os.Stdin)
	log.Info("Enter pastebin URL: ")

	text, _ := reader.ReadString('\n')
	url = strings.TrimSuffix(strings.TrimSuffix(text, "\n"), "\r")
	latestAnswer = getCurrentWebhook()

	fmt.Println()
	log.Info("Current content is: ")
	log.Info(latestAnswer)
	log.Info("Will now automatically refresh...")

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
		log.Warning("Answer changed:")
		log.Info(latestAnswer)
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

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Info("Deleted webhook [" + resp.Status + "]")
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
