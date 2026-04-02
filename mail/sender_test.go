package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
	"pxsemic.com/simplebank/util"
)

func TestSender(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config, err := util.LoadConfig("../")
	require.NoError(t, err)
	require.NotNil(t, config)
	// Test the sender
	sender := NewYMailSend(config.EmailSenderName, config.EmailSenderEmail, config.EmailSenderPassword)
	subject := "A test mail"
	content := `
		<h1>Hello, World!</h1>
		<p>This is a test mail. <a href="https://www.baidu.com">Click here</a></p>
	`
	to := []string{"1539490160@qq.com"}
	attachFiles := []string{"../app.env"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
