package aws_mg

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/george012/aws_mg/aws_mg_common"
	"github.com/george012/aws_mg/aws_mg_ec2"
	"github.com/george012/aws_mg/aws_mg_model"
	"sync"
)

type AWSManager struct {
	Region            aws_mg_common.AWSRegion
	AwsAk             string
	AwsSk             string
	AWSConfig         *aws.Config
	InstancePreConfig *aws_mg_model.AWSInstancePreConfig
	EC2Client         *ec2.Client
}

var (
	aws_mg_once    sync.Once
	currentRegion  aws_mg_common.AWSRegion
	currentManager *AWSManager
)

func instanceOnce() *AWSManager {
	aws_mg_once.Do(func() {
		currentManager = &AWSManager{}
	})

	return currentManager
}

func (aws_mg *AWSManager) listInstance() {
	aws_mg_ec2.ListInstanceFromAWSManager(aws_mg.Region, aws_mg.AWSConfig, aws_mg.EC2Client)

}

func (aws_mg *AWSManager) deleteInstance() {

}

func (aws_mg *AWSManager) createInstanceWithAMIId(ami_id string, end_func func(result_info interface{}, err error)) {
	aws_mg_ec2.CreateInstanceFromAWSManager(aws_mg.Region, aws_mg.AWSConfig, aws_mg.EC2Client, ami_id, end_func)
}

// NewEC2WithRegion 创建EC2
func NewEC2WithRegion(region aws_mg_common.AWSRegion, ami_id string, end_func func(result_info interface{}, err error)) {
	currentRegion = region

	instanceOnce().createInstanceWithAMIId(ami_id, end_func)
}

// SetupAWSManager 初始化工具，仅运行一次。
func SetupAWSManager(ak string, sk string, end_func func(err error)) {

	instanceOnce().Region = currentRegion

	instanceOnce().AwsAk = ak
	instanceOnce().AwsSk = sk

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(currentManager.AwsAk, currentManager.AwsSk, "")),
		config.WithRegion(currentManager.Region.String()),
	)
	if err != nil {
		end_func(errors.New(fmt.Sprintf("无法配置AWS SDK: %s", err.Error())))
		return
	}

	currentManager.AWSConfig = &cfg
	end_func(nil)
}
