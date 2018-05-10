srv
===

"srv" is an HTTP server for helping development which can serve static files or 
prints incoming requests.

"srv" supports TLS (https) connections and comes with a self signed TLS cerificate
generator based on [a generator](1) in golang respository

Only tested on Linux but releases page have builds on different systems and 
architectures that go supports.

Note that this command should not be used in production.

__Features:__
- Static file server
- A server for printing incoming requests
- TLS (https) support
- Self-signed X.509 certificate file generator for TLS server

Installation
------------

Download prebuild executables:

- Download the latest executable from [Github Releases](2) page. 
- Unzip it. 

If you have Go installed, you may use go get:

    go get github.com/ozgio/srv

Usage
-----

For general help

    srv --help

#### Global Flags:

    -c, --cert string   Path to cert file for https server
    -o, --host string   Host name or address (default "127.0.0.1")
    -k, --key string    Path to key file for https server
    -p, --port int      Port to listen (default 8010)


### srv files

Static file server with directory listing

#### Usage:
  srv files [flags]

#### Flags:

    -r, --root string   Root path for server (default "./")
    -h, --help          help for files

#### Examples:

    #serve the files at current working directory
    srv files

    #serve the files at ~/files using port 80
    srv files --port=80 --root=~/files

    #create https server using specified pem files
    srv files --port=443 --cert=cert.pem --key=key.pem


### srv mirror

Prints incoming request in plain text or json format

#### Usage:

    srv mirror [flags]

#### Flags:

    -h, --help   help for mirror
  

#### Examples:

    #start server
    srv mirror
    
    #json responses
    srv mirror --port=80 --json


### srv generate

Generates key and cert files for https server. Keep in mind that these are
only meant for development. These files can be used with "srv files" or
"srv mirror" commands later

#### Usage:

    srv generate [flags]

#### Flags:
  
    -h, --help   help for generate

#### Examples:

    srv generate --key=key.pem --cert=cert.pem

    #create https server using generated certificate
    srv files --cert=cert.pem --key=key.pem


Development
-----------

### Prerequisites:

- dep: https://golang.github.io/dep/
- make: https://www.gnu.org/software/make/

### Setup

    # Get dependencies
    make dep

For online code documentation see [Godocs](3).

TODO
-----
- CORS support
- Auto generate README from help commands

[1]: https://golang.org/src/crypto/tls/generate_cert.go
[2]: https://github.com/ozgio/srv/releases
[3]: https://godoc.org/github.com/ozgio/srv