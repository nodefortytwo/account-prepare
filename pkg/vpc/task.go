package vpc

import (
	"time"

	"github.com/nodefortytwo/account-prepare/pkg/regions"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

func Run() error {
	log.Info("Starting to delete default VPCs ")
	for _, region := range regions.GetAllRegions() {
		log.Debugf("looking default VPC Deletion in %s", region)
		sess, err := session.NewSession()
		if err != nil {
			return err
		}
		svc := ec2.New(sess, aws.NewConfig().WithRegion(region))

		defaultVPC, err := getDefaultVPC(svc)
		if err != nil {
			return err
		}

		if defaultVPC == nil {
			log.Debugf("no default VPC in %s", region)
			continue
		}
		log.Warnf("deleting VPC %s and all its dependencies", aws.StringValue(defaultVPC.VpcId))
		err = deleteDependentInternetGateways(defaultVPC, svc)
		if err != nil {
			return err
		}

		err = deleteVPCSubnets(defaultVPC, svc)
		if err != nil {
			return err
		}

		_, err = svc.DeleteVpc(&ec2.DeleteVpcInput{VpcId: defaultVPC.VpcId})
		if err != nil {
			return err
		}
	}
	log.Warn("all default VPCs deleted")
	return nil
}

func deleteDependentInternetGateways(vpc *ec2.Vpc, svc *ec2.EC2) error {
	filter := &ec2.Filter{
		Name:   aws.String("attachment.vpc-id"),
		Values: []*string{vpc.VpcId},
	}

	log.Debug("Looking for dependent internet gateways")
	igResult, err := svc.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{Filters: []*ec2.Filter{filter}})
	if err != nil {
		return err
	}

	for _, ig := range igResult.InternetGateways {
		log.Infof("Deleting internet gateway: %s", aws.StringValue(ig.InternetGatewayId))

		_, err := svc.DetachInternetGateway(&ec2.DetachInternetGatewayInput{VpcId: vpc.VpcId, InternetGatewayId: ig.InternetGatewayId})
		if err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
		_, err = svc.DeleteInternetGateway(&ec2.DeleteInternetGatewayInput{InternetGatewayId: ig.InternetGatewayId})
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteVPCSubnets(vpc *ec2.Vpc, svc *ec2.EC2) error {
	filter := &ec2.Filter{
		Name:   aws.String("vpc-id"),
		Values: []*string{vpc.VpcId},
	}
	log.Debug("Looking for dependent subnets")
	snResult, err := svc.DescribeSubnets(&ec2.DescribeSubnetsInput{Filters: []*ec2.Filter{filter}})
	if err != nil {
		return err
	}

	for _, sn := range snResult.Subnets {
		log.Infof("Deleting subnet: %s", aws.StringValue(sn.SubnetId))
		_, err := svc.DeleteSubnet(&ec2.DeleteSubnetInput{SubnetId: sn.SubnetId})
		if err != nil {
			return err
		}
	}
	return nil
}

func getDefaultVPC(svc *ec2.EC2) (*ec2.Vpc, error) {
	filter := &ec2.Filter{
		Name:   aws.String("isDefault"),
		Values: []*string{aws.String("true")},
	}

	result, err := svc.DescribeVpcs(&ec2.DescribeVpcsInput{Filters: []*ec2.Filter{filter}})
	if err != nil {
		return nil, err
	}

	if len(result.Vpcs) == 0 {
		return nil, nil
	}

	return result.Vpcs[0], nil
}
