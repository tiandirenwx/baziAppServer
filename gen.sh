#trpc create -protofile= <-rpconly>
#trpc命令可以由https://github.com/trpc-group/trpc-cmdline编译得到
trpc create --protodir=.  --protofile=./proto/bazi.proto --output="./trpc-go" --goversion="1.21.6" --validate=true --alias -f