package ripe_atlas

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/innogames/slack-bot/v2/bot"
	"github.com/innogames/slack-bot/v2/bot/matcher"
	"github.com/innogames/slack-bot/v2/bot/msg"
	"github.com/innogames/slack-bot/v2/client"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"io"
	"net/http"
	"net/netip"
	"time"
)

type tracerouteCommand struct {
	bot.BaseCommand
	cfg Config
}

func (c *tracerouteCommand) GetMatcher() matcher.Matcher {
	return matcher.NewGroupMatcher(
		matcher.NewRegexpMatcher(`traceroute (?P<TGT>.*)`, c.traceroute),
	)
}

func (c *tracerouteCommand) traceroute(match matcher.Result, message msg.Message) {
	destination := match.GetString("TGT")

	c.AddReaction(":stopwatch:", message)
	defer c.RemoveReaction(":stopwatch:", message)

	var af int
	address, err := netip.ParseAddr(destination)
	if err != nil {
		af = 6
	} else {
		if address.Is4() {
			af = 4
		} else {
			af = 6
		}
	}

	jsonData, _ := json.Marshal(MeasurementRequest{
		Definitions: []MeasurementDefinition{
			{
				Af:             af,
				Target:         destination,
				Description:    fmt.Sprintf("Slackbot measurement to %s", destination),
				Type:           "traceroute",
				Protocol:       "ICMP",
				Packets:        3,
				ResolveOnProbe: false,
				Paris:          0,
				IsPublic:       true,
			},
		},
		Probes: []Probes{
			{
				Type:      "area",
				Value:     "WW",
				Requested: 1,
			},
		},
		IsOneOff: true,
	})

	url := fmt.Sprintf("%s/measurements", c.cfg.APIURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.ReplyError(message, errors.Wrap(err, "Request creation returned an err"))
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Key "+c.cfg.APIKey)

	response, err := client.GetHTTPClient().Do(req)
	if err != nil {
		c.ReplyError(message, errors.Wrap(err, "HTTP Client Error"))
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 300 {
		c.ReplyError(message, errors.New(fmt.Sprintf("Api call returned an err: %d", response.StatusCode)))
		return
	}

	body, _ := io.ReadAll(response.Body)

	var measurementResult MeasurementResult
	err = json.Unmarshal(body, &measurementResult)

	if err != nil {
		c.ReplyError(message, errors.Wrap(err, "Error unmarshalling MeasurementResult"))
		return
	}

	//replyRef := c.SendMessage(
	c.SendMessage(
		message,
		fmt.Sprintf("Measurement created: https://atlas.ripe.net/measurements/%d\n", measurementResult.Measurements[0]),
		slack.MsgOptionTS(message.GetTimestamp()),
	)

	messageUpdates := make(chan string, 2)

	go func() {
		defer close(messageUpdates)

		subscribeUrl := fmt.Sprintf("https://atlas-stream.ripe.net/stream/?streamType=result&msm=%d", measurementResult.Measurements[0])

		client := http.Client{Timeout: 240 * time.Second}
		response, err = client.Get(subscribeUrl)
		defer response.Body.Close()
		fileScanner := bufio.NewScanner(response.Body)
		fileScanner.Split(bufio.ScanLines)
	scanner:
		for fileScanner.Scan() {
			line := fileScanner.Text()

			var streamResponse StreamingResponse
			err = json.Unmarshal([]byte(line), &streamResponse)
			if err != nil {
				log.Errorf("Error unmarshaling streamResponse: %s", err)
			}

			switch streamResponse.Type {
			case "atlas_subscribed":
				log.Debugf("Succesfully subscribed to measurement")
			case "atlas_result":
				srp := streamResponse.Payload
				messageUpdates <- fmt.Sprintf("%s", srp)
				response.Body.Close()
				break scanner
			}
		}
	}()

	for delta := range messageUpdates {
		c.SendMessage(
			message,
			delta,
			slack.MsgOptionTS(message.GetTimestamp()),
		)
	}
}

func (c *tracerouteCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "traceroute <destination>",
			Description: "Sends a traceroute to the given destination",
			Category:    category,
			Examples: []string{
				"traceroute 8.8.8.8",
			},
		},
	}
}
