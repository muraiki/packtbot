package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/belak/irc"
	"gopkg.in/yaml.v2"
	"crypto/tls"
)

type config struct {
	Botname, Channel, Server, Secureserver string
}

func getTitle() (string, error) {
	resp, err := http.Get("https://www.packtpub.com/packt/offers/free-learning")
	if err != nil {
		return "", errors.New("Error getting url")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Error reading response body")
	}

	re, _ := regexp.Compile(`<div class="dotd-title">\s*<h2>\s+\b(.+)\b\s+`)

	return re.FindStringSubmatch(string(body))[1], nil
}

func getConfig(c *config) error {
	f, err := ioutil.ReadFile("packtbot.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		return err
	}
	return nil
}

type ircBundle struct {
	Client *irc.Client
	Event *irc.Event
}

func (ib ircBundle) Message (s string) {
	ib.Client.Reply(ib.Event, s)
}

func handleGetPackt(ch chan ircBundle) {
	for {
		ib := <-ch
		title, err := getTitle()
		if err != nil {
			log.Println(err.Error())
			return
		}
		ib.Message(fmt.Sprintf("The current free book is: '%s'. Download it at https://www.packtpub.com/packt/offers/free-learning", title))
	}
}

func main() {
	var conf config

	if err := getConfig(&conf); err != nil {
		log.Println(err.Error())
		return
	}

	handler := irc.NewBasicMux()
	client := irc.NewClient(
		irc.HandlerFunc(handler.HandleEvent), conf.Botname, conf.Botname, conf.Botname, "")

	getPacktChan := make(chan ircBundle)
	go handleGetPackt(getPacktChan)

	// event 001 is received once fully connected
	handler.Event("001", func(c *irc.Client, e *irc.Event) {
		c.Write(fmt.Sprintf("JOIN %s", conf.Channel))
		log.Printf(fmt.Sprintf("Joined channel %s", conf.Channel))
	})

	handler.Event("PRIVMSG", func(c *irc.Client, e *irc.Event) {
		if e.Trailing() == "!packt" {
			getPacktChan <- ircBundle{c, e}
		}
	})

	if conf.Server != "" {
		log.Printf(fmt.Sprintf("Connecting to server (unencrypted): %s", conf.Server))
		err := client.Dial(conf.Server)
		if err != nil {
			log.Fatalln(err)
		}
	} else if conf.Secureserver != "" {
		log.Printf(fmt.Sprintf("Connecting to server using TLS: %s", conf.Secureserver))
		tlc := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         conf.Secureserver,
		}
		err := client.DialTLS(conf.Secureserver, tlc)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Println("No server specified")
	}
}
