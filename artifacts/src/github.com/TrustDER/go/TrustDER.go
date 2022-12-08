package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"
	 "math"
	 "math/big"
	 "io"
	 "os"
	 "encoding/base64"
	"net"
	"bufio"
	 "golang.org/x/crypto/ed25519"
	 crand "crypto/rand"
	 "math/rand"
	"encoding/pem"
	"github.com/mikesmitty/edkey"
	"golang.org/x/crypto/nacl/secretbox"
	 "golang.org/x/crypto/ssh"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "1234"
	SERVER_TYPE = "tcp"
	order       = 1000
)

// SmartContract Define the Smart Contract structure
type SmartContract struct {
}

// Device :  Define the Device structure, with 2 properties.  Structure tags are used by encoding/json library
type Device struct {
	ID   string `json:"id"`
	VerificationStatus  string `json:"verificationstatus"`
}

// Init ;  Method for initializing smart contract
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

var logger = flogging.MustGetLogger("TrustDER_cc")

// Invoke :  Method for INVOKING smart contract
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	logger.Infof("Function name is:  %d", function)
	logger.Infof("Args length is : %d", len(args))

	switch function {
	case "test":
		return s.test(APIstub, args)
	case "queryDevice":
		return s.queryDevice(APIstub, args)
	case "initLedger":
		return s.initLedger(APIstub)
	case "createDevice":
		return s.createDevice(APIstub, args)
	case "queryAllDevices":
		return s.queryAllDevices(APIstub)
	case "changeDeviceID":
		return s.changeDeviceID(APIstub, args)
	case "getHistoryForAsset":
		return s.getHistoryForAsset(APIstub, args)
	case "queryDevicesByID":
		return s.queryDevicesByID(APIstub, args)
	case "provisionID":
		return s.provisionID(APIstub)
	case "keymaker":
		return s.keymaker(APIstub)
	case "sendShard":
		return s.sendShard(APIstub, args)
	case "listenShard":
		return s.listenShard(APIstub)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}

func (s *SmartContract) queryDevice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	DeviceAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(DeviceAsBytes)
}

func (s *SmartContract) test(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	DeviceAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(DeviceAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	devices := []Device{
		Device{ID: "01", VerificationStatus: "1"},
		Device{ID: "02", VerificationStatus: "1"},
		Device{ID: "03", VerificationStatus: "1"},
		Device{ID: "04", VerificationStatus: "1"},
	}

	i := 0
	for i < len(devices) {
		deviceAsBytes, _ := json.Marshal(devices[i])
		APIstub.PutState("DEVICE"+strconv.Itoa(i), deviceAsBytes)
		i = i + 1
	}
	return shim.Success(nil)
}

func (s *SmartContract) createDevice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var device = Device{ID: args[1], VerificationStatus: args[2]}

	deviceAsBytes, _ := json.Marshal(device)
	APIstub.PutState(args[0], deviceAsBytes)

	indexName := "id~key"
	NameIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{device.ID, args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	APIstub.PutState(NameIndexKey, value)

	return shim.Success(deviceAsBytes)
}

func (S *SmartContract) queryDevicesByID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments")
	}
	id_device := args[0]

	IdResultIterator, err := APIstub.GetStateByPartialCompositeKey("id~key", []string{id_device})
	if err != nil {
		return shim.Error(err.Error())
	}

	defer IdResultIterator.Close()

	var i int
	var id string

	var devices []byte
	bArrayMemberAlreadyWritten := false

	devices = append([]byte("["))

	for i = 0; IdResultIterator.HasNext(); i++ {
		responseRange, err := IdResultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}

		id = compositeKeyParts[1]
		assetAsBytes, err := APIstub.GetState(id)

		if bArrayMemberAlreadyWritten == true {
			newBytes := append([]byte(","), assetAsBytes...)
			devices = append(devices, newBytes...)

		} else {
			devices = append(devices, assetAsBytes...)
		}

		fmt.Printf("Found a asset for index : %s asset id : ", objectType, compositeKeyParts[0], compositeKeyParts[1])
		bArrayMemberAlreadyWritten = true

	}

	devices = append(devices, []byte("]")...)

	return shim.Success(devices)
}

func (s *SmartContract) queryAllDevices(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "DEVICE0"
	endKey := "DEVICE999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllDevices:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeDeviceID (APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	deviceAsBytes, _ := APIstub.GetState(args[0])
	device := Device{}

	json.Unmarshal(deviceAsBytes, &device)
	device.ID = args[1]

	deviceAsBytes, _ = json.Marshal(device)
	APIstub.PutState(args[0], deviceAsBytes)

	return shim.Success(deviceAsBytes)
}

func (t *SmartContract) getHistoryForAsset(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	deviceName := args[0]

	resultsIterator, err := stub.GetHistoryForKey(deviceName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForAsset returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) provisionID(APIstub shim.ChaincodeStubInterface) sc.Response{
	pubKey, privKey, _ := ed25519.GenerateKey(crand.Reader)
	publicKey, _ := ssh.NewPublicKey(pubKey)

	pemKey := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(privKey), // <- marshals ed25519 correctly
	}
	privateKey := pem.EncodeToMemory(pemKey)
	authorizedKey := ssh.MarshalAuthorizedKey(publicKey)

	f, err := os.Create("/tmp/testID")
	if err!=nil{
		return shim.Error(err.Error())
	}
	defer f.Close()

	_, err = f.Write(privateKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	f2, err := os.Create("/tmp/testID.pub")
	if err!=nil{
		return shim.Error(err.Error())
	}
	defer f2.Close()

	_, err = f2.Write(authorizedKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (s *SmartContract) keymaker(APIstub shim.ChaincodeStubInterface) sc.Response {

	rand.Seed(time.Now().UnixNano())
	randomKey := rand.Int63n(order * 10)

	poly := genEquation(randomKey, 2)
	shardX, shardY := genShards(poly)

	encoded, err := encryptKey(randomKey)

	if err != nil {
		shim.Error(err.Error())
	}
	shards := distShards(shardX, shardY, encoded)

	f, err := os.Create("/tmp/data.txt")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer f.Close()

	for _, shard := range shards {
		n, err := f.Write(shard)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Printf("wrote %d bytes\n", n)
	}

	result := "Successfully stored shards at " + f.Name()
	return shim.Success([]byte(result))
}

func genEquation(randomKey int64, k int64) []int64 {
	rand.Seed(time.Now().UnixNano())
	var poly []int64
	for i := int64(0); i < k; i++ {
		// poly = append(poly, rand.Int63()%order)
		poly = append(poly, rand.Int63n(order/10))
	}
	poly = append(poly, randomKey)
	return poly
}

func genShards(poly []int64) ([]float64, []float64) {
	rand.Seed(time.Now().UnixNano())
	var dist []float64
	// x := rand.Int63() % order
	x := rand.Int63n(order)
	for i := 0; i < 5; i++ {
		dist = append(dist, float64(x))
		// x = rand.Int63() % order
		x = rand.Int63n(order)
	}
	y := solveWithX(poly, dist)

	return dist, y
}

func solveWithX(poly []int64, x []float64) []float64 {
	var yvals []float64
	for _, v := range x {
		y := float64(0)
		y = float64(poly[0])*math.Pow(v, 2) + float64(poly[1])*v + float64(poly[2])

		yvals = append(yvals, y)
	}

	return yvals
}

func encryptKey(randomKey int64) (string, error) {
	plaintext, err := os.ReadFile("/tmp/testID")
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}

	key := big.NewInt(randomKey)
	var keybuf [32]byte
	key.FillBytes(keybuf[:])

	var nonce [24]byte
	if _, err := io.ReadFull(crand.Reader, nonce[:]); err != nil {
		return "", fmt.Errorf("%s", err)
	}
	encrypted := secretbox.Seal(nonce[:], plaintext, &nonce, &keybuf)
	encoded := base64.StdEncoding.EncodeToString(encrypted)

	return encoded, nil
}

func distShards(shardsX, shardsY []float64, cipher string) [][]byte {
	var shards [][]byte
	cipher = cipher + "\n"
	for i, _ := range shardsX {
		x := big.NewInt(int64(shardsX[i]))
		var xbuf [8]byte
		x.FillBytes(xbuf[:])
		y := big.NewInt(int64(shardsY[i]))
		var ybuf [8]byte
		y.FillBytes(ybuf[:])
		var buf []byte
		buf = append(buf, xbuf[:]...)
		buf = append(buf, ybuf[:]...)
		buf = append(buf, cipher...)
		shards = append(shards, buf)
	}

	return shards
}

func (s *SmartContract) listenShard(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Start listening on port " + SERVER_PORT + "...")

	// listen on port 8000
	ln, err := net.Listen(SERVER_TYPE, ":"+SERVER_PORT)
	if err != nil {
		return shim.Error(err.Error())
	}
	// accept connection
	conn, err := ln.Accept()
	if err != nil {
		return shim.Error(err.Error())
	}
	// run loop forever (or until ctrl-c)
	//for {
	// get message, output
	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		return shim.Error(err.Error())
	}
	//}
	f, err := os.Create("/tmp/shard.txt")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer f.Close()
	f.Write(buffer)
	conn.Close()
	result := "Shard written to file " + f.Name()
	return shim.Success([]byte(result))
}

func (s *SmartContract) sendShard(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Please send the listening server name!")
	}

	SERVER_NAME := args[0]

	fp, err := os.Open("/tmp/data.txt")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer fp.Close()
	fileScanner := bufio.NewScanner(fp)
	fileScanner.Split(bufio.ScanLines)
	var text []string
	for fileScanner.Scan() {
		text = append(text, fileScanner.Text())
	}
	//for _, reader := range text {
	//	// what to send?
	//	shard := []byte(reader)
	//	x := big.NewInt(0)
	//	y := big.NewInt(0)
	//	x = x.SetBytes(shard[:8])
	//	y = y.SetBytes(shard[8:16])
	//	fmt.Println(x.String() + " " + y.String())
	//}
	conn, err := net.Dial(SERVER_TYPE, SERVER_NAME+":"+SERVER_PORT)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer conn.Close()
	//for _, reader := range text {
	//	// what to send?
	//	fmt.Print("Text to send: ")
	//	// send to server
	//	conn.Write([]byte(reader + "\n"))
	//}
	conn.Write([]byte(text[0] + "\n"))
	return shim.Success([]byte("Shards sent"))
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
