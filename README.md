# GWARR
Golang Webhooks (for) \*arr

# Description 

I was unhappy with the current state of \*arr notifications sent to Slack via the inbuilt notifications.

I didn't want to learn how to work on the \*arr codebase initially, so I just decided to write a little proxy.

It sits on a host you own, you point your \*arr notifications at it, and it makes modern well formatted Slack messages.

# Features

* Updates messages for certain state changes (added -> grabbed -> deleted)
 * Redis is used as the cache now. (There might be some bugs)
* Emojis :star:
* Links back to your \*arr instance
* Adds some (currently very limited) metadata to the message
* Unfurls
* a prom /metrics endpoint

# Running

* Create a Slack application [here](https://api.slack.com/apps)
* Create a token with the `chat:write` scope
* Install the application to your workspace
* Take note of the token. It should begin with `xoxb`
* Add the application to the channel of your choosing and note the channel ID
* Add the following variables to your environment:
```bash
GWARR_RADARR_URL='<url with proto of radarr instance>'
GWARR_SLACK_CHANNEL_ID='<channel id>'
GWARR_SLACK_BOT_TOKEN='<bot token>'
```
* Build the binary `go build cmd/gwarr/gwarr.go`
* Run GWARR `./gwarr`

# Planned

* Fix `golangci-lint` errors
* Handle \*arr test notifications
* Add some images
* Add other \*arrs support
* More testing
* Refine the message types
* Clean up the code
 * It works, but it needs some TLC
* More configuration options
* A docker container
* Add version details
* Add Makefile/build stuff
* Add Github Actions
* Improve logging
