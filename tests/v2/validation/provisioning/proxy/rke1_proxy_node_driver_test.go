package proxy

import (
	"testing"

	"github.com/rancher/rancher/tests/framework/clients/corral"
	"github.com/rancher/rancher/tests/framework/clients/rancher"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/clusters/kubernetesversions"
	"github.com/rancher/rancher/tests/framework/extensions/provisioninginput"
	"github.com/rancher/rancher/tests/framework/extensions/users"
	password "github.com/rancher/rancher/tests/framework/extensions/users/passwordgenerator"
	"github.com/rancher/rancher/tests/framework/pkg/config"
	"github.com/rancher/rancher/tests/framework/pkg/environmentflag"
	namegen "github.com/rancher/rancher/tests/framework/pkg/namegenerator"
	"github.com/rancher/rancher/tests/framework/pkg/session"
	"github.com/rancher/rancher/tests/v2/validation/pipeline/rancherha/corralha"
	"github.com/rancher/rancher/tests/v2/validation/provisioning/permutations"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RKE1ProxyProvisioningTestSuite struct {
	suite.Suite
	client             *rancher.Client
	standardUserClient *rancher.Client
	session            *session.Session
	provisioningConfig *provisioninginput.Config
	EnvVar             management.EnvVar
	corralImage        string
	corralAutoCleanup  bool
}

func (r *RKE1ProxyProvisioningTestSuite) TearDownSuite() {
	r.session.Cleanup()
}

func (r *RKE1ProxyProvisioningTestSuite) SetupSuite() {
	testSession := session.NewSession()
	r.session = testSession

	corralRancherHA := new(corralha.CorralRancherHA)
	config.LoadConfig(corralha.CorralRancherHAConfigConfigurationFileKey, corralRancherHA)

	r.provisioningConfig = new(provisioninginput.Config)
	config.LoadConfig(provisioninginput.ConfigurationFileKey, r.provisioningConfig)

	client, err := rancher.NewClient("", testSession)
	require.NoError(r.T(), err)

	r.client = client

	r.provisioningConfig.RKE1KubernetesVersions, err = kubernetesversions.Default(
		r.client, clusters.RKE1ClusterType.String(), r.provisioningConfig.RKE1KubernetesVersions)
	require.NoError(r.T(), err)

	enabled := true
	var testuser = namegen.AppendRandomString("testuser-")
	var testpassword = password.GenerateUserPassword("testpass-")
	user := &management.User{
		Username: testuser,
		Password: testpassword,
		Name:     testuser,
		Enabled:  &enabled,
	}

	newUser, err := users.CreateUserWithRole(client, user, "user")
	require.NoError(r.T(), err)

	newUser.Password = user.Password

	standardUserClient, err := client.AsUser(newUser)
	require.NoError(r.T(), err)

	r.standardUserClient = standardUserClient

	listOfCorrals, err := corral.ListCorral()
	require.NoError(r.T(), err)

	corralConfig := corral.Configurations()
	err = corral.SetupCorralConfig(corralConfig.CorralConfigVars, corralConfig.CorralConfigUser, corralConfig.CorralSSHPath)
	require.NoError(r.T(), err)

	corralPackage := corral.PackagesConfig()
	r.corralImage = corralPackage.CorralPackageImages[corralPackageProxyCustomClusterName]
	r.corralAutoCleanup = corralPackage.HasCleanup

	_, corralExist := listOfCorrals[corralRancherHA.Name]
	if corralExist {
		bastionIP, err := corral.GetCorralEnvVar(corralRancherHA.Name, corralRegistryIP)
		require.NoError(r.T(), err)

		err = corral.UpdateCorralConfig(corralBastionIP, bastionIP)
		require.NoError(r.T(), err)

		r.EnvVar.Name = "HTTP_PROXY"
		r.EnvVar.Value = bastionIP + ":3219"
		r.provisioningConfig.AgentEnvVarsRKE1 = append(r.provisioningConfig.AgentEnvVarsRKE1, r.EnvVar)

		r.EnvVar.Name = "HTTPS_PROXY"
		r.EnvVar.Value = bastionIP + ":3219"
		r.provisioningConfig.AgentEnvVarsRKE1 = append(r.provisioningConfig.AgentEnvVarsRKE1, r.EnvVar)

		r.EnvVar.Name = "NO_PROXY"
		r.EnvVar.Value = "localhost,127.0.0.1,0.0.0.0,10.0.0.0/8,cattle-system.svc"
		r.provisioningConfig.AgentEnvVarsRKE1 = append(r.provisioningConfig.AgentEnvVarsRKE1, r.EnvVar)
	}
}

func (r *RKE1ProxyProvisioningTestSuite) TestProvisioningRKE1ClusterProxy() {
	nodeRolesAll := []provisioninginput.NodePools{provisioninginput.AllRolesNodePool}

	tests := []struct {
		name      string
		nodePools []provisioninginput.NodePools
		client    *rancher.Client
		runFlag   bool
	}{
		//{"5 Node all roles " + provisioninginput.AdminClientName.String(), nodeRoles0, r.client},
		{"3 nodes - 1 role per node " + provisioninginput.AdminClientName.String(), nodeRolesAll, r.client, r.client.Flags.GetValue(environmentflag.Short)},
	}
	for _, tt := range tests {
		if !tt.runFlag {
			r.T().Logf("SKIPPED")
			continue
		}
		provisioningConfig := *r.provisioningConfig
		provisioningConfig.NodePools = tt.nodePools
		permutations.RunTestPermutations(&r.Suite, tt.name, tt.client, &provisioningConfig, permutations.RKE1ProvisionCluster, nil, nil)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRKE1ProxyProvisioningTestSuite(t *testing.T) {
	suite.Run(t, new(RKE1ProxyProvisioningTestSuite))
}
