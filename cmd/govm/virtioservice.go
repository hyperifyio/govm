// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/diskfs/go-diskfs"
	disk2 "github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/diskfs/go-diskfs/filesystem/iso9660"
	"golang.org/x/crypto/pbkdf2"
	"libvirt.org/go/libvirt"
)

type VirtioService struct {
	system         string
	imagesPath     string
	volumesPath    string
	interfaceType  string
	defaultNetwork string
}

func NewVirtioService(
	system, imagesPath, volumesPath, interfaceType, defaultNetwork string,
) *VirtioService {
	return &VirtioService{
		system:         system,
		imagesPath:     imagesPath,
		volumesPath:    volumesPath,
		interfaceType:  interfaceType,
		defaultNetwork: defaultNetwork,
	}
}

func (s *VirtioService) Start() error {
	log.Printf("Connecting libvirt to %s", s.system)
	conn, err := libvirt.NewConnect(s.system)
	if err != nil {
		return fmt.Errorf("start: Could not connect to libvirt: %v", err)
	}
	defer conn.Close()
	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	if err != nil {
		return fmt.Errorf("start: Could not list domains: %v", err)
	}
	fmt.Printf("%d running domains\n", len(doms))
	return nil
}

func (s *VirtioService) Stop() error {
	return nil
}

func (s *VirtioService) AddServer(
	name string,
) (*ServerModel, error) {

	fmt.Println("AddServer: Connecting to libvirt to add domain: ", name)

	// Connect to the local libvirt daemon
	conn, err := libvirt.NewConnect(s.system)
	if err != nil {
		return nil, fmt.Errorf("AddServer: Could not connect to libvirt: %v", err)
	}
	defer conn.Close()

	const domainType string = "qemu"
	const memory string = "1048576"
	const vcpu string = "1"
	const domainOsArchType string = "x86_64"
	const imageArch string = "amd64"
	const machine string = "pc-i440fx-9.0"
	const osType string = "hvm"
	const bootDev string = "hd"
	const imageType string = "qcow2"
	const vncPort string = "-1"
	const vncListen string = "127.0.0.1"
	const username string = "admin"
	const diskDevice string = "vda"
	const networkAddress string = "192.168.123.2"
	const networkGateway string = "192.168.123.1"
	const networkNetmask string = "255.255.255.0"
	const networkPrefix string = "24"

	// TODO: Generate new mac
	const macAddress string = "52:54:00:6b:3c:58"

	log.Printf("AddServer: Network address is %s/%s", networkAddress, networkPrefix)

	vncPassword, err := generatePassword(8)
	if err != nil {
		return nil, fmt.Errorf("AddServer: failed to generate vnc password: %v", err)
	}
	log.Printf("AddServer: Created VNC with password %s", vncPassword)

	userPassword, err := generatePassword(12)
	if err != nil {
		return nil, fmt.Errorf("AddServer: failed to generate user password: %v", err)
	}
	log.Printf("AddServer: Created user %s with password %s", username, userPassword)

	passwordSalt, err := generateSalt()
	if err != nil {
		return nil, fmt.Errorf("AddServer: failed to generate password salt: %v", err)
	}

	encryptedPassword, err := encryptPassword(userPassword, passwordSalt)
	if err != nil {
		return nil, fmt.Errorf("AddServer: failed to encrypt password: %v", err)
	}

	interfaceType := s.interfaceType

	imageFile := s.imagesPath + "/debian-12-genericcloud-" + imageArch + "." + imageType
	diskFile := s.volumesPath + "/" + name + "/" + name + "-" + diskDevice + "." + imageType
	ciDataFile := s.volumesPath + "/" + name + "/" + name + "-cidata.iso"

	// Copy the image to the destination directory
	destinationImagePath, err := copyImageFile(imageFile, diskFile)
	if err != nil {
		return nil, fmt.Errorf("AddServer: failed to copy image file: %v", err)
	}

	log.Printf("Image file copied to: %s", destinationImagePath)

	var interfaceXML string = ""
	if interfaceType == "network" {
		ifSourceNetwork := s.defaultNetwork
		interfaceXML = `
    <interface type='network'>
      <mac address='` + macAddress + `'/>
      <source network='` + ifSourceNetwork + `'/>
      <model type='virtio'/>
      <address type='pci' domain='0x0000' bus='0x00' slot='0x03' function='0x0'/>
    </interface>`
	} else if interfaceType == "user" {
		interfaceXML = `<interface type='user'>
	<model type='virtio'/>
	<mac address="` + macAddress + `"/>
	<ip family='ipv4' address='` + networkAddress + `' prefix='` + networkPrefix + `'/>
</interface>`
	}

	diskXML := `<disk type='file' device='disk'>
      <driver name='qemu' type='` + imageType + `'/>
      <source file='` + diskFile + `'/>
      <target dev='` + diskDevice + `' bus='virtio'/>
      <address type='pci' domain='0x0000' bus='0x00' slot='0x04' function='0x0'/>
    </disk>`

	cloudInitXML := `<disk type='file' device='cdrom'>
      <driver name='qemu' type='raw'/>
      <source file='` + ciDataFile + `'/>
      <target dev='hdb' bus='ide'/>
      <readonly/>
    </disk>`

	graphicsXML := `<graphics type='vnc' port='` + vncPort + `' listen='` + vncListen + `'  passwd='` + vncPassword + `'/>`

	// Define the domain XML
	domainXML := `
<domain type='` + domainType + `'>
  <name>` + name + `</name>
  <memory unit='KiB'>` + memory + `</memory>
  <vcpu placement='static'>` + vcpu + `</vcpu>
  <os>
    <type arch='` + domainOsArchType + `' machine='` + machine + `'>` + osType + `</type>
    <boot dev='` + bootDev + `'/>
  </os>
  <devices>
` + diskXML + `
` + cloudInitXML + `
` + interfaceXML + `
` + graphicsXML + `
  </devices>
</domain>`

	// Define Cloud-Init configuration
	metaData := `instance-id: ` + name + `
local-hostname: ` + name

	userData := `#cloud-config
users:
  - name: ` + username + `
    passwd: ` + encryptedPassword + `
    sudo: ['ALL=(ALL) NOPASSWD:ALL']
    groups: sudo
    shell: /bin/bash`

	networkConfig := `
version: 2
ethernets:
  interface0:
    match:
      macaddress: "` + macAddress + `"
    set-name: interface0
    addresses:
      - ` + networkAddress + `/` + networkNetmask + `
    gateway4: ` + networkGateway + `
`

	// Create Cloud-Init ISO
	err = createCloudInitISO(ciDataFile, metaData, userData, networkConfig)
	if err != nil {
		return nil, fmt.Errorf("AddServer: failed to create Cloud-Init ISO: %v", err)
	}
	log.Printf("Cloud-Init ISO created successfully at %s", ciDataFile)

	// Create the domain
	domain, err := conn.DomainDefineXML(domainXML)
	if err != nil {
		return nil, fmt.Errorf("AddServer: failed to define domain: %v", err)
	}
	defer domain.Free()

	model, err := getServerModel(domain)
	if err != nil {
		return nil, fmt.Errorf("AddServer: failed to get domain data: %v", err)
	}

	fmt.Println("Domain created successfully: ", model.Name)
	return model, nil
}

func (s *VirtioService) GetServerList() ([]*ServerModel, error) {
	var servers []*ServerModel
	log.Printf("Connecting libvirt to %s", s.system)
	conn, err := libvirt.NewConnect(s.system)
	if err != nil {
		return nil, fmt.Errorf("GetServerList: failed to connect to libvirt: %v", err)
	}
	defer conn.Close()
	list, err := conn.ListAllDomains(0)
	if err != nil {
		return nil, fmt.Errorf("GetServerList: failed to list domains from libvirt: %v", err)
	}
	for _, item := range list {
		defer item.Free()
		model, err := getServerModel(&item)
		if err != nil {
			return nil, fmt.Errorf("AddServer: failed to get domain data: %v", err)
		}
		servers = append(servers, model)
	}
	return servers, nil
}

func (s *VirtioService) FindServer(targetName string) (*ServerModel, error) {

	log.Printf("FindServer: Connecting libvirt to %s", s.system)
	conn, err := libvirt.NewConnect(s.system)
	if err != nil {
		return nil, fmt.Errorf("FindServer: failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

	item, err := conn.LookupDomainByName(targetName)
	if err != nil {
		libvirtError, ok := err.(libvirt.Error)
		if ok && libvirtError.Code == libvirt.ERR_NO_DOMAIN {
			return nil, nil
		} else {
			return nil, fmt.Errorf("FindServer: Failed to find the domain: %s: %v", targetName, err)
		}
	}
	if item == nil {
		return nil, fmt.Errorf("FindServer: Failed to find the domain by name: %s", targetName)
	}
	defer item.Free()

	model, err := getServerModel(item)
	if err != nil {
		return nil, fmt.Errorf("FindServer: failed to get domain data: %v", err)
	}
	return model, nil

}

func (s *VirtioService) DeployServer(name string) (*ServerModel, error) {
	server, err := s.FindServer(name)
	if err != nil {
		return nil, fmt.Errorf("DeployServer: failed to find the server: error: %v", err)
	}
	if server == nil {
		return nil, fmt.Errorf("DeployServer: failed to find the server: not found")
	}
	if server.Status == UninitializedServerStatusCode {
		server.Status = DeployingServerStatusCode
		time.AfterFunc(3*time.Second, func() {
			server.Status = StoppedServerStatusCode
		})
	}
	return server, nil
}

func (s *VirtioService) StartServer(name string) (*ServerModel, error) {

	log.Printf("StartServer: Connecting libvirt to %s", s.system)
	conn, err := libvirt.NewConnect(s.system)
	if err != nil {
		return nil, fmt.Errorf("StartServer: failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

	item, err := conn.LookupDomainByName(name)
	if err != nil {
		return nil, fmt.Errorf("StartServer: Failed to find the domain: %s: %v", name, err)
	}
	if item == nil {
		return nil, fmt.Errorf("StartServer: Failed to find the domain by name: %s", name)
	}
	defer item.Free()

	err = item.Create()
	if err != nil {
		return nil, fmt.Errorf("StartServer: failed to start the domain: %v", err)
	}
	defer item.Free()

	model, err := getServerModel(item)
	if err != nil {
		return nil, fmt.Errorf("StartServer: failed to get domain data: %v", err)
	}
	return model, nil
}

func (s *VirtioService) StopServer(name string) (*ServerModel, error) {

	log.Printf("StopServer: Connecting libvirt to %s", s.system)
	conn, err := libvirt.NewConnect(s.system)
	if err != nil {
		return nil, fmt.Errorf("StopServer: failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

	item, err := conn.LookupDomainByName(name)
	if err != nil {
		return nil, fmt.Errorf("StopServer: Failed to find the domain: %s: %v", name, err)
	}
	if item == nil {
		return nil, fmt.Errorf("StopServer: Failed to find the domain by name: %s", name)
	}
	defer item.Free()

	err = item.Shutdown()
	if err != nil {
		return nil, fmt.Errorf("StopServer: failed to start the domain: %v", err)
	}
	defer item.Free()

	model, err := getServerModel(item)
	if err != nil {
		return nil, fmt.Errorf("StopServer: failed to get domain data: %v", err)
	}
	return model, nil
}

func (s *VirtioService) RestartServer(name string) (*ServerModel, error) {

	log.Printf("RestartServer: Connecting libvirt to %s", s.system)
	conn, err := libvirt.NewConnect(s.system)
	if err != nil {
		return nil, fmt.Errorf("RestartServer: failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

	item, err := conn.LookupDomainByName(name)
	if err != nil {
		return nil, fmt.Errorf("RestartServer: Failed to find the domain: %s: %v", name, err)
	}
	if item == nil {
		return nil, fmt.Errorf("RestartServer: Failed to find the domain by name: %s", name)
	}
	// defer item.Free() // Handled at gracefulRestart

	go gracefulRestart(item, name)

	model, err := getServerModel(item)
	if err != nil {
		return nil, fmt.Errorf("RestartServer: failed to get domain data: %v", err)
	}
	return model, nil
}

func (s *VirtioService) DeleteServer(name string) (*ServerModel, error) {

	log.Printf("DeleteServer: Connecting libvirt to %s", s.system)
	conn, err := libvirt.NewConnect(s.system)
	if err != nil {
		return nil, fmt.Errorf("DeleteServer: failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

	item, err := conn.LookupDomainByName(name)
	if err != nil {
		return nil, fmt.Errorf("DeleteServer: Failed to find the domain: %s: %v", name, err)
	}
	if item == nil {
		return nil, fmt.Errorf("DeleteServer: Failed to find the domain by name: %s", name)
	}
	defer item.Free()

	model, err := getServerModel(item)
	if err != nil {
		return nil, fmt.Errorf("DeleteServer: failed to get domain data: %v", err)
	}

	err = item.Undefine()
	if err != nil {
		return nil, fmt.Errorf("DeleteServer: failed to delete the domain: %v", err)
	}

	log.Printf("Domain deleted successfully: %s", name)
	model.Status = DeletedServerStatusCode
	return model, nil
}

var _ ServerService = &VirtioService{}

func getServerModel(item *libvirt.Domain) (*ServerModel, error) {
	state, _, err := item.GetState()
	if err != nil {
		return nil, fmt.Errorf("failed to get domain state: %v", err)
	}
	name, err := item.GetName()
	if err != nil {
		return nil, fmt.Errorf("failed to get domain name: %v", err)
	}
	return NewServerModel(name, domainStateToServerStatusCode(state)), nil
}

func domainStateToServerStatusCode(state libvirt.DomainState) ServerStatusCode {
	switch state {
	case libvirt.DOMAIN_NOSTATE:
		return UninitializedServerStatusCode
	case libvirt.DOMAIN_RUNNING:
		return StartedServerStatusCode
	case libvirt.DOMAIN_BLOCKED:
		return BlockedServerStatusCode
	case libvirt.DOMAIN_PAUSED:
		return PausedServerStatusCode
	case libvirt.DOMAIN_SHUTDOWN:
		return StoppingServerStatusCode
	case libvirt.DOMAIN_SHUTOFF:
		return StoppedServerStatusCode
	case libvirt.DOMAIN_CRASHED:
		return CrashedServerStatusCode
	case libvirt.DOMAIN_PMSUSPENDED:
		return SuspendedServerStatusCode
	default:
		return UnknownServerStatusCode
	}
}

// gracefulRestart performs a graceful restart of the domain
func gracefulRestart(domain *libvirt.Domain, domainName string) {

	// Gracefully stop the domain
	err := domain.Shutdown()
	if err != nil {
		log.Printf("gracefulRestart: Failed to shut down the domain: %v", err)
		return
	}
	log.Printf("gracefulRestart: Domain '%s' shutdown signal sent", domainName)

	// Wait for the domain to shut down completely
	for {
		state, _, err := domain.GetState()
		if err != nil {
			log.Printf("gracefulRestart: ERROR: Failed to get domain state: %v", err)
			return
		}
		if state == libvirt.DOMAIN_SHUTOFF {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Start the domain again
	err = domain.Create()
	if err != nil {
		log.Printf("gracefulRestart: ERROR: Failed to start the domain: %v", err)
		return
	}
	log.Printf("gracefulRestart: Domain '%s' restarted successfully", domainName)
}

func copyImageFile(sourcePath, destinationFile string) (string, error) {

	destinationDir := filepath.Dir(destinationFile)
	destinationBasename := filepath.Base(destinationFile)
	destinationPath := filepath.Join(destinationDir, destinationBasename)

	if err := os.MkdirAll(destinationDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create destination directory: %w", err)
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return "", fmt.Errorf("failed to open source image: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destinationPath)
	if err != nil {
		return "", fmt.Errorf("failed to create image: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return "", fmt.Errorf("failed to copy image: %w", err)
	}

	return destinationPath, nil
}

func createCloudInitISO(isoPath, metaData, userData, networkConfig string) error {

	// Create the ISO file
	disk, err := diskfs.Create(isoPath, 38912, diskfs.Raw, 2048)
	if err != nil {
		return fmt.Errorf("createCloudInitISO: failed to create ISO file: %v", err)
	}

	// Create an ISO9660 filesystem with default settings
	fs, err := disk.CreateFilesystem(disk2.FilesystemSpec{
		// Partition:   0,
		FSType: filesystem.TypeISO9660,
		// VolumeLabel: "",
		// WorkDir:     "",
	})
	if err != nil {
		return fmt.Errorf("createCloudInitISO: failed to create ISO9660 filesystem: %v", err)
	}

	// Cast the filesystem to ISO9660 to add files
	isoFs := fs.(*iso9660.FileSystem)

	// Add files to the ISO filesystem
	files := map[string]string{
		"/meta-data":      metaData,
		"/user-data":      userData,
		"/network-config": networkConfig,
	}

	for path, content := range files {
		file, err := isoFs.OpenFile(path, os.O_CREATE|os.O_RDWR)
		if err != nil {
			return fmt.Errorf("createCloudInitISO: failed to open file '%s' in ISO filesystem: %v", path, err)
		}
		_, err = file.Write([]byte(content))
		if err != nil {
			return fmt.Errorf("createCloudInitISO: failed to write to file '%s' in ISO filesystem: %v", path, err)
		}
		file.Close()
	}

	// Finalize the ISO filesystem
	err = isoFs.Finalize(iso9660.FinalizeOptions{
		RockRidge:        true,
		VolumeIdentifier: "CIDATA",
	})
	if err != nil {
		return fmt.Errorf("createCloudInitISO: failed to finalize ISO filesystem: %v", err)
	}

	fmt.Printf("ISO file created successfully at '%s'\n", isoPath)
	return nil
}

func generateSalt() (string, error) {
	salt := make([]byte, 16) // 16 bytes salt for SHA-512 crypt
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %v", err)
	}
	return base64.RawStdEncoding.EncodeToString(salt), nil
}

func encryptPassword(password, salt string) (string, error) {
	// Use pbkdf2 to hash the password with SHA-512
	hash := pbkdf2.Key([]byte(password), []byte(salt), 10000, sha512.Size, sha512.New)
	// Format it into the required $6$<salt>$<hash> format
	return fmt.Sprintf("$6$%s$%s", salt, base64.RawStdEncoding.EncodeToString(hash)), nil
}
