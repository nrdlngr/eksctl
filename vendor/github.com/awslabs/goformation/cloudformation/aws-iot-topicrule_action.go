package cloudformation

import (
	"encoding/json"
)

// AWSIoTTopicRule_Action AWS CloudFormation Resource (AWS::IoT::TopicRule.Action)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-action.html
type AWSIoTTopicRule_Action struct {

	// CloudwatchAlarm AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-action.html#cfn-iot-topicrule-action-cloudwatchalarm
	CloudwatchAlarm *AWSIoTTopicRule_CloudwatchAlarmAction `json:"CloudwatchAlarm,omitempty"`

	// CloudwatchMetric AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-action.html#cfn-iot-topicrule-action-cloudwatchmetric
	CloudwatchMetric *AWSIoTTopicRule_CloudwatchMetricAction `json:"CloudwatchMetric,omitempty"`

	// DynamoDB AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-action.html#cfn-iot-topicrule-action-dynamodb
	DynamoDB *AWSIoTTopicRule_DynamoDBAction `json:"DynamoDB,omitempty"`

	// DynamoDBv2 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-action.html#cfn-iot-topicrule-action-dynamodbv2
	DynamoDBv2 *AWSIoTTopicRule_DynamoDBv2Action `json:"DynamoDBv2,omitempty"`

	// Elasticsearch AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-action.html#cfn-iot-topicrule-action-elasticsearch
	Elasticsearch *AWSIoTTopicRule_ElasticsearchAction `json:"Elasticsearch,omitempty"`

	// Firehose AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-action.html#cfn-iot-topicrule-action-firehose
	Firehose *AWSIoTTopicRule_FirehoseAction `json:"Firehose,omitempty"`

	// Kinesis AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-action.html#cfn-iot-topicrule-action-kinesis
	Kinesis *AWSIoTTopicRule_KinesisAction `json:"Kinesis,omitempty"`

	// Lambda AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-action.html#cfn-iot-topicrule-action-lambda
	Lambda *AWSIoTTopicRule_LambdaAction `json:"Lambda,omitempty"`

	// Republish AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-action.html#cfn-iot-topicrule-action-republish
	Republish *AWSIoTTopicRule_RepublishAction `json:"Republish,omitempty"`

	// S3 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-action.html#cfn-iot-topicrule-action-s3
	S3 *AWSIoTTopicRule_S3Action `json:"S3,omitempty"`

	// Sns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-action.html#cfn-iot-topicrule-action-sns
	Sns *AWSIoTTopicRule_SnsAction `json:"Sns,omitempty"`

	// Sqs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-action.html#cfn-iot-topicrule-action-sqs
	Sqs *AWSIoTTopicRule_SqsAction `json:"Sqs,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTTopicRule_Action) AWSCloudFormationType() string {
	return "AWS::IoT::TopicRule.Action"
}

func (r *AWSIoTTopicRule_Action) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
