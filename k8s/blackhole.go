package k8s

import (
	"errors"
	"k8s.io/api/extensions/v1beta1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BlackHole struct {}
var bhError = errors.New("BlackHole kuber")

func (bh *BlackHole)Deps() DepAPI { return bh }
func (bh *BlackHole)Notify(_ *EventHandlers, _ interface{}) {}

func (_ *BlackHole)Create(_ *v1beta1.Deployment) (*v1beta1.Deployment, error) {
	return nil, bhError
}

func (_ *BlackHole)Get(_ string, _ meta.GetOptions) (*v1beta1.Deployment, error) {
	return nil, bhError
}

func (_ *BlackHole)Update(_ *v1beta1.Deployment) (*v1beta1.Deployment, error) {
	return nil, bhError
}

func (_ *BlackHole)Delete(_ string, _ *meta.DeleteOptions) error {
	return nil
}
