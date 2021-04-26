package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/hyperyuri/server-stream-grpc-go/pb"
)

type TestServiceServer struct {
	pb.UnimplementedTestServiceServer
}

func (s *TestServiceServer) Download(req *pb.Request, responseStream pb.TestService_DownloadServer) error {
	file, err := os.Open(req.GetFileName())
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()

	const fileChunk = 1 * (1 << 20) // 1 MB, change this to your requirement

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)

	for i := uint64(0); i < totalPartsNum; i++ {

		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		resp := &pb.Response{
			FileChunk: partBuffer,
			Proccess:  int32(i),
			Total:     int32(totalPartsNum),
		}

		err = responseStream.Send(resp)
		if err != nil {
			log.Println("error while sending chunk:", err)
			return err
		}
	}
	return nil
}
