package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

// InitSession :
func InitSession() *session.Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return sess
}
