package scanner

import (
	"testing"

	"github.com/sixt/hardcoded-slack-webhook-alerter/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectSlackWebhooks(t *testing.T) {
	conf, err := config.LoadConfig("../testdata/test_config.yaml")
	require.NoError(t, err)

	scanner := New(conf)
	scanner.rootConfigPath = "../scanner/root_config.yaml"
	scanner.Scan()

	slackChans := []string{}
	for _, r := range scanner.Results {
		for c := range r.Channels {
			slackChans = append(slackChans, c)
		}
	}

	assert.Len(t, scanner.Results, 2)
	assert.Len(t, slackChans, 1)
}
