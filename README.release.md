## Installing govm -- a virtual server manager

Use `./govm --install-private-key SECRET_KEY --install-output=/opt/govm/govm` to
decrypt and install the game server executable.

For systemd configuration, install the provided service file to
`/etc/systemd/system/govm.service`.

Use `/opt/govm/govm --init-private-key` to generate a private key for the 
server.

Configure the env by editing the file `/etc/govm/env`.
