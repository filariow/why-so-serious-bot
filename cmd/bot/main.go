package main

import (
	"fmt"
	"log"
	"os"

	"github.com/FrancescoIlario/why-so-serious-bot/internal/bot"
	"github.com/FrancescoIlario/why-so-serious-bot/internal/conf"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/envext"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssface"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssformrecognizer"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssmoderator"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wsssentiment"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wsstranslator"
	"github.com/FrancescoIlario/why-so-serious-bot/pkg/wssvision"
	"github.com/joho/godotenv"
)

func main() {
	loadEnvs() // load the env vars from .env file

	conf, err := getConfigurations()
	if err != nil {
		log.Fatalln(err)
	}

	fbot, err := bot.New(*conf)
	if err != nil {
		log.Printf("can not instantiate bot: %v", err)
	}

	// start bot
	fbot.Start()

	// wait undefinetly
	shutdown := make(chan struct{})
	<-shutdown
}

func getConfigurations() (*bot.Configuration, error) {
	c := bot.Configuration{}

	{ // Telegram Configuration
		pollerInterval, err := conf.GetPollerInterval()
		if err != nil {
			return nil, fmt.Errorf("error retrieving poller interval: %v", err)
		}
		c.PollerInterval = *pollerInterval

		token, err := conf.GetToken()
		if err != nil {
			return nil, fmt.Errorf("error retrieving Telegram token: %v", err)
		}
		c.Token = *token
	}
	{ // Vision: Face
		faceConf, err := wssface.BuildConfigurationFromEnvs()
		if err != nil {
			log.Printf("error retrieving face service configuration: %v", err)
		}
		c.FaceConf = faceConf
	}
	{ // Vision API
		visionConf, err := wssvision.BuildConfigurationFromEnvs()
		if err != nil {
			log.Printf("error retrieving vision service configuration: %v", err)
		}
		c.VisionConf = visionConf
	}
	{ // Vision: Form Recognizer
		formRecognizerConf, err := wssformrecognizer.BuildConfigurationFromEnvs()
		if err != nil {
			log.Printf("error retrieving form recognizer service configuration: %v", err)
		}
		c.FormRecognizerConf = formRecognizerConf
	}
	{ // Language Text Analytics
		textAnalyticsConf, err := wsssentiment.BuildConfigurationFromEnvs()
		if err != nil {
			log.Printf("error retrieving text analitycs service configuration: %v", err)
		}
		c.TextAnalyticsConf = textAnalyticsConf
	}
	{ // Language: Translator
		translatorConf, err := wsstranslator.BuildConfigurationFromEnvs()
		if err != nil {
			log.Printf("error retrieving translator service configuration: %v", err)
		}
		c.TranslatorConf = translatorConf
	}
	{ // Decision: Moderator
		moderatorConf, err := wssmoderator.BuildConfigurationFromEnvs()
		if err != nil {
			log.Printf("error retrieving moderator service configuration: %v", err)
		}
		c.ModeratorConf = moderatorConf
	}

	return &c, nil
}

// AppEnvKey Environment variable key where is stored the environment to use for the app
// execution. default is `development`
const AppEnvKey = "WSSBOT_ENV"

func loadEnvs() {
	env := envext.GetEnvOrDefault(AppEnvKey, "development")
	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}

	envfile := ".env." + env
	if err := godotenv.Load(envfile); err != nil {
		cwd, _ := os.Getwd()
		log.Printf("error loading file %v (%s): %v", envfile, cwd, err)
	}

	if err := godotenv.Load(); err != nil { // The Original .env
		cwd, _ := os.Getwd()
		log.Printf("error loading .env (%v) file: %v", cwd, err)
	}
}
