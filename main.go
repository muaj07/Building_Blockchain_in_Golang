package main 
import
(
"github.com/muaj07/transport/network"
"github.com/muaj07/transport/core"
"github.com/muaj07/transport/crypto"
"github.com/go-kit/log"
"time"
"net"
//"math/rand"
"bytes"
//"fmt"
//"strconv"
)

var Logger log.Logger



// main is the entry point of the program
func main() {
    // Generate a new private key
    privKey := crypto.GeneratePrivateKey()

    // Create a server instance for the local node
    localNode := makeServer("LOCAL_NODE", &privKey, ":3000", []string{":4000"})

    // Start the local node server in a separate goroutine
    go localNode.Start()

	remoteNode := makeServer("REMOTE_NODE", nil, ":4000", []string{":5000"})
	go remoteNode.Start()

	remoteNodeB := makeServer("REMOTE_NODE_B", nil, ":5000", nil)
	go remoteNodeB.Start()

    // Block the main thread to keep the program running indefinitely
    select {}
}


// makeServer creates and returns a new network server with the specified ID, transport, and private key.
func makeServer(id string, privkey *crypto.PrivateKey, addr string, seedNodes []string) *network.Server{

    // Set the server options.
    opts := network.ServerOpts{
		SeedNodes: seedNodes,
		ListenAddr: addr, //source address
        ID:         id,
        PrivateKey: privkey,
    }

    // Create the new server.
    s, err := network.NewServer(opts)

    // Log any errors that occurred during server creation.
    if err != nil {
        Logger.Log(
            "Error", err,
        )
    }
    // Return the new server.
    return s
}


// tcpTester tests a TCP connection by dialing to port 3000 and sending a message.
func tcpTester() {
    // Dial TCP connection to port :3000
    conn, err := net.Dial("tcp", ":3000")
    if err != nil {
        panic(err)
    }

	privKey := crypto.GeneratePrivateKey()
	data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f}
	tx := core.NewTransaction(data)
	tx.SetFirstSeen(time.Now().UnixNano())
	tx.Sign(privKey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err !=nil {
		panic(err)
	}
	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	// Send "Hello there!" message to TCP server
    _, err = conn.Write(msg.Bytes())
    if err != nil {
        panic(err)
    }
}





// func main() {

// //List of Remote servers
// initRemoteServers(transports)
// localNode := transports[0]
// lateNode := transports[1]
// //remoteNodeA := transports[1]

// // go func(){
// // 	for {
// // 		if err := SendTransaction(remoteNodeA, localNode.Addr()); err!=nil {
// // 			Logger.Log(
// // 				"Error", err,
// // 			)
// // 		}
// // 		time.Sleep(2 * time.Second)
// // 	}
// // }()

// //out-of-sync Server
// go func(){
// 	time.Sleep(9 * time.Second)
// 	lateServer := makeServer(string(lateNode.Addr()), lateNode, nil)
// 	go lateServer.Start()
// }()



// privKey := crypto.GeneratePrivateKey()
// //configure the local server options
// localServer := makeServer("LOCAL", localNode, &privKey)
// localServer.Start()
// }

// // configure and start the remote servers
//  func initRemoteServers(trs []network.Transport) {
// 	for i := 0; i<len(trs); i++ {
// 		id := fmt.Sprintf("REMOTE_%d", i)
// 		s := makeServer(id, trs[i], nil)
// 		go s.Start()
// 	}
//  }


// // makeServer creates and returns a new network server with the specified ID, transport, and private key.
// func makeServer(id string, tr network.Transport, privkey *crypto.PrivateKey) *network.Server {

//     // Set the server options.
//     opts := network.ServerOpts{
// 		Transport: tr,
//         ID:         id,
//         Transports: transports,
//         PrivateKey: privkey,
//     }

//     // Create the new server.
//     s, err := network.NewServer(opts)

//     // Log any errors that occurred during server creation.
//     if err != nil {
//         Logger.Log(
//             "Error", err,
//         )
//     }

//     // Return the new server.
//     return s
// }


// func SendTransaction(tr network.Transport, to network.NetAddr) error{
// 	privKey := crypto.GeneratePrivateKey()
// 	tx := core.NewTransaction(contract())
// 	tx.SetFirstSeen(time.Now().UnixNano())
// 	tx.Sign(privKey)
// 	buf := &bytes.Buffer{}
// 	if err := tx.Encode(core.NewGobTxEncoder(buf)); err !=nil {
// 		return err
// 	}
// 	//NewMessage in the "rpc.go" file of the network
// 	//NewMessage contains the MessageType and Encoded data of txs
// 	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())
// 	// SendMessage is inside "local_transport.go" file of the Transport folder
// 	// msg.Bytes() is a method in "rpc.go" file that returns gob-encoded 
// 	// byte slice of the message (msg).
// 	return tr.SendMessage(to, msg.Bytes())
// }

// func contract() []byte{
// 	pushFoo := []byte{0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0xae}
// 	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0b, 0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0x0f}
// 	data = append(data, pushFoo...)
// 	return data
// }