package k8s

import (
	"errors"
	"encoding/base64"
	"k8s.io/client-go/rest"
	"github.com/unectio/util"
	"k8s.io/client-go/kubernetes"
	"k8s.io/api/extensions/v1beta1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	InClusterConfig	=	"$incluster"
)

type DepAPI interface {
	Create(*v1beta1.Deployment) (*v1beta1.Deployment, error)
	Get(string, meta.GetOptions) (*v1beta1.Deployment, error)
	Update(*v1beta1.Deployment) (*v1beta1.Deployment, error)
	Delete(string, *meta.DeleteOptions) error
}

type Client interface {
	Deps() DepAPI
	Notify(*EventHandlers, interface{})
}

type KubeNsClient struct {
	ns	string
	c	*kubernetes.Clientset
}

func (kc *KubeNsClient)Deps() DepAPI {
	return kc.c.Extensions().Deployments(kc.ns)
}

type SaConfig struct {
	Host		string		`yaml:"host"`
	Token		string		`yaml:"token"`
	CaCert		string		`yaml:"ca_cert"`
}

func Connect(cfg, ns string) (Client, error) {
	var kconf *rest.Config

	if cfg == InClusterConfig {
		var err error

		kconf, err = rest.InClusterConfig()
		if err != nil {
			return nil, errors.New("Error making in-cluster config: " + err.Error())
		}
	} else {
		var sacfg SaConfig

		err := util.LoadYAML(cfg, &sacfg)
		if err != nil {
			return nil, errors.New("Error loading sa config: " + err.Error())
		}

		cert, err := base64.StdEncoding.DecodeString(sacfg.CaCert)
		if err != nil {
			return nil, errors.New("Error base64-decodinf cert: " + err.Error())
		}

		kconf = &rest.Config{
			Host:			sacfg.Host,
			BearerToken:		sacfg.Token,
			TLSClientConfig:	rest.TLSClientConfig{CAData: cert},
		}
	}

	return connectTo(kconf, ns)
}

func connectTo(kconf *rest.Config, ns string) (Client, error) {
	var err error

	knc := &KubeNsClient{}

	knc.c, err = kubernetes.NewForConfig(kconf)
	if err != nil {
		return nil, errors.New("Error making k8s client: " + err.Error())
	}

	knc.ns = ns
	return knc, nil
}
