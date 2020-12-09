package tests

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"golang.org/x/crypto/ssh"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/resourcegroups"
)

const (
	awsbiImageTag = "epiphanyplatform/awsbi:0.0.1"
	awsksImageTag = "epiphanyplatform/awsks:0.0.1"
	awsTagName  = "resource_group"
	awsTagValue = "bi-module"
	moduleName  = "bi-module"
	awsRegion   = "eu-central-1"
	sshKeyName  = "vms_rsa"
	retries     = 30
)

func TestInit(t *testing.T) {
	tests := []struct{
		name               string
		initParams         []string
		stateLocation      string
		stateContent       string
		wantOutput         string
		wantConfigLocation string
		wantConfigContent  string
		wantStateContent   string
	}{
		{
			name: "init with defaults",
			initParams: nil,
			stateLocation: "state.yml",
			stateContent: ``,
			wantOutput: `
#AWSKS | setup | ensure required directories
#AWSKS | ensure-state-file | checks if state file exists
#AWSKS | template-config-file | will template config file (and backup previous if exists)
#AWSKS | template-config-file | will replace arguments with values from state file
#AWSKS | initialize-state-file | will initialize state file
#AWSKS | display-config-file | config file content is:
kind: awsks-config
awsks:
  name: epiphany
  vpc_id: unset
  region: eu-central-1
  subnet_ids: null
  private_route_table_id: unset
`,
			wantConfigLocation: "awsks/awsks-config.yml",
			wantConfigContent: `
kind: awsks-config
awsks:
  name: epiphany
  vpc_id: unset
  region: eu-central-1
  subnet_ids: null
  private_route_table_id: unset
`,
			wantStateContent: `
kind: state
awsks:
  status: initialized
`,
		},
		{
			name: "init with variables",
			initParams: []string{"M_NAME=value1", "M_VPC_ID=value2", "M_REGION=value3", "M_SUBNET_IDS=value4"},
			stateLocation: "state.yml",
			stateContent: ``,
			wantOutput: `
#AWSKS | setup | ensure required directories
#AWSKS | ensure-state-file | checks if state file exists
#AWSKS | template-config-file | will template config file (and backup previous if exists)
#AWSKS | template-config-file | will replace arguments with values from state file
#AWSKS | initialize-state-file | will initialize state file
#AWSKS | display-config-file | config file content is:
kind: awsks-config
awsks:
  name: value1
  vpc_id: value2
  region: value3
  subnet_ids: value4
  private_route_table_id: unset
`,
			wantConfigLocation: "awsks/awsks-config.yml",
			wantConfigContent: `
kind: awsks-config
awsks:
  name: value1
  vpc_id: value2
  region: value3
  subnet_ids: value4
  private_route_table_id: unset
`,
			wantStateContent: `
kind: state
awsks:
  status: initialized
`,
		},
		{
			name: "init with state",
			initParams: nil,
			stateLocation: "state.yml",
			stateContent: `
kind: state
awsbi:
  status: applied
  name: epiphany
  instance_count: 0
  region: eu-central-1
  use_public_ip: false
  force_nat_gateway: true
  rsa_pub_path: "/shared/vms_rsa.pub"
  output:
    private_ip.value: []
    public_ip.value: []
    public_subnet_id.value: subnet-0137cf1e7921c1551
    vpc_id.value: vpc-0baa2c4e9e48e608c
`,
			wantOutput: `
#AWSKS | setup | ensure required directories
#AWSKS | ensure-state-file | checks if state file exists
#AWSKS | template-config-file | will template config file (and backup previous if exists)
#AWSKS | template-config-file | will replace arguments with values from state file
#AWSKS | initialize-state-file | will initialize state file
#AWSKS | display-config-file | config file content is:
kind: awsks-config
awsks:
  name: epiphany
  vpc_id: vpc-0baa2c4e9e48e608c
  region: eu-central-1
  subnet_ids: null
  private_route_table_id: unset
`,
			wantConfigLocation: "awsks/awsks-config.yml",
			wantConfigContent: `
kind: awsks-config
awsks:
  name: epiphany
  vpc_id: vpc-0baa2c4e9e48e608c
  region: eu-central-1
  subnet_ids: null
  private_route_table_id: unset
`,
			wantStateContent: `
kind: state
awsbi:
  status: applied
  name: epiphany
  instance_count: 0
  region: eu-central-1
  use_public_ip: false
  force_nat_gateway: true
  rsa_pub_path: "/shared/vms_rsa.pub"
  output:
    private_ip.value: []
    public_ip.value: []
    public_subnet_id.value: subnet-0137cf1e7921c1551
    vpc_id.value: vpc-0baa2c4e9e48e608c
awsks:
  status: initialized
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sharedPath := setupOutput(t, "init")
			defer cleanupOutput(sharedPath)

			stateLocation := path.Join(sharedPath, tt.stateLocation)
			if err := ioutil.WriteFile(stateLocation, []byte(normStr(tt.stateContent)), 0644); err != nil {
				t.Fatalf("wasnt able to save state file: %s", err)
			}

			command := []string{"init"}
			command = append(command, tt.initParams...)

			runOpts := &docker.RunOptions{
				Command: command,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			output := docker.Run(t, awsksImageTag, runOpts)
			if diff := deep.Equal(normStr(output), normStr(tt.wantOutput)); diff != nil {
				t.Error(diff)
			}

			configLocation := path.Join(sharedPath, tt.wantConfigLocation)
			if _, err := os.Stat(configLocation); os.IsNotExist(err) {
				t.Fatalf("missing expected file: %s", configLocation)
			}

			gotConfigContent, err := ioutil.ReadFile(configLocation)
			if err != nil {
				t.Errorf("wasnt able to read form output file: %v", err)
			}

			if diff := deep.Equal(normStr(string(gotConfigContent)), normStr(tt.wantConfigContent)); diff != nil {
				t.Error(diff)
			}

			gotStateContent, err := ioutil.ReadFile(stateLocation)
			if err != nil {
				t.Errorf("wasnt able to read form state file: %v", err)
			}

			if diff := deep.Equal(normStr(string(gotStateContent)), normStr(tt.wantStateContent)); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestPlan(t *testing.T) {
	awsAccessKey, awsSecretKey := getAwsCreds(t)
	sharedPath := setupOutput(t, "plan")
	setupPlan(t, "plan", sharedPath, awsAccessKey, awsSecretKey)

	tests := []struct{
		name                   string
		initParams             []string
		wantPlanOutputLastLine string
		wantTfPlanLocation     string
	}{
		{
			name: "plan",
			initParams: []string{"M_NAME=eks-module-tests-plan"},
			wantPlanOutputLastLine: `Plan: 29 to add, 0 to change, 0 to destroy.`,
			wantTfPlanLocation: "awsks/terraform-apply.tfplan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initCommand := []string{"init"}
			initCommand = append(initCommand, tt.initParams...)

			initOpts := &docker.RunOptions{
				Command: initCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			docker.Run(t, awsksImageTag, initOpts)

			planCommand := []string{"plan",
				fmt.Sprintf("M_AWS_ACCESS_KEY=%s", awsAccessKey),
				fmt.Sprintf("M_AWS_SECRET_KEY=%s", awsSecretKey),
			}

			planOpts := &docker.RunOptions{
				Command: planCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			gotPlanOutput := docker.Run(t, awsksImageTag, planOpts)
			gotPlanOutputLastLine, err := getLastLineFromMultilineString(gotPlanOutput)
			if err != nil {
				t.Fatalf("reading last line from multiline failed with: %v", err)
			}

			if diff := deep.Equal(gotPlanOutputLastLine, tt.wantPlanOutputLastLine); diff != nil {
				t.Error(diff)
			}

			tfPlanLocation := path.Join(sharedPath, tt.wantTfPlanLocation)
			if _, err := os.Stat(tfPlanLocation); os.IsNotExist(err) {
				t.Fatalf("missing tfplan file: %s", tfPlanLocation)
			}
		})
	}

	cleanupPlan(t, sharedPath, awsAccessKey, awsSecretKey)
	cleanupOutput(sharedPath)
}

func TestApply(t *testing.T) {
	awsAccessKey, awsSecretKey := getAwsCreds(t)
	sharedPath := setupOutput(t, "apply")
	setupPlan(t, "apply", sharedPath, awsAccessKey, awsSecretKey)

	tests := []struct{
		name       string
		initParams []string
	}{
		{
			name: "apply",
			initParams: []string{"M_NAME=eks-module-tests-apply"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initCommand := []string{"init"}
			initCommand = append(initCommand, tt.initParams...)

			initOpts := &docker.RunOptions{
				Command: initCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			docker.Run(t, awsksImageTag, initOpts)

			planCommand := []string{"plan",
				fmt.Sprintf("M_AWS_ACCESS_KEY=%s", awsAccessKey),
				fmt.Sprintf("M_AWS_SECRET_KEY=%s", awsSecretKey),
			}

			planOpts := &docker.RunOptions{
				Command: planCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			docker.Run(t, awsksImageTag, planOpts)

			applyCommand := []string{"apply",
				fmt.Sprintf("M_AWS_ACCESS_KEY=%s", awsAccessKey),
				fmt.Sprintf("M_AWS_SECRET_KEY=%s", awsSecretKey),
			}

			applyOpts := &docker.RunOptions{
				Command: applyCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			docker.Run(t, awsksImageTag, applyOpts)

			kubeconfigCommand := []string{"kubeconfig"}

			kubeconfigOpts := &docker.RunOptions{
				Command: kubeconfigCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			docker.Run(t, awsksImageTag, kubeconfigOpts)

			kubectlOpts := &k8s.KubectlOptions{
				ConfigPath: fmt.Sprintf("%s/kubeconfig", sharedPath),
			}

			k8s.RunKubectl(t, kubectlOpts, "get", "all", "-A")

			planDestroyCommand := []string{"plan-destroy",
				fmt.Sprintf("M_AWS_ACCESS_KEY=%s", awsAccessKey),
				fmt.Sprintf("M_AWS_SECRET_KEY=%s", awsSecretKey),
			}

			planDestroyOpts := &docker.RunOptions{
				Command: planDestroyCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			docker.Run(t, awsksImageTag, planDestroyOpts)

			destroyCommand := []string{"destroy",
				fmt.Sprintf("M_AWS_ACCESS_KEY=%s", awsAccessKey),
				fmt.Sprintf("M_AWS_SECRET_KEY=%s", awsSecretKey),
			}

			destroyOpts := &docker.RunOptions{
				Command: destroyCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			docker.Run(t, awsksImageTag, destroyOpts)
		})
	}

	cleanupPlan(t, sharedPath, awsAccessKey, awsSecretKey)
	cleanupOutput(sharedPath)
}

func setupPlan(t *testing.T, suffix, sharedPath, awsAccessKey, awsSecretKey string) {
	if err := generateRsaKeyPair(sharedPath, "test_vms_rsa"); err != nil {
		t.Fatalf("wasnt able to create rsa file: %s", err)
	}

	initCommand := []string{
		"init",
		"M_VMS_COUNT=0",
		"M_PUBLIC_IPS=false",
		fmt.Sprintf("M_NAME=eks-module-tests-%s", suffix),
		"M_VMS_RSA=test_vms_rsa",
	}

	initOpts := &docker.RunOptions{
		Command: initCommand,
		Remove:  true,
		Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
	}

	docker.Run(t, awsbiImageTag, initOpts)

	planCommand := []string{"plan",
		fmt.Sprintf("M_AWS_ACCESS_KEY=%s", awsAccessKey),
		fmt.Sprintf("M_AWS_SECRET_KEY=%s", awsSecretKey),
	}

	planOpts := &docker.RunOptions{
		Command: planCommand,
		Remove:  true,
		Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
	}

	docker.Run(t, awsbiImageTag, planOpts)

	applyCommand := []string{"apply",
		fmt.Sprintf("M_AWS_ACCESS_KEY=%s", awsAccessKey),
		fmt.Sprintf("M_AWS_SECRET_KEY=%s", awsSecretKey),
	}

	applyOpts := &docker.RunOptions{
		Command: applyCommand,
		Remove:  true,
		Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
	}

	docker.Run(t, awsbiImageTag, applyOpts)
}

func cleanupPlan(t *testing.T, sharedPath, awsAccessKey, awsSecretKey string) {
	planDestroyCommand := []string{"plan-destroy",
		fmt.Sprintf("M_AWS_ACCESS_KEY=%s", awsAccessKey),
		fmt.Sprintf("M_AWS_SECRET_KEY=%s", awsSecretKey),
	}

	planDestroyOpts := &docker.RunOptions{
		Command: planDestroyCommand,
		Remove:  true,
		Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
	}

	docker.Run(t, awsbiImageTag, planDestroyOpts)

	destroyCommand := []string{"destroy",
		fmt.Sprintf("M_AWS_ACCESS_KEY=%s", awsAccessKey),
		fmt.Sprintf("M_AWS_SECRET_KEY=%s", awsSecretKey),
	}

	destroyOpts := &docker.RunOptions{
		Command: destroyCommand,
		Remove:  true,
		Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
	}

	docker.Run(t, awsbiImageTag, destroyOpts)
}

func setupOutput(t *testing.T, suffix string) (string) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("setupOutput() failed with: %v", err)
	}
	p := path.Join(wd, fmt.Sprintf("%s-%s", "shared", suffix))
	cleanupOutput(p)
	err = os.MkdirAll(p, os.ModePerm)
	if err != nil {
		t.Fatalf("setupOutput() failed with: %v", err)
	}
	return p
}

func cleanupOutput(sharedPath string) error {
	return os.RemoveAll(sharedPath)
}

func normStr(s string) string {
	return strings.TrimSpace(s)
}

func getLastLineFromMultilineString(s string) (string, error) {
	in := strings.NewReader(s)
	reader := bufio.NewReader(in)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		if err == io.EOF {
			return string(line), nil
		}
	}
}

func generateRsaKeyPair(directory, name string) error {
	privateRsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateRsaKey)}
	privateKeyBytes := pem.EncodeToMemory(pemBlock)

	publicRsaKey, err := ssh.NewPublicKey(&privateRsaKey.PublicKey)
	if err != nil {
		return err
	}
	publicKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	err = ioutil.WriteFile(path.Join(directory, name), privateKeyBytes, 0600)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(directory, fmt.Sprintf("%s.pub", name)), publicKeyBytes, 0644)
}

func getAwsCreds(t *testing.T) (awsAccessKey, awsSecretKey string) {
	awsAccessKey = os.Getenv("AWS_ACCESS_KEY")
	if len(awsAccessKey) == 0 {
		t.Fatalf("expected non-empty AWS_ACCESS_KEY environment variable")
	}

	awsSecretKey = os.Getenv("AWS_SECRET_KEY")
	if len(awsSecretKey) == 0 {
		t.Fatalf("expected non-empty AWS_SECRET_KEY environment variable")
	}

	return
}





func cleanupAWSResources(t *testing.T) {
	newSession, errSession := session.NewSession(&aws.Config{Region: aws.String(awsRegion)})
	if errSession != nil {
		t.Fatalf("Cannot get session: %s", errSession)
	}

	rgClient := resourcegroups.New(newSession)

	rgName := moduleName + "-rg"
	kpName := moduleName + "-kp"

	rgResourcesList, errResourcesList := rgClient.ListGroupResources(&resourcegroups.ListGroupResourcesInput{
		GroupName: aws.String(rgName),
	})

	if errResourcesList != nil {
		if aerr, ok := errResourcesList.(awserr.Error); ok {
			t.Log(aerr.Code())
			if aerr.Code() == "NotFoundException" {
				t.Log("Resource group: ", rgName, " not found.")
			} else {
				t.Log("Resource group: Cannot get list of resources: ", errResourcesList)
			}
		} else {
			t.Log("Resource group: There was an error: ", errResourcesList.Error())
		}
	}

	resourcesTypesToRemove := []string{"Instance", "SecurityGroup", "NatGateway", "EIP", "InternetGateway", "Subnet", "RouteTable", "VPC"}

	for _, resourcesTypeToRemove := range resourcesTypesToRemove {

		filtered := make([]*resourcegroups.ResourceIdentifier, 0)
		for _, element := range rgResourcesList.ResourceIdentifiers {
			s := strings.Split(*element.ResourceType, ":")
			if s[4] == resourcesTypeToRemove {
				filtered = append(filtered, element)
			}

		}

		switch resourcesTypeToRemove {
		case "Instance":
			t.Log("Instance.")
			removeEc2s(t, newSession, filtered)	
		case "EIP":
			t.Log("Releasing public EIPs.")
			releaseAddresses(t, newSession)
		case "RouteTable":
			t.Log("RouteTable.")
			removeRouteTables(t, newSession, filtered)
		case "InternetGateway":
			t.Log("InternetGateway.")
			removeInternetGateway(t, newSession, filtered)
		case "SecurityGroup":
			t.Log("SecurityGroup.")
			removeSecurityGroup(t, newSession, filtered)
		case "NatGateway":
			t.Log("NatGateway.")
			removeNatGateways(t, newSession, filtered)
		case "Subnet":
			t.Log("Subnet.")
			removeSubnet(t, newSession, filtered)
		case "VPC":
			t.Log("VPC.")
			removeVpc(t, newSession, filtered)
		}
	}

	removeResourceGroup(t, newSession, rgName)
	removeKeyPair(t, newSession, kpName)
}

func removeEc2s(t *testing.T, session *session.Session, ec2sToRemove []*resourcegroups.ResourceIdentifier) {
	ec2Client := ec2.New(session)

	for _, ec2ToRemove := range ec2sToRemove {

		ec2ToRemoveID := strings.Split(*ec2ToRemove.ResourceArn, "/")[1]
		t.Log("EC2: Removing instance with ID: ", ec2ToRemoveID)

		ec2DescInp := &ec2.DescribeInstancesInput{
			InstanceIds: []*string{&ec2ToRemoveID},
		}

		outDesc, errDesc := ec2Client.DescribeInstances(ec2DescInp)
		if errDesc != nil {
			t.Fatalf("EC2: Describe error: %s", errDesc)
		}
		t.Log("EC2: Describe output: ", outDesc)

		if outDesc.Reservations != nil {

			instancesToTerminateInp := &ec2.TerminateInstancesInput{
				InstanceIds: []*string{&ec2ToRemoveID},
			}

			outputTerm, errTerm := ec2Client.TerminateInstances(instancesToTerminateInp)
			if errTerm != nil {
				t.Fatalf("EC2: Terminate error: %s", outputTerm)
			}
			t.Log("EC2: Terminate output: ", outputTerm)

			errWait := ec2Client.WaitUntilInstanceTerminated(ec2DescInp)
			if errWait != nil {
				t.Fatalf("EC2: Waiting for termination error: %s", errWait)
			}
		}

	}
}

func removeRouteTables(t *testing.T, session *session.Session, rtsToRemove []*resourcegroups.ResourceIdentifier) {
	ec2Client := ec2.New(session)

	for _, rtToRemove := range rtsToRemove {
		rtIDToRemove := strings.Split(*rtToRemove.ResourceArn, "/")[1]
		t.Log("RouteTable: rtIDToRemove: ", rtIDToRemove)

		rtToDeleteInp := &ec2.DeleteRouteTableInput{
			RouteTableId: &rtIDToRemove,
		}

		output, err := ec2Client.DeleteRouteTable(rtToDeleteInp)

		if err != nil {
			t.Fatalf("RouteTable: Deleting route table error: %s", err)
		}

		t.Log("RouteTable: Deleting route table: ", output)
	}
}

func removeSecurityGroup(t *testing.T, session *session.Session, sgsToRemove []*resourcegroups.ResourceIdentifier) {
	ec2Client := ec2.New(session)

	for _, sgToRemove := range sgsToRemove {
		sgIDToRemove := strings.Split(*sgToRemove.ResourceArn, "/")[1]
		t.Log("Security Group: sgIdToRemove: ", sgIDToRemove)

		secGrpInp := &ec2.DeleteSecurityGroupInput{GroupId: &sgIDToRemove}

		output, err := ec2Client.DeleteSecurityGroup(secGrpInp)
		if err != nil {
			t.Fatalf("Security Group: Deleting security group error: %s", err)
		}

		t.Log("Security Group: Deleting security group: ", output)
	}
}

func removeInternetGateway(t *testing.T, session *session.Session, igsToRemove []*resourcegroups.ResourceIdentifier) {
	ec2Client := ec2.New(session)

	for _, igToRemove := range igsToRemove {
		igIDToRemove := strings.Split(*igToRemove.ResourceArn, "/")[1]
		t.Log("Internet Gateway: igIdToRemove: ", igIDToRemove)

		igDescribeInp := &ec2.DescribeInternetGatewaysInput{
			InternetGatewayIds: []*string{&igIDToRemove},
		}

		descOut, descErr := ec2Client.DescribeInternetGateways(igDescribeInp)

		if descErr != nil {
			t.Fatalf("Internet Gateway: Describing internet gateway error: %s", descErr)
		}
		t.Log("Internet Gateway: Describing internet gateway: ", descOut)
		vpcID := *descOut.InternetGateways[0].Attachments[0].VpcId

		igDetachInp := &ec2.DetachInternetGatewayInput{
			InternetGatewayId: &igIDToRemove,
			VpcId:             &vpcID,
		}

		detachOut, detachErr := ec2Client.DetachInternetGateway(igDetachInp)
		if detachErr != nil {
			t.Fatalf("Internet Gateway: Detaching internet gateway error: %s", detachErr)
		}
		t.Log("Internet Gateway: Detaching internet gateway: ", detachOut)

		igDeleteInp := &ec2.DeleteInternetGatewayInput{
			InternetGatewayId: &igIDToRemove,
		}

		delOut, delErr := ec2Client.DeleteInternetGateway(igDeleteInp)
		if delErr != nil {
			t.Fatalf("Internet Gateway: Deleting internet gateway error: %s", delErr)
		}
		t.Log("Internet Gateway: Deleting internet gateway: ", delOut)
	}
}

func removeNatGateways(t *testing.T, session *session.Session, ngsToRemove []*resourcegroups.ResourceIdentifier) {
	ec2Client := ec2.New(session)
	for _, ngToRemove := range ngsToRemove {
		ngIDToRemove := strings.Split(*ngToRemove.ResourceArn, "/")[1]
		t.Log("Nat Gateway: ngIdToRemove: ", ngIDToRemove)
		removeSingleNatGatewayWithRetries(t, ec2Client, ngIDToRemove)
	}
}

func removeSingleNatGatewayWithRetries(t *testing.T, ec2Client *ec2.EC2, ngIDToRemove string) {
	found := true
	for retry := 0; retry <= retries && found; retry++ {
		found = describeNatGateway(t, ec2Client, ngIDToRemove)

		if found == false {
			continue
		}

		found = removeNatGateway(t, ec2Client, ngIDToRemove)

		if found == false {
			continue
		}

		waitForNatGatewayDelete(t, ec2Client, ngIDToRemove)

		t.Log("Nat Gateway: Deleting NAT Gateway. ", ngIDToRemove, " Retry: ", retry)
		time.Sleep(5 * time.Second)
	}
}

func describeNatGateway(t *testing.T, ec2Client *ec2.EC2, ngIDToDescribe string) bool {
	descInp := &ec2.DescribeNatGatewaysInput{
		NatGatewayIds: []*string{&ngIDToDescribe},
	}

	outDesc, errDesc := ec2Client.DescribeNatGateways(descInp)
	if errDesc != nil {
		t.Log(errDesc)
		if aerr, ok := errDesc.(awserr.Error); ok {
			if aerr.Code() == "NatGatewayNotFound" {
				t.Log("Nat Gateway: Nat Gateway not found.")
				return false
			} else {
				t.Fatalf("Nat Gateway: Describe error: %s", errDesc)
			}
		} else {
			t.Fatalf("Nat Gateway: There was an error: %s", errDesc.Error())
		}
	}
	t.Log("Nat Gateway: Describe output: ", outDesc)

	if outDesc.NatGateways == nil || *outDesc.NatGateways[0].State == "deleted" {
		t.Log("Nat Gateway: Element not found or has been already deleted.")
		return false
	}
	return true
}

func removeNatGateway(t *testing.T, ec2Client *ec2.EC2, ngIDToRemove string) bool {
	ngDelInp := &ec2.DeleteNatGatewayInput{
		NatGatewayId: &ngIDToRemove,
	}

	_, err := ec2Client.DeleteNatGateway(ngDelInp)

	if err != nil {
		t.Log("Nat Gateway: Error: ", err)
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "NatGatewayNotFound" {
				t.Log("Nat Gateway: Element not found.", err)
				return false
			}
			if aerr.Code() != "ResourceNotReady" {
				t.Fatalf("Nat Gateway: Deleting NAT Gateway: %s", err)
			}
		} else {
			t.Fatalf("Nat Gateway: Deleting NAT Gateway: %s", err.Error())
		}

	}
	return true
}

func waitForNatGatewayDelete(t *testing.T, ec2Client *ec2.EC2, ngIDToWait string) {
	descInp := &ec2.DescribeNatGatewaysInput{
		NatGatewayIds: []*string{&ngIDToWait},
	}
	errWait := ec2Client.WaitUntilNatGatewayAvailable(descInp)
	if errWait != nil {
		if aerr, ok := errWait.(awserr.Error); ok {
			if aerr.Code() != "ResourceNotReady" {
				t.Fatalf("Nat Gateway: Wait error: %s", errWait)
			}
		} else {
			t.Fatalf("Nat Gateway: There was an error: %s", errWait.Error())
		}
	}
}

func removeSubnet(t *testing.T, session *session.Session, subnetsToRemove []*resourcegroups.ResourceIdentifier) {
	ec2Client := ec2.New(session)
	for _, subnetToRemove := range subnetsToRemove {
		subnetIDToRemove := strings.Split(*subnetToRemove.ResourceArn, "/")[1]
		t.Log("Subnet: subnetIdToRemove: ", subnetIDToRemove)

		subnetInp := &ec2.DeleteSubnetInput{
			SubnetId: &subnetIDToRemove,
		}

		output, err := ec2Client.DeleteSubnet(subnetInp)
		if err != nil {
			t.Fatalf("Subnet: Deleting subnet error: %s", err)
		}
		t.Log("Subnet: Deleting subnet: ", output)
	}
}

func removeVpc(t *testing.T, session *session.Session, vpcsToRemove []*resourcegroups.ResourceIdentifier) {
	ec2Client := ec2.New(session)
	for _, vpcToRemove := range vpcsToRemove {
		vpcIDToRemove := strings.Split(*vpcToRemove.ResourceArn, "/")[1]
		t.Log("VPC: vpcIdToRemove: ", vpcIDToRemove)

		vpcToDeleteInp := &ec2.DeleteVpcInput{
			VpcId: &vpcIDToRemove,
		}

		output, err := ec2Client.DeleteVpc(vpcToDeleteInp)
		if err != nil {
			t.Log("VPC: Delete VPC error: ", err)
		}
		t.Log("VPC: Delete VPC: ", output)
	}
}

func removeKeyPair(t *testing.T, session *session.Session, kpName string) {
	ec2Client := ec2.New(session)

	removeKeyInp := &ec2.DeleteKeyPairInput{
		KeyName: &kpName,
	}

	output, err := ec2Client.DeleteKeyPair(removeKeyInp)
	if err != nil {
		t.Fatalf("Key Pair: Deleting key pair error: %s", err)
	}
	t.Log("Key Pair: Deleting key pair: ", output)
}

func releaseAddresses(t *testing.T, session *session.Session) {
    ec2Client := ec2.New(session)

    eipDescInp := &ec2.DescribeAddressesInput {
        Filters: []*ec2.Filter{
            {
                Name: aws.String("tag:" + awsTagName),
                Values: []*string{
                    aws.String(awsTagValue),
                },
            },
        },
    }

    describeEips, err := ec2Client.DescribeAddresses(eipDescInp)
    if err != nil {
        t.Fatalf("EIP: Cannot get EIP list: %s", err)
    }

    for _, eip := range describeEips.Addresses {

        t.Log("EIP: Releasing EIP with AllocationId: ", *eip.AllocationId)

        eipToReleaseInp := &ec2.ReleaseAddressInput{
            AllocationId: eip.AllocationId,
        }

        found := true
        for retry := 0; retry <= retries && found; retry++ {
            _, err := ec2Client.ReleaseAddress(eipToReleaseInp)
            if err != nil {
                if aerr, ok := err.(awserr.Error); ok {
                    if aerr.Code() == "InvalidAllocationID.NotFound" {
                        t.Log("EIP: Element not found.", err)
                        found = false
                        continue
                    }
                    if aerr.Code() != "AuthFailure" && aerr.Code() != "InvalidAllocationID.NotFound" {
                        t.Fatalf("EIP: Releasing EIP error: %s", err)
                    }
                } else {
                    t.Fatalf("EIP: There was an error: %s", err.Error())
                }
            }
            t.Log("EIP: Releasing EIP. Retry: ", retry)
            time.Sleep(5 * time.Second)
        }
    }
}

func removeResourceGroup(t *testing.T, session *session.Session, rgToRemoveName string) {
	rgClient := resourcegroups.New(session)

	t.Log("Resource Group: Removing resource group: ", rgToRemoveName)
	rgDelInp := resourcegroups.DeleteGroupInput{
		GroupName: aws.String(rgToRemoveName),
	}
	rgDelOut, rgDelErr := rgClient.DeleteGroup(&rgDelInp)
	if rgDelErr != nil {
		if aerr, ok := rgDelErr.(awserr.Error); ok {
			if aerr.Code() == "NotFoundException" {
				t.Log("Resource Group: Resource group not found. ")
			} else {
				t.Fatalf("Resource Group: Deleting resource group error: %s", rgDelErr)
			}
		} else {
			t.Fatalf("Resource Group: There was an error: %s", rgDelErr.Error())
		}

	} else {
		t.Log("Resource Group: Deleting resource group: ", rgDelOut)
	}
}
