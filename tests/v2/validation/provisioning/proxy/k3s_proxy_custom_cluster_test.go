package proxy

import (
	"testing"

	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/rancher/rancher/tests/framework/clients/corral"
	"github.com/rancher/rancher/tests/framework/clients/rancher"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/clusters/kubernetesversions"
	"github.com/rancher/rancher/tests/framework/extensions/provisioninginput"
	"github.com/rancher/rancher/tests/framework/extensions/users"
	password "github.com/rancher/rancher/tests/framework/extensions/users/passwordgenerator"
	"github.com/rancher/rancher/tests/framework/pkg/config"
	namegen "github.com/rancher/rancher/tests/framework/pkg/namegenerator"
	"github.com/rancher/rancher/tests/framework/pkg/session"
	"github.com/rancher/rancher/tests/v2/validation/pipeline/rancherha/corralha"
	"github.com/rancher/rancher/tests/v2/validation/provisioning/permutations"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ProxyK3SCustomClusterTestSuite struct {
	suite.Suite
	client             *rancher.Client
	standardUserClient *rancher.Client
	session            *session.Session
	corralPackage      *corral.Packages
	clustersConfig     *provisioninginput.Config
	EnvVar             rkev1.EnvVar
	corralImage        string
	corralAutoCleanup  bool
}

func (k *ProxyK3SCustomClusterTestSuite) TearDownSuite() {
	k.session.Cleanup()
}

func (k *ProxyK3SCustomClusterTestSuite) SetupSuite() {
	testSession := session.NewSession()
	k.session = testSession

	corralRancherHA := new(corralha.CorralRancherHA)
	config.LoadConfig(corralha.CorralRancherHAConfigConfigurationFileKey, corralRancherHA)

	k.clustersConfig = new(provisioninginput.Config)
	config.LoadConfig(provisioninginput.ConfigurationFileKey, k.clustersConfig)

	client, err := rancher.NewClient("", testSession)
	require.NoError(k.T(), err)

	k.client = client

	k.clustersConfig.K3SKubernetesVersions, err = kubernetesversions.Default(
		k.client, clusters.K3SClusterType.String(), k.clustersConfig.K3SKubernetesVersions)
	require.NoError(k.T(), err)

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
	require.NoError(k.T(), err)

	newUser.Password = user.Password

	standardUserClient, err := client.AsUser(newUser)
	require.NoError(k.T(), err)

	k.standardUserClient = standardUserClient

	listOfCorrals, err := corral.ListCorral()
	require.NoError(k.T(), err)

	corralConfig := corral.Configurations()
	err = corral.SetupCorralConfig(corralConfig.CorralConfigVars, corralConfig.CorralConfigUser, corralConfig.CorralSSHPath)
	require.NoError(k.T(), err)

	corralPackage := corral.PackagesConfig()
	k.corralImage = corralPackage.CorralPackageImages[corralPackageAirgapCustomClusterName]
	k.corralAutoCleanup = corralPackage.HasCleanup

	_, corralExist := listOfCorrals[corralRancherHA.Name]
	if corralExist {
		bastionIP, err := corral.GetCorralEnvVar(corralRancherHA.Name, corralRegistryIP)
		require.NoError(k.T(), err)

		err = corral.UpdateCorralConfig(corralBastionIP, bastionIP)
		require.NoError(k.T(), err)

		k.EnvVar.Name = "HTTP_PROXY"
		k.EnvVar.Value = bastionIP + ":3219"
		k.clustersConfig.AgentEnvVars = append(k.clustersConfig.AgentEnvVars, k.EnvVar)

		k.EnvVar.Name = "HTTPS_PROXY"
		k.EnvVar.Value = bastionIP + ":3219"
		k.clustersConfig.AgentEnvVars = append(k.clustersConfig.AgentEnvVars, k.EnvVar)

		k.EnvVar.Name = "NO_PROXY"
		k.EnvVar.Value = "localhost,127.0.0.1,0.0.0.0,10.0.0.0/8,cattle-system.svc"
		k.clustersConfig.AgentEnvVars = append(k.clustersConfig.AgentEnvVars, k.EnvVar)
	}

}

func (k *ProxyK3SCustomClusterTestSuite) TestProvisioningK3SCustomClusterProxy() {
	k.clustersConfig.MachinePools = []provisioninginput.MachinePools{provisioninginput.AllRolesMachinePool}

	tests := []struct {
		name   string
		client *rancher.Client
	}{
		{provisioninginput.AdminClientName.String() + "-" + permutations.K3SAirgapCluster + "-", k.client},
	}
	for _, tt := range tests {
		permutations.RunTestPermutations(&k.Suite, tt.name, tt.client, k.clustersConfig, permutations.K3SAirgapCluster, nil, k.corralPackage)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestProxyCustomClusterK3SProvisioningTestSuite(t *testing.T) {
	suite.Run(t, new(ProxyK3SCustomClusterTestSuite))
}
