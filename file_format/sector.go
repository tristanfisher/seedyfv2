package file_format

// 64-bit value representing number of 100 nanoseconds since January 1, 1601
// 	typedef struct tagFILETIME {
//		DWORD dwLowDateTime;
//		DWORD dwHighDateTime;
// } FILETIME;

//
// aliased types for ease of referencing against documentation
//

// USHORT is 16-bit / 2 byte unsigned: ushort -> uint16
type USHORT uint16

// ULONG is 64-bit unsigned: ulong -> uint64, but defined as 32-bit in the spec (4 bytes) -> uint32
type ULONG uint32 // note - per spec a 32 not 64

// CHAR is 8-bit alias: char -> byte
type CHAR byte

// WORD microsoft type.  note 32 bit is used in this file format (16-bit unsigned: word -> uint16)
type WORD uint16

// DWORD microsoft type.  double word 32-bit unsigned: dword -> uint32
type DWORD uint32

type WCHAR uint16

// SECT / SID are "ulong" types that are uint32.  These represent sectors.
type SECT ULONG
type SID ULONG

type OFFSET uint16
type FSOFFSET uint16

// /aliases

//
// Sector values are preserved as hexadecimal representations as these are used for comparison
// against values and documentation.
//

// REGSECTLow to REGSECTHigh (0x00000000 - 0xFFFFFFF9) are regular sector numbers
const REGSECTLow = SECT(0x00000000)
const REGSECTHigh SECT = 0xFFFFFFF9

// MAXREGSECT is the maximum regular sector number
// maximum directory entry ID
const MAXREGSECT = 0xFFFFFFFA

//
// reserved sectors.  not available for use to identify the location of sectors
//

const RESERVED SECT = 0xFFFFFFFB   // Reserved for future use.
const DIFSECT SECT = 0xFFFFFFFC    // denotes a DIFAT sector in a FAT
const FATSECT SECT = 0xFFFFFFFD    // denotes a FAT sector in a FAT
const ENDOFCHAIN SECT = 0xFFFFFFFE // end of a virtual stream chain
const FREESECT SID = 0xFFFFFFFF    // unallocated directory entry (2 ^ 32 - 1)
