
DEBIAN_VERSION=bookworm
DEBIAN_VERSION_NUMBER=12
DEBIAN_IMAGE_VERSION=20240717-1811
DEBIAN_IMAGE_TYPE=genericcloud
DEBIAN_ARCH=amd64

all: download

clean:
	rm -f debian-${DEBIAN_VERSION_NUMBER}-genericcloud-${DEBIAN_ARCH}.qcow2 \
	      debian-${DEBIAN_VERSION_NUMBER}-genericcloud-${DEBIAN_ARCH}.json

download: debian-${DEBIAN_VERSION_NUMBER}-genericcloud-${DEBIAN_ARCH}.qcow2

debian-${DEBIAN_VERSION_NUMBER}-genericcloud-${DEBIAN_ARCH}.qcow2: debian-${DEBIAN_VERSION_NUMBER}-genericcloud-${DEBIAN_ARCH}.json
	wget -O debian-${DEBIAN_VERSION_NUMBER}-${DEBIAN_IMAGE_TYPE}-${DEBIAN_ARCH}.qcow2 \
		https://cloud.debian.org/images/cloud/${DEBIAN_VERSION}/${DEBIAN_IMAGE_VERSION}/debian-${DEBIAN_VERSION_NUMBER}-${DEBIAN_IMAGE_TYPE}-${DEBIAN_ARCH}-${DEBIAN_IMAGE_VERSION}.qcow2

debian-${DEBIAN_VERSION_NUMBER}-genericcloud-${DEBIAN_ARCH}.json:
	wget -O debian-${DEBIAN_VERSION_NUMBER}-${DEBIAN_IMAGE_TYPE}-${DEBIAN_ARCH}.json \
		https://cloud.debian.org/images/cloud/${DEBIAN_VERSION}/${DEBIAN_IMAGE_VERSION}/debian-${DEBIAN_VERSION_NUMBER}-${DEBIAN_IMAGE_TYPE}-${DEBIAN_ARCH}-${DEBIAN_IMAGE_VERSION}.json
