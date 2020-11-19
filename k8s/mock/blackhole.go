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
	"errors"

	v1 "k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/unectio/util/k8s"
)

type BlackHole struct{}

var bhError = errors.New("BlackHole kuber")

func (bh *BlackHole) Deps() k8s.DepAPI                           { return bh }
func (bh *BlackHole) Notify(_ *k8s.EventHandlers, _ interface{}) {}

func (_ *BlackHole) Create(_ context.Context, _ *v1.Deployment, _ metaV1.CreateOptions) (*v1.Deployment, error) {
	return nil, bhError
}

func (_ *BlackHole) Get(_ context.Context, _ string, _ metaV1.GetOptions) (*v1.Deployment, error) {
	return nil, bhError
}

func (_ *BlackHole) Update(_ context.Context, _ *v1.Deployment, _ metaV1.UpdateOptions) (*v1.Deployment, error) {
	return nil, bhError
}

func (_ *BlackHole) Delete(_ context.Context, _ string, _ metaV1.DeleteOptions) error {
	return nil
}
