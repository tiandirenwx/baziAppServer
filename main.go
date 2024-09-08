package main

import (
	"context"
	"trpc.app/app/baziAppServer/service/paipan"
	pb "trpc.bazi.paipan"

	_ "trpc.group/trpc-go/trpc-filter/debuglog"
	_ "trpc.group/trpc-go/trpc-filter/recovery"
	_ "trpc.group/trpc-go/trpc-filter/validation"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
	_ "trpc.group/trpc-go/trpc-go/server"
)

//	func RunBaziTrpcServer(ctx context.Context) {
//		trpc.ServerConfigPath = "./config/trpc_go.yaml"
//		s := trpc.NewServer()
//		pb.RegisterBaziPaipanService(s.Service("trpc.bazi.paipan.BaziPaipan"), &paipan.BaziPaipanImpl{})
//		if err := s.Serve(); err != nil {
//			log.FatalContext(ctx, err)
//		}
//	}
//
//	func main() {
//		ctx, cancel := context.WithCancel(context.Background())
//		defer cancel()
//		var wg sync.WaitGroup
//		wg.Add(1)
//		go func() {
//			sigChan := make(chan os.Signal, 1)
//			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
//			<-sigChan
//			//cancel()
//		}()
//
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			RunBaziTrpcServer(ctx)
//		}()
//
//		<-ctx.Done() // 阻塞，直到上下文被取消
//		wg.Wait()
//
// }
func main() {
	trpc.ServerConfigPath = "./config/trpc_go.yaml"
	s := trpc.NewServer()
	pb.RegisterBaziPaipanService(s.Service("bazi.paipan.BaziPaipan_trpc"), &paipan.BaziPaipanImpl{})
	if err := s.Serve(); err != nil {
		log.FatalContext(context.Background(), err)
	}
}
