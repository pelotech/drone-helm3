package run

import (
	"fmt"
	"github.com/pelotech/drone-helm3/internal/env"
	"github.com/stretchr/testify/suite"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"text/template"
)

type InitKubeTestSuite struct {
	suite.Suite
}

func TestInitKubeTestSuite(t *testing.T) {
	suite.Run(t, new(InitKubeTestSuite))
}

func (suite *InitKubeTestSuite) TestNewInitKube() {
	cfg := env.Config{
		SkipTLSVerify:  true,
		Certificate:    "cHJvY2xhaW1zIHdvbmRlcmZ1bCBmcmllbmRzaGlw",
		APIServer:      "98.765.43.21",
		ServiceAccount: "greathelm",
		KubeToken:      "b2YgbXkgYWZmZWN0aW9u",
		Stderr:         &strings.Builder{},
		Debug:          true,
	}

	init := NewInitKube(cfg, "conf.tpl", "conf.yml")
	suite.Equal(kubeValues{
		SkipTLSVerify:  true,
		Certificate:    "cHJvY2xhaW1zIHdvbmRlcmZ1bCBmcmllbmRzaGlw",
		APIServer:      "98.765.43.21",
		ServiceAccount: "greathelm",
		Token:          "b2YgbXkgYWZmZWN0aW9u",
	}, init.values)
	suite.Equal("conf.tpl", init.templateFilename)
	suite.Equal("conf.yml", init.configFilename)
	suite.NotNil(init.config)
}

func (suite *InitKubeTestSuite) TestPrepareExecute() {
	templateFile, err := tempfile("kubeconfig********.yml.tpl", `
certificate: {{ .Certificate }}
namespace: {{ .Namespace }}
`)
	defer os.Remove(templateFile.Name())
	suite.Require().Nil(err)

	configFile, err := tempfile("kubeconfig********.yml", "")
	defer os.Remove(configFile.Name())
	suite.Require().Nil(err)

	cfg := env.Config{
		APIServer:   "Sysadmin",
		Certificate: "CCNA",
		KubeToken:   "Aspire virtual currency",
		Namespace:   "Cisco",
	}
	init := NewInitKube(cfg, templateFile.Name(), configFile.Name())
	err = init.Prepare()
	suite.Require().Nil(err)

	suite.IsType(&template.Template{}, init.template)
	suite.NotNil(init.configFile)

	err = init.Execute()
	suite.Require().Nil(err)

	conf, err := ioutil.ReadFile(configFile.Name())
	suite.Require().Nil(err)

	want := `
certificate: CCNA
namespace: Cisco
`
	suite.Equal(want, string(conf))
}

func (suite *InitKubeTestSuite) TestExecuteGeneratesConfig() {
	configFile, err := tempfile("kubeconfig********.yml", "")
	defer os.Remove(configFile.Name())
	suite.Require().NoError(err)

	cfg := env.Config{
		APIServer:      "https://kube.cluster/peanut",
		ServiceAccount: "chef",
		KubeToken:      "eWVhaCB3ZSB0b2tpbic=",
		Certificate:    "d293LCB5b3UgYXJlIHNvIGNvb2wgZm9yIHNtb2tpbmcgd2VlZCDwn5mE",
		Namespace:      "marshmallow",
	}
	init := NewInitKube(cfg, "../../assets/kubeconfig.tpl", configFile.Name()) // the actual kubeconfig template
	suite.Require().NoError(init.Prepare())
	suite.Require().NoError(init.Execute())

	contents, err := ioutil.ReadFile(configFile.Name())
	suite.Require().NoError(err)

	// each setting should be reflected in the generated file
	expectations := []string{
		"namespace: marshmallow",
		"server: https://kube.cluster/peanut",
		"user: chef",
		"name: chef",
		"token: eWVhaCB3ZSB0b2tpbic",
		"certificate-authority-data: d293LCB5b3UgYXJlIHNvIGNvb2wgZm9yIHNtb2tpbmcgd2VlZCDwn5mE",
	}
	for _, expected := range expectations {
		suite.Contains(string(contents), expected)
	}

	// the generated config should be valid yaml, with no repeated keys
	conf := map[string]interface{}{}
	suite.NoError(yaml.UnmarshalStrict(contents, &conf))

	// test the other branch of the certificate/SkipTLSVerify conditional
	init.values.SkipTLSVerify = true
	init.values.Certificate = ""

	suite.Require().NoError(init.Prepare())
	suite.Require().NoError(init.Execute())
	contents, err = ioutil.ReadFile(configFile.Name())
	suite.Require().NoError(err)
	suite.Contains(string(contents), "insecure-skip-tls-verify: true")

	conf = map[string]interface{}{}
	suite.NoError(yaml.UnmarshalStrict(contents, &conf))
}

func (suite *InitKubeTestSuite) TestPrepareParseError() {
	templateFile, err := tempfile("kubeconfig********.yml.tpl", `{{ NonexistentFunction }}`)
	defer os.Remove(templateFile.Name())
	suite.Require().Nil(err)

	cfg := env.Config{
		APIServer:   "Sysadmin",
		Certificate: "CCNA",
		KubeToken:   "Aspire virtual currency",
	}
	init := NewInitKube(cfg, templateFile.Name(), "")
	err = init.Prepare()
	suite.Error(err)
	suite.Regexp("could not load kubeconfig .* function .* not defined", err)
}

func (suite *InitKubeTestSuite) TestPrepareNonexistentTemplateFile() {
	cfg := env.Config{
		APIServer:   "Sysadmin",
		Certificate: "CCNA",
		KubeToken:   "Aspire virtual currency",
	}
	init := NewInitKube(cfg, "/usr/foreign/exclude/kubeprofig.tpl", "")
	err := init.Prepare()
	suite.Error(err)
	suite.Regexp("could not load kubeconfig .* no such file or directory", err)
}

func (suite *InitKubeTestSuite) TestPrepareCannotOpenDestinationFile() {
	templateFile, err := tempfile("kubeconfig********.yml.tpl", "hurgity burgity")
	defer os.Remove(templateFile.Name())
	suite.Require().Nil(err)
	cfg := env.Config{
		APIServer:   "Sysadmin",
		Certificate: "CCNA",
		KubeToken:   "Aspire virtual currency",
	}
	init := NewInitKube(cfg, templateFile.Name(), "/usr/foreign/exclude/kubeprofig")

	err = init.Prepare()
	suite.Error(err)
	suite.Regexp("could not open .* for writing: .* no such file or directory", err)
}

func (suite *InitKubeTestSuite) TestPrepareRequiredConfig() {
	templateFile, err := tempfile("kubeconfig********.yml.tpl", "hurgity burgity")
	defer os.Remove(templateFile.Name())
	suite.Require().Nil(err)

	configFile, err := tempfile("kubeconfig********.yml", "")
	defer os.Remove(configFile.Name())
	suite.Require().Nil(err)

	// initial config with all required fields present
	cfg := env.Config{
		APIServer:   "Sysadmin",
		Certificate: "CCNA",
		KubeToken:   "Aspire virtual currency",
	}

	init := NewInitKube(cfg, templateFile.Name(), configFile.Name())
	suite.NoError(init.Prepare()) // consistency check; we should be starting in a happy state

	init.values.APIServer = ""
	suite.Error(init.Prepare(), "APIServer should be required.")

	init.values.APIServer = "Sysadmin"
	init.values.Token = ""
	suite.Error(init.Prepare(), "Token should be required.")
}

func (suite *InitKubeTestSuite) TestPrepareDefaultsServiceAccount() {
	templateFile, err := tempfile("kubeconfig********.yml.tpl", "hurgity burgity")
	defer os.Remove(templateFile.Name())
	suite.Require().Nil(err)

	configFile, err := tempfile("kubeconfig********.yml", "")
	defer os.Remove(configFile.Name())
	suite.Require().Nil(err)

	cfg := env.Config{
		APIServer:   "Sysadmin",
		Certificate: "CCNA",
		KubeToken:   "Aspire virtual currency",
	}
	init := NewInitKube(cfg, templateFile.Name(), configFile.Name())

	init.Prepare()
	suite.Equal("helm", init.values.ServiceAccount)
}

func (suite *InitKubeTestSuite) TestDebugOutput() {
	templateFile, err := tempfile("kubeconfig********.yml.tpl", "hurgity burgity")
	defer os.Remove(templateFile.Name())
	suite.Require().Nil(err)

	configFile, err := tempfile("kubeconfig********.yml", "")
	defer os.Remove(configFile.Name())
	suite.Require().Nil(err)

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	cfg := env.Config{
		APIServer: "http://my.kube.server/",
		KubeToken: "QSBzaW5nbGUgcm9zZQ==",
		Debug:     true,
		Stdout:    stdout,
		Stderr:    stderr,
	}
	init := NewInitKube(cfg, templateFile.Name(), configFile.Name())
	suite.NoError(init.Prepare())

	suite.Contains(stderr.String(), fmt.Sprintf("loading kubeconfig template from %s\n", templateFile.Name()))
	suite.Contains(stderr.String(), fmt.Sprintf("truncating kubeconfig file at %s\n", configFile.Name()))

	suite.NoError(init.Execute())
	suite.Contains(stderr.String(), fmt.Sprintf("writing kubeconfig file to %s\n", configFile.Name()))
}

func tempfile(name, contents string) (*os.File, error) {
	file, err := ioutil.TempFile("", name)
	if err != nil {
		return nil, err
	}
	_, err = file.Write([]byte(contents))
	if err != nil {
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}
	return file, nil
}
