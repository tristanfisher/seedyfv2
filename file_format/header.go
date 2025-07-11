package file_format

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
)

// LittleEndian as a struct exists to (ab)use a structs to create a namespace/static methods
type LittleEndian struct {
	ErrorLog *log.Logger
}

func (le LittleEndian) USHORT(b []byte) (USHORT, error) {
	if len(b) < 2 {
		return 0, errors.New("not enough bytes provided for conversion")
	}
	return USHORT(binary.LittleEndian.Uint16(b)), nil
}

func (le LittleEndian) ULONG(b []byte) (ULONG, error) {
	if len(b) < 4 {
		return 0, errors.New("not enough bytes provided for conversion")
	}
	return ULONG(binary.LittleEndian.Uint32(b)), nil
}

func (le LittleEndian) FSINDEX(b []byte) (FSINDEX, error) {
	v, err := le.ULONG(b)
	return FSINDEX(v), err
}

func (le LittleEndian) SECT(b []byte) (SECT, error) {
	v, err := le.ULONG(b)
	return SECT(v), err
}

func (le LittleEndian) DFSIGNATURE(b []byte) (DFSIGNATURE, error) {
	v, err := le.ULONG(b)
	return DFSIGNATURE(v), err
}

// SectFatSize is the first 109 FAT sectors, which are SECT (uint32)
const SectFatSize = 109

func (le LittleEndian) SECTFAT(b []byte) ([109]SECT, error) {

	// see "ULONG" note. we're using uint32 as cdfv2 uses 32-bit size uinsigned int (4 bytes == 32 bit)
	const chunkSize = 4

	// split up byte region into chunks
	var sectors [109]SECT

	// chunk into 4 bytes, then convert to SECT (32 bit ulong -> uint32)
	chunk := []byte{}
	sectId := 0
	for i, nibble := range b {
		// split into n-sized chunks.  don't split on zeroth position
		if i != 0 && i%chunkSize == 0 {
			// we have enough data to convert to ulong / uint32
			sect, err := le.SECT(chunk)
			if err != nil {
				return [109]SECT{}, err
			}

			// note 4294967295 => 0xffffffff
			// e.g. 00000058: ffff ffff ffff ffff  ........

			sectors[sectId] = sect
			sectId++
			// clear out chunk container
			chunk = []byte{}
		}
		chunk = append(chunk, nibble)
	}

	return sectors, nil
}

//
// aliased types for ease of referencing against documentation
//

// FSINDEX is an alias for a ULONG (64-bit unsigned)
type FSINDEX ULONG

type DFSIGNATURE ULONG

// GUID -> 16 byte structure with different representations. Microsoft type.
//
//	 typedef GUID CLSID // 16 bytes
//	https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-dtyp/4926e530-816e-41c2-b251-ec5c7aca018a
type GUID [16]byte

// /aliases

// StructuredStorageHeader contains the information
// required for instantiating and parsing a compound file.
//
// Comments predominantly from:
// - Advanced Authoring Format (AAF) Low-Level Container Specification v1.0.1
// - [MS-CFB] - v20240423
type StructuredStorageHeader struct {

	// [offset from start (bytes in hex), length (bytes)]

	// [00H,08] {0xd0, 0xcf, 0x11, 0xe0, 0xa1, 0xb1, 0x1a, 0xe1} for current version
	// decimal offset: 0
	ABSig [8]byte

	// [08H,16] reserved must be zero (WriteClassStg/
	// GetClassFile uses root directory class id)
	// decimal offset: 08
	Clsid GUID

	// [18H,02] minor version of the format: 33 is
	// written by reference implementation
	// decimal offset: 24
	MinorVersion USHORT

	// [1AH,02] major version of the dll/format: 3 for
	// 512-byte sectors, 4 for 4 KB sectors
	// aka DLL Version in some documentation
	// decimal offset: 26
	MajorVersion USHORT

	// [1CH,02] 0xFFFE: indicates Intel byte-ordering
	// decimal offset: 28
	ByteOrder USHORT

	// [1EH,02] size of sectors in power-of-two;
	// typically 9 indicating 512-byte sectors
	// if major version is 3, must be 0x0009
	// if major version is 4, must be 0x000c
	// decimal offset: 30
	SectorShift USHORT

	// [20H,02] size of mini-sectors in power-of-two;
	// typically 6 indicating 64-byte mini-sectors
	// decimal offset: 32
	MiniSectorShift USHORT

	// [22H,02] reserved, must be zero
	// decimal offset: 34
	Reserved USHORT

	// [24H,04] reserved, must be zero
	// decimal offset: 36
	Reserved2 ULONG

	// [28H,04] must be zero for 512-byte sectors,
	// number of SECTs in directory chain for 4 KB
	// sectors
	// decimal offset: 40
	CSectDir ULONG

	// [2CH,04] number of SECTs in the FAT chain
	// decimal offset: 44
	CSectFat FSINDEX

	// [30H,04] first SECT in the directory chain
	// decimal offset: 48
	SectDirStart SECT

	// [34H,04] signature used for transactions; must
	// be zero. The reference implementation
	// does not support transactions
	// decimal offset: 52
	Signature DFSIGNATURE

	// [38H,04] maximum size for a mini stream;
	// typically 4096 bytes
	// decimal offset: 56
	MiniSectorCutoff ULONG

	// [3CH,04] first SECT in the MiniFAT chain
	// decimal offset: 60
	SectMiniFatStart SECT

	// [40H,04] number of SECTs in the MiniFAT chain
	// decimal offset: 64
	CSectMiniFat FSINDEX

	// [44H,04] first SECT in the DIFAT chain
	// decimal offset: 68
	SectDifStart SECT

	// [48H,04] number of SECTs in the DIFAT chain
	// decimal offset: 72
	CSectDif FSINDEX

	// [4CH,436] the SECTs of first 109 FAT sectors
	// decimal offset: 76
	DIFAT [109]SECT
}

func (ssh StructuredStorageHeader) String() string {
	return fmt.Sprintf(""+
		"<ABSig: %x ; "+
		"Clsid: %x ; "+
		"MinorVersion: %x ; "+
		"MajorVersion: %x ; "+
		"ByteOrder: %x ; "+
		"SectorShift: %x ; "+
		"MiniSectorShift: %x ; "+
		"Reserved: %x ; "+
		"Reserved2: %x ; "+
		"CSectDir: %x ; "+
		"CSectFat: %x ; "+
		"SectDirStart: %x ; "+
		"Signature: %x ; "+
		"MiniSectorCutoff: %x ; "+
		"SectMiniFatStart: %x ; "+
		"CSectMiniFat: %x ; "+
		"SectDifStart: %x ; "+
		"CSectDif: %x ; "+
		"DIFAT: %x>",
		ssh.ABSig, ssh.Clsid, ssh.MinorVersion, ssh.MajorVersion, ssh.ByteOrder, ssh.SectorShift, ssh.MiniSectorShift, ssh.Reserved, ssh.Reserved2,
		ssh.CSectDir, ssh.CSectFat, ssh.SectDirStart, ssh.Signature, ssh.MiniSectorCutoff, ssh.SectMiniFatStart, ssh.CSectMiniFat, ssh.SectDifStart, ssh.CSectDif, ssh.DIFAT)
}

// SectorIndexTable maps sect numbers to the next sector in the chain or special value (e.g. ffffffff for empty/free sector)
// unsure of return value.
// this is a double indirect FAT lookup
func (ssh StructuredStorageHeader) SectorIndexTable() (map[SECT]uint32, error) {

	ret := map[SECT]uint32{}

	// SECT and SID are both expected to be ULONG
	FREESECTULong := ULONG(FREESECT)
	for idx, v := range ssh.DIFAT {
		vULong := ULONG(v)
		if vULong == FREESECTULong {
			continue
		}
		ret[v] = 0

		if idx == len(ssh.DIFAT)-1 && vULong != 0 {
			return ret, errors.New("sector chains for double-indirect fat not supported.  please request improvements.")
		}

	}

	return ret, nil
}
