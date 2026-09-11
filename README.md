# Snapcontrol

Snapcontrol is a little CLI tool to control snapserver. It aims to be a full implementation of the snapcast API.

## Usage
Note: `<id>` for clients and groups can be found in the output of `server status`

### Commands:

| Command                                               | Description                      |
|-------------------------------------------------------|----------------------------------|
| `snapcontrol client status <id\|name>`                | Show status of a specific client |
| `snapcontrol client volume <id\|name> <volume>` 	     | Set volume of a specific client  |
| `snapcontrol client name <id\|name> <name>`	          | Change name of a specific client |
| `snapcontrol client latency <id\|name> <latency>`	    | Set latency of a specific client |
| `snapcontrol group status <id>`                       | Show status of aspecific group   |
| `snapcontrol group mute <id> <true\|false>`           | Mute a specific group            |
| `snapcontrol group stream <id> <streamId>`            | Assign a stream to a group       |
| `snapcontrol group clients <id> <clientId\|name>...`  | Replace a group's clients         |
| `snapcontrol group name <id> <name>`                  | Change name of a specific group  |
| `snapcontrol server status`	                          | Show all groups and clients      |
| `snapcontrol server version`                          | 	Show RPC version of server      |
| `snapcontrol server deleteclient <id>`                |                                  |
| `snapcontrol stream add <streamUri>`                  | Add a stream                      |
| `snapcontrol stream remove <id>`                      | Remove a stream                   |
| `snapcontrol stream control <id> <command> [key=value ...]` | Control a stream           |
| `snapcontrol stream property <id> <property> <jsonValue>` | Set a stream property          |

## Installation
### Download .deb packages (Version 0.5.1)

* [amd64](http://deb.flupps.net/pool/main/s/snapcontrol/snapcontrol_0.5.1_amd64.deb)
* [arm64](http://deb.flupps.net/pool/main/s/snapcontrol/snapcontrol_0.5.1_arm64.deb)
* [armhf](http://deb.flupps.net/pool/main/s/snapcontrol/snapcontrol_0.5.1_armhf.deb)
* [i386](http://deb.flupps.net/pool/main/s/snapcontrol/snapcontrol_0.5.1_i386.deb)

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
