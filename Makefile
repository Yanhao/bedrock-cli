default: bedrock-cli

bedrock-cli: main.go
	make -C proto

	cp proto/proxy.pb.go clients/proxy/
	cp proto/proxy_grpc.pb.go clients/proxy/

	rm -f proto/*.cc proto/*.h
	go build -race

clean:
	rm -f bedrock-cli
	make clean -C proto

.PHONY: clean proto bedrock-cli
