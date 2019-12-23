package run

import (
	"io/ioutil"
	"os"
	"text/template"
	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
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
		ConfigFile:   configFile.Name(),
	}
	cfg := Config{
		Namespace: "Cisco",
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
		ConfigFile:   "/usr/foreign/exclude/kubeprofig",
	}

	cfg := Config{}
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
		ConfigFile:   configFile.Name(),
	}

	cfg := Config{}

	suite.NoError(init.Prepare(cfg)) // consistency check; we should be starting in a happy state

	init.APIServer = ""
	suite.Error(init.Prepare(cfg), "APIServer should be required.")

	init.APIServer = "Sysadmin"
	init.Token = ""
	suite.Error(init.Prepare(cfg), "Token should be required.")

	init.Token = "Aspire virtual currency"
	init.Certificate = ""
	suite.Error(init.Prepare(cfg), "Certificate should be required.")

	init.SkipTLSVerify = true
	suite.NoError(init.Prepare(cfg), "Certificate should not be required if SkipTLSVerify is true")
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
		ConfigFile:   configFile.Name(),
	}

	cfg := Config{}

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
