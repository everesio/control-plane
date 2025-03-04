package upgrade

import (
	"errors"
	"testing"

	"github.com/kyma-incubator/hydroform/install/installation"
	installationMocks "github.com/kyma-project/control-plane/components/provisioner/internal/installation/mocks"
	"github.com/kyma-project/control-plane/components/provisioner/internal/model"
	"github.com/kyma-project/control-plane/components/provisioner/internal/util"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	kubeconfig = `apiVersion: v1
clusters:
- cluster:
    server: https://192.168.64.4:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURBRENDQWVpZ0F3SUJBZ0lCQWpBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwdGFXNXAKYTNWaVpVTkJNQjRYRFRFNU1URXhOekE0TXpBek1sb1hEVEl3TVRFeE56QTRNekF6TWxvd01URVhNQlVHQTFVRQpDaE1PYzNsemRHVnRPbTFoYzNSbGNuTXhGakFVQmdOVkJBTVREVzFwYm1scmRXSmxMWFZ6WlhJd2dnRWlNQTBHCkNTcUdTSWIzRFFFQkFRVUFBNElCRHdBd2dnRUtBb0lCQVFDNmY2SjZneElvL2cyMHArNWhybklUaUd5SDh0VW0KWGl1OElaK09UKyt0amd1OXRneXFnbnNsL0dDT1Q3TFo4ejdOVCttTEdKL2RLRFdBV3dvbE5WTDhxMzJIQlpyNwpDaU5hK3BBcWtYR0MzNlQ2NEQyRjl4TEtpVVpuQUVNaFhWOW1oeWVCempscTh1NnBjT1NrY3lJWHRtdU9UQUVXCmErWlp5UlhOY3BoYjJ0NXFUcWZoSDhDNUVDNUIrSm4rS0tXQ2Y1Nm5KZGJQaWduRXh4SFlaMm9TUEc1aXpkbkcKZDRad2d0dTA3NGttaFNtNXQzbjgyNmovK29tL25VeWdBQ24yNmR1K21aZzRPcWdjbUMrdnBYdUEyRm52bk5LLwo5NWErNEI3cGtNTER1bHlmUTMxcjlFcStwdHBkNUR1WWpldVpjS1Bxd3ZVcFUzWVFTRUxVUzBrUkFnTUJBQUdqClB6QTlNQTRHQTFVZER3RUIvd1FFQXdJRm9EQWRCZ05WSFNVRUZqQVVCZ2dyQmdFRkJRY0RBUVlJS3dZQkJRVUgKQXdJd0RBWURWUjBUQVFIL0JBSXdBREFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBQ3JnbExWemhmemZ2aFNvUgowdWNpNndBZDF6LzA3bW52MDRUNmQyTkpjRG80Uzgwa0o4VUJtRzdmZE5qMlJEaWRFbHRKRU1kdDZGa1E1TklOCk84L1hJdENiU0ZWYzRWQ1NNSUdPcnNFOXJDajVwb24vN3JxV3dCbllqYStlbUVYOVpJelEvekJGU3JhcWhud3AKTkc1SmN6bUg5ODRWQUhGZEMvZWU0Z2szTnVoV25rMTZZLzNDTTFsRkxlVC9Cbmk2K1M1UFZoQ0x3VEdmdEpTZgorMERzbzVXVnFud2NPd3A3THl2K3h0VGtnVmdSRU5RdTByU2lWL1F2UkNPMy9DWXdwRTVIRFpjalM5N0I4MW0yCmVScVBENnVoRjVsV3h4NXAyeEd1V2JRSkY0WnJzaktLTW1CMnJrUnR5UDVYV2xWZU1mR1VjbFdjc1gxOW91clMKaWpKSTFnPT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBdW4raWVvTVNLUDROdEtmdVlhNXlFNGhzaC9MVkpsNHJ2Q0dmamsvdnJZNEx2YllNCnFvSjdKZnhnamsreTJmTSt6VS9waXhpZjNTZzFnRnNLSlRWUy9LdDlod1dhK3dvald2cVFLcEZ4Z3Qrayt1QTkKaGZjU3lvbEdad0JESVYxZlpvY25nYzQ1YXZMdXFYRGtwSE1pRjdacmprd0JGbXZtV2NrVnpYS1lXOXJlYWs2bgo0Ui9BdVJBdVFmaVovaWlsZ24rZXB5WFd6NG9KeE1jUjJHZHFFanh1WXMzWnhuZUdjSUxidE8rSkpvVXB1YmQ1Ci9OdW8vL3FKdjUxTW9BQXA5dW5idnBtWU9EcW9ISmd2cjZWN2dOaFo3NXpTdi9lV3Z1QWU2WkRDdzdwY24wTjkKYS9SS3ZxYmFYZVE3bUkzcm1YQ2o2c0wxS1ZOMkVFaEMxRXRKRVFJREFRQUJBb0lCQVFDTEVFa3pXVERkYURNSQpGb0JtVGhHNkJ1d0dvMGZWQ0R0TVdUWUVoQTZRTjI4QjB4RzJ3dnpZNGt1TlVsaG10RDZNRVo1dm5iajJ5OWk1CkVTbUxmU3VZUkxlaFNzaTVrR0cwb1VtR3RGVVQ1WGU3cWlHMkZ2bm9GRnh1eVg5RkRiN3BVTFpnMEVsNE9oVkUKTzI0Q1FlZVdEdXc4ZXVnRXRBaGJ3dG1ERElRWFdPSjcxUEcwTnZKRHIwWGpkcW1aeExwQnEzcTJkZTU2YmNjawpPYzV6dmtJNldrb0o1TXN0WkZpU3pVRDYzN3lIbjh2NGd3cXh0bHFoNWhGLzEwV296VmZqVGdWSG0rc01ZaU9SCmNIZ0dMNUVSbDZtVlBsTTQzNUltYnFnU1R2NFFVVGpzQjRvbVBsTlV5Yksvb3pPSWx3RjNPTkJjVVV6eDQ1cGwKSHVJQlQwZ1JBb0dCQU9SR2lYaVBQejdsay9Bc29tNHkxdzFRK2hWb3Yvd3ovWFZaOVVkdmR6eVJ1d3gwZkQ0QgpZVzlacU1hK0JodnB4TXpsbWxYRHJBMklYTjU3UEM3ZUo3enhHMEVpZFJwN3NjN2VmQUN0eDN4N0d0V2pRWGF2ClJ4R2xDeUZxVG9LY3NEUjBhQ0M0Um15VmhZRTdEY0huLy9oNnNzKys3U2tvRVMzNjhpS1RiYzZQQW9HQkFORW0KTHRtUmZieHIrOE5HczhvdnN2Z3hxTUlxclNnb2NmcjZoUlZnYlU2Z3NFd2pMQUs2ZHdQV0xWQmVuSWJ6bzhodApocmJHU1piRnF0bzhwS1Q1d2NxZlpKSlREQnQxYmhjUGNjWlRmSnFmc0VISXc0QW5JMVdRMlVzdzVPcnZQZWhsCmh0ek95cXdBSGZvWjBUTDlseTRJUHRqbXArdk1DQ2NPTHkwanF6NWZBb0dCQUlNNGpRT3hqSkN5VmdWRkV5WTMKc1dsbE9DMGdadVFxV3JPZnY2Q04wY1FPbmJCK01ZRlBOOXhUZFBLeC96OENkVyszT0syK2FtUHBGRUdNSTc5cApVdnlJdUxzTGZMZDVqVysyY3gvTXhaU29DM2Z0ZmM4azJMeXEzQ2djUFA5VjVQQnlUZjBwRU1xUWRRc2hrRG44CkRDZWhHTExWTk8xb3E5OTdscjhMY3A2L0FvR0FYNE5KZC9CNmRGYjRCYWkvS0lGNkFPQmt5aTlGSG9iQjdyVUQKbTh5S2ZwTGhrQk9yNEo4WkJQYUZnU09ENWhsVDNZOHZLejhJa2tNNUVDc0xvWSt4a1lBVEpNT3FUc3ZrOThFRQoyMlo3Qy80TE55K2hJR0EvUWE5Qm5KWDZwTk9XK1ErTWRFQTN6QzdOZ2M3U2U2L1ZuNThDWEhtUmpCeUVTSm13CnI3T1BXNDhDZ1lBVUVoYzV2VnlERXJxVDBjN3lIaXBQbU1wMmljS1hscXNhdC94YWtobENqUjZPZ2I5aGQvNHIKZm1wUHJmd3hjRmJrV2tDRUhJN01EdDJrZXNEZUhRWkFxN2xEdjVFT2k4ZG1uM0ZPNEJWczhCOWYzdm52MytmZwpyV2E3ZGtyWnFudU12cHhpSWlqOWZEak9XbzdxK3hTSFcxWWdSNGV2Q1p2NGxJU0FZRlViemc9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
`
)

func TestUpgradeKymaStep_Run(t *testing.T) {
	t.Run("should return error when kubeconfig is nil", func(t *testing.T) {
		//given
		upgradeStep := NewUpgradeKymaStep(nil, nextStageName, 0)

		cluster := model.Cluster{Kubeconfig: nil}

		//when
		_, err := upgradeStep.Run(cluster, model.Operation{}, logrus.New())

		//then
		require.Error(t, err)
	})

	t.Run("should return error when kubeconfig is unparsable", func(t *testing.T) {
		//given
		upgradeStep := NewUpgradeKymaStep(nil, nextStageName, 0)

		kubeconfig := "wrong konfig"
		cluster := model.Cluster{Kubeconfig: &kubeconfig}

		//when
		_, err := upgradeStep.Run(cluster, model.Operation{}, logrus.New())

		//then
		require.Error(t, err)
	})

	t.Run("should return error when failed to check installation CR state", func(t *testing.T) {
		//given
		installationClient := &installationMocks.Service{}
		installationClient.On("CheckInstallationState", mock.Anything).Return(installation.InstallationState{}, errors.New("some error"))

		upgradeStep := NewUpgradeKymaStep(installationClient, nextStageName, 0)

		cluster := model.Cluster{Kubeconfig: util.StringPtr(kubeconfig)}

		//when
		_, err := upgradeStep.Run(cluster, model.Operation{}, logrus.New())

		//then
		require.Error(t, err)
	})

	t.Run("should return error when installation CR is not present on the cluster", func(t *testing.T) {
		//given
		installationClient := &installationMocks.Service{}
		installationClient.On("CheckInstallationState", mock.Anything).Return(installation.InstallationState{State: "NoInstallation"}, nil)

		upgradeStep := NewUpgradeKymaStep(installationClient, nextStageName, 0)

		cluster := model.Cluster{Kubeconfig: util.StringPtr(kubeconfig)}

		//when
		_, err := upgradeStep.Run(cluster, model.Operation{}, logrus.New())

		//then
		require.Error(t, err)
	})

	t.Run("should return next step when upgrade already in progress", func(t *testing.T) {
		//given
		installationClient := &installationMocks.Service{}
		installationClient.On("CheckInstallationState", mock.Anything).Return(installation.InstallationState{}, installation.InstallationError{ShortMessage: "upgrade in progress", Recoverable: true})

		upgradeStep := NewUpgradeKymaStep(installationClient, nextStageName, 0)

		cluster := model.Cluster{Kubeconfig: util.StringPtr(kubeconfig)}

		//when
		result, err := upgradeStep.Run(cluster, model.Operation{}, logrus.New())

		//then
		require.NoError(t, err)
		assert.Equal(t, nextStageName, result.Stage)
	})

	t.Run("should trigger upgrade when installation is in unrecoverable error state", func(t *testing.T) {
		//given
		installationClient := &installationMocks.Service{}
		installationClient.On("CheckInstallationState", mock.Anything).Return(installation.InstallationState{}, installation.InstallationError{ShortMessage: "error", Recoverable: false})
		installationClient.On("TriggerUpgrade", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		upgradeStep := NewUpgradeKymaStep(installationClient, nextStageName, 0)

		cluster := model.Cluster{Kubeconfig: util.StringPtr(kubeconfig)}

		//when
		result, err := upgradeStep.Run(cluster, model.Operation{}, logrus.New())

		//then
		require.NoError(t, err)
		assert.Equal(t, nextStageName, result.Stage)
	})

	t.Run("should return error when fail to trigger update", func(t *testing.T) {
		//given
		installationClient := &installationMocks.Service{}
		installationClient.On("CheckInstallationState", mock.Anything).Return(installation.InstallationState{State: "Installed"}, nil)
		installationClient.On("TriggerUpgrade", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("bad error"))

		upgradeStep := NewUpgradeKymaStep(installationClient, nextStageName, 0)

		cluster := model.Cluster{Kubeconfig: util.StringPtr(kubeconfig)}

		//when
		_, err := upgradeStep.Run(cluster, model.Operation{}, logrus.New())

		//then
		require.Error(t, err)
	})

	t.Run("should return next step when upgrade successfully triggered", func(t *testing.T) {
		//given
		installationClient := &installationMocks.Service{}
		installationClient.On("CheckInstallationState", mock.Anything).Return(installation.InstallationState{State: "Installed"}, nil)
		installationClient.On("TriggerUpgrade", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		upgradeStep := NewUpgradeKymaStep(installationClient, nextStageName, 0)

		cluster := model.Cluster{Kubeconfig: util.StringPtr(kubeconfig)}

		//when
		result, err := upgradeStep.Run(cluster, model.Operation{}, logrus.New())

		//then
		require.NoError(t, err)
		assert.Equal(t, nextStageName, result.Stage)
	})

	t.Run("should return next step when installation state in progress", func(t *testing.T) {
		//given
		installationClient := &installationMocks.Service{}
		installationClient.On("CheckInstallationState", mock.Anything).Return(installation.InstallationState{State: "In Progress"}, nil)

		upgradeStep := NewUpgradeKymaStep(installationClient, nextStageName, 0)

		cluster := model.Cluster{Kubeconfig: util.StringPtr(kubeconfig)}

		//when
		result, err := upgradeStep.Run(cluster, model.Operation{}, logrus.New())

		//then
		require.NoError(t, err)
		assert.Equal(t, nextStageName, result.Stage)
	})
}
