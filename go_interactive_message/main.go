package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
)

var (
	signingSecret = "YOUR_SIGNING_SECRET_HERE!"
	tokenBotUser  = "YOUR_BOT_USER_OAUTH_ACCESS_TOKEN_HERE!"

	reqButtonPushedAction          = "buttonPushedAction"
	reqOrderModalSubmission        = "orderModalSubmission"
	reqConfirmationModalSubmission = "confirmationModalSubmission"
	reqUnknown                     = "unknown"

	burgers = map[string]string{
		"hamburger":     "Hamburger",
		"cheese_burger": "Cheese Burger",
		"blt_burger":    "BLT Burger",
		"big_burger":    "Big burger",
		"king_burger":   "King burger",
	}
)

type privateMeta struct {
	ChannelID string `json:"channel_id"`
	order
}

type order struct {
	Menu   string `json:"order_menu"`
	Steak  string `json:"order_steak"`
	Note   string `json:"order_note"`
	Amount string `json:"order_amount"`
}

func main() {
	// NOTE: In this example, we use 4 handlers. You should see what you want to know.
	// 1. Receive a message to call a bot and send an interactive message with button -> handleEventRequest()
	// 2. Receive a button pushed message and send an order modal -> handleButtonPushedRequest()
	// 3. Receive an order modal submission message and send a confirmation modal -> handleOrderModalSubmissionRequest()
	// 4. Receive a confirmation modal submission message and send a complession message -> handleConfirmationModalSubmissionRequest()

	lambda.Start(handleInteractiveRequest)
}

func handleInteractiveRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Verify the request.
	if err := verify(request, signingSecret); err != nil {
		log.Printf("[ERROR] Failed to verify request: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 200}, nil
	}

	// Parse the request
	payload, err := url.QueryUnescape(request.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to unescape: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 200}, nil
	}
	payload = strings.Replace(payload, "payload=", "", 1)

	var message slack.InteractionCallback
	if err := json.Unmarshal([]byte(payload), &message); err != nil {
		log.Printf("[ERROR] Failed to unmarshal json: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 200}, nil
	}

	// Identify the request type and dispatch message to appropreate handlers.
	switch identifyRequestType(message) {
	case reqButtonPushedAction:
		res, err := handleButtonPushedRequest(message)
		if err != nil {
			log.Printf("[ERROR] Failed to handle button pushed action: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 200}, nil
		}
		return res, nil
	case reqOrderModalSubmission:
		res, err := handleOrderSubmissionRequest(message)
		if err != nil {
			log.Printf("[ERROR] Failed to handle order submission: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 200}, nil
		}
		return res, nil
	case reqConfirmationModalSubmission:
		res, err := handleConfirmationModalSubmissionRequest(message)
		if err != nil {
			log.Printf("[ERROR] Failed to handle confirmation modal submission: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 200}, nil
		}
		return res, nil
	default:
		log.Printf("[ERROR] unknown request type: %v", message.Type)
		return events.APIGatewayProxyResponse{StatusCode: 200}, nil
	}
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

// identifyRequestType returns the request type of a slack message.
func identifyRequestType(message slack.InteractionCallback) string {

	// Check if the request is button pushed message.
	if message.Type == slack.InteractionTypeBlockActions && message.View.Hash == "" {
		return reqButtonPushedAction
	}

	// Check if the request is order modal submission.
	if message.Type == slack.InteractionTypeViewSubmission && strings.Contains(message.View.CallbackID, reqOrderModalSubmission) {
		return reqOrderModalSubmission
	}

	// Check if the request is confirmation modal submission.
	if message.Type == slack.InteractionTypeViewSubmission && strings.Contains(message.View.CallbackID, reqConfirmationModalSubmission) {
		return reqConfirmationModalSubmission
	}

	return reqUnknown
}
