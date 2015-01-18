## Overview
This is the server implementation to match the [Client Library](http://github.com/conorbrady/dfs-client).

## Setup Instructions
1. Setup Go >=1.3.3
	1. [Install Go](https://golang.org/doc/install)
	2. Navigate to a directory that will contain your Go directory
	3. Create the folder hierarchy

			mkdir -p go/{src,bin,pkg}
	4. Set this directory to be your GOPATH environment variable from 3.

			export GOPATH="Your/go/directory"
	5. Export go/bin to your bath PATH

			export PATH=$PATH:$GOPATH/bin

4. Download the code

		go get github.com/conorbrady/distributed-file-system

## Running

	go install github.com/conorbrady/distributed-file-system
### Private Key
All server instances must share a 32 byte private key for security between the
system, this should be included as a 32 byte binary file as a command line
argument. An example of how to generate such a key is:

	dd if=/dev/urandom of=private.key bs=1 count=32
### Working Directory
Each instance of the server must run in its own working directory, each will
generate:

* **uuid.txt** - Used for persistant identification

* **out.log** - A log of the servers operations

* **.sqlite** - A database used by each server for its own purposes

* In addition the Fileserver will generate **storage/** for storing files

### Authentication Server

	distributed-file-system -port <address> -mode AS -key path/to/private.key

We are presented with:
```
Please Select:
1. View Users
2. Add User
3. Delete User
```

We must add at least one user to the system for the server to authenticate any
requests. Users persist from between launches in **auth.sqlite**. The server is
always listening for requests when running regardless of the state of the prompt.

### Directory Server

	distributed-file-system -port <address> -mode DS -key path/to/private.key

We are presented with:
```
Please Select:
1. View File Servers
2. Add File Server
```

The file servers must be alive on the network when they are added. The Directory
server will ping them for identification. If they cannot be found it will give
an error and not add them to the list. Files get randomly assigned to the file
servers when a new one is requested. The file servers and the files they contain
are saved in **locate.sqlite**.

### File Server

	distributed-file-system -port <address> -mode FS -key path/to/private.key

Unlike the directory and authentication servers, where only one can be specified,
many file servers can be run up simultaneously. Each must be registered with the
Directory Server and it will distribute out the files accordingly across the file
servers. File servers have to prompt and just run as processing responding to read/write
requests.
