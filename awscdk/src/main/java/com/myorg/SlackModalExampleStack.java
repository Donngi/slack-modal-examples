package com.myorg;

import software.amazon.awscdk.core.Construct;
import software.amazon.awscdk.core.Stack;
import software.amazon.awscdk.core.StackProps;
import software.amazon.awscdk.services.apigateway.LambdaRestApi;
import software.amazon.awscdk.services.lambda.Code;
import software.amazon.awscdk.services.lambda.Function;
import software.amazon.awscdk.services.lambda.Runtime;

public class SlackModalExampleStack extends Stack {
    public SlackModalExampleStack(final Construct scope, final String id) {
        this(scope, id, null);
    }

    public SlackModalExampleStack(final Construct scope, final String id, final StackProps props) {
        super(scope, id, props);

        // Lamnda - event handler
        final Function eventLambda = Function.Builder.create(this, "EventHandler")
            .runtime(Runtime.GO_1_X)
            .code(Code.fromAsset("../go_event_message/bin"))
            .handler("main")
            .build();

        // Lamnda - interactive handler
        final Function interactiveLambda = Function.Builder.create(this, "InteractiveHandler")
            .runtime(Runtime.GO_1_X)
            .code(Code.fromAsset("../go_interactive_message/bin"))
            .handler("main")
            .build();

        // API Gateway
        LambdaRestApi.Builder.create(this, "SlackExampleEventEndpoint")
            .handler(eventLambda)
            .build();

        LambdaRestApi.Builder.create(this, "SlackExampleInteractiveHandler")
            .handler(interactiveLambda)
            .build();
    }
}
