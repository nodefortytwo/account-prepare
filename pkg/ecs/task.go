package ecs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	log "github.com/sirupsen/logrus"
)

func Run() error {

	log.Info("Setting ECS Defaults")

	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	svc := ecs.New(sess, aws.NewConfig().WithRegion("eu-west-1"))

	accountSettingValue := "enabled"
	accountSettingNames := []string{"serviceLongArnFormat", "taskLongArnFormat", "containerInstanceLongArnFormat"}

	for _, accountSettingName := range accountSettingNames {
		_, err := svc.PutAccountSettingDefault(&ecs.PutAccountSettingDefaultInput{
			Name:  aws.String(accountSettingName),
			Value: aws.String(accountSettingValue),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
