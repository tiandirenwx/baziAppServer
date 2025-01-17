// Code generated by trpc-go/trpc-cmdline v1.0.6. DO NOT EDIT.
// source: proto/bazi.proto

package trpc_bazi_paipan

import (
	"context"
	"errors"
	"fmt"

	_ "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/codec"
	_ "trpc.group/trpc-go/trpc-go/http"
	"trpc.group/trpc-go/trpc-go/server"
)

// START ======================================= Server Service Definition ======================================= START

// BaziPaipanService defines service.
type BaziPaipanService interface {
	CreateBaziPaipan(ctx context.Context, req *PaiPanRequest) (*CreatePaiPanRsp, error)

	RenderPaiPanImage(ctx context.Context, req *PaiPanRequest) (*RenderPaiPanImageRsp, error)
}

func BaziPaipanService_CreateBaziPaipan_Handler(svr interface{}, ctx context.Context, f server.FilterFunc) (interface{}, error) {
	req := &PaiPanRequest{}
	filters, err := f(req)
	if err != nil {
		return nil, err
	}
	handleFunc := func(ctx context.Context, reqbody interface{}) (interface{}, error) {
		return svr.(BaziPaipanService).CreateBaziPaipan(ctx, reqbody.(*PaiPanRequest))
	}

	var rsp interface{}
	rsp, err = filters.Filter(ctx, req, handleFunc)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func BaziPaipanService_RenderPaiPanImage_Handler(svr interface{}, ctx context.Context, f server.FilterFunc) (interface{}, error) {
	req := &PaiPanRequest{}
	filters, err := f(req)
	if err != nil {
		return nil, err
	}
	handleFunc := func(ctx context.Context, reqbody interface{}) (interface{}, error) {
		return svr.(BaziPaipanService).RenderPaiPanImage(ctx, reqbody.(*PaiPanRequest))
	}

	var rsp interface{}
	rsp, err = filters.Filter(ctx, req, handleFunc)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

// BaziPaipanServer_ServiceDesc descriptor for server.RegisterService.
var BaziPaipanServer_ServiceDesc = server.ServiceDesc{
	ServiceName: "trpc.bazi.paipan.BaziPaipan",
	HandlerType: ((*BaziPaipanService)(nil)),
	Methods: []server.Method{
		{
			Name: "/trpc.bazi.paipan.BaziPaipan/CreateBaziPaipan",
			Func: BaziPaipanService_CreateBaziPaipan_Handler,
		},
		{
			Name: "/trpc.bazi.paipan.BaziPaipan/RenderPaiPanImage",
			Func: BaziPaipanService_RenderPaiPanImage_Handler,
		},
	},
}

// RegisterBaziPaipanService registers service.
func RegisterBaziPaipanService(s server.Service, svr BaziPaipanService) {
	if err := s.Register(&BaziPaipanServer_ServiceDesc, svr); err != nil {
		panic(fmt.Sprintf("BaziPaipan register error:%v", err))
	}
}

// START --------------------------------- Default Unimplemented Server Service --------------------------------- START

type UnimplementedBaziPaipan struct{}

func (s *UnimplementedBaziPaipan) CreateBaziPaipan(ctx context.Context, req *PaiPanRequest) (*CreatePaiPanRsp, error) {
	return nil, errors.New("rpc CreateBaziPaipan of service BaziPaipan is not implemented")
}
func (s *UnimplementedBaziPaipan) RenderPaiPanImage(ctx context.Context, req *PaiPanRequest) (*RenderPaiPanImageRsp, error) {
	return nil, errors.New("rpc RenderPaiPanImage of service BaziPaipan is not implemented")
}

// END --------------------------------- Default Unimplemented Server Service --------------------------------- END

// END ======================================= Server Service Definition ======================================= END

// START ======================================= Client Service Definition ======================================= START

// BaziPaipanClientProxy defines service client proxy
type BaziPaipanClientProxy interface {
	CreateBaziPaipan(ctx context.Context, req *PaiPanRequest, opts ...client.Option) (rsp *CreatePaiPanRsp, err error)

	RenderPaiPanImage(ctx context.Context, req *PaiPanRequest, opts ...client.Option) (rsp *RenderPaiPanImageRsp, err error)
}

type BaziPaipanClientProxyImpl struct {
	client client.Client
	opts   []client.Option
}

var NewBaziPaipanClientProxy = func(opts ...client.Option) BaziPaipanClientProxy {
	return &BaziPaipanClientProxyImpl{client: client.DefaultClient, opts: opts}
}

func (c *BaziPaipanClientProxyImpl) CreateBaziPaipan(ctx context.Context, req *PaiPanRequest, opts ...client.Option) (*CreatePaiPanRsp, error) {
	ctx, msg := codec.WithCloneMessage(ctx)
	defer codec.PutBackMessage(msg)
	msg.WithClientRPCName("/trpc.bazi.paipan.BaziPaipan/CreateBaziPaipan")
	msg.WithCalleeServiceName(BaziPaipanServer_ServiceDesc.ServiceName)
	msg.WithCalleeApp("bazi")
	msg.WithCalleeServer("paipan")
	msg.WithCalleeService("BaziPaipan")
	msg.WithCalleeMethod("CreateBaziPaipan")
	msg.WithSerializationType(codec.SerializationTypePB)
	callopts := make([]client.Option, 0, len(c.opts)+len(opts))
	callopts = append(callopts, c.opts...)
	callopts = append(callopts, opts...)
	rsp := &CreatePaiPanRsp{}
	if err := c.client.Invoke(ctx, req, rsp, callopts...); err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *BaziPaipanClientProxyImpl) RenderPaiPanImage(ctx context.Context, req *PaiPanRequest, opts ...client.Option) (*RenderPaiPanImageRsp, error) {
	ctx, msg := codec.WithCloneMessage(ctx)
	defer codec.PutBackMessage(msg)
	msg.WithClientRPCName("/trpc.bazi.paipan.BaziPaipan/RenderPaiPanImage")
	msg.WithCalleeServiceName(BaziPaipanServer_ServiceDesc.ServiceName)
	msg.WithCalleeApp("bazi")
	msg.WithCalleeServer("paipan")
	msg.WithCalleeService("BaziPaipan")
	msg.WithCalleeMethod("RenderPaiPanImage")
	msg.WithSerializationType(codec.SerializationTypePB)
	callopts := make([]client.Option, 0, len(c.opts)+len(opts))
	callopts = append(callopts, c.opts...)
	callopts = append(callopts, opts...)
	rsp := &RenderPaiPanImageRsp{}
	if err := c.client.Invoke(ctx, req, rsp, callopts...); err != nil {
		return nil, err
	}
	return rsp, nil
}

// END ======================================= Client Service Definition ======================================= END
