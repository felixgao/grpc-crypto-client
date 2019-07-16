package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	pb "github.com/felixgao/grpc-crypto-proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// var (
// 	// ErrClosed is the error when the client pool is closed
// 	ErrClosed = errors.New("grpc pool: client pool is closed")
// 	// ErrTimeout is the error when the client pool timed out
// 	ErrTimeout = errors.New("grpc pool: client pool timed out")
// 	// ErrAlreadyClosed is the error when the client conn was already closed
// 	ErrAlreadyClosed = errors.New("grpc pool: the connection was already closed")
// 	// ErrFullPool is the error when the pool is already full
// 	ErrFullPool = errors.New("grpc pool: closing a ClientConn into a full pool")
// )

// Factory is a function type creating a grpc client
// type Factory func() (*grpc.ClientConn, error)

// // Pool is the grpc client pool
// type Pool struct {
// 	clients         chan ClientConn
// 	factory         Factory
// 	idleTimeout     time.Duration
// 	maxLifeDuration time.Duration
// 	mu              sync.RWMutex
// }

// // ClientConn is the wrapper for a grpc client conn
// type ClientConn struct {
// 	*grpc.ClientConn
// 	pool          *Pool
// 	timeUsed      time.Time
// 	timeInitiated time.Time
// 	unhealthy     bool
// }

func getIDPSConnection(host string) (*grpc.ClientConn, error) {
	wd, _ := os.Getwd()
	parentDir := filepath.Dir(wd)
	certFile := filepath.Join(parentDir, "keys", "cert.pem")
	creds, _ := credentials.NewClientTLSFromFile(certFile, "")
	return grpc.Dial(host, grpc.WithTransportCredentials(creds))
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func idpsEncrypt(ctx context.Context, client pb.CryptoClient, jobID string, data string) (*pb.EncryptedText, error) {
	defer timeTrack(time.Now(), "Encrypt")
	encryptedResponse, err := client.EncryptData(ctx, &pb.TrackableDecryptedRequest{JobId: jobID, Data: &pb.DecryptedText{TextData: data}})
	if err != nil {
		eMsg := fmt.Sprintf("Failed to Decrypt data for [JobID=%s]", jobID)
		log.Fatal(eMsg, err)
		return nil, errors.New(eMsg)
	}
	return encryptedResponse, nil
}

func idpsDecrypt(ctx context.Context, client pb.CryptoClient, jobID string, data string) (*pb.DecryptedText, error) {
	defer timeTrack(time.Now(), "Decrypt")
	decryptedResponse, err := client.DecryptData(ctx, &pb.TrackableEncryptedRequest{JobId: jobID, Data: &pb.EncryptedText{EncryptedData: data}})
	if err != nil {
		eMsg := fmt.Sprintf("Failed to Decrypt data for [JobID=%s]", jobID)
		log.Fatal(eMsg, err)
		return nil, errors.New(eMsg)
	}
	return decryptedResponse, nil
}

// NewCryptoClient Function to create a new client.
// serves as a factory method
func NewCryptoClient(conn *grpc.ClientConn) pb.CryptoClient {
	return pb.NewCryptoClient(conn)
}
