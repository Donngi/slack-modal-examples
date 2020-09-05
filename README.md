# Slack Modal Example
![DEMO](https://github.com/nicoJN/slack-modal-examples/images/sample.gif)

This is a simple example of slack modal application written in Go. It contains fundamental functions like 
- Authorization
- Sending a interactice message
- Sending a modal
- Updating a modal
- Validating a modal input value

Feel free to fork this repository and customize for your own!

**see more detail (in Japanese)**: [NICOA](https://jimon.info/slack-modal-example-go/)

## Usage
If you want to try this app, set tokens in two files.
- go_event_message/main.go
- go_interactive_message/main.go

```
var (
	signingSecret = "YOUR_SIGNING_SECRET_HERE!"
	tokenBotUser  = "YOUR_BOT_USER_OAUTH_ACCESS_TOKEN_HERE!"
)
```

This example includes awscdk setting files. You can easily deploy with AWS CDK.

At the project root, enter these commands.

```
$ make deploy OPT="--profile YOUR_AWS_PROFILE_HERE!!!"
```

or directly execute cdk commands at root/awsdk.

```
$ cdk bootstrap --profile YOUR_AWS_PROFILE_HERE!!!
$ cdk deploy --profile YOUR_AWS_PROFILE_HERE!!!
```

## License
MIT