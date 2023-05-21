package main 
import

(
"github.com/muaj07/transport/network"
"github.com/muaj07/transport/core"
"github.com/muaj07/transport/crypto"
"time"
"math/rand"
"bytes"
"github.com/sirupsen/logrus"
"strconv"
)

func main() {
trLocal := network.NewLocalTransport("LOCAL")
trRemote := network.NewLocalTransport("REMOTE")
trLocal.Connect(trRemote)
trRemote.Connect(trLocal)

// start a goroutine that sends a message every 3 seconds
go func() {
	for {
		if err := SendTransaction(trRemote, trLocal.Addr()); err!= nil{
			logrus.Error(err)
		} 
		time.Sleep(3 * time.Second)
	}
}()

// configure server options
opts := network.ServerOpts{
	Transports: []network.Transport{trLocal},
}

// create and start the server
s := network.NewServer(opts)
s.Start()
}

func SendTransaction(tr network.Transport, to network.NetAddr) error {
	privKey := crypto.GeneratePrivateKey()
	data:= []byte(strconv.FormatInt(int64(rand.Intn(1000000)), 10))
	tx := core.NewTransaction(data)
	//tx.SetFirstSeen(time.Now().UnixNano())
	tx.Sign(privKey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err !=nil {
		return err
	}
	// NewMessage in the rpc.go file of the network
	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())
	// SendMessage is inside local_transport.go file of the Transport folder
	return tr.SendMessage(to, msg.Bytes())
}