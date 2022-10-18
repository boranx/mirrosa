package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/mjlshen/mirrosa/pkg/rosa"

	"github.com/mjlshen/mirrosa/pkg/mirrosa"
)

func main() {
	clusterId := flag.String("cluster-id", "", "Cluster ID")
	flag.Parse()

	if *clusterId == "" {
		panic("cluster id must not be empty")
	}

	mirrosa, err := mirrosa.NewClient(context.TODO(), *clusterId)
	if err != nil {
		panic(err)
	}

	if err := ValidateAll(context.TODO(), mirrosa); err != nil {
		panic(err)
	}

	fmt.Printf("%s, \"Mirror mirror on the wall, who's the fairest of them all?\"\n%+v\n", mirrosa.ClusterInfo.Name, *mirrosa.ClusterInfo)
}

// ValidateAll runs Validate against all known ROSA components
func ValidateAll(ctx context.Context, c *mirrosa.Client) error {
	vpc := rosa.NewVpc(c.Cluster, ec2.NewFromConfig(c.AwsConfig))
	vpcId, err := c.ValidateComponent(ctx, vpc)
	if err != nil {
		fmt.Println(vpc.Documentation())
		return err
	}

	c.ClusterInfo.VpcId = vpcId

	privateHz := rosa.NewPrivateHostedZone(c.Cluster, route53.NewFromConfig(c.AwsConfig), c.ClusterInfo.VpcId)
	privateHzId, err := c.ValidateComponent(ctx, privateHz)
	if err != nil {
		fmt.Println(privateHz.Documentation())
		return err
	}

	c.ClusterInfo.PrivateHostedZoneId = privateHzId

	return nil
}
