package restmux

import (
	"context"
	"net/url"
)

type ColRO struct {}

func(_ ColRO)Add(ctx context.Context, i Image) (Object, Error) { return nil, errBadm }
func(_ ColRO)Del(ctx context.Context, o Object) Error { return errBadm }

type ObjRO struct {}
func (_ ObjRO)Upd(ctx context.Context, _ url.Values, i Image) Error { return errBadm }

type PropRO struct {}

func (_ PropRO)Set(ctx context.Context, o Object, i Image) Error { return errBadm }
func (_ PropRO)Del(ctx context.Context, o Object) Error { return errBadm }
func (_ PropRO)Image() Image { return nil }
