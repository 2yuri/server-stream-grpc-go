package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/hyperyuri/server-stream-grpc-go/pb"
	"google.golang.org/grpc"
)

type Products struct {
	CompanyCod      string `json:"EMPRESA,omitempty"`
	Name            string `json:"DESCRICAO,omitempty"`
	BarCode         string `json:"CODIGOBARRA,omitempty"`
	Type            string `json:"EMBALAGEM,omitempty"`
	EstRetirar      string `json:"ESTARETIRAR,omitempty"`
	EstCondicional  string `json:"ESTCONDICIONAL,omitempty"`
	EstDisponivel   string `json:"ESTDISPONIVEL,omitempty"`
	EstReservado    string `json:"ESTRESERVADO,omitempty"`
	EstTransito     string `json:"ESTTRANSITO,omitempty"`
	EstVendaExterna string `json:"ESTVENDAEXTERNA,omitempty"`
}

type ArrayOfProducts struct {
	Value []Products `json:"VAL"`
}

type Company struct {
	Code     string   `json:"CODIGO"`
	Cnpj     string   `json:"CPFCNPJ"`
	Name     string   `json:"NOMEFANTASIA"`
	City     string   `json:"CIDADE"`
	State    string   `json:"ESTADO"`
	Products Products `json:"PRODUTOS,omitempty"`
}

type ArrayOfCompany struct {
	Value []Company `json:"VAL"`
}

func main() {
	start := time.Now()
	conn, err := grpc.Dial("10.0.242.9:7001", grpc.WithInsecure())
	if err != nil {
		log.Fatal("client could connect to grpc service:", err)
	}
	defer conn.Close()
	c := pb.NewTestServiceClient(conn)

	request := &pb.Request{
		FileName: "arquivo.txt",
	}
	fileStreamResponse, err := c.Download(context.TODO(), request)
	if err != nil {
		log.Println("error downloading:", err)
		return
	}

	b := new(bytes.Buffer)

	for {
		size, err := fileStreamResponse.Recv()

		if err == io.EOF {
			log.Println("received all chunks")

			log.Println("creating file with chunks receive")
			ioutil.WriteFile(request.FileName, b.Bytes(), 0755|os.ModeAppend)
			log.Println("created")

			break
		}
		if err != nil {
			log.Println("err receiving chunk:", err)
			break
		}

		log.Printf("processo: %v de %v ", size.GetProccess()+1, size.GetTotal())
		b.Write(size.FileChunk)
	}

	log.Print("\nProccess time: ", time.Since(start))
}
