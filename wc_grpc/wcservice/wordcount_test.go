package wcservice

import (
	"context"
	"log"
	"net"
	"testing"
	pb "wc_grpc/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterWordCountServieceServer(s, &WcServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestWc(t *testing.T) {

	//Dial a connection to grpc Server
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	//Create new Client
	c := pb.NewWordCountServieceClient(conn)

	//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer conn.Close()

	//Result
	response, err := c.WordCount(ctx, &pb.TextRequest{Text: "My name is Anibrata, I am Anibrata."})
	if err != nil {
		t.Fatal("Could not count word: \n", err)
	}
	t.Log("WordCount:\n")
	for _, value := range response.WcList {
		t.Logf(`"Word: %s	Count: %d`, value.Word, value.Count)
	}
}
