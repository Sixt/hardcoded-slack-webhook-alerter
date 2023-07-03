package testdata

var blah = "blah"

type NotificationConfig struct {
	Url string `yaml:"url"`
}

var C = NotificationConfig{
	Url: "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXABC",
}