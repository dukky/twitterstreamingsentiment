package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/cdipaolo/sentiment"
)

type configuration struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

var search = flag.String("terms", "brexit", "The term(s) to search for.")

func main() {
	flag.Parse()
	config, err := getConfigFromFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Set up twitter api
	api := getAPI(config)
	defer api.Close()
	v := url.Values{}
	v.Set("track", *search)
	stream := api.PublicStreamFilter(v)

	// Set up sentiment model
	model, err := sentiment.Restore()
	if err != nil {
		log.Fatal(err)
	}

	totalSentiment := []int{}

	// Start processing tweets
	for data := range stream.C {

		switch tweet := data.(type) {
		case anaconda.Tweet:
			if tweet.RetweetedStatus == nil {
				sentimentAnalysis := model.SentimentAnalysis(tweet.Text, sentiment.English)
				//fmt.Printf("Sentiment: %d, Tweet: %s\n", sentimentAnalysis.Score, tweet.Text)
				totalSentiment = append(totalSentiment, int(sentimentAnalysis.Score))
				length := len(totalSentiment)
				sum := 0
				for _, x := range totalSentiment {
					sum += x
				}
				fmt.Printf("Total tweets: %d, Average sentiment: %f               \r", length, float32(sum)/float32(length))
			}
		case anaconda.LimitNotice:
			fmt.Printf("Skipped %d tweets.\n", tweet.Track)
		default:
			log.Println(tweet)
		}
	}
}

func getConfigFromFile(path string) (configuration, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return configuration{}, err
	}
	decoder := json.NewDecoder(file)
	config := configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}
	return config, nil

}

func getAPI(c configuration) *anaconda.TwitterApi {
	anaconda.SetConsumerKey(c.ConsumerKey)
	anaconda.SetConsumerSecret(c.ConsumerSecret)
	var api = anaconda.NewTwitterApi(c.AccessToken, c.AccessTokenSecret)
	return api
}
