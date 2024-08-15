package main

import (
	"encoding/hex"
	"fmt"
	"github.com/theQRL/go-qrllib/common"
	"github.com/theQRL/go-qrllib/xmss"
	"github.com/theQRL/walletd-rest-proxy/generated"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func main() {
	opts := []grpc.DialOption{grpc.WithInsecure()}
	dial, err := grpc.Dial("mainnet-3.automated.theqrl.org:19009", opts...)
	if err != nil {
		log.Println("dial err , the err is ", err)
		return
	}
	ctx := context.Background()
	client := generated.NewPublicAPIClient(dial)
	seed, _ := hex.DecodeString("eddc73f97c74acc797a4f113bcbc34475ddeac43e0db363e92e72153bb99f5cda11fefb042e4a985dfff6b52b1528e5d")
	seedBytes := seed[:]

	// Create an array of fixed size
	var uint8Array [common.SeedSize]uint8

	// Copy the slice into the array
	copy(uint8Array[:], seedBytes)

	xmss := xmss.NewXMSSFromSeed(uint8Array, 4, xmss.SHAKE_128, common.SHA256_2X)

	pk := xmss.GetPK()
	log.Println("len(pk ) is ", len(pk), " the pk is ", pk, " the pk hex is ", hex.EncodeToString(pk[:]))
	//expectedAddress := "0
	coins, err := client.TransferCoins(ctx, &generated.TransferCoinsReq{
		MasterAddr:  GetAddressByte("Q010200b70a623213a24112766d6a85648ddffa347373c58bdfee02619d25505c0c99057cc7d4e5"),
		AddressesTo: [][]byte{GetAddressByte("Q01020036c064daeacb6c7705fd98e833cf4542e5b488251e5446a3bdf43ec16a44cd9395fd7f9c")},
		Amounts:     []uint64{100000},
		Fee:         100000,
		XmssPk:      pk[:],
	})
	address := xmss.GetLegacyAddress()

	fmt.Println("the address is ", hex.EncodeToString(address[:]))

	if err != nil {
		log.Println("transfer coin err", err)
		return
	}
	log.Println("the coins is ", coins)
}
