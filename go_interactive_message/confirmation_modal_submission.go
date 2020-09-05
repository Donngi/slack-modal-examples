package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack"
)

func handleConfirmationModalSubmissionRequest(message slack.InteractionCallback) (events.APIGatewayProxyResponse, error) {
	// Validate a message.
	if err := validateChip(message); err != nil {
		// Create validation failed response.
		errors := map[string]string{
			"block_id_chip": "[ERROR] Please enter a number.",
		}

		resAction := slack.NewErrorsViewSubmissionResponse(errors)
		bytes, err := json.Marshal(resAction)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 200}, fmt.Errorf("failed to marshal a validation failed message: %w", err)
		}

		return events.APIGatewayProxyResponse{
			StatusCode:      200,
			IsBase64Encoded: false,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: string(bytes),
		}, nil
	}

	// Get private metadata
	var privateMeta privateMeta
	if err := json.Unmarshal([]byte(message.View.PrivateMetadata), &privateMeta); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 200}, fmt.Errorf("failed to unmarshal private metadata: %w", err)
	}

	// Send a complession message.
	// - Create message options
	option, err := createOption(message, privateMeta)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 200}, fmt.Errorf("failed to create message options: %w", err)
	}

	// - Post a message
	api := slack.New(tokenBotUser)
	if _, _, err := api.PostMessage(privateMeta.ChannelID, option); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 200}, fmt.Errorf("failed to send a message: %w", err)
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func validateChip(message slack.InteractionCallback) error {
	// Get an input value.
	chip := message.View.State.Values["block_id_chip"]["action_id_chip"].Value

	// Chech if the value is number or not.
	if _, err := strconv.ParseFloat(chip, 64); err != nil {
		return err
	}
	return nil
}

func createOption(message slack.InteractionCallback, privateMeta privateMeta) (slack.MsgOption, error) {

	// Text section
	titleText := slack.NewTextBlockObject("mrkdwn", ":hamburger: *Thank you for your order !!*", false, false)
	titleTextSection := slack.NewSectionBlock(titleText, nil, nil)

	// Divider
	dividerBlock := slack.NewDividerBlock()

	// Text section
	sMenuText := slack.NewTextBlockObject("mrkdwn", "*Menu*\n"+burgers[privateMeta.Menu], false, false)
	sMenuTextSection := slack.NewSectionBlock(sMenuText, nil, nil)

	// Text section
	sSteakText := slack.NewTextBlockObject("mrkdwn", "*How do you like your steak?*\n"+privateMeta.Steak, false, false)
	sSteakTextSection := slack.NewSectionBlock(sSteakText, nil, nil)

	// Text section
	sNoteText := slack.NewTextBlockObject("mrkdwn", "*Anything else you want to tell us?*\n"+privateMeta.Note, false, false)
	sNoteTextSection := slack.NewSectionBlock(sNoteText, nil, nil)

	// Text section
	amount, err := strconv.ParseFloat(privateMeta.Amount, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to convert amount to float64: %w", err)
	}

	chip, err := strconv.ParseFloat(message.View.State.Values["block_id_chip"]["action_id_chip"].Value, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to convert amount to float64: %w", err)
	}

	amountText := slack.NewTextBlockObject("mrkdwn", "*Total amount :moneybag:*\n$ "+strconv.FormatFloat(amount+chip, 'f', 2, 64), false, false)
	amountTextSection := slack.NewSectionBlock(amountText, nil, nil)

	// Blocks
	blocks := slack.MsgOptionBlocks(
		titleTextSection,
		dividerBlock,
		sMenuTextSection,
		sSteakTextSection,
		sNoteTextSection,
		dividerBlock,
		amountTextSection,
	)
	return blocks, nil
}
