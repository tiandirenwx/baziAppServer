package paipan

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	pb "trpc.bazi.paipan"
	_ "trpc.group/trpc-go/trpc-go/http"
)

//go:generate go mod tidy
//go:generate mockgen -destination=stub/trpc.bazi.paipan/bazi_mock.go -package=trpc_bazi_paipan -self_package=trpc.bazi.paipan --source=stub/trpc.bazi.paipan/bazi.trpc.go

func Test_baziPaipanImpl_CreateBaziPaipan(t *testing.T) {
	// Start writing mock logic.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	baziPaipanService := pb.NewMockBaziPaipanService(ctrl)
	var inorderClient []*gomock.Call
	// Expected behavior.
	m := baziPaipanService.EXPECT().CreateBaziPaipan(gomock.Any(), gomock.Any()).AnyTimes()
	m.DoAndReturn(func(ctx context.Context, req *pb.PaiPanRequest) (*pb.CreatePaiPanRsp, error) {
		s := &BaziPaipanImpl{}
		return s.CreateBaziPaipan(ctx, req)
	})
	gomock.InOrder(inorderClient...)

	// Start writing unit test logic.
	type args struct {
		ctx context.Context
		req *pb.PaiPanRequest
		rsp *pb.CreatePaiPanRsp
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rsp *pb.CreatePaiPanRsp
			var err error
			if rsp, err = baziPaipanService.CreateBaziPaipan(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("baziPaipanImpl.CreateBaziPaipan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(rsp, tt.args.rsp) {
				t.Errorf("baziPaipanImpl.CreateBaziPaipan() rsp got = %v, want %v", rsp, tt.args.rsp)
			}
		})
	}
}

func Test_baziPaipanImpl_RenderPaiPanImage(t *testing.T) {
	// Start writing mock logic.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	baziPaipanService := pb.NewMockBaziPaipanService(ctrl)
	var inorderClient []*gomock.Call
	// Expected behavior.
	m := baziPaipanService.EXPECT().RenderPaiPanImage(gomock.Any(), gomock.Any()).AnyTimes()
	m.DoAndReturn(func(ctx context.Context, req *pb.PaiPanRequest) (*pb.RenderPaiPanImageRsp, error) {
		s := &BaziPaipanImpl{}
		return s.RenderPaiPanImage(ctx, req)
	})
	gomock.InOrder(inorderClient...)

	// Start writing unit test logic.
	type args struct {
		ctx context.Context
		req *pb.PaiPanRequest
		rsp *pb.RenderPaiPanImageRsp
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rsp *pb.RenderPaiPanImageRsp
			var err error
			if rsp, err = baziPaipanService.RenderPaiPanImage(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("baziPaipanImpl.RenderPaiPanImage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(rsp, tt.args.rsp) {
				t.Errorf("baziPaipanImpl.RenderPaiPanImage() rsp got = %v, want %v", rsp, tt.args.rsp)
			}
		})
	}
}
