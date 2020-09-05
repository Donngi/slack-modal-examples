.PHONY: build, deploy, tidy

build:
	cd "$(PWD)/go_event_message" && make build
	cd "$(PWD)/go_interactive_message" && make build
	cd "$(PWD)/awscdk" && mvn package
	
deploy: build
	cd "$(PWD)/awscdk" && cdk bootstrap ${OPT}
	cd "$(PWD)/awscdk" && cdk deploy ${OPT}

tidy:
	cd "$(PWD)/go_event_message" && make tidy
	cd "$(PWD)/go_interactive_message" && make tidy

update-dependencies:
	cd "$(PWD)/go_event_message" && go get -u
	cd "$(PWD)/go_interactive_message" && go get -u
