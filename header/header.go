package header

import (
	"errors"
	"io"
	"seedyfv2/file_format"
)

const CDFV2HeaderBytes = 512

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
	// try to get endiannesss
	// spec states that this should be little-endian, but confirm
	// 1C => 28
	// we expect 0xFFFE for intel / apple silicon / little-endian
	byteOrder, err := le.SHORT(headerBytes[28:30])
	if err != nil {
		return file_format.StructuredStorageHeader{}, err
	}
	
	header := file_format.StructuredStorageHeader{
		// 0->8
		// d0cf 11e0 a1b1 1ae1
		ABSig: [8]byte(headerBytes[0:8]), // i think this will stay on the stack without needing a new reg
		// 8->16
		// 0000 0000 0000 0000 0000 0000 0000 0000
		Clsid: file_format.CLSID{},
		//
		// 3e00 - 0 padded
		MinorVersion: 0,
		// 0300
		DllVersion:       0,
		ByteOrder:        byteOrder,
		SectorShift:      0,
		MiniSectorShift:  0,
		Reserved:         0,
		Reserved2:        0,
		CSectDir:         0,
		CSectFat:         0,
		SectDirStart:     0,
		Signature:        0,
		MiniSectorCutoff: 0,
		SectMiniFatStart: 0,
		CSectMiniFat:     0,
		SectDifStart:     0,
		CSectDif:         0,
		SectFat:          [109]file_format.SECT{},
	}

	return header, nil
}
