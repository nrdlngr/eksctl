package cloudformation

import (
	"encoding/json"
)

// AWSElasticLoadBalancingV2LoadBalancer_SubnetMapping AWS CloudFormation Resource (AWS::ElasticLoadBalancingV2::LoadBalancer.SubnetMapping)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-loadbalancer-subnetmapping.html
type AWSElasticLoadBalancingV2LoadBalancer_SubnetMapping struct {

	// AllocationId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-loadbalancer-subnetmapping.html#cfn-elasticloadbalancingv2-loadbalancer-subnetmapping-allocationid
	AllocationId *Value `json:"AllocationId,omitempty"`

	// SubnetId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-loadbalancer-subnetmapping.html#cfn-elasticloadbalancingv2-loadbalancer-subnetmapping-subnetid
	SubnetId *Value `json:"SubnetId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticLoadBalancingV2LoadBalancer_SubnetMapping) AWSCloudFormationType() string {
	return "AWS::ElasticLoadBalancingV2::LoadBalancer.SubnetMapping"
}

func (r *AWSElasticLoadBalancingV2LoadBalancer_SubnetMapping) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
