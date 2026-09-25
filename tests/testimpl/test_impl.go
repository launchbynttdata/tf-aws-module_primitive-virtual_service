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

type serviceVerification struct {
	client      *appmesh.Client
	serviceArn  string
	serviceName string
	meshName    string
}

// TestComposableVirtualService runs read-only assertions first, then performs
// a small mutating operation (temporary tag add/remove) to prove write
// behavior.
func TestComposableVirtualService(t *testing.T, ctx types.TestContext) {
	verification := verifyServiceReadOnly(t, ctx)
	runServiceTagWriteProbe(t, verification.client, verification.serviceArn)
}

// TestComposableVirtualServiceReadOnly validates the deployed virtual service
// via read-only SDK calls only. It shares verifyServiceReadOnly with the
// functional test but never invokes a mutating operation.
func TestComposableVirtualServiceReadOnly(t *testing.T, ctx types.TestContext) {
	verifyServiceReadOnly(t, ctx)
}

func verifyServiceReadOnly(t *testing.T, ctx types.TestContext) serviceVerification {
	t.Helper()

	appmeshClient := appmesh.NewFromConfig(GetAWSConfig(t))
	serviceName := terraform.OutputContext(t, context.Background(), ctx.TerratestTerraformOptions(), "virtual_service_name")
	meshName := terraform.OutputContext(t, context.Background(), ctx.TerratestTerraformOptions(), "mesh_name")

	_, err := appmeshClient.DescribeMesh(context.TODO(), &appmesh.DescribeMeshInput{MeshName: &meshName})
	require.NoErrorf(t, err, "error getting mesh description, %v", err)

	output, err := appmeshClient.DescribeVirtualService(context.TODO(), &appmesh.DescribeVirtualServiceInput{
		MeshName:           &meshName,
		VirtualServiceName: &serviceName,
	})
	require.NoErrorf(t, err, "unable to describe virtual service, %v", err)
	virtualService := output.VirtualService

	t.Run("TestDoesServiceExist", func(t *testing.T) {
		require.Equal(t, "ACTIVE", string(virtualService.Status.Status), "Expected virtual service to be active")
	})

	return serviceVerification{
		client:      appmeshClient,
		serviceArn:  *virtualService.Metadata.Arn,
		serviceName: serviceName,
		meshName:    meshName,
	}
}

// runServiceTagWriteProbe proves the module's write path by adding a
// temporary tag to the virtual service, verifying it took effect, then
// removing it. It must only be called from the functional (non-readonly)
// test path.
func runServiceTagWriteProbe(t *testing.T, client *appmesh.Client, serviceArn string) {
	t.Run("CanTagAndUntagVirtualService", func(t *testing.T) {
		const probeKey = "lcaf-readonly-probe"
		const probeValue = "terratest"

		_, err := client.TagResource(context.TODO(), &appmesh.TagResourceInput{
			ResourceArn: &serviceArn,
			Tags: []appmeshtypes.TagRef{
				{Key: aws.String(probeKey), Value: aws.String(probeValue)},
			},
		})
		require.NoErrorf(t, err, "unable to tag virtual service, %v", err)

		defer func() {
			_, err := client.UntagResource(context.TODO(), &appmesh.UntagResourceInput{
				ResourceArn: &serviceArn,
				TagKeys:     []string{probeKey},
			})
			require.NoErrorf(t, err, "unable to untag virtual service, %v", err)
		}()

		tagsOutput, err := client.ListTagsForResource(context.TODO(), &appmesh.ListTagsForResourceInput{
			ResourceArn: &serviceArn,
		})
		require.NoErrorf(t, err, "unable to list tags for virtual service, %v", err)

		var found bool
		for _, tag := range tagsOutput.Tags {
			if aws.ToString(tag.Key) == probeKey {
				require.Equal(t, probeValue, aws.ToString(tag.Value), "Expected probe tag value to be %s, but got %s", probeValue, aws.ToString(tag.Value))
				found = true
				break
			}
		}
		require.True(t, found, "Expected probe tag %s to be present after TagResource", probeKey)
	})
}

func GetAWSConfig(t *testing.T) (cfg aws.Config) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	require.NoErrorf(t, err, "unable to load SDK config, %v", err)
	return cfg
}
