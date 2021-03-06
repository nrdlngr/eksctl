package cloudformation

import (
	"encoding/json"
)

// AWSEC2Instance_InstanceIpv6Address AWS CloudFormation Resource (AWS::EC2::Instance.InstanceIpv6Address)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-instance-instanceipv6address.html
type AWSEC2Instance_InstanceIpv6Address struct {

	// Ipv6Address AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-instance-instanceipv6address.html#cfn-ec2-instance-instanceipv6address-ipv6address
	Ipv6Address *Value `json:"Ipv6Address,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2Instance_InstanceIpv6Address) AWSCloudFormationType() string {
	return "AWS::EC2::Instance.InstanceIpv6Address"
}

func (r *AWSEC2Instance_InstanceIpv6Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
