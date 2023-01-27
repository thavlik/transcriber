package transcriber

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func AWSConfigFromEnv() *aws.Config {
	config := aws.NewConfig()
	if v, ok := os.LookupEnv("S3_ENDPOINT"); ok {
		config.Endpoint = aws.String(v)
	}
	return config
}

func AWSSession() *session.Session {
	return session.Must(session.NewSession(AWSConfigFromEnv()))
}

func duplicate(buf []byte) []byte {
	dup := make([]byte, len(buf))
	copy(dup, buf)
	return dup
}
