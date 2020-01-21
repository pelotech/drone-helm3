package run

import (
	"fmt"
	"github.com/pelotech/drone-helm3/internal/env"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type RepoCertsTestSuite struct {
	suite.Suite
}

func TestRepoCertsTestSuite(t *testing.T) {
	suite.Run(t, new(RepoCertsTestSuite))
}

func (suite *RepoCertsTestSuite) TestNewRepoCerts() {
	cfg := env.Config{
		RepoCertificate:   "bGljZW5zZWQgYnkgdGhlIFN0YXRlIG9mIE9yZWdvbiB0byBwZXJmb3JtIHJlcG9zc2Vzc2lvbnM=",
		RepoCACertificate: "T3JlZ29uIFN0YXRlIExpY2Vuc3VyZSBib2FyZA==",
	}
	rc := newRepoCerts(cfg)
	suite.Require().NotNil(rc)
	suite.Equal("bGljZW5zZWQgYnkgdGhlIFN0YXRlIG9mIE9yZWdvbiB0byBwZXJmb3JtIHJlcG9zc2Vzc2lvbnM=", rc.cert)
	suite.Equal("T3JlZ29uIFN0YXRlIExpY2Vuc3VyZSBib2FyZA==", rc.caCert)
}

func (suite *RepoCertsTestSuite) TestWrite() {
	cfg := env.Config{
		RepoCertificate:   "bGljZW5zZWQgYnkgdGhlIFN0YXRlIG9mIE9yZWdvbiB0byBwZXJmb3JtIHJlcG9zc2Vzc2lvbnM=",
		RepoCACertificate: "T3JlZ29uIFN0YXRlIExpY2Vuc3VyZSBib2FyZA==",
	}
	rc := newRepoCerts(cfg)
	suite.Require().NotNil(rc)

	suite.NoError(rc.write())
	defer os.Remove(rc.certFilename)
	defer os.Remove(rc.caCertFilename)
	suite.NotEqual("", rc.certFilename)
	suite.NotEqual("", rc.caCertFilename)

	cert, err := ioutil.ReadFile(rc.certFilename)
	suite.Require().NoError(err)
	caCert, err := ioutil.ReadFile(rc.caCertFilename)
	suite.Require().NoError(err)
	suite.Equal("licensed by the State of Oregon to perform repossessions", string(cert))
	suite.Equal("Oregon State Licensure board", string(caCert))
}

func (suite *RepoCertsTestSuite) TestFlags() {
	rc := newRepoCerts(env.Config{})
	suite.Equal([]string{}, rc.flags())
	rc.certFilename = "hurgityburgity"
	suite.Equal([]string{"--cert-file", "hurgityburgity"}, rc.flags())
	rc.caCertFilename = "honglydongly"
	suite.Equal([]string{"--cert-file", "hurgityburgity", "--ca-file", "honglydongly"}, rc.flags())
}

func (suite *RepoCertsTestSuite) TestDebug() {
	stderr := strings.Builder{}
	cfg := env.Config{
		RepoCertificate:   "bGljZW5zZWQgYnkgdGhlIFN0YXRlIG9mIE9yZWdvbiB0byBwZXJmb3JtIHJlcG9zc2Vzc2lvbnM=",
		RepoCACertificate: "T3JlZ29uIFN0YXRlIExpY2Vuc3VyZSBib2FyZA==",
		Stderr:            &stderr,
		Debug:             true,
	}
	rc := newRepoCerts(cfg)
	suite.Require().NotNil(rc)

	suite.NoError(rc.write())
	defer os.Remove(rc.certFilename)
	defer os.Remove(rc.caCertFilename)

	suite.Contains(stderr.String(), fmt.Sprintf("writing repo certificate to %s", rc.certFilename))
	suite.Contains(stderr.String(), fmt.Sprintf("writing repo ca certificate to %s", rc.caCertFilename))
}
