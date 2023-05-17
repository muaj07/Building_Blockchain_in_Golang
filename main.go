package main 
import

(
"github.com/muaj07/transport/network"
"time"
)

func main() {
trLocal := network.NewLocalTransport("LOCAL")
trRemote := network.NewLocalTransport("REMOTE")
trLocal.Connect(trRemote)
trRemote.Connect(trLocal)

go func()	{
	for {
		trRemote.SendMessage(trLocal.Addr(), []byte("Hello Enoda"))
		time.Sleep(3 * time.Second)
	}
}()

opts := network.ServerOpts {
	Transports: []network.Transport{trLocal},
}
s := network.NewServer(opts)
s.Start()
}