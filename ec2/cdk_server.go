package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"os"
)

type CdkServerStackProps struct {
	awscdk.StackProps
}

func NewCdkServerStack(scope constructs.Construct, id string, props *CdkServerStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// create an ec2 instance here
	defaultVpc := awsec2.Vpc_FromLookup(stack, jsii.String("DefaultVpc"), &awsec2.VpcLookupOptions{IsDefault: jsii.Bool(true)})
	sg := awsec2.NewSecurityGroup(stack, jsii.String("app_sg"), &awsec2.SecurityGroupProps{Vpc: defaultVpc, AllowAllOutbound: jsii.Bool(true)})
	sg.AddIngressRule(
		awsec2.Peer_AnyIpv4(),
		awsec2.NewPort(&awsec2.PortProps{
			Protocol:             awsec2.Protocol_UDP,
			FromPort:             jsii.Number(500),
			ToPort:               jsii.Number(500),
			StringRepresentation: jsii.String("udp incoming"),
		}),
		jsii.String("allow udp 500"),
		jsii.Bool(false),
	)
	sg.AddIngressRule(
		awsec2.Peer_AnyIpv4(),
		awsec2.NewPort(&awsec2.PortProps{
			Protocol:             awsec2.Protocol_UDP,
			FromPort:             jsii.Number(4500),
			ToPort:               jsii.Number(4500),
			StringRepresentation: jsii.String("udp incoming"),
		}),
		jsii.String("allow udp 4500"),
		jsii.Bool(false),
	)
	sg.AddIngressRule(
		awsec2.Peer_AnyIpv4(),
		awsec2.NewPort(&awsec2.PortProps{
			Protocol:             awsec2.Protocol_TCP,
			FromPort:             jsii.Number(22),
			ToPort:               jsii.Number(22),
			StringRepresentation: jsii.String("ssh incoming"),
		}),
		jsii.String("allow tcp 22"),
		jsii.Bool(false),
	)
	awsec2.NewInstance(stack, jsii.String("testservervpn"), &awsec2.InstanceProps{
		Vpc:          defaultVpc,
		InstanceName: jsii.String("testservervpn"),
		InstanceType: awsec2.InstanceType_Of(awsec2.InstanceClass_MEMORY5, awsec2.InstanceSize_LARGE),
		MachineImage: awsec2.MachineImage_LatestAmazonLinux(&awsec2.AmazonLinuxImageProps{
			Generation: awsec2.AmazonLinuxGeneration_AMAZON_LINUX_2,
		}),
		VpcSubnets:    &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PUBLIC},
		SecurityGroup: sg,
		UserData: awsec2.MultipartUserData_ForLinux(&awsec2.LinuxUserDataOptions{
			Shebang: jsii.String("wget https://git.io/vpnsetup -qO vpn.sh && sudo sh vpn.sh"),
		}),
		KeyName: jsii.String("macYu"),
	})
	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewCdkServerStack(app, "CdkServerStack", &CdkServerStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
