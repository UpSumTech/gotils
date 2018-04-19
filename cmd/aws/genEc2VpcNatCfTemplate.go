package aws

import (
	"github.com/awslabs/goformation/cloudformation"
	"github.com/sumanmukherjee03/gotils/cmd/utils"
)

func GenEc2VpcNatCfTemplate() (string, error) {
	var data string
	template := cloudformation.NewTemplate()

	template.Resources["MySNSTopic"] = &cloudformation.AWSSNSTopic{
		DisplayName: "test-sns-topic-display-name",
		TopicName:   "test-sns-topic-name",
		Subscription: []cloudformation.AWSSNSTopic_Subscription{
			cloudformation.AWSSNSTopic_Subscription{
				Endpoint: "test-sns-topic-subscription-endpoint",
				Protocol: "test-sns-topic-subscription-protocol",
			},
		},
	}

	template.Resources["MyRoute53HostedZone"] = &cloudformation.AWSRoute53HostedZone{
		Name: "example.com",
	}

	data, err := utils.ToJson(template)
	if err != nil {
		return data, err
	}
	return data, nil
}
