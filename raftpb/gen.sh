go get -v github.com/gogo/protobuf/protoc-gen-gogofaster
protoc --gogofaster_out=. --gogofaster_opt=paths=source_relative -I=. raft.proto