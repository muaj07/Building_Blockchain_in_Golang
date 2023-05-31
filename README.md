# Building_Blockchain_in_Golang

This repository contains code for a basic blockchain written in Golang, which can be used for educational purposes.

## Requirements
> [makefile](https://www.gnu.org/software/make/manual/make.html)
> 
> [Go 1.18 or above](https://go.dev/)

## Running the blockchain
> Run the following command:
```
make run
```
## Details of the code

The **main.go** file contains example code for bootstrapping the blockchain. Here are the main steps for kicking-off the network:
> 1. use the __makeServer__ function to create a instant of server. For instance ***makeServer("LOCAL_NODE", &privKey, ":3000", []string{":4000"})***.
> This will instantiate **Server** instant, using the code in the __network__ package/folder.
> 2. start a _Goroutine_ of the new **Server** instant, for instance, ***go localNode.Start()***
> 3. __select{}__ block the main thread to keep the program running indefinitely.
>
> For more details about each folder/package, check the ***README.md*** in these folders/packages.

