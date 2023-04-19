package message

//go:generate protoc -I . --go_out=plugins=grpc:. message.proto
//go:generate protoc -I . --go-http_out=paths=source_relative:. message.proto
//go:generate protoc -I . --doc_out=:. --doc_opt=html,proto.html message.proto
