package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack/slackevents"
	"github.com/slack-go/slack"
)

var (
	signingSecret = "YOUR_SIGNING_SECRET_HERE!"
	tokenBotUser  = "YOUR_BOT_USER_OAUTH_ACCESS_TOKEN_HERE!"
)

func main() {
	// NOTE: In this example, we use 4 handlers. You should see what you want to know.
	// 1. Receive a message to call a bot and send an interactive message with button -> handleEventRequest()
	// 2. Receive a button pushed message and send an order modal -> handleButtonPushedRequest()
	// 3. Receive an order modal submission message and send a confirmation modal -> handleOrderModalSubmissionRequest()
	// 4. Receive a confirmation modal submission message and send a complession message -> handleConfirmationModalSubmissionRequest()

	lambda.Start(handleEventRequest)
}

func handleEventRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Verify the request.
	if err := verify(request, signingSecret); err != nil {
		log.Printf("[ERROR] Failed to verify request: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 200}, nil
	}

	// Parse event.
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(request.Body), slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Printf("[ERROR] Failed to parse request body: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 200}, nil
	}

	// Check if the request type is URL Verification. This logic is only called from slack developer's console when you set up your app.
	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		if err := json.Unmarshal([]byte(request.Body), &r); err != nil {
			log.Printf("[ERROR] Failed to unmarshal json: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 200}, nil
		}
		return events.APIGatewayProxyResponse{Body: r.Challenge, StatusCode: 200}, nil
	}

	// Verify the request type.
	if eventsAPIEvent.Type != slackevents.CallbackEvent {
		log.Printf("[ERROR] Unexpected event type: expect = CallbackEvent , actual = %v", eventsAPIEvent.Type)
		return events.APIGatewayProxyResponse{StatusCode: 200}, nil
	}

	// Verify the event type.
	switch ev := eventsAPIEvent.InnerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:

		// Create a shop list.
		list := createShopListBySDK()

		// Send a shop list to slack channel.
		api := slack.New(tokenBotUser)
		if _, _, err := api.PostMessage(ev.Channel, list); err != nil {
			log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 200}, nil
		}

	default:
		return events.APIGatewayProxyResponse{StatusCode: 200}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

// createShopListBySDK returns a message option which contains shop infomation.
func createShopListBySDK() slack.MsgOption {
	// Top text
	descText := slack.NewTextBlockObject("mrkdwn", "What do you want to have?", false, false)
	descTextSection := slack.NewSectionBlock(descText, nil, nil)

	// Divider
	dividerBlock := slack.NewDividerBlock()

	// Shops
	// - Hamburger
	hamburgerButtonText := slack.NewTextBlockObject("plain_text", "Order", true, false)
	hamburgerButtonElement := slack.NewButtonBlockElement("actionIDHamburger", "hamburger", hamburgerButtonText)
	hamburgerAccessory := slack.NewAccessory(hamburgerButtonElement)
	hamburgerSectionText := slack.NewTextBlockObject("mrkdwn", ":hamburger: *Hungryman Hamburgers*\nOnly for the hungriest of the hungry.", false, false)
	hamburgerSection := slack.NewSectionBlock(hamburgerSectionText, nil, hamburgerAccessory)

	// - Sushi
	sushiButtonText := slack.NewTextBlockObject("plain_text", "Order", true, false)
	sushiButtonElement := slack.NewButtonBlockElement("actionIDSushi", "sushi", sushiButtonText)
	sushiAccessory := slack.NewAccessory(sushiButtonElement)
	sushiSectionText := slack.NewTextBlockObject("mrkdwn", ":sushi: *Ace Wasabi Rock-n-Roll Sushi Bar*\nFresh raw wish and wasabi.", false, false)
	sushiSection := slack.NewSectionBlock(sushiSectionText, nil, sushiAccessory)

	// - Ramen
	ramenButtonText := slack.NewTextBlockObject("plain_text", "Order", true, false)
	ramenButtonElement := slack.NewButtonBlockElement("actionIDRamen", "ramen", ramenButtonText)
	ramenAccessory := slack.NewAccessory(ramenButtonElement)
	ramenSectionText := slack.NewTextBlockObject("mrkdwn", ":ramen: *Sazanami Ramen*\nWhy don't you try Japanese soul food?", false, false)
	ramenSection := slack.NewSectionBlock(ramenSectionText, nil, ramenAccessory)

	// Blocks
	blocks := slack.MsgOptionBlocks(descTextSection, dividerBlock, hamburgerSection, sushiSection, ramenSection)

	return blocks
}

// verify returns the result of slack signing secret verification.
func verify(request events.APIGatewayProxyRequest, sc string) error {
	body := request.Body
	header := http.Header{}
	for k, v := range request.Headers {
		header.Set(k, v)
	}

	sv, err := slack.NewSecretsVerifier(header, sc)
	if err != nil {
		return err
	}

	sv.Write([]byte(body))
	return sv.Ensure()
}
