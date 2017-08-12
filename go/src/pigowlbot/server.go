package main

import (
	"log"
	"net/http"
	"pigowlbot/api"
	"pigowlbot/database"
	"pigowlbot/private"
	"pigowlbot/sort"
	"strconv"
	"strings"
	"time"

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

func formatDownloadsMessage(sortedMap *sort.SortedMap) string {
	var result []string
	for _, v := range sortedMap.Keys {
		result = append(result, v+", "+strconv.Itoa(sortedMap.Original[v]))
	}
	if len(result) > 0 {
		return strings.Join(result, "\n")
	}
	return "There were not any downloads :'("
}

func formatDiffDownloadsMessage(updatedMap *sort.SortedMap, dailyMap *sort.SortedMap) string {
	var result []string
	for _, v := range updatedMap.Keys {
		result = append(result, v+", "+strconv.Itoa(updatedMap.Original[v])+" ("+strconv.Itoa(dailyMap.Original[v])+")")
	}
	if len(result) > 0 {
		return strings.Join(result, "\n")
	}
	return "There were not any downloads :'("
}

func getDownloads(period int64) *sort.SortedMap {
	packageIDNameMap := getPackagesName()
	packsStatResponse := api.GetPackagesStatistics()

	downloadsMap := make(map[string]int)
	for _, packStat := range packsStatResponse.PacksStat {
		if packStat.Timestamp >= period {
			downloadsMap[packageIDNameMap[packStat.ID]]++
		}
	}

	return sort.SortedKeys(downloadsMap)
}

func timer(channel chan int) {
	for {
		time.Sleep(10 * time.Minute)
		//time.Sleep(time.Second * 5)
		channel <- 0
	}
}

func subscribe(bot *tgbotapi.BotAPI, chatID int64) {
	lastDownloadsTimestamp := time.Now().Unix()
	timerChannel := make(chan int)
	go timer(timerChannel)
	for _ = range timerChannel {
		downloads := getDownloads(lastDownloadsTimestamp)
		lastDownloadsTimestamp = time.Now().Unix()

		if len(downloads.Keys) > 0 {
			msg := tgbotapi.NewMessage(chatID, formatDiffDownloadsMessage(downloads, getDownloads(time.Now().Truncate(24*time.Hour).Unix())))
			bot.Send(msg)
		}
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(private.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://pigowl.com:88/"))
	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/")
	go http.ListenAndServeTLS(":88", "fullchain.pem", "privkey.pem", nil)

	//subscribers := make(map[int64]bool)
	dbc := database.Init()
	subscribers := dbc.GetSubscribers()
	for _, subscriber := range subscribers {
		go subscribe(bot, subscriber)
	}

	for update := range updates {
		command := update.Message.Command()
		switch command {
		case "getpackages":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, getPackages())
			bot.Send(msg)
		case "getweeklydownloads":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, formatDownloadsMessage(getDownloads(time.Now().Add(-7*24*time.Hour).Truncate(24*time.Hour).Unix())))
			bot.Send(msg)
		case "getdailydownloads":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, formatDownloadsMessage(getDownloads(time.Now().Truncate(24*time.Hour).Unix())))
			bot.Send(msg)
		case "getalldownloads":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, formatDownloadsMessage(getDownloads(0)))
			bot.Send(msg)
		case "subscribe":
			//_, exist := subscribers[update.Message.Chat.ID]
			exist := dbc.CheckSubscriber(update.Message.Chat.ID)
			if !exist {
				go subscribe(bot, update.Message.Chat.ID)
				//subscribers[update.Message.Chat.ID] = true
				dbc.AddSubscriber(update.Message.Chat.ID)
			}
		}
	}
}
