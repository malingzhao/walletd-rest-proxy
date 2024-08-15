package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/theQRL/walletd-rest-proxy/generated"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func main() {
	decodeString, _ := hex.DecodeString("0106008009df1bff3fc861fbe1bf5e7a3a93d85092c0584b6bad6820969d5fe2fd3a59c1465b23")
	fmt.Println(base64.StdEncoding.EncodeToString(decodeString))
	opts := []grpc.DialOption{grpc.WithInsecure()}
	dial, err := grpc.Dial("mainnet-3.automated.theqrl.org:19009", opts...)
	if err != nil {
		log.Println("dial err , the err is ", err)
		return
	}

	ctx := context.Background()
	client := generated.NewPublicAPIClient(dial)

	balance, err := client.GetBalance(ctx, &generated.GetBalanceReq{Address: decodeString})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("the balance is ", balance)

}
