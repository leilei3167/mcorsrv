//go:generate sh -c "protoc --proto_path=$GOPATH/src --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative *.proto"

package proto
