package run

import (
	"github.com/go-yaml/yaml"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"testing"
	"text/template"
)

type InitKubeTestSuite struct {
	suite.Suite
}

func TestInitKubeTestSuite(t *testing.T) {
	suite.Run(t, new(InitKubeTestSuite))
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

	init := InitKube{
		APIServer:    "Sysadmin",
		Certificate:  "CCNA",
		Token:        "Aspire virtual currency",
		TemplateFile: templateFile.Name(),
	}
	cfg := Config{
		Namespace:  "Cisco",
		KubeConfig: configFile.Name(),
	}
	err = init.Prepare(cfg)
	suite.Require().Nil(err)

	suite.IsType(&template.Template{}, init.template)
	suite.NotNil(init.configFile)

	err = init.Execute(cfg)
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

	cfg := Config{
		KubeConfig: configFile.Name(),
		Namespace:  "marshmallow",
	}
	init := InitKube{
		TemplateFile:   "../../assets/kubeconfig.tpl", // the actual kubeconfig template
		APIServer:      "https://kube.cluster/peanut",
		ServiceAccount: "chef",
		Token:          "eWVhaCB3ZSB0b2tpbic=",
		Certificate:    "d293LCB5b3UgYXJlIHNvIGNvb2wgZm9yIHNtb2tpbmcgd2VlZCDwn5mE",
	}
	suite.Require().NoError(init.Prepare(cfg))
	suite.Require().NoError(init.Execute(cfg))

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
	init.SkipTLSVerify = true
	init.Certificate = ""

	suite.Require().NoError(init.Prepare(cfg))
	suite.Require().NoError(init.Execute(cfg))
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

	init := InitKube{
		APIServer:    "Sysadmin",
		Certificate:  "CCNA",
		Token:        "Aspire virtual currency",
		TemplateFile: templateFile.Name(),
	}
	err = init.Prepare(Config{})
	suite.Error(err)
	suite.Regexp("could not load kubeconfig .* function .* not defined", err)
}

func (suite *InitKubeTestSuite) TestPrepareNonexistentTemplateFile() {
	init := InitKube{
		APIServer:    "Sysadmin",
		Certificate:  "CCNA",
		Token:        "Aspire virtual currency",
		TemplateFile: "/usr/foreign/exclude/kubeprofig.tpl",
	}
	err := init.Prepare(Config{})
	suite.Error(err)
	suite.Regexp("could not load kubeconfig .* no such file or directory", err)
}

func (suite *InitKubeTestSuite) TestPrepareCannotOpenDestinationFile() {
	templateFile, err := tempfile("kubeconfig********.yml.tpl", "hurgity burgity")
	defer os.Remove(templateFile.Name())
	suite.Require().Nil(err)
	init := InitKube{
		APIServer:    "Sysadmin",
		Certificate:  "CCNA",
		Token:        "Aspire virtual currency",
		TemplateFile: templateFile.Name(),
	}

	cfg := Config{
		KubeConfig: "/usr/foreign/exclude/kubeprofig",
	}
	err = init.Prepare(cfg)
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
	init := InitKube{
		APIServer:    "Sysadmin",
		Certificate:  "CCNA",
		Token:        "Aspire virtual currency",
		TemplateFile: templateFile.Name(),
	}

	cfg := Config{
		KubeConfig: configFile.Name(),
	}

	suite.NoError(init.Prepare(cfg)) // consistency check; we should be starting in a happy state

	init.APIServer = ""
	suite.Error(init.Prepare(cfg), "APIServer should be required.")

	init.APIServer = "Sysadmin"
	init.Token = ""
	suite.Error(init.Prepare(cfg), "Token should be required.")
}

func (suite *InitKubeTestSuite) TestPrepareDefaultsServiceAccount() {
	templateFile, err := tempfile("kubeconfig********.yml.tpl", "hurgity burgity")
	defer os.Remove(templateFile.Name())
	suite.Require().Nil(err)

	configFile, err := tempfile("kubeconfig********.yml", "")
	defer os.Remove(configFile.Name())
	suite.Require().Nil(err)

	init := InitKube{
		APIServer:    "Sysadmin",
		Certificate:  "CCNA",
		Token:        "Aspire virtual currency",
		TemplateFile: templateFile.Name(),
	}

	cfg := Config{
		KubeConfig: configFile.Name(),
	}

	init.Prepare(cfg)
	suite.Equal("helm", init.ServiceAccount)
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
