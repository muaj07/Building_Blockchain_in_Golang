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

# Building Blockchain in Golang

This repository contains a simple implementation of a blockchain network using Golang. The main focus of this project is to understand the basic concepts of blockchain and how they can be implemented using the Go programming language.

## Features

- Generation of private keys for securing data.
- Creation of a local node server.
- Communication between nodes in the network.
- Sending transactions over the network.
- TCP connection testing.

## Getting Started

### Prerequisites

- Go (version 1.15 or later)

### Dependencies

- github.com/go-kit/log
- github.com/muaj07/transport/core
- github.com/muaj07/transport/crypto
- github.com/muaj07/transport/network

### Running the Code

1. Clone the repository:

   ```sh
   git clone https://github.com/muaj07/Building_Blockchain_in_Golang.git
   cd Building_Blockchain_in_Golang

2. Run the main.go file:
   
   go run main.go
   
This will start a local node server and two remote nodes. The local node communicates with the remote nodes over the network.

### Code Structure
**main.go**: This is the entry point of the program. It contains the main function which initializes the local node server and remote nodes. It also contains helper functions for creating servers and testing TCP connections.

__go.mod and go.sum__: These files are used to handle dependencies.

**Makefile**: This file can be used to automate the build process.

### Contributing
Contributions are welcome! Feel free to fork the repository and submit pull requests.

### License
This project is open source and available under the MIT License.

### Contact
GitHub: @muaj07

### Acknowledgments
- The Go Programming Language
- Blockchain Technology

