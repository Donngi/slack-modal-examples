package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack"
)

func handleButtonPushedRequest(message slack.InteractionCallback) (events.APIGatewayProxyResponse, error) {
	// Get selected value
	shop := message.ActionCallback.BlockActions[0].Value

	switch shop {
	case "hamburger":
		// Create an order modal.
		// - apperance
		modal := createOrderModalBySDK()

		// You can also create a modal apperance by using JSON.
		// modal, err := createOrderModalByJSON()
		// if err != nil {
		// 	return events.APIGatewayProxyResponse{StatusCode: 200}, fmt.Errorf("failed to create modal: %w", err)
		// }

		// - metadata : CallbackID
		modal.CallbackID = reqOrderModalSubmission

		// - metadata : ExternalID
		modal.ExternalID = message.User.ID + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

		// - metadata : PrivateMeta
		params := privateMeta{
			ChannelID: message.Channel.ID,
		}
		bytes, err := json.Marshal(params)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 200}, fmt.Errorf("failed to marshal private metadata: %w", err)
		}
		modal.PrivateMetadata = string(bytes)

		// Send the view to slack
		api := slack.New(tokenBotUser)
		if _, err := api.OpenView(message.TriggerID, *modal); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 200}, fmt.Errorf("failed to open modal: %w", err)
		}

	case "sushi":
		// In this example, we ignore this case.
	case "ramen":
		// In this example, we ignore this case.
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

// createOrderModalBySDK makes a modal view by using slack-go/slack
func createOrderModalBySDK() *slack.ModalViewRequest {
	// Text section
	shopText := slack.NewTextBlockObject("mrkdwn", ":hamburger: *Hey! Thank you for choosing us! We'll promise you to be full.*", false, false)
	shopTextSection := slack.NewSectionBlock(shopText, nil, nil)

	// Divider
	dividerBlock := slack.NewDividerBlock()

	// Input with radio buttons
	optHamburgerText := slack.NewTextBlockObject("plain_text", burgers["hamburger"] /*"Hamburger"*/, false, false)
	optHamburgerObj := slack.NewOptionBlockObject("hamburger", optHamburgerText)

	optCheeseText := slack.NewTextBlockObject("plain_text", burgers["cheese_burger"] /*"Cheese burger"*/, false, false)
	optCheeseObj := slack.NewOptionBlockObject("cheese_burger", optCheeseText)

	optBLTText := slack.NewTextBlockObject("plain_text", burgers["blt_burger"] /*"BLT burger"*/, false, false)
	optBLTObj := slack.NewOptionBlockObject("blt_burger", optBLTText)

	optBigText := slack.NewTextBlockObject("plain_text", burgers["big_burger"] /*"Big burger"*/, false, false)
	optBigObj := slack.NewOptionBlockObject("big_burger", optBigText)

	optKingText := slack.NewTextBlockObject("plain_text", burgers["king_burger"] /*"King burger"*/, false, false)
	optKingObj := slack.NewOptionBlockObject("king_burger", optKingText)

	menuElement := slack.NewRadioButtonsBlockElement("action_id_menu", optHamburgerObj, optCheeseObj, optBLTObj, optBigObj, optKingObj)

	menuLabel := slack.NewTextBlockObject("plain_text", "Which one you want to have?", false, false)
	menuInput := slack.NewInputBlock("block_id_menu", menuLabel, menuElement)

	// Input with static_select
	optWellDoneText := slack.NewTextBlockObject("plain_text", "well done", false, false)
	optWellDoneObj := slack.NewOptionBlockObject("well_done", optWellDoneText)

	optMediumText := slack.NewTextBlockObject("plain_text", "medium", false, false)
	optMediumObj := slack.NewOptionBlockObject("medium", optMediumText)

	optRareText := slack.NewTextBlockObject("plain_text", "rare", false, false)
	optRareObj := slack.NewOptionBlockObject("rare", optRareText)

	optBlueText := slack.NewTextBlockObject("plain_text", "blue", false, false)
	optBlueObj := slack.NewOptionBlockObject("blue", optBlueText)

	steakInputElement := slack.NewOptionsSelectBlockElement("static_select", nil, "action_id_steak", optWellDoneObj, optMediumObj, optRareObj, optBlueObj)

	steakLabel := slack.NewTextBlockObject("plain_text", "How do you like your steak?", false, false)
	steakInput := slack.NewInputBlock("block_id_steak", steakLabel, steakInputElement)

	// Input with plain_text_input
	noteText := slack.NewTextBlockObject("plain_text", "Anything else you want to tell us?", false, false)
	noteInputElement := slack.NewPlainTextInputBlockElement(nil, "action_id_note")
	noteInputElement.Multiline = true
	noteInput := slack.NewInputBlock("block_id_note", noteText, noteInputElement)
	noteInput.Optional = true

	// Blocks
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			shopTextSection,
			dividerBlock,
			menuInput,
			steakInput,
			noteInput,
		},
	}

	// ModalView
	modal := slack.ModalViewRequest{
		Type:   slack.ViewType("modal"),
		Title:  slack.NewTextBlockObject("plain_text", "Hungryman Hamburgers", false, false),
		Close:  slack.NewTextBlockObject("plain_text", "Cancel", false, false),
		Submit: slack.NewTextBlockObject("plain_text", "Submit", false, false),
		Blocks: blocks,
	}

	return &modal
}

// createOrderModalByJSON makes a modal view by using JSON
func createOrderModalByJSON() (*slack.ModalViewRequest, error) {

	// modal JOSN
	j := `
{
	"type": "modal",
	"submit": {
		"type": "plain_text",
		"text": "Submit",
		"emoji": true
	},
	"close": {
		"type": "plain_text",
		"text": "Cancel",
		"emoji": true
	},
	"title": {
		"type": "plain_text",
		"text": "Hungryman Hamburgers",
		"emoji": true
	},
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": ":hamburger: *Hey! Thank you for choosing us! We'll promise you to be full.*"
			}
		},
		{
			"type": "divider"
		},
		{
			"type": "input",
			"block_id": "block_id_menu",
			"label": {
				"type": "plain_text",
				"text": "Which one you want to have?",
				"emoji": true
			},
			"element": {
				"type": "radio_buttons",
				"action_id": "action_id_menu",
				"options": [
					{
						"text": {
							"type": "plain_text",
							"text": "Hamburger",
							"emoji": true
						},
						"value": "hamburger"
					},
					{
						"text": {
							"type": "plain_text",
							"text": "Cheese Burger",
							"emoji": true
						},
						"value": "cheese_burger"
					},
					{
						"text": {
							"type": "plain_text",
							"text": "BLT Burger",
							"emoji": true
						},
						"value": "blt_burger"
					},
					{
						"text": {
							"type": "plain_text",
							"text": "Big Burger",
							"emoji": true
						},
						"value": "big_burger"
					},
					{
						"text": {
							"type": "plain_text",
							"text": "King Burger",
							"emoji": true
						},
						"value": "king_burger"
					}
				]
			}
		},
		{
			"type": "input",
			"block_id": "block_id_steak",
			"element": {
				"type": "static_select",
				"action_id": "action_id_steak",
				"placeholder": {
					"type": "plain_text",
					"text": "Select ...",
					"emoji": true
				},
				"options": [
					{
						"text": {
							"type": "plain_text",
							"text": "well done",
							"emoji": true
						},
						"value": "well_done"
					},
					{
						"text": {
							"type": "plain_text",
							"text": "medium",
							"emoji": true
						},
						"value": "medium"
					},
					{
						"text": {
							"type": "plain_text",
							"text": "rare",
							"emoji": true
						},
						"value": "rare"
					},
					{
						"text": {
							"type": "plain_text",
							"text": "blue",
							"emoji": true
						},
						"value": "blue"
					}
				]
			},
			"label": {
				"type": "plain_text",
				"text": "How do you like your steak? ",
				"emoji": true
			}
		},
		{
			"type": "input",
			"block_id": "block_id_note",
			"label": {
				"type": "plain_text",
				"text": "Anything else you want to tell us?",
				"emoji": true
			},
			"element": {
				"type": "plain_text_input",
				"action_id": "action_id_note",
				"multiline": true
			},
			"optional": true
		}
	]
}`

	var modal slack.ModalViewRequest
	if err := json.Unmarshal([]byte(j), &modal); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return &modal, nil
}
