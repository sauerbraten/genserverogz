package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

/*
struct octaheader
{
    char magic[4];              // "OCTA"
    int version;                // any >8bit quantity is little endian
    int headersize;             // sizeof(header)
    int worldsize;
    int numents;
    int numpvs;
    int lightmaps;
    int blendmap;
    int numvars;
    int numvslots;
};

struct compatheader             // map file format header
{
    char magic[4];              // "OCTA"
    int version;                // any >8bit quantity is little endian
    int headersize;             // sizeof(header)
    int worldsize;
    int numents;
    int numpvs;
    int lightmaps;
    int lightprecision, lighterror, lightlod;
    uchar ambient;
    uchar watercolour[3];
    uchar blendmap;
    uchar lerpangle, lerpsubdiv, lerpsubdivsize;
    uchar bumperror;
    uchar skylight[3];
    uchar lavacolour[3];
    uchar waterfallcolour[3];
    uchar reserved[10];
    char maptitle[128];
};
*/

type Header struct {
	Magic      [4]byte
	Version    int32
	HeaderSize int32
	_          int32 // worldsize
	NumEnts    int32
	_          int32 // numpvs
	_          int32 // lightmaps
}

func readHeader(stream io.Reader) (*Header, int, error) {
	hdr := Header{}
	err := binary.Read(stream, binary.LittleEndian, &hdr)
	if err != nil {
		return nil, 0, fmt.Errorf("reading OGZ header: %w", err)
	}

	if string(hdr.Magic[:]) != "OCTA" {
		return nil, 0, fmt.Errorf("reading OGZ header: wrong magic (not a OGZ file?)")
	}

	if hdr.HeaderSize < 36 {
		return nil, 0, fmt.Errorf("reading OGZ header: too short (%d bytes, need at least 36)", hdr.HeaderSize)
	}

	// magic, version, headersize, worldsize, numents, numpvs, lightmaps are read (7*4 bytes)

	if hdr.Version < 29 {
		// headersize seems to have been a uint8 in the old format, but
		// we parsed it (together with 3 obsolete bytes) as int32!
		hdr.HeaderSize = int32(uint8(hdr.HeaderSize))
	}

	buf := make([]byte, hdr.HeaderSize-(7*4))
	_, err = stream.Read(buf)
	if err != nil {
		return nil, 0, fmt.Errorf("reading remaining OGZ header: %w", err)
	}

	if hdr.Version < 29 {
		// OGZs before version 29 don't have a variables section
		return &hdr, 0, nil
	}

	return &hdr, int(binary.LittleEndian.Uint32(buf[4:])), nil
}

func readMapVars(stream io.Reader, numVars int) (map[string]interface{}, error) {
	type MapVarHeader struct {
		Type          byte
		IdentifierLen uint16
	}

	vars := map[string]interface{}{}
	for i := 0; i < numVars; i++ {
		mvh := MapVarHeader{}
		err := binary.Read(stream, binary.LittleEndian, &mvh)
		if err != nil {
			return nil, fmt.Errorf("reading map variable header at pos %d: %w", i, err)
		}

		identifier := make([]byte, mvh.IdentifierLen)
		err = binary.Read(stream, binary.LittleEndian, identifier)
		if err != nil {
			return nil, fmt.Errorf("reading map variable identifier at pos %d: %w", i, err)
		}

		var (
			ival int32
			fval float32
			sval []byte

			p interface{}
		)

		switch mvh.Type {
		case 0:
			p = &ival
		case 1:
			p = &fval
		case 2:
			var vlen uint16
			err = binary.Read(stream, binary.LittleEndian, &vlen)
			if err != nil {
				return nil, fmt.Errorf("reading string value length of map variable at pos %d: %w", i, err)
			}
			sval = make([]byte, vlen)
			p = &sval
		}

		err = binary.Read(stream, binary.LittleEndian, p)
		if err != nil {
			return nil, fmt.Errorf("reading map variable value at pos %d: %w", i, err)
		}

		switch mvh.Type {
		case 0:
			vars[string(identifier)] = ival
		case 1:
			vars[string(identifier)] = fval
		case 2:
			vars[string(identifier)] = string(sval)
		}
	}

	return vars, nil
}

func readGameIdentifier(stream io.Reader, version int32) (string, error) {
	if version < 16 {
		return "fps", nil
	}

	buf := make([]byte, 1)
	n, err := stream.Read(buf)
	if err != nil {
		return "", fmt.Errorf("reading game identifier length: %w", err)
	}
	if n != 1 {
		return "", fmt.Errorf("reading game identifier length: tried to read 1 byte, read %d", n)
	}

	identifierLen := buf[0]

	buf = make([]byte, identifierLen+1) // read terminating zero byte
	n, err = stream.Read(buf)
	if err != nil {
		return "", fmt.Errorf("reading game identifier: %w", err)
	}
	if n != int(identifierLen+1) {
		return "", fmt.Errorf("reading game identifier: tried to read %d bytes, read %d", identifierLen, n)
	}

	if buf[identifierLen] != 0x00 {
		return "", fmt.Errorf("reading game identifier: not zero-byte terminated (last byte is %d)", buf[identifierLen])
	}

	return string(buf[:identifierLen]), nil
}

// reads extraentinfosize (aka einfosize, aka eif) and skips extras, returns extraentinfosize
func readExtraEntInfoLen(stream io.Reader, version int32) (int, error) {
	if version < 16 {
		return 0, nil
	}

	buf := make([]uint16, 2)
	err := binary.Read(stream, binary.LittleEndian, buf)
	if err != nil {
		return 0, fmt.Errorf("reading extraentinfosize: %w", err)
	}

	discard := make([]byte, int(buf[1]))
	_, err = stream.Read(discard)
	if err != nil {
		return 0, fmt.Errorf("skipping extras section: %w", err)
	}

	return int(buf[0]), nil
}

func skipMostRecentlyUsed(stream io.Reader, version int32) error {
	if version < 14 {
		discard := make([]byte, 256)
		_, err := stream.Read(discard)
		if err != nil {
			return fmt.Errorf("skipping 'most recently used' section: %w", err)
		}
		return nil
	}

	var numMRU uint16
	err := binary.Read(stream, binary.LittleEndian, &numMRU)
	if err != nil {
		return fmt.Errorf("reading size of 'most recently used' section: %w", err)
	}

	discard := make([]byte, int(numMRU*2)) // one MRU ~ 16bit
	_, err = stream.Read(discard)
	if err != nil {
		return fmt.Errorf("skipping 'most recently used' section: %w", err)
	}

	return nil
}

type Entity struct {
	X, Y, Z                           float32 // position
	Attr1, Attr2, Attr3, Attr4, Attr5 int16
	Type                              byte
	Reserved                          byte
}

func (e *Entity) String() string {
	return fmt.Sprintf("type: %2d, attrs: %3d %3d %3d %3d %3d, pos: %f %f %f", e.Type, e.Attr1, e.Attr2, e.Attr3, e.Attr4, e.Attr5, e.X, e.Y, e.Z)
}

func readEnts(stream io.Reader, numEnts int32, extraEntInfoLen int) ([]Entity, error) {
	ents := make([]Entity, 0, numEnts)
	for i := 0; i < int(numEnts); i++ {
		ent := Entity{}
		err := binary.Read(stream, binary.LittleEndian, &ent)
		if err != nil {
			return nil, fmt.Errorf("reading ent at position %d: %w", i, err)
		}
		ents = append(ents, ent)
		if extraEntInfoLen > 0 {
			discard := make([]byte, extraEntInfoLen)
			_, err = stream.Read(discard)
			if err != nil {
				return nil, fmt.Errorf("skipping extra entity info at position %d: %w", i, err)
			}
		}
	}

	return ents, nil
}

// writes a minimal OGZ only containing the ents data
// (v29, since that doesn't have numvslots and is thus 4 bytes smaller than v30+)
//
// the file will contain the following fields:
//
// length in bytes | contents
//              36 | header
//               5 | gameident: len("fps"), "fps", 0x00
//               2 | extraentinfosize: 0x00 0x00
//               2 | len(extras): 0x00 0x00
//               2 | nummru: 0x00 0x00
//  len(ents) * 24 | ents
//
// each ent is 24 bytes in length:
//
// length in bytes | contents
//             3*4 | x, y, z
//             2*5 | attributes 1-5
//               1 | type
//               1 | reserved
func writeEntsOnlyOGZ(ents []Entity, w io.Writer) error {
	err := binary.Write(w, binary.LittleEndian, []byte{'O', 'C', 'T', 'A'})
	if err != nil {
		return fmt.Errorf("writing magic: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, []int32{
		29,               // version
		36,               // headersize
		1,                // worldsize
		int32(len(ents)), // numents
		0,                // numpvs
		0,                // lightmaps
		0,                // blendmap
		0,                // numvars
	})
	if err != nil {
		return fmt.Errorf("writing header: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, []byte{
		3, // len(fps)
		'f', 'p', 's',
		0,
	})
	if err != nil {
		return fmt.Errorf("writing game identifier: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, []uint16{
		0, // extraentinfosize
		0, // len(extras)
		0, // nummru
	})
	if err != nil {
		return fmt.Errorf("writing extras/MRU size: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, ents)
	if err != nil {
		return fmt.Errorf("writing entities: %w", err)
	}

	return nil
}

var (
	printVersion        = flag.Bool("version", false, "print file format version")
	printVars           = flag.Bool("vars", false, "print map vars (version 29+ only)")
	printGameIdentifier = flag.Bool("game", false, "print game identifier")
	printEnts           = flag.Bool("ents", false, "print map entities")
)

func init() {
	flag.CommandLine.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Reads uncompressed OGZ file on stdin and writes a minimal uncompressed OGZ file containing only the entity data to stdout.\n"+
				"Specify any print flag to only print the requested fields to stdout instead of the shrunk map data.\n",
		)
		flag.Usage()
	}
}

func main() {
	flag.Parse()

	hdr, numVars, err := readHeader(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing input stream: %v\n", err)
		os.Exit(1)
	}
	if *printVersion {
		fmt.Printf("OGZ file format version: %d\n", hdr.Version)
	}

	vars, err := readMapVars(os.Stdin, numVars)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing input stream: %v\n", err)
		os.Exit(1)
	}
	if *printVars {
		fmt.Printf("map variables: (%d)\n", len(vars))
		for k, v := range vars {
			switch _v := v.(type) {
			case string:
				v = strings.ReplaceAll(_v, "\f", "\\f")
			}
			fmt.Printf("  %s = %v\n", k, v)
		}
	}

	gameIdentifier, err := readGameIdentifier(os.Stdin, hdr.Version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing input stream: %v\n", err)
		os.Exit(1)
	}
	if *printGameIdentifier {
		fmt.Printf("game: %s\n", gameIdentifier)
	}

	extraEntInfoLen, err := readExtraEntInfoLen(os.Stdin, hdr.Version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing input stream: %v\n", err)
		os.Exit(1)
	}

	err = skipMostRecentlyUsed(os.Stdin, hdr.Version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing input stream: %v\n", err)
		os.Exit(1)
	}

	ents, err := readEnts(os.Stdin, hdr.NumEnts, extraEntInfoLen)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing input stream: %v\n", err)
		os.Exit(1)
	}
	if *printEnts {
		fmt.Printf("map entities: (%d)\n", len(ents))
		for i, e := range ents {
			fmt.Printf("  %4d: %s\n", i, e.String())
		}
	}

	if *printVersion || *printVars || *printGameIdentifier || *printEnts {
		return
	}

	err = writeEntsOnlyOGZ(ents, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "writing stripped OGZ: %v\n", err)
		os.Exit(1)
	}
}
