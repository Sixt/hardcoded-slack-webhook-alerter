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

Exclusions for channels and webhooks can be added to `fullTextValueExcludePatterns` in `scanner/root_config.yaml`.

This tool is assuming that you the scan paths you provide contain GitHub repositories at the top level, each named after the repository. For example if you provide the directory "/Users/me/go", it expects each subdirectory to be a Git repo (this is how it decides the repository name and how it builds the GitHub links). If you don't structure things this way, the scanner will still find results, but the links won't be valid. Also note that when building the links, it's using `master`, which may not always be correct.

### Scripts
Since this tool requires you to have a lot of repositories cloned and updated, there are a few helpful scripts included in this repository. In `scripts/fetch-services.sh`, you can add a list of GitHub repos and your organization name, and it'll clone all of them into the current directory. You can update the list of directories in `scripts/update-all-git-repos.sh` to do a git pull of master/main in every repo you have cloned.