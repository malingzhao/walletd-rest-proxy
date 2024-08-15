package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
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
	decodeString, _ := hex.DecodeString("010200b70a623213a24112766d6a85648ddffa347373c58bdfee02619d25505c0c99057cc7d4e5")
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
		log.Println("get balance err ", err)
		return
	}
	fmt.Println("the balance is ", balance)

	height, err := client.GetHeight(ctx, &generated.GetHeightReq{})
	if err != nil {
		log.Println("get height err ", err)
		return
	}

	fmt.Println("the height is ", height)

	seed, _ := hex.DecodeString("")
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
		//MasterAddr:  GetAddressByte("Q010200b70a623213a24112766d6a85648ddffa347373c58bdfee02619d25505c0c99057cc7d4e5"),
		AddressesTo: [][]byte{GetAddressByte("Q01020036c064daeacb6c7705fd98e833cf4542e5b488251e5446a3bdf43ec16a44cd9395fd7f9c")},
		Amounts:     []uint64{1000000},
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
	tx := coins.ExtendedTransactionUnsigned.GetTx()
	dataBytes := make([]byte, 0)
	dataBytes = append(dataBytes, tx.GetMasterAddr()...)
	dataBytes = append(dataBytes, Uint64ToBytes(tx.GetFee())...)
	dataBytes = append(dataBytes, tx.GetTransfer().GetMessageData()...)
	for index, _ := range tx.GetTransfer().AddrsTo {
		dataBytes = append(dataBytes, tx.GetTransfer().GetAddrsTo()[index]...)
		dataBytes = append(dataBytes, Uint64ToBytes(tx.GetTransfer().GetAmounts()[index])...)
	}

	sum256 := sha256.Sum256(dataBytes)
	sign, err := xmss.Sign(sum256[:])
	if err != nil {
		log.Println("the sign err is ", err)
		return
	}

	tx.Signature = sign

	transaction, err := client.PushTransaction(context.TODO(), &generated.PushTransactionReq{TransactionSigned: tx})
	if err != nil {
		log.Println("push tx err ", err)
		return
	}
	log.Println("the tx result is ", transaction)
	log.Println("the tx hash is ", hex.EncodeToString(transaction.GetTxHash()))
	_, err = client.GetTransaction(ctx, &generated.GetTransactionReq{TxHash: transaction.TxHash})

	if err != nil {
		log.Println("get tx err", err)
		return
	}
	//log.Println("the newTx is ", newTx)
}

func GetAddressByte(address string) []byte {
	address = address[1:]
	decodeString, _ := hex.DecodeString(address)
	return decodeString
}

func Uint64ToBytes(fee uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, fee)
	return bytes
}
