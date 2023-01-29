package util

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func AWSConfigFromEnv() *aws.Config {
	config := aws.NewConfig()
	if v, ok := os.LookupEnv("AWS_REGION"); ok {
		config.Region = aws.String(v)
	} else {
		config.Region = aws.String("us-east-1")
	}
	return config
}

func AWSSession() *session.Session {
	return session.Must(session.NewSession(AWSConfigFromEnv()))
}

func Duplicate(buf []byte) []byte {
	dup := make([]byte, len(buf))
	copy(dup, buf)
	return dup
}
