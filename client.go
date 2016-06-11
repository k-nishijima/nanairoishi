package nanairoishi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func GetMyIP() (string, error) {
	resp, httpErr := http.Get("http://checkip.amazonaws.com/")
	if httpErr != nil {
		return "", httpErr
	}
	defer resp.Body.Close()

	byteArray, ioErr := ioutil.ReadAll(resp.Body)
	if ioErr != nil {
		return "", ioErr
	}

	return strings.TrimRight(string(byteArray), "\n"), nil
}

func getEC2(config SGConfig) *ec2.EC2 {
	cred := credentials.NewSharedCredentials("", config.Profile)
	cfg := aws.NewConfig().WithRegion(config.Region).WithCredentials(cred)
	return ec2.New(session.New(cfg))
}

func AddRule(dryRun bool, config SGConfig) error {
	svc := getEC2(config)
	params := &ec2.AuthorizeSecurityGroupIngressInput{
		DryRun:  aws.Bool(dryRun),
		GroupId: aws.String(config.ID),
		IpPermissions: []*ec2.IpPermission{
			{ // Required
				IpProtocol: aws.String("TCP"),
				IpRanges: []*ec2.IpRange{
					{ // Required
						CidrIp: aws.String(config.IP + "/32"),
					},
				},
				FromPort: aws.Int64(config.Port),
				ToPort:   aws.Int64(config.Port),
			},
		},
	}
	resp, err := svc.AuthorizeSecurityGroupIngress(params)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(resp)
	return nil
}

func RemoveRule(dryRun bool, config SGConfig) error {
	svc := getEC2(config)
	params := &ec2.RevokeSecurityGroupIngressInput{
		DryRun:  aws.Bool(dryRun),
		GroupId: aws.String(config.ID),
		IpPermissions: []*ec2.IpPermission{
			{ // Required
				IpProtocol: aws.String("TCP"),
				IpRanges: []*ec2.IpRange{
					{ // Required
						CidrIp: aws.String(config.IP + "/32"),
					},
				},
				FromPort: aws.Int64(config.Port),
				ToPort:   aws.Int64(config.Port),
			},
		},
	}
	resp, err := svc.RevokeSecurityGroupIngress(params)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(resp)
	return nil
}
