# okta-go-event-hook

This is a sample okta event hook written in golang.  https://developer.okta.com/docs/concepts/event-hooks/.  

And is largely based on this blog https://developer.okta.com/blog/2020/07/20/easy-user-sync-hooks.

1) You will need an okta org.  https://developer.okta.com/signup/

2) Next set up a an event hook in your org. https://developer.okta.com/docs/guides/set-up-event-hook/overview/

3) The event hook will be triggered by creating a user in your okta org. https://developer.okta.com/code/rest/

5) The event hook will reach out to your REST based API (webservice) with an endpoint of /userTransfer (case sensitive).

6) The sample app needs to be hosted in a publicly access address with an ssl enable port (https) 
  -This sample includes a makefile and toml file set up specifically for https://www.netlify.com/.

7) In your API, it will need to have a GET request that will look for 'x-okta-verification-challenge' in the header (all lower case).

8) You will respond in the body of that GET request with whatever value that came with the 'x-okta-verification-challenge' header (this is just to verify your endpoint exists).

9) All subsequent request from the okta event hook will come via POST to the /userTransfer and you can put your business logic there.  (For this sample, it will create a random
user in another org.  This was done to save a little time as we need to show that it works but I didn't want to set up a db, make a GET call to get user profile info and to create another GET endpoint, which is what was done in the blog.) 

