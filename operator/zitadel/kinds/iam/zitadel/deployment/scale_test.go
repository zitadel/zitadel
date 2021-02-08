package deployment

import (
	"errors"
	"testing"

	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/orbos/pkg/labels/mocklabels"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeployment_Scale1(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"
	replicaCount := 1

	client.EXPECT().ScaleDeployment(namespace, mocklabels.NameVal, replicaCount).Times(1).Return(nil)
	scaleFunc := GetScaleFunc(monitor, namespace, mocklabels.Name)
	assert.NotNil(t, scaleFunc)
	ensure := scaleFunc(replicaCount)
	assert.NotNil(t, ensure)
	assert.NoError(t, ensure(client))
}

func TestDeployment_Scale2(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"
	replicaCount := 0

	client.EXPECT().ScaleDeployment(namespace, mocklabels.NameVal, replicaCount).Times(1).Return(nil)
	scaleFunc := GetScaleFunc(monitor, namespace, mocklabels.Name)
	assert.NotNil(t, scaleFunc)
	ensure := scaleFunc(replicaCount)
	assert.NotNil(t, ensure)
	assert.NoError(t, ensure(client))
}

func TestDeployment_Scale3(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"
	replicaCount := 3

	client.EXPECT().ScaleDeployment(namespace, mocklabels.NameVal, replicaCount).Times(1).Return(nil)
	scaleFunc := GetScaleFunc(monitor, namespace, mocklabels.Name)
	assert.NotNil(t, scaleFunc)
	ensure := scaleFunc(replicaCount)
	assert.NotNil(t, ensure)
	assert.NoError(t, ensure(client))
}

func TestDeployment_ScaleFailure1(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"
	replicaCount := 0

	client.EXPECT().ScaleDeployment(namespace, mocklabels.NameVal, replicaCount).Times(1).Return(errors.New("fail"))
	scaleFunc := GetScaleFunc(monitor, namespace, mocklabels.Name)
	assert.NotNil(t, scaleFunc)
	ensure := scaleFunc(replicaCount)
	assert.NotNil(t, ensure)
	assert.Error(t, ensure(client))
}
