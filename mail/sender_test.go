package mail

import (
	"testing"

	"github.com/mativm02/bank_system/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestSendEmailWithGmail test in short mode")
	}
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Test subject"
	content := `
	<h1>Test content</h1>
	<p>This is a test message from <a href="github.com/mativm02"> Tech School </a></p>
	`

	to := []string{"mancute2010@gmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)

}
