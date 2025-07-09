package header

import (
	"errors"
	"io"
	"seedyfv2/file_format"
)

const CDFV2HeaderBytes = 512

type ByteOrderError struct {
	Msg string
}

func (bo ByteOrderError) Error() string {
	return bo.Msg
}

// getHeader handles decoding CDFV2 file headers
func GetHeader(r io.Reader) (file_format.StructuredStorageHeader, error) {

	headerBytes := make([]byte, CDFV2HeaderBytes)

	bRead, err := io.ReadFull(r, headerBytes)
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}
	if bRead < CDFV2HeaderBytes {
		return file_format.StructuredStorageHeader{}, errors.New("invalid header, too short")
	}

	//fmt.Println(util.Btoh(headerBytes[0:7]))

	//classIdBytes := make([]byte)

	le := file_format.LittleEndian{}
	// try to get endianness. spec states that this should be little-endian, but confirm
	// we expect 0xFFFE for intel / apple silicon / little-endian
	// we must retrieve this before parsing other numeric fields
	byteOrder, err := le.USHORT(headerBytes[28:30])
	if err != nil {
		return file_format.StructuredStorageHeader{}, errors.Join(err, ByteOrderError{Msg: "unable to determine byte order"})
	}

	//
	minorVersion, err := le.USHORT(headerBytes[24:26])
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}

	// if DLLVersion / Major Version is 0x0003 (version 3) or 0x0004 (version 4), minor version SHOULD be 0x003E
	dllVersion, err := le.USHORT(headerBytes[26:28])
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}

	sectorShift, err := le.USHORT(headerBytes[30:32])
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}

	miniSectorShift, err := le.USHORT(headerBytes[32:34])
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}

	reserved, err := le.USHORT(headerBytes[34:36])
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}

	reserved2, err := le.ULONG(headerBytes[36:40])
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}

	cSectDir, err := le.ULONG(headerBytes[40:44])
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}

	cSectFat, err := le.FSINDEX(headerBytes[44:48])
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}

	sectDirStart, err := le.SECT(headerBytes[48:52])
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}

	signature, err := le.DFSIGNATURE(headerBytes[52:56])
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}

	miniSectorCutoff, err := le.ULONG(headerBytes[56:60])
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}

	sectMiniFatStart, err := le.SECT(headerBytes[60:64])
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}

	cSectMiniFat, err := le.FSINDEX(headerBytes[64:68])
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}

	sectDifStart, err := le.SECT(headerBytes[68:72])
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}

	cSectDif, err := le.FSINDEX(headerBytes[72:76])
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}

	sectFatEnd := 76 + 436
	sectFat, err := le.SECTFAT(headerBytes[76:sectFatEnd])
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}
	
	// xxd -c 8 my_input is useful for parsing this file
	header := file_format.StructuredStorageHeader{
		// e.g. d0cf 11e0 a1b1 1ae1
		ABSig: [8]byte(headerBytes[0:8]), // i think this will stay on the stack without needing a new reg

		// 0000 0000 0000 0000 0000 0000 0000 0000
		Clsid: file_format.GUID(headerBytes[8:24]),

		// 3e00
		MinorVersion: minorVersion,

		// 3 for 512 byte sectors, 4 for 4KB sectors
		// 0300
		DllVersion: dllVersion,

		ByteOrder: byteOrder,

		SectorShift: sectorShift,

		MiniSectorShift: miniSectorShift,

		Reserved: reserved,

		Reserved2: reserved2,

		CSectDir: cSectDir,

		CSectFat: cSectFat,

		SectDirStart: sectDirStart,

		Signature: signature,

		MiniSectorCutoff: miniSectorCutoff,

		SectMiniFatStart: sectMiniFatStart,

		CSectMiniFat: cSectMiniFat,

		SectDifStart: sectDifStart,

		CSectDif: cSectDif,

		SectFat: sectFat,
	}

	return header, nil
}
