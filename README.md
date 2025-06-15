# Snapcontrol

Snapcontrol is a little CLI tool to control snapserver. It aims to be a full implementation of the snapcast API.

## Usage
Note: `<id>` for clients and groups can be found in the output of `server status`
```
Commands:
client status <id>
	Show status of a specific client

client volume <id> <volume>
	Set volume of a specific client

client name <id> <name>
	Change name of a specific client

client latency <id> <latency>
	Set latency of a specific client

group status <id>
group mute <id> <true|false>
group clients <id>
group name <id> <name>

server status
	Show all groups and clients
server version
	Show RPC version of server
server deleteclient <id>

stream addstream <id>
stream removestream <id>
stream control <id> <play|pause|stop
```

## Installation
### Download .deb packages (Version 0.4.1)

* [amd64](http://deb.flupps.net/pool/main/s/snapcontrol/snapcontrol_0.4.1_amd64.deb)
* [arm64](http://deb.flupps.net/pool/main/s/snapcontrol/snapcontrol_0.4.1_arm64.deb)
* [armhf](http://deb.flupps.net/pool/main/s/snapcontrol/snapcontrol_0.4.1_armhf.deb)
* [i386](http://deb.flupps.net/pool/main/s/snapcontrol/snapcontrol_0.4.1_i386.deb)

### Installation via apt

Debian packages for all of the above architectures are available on the repository at deb.flupps.net

```
sudo echo 'deb http://deb.flupps.net/ stable main' > /etc/apt/sources.list.d/flupps.list
apt update
apt install snapcontrol
```

### Building snapcontrol
```aiignore
go build
```
