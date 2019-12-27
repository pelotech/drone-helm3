package run

import (
	"github.com/stretchr/testify/suite"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"testing"
)

type KubeconfigTestSuite struct {
	suite.Suite
	configFile *os.File
	initKube   InitKube
}

func (suite *KubeconfigTestSuite) BeforeTest(_, _ string) {
	file, err := ioutil.TempFile("", "kubeconfig********.yml")
	suite.Require().NoError(err)
	file.Close()
	suite.configFile = file

	// set up an InitKube with the bare minimum configuration
	suite.initKube = InitKube{
		ConfigFile:   file.Name(),
		TemplateFile: "../../assets/kubeconfig.tpl", // the actual kubeconfig template
		APIServer:    "a",
		Token:        "b",
	}
}

func (suite *KubeconfigTestSuite) AfterTest(_, _ string) {
	if suite.configFile != nil {
		os.Remove(suite.configFile.Name())
	}
}

func TestKubeconfigTestSuite(t *testing.T) {
	suite.Run(t, new(KubeconfigTestSuite))
}

func (suite *KubeconfigTestSuite) TestSetsNamespace() {
	cfg := Config{
		Namespace: "marshmallow",
	}
	contents := suite.generateKubeconfig(cfg)
	suite.Contains(contents, "namespace: marshmallow")
}

func (suite *KubeconfigTestSuite) TestSetsAPIServer() {
	suite.initKube.APIServer = "https://kube.cluster/peanut"
	contents := suite.generateKubeconfig(Config{})
	suite.Contains(contents, "server: https://kube.cluster/peanut")
}

func (suite *KubeconfigTestSuite) TestSetsServiceAccount() {
	suite.initKube.ServiceAccount = "chef"
	contents := suite.generateKubeconfig(Config{})
	suite.Contains(contents, "user: chef")
	suite.Contains(contents, "name: chef")
}

func (suite *KubeconfigTestSuite) TestSetsToken() {
	suite.initKube.Token = "eWVhaCB3ZSB0b2tpbic"
	contents := suite.generateKubeconfig(Config{})
	suite.Contains(contents, "token: eWVhaCB3ZSB0b2tpbic")
}

func (suite *KubeconfigTestSuite) TestSetsCertificate() {
	suite.initKube.Certificate = "d293LCB5b3UgYXJlIHNvIGNvb2wgZm9yIHNtb2tpbmcgd2VlZCDwn5mE"
	contents := suite.generateKubeconfig(Config{})
	suite.Contains(contents, "certificate-authority-data: d293LCB5b3UgYXJlIHNvIGNvb2wgZm9yIHNtb2tpbmcgd2VlZCDwn5mE")
}

func (suite *KubeconfigTestSuite) TestSetsSkipTLSVerify() {
	suite.initKube.SkipTLSVerify = true
	contents := suite.generateKubeconfig(Config{})
	suite.Contains(contents, "insecure-skip-tls-verify: true")
}

func (suite *KubeconfigTestSuite) generateKubeconfig(cfg Config) string {
	suite.Require().NoError(suite.initKube.Prepare(cfg))
	suite.Require().NoError(suite.initKube.Execute(cfg))

	contents, err := ioutil.ReadFile(suite.configFile.Name())
	suite.Require().NoError(err)

	conf := map[string]interface{}{}
	suite.NoError(yaml.UnmarshalStrict(contents, &conf))

	return string(contents)
}
