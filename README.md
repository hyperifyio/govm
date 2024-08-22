# @hyperifyio/govm

This is the source code for the GoVM Virtual Manager. It includes the server written in Go and a web frontend.

For more information about the progress of the project, see [GoVM Project](https://github.com/hyperifyio/project-govm/issues/1).

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

### DevOps

Our devops pipelines use following secrets:

`RELEASE_PAT` is a GitHub Personal Access Token.
See [Managing your personal access tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens).

`RELEASE_PRIVATE_KEY` is a private random passphrase which is used to encrypt 
releases using [goselfshield](https://github.com/hyperifyio/goselfshield).
