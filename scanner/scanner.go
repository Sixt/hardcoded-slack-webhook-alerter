package scanner

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/sixt/hardcoded-slack-webhook-alerter/config"
	"github.com/sixt/hardcoded-slack-webhook-alerter/model"
	credDectConf "github.com/ynori7/credential-detector/config"
	"github.com/ynori7/credential-detector/parser"
)

const (
	githubUrl = "https://github.com/%s/%s/tree/master%s"
	slackUrl  = "https://hooks.slack.com"
)

// Scanner is used to scan the configured paths for hardcoded slack webhooks and channels
type Scanner struct {
	rootConfigPath    string
	conf              *config.Config
	slackChannelRegex *regexp.Regexp
	urlRegex          *regexp.Regexp //used to extract the webhook in case the result contains the whole line

	Results map[string]model.Result
}

// New returns a new scanner
func New(conf *config.Config) *Scanner {
	return &Scanner{
		conf:              conf,
		rootConfigPath:    "scanner/root_config.yaml",
		Results:           make(map[string]model.Result),
		slackChannelRegex: regexp.MustCompile(conf.ChannelPattern),
		urlRegex:          regexp.MustCompile(`(https:\/\/[a-zA-Z0-9\.\/]*)`),
	}
}

// Scan scans the configured paths for hardcoded slack webhooks and channels
func (s *Scanner) Scan() {
	s.CollectSlackWebhooks()
	s.CollectSlackChannels()
}

// CollectSlackWebhooks scans the given paths for hardcoded slack webhooks
func (s *Scanner) CollectSlackWebhooks() {
	conf, err := credDectConf.LoadConfig("", s.rootConfigPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, dir := range s.conf.Directories {
		p := parser.NewParser(conf)
		if err := p.Scan(dir); err != nil {
			log.Fatal(err.Error())
		}

		for _, res := range p.Results {
			repoName := s.getRepoNameFromFilePath(res.File, dir)
			webhook := s.trimWebhook(res.Value)
			key := s.getUniqueRepoAndWebhookString(repoName, webhook)

			if _, ok := s.Results[key]; !ok {
				s.Results[key] = model.Result{
					Repo:     repoName,
					Webhook:  webhook,
					URL:      s.getGithubUrlFromFilePath(res.File, dir),
					Files:    []string{res.File},
					Channels: make(map[string]struct{}),
				}
			} else {
				r := s.Results[key]
				r.Files = append(r.Files, res.File)
				s.Results[key] = r
			}
		}
	}
}

// CollectSlackChannels scans all the files with results to see if they also contain a slack channel
func (s *Scanner) CollectSlackChannels() {
	conf, err := credDectConf.LoadConfig("", s.rootConfigPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	conf.ValueMatchPatterns = []credDectConf.ValueMatchPattern{
		{
			Name:    "Slack Channel",
			Pattern: s.conf.ChannelPattern,
		},
	}

	filePaths := map[string]string{}
	for key, res := range s.Results {
		for _, f := range res.Files {
			filePaths[f] = key
		}
	}
	for dir, key := range filePaths {
		p := parser.NewParser(conf)
		if err := p.Scan(dir); err != nil {
			log.Fatal(err.Error())
		}

		for _, res := range p.Results {
			if _, ok := s.Results[key]; ok {
				channel := s.trimSlackChannel(res.Value)
				if channel != "" {
					r := s.Results[key]
					r.Channels[channel] = struct{}{}
					s.Results[key] = r
				}
			}
		}
	}
}

// getUniqueRepoAndWebhookString returns a string that is unique for a given repo and webhook
func (s *Scanner) getUniqueRepoAndWebhookString(repoName string, webhook string) string {
	return repoName + "_" + webhook
}

// getRepoNameFromFilePath returns the repo name from the given file path
func (s *Scanner) getRepoNameFromFilePath(path string, basePath string) string {
	repoName := strings.Replace(path, basePath, "", 1)
	repoName = strings.TrimPrefix(repoName, "/")

	slash := strings.Index(repoName, "/")
	if slash == -1 { //happens during tests
		return "."
	}
	return repoName[0:slash]
}

// getGithubUrlFromFilePath returns the github url for the given file path
func (s *Scanner) getGithubUrlFromFilePath(path string, basePath string) string {
	repoName := strings.Replace(path, basePath, "", 1)
	repoName = strings.TrimPrefix(repoName, "/")

	slash := strings.Index(repoName, "/")
	if slash == -1 { //happens during tests
		return fmt.Sprintf(githubUrl, s.conf.GithubOrg, repoName, "")
	}

	filePath := repoName[slash:]
	repoName = repoName[0:slash]

	return fmt.Sprintf(githubUrl, s.conf.GithubOrg, repoName, filePath)
}

func (s *Scanner) trimWebhook(h string) string {
	h = strings.Trim(h, " \t\n\r;\",'") //trim the value of any whitespace or quotes
	if !strings.HasPrefix(h, slackUrl) {
		matches := s.urlRegex.FindStringSubmatch(h)
		if len(matches) > 0 {
			h = matches[1]
		}
	}
	return h
}

func (s *Scanner) trimSlackChannel(c string) string {
	matches := s.slackChannelRegex.FindStringSubmatch(c)
	if len(matches) > 0 {
		return matches[1]
	}

	return ""
}
