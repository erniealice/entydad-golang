module github.com/erniealice/entydad-golang

go 1.25.1

require github.com/erniealice/pyeza-golang v0.0.8-alpha

require (
	github.com/erniealice/esqyma v0.0.0
	github.com/erniealice/hybra-golang v0.0.0-00010101000000-000000000000
)

require (
	github.com/yuin/goldmark v1.7.16 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251002232023-7c0ddcbb5797 // indirect
	google.golang.org/grpc v1.75.1 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace leapfor.xyz/copya => ../../../master-monorepo-v2/packages/copya

replace github.com/erniealice/esqyma => ../esqyma-ryta

replace github.com/erniealice/lyngua => ../lyngua-ryta

replace leapfor.xyz/vya => ../../../master-monorepo-v2/packages/vya

replace github.com/erniealice/hybra-golang => ../hybra-golang-ryta

replace github.com/erniealice/pyeza-golang => ../pyeza-golang-ryta
