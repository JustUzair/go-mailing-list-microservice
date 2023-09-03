package main

import (
	"context"
	"log"
	pb "mailinglist/proto"
	"time"

	"github.com/alexflint/go-arg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func logResponse(res *pb.EmailResponse, err error) {
	if err != nil {
		log.Fatalf("\terror : %v", err)
	}

	if res.EmailEntry == nil {
		log.Printf("\t email not found\n")
	} else {
		log.Printf("\t response : %v\n", res.EmailEntry)

	}
}

func createEmail(client pb.MailingListServiceClient, addr string) *pb.EmailEntry {
	log.Println("\t client : create email")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := client.CreateEmail(ctx, &pb.CreateEmailRequest{EmailAddr: addr})
	logResponse(res, err)
	return res.EmailEntry
}

func getEmail(client pb.MailingListServiceClient, addr string) *pb.EmailEntry {
	log.Println("\t client : get email")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := client.GetEmail(ctx, &pb.GetEmailRequest{EmailAddr: addr})
	logResponse(res, err)
	return res.EmailEntry
}

func getEmailBatch(client pb.MailingListServiceClient, count, page int) {
	log.Println("\t client : get email batch")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := client.GetEmailBatch(ctx, &pb.GetEmailBatchRequest{Page: int32(page), Count: int32(count)})
	if err != nil {
		log.Fatalf("\terror : %v", err)
	}
	log.Println("response : ")
	for i := 0; i < len(res.EmailEntries); i++ {
		log.Printf("\t item [%v of %v] : %s\n", i+1, len(res.EmailEntries), res.EmailEntries[i])
	}
}

func updateEmail(client pb.MailingListServiceClient, entry pb.EmailEntry) *pb.EmailEntry {
	log.Println("\t client : update email")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := client.UpdateEmail(ctx, &pb.UpdateEmailRequest{EmailEntry: &entry})
	logResponse(res, err)
	return res.EmailEntry
}

func deleteEmail(client pb.MailingListServiceClient, addr string) *pb.EmailEntry {
	log.Println("\t client : delete email")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := client.DeleteEmail(ctx, &pb.DeleteEmailRequest{EmailAddr: addr})
	logResponse(res, err)
	return res.EmailEntry
}

var args struct {
	GrpcAddr string `arg:"env:MAILINGLIST_GRPC_ADDR"`
}

func main() {
	arg.MustParse(&args)
	if args.GrpcAddr == "" {
		args.GrpcAddr = ":8081"
	}
	conn, err := grpc.Dial(args.GrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("client : could not connect: %v\n", err)
	}
	defer conn.Close()
	client := pb.NewMailingListServiceClient(conn)
	var confirmedAt int64 = 10000
	var optOut bool = true

	newEmail := createEmail(client, "pqr@yahoo.in")
	newEmail.ConfirmedAt = &confirmedAt
	newEmail.OptOut = &optOut

	updateEmail(client, *newEmail)
	deleteEmail(client, *newEmail.Email)
	getEmailBatch(client, 5, 1)
}
