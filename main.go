package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	log "github.com/sirupsen/logrus"

	"github.com/nodefortytwo/account-prepare/pkg/vpc"
)

func main() {

	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	svc := sts.New(sess)
	out, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("running on account: %s", aws.StringValue(out.Account))
	if !confirm("Is that correct?") {
		log.Warning("user aborted")
		os.Exit(0)
	}
	err = vpc.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func confirm(question string) bool {
	var s string

	log.Infof("%s (y/N): ", question)
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}
