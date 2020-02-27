package main

import (
	"crypto/sha1"
	"log"
	"os"
	"syahrultorrent/serialize"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	// Lets open a bittorent file and the the correct bencodeTorrent struct
	fileReader, err := os.Open("./files/debian-10.3.0-amd64-netinst.iso.torrent")
	if err != nil {
		log.Println(err)
	}
	bencodeTorrent, err := serialize.Open(fileReader)
	if err != nil {
		log.Println(err)
	}
	torrentFile, err := bencodeTorrent.ToTorrentFile()
	if err != nil {
		log.Println(err)
	}
	// Generate a random peer id
	peerID := sha1.Sum([]byte("Hello World!"))
	res, err := torrentFile.RequestPeers(peerID, 8768)
	spew.Dump(res)
}
