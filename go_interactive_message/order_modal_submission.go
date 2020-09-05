package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack"
)

func handleOrderSubmissionRequest(message slack.InteractionCallback) (events.APIGatewayProxyResponse, error) {
	// Get the selected information.
	// - radio button
	menu := message.View.State.Values["block_id_menu"]["action_id_menu"].SelectedOption.Value

	// - static_select
	steak := message.View.State.Values["block_id_steak"]["action_id_steak"].SelectedOption.Value

	// - text
	note := message.View.State.Values["block_id_note"]["action_id_note"].Value

	// Create a confirmation modal.
	// - apperance
	modal := createConfirmationModalBySDK(menu, steak, note)

	// - metadata : CallbackID
	modal.CallbackID = reqConfirmationModalSubmission

	// - metadata : ExternalID
	modal.ExternalID = message.User.ID + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	// - metadata : PrivateMeta
	//   - Get private metadata of a message
	var pMeta privateMeta
	if err := json.Unmarshal([]byte(message.View.PrivateMetadata), &pMeta); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 200}, fmt.Errorf("failed to unmarshal private metadata: %w", err)
	}

	//   - Create new private metadata
	params := privateMeta{
		ChannelID: pMeta.ChannelID,
		order: order{
			Menu:   menu,
			Steak:  steak,
			Note:   note,
			Amount: "700",
		},
	}

	pBytes, err := json.Marshal(params)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 200}, fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	modal.PrivateMetadata = string(pBytes)

	// Create response
	resAction := slack.NewUpdateViewSubmissionResponse(modal)
	rBytes, err := json.Marshal(resAction)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 200}, fmt.Errorf("failed to marshal json: %w", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(rBytes),
	}, nil
}

func createConfirmationModalBySDK(menu, steak, note string) *slack.ModalViewRequest {

	// Create a modal.
	// - Text section
	titleText := slack.NewTextBlockObject("mrkdwn", ":wave: *Order confirmation*", false, false)
	titleTextSection := slack.NewSectionBlock(titleText, nil, nil)

	// Divider
	dividerBlock := slack.NewDividerBlock()

	// - Text section
	sMenuText := slack.NewTextBlockObject("mrkdwn", "*Menu :hamburger:*\n"+burgers[menu], false, false)
	sMenuTextSection := slack.NewSectionBlock(sMenuText, nil, nil)

	// - Text section
	sSteakText := slack.NewTextBlockObject("mrkdwn", "*How do you like your steak?*\n"+steak, false, false)
	sSteakTextSection := slack.NewSectionBlock(sSteakText, nil, nil)

	// - Text section
	sNoteText := slack.NewTextBlockObject("mrkdwn", "*Anything else you want to tell us?*\n"+note, false, false)
	sNoteTextSection := slack.NewSectionBlock(sNoteText, nil, nil)

	// - Text section
	amountText := slack.NewTextBlockObject("mrkdwn", "*Amount :moneybag:*\n$ 700", false, false)
	amountTextSection := slack.NewSectionBlock(amountText, nil, nil)

	// - Input with plain_text_input
	chipText := slack.NewTextBlockObject("plain_text", "Chip ($)", false, false)
	chipInputElement := slack.NewPlainTextInputBlockElement(nil, "action_id_chip")
	chipInput := slack.NewInputBlock("block_id_chip", chipText, chipInputElement)
	chipHintText := slack.NewTextBlockObject("plain_text", "Thank you for your kindness!", false, false)
	chipInput.Hint = chipHintText
	chipInput.Optional = true

	// Blocks
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			titleTextSection,
			dividerBlock,
			sMenuTextSection,
			sSteakTextSection,
			sNoteTextSection,
			dividerBlock,
			amountTextSection,
			chipInput,
		},
	}

	// ModalView
	modal := slack.ModalViewRequest{
		Type:   slack.ViewType("modal"),
		Title:  slack.NewTextBlockObject("plain_text", "Hungryman Hamburgers", false, false),
		Close:  slack.NewTextBlockObject("plain_text", "Cancel", false, false),
		Submit: slack.NewTextBlockObject("plain_text", "Order!", false, false),
		Blocks: blocks,
	}

	return &modal
}
