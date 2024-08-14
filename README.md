# @hyperifyio/govm

Source code for govm project

## Clone 

```
git clone --recurse-submodules -j8 git@github.com:hyperifyio/govm.git
```

...or update:

```
git submodule update
```

## Requirements for MacOS

For MacOS, you need to install `libvirt` and QEMU version 2.12 or newer.

```
brew install dbus-glib libvirt
```

Read more about [libvirt MacOS support](https://libvirt.org/macos.html).

Then either install it as a service:

`brew services start libvirt`

..or run manually: `/opt/homebrew/opt/libvirt/sbin/libvirtd -f /opt/homebrew/etc/libvirt/libvirtd.conf`

Then use `GOVM_SYSTEM='qemu:///session'` when starting the govm.

## Starting the server with Docker for development

```bash
docker-compose build && docker-compose up
```

Once started, the server is available at http://localhost:8080

## Starting the server from localhost

You can start the server locally like this:

```
PRIVATE_KEY=9ca549e8e80e363cb92b99936dd869c65eca7f474d2b595a72d5e9a2d79eff61 \
./govm
```

The command above works if you have our development Docker setup running with 
default settings.

## Manual testing with Curl

### Starting a virtual server

Request body:

```json
{
}
```

Command: 

```bash
curl -i -d '{}' http://localhost:3001
```

Response:

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 07 Apr 2024 23:41:23 GMT
Content-Length: 436
```

```json
{
}
```

### Starting a server

Request body:

```json
{
}
```

Command:
```bash
curl -i -d '{}' http://localhost:3001
```

Response:

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 07 Apr 2024 23:42:07 GMT
Content-Length: 436
```

```json
{
}
```
