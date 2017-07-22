package main

import (
	"log"
	"net/http"
	"pigowlbot/api"
	"pigowlbot/private"
	"pigowlbot/sort"
	"strconv"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

func getPackages() string {
	response := api.GetPackages()

	var parts []string
	for _, pack := range response.Packs {
		parts = append(parts, pack.Pack.Name)
	}
	return strings.Join(parts, "\n")
}

func getPackagesName() map[int]string {
	packageIdNameMap := make(map[int]string)

	packsResponse := api.GetPackages()
	for _, pack := range packsResponse.Packs {
		packageIdNameMap[pack.Pack.ID] = pack.Pack.Name
	}
	return packageIdNameMap
}

func getDownloads(period int64) string {
	packageIdNameMap := getPackagesName()
	packsStatResponse := api.GetPackagesStatistics()

	downloadsMap := make(map[string]int)
	for _, packStat := range packsStatResponse.PacksStat {
		if packStat.Timestamp >= period {
			downloadsMap[packageIdNameMap[packStat.ID]]++
		}
	}

	var result []string
	for _, v := range sort.SortedKeys(downloadsMap) {
		result = append(result, v+", "+strconv.Itoa(downloadsMap[v]))
	}
	if len(result) > 0 {
		return strings.Join(result, "\n")
	}
	return "There were not any downloads :'("
}

func main() {
	bot, err := tgbotapi.NewBotAPI(private.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://www.pigowl.com:8443/"+bot.Token, "fullchain.pem"))
	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS(":8443", "fullchain.pem", "privkey.pem", nil)

	for update := range updates {
		log.Printf("%+v\n", update)
	}
}
