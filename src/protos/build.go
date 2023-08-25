package protoc

//go:generate protoc --proto_path=. --go_out=. --go_opt=paths=source_relative api.proto
