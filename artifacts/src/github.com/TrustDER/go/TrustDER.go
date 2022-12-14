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
	// Port for the shard distribution. This is for the step where a
	// peer distributes its shards to all the other trusted peers
	SERVER_PORT = "1234"
	// Port for getting shards to the verifier from k trusted peers
	// and the peer that needs to be verified
	SERVER_PORT2 = "1235"
	// Server type is tcp as we're opening up a TCP socket
	SERVER_TYPE = "tcp"
	// The domain should be changed to match the name of the chaincode
	// Shoud follow syntax: ".org1.secretidentity.com-<chaincode_name>_1"
	SERVER_DOMAIN = ".org1.secretidentity.com-secureID_1"
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
	case "sendSingleShard":
		return s.sendSingleShard(APIstub)
	case "getShardsK":
		return s.getShardsK(APIstub)
	case "keychecker":
		return s.keychecker(APIstub, args)
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

// Smart contract that provisions an identity for a peer. It generates an
// ed25519 private/public key pair and stores them in /tmp
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

// Smart contract that implements the keymaker protocol. It generates a
// random key, generates an equation using the random key and generates
// shards lying on the equation. It then encrypts the private key of the
// peer with the random key and appends that to the generated shards
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

// Helper function that takes the random key and k as arguments and
// generates a polynomial equation of degree k with the random key
// as the intercept.
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

// Helper function that takes a slice containing polynomial constants
// and returns two slices representing (x,y) points lying on the line
// formed by the constants.
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

// Helper function that takes a slice containing polynimial constants
// and a slice containing x values as input and returns a slice
// containing values obtained from solving the equation using each x
func solveWithX(poly []int64, x []float64) []float64 {
	var yvals []float64
	for _, v := range x {
		y := float64(0)
		y = float64(poly[0])*math.Pow(v, 2) + float64(poly[1])*v + float64(poly[2])

		yvals = append(yvals, y)
	}

	return yvals
}

// Helper function that takes the random key as input and encrypts the
// private key/ID of the peer and returns the cipher to the caller
func encryptKey(randomKey int64) (string, error) {
	plaintext, err := os.ReadFile("/tmp/testID")
	if err != nil {
		return shim.Error(err)
	}

	key := big.NewInt(randomKey)
	var keybuf [32]byte
	key.FillBytes(keybuf[:])

	var nonce [24]byte
	if _, err := io.ReadFull(crand.Reader, nonce[:]); err != nil {
		return shim.Error(err)
	}
	encrypted := secretbox.Seal(nonce[:], plaintext, &nonce, &keybuf)
	encoded := base64.StdEncoding.EncodeToString(encrypted)

	return encoded, nil
}

// Helper function that takes generated x and y slices and the cipher
// text as input. It converts the (x,y) points to a 16-byte (8,8)
// representation and appends the ciphertext to it. It then performs
// a base64 encode on the obtained byte-string and returns the resulting
// slice of byte slices as the shards to be stored
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

// Smart contract to receive a single shard from a peer. It stores the
// received shards at /tmp/shard.txt
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

// Smart contract to distribute the shards stored on a peer to all
// other trusted peers in the network including itself
func (s *SmartContract) sendShard(APIstub shim.ChaincodeStubInterface) sc.Response {

	// List of names of trusted peers and the peer itself
	peers := [...]string{"dev-peer4", "dev-peer3", "dev-peer2", "dev-peer1", "dev-peer0"}

	// Reading from the file containing all shards of the peer
	fp, err := os.Open("/tmp/data.txt")
	if err != nil {
		return shim.Error(err)
	}
	defer fp.Close()
	fileScanner := bufio.NewScanner(fp)
	fileScanner.Split(bufio.ScanLines)
	var text []string
	// Reading line by line and storing them in a slice
	for fileScanner.Scan() {
		text = append(text, fileScanner.Text())
	}

	// Making a slice of Connection objects to connect to all the trusted peers
	var conns = make([]net.Conn, 5)
	// Dialing to each peer
	for i, peer := range peers {
		conn, err := net.Dial(SERVER_TYPE, peer+SERVER_DOMAIN+":"+SERVER_PORT)
		if err != nil {
			return shim.Error(err)
		}
		conns[i] = conn
	}

	// Sending one shard to each peer where peer4 receives shard0 and so on
	for i, conn := range conns {
		_, err = conn.Write([]byte(text[i] + "\n"))
		if err != nil {
			return shim.Error(err)
		}
		err = conn.Close()
		if err != nil {
			return shim.Error(err)
		}
	}
	return shim.Success([]byte("Shards sent"))
}

// Smart contract to receive shards sent by k trusted peers and by
// the device to be verified. The shards are stored as shard[i].txt
func (s *SmartContract) getShardsK(APIstub shim.ChaincodeStubInterface) sc.Response {
	var shards []string

	// listen on port SERVER_PORT2
	ln, err := net.Listen(SERVER_TYPE, ":"+SERVER_PORT2)
	if err != nil {
		return shim.Error(err)
	}
	defer ln.Close()
	// Receive 4 shards from 4 peers
	for i := 0; i < 4; i++ {
		// accept connection
		conn, err := ln.Accept()
		if err != nil {
			return shim.Error(err)
		}
		// Buffer to read the shard
		buffer := make([]byte, 1024)
		_, err = conn.Read(buffer)
		if err != nil {
			return shim.Error(err)
		}
		// Forming the file name
		fname := "/tmp/shard" + string(i+48) + ".txt"
		shards = append(shards, string(buffer[:]))
		shards = append(shards, "\n")
		// Creating a file with the formed name
		f, err := os.Create(fname)
		if err != nil {
			return shim.Error(err)
		}
		// Writing to the file
		f.Write(buffer)
		f.Close()
		conn.Close()
	}
	result := "Shards for verification received"
	return shim.Success([]byte(result))
}

// Smart contract to send a single shard to a verifier. It is run on
// k trusted peers to send their shards to the verifier peer
func (s *SmartContract) sendSingleShard(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Name of the verifier peer
	verifier := "dev-peer1"

	// Reading from the file containing the shard on the peer
	fp, err := os.Open("/tmp/shard.txt")
	if err != nil {
		return shim.Error(err)
	}
	defer fp.Close()
	fileScanner := bufio.NewScanner(fp)
	fileScanner.Split(bufio.ScanLines)
	var text []string
	// Reading the first line of the file as this is the line
	// containing the shard
	for fileScanner.Scan() {
		text = append(text, fileScanner.Text())
	}
	// Dialing to the verifier shard
	conn, err := net.Dial(SERVER_TYPE, verifier+SERVER_DOMAIN+":"+SERVER_PORT2)
	if err != nil {
		return shim.Error(err)
	}
	// Sending the shard to the verifier
	_, err = conn.Write([]byte(text[0] + "\n"))
	if err != nil {
		return shim.Error(err)
	}
	err = conn.Close()
	if err != nil {
		return shim.Error(err)
	}
	return shim.Success([]byte("Shards sent"))
}

// Smart contract to decode the shards and perform the verification of
// an untrusted peer on the verifier using Lagrange interpolation formula
func (s *SmartContract) keychecker(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	// Read and decode all the shards
	x_shards, y_shards, err := decodeShards()
	if err != nil {
		return shim.Error(err)
	}
	var result string
	// Perform the verification and get appropriate string depending on the result
	if verify(x_shards, y_shards) {
		result = "verified"
	} else {
		result = "not verified"
	}
	// Modify the ledger with the result
	err = stub.PutState(args[0], []byte(result))
	if err != nil {
		return shim.Error(err)
	}
	return shim.Success([]byte(result))
}

// Helper function to read all the shards residing on the verifier and
// decode them to get (x,y) points and return 2 slices, one containing
// x values and the other containing y values
func decodeShards() ([]float64, []float64, error) {
	var shards = make([]string, 4)

	var x_shards = make([]float64, 4)
	var y_shards = make([]float64, 4)
	for i := 0; i < 4; i++ {
		// Read shard from file
		fname := "/tmp/shard" + string(i+48) + ".txt"
		fp, err := os.Open(fname)
		if err != nil {
			return shim.Error(err)
		}
		fileScanner := bufio.NewScanner(fp)
		var text []string
		fileScanner.Split(bufio.ScanLines)
		for fileScanner.Scan() {
			text = append(text, fileScanner.Text())
		}
		// Store shard in slice
		shards[i] = text[0]
		err = fp.Close()
		if err != nil {
			return shim.Error(err)
		}
	}
	for i, shard := range shards {
		// Decode shard from base64
		decoded, err := base64.StdEncoding.DecodeString(shard)
		if err != nil {
			return shim.Error(err)
		}
		x := big.NewInt(0)
		y := big.NewInt(0)
		// Extract first 8 bytes of shard as x
		x = x.SetBytes(decoded[:8])
		// Extract bytes 8 to 16 of shard as y
		y = y.SetBytes(decoded[8:16])
		x_shards[i] = float64(x.Int64())
		y_shards[i] = float64(y.Int64())
	}
	return x_shards, y_shards, nil
}

// Helper function that takes 2 slices containing x values and y
// values as input, performs the verification for the unverified
// peer and returns a Boolean output
func verify(shardX []float64, shardY []float64) bool {
	// Set the last shard in the slice as the shard of the unverified peer
	device_x := shardX[len(shardX)-1]
	device_y := math.Round(shardY[len(shardY)-1])

	// Remove above shard from the slices
	shardX = shardX[:len(shardX)-1]
	shardY = shardY[:len(shardY)-1]

	// Calculate interpolated y value of the unverified peer's x using
	// Lagrange interpolation with the other shards
	interpolatedY := lagrangePoly(device_x, shardX, shardY)

	// Check if interpolated y value matches the shard y value
	return device_y == interpolatedY
}

// Helper function that implements Lagrange interpolation formula. It takes an
// x value and 2 slices, one containing x values and the other containing y values
// as input. It implements Lagrange interpolation on the x value using the (x,y)
// values and returns the resulting y value
func lagrangePoly(device_x float64, shardsX []float64, shardsY []float64) float64 {
	sum := float64(0)
	for i, v := range shardsY {
		prod := float64(1)
		for j, x := range shardsX {
			if i != j {
				num := device_x - x
				den := shardsX[i] - x
				div := num / den
				prod = prod * div
			}
		}
		sum = sum + v*prod
	}

	return math.Round(sum)
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
