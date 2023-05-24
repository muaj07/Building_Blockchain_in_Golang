package main 
import
(
"github.com/muaj07/transport/network"
"github.com/muaj07/transport/core"
"github.com/muaj07/transport/crypto"
"github.com/go-kit/log"
"time"
"math/rand"
"bytes"
"fmt"
"strconv"
)

var Logger log.Logger

func main() {
trLocal := network.NewLocalTransport("LOCAL")
trRemoteA := network.NewLocalTransport("REMOTE_A")
trRemoteB := network.NewLocalTransport("REMOTE_B")
trRemoteC := network.NewLocalTransport("REMOTE_C")

trLocal.Connect(trRemoteA)
trLocal.Connect(trRemoteB)
trLocal.Connect(trRemoteC)

trRemoteA.Connect(trLocal)
trRemoteA.Connect(trRemoteB)
trRemoteA.Connect(trRemoteC)

trRemoteB.Connect(trLocal)
trRemoteB.Connect(trRemoteA)
trRemoteB.Connect(trRemoteC)

trRemoteC.Connect(trLocal)
trRemoteC.Connect(trRemoteA)
trRemoteC.Connect(trRemoteB)


//List of remote servers
initRemoteServers([]network.Transport{trRemoteA, trRemoteB, trRemoteC})

go func(){
	for {
		if err := SendTransaction(trRemoteA, trLocal.Addr()); err!=nil {
			Logger.Log(
				"Error", err,
			)
		}
		time.Sleep(2 * time.Second)
	}
}()
// out-of-sync Server
go func(){
	time.Sleep(9 * time.Second)
	trLate := network.NewLocalTransport("REMOTE_LATE")
	trRemoteC.Connect(trLate)
	lateServer := makeServer(string(trLate.Addr()), trLate, nil)
	go lateServer.Start()

}()



privKey := crypto.GeneratePrivateKey()
//configure the local server options
localServer := makeServer("LOCAL", trLocal, &privKey)
localServer.Start()
}

// configure and start the remote servers
 func initRemoteServers(trs []network.Transport) {
	for i := 0; i<len(trs); i++ {
		id := fmt.Sprint("REMOTE_%d", i)
		s := makeServer(id, trs[i], nil)
		go s.Start()
	}
 }

func makeServer(id string, tr network.Transport, pk *crypto.PrivateKey) *network.Server{
	opts := network.ServerOpts{
		ID: id,
		Transports: []network.Transport{tr},
		PrivateKey: pk,
	}
	s, err := network.NewServer(opts)
	if err != nil {
		Logger.Log(
			"Error", err,
		)
	}
	return s
}

func SendTransaction(tr network.Transport, to network.NetAddr) error {
	privKey := crypto.GeneratePrivateKey()
	// Generate a random number (based 10) and convert it to string
	data:= []byte(strconv.FormatInt(int64(rand.Intn(1000000)), 10))
	tx := core.NewTransaction(data)
	tx.SetFirstSeen(time.Now().UnixNano())
	tx.Sign(privKey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err !=nil {
		return err
	}
	// NewMessage in the rpc.go file of the network
	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())
	// SendMessage is inside local_transport.go file of the Transport folder
	// msg.Bytes() is a method in rpc.go file
	return tr.SendMessage(to, msg.Bytes())
}