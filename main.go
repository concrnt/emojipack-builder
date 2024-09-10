package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type EmojiPack struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Credit      string  `json:"credits"`
	IconURL     string  `json:"iconURL"`
	Emojis      []Emoji `json:"emojis"`
}

type Emoji struct {
	Shortcode string `json:"shortcode"`
	ImageURL  string `json:"imageURL"`
	aliases   []string
}

func main() {

	BaseURL := os.Getenv("BASE_URL")
	if BaseURL == "" {
		fmt.Println("Warn: BASE_URL is not set")
	}

	var targets = []string{}
	entries, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			dirname := entry.Name()
			if dirname[0] == '.' {
				continue
			}
			targets = append(targets, dirname)
		}
	}

	for _, target := range targets {
		fmt.Println("Processing: ", target)
		// Read Metadata File
		metadataFile, err := os.ReadFile(target + "/metadata.json")
		if err != nil {
			panic(err)
		}

		var metadata EmojiPack
		err = json.Unmarshal(metadataFile, &metadata)
		if err != nil {
			panic(err)
		}

		// Get the list of files in the target directory
		files, err := os.ReadDir(target)
		if err != nil {
			panic(err)
		}

		emojis := make([]Emoji, 0)
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			split := strings.Split(file.Name(), ".")
			basename := split[0]

			extension := split[len(split)-1]
			if extension == "json" {
				continue
			}

			names := strings.Split(basename, "-")
			aliases := make([]string, 0)
			if len(names) > 1 {
				aliases = names[1:]
			}

			emojis = append(
				emojis,
				Emoji{
					Shortcode: basename,
					ImageURL:  BaseURL + "/" + target + "/" + file.Name(),
					aliases:   aliases,
				},
			)
		}

		EmojiPack := EmojiPack{
			Name:        metadata.Name,
			Description: metadata.Description,
			Credit:      metadata.Credit,
			IconURL:     BaseURL + "/" + target + "/icon.png",
			Emojis:      emojis,
		}

		// Convert the struct to JSON
		emojiJSON, err := json.Marshal(EmojiPack)
		if err != nil {
			panic(err)
		}

		// Write the JSON to a file
		f, err := os.Create(target + "/emojis.json")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		f.Write(emojiJSON)
		f.Sync()
	}
}
