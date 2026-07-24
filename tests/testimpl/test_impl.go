package common

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/appmesh"
	appmeshtypes "github.com/aws/aws-sdk-go-v2/service/appmesh/types"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/require"
)

// describeVirtualService fetches the virtual service under test from the
// AppMesh API using the mesh_name and virtual_service_name Terraform outputs.
func describeVirtualService(t *testing.T, ctx types.TestContext) *appmeshtypes.VirtualServiceData {
	appmeshClient := appmesh.NewFromConfig(GetAWSConfig(t))
	serviceName := terraform.Output(t, ctx.TerratestTerraformOptions(), "virtual_service_name")
	meshName := terraform.Output(t, ctx.TerratestTerraformOptions(), "mesh_name")

	_, err := appmeshClient.DescribeMesh(context.TODO(), &appmesh.DescribeMeshInput{MeshName: &meshName})
	require.NoError(t, err, "Error getting mesh description")

	output, err := appmeshClient.DescribeVirtualService(context.TODO(), &appmesh.DescribeVirtualServiceInput{
		MeshName:           &meshName,
		VirtualServiceName: &serviceName,
	})
	require.NoError(t, err, "Unable to describe virtual service")
	return output.VirtualService
}

// TestComposableVirtualService is the functional test implementation. It
// verifies the deployed virtual service is active via the AppMesh API.
func TestComposableVirtualService(t *testing.T, ctx types.TestContext) {
	virtualService := describeVirtualService(t, ctx)

	t.Run("TestDoesServiceExist", func(t *testing.T) {
		require.Equal(t, "ACTIVE", string(virtualService.Status.Status), "Expected virtual service to be active")
	})
}

// TestComposableVirtualServiceReadonly is the readonly test implementation.
// It performs the same read-only verification as TestComposableVirtualService
// via the AppMesh API and must not create, update, or delete any resources.
func TestComposableVirtualServiceReadonly(t *testing.T, ctx types.TestContext) {
	virtualService := describeVirtualService(t, ctx)

	t.Run("TestDoesServiceExist", func(t *testing.T) {
		require.Equal(t, "ACTIVE", string(virtualService.Status.Status), "Expected virtual service to be active")
	})
}

func GetAWSConfig(t *testing.T) (cfg aws.Config) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	require.NoErrorf(t, err, "unable to load SDK config, %v", err)
	return cfg
}
