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

	// Start the build instance
	fmt.Println("launching with AMI ID", ami)
	out, err := svc.RunInstances(context.Background(), &ec2.RunInstancesInput{
		MaxCount: aws.Int32(1),
		MinCount: aws.Int32(1),
		ImageId:  aws.String(ami),
	})
	if err != nil {
		return err
	}

	fmt.Println("started instance", *out.Instances[0].InstanceId)

	return nil
}
