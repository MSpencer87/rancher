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
	"github.com/rancher/rancher/tests/framework/pkg/environmentflag"
	namegen "github.com/rancher/rancher/tests/framework/pkg/namegenerator"
	"github.com/rancher/rancher/tests/framework/pkg/session"
	"github.com/rancher/rancher/tests/v2/validation/pipeline/rancherha/corralha"
	"github.com/rancher/rancher/tests/v2/validation/provisioning/permutations"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type K3SProxyTestSuite struct {
	suite.Suite
	client             *rancher.Client
	standardUserClient *rancher.Client
	session            *session.Session
	clustersConfig     *provisioninginput.Config
	EnvVar             rkev1.EnvVar
	corralImage        string
	corralAutoCleanup  bool
}

func (k *K3SProxyTestSuite) TearDownSuite() {
	k.session.Cleanup()
}

func (k *K3SProxyTestSuite) SetupSuite() {
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

func (k *K3SProxyTestSuite) TestProvisioningK3SClusterProxy() {
	nodeRolesAll := []provisioninginput.MachinePools{provisioninginput.AllRolesMachinePool}
	tests := []struct {
		name         string
		machinePools []provisioninginput.MachinePools
		client       *rancher.Client
		runFlag      bool
	}{
		{"1 Node all roles " + provisioninginput.AdminClientName.String(), nodeRolesAll, k.client, k.client.Flags.GetValue(environmentflag.Short)},
	}

	for _, tt := range tests {
		if !tt.runFlag {
			k.T().Logf("SKIPPED")
			continue
		}
		provisioningConfig := *k.clustersConfig
		provisioningConfig.MachinePools = tt.machinePools
		permutations.RunTestPermutations(&k.Suite, tt.name, tt.client, &provisioningConfig, permutations.K3SProvisionCluster, nil, nil)
	}
}

func TestK3SProxyTestSuite(t *testing.T) {
	suite.Run(t, new(K3SProxyTestSuite))
}
