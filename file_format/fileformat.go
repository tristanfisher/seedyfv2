package file_format

import (
	"encoding/binary"
	"errors"
	"log"
)

// LittleEndian as a struct exists to (ab)use a structs to create a namespace/static methods
type LittleEndian struct {
	ErrorLog *log.Logger
}

func (le LittleEndian) SHORT(b []byte) (USHORT, error) {
	if len(b) < 2 {
		return 0, errors.New("not enough bytes provided for u conversion")
	}
	return USHORT(binary.LittleEndian.Uint16(b)), nil
}

// StructuredStorageHeader contains the information
// required for instantiating and parsing a compound file.
//
// Comments from:
//
//	Advanced Authoring Format (AAF)
//	Low-Level Container Specification v1.0.1
//
// [offset from start (bytes), length (bytes)]

/*
Sector types

typedef unsigned long ULONG  // 4 bytes
typedef unsigned short USHORT // 2 bytes
typedef short OFFSET // 2 bytes
typedef ULONG SECT // 4 bytes
typedef ULONG FSINDEX // 4 bytes
typedef USHORT FSOFFSET // 2 bytes
typedef USHORT WCHAR // 2 bytes
typedef ULONG DFSIGNATURE // 4 bytes
typedef unsigned char BYTE // 1 byte
typedef unsigned short WORD // 2 bytes
typedef unsigned long DWORD // 4 bytes
typedef ULONG SID // 4 bytes
typedef GUID CLSID // 16 bytes

// 64-bit value representing number of 100 nanoseconds since January 1, 1601
typedef struct tagFILETIME {
DWORD dwLowDateTime;
DWORD dwHighDateTime;
} FILETIME;


const SECT MAXREGSECT = 0xFFFFFFFA; // maximum SECT
const SECT DIFSECT = 0xFFFFFFFC; // denotes a DIFAT sector in a FAT
const SECT FATSECT = 0xFFFFFFFD; // denotes a FAT sector in a FAT
const SECT ENDOFCHAIN = 0xFFFFFFFE; // end of a virtual stream chain
const SECT FREESECT = 0xFFFFFFFF; // unallocated sector
const SID MAXREGSID = 0xFFFFFFFA; // maximum directory entry ID
const SID NOSTREAM = 0xFFFFFFFF; // unallocated directory entry
*/

// c -> go numeric types
// 16-bit unsigned: ushort -> uint16
// 64-bit unsigned: ulong -> uint64, but defined as 32-bit in the spec (4 bytes) -> uint32
// 8-bit alias: char -> byte
//
// microsoft types
// https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-dtyp/432cd7e7-6276-4c82-87fc-50bbcbd5ffa0
//
// double word 32-bit unsigned: dword -> uint32
// word 16-bit unsigned: word -> uint16
// GUID -> 16 byte structure with different representations
// 	https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-dtyp/4926e530-816e-41c2-b251-ec5c7aca018a
// SECT / SID are "ulong" types that are uint32

//
// aliased types for ease of referencing against documentation
//

type USHORT uint16

type ULONG uint32 // note - per spec a 32 not 64
type CHAR byte
type FSINDEX ULONG
type SECT ULONG
type DFSIGNATURE ULONG

type CLSID [16]byte

type StructuredStorageHeader struct {
	// [00H,08] {0xd0, 0xcf, 0x11, 0xe0, 0xa1, 0xb1,
	// 0x1a, 0xe1} for current version
	//ABSig [8]byte
	ABSig [8]byte

	// [08H,16] reserved must be zero (WriteClassStg/
	// GetClassFile uses root directory class id)
	Clsid CLSID

	// [18H,02] minor version of the format: 33 is
	// written by reference implementation
	MinorVersion USHORT

	// [1AH,02] major version of the dll/format: 3 for
	// 512-byte sectors, 4 for 4 KB sectors
	DllVersion USHORT

	// [1CH,02] 0xFFFE: indicates Intel byte-ordering
	ByteOrder USHORT

	// [1EH,02] size of sectors in power-of-two;
	// typically 9 indicating 512-byte sectors
	SectorShift USHORT

	// [20H,02] size of mini-sectors in power-of-two;
	// typically 6 indicating 64-byte mini-sectors
	MiniSectorShift USHORT

	// [22H,02] reserved, must be zero
	Reserved USHORT

	// [24H,04] reserved, must be zero
	Reserved2 ULONG

	// [28H,04] must be zero for 512-byte sectors,
	// number of SECTs in directory chain for 4 KB
	// sectors
	CSectDir ULONG

	// [2CH,04] number of SECTs in the FAT chain
	CSectFat FSINDEX

	// [30H,04] first SECT in the directory chain
	SectDirStart SECT

	// [34H,04] signature used for transactions; must
	// be zero. The reference implementation
	// does not support transactions
	Signature DFSIGNATURE

	// [38H,04] maximum size for a mini stream;
	// typically 4096 bytes
	MiniSectorCutoff ULONG

	// [3CH,04] first SECT in the MiniFAT chain
	SectMiniFatStart SECT

	// [40H,04] number of SECTs in the MiniFAT chain
	CSectMiniFat FSINDEX

	// [44H,04] first SECT in the DIFAT chain
	SectDifStart SECT

	// [48H,04] number of SECTs in the DIFAT chain
	CSectDif FSINDEX

	// [4CH,436] the SECTs of first 109 FAT sectors
	SectFat [109]SECT
}
