package username_test

import (
	"github.com/proofpoint/kapprover/inspectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	certificates "k8s.io/api/certificates/v1beta1"
	"k8s.io/client-go/kubernetes"
	"testing"

	_ "github.com/proofpoint/kapprover/inspectors/group"
)

var (
	client *kubernetes.Clientset
)

func TestInspect(t *testing.T) {
	inspector, exists := inspectors.Get("username")
	if !exists {
		t.Fatal("Expected inspectors.Get(\"username\") to exist, did not")
	}

	for username, expectedMessage := range map[string]string{
		"kubelet-bootstrap": "",
		"someone-else":      "Requesting user is not kubelet-bootstrap",
	} {
		assertInspectionResult(t, inspector, username, expectedMessage)
	}
}

func TestInspectConfigured(t *testing.T) {
	inspector, exists := inspectors.Get("username")
	require.True(t, exists, "inspectors.Get(\"username\") to exist")

	inspector, err := inspector.Configure("some-user")
	assert.NoError(t, err, "Configure")

	for username, expectedMessage := range map[string]string{
		"some-user":         "",
		"kubelet-bootstrap": "Requesting user is not some-user",
	} {
		assertInspectionResult(t, inspector, username, expectedMessage)
	}
}

func assertInspectionResult(t *testing.T, inspector inspectors.Inspector, username string, expectedMessage string) {
	request := certificates.CertificateSigningRequest{
		Spec: certificates.CertificateSigningRequestSpec{
			Username: username,
			Groups: []string{
				"someRandomGroup",
				"someOtherGroup",
			},
		},
	}
	message, err := inspector.Inspect(client, &request)
	assert.Equal(t, expectedMessage, message, "Username %s", username)
	assert.NoError(t, err, "Username %s", username)
}
