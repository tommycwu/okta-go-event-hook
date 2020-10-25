# okta-go-event-hook

This ia sample okta event hook written in golang.  https://developer.okta.com/docs/concepts/event-hooks/

1) You will need an okta org.  https://developer.okta.com/signup/

2) Set up a an Event hook in your org. https://developer.okta.com/docs/guides/set-up-event-hook/overview/

![](https://d33wubrfki0l68.cloudfront.net/cee289001a7406f9c0907efdef9f055540754837/a26fb/assets-jekyll/blog/easy-user-sync-hooks/event-hook-create-4e5b321eeabab9e2de3064fef5805ef9bce2d25d3e6a50404ec898694ff79d7a.png)

3) The event hook will be invoked by creating a user in your okta org. 

4) The sample app needs to be hosted in a publicly access address with an ssl enable port (https) 
  - This sample includes a makefile and toml file set up specificly for https://www.netlify.com/.

