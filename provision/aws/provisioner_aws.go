package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/tmlbl/rem/config"
)

type Provisioner struct{}

func (p *Provisioner) getAMI(base *config.Base, svc *ec2.Client) (string, error) {
	// If it's already an AMI ID, just use that
	if strings.HasPrefix(base.From, "ami-") {
		return base.From, nil
	}

	// Otherwise assume it is the AMI name, which can be the same in
	// different regions
	out, err := svc.DescribeImages(context.Background(), &ec2.DescribeImagesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("name"),
				Values: []string{base.From},
			},
		},
	})

	if err != nil {
		return "", err
	}

	if len(out.Images) < 1 {
		return "", fmt.Errorf("no AMI found with name: %s", base.From)
	}

	return *out.Images[0].ImageId, nil
}

const securityGroupName = "rem-builder"

// Ensure that the security group exists allowing us to connect via SSH
func (p *Provisioner) ensureSecurityGroup(ctx context.Context, svc *ec2.Client) (string, error) {
	out, err := svc.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("group-name"),
				Values: []string{securityGroupName},
			},
		},
	})
	if err != nil {
		return "", err
	}
	if len(out.SecurityGroups) == 0 {
		grp, err := svc.CreateSecurityGroup(ctx, &ec2.CreateSecurityGroupInput{
			GroupName:   aws.String(securityGroupName),
			Description: aws.String("security group for rem base builds"),
		})
		if err != nil {
			return "", err
		}
		fmt.Println("created security group", grp.GroupId)

		// Add the ingress rule
		_, err = svc.AuthorizeSecurityGroupIngress(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
			IpProtocol: aws.String("tcp"),
			CidrIp:     aws.String("0.0.0.0/0"),
			FromPort:   aws.Int32(22),
			ToPort:     aws.Int32(22),
			GroupId:    grp.GroupId,
		})
		if err != nil {
			return "", err
		}
		return *grp.GroupId, nil
	} else {
		fmt.Println("found existing security group", out.SecurityGroups[0].GroupId)
		return *out.SecurityGroups[0].GroupId, nil
	}
}

func (p *Provisioner) Build(base *config.Base) error {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}
	svc := ec2.NewFromConfig(cfg)

	ami, err := p.getAMI(base, svc)
	if err != nil {
		return err
	}

	// Get the security group ID
	sid, err := p.ensureSecurityGroup(context.Background(), svc)
	if err != nil {
		return err
	}

	// Start the build instance
	out, err := svc.RunInstances(context.Background(), &ec2.RunInstancesInput{
		MaxCount:         aws.Int32(1),
		MinCount:         aws.Int32(1),
		ImageId:          aws.String(ami),
		SecurityGroupIds: []string{sid},
	})
	if err != nil {
		return err
	}

	fmt.Println("started instance", *out.Instances[0].InstanceId)

	return nil
}
