package api

import (
	"context"
	"encoding/json"
	"net/url"

	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	log "github.com/sirupsen/logrus"
)

func GetSecretByIDField(arn, field string, uriEncoded bool) interface{} {
	client := sm.NewFromConfig(getCfg())
	output, err := client.GetSecretValue(context.TODO(),
		&sm.GetSecretValueInput{
			SecretId: &arn,
		})

	if err != nil {
		log.WithFields(log.Fields{
			"reason": "get secret",
			"id":     arn,
		}).Fatal(err)
	}

	secret := *output.SecretString
	var secretMap map[string]interface{}

	err = json.Unmarshal([]byte(secret), &secretMap)
	if err != nil {
		log.WithFields(log.Fields{
			"reason": "parse JSON",
			"id":     arn,
		}).Fatal(err)
	}

	value, ok := secretMap[field]
	if !ok {
		log.WithFields(log.Fields{
			"reason": "read value",
			"id":     arn,
			"field":  field,
		}).Fatalf("cannot find field %s in secret %s", field, arn)
	}

	if strValue, isString := value.(string); isString && uriEncoded {
		value = url.QueryEscape(strValue)
	}

	return value
}
