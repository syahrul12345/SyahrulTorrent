package serialize

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/jackpal/bencode-go"
)

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

// TorrentFile is the torrent file that holds information of the file to be downloaded
type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHash   [][20]byte
	PieceLength int
	Length      int
	Name        string
}

// Hashes the bencodeInfo
func (bencodeInfo *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *bencodeInfo)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

// Split the pieces into piecehashes
func (bencodeInfo *bencodeInfo) splitHash() ([][20]byte, error) {
	hashlen := 20
	buf := []byte(bencodeInfo.Pieces)
	if len(buf)%hashlen != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashlen
	hashes := make([][20]byte, numHashes)
	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashlen:(i+1)*hashlen])
	}
	return hashes, nil
}

// Open parses a torrent file
func Open(r io.Reader) (*bencodeTorrent, error) {
	bto := bencodeTorrent{}
	err := bencode.Unmarshal(r, &bto)
	if err != nil {
		return nil, err
	}
	return &bto, nil
}

// Convert the bencodeTorrent to a TorrentFile type
func (bto *bencodeTorrent) toTorrentFile() (*TorrentFile, error) {
	// Create an empty torrentFile struct
	infoHash, err := bto.Info.hash()
	if err != nil {
		return nil, err
	}
	pieceHash, err := bto.Info.splitHash()
	if err != nil {
		return nil, err
	}
	return &TorrentFile{
		Announce:    bto.Announce,
		InfoHash:    infoHash,
		PieceHash:   pieceHash,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}, nil
}
