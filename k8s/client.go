/////////////////////////////////////////////////////////////////////////////////
//
// Copyright (C) 2019-2020, Unectio Inc, All Right Reserved.
//
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
/////////////////////////////////////////////////////////////////////////////////

package k8s

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/unectio/util"
	v1 "k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	InClusterConfig = "$incluster"
)

type DepAPI interface {
	Create(ctx context.Context, deployment *v1.Deployment, opts metaV1.CreateOptions) (*v1.Deployment, error)
	Get(ctx context.Context, name string, opts metaV1.GetOptions) (*v1.Deployment, error)
	Update(ctx context.Context, deployment *v1.Deployment, opts metaV1.UpdateOptions) (*v1.Deployment, error)
	Delete(ctx context.Context, name string, opts metaV1.DeleteOptions) error
}

type Client interface {
	Deps() DepAPI
	Notify(*EventHandlers, interface{})
}

type KubeNsClient struct {
	ns string
	c  *kubernetes.Clientset
}

func (kc *KubeNsClient) Deps() DepAPI {
	return kc.c.AppsV1().Deployments(kc.ns)
}

type SaConfig struct {
	Host   string `yaml:"host"`
	Token  string `yaml:"token"`
	CaCert string `yaml:"ca_cert"`
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
			Host:            sacfg.Host,
			BearerToken:     sacfg.Token,
			TLSClientConfig: rest.TLSClientConfig{CAData: cert},
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
