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
| `snapcontrol group clients <id>`                      | Show clients of a specific group |
| `snapcontrol group name <id> <name>`                  | Change name of a specific group  |
| `snapcontrol server status`	                          | Show all groups and clients      |
| `snapcontrol server version`                          | 	Show RPC version of server      |
| `snapcontrol server deleteclient <id>`                |                                  |
| `snapcontrol stream addstream <id>`                   |                                  |
| `snapcontrol stream removestream <id>`                |                                  |
| `snapcontrol stream control <id> <play\|pause\|stop>` |                                  |

## Installation
### Download .deb packages (Version 0.4.2)

* [amd64](http://deb.flupps.net/pool/main/s/snapcontrol/snapcontrol_0.4.2_amd64.deb)
* [arm64](http://deb.flupps.net/pool/main/s/snapcontrol/snapcontrol_0.4.2_arm64.deb)
* [armhf](http://deb.flupps.net/pool/main/s/snapcontrol/snapcontrol_0.4.2_armhf.deb)
* [i386](http://deb.flupps.net/pool/main/s/snapcontrol/snapcontrol_0.4.2_i386.deb)

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
