package test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/appmesh"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

var standardTags = map[string]string{
	"provisioner": "Terraform",
}

const (
	base            = "../../examples/"
	testVarFileName = "/test.tfvars"
	caModule        = "module.private_ca"
)

// var standardTags = map[string]string{
// 	"provisioner": "Terraform",
// }

func TestAppMeshVirtualService(t *testing.T) {
	t.Parallel()
	stage := test_structure.RunTestStage

	files, err := os.ReadDir(base)
	if err != nil {
		assert.Error(t, err)
	}
	for _, file := range files {
		dir := base + file.Name()
		if file.IsDir() {
			defer stage(t, "teardown_appmesh_virtualservice", func() { tearDownAppMeshVirtualService(t, dir) })
			stage(t, "setup_and_test_appmesh_virtualservice", func() { setupAndTestAppMeshVirtualService(t, dir) })
		}
	}
}

func setupAndTestAppMeshVirtualService(t *testing.T, dir string) {
	varsFilePath := []string{dir + testVarFileName}

	terraformOptionsCA := &terraform.Options{
		TerraformDir: dir,
		Targets:      []string{caModule},
		VarFiles:     varsFilePath,
		Logger:       logger.Discard,

		// Disable colors in Terraform commands so its easier to parse stdout/stderr
		NoColor: true,
	}
	terraformOptionsCA.VarFiles = varsFilePath

	terraformOptions := &terraform.Options{
		TerraformDir: dir,
		VarFiles:     []string{dir + testVarFileName},
		NoColor:      true,
		Logger:       logger.Discard,
	}

	test_structure.SaveTerraformOptions(t, dir, terraformOptions)

	terraform.InitAndApply(t, terraformOptionsCA)
	// sleep for 3 minutes for the CA status change to ISSUED
	time.Sleep(3 * time.Minute)
	terraform.InitAndApply(t, terraformOptions)
	expectedPatternARN := "^arn:aws:appmesh:[a-z]{2}-[a-z]+-[0-9]{1}:[0-9]{12}:mesh/[a-zA-Z0-9-]+/virtualService/[a-zA-Z0-9-]+$"

	actualVirtualServiceId := terraform.Output(t, terraformOptions, "virtual_service_id")
	assert.NotEmpty(t, actualVirtualServiceId, "Virtual Service Id is empty")
	actualServiceARN := terraform.Output(t, terraformOptions, "service_arn")
	assert.Regexp(t, expectedPatternARN, actualServiceARN, "ARN does not match expected pattern")
	actualRandomId := terraform.Output(t, terraformOptions, "random_int")
	assert.NotEmpty(t, actualRandomId, "Random ID is empty")
	actualVirtualNodeName := terraform.Output(t, terraformOptions, "virtual_node_name")
	assert.NotEmpty(t, actualVirtualNodeName, "Virtual Node Name is empty")
	actualVpcID := terraform.Output(t, terraformOptions, "vpc_id")
	assert.NotEmpty(t, actualVpcID, "Vpc ID  is empty")

	expectedNamePrefix := terraform.GetVariableAsStringFromVarFile(t, dir+testVarFileName, "naming_prefix")
	expectedMeshName := expectedNamePrefix + "-app-mesh-" + actualRandomId
	expectedVirtualServiceName := expectedNamePrefix + "-vsvc-" + actualRandomId

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(os.Getenv("AWS_PROFILE")),
	)
	if err != nil {
		assert.Error(t, err, "can't connect to aws")
	}

	client := appmesh.NewFromConfig(cfg)
	input := &appmesh.DescribeVirtualServiceInput{
		VirtualServiceName: aws.String(expectedVirtualServiceName),
		MeshName:           aws.String(expectedMeshName),
	}
	result, err := client.DescribeVirtualService(context.TODO(), input)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("The Expected Virtual Service was not found %s", err.Error()))

	}
	virtualService := result.VirtualService
	expectedId := *virtualService.Metadata.Uid
	expectedArn := *virtualService.Metadata.Arn

	assert.Regexp(t, expectedPatternARN, actualServiceARN, "SVC ARN does not match expected pattern")
	assert.Equal(t, expectedArn, actualServiceARN, "SVC ARN does not match")
	assert.Equal(t, expectedId, actualVirtualServiceId, "svc id does not match")
	checkTagsMatch(t, dir, actualServiceARN, client)

}

func checkTagsMatch(t *testing.T, dir string, actualARN string, client *appmesh.Client) {
	expectedTags, err := terraform.GetVariableAsMapFromVarFileE(t, dir+testVarFileName, "tags")
	if err == nil {
		result2, errListTags := client.ListTagsForResource(context.TODO(), &appmesh.ListTagsForResourceInput{ResourceArn: aws.String(actualARN)})
		if errListTags != nil {
			assert.Error(t, errListTags, "Failed to retrieve tags from AWS")
		}
		// convert AWS Tag[] to map so we can compare
		actualTags := map[string]string{}
		for _, tag := range result2.Tags {
			actualTags[*tag.Key] = *tag.Value
		}

		// add the standard tags to the expected tags
		for k, v := range standardTags {
			expectedTags[k] = v
		}
		expectedTags["env"] = actualTags["env"]
		assert.True(t, reflect.DeepEqual(actualTags, expectedTags), fmt.Sprintf("tags did not match, expected: %v\nactual: %v", expectedTags, actualTags))
	}
}

func tearDownAppMeshVirtualService(t *testing.T, dir string) {
	terraformOptions := test_structure.LoadTerraformOptions(t, dir)
	terraformOptions.Logger = logger.Discard
	terraform.Destroy(t, terraformOptions)

}
