package common

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/appmesh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/require"
)

func TestVirtualService(t *testing.T, ctx types.TestContext) {
	appmeshClient := appmesh.NewFromConfig(GetAWSConfig(t))
	serviceName := terraform.Output(t, ctx.TerratestTerraformOptions(), "virtual_service_name")
	meshName := terraform.Output(t, ctx.TerratestTerraformOptions(), "mesh_name")

	_, err := appmeshClient.DescribeMesh(context.TODO(), &appmesh.DescribeMeshInput{MeshName: &meshName})
	if err != nil {
		t.Errorf("Error getting mesh description: %v", err)
	}

	output, err := appmeshClient.DescribeVirtualService(context.TODO(), &appmesh.DescribeVirtualServiceInput{
		MeshName:           &meshName,
		VirtualServiceName: &serviceName,
	})
	if err != nil {
		t.Errorf("Unable to describe virtual service, %v", err)
	}
	virtualService := output.VirtualService

	t.Run("TestDoesServiceExist", func(t *testing.T) {
		require.Equal(t, "ACTIVE", string(virtualService.Status.Status), "Expected virtual service to be active")
	})
}

func GetAWSConfig(t *testing.T) (cfg aws.Config) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	require.NoErrorf(t, err, "unable to load SDK config, %v", err)
	return cfg
}
