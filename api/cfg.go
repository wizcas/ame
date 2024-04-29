package api

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	log "github.com/sirupsen/logrus"
)

func getCfg() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.WithField("reason", "AWS config").Fatal(err)
	}
	return cfg
}
