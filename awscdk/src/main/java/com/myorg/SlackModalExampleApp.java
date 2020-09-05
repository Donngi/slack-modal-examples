package com.myorg;

import software.amazon.awscdk.core.App;

public class SlackModalExampleApp {
    public static void main(final String[] args) {
        App app = new App();
        
        new SlackModalExampleStack(app, "SlackModalExampleStack");

        app.synth();
    }
}
