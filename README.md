# Hardcoded Slack Webhook Alerter
This tool can be used to scan for hardcoded Slack webhooks and send an alert to each one about the dangers of hardcoded credentials. This tool was built using [Credential-Detector](https://github.com/ynori7/credential-detector).

### How it works
1. It combs all repos locally cloned in the specified directories and collects the Slack webhooks.
2. Then for each file where it found a webhook, it checks if it finds any slack channels hardcoded as well and links them to the webhook/repo combination.
3. Then it iterates all the webhook/repo combos it found
    - If there's no channel associated with it, then it just sends a payload to the webhook and lets it go to whatever default channel is configured.
    - If it does have channels associated with the webhook/repo, it sends the payload to each of them.

### Usage

```
go run main.go
```

The following configuration should be set in config/config.yaml: 

|Config Option|Description|Values|
|-------------|-----------|------|
|channel_pattern|The pattern used to detect Slack channels. The default pattern is looking for things shaped like this: #something_somethingelse (with an underscore)|A regular expression string|
|directories|A list of directories to scan for hard-coded Slack webhooks|A list of strings|
|github_org|The name of your Github organization (used for building links)|A string|
|dry_run|When true, it'll output all the logs without actually sending anything to Slack|Boolean|
