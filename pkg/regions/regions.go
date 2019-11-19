package regions

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func GetAllRegions() []string {
	var regions []string

	resolver := endpoints.DefaultResolver()
	partitions := resolver.(endpoints.EnumPartitions).Partitions()

	for _, p := range partitions {
		if p.ID() == "aws" {
			for _, region := range p.Regions() {
				if haveAccess(region.ID()) {
					regions = append(regions, region.ID())
				}
			}
		}

	}

	return regions
}

func haveAccess(region string) bool {
	sess, err := session.NewSession()
	if err != nil {
		return false
	}
	svc := sts.New(sess, aws.NewConfig().WithRegion(region))
	_, err = svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return false
	}

	return true
}
