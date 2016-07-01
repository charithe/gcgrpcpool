PROTO_DIR=gcgrpc
PROTO_BASENAME=gcgrpc
PROTO_TARGET=$(PROTO_DIR)/$(PROTO_BASENAME).pb.go

all: $(PROTO_TARGET) 

$(PROTO_TARGET):
	protoc -I $(PROTO_DIR) --go_out=plugins=grpc:$(PROTO_DIR) $(PROTO_DIR)/$(PROTO_BASENAME).proto 

test:
	go test

clean:
	-rm -rf $(PROTO_TARGET)
