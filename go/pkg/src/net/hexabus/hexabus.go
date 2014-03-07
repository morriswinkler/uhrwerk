package hexabus

import "fmt"
import "bytes"
import "encoding/binary"


// Constants
const (
	/* Header */

	// The UDP Data of a Hexabus Packet starts with the Bytes 0x48 0x58 0x30 0x43 
	// (HX0C) to identify it as a Hexabus Packet
	
	HXB_HEADER0 = 0x48
	HXB_HEADER1 = 0x58
	HXB_HEADER2 = 0x30
	HXB_HEADER3 = 0x43

	/* Boolean values */      

	// boolean false
	HXB_FALSE = 0x00
	
	// boolean true
	HXB_TRUE = 0x01

	/* Packet types */
	
	// Hexabus Error Packet 
	// An error occured -- check the error code field for more information
	HXB_PTYPE_ERROR = 0x00
	
	// Hexabus Info Packet
	// Endpoint provides information
	HXB_PTYPE_INFO = 0x01

	// Hexabus Query Packet
	// Endpoint is requested to provide information
	HXB_PTYPE_QUERY = 0x02
	
	// Hexabus Write Packet
	// Endpoint is requested to set its value
	HXB_PTYPE_WRITE = 0x04

	// Hexabus EpInfo Packet
	// Endpoint metadata
	HXB_PTYPE_EPINFO = 0x09

	// Hexabus EpQuery Packet
	// Request endpoint metadata
	HXB_PTYPE_EPQUERY = 0x09

	/* Flags */

	// Hexabus Flag "No Flag set" 
	HXB_FLAG_NONE = 0x00

	/* Data types */
	
	// Hexabus Data Type "No data at all"
	HXB_DTYPE_UNDEFINED = 0x00
	
	// Data type Bool
	HXB_DTYPE_BOOL = 0x01

	// Data type uint8
	HXB_DTYPE_UINT8 = 0x02

	// Data type uint32
	HXB_DTYPE_UINT32 = 0x03
	
	// Data type date/time 
	HXB_DTYPE_DATETIME = 0x04

	// Data type float
	HXB_DTYPE_FLOAT = 0x05
	
	// Data type 128String 
	// char string with 128 bytes, must be 0 terminated 
	HXB_DTYPE_128STRING = 0x06
	// max char length = 127 + 0 termination
	HXB_STRING_PACKET_MAX_BUFFER_LENGTH = 127

	// Data type timestamp
	// in secondes since device was booted up, 32 bit unsigned integer (4 bytes)
	HXB_DTYPE_TIMESTAMP = 0x07

	// Data type 66bytes 
	// 66 bytes of raw binary data
	HXB_DTYPE_66BYTES = 0x08

	// Data type 16bytes
	// 16 bytes of raw binary data
	HXB_DTYPE_16BYTES = 0x09

	/* Error codes */

	// reserved: No error
	HXB_ERR_SUCESS = 0x00

	// A request for an endpoint which does not exist on the device was received
	HXB_ERR_UNKNOWNEID = 0x01
	
	// WRITE was received for a readonly endpoint
	HXB_ERR_WRITEREADONLY = 0x02

	// A packet failed the CRC check 
	// TODO How can we find out what information was lost?
	HXB_ERR_CRCFAILED = 0x03

	// A packet with a datatype that does not fit the endpoint was received
	HXB_ERR_DATATYPE = 0x04

	// A value was encountered that cannot be interpreted
	HXB_ERR_INVALID_VALUE = 0x05
)

func addHeader(packet []byte) {
	packet[0], packet[1], packet[2], packet[3] = HXB_HEADER0, HXB_HEADER1, HXB_HEADER2, HXB_HEADER3
}

func addData(packet []byte, data interface{}) []byte { 	
	// set datatype
	switch data := data.(type) {
	case bool:
		packet[10] = HXB_DTYPE_BOOL
		if data == true {
			packet = append(packet, HXB_TRUE)
		} else {
			packet = append(packet, HXB_FALSE)
		}
	case uint8:
		packet[10] = HXB_DTYPE_UINT8
		packet = append(packet, data)
	case uint32:
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, data)
		if err != nil {
			panic(fmt.Errorf("binary.Write failed:", err))
		}
		packet[10] = HXB_DTYPE_UINT32
		packet = append(packet, buf.Bytes()...)
	// DateTime: holds HXB_DTYPE_DATETIME data ; needs testing
	case DateTime:
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, data)
		if err != nil {
			panic(fmt.Errorf("binary.Write failed:", err))
		}
		packet[10] = HXB_DTYPE_DATETIME
		packet = append(packet, buf.Bytes()...)
	case float32:
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, data)
		if err != nil {
			panic(fmt.Errorf("binary.Write failed:", err))
		}
		packet[10] = HXB_DTYPE_FLOAT
		packet = append(packet, buf.Bytes()...)
	case string:
		if len(data) > HXB_STRING_PACKET_MAX_BUFFER_LENGTH { 
			panic(fmt.Errorf("max string length 127 exeeded for string: %s", data))
		} else {
			// TODO: check if 0 termination in string is right that way
			packet[10] = HXB_DTYPE_128STRING
			packet = append(packet, data...)
			packet = append(packet, byte(0))
		}
		// TIMESTAMP: intended for type syscall.Sysinfo_t.Uptime not working
	case Timestamp:
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, data)
		if err != nil {
			panic(fmt.Errorf("binary.Write failed:", err))
		}
		packet[10] = HXB_DTYPE_TIMESTAMP
		packet = append(packet, buf.Bytes()...)
	case []byte:
		// TODO: check if padding is the intended behavior
		if len(data) == 16 {
			packet[10] = HXB_DTYPE_16BYTES
			packet = append(packet, data...)
		} else if len(data) == 66 {
			packet[10] = HXB_DTYPE_66BYTES
			packet = append(packet, data...)
		} else {
			panic(fmt.Errorf("only 16 or 66 bytes of raw data are allowed length %d is not supported", len(data)))
		}
	default:
		packet[10] = HXB_DTYPE_UNDEFINED
		panic(fmt.Errorf("unsupported payload type: %T", data))
	}
	
	return packet
	
}

// calculate crc16 kermit variant
// this code was translated from a php snippet found on http://www.lammertbies.nl/forum/viewtopic.php?t=1253
func crc16_KERMIT(packet []byte) uint16 {
	var crc uint16
        for _, v := range packet {
                crc = crc ^ uint16(v)
                for y := 0; y < 8; y++ {
                        if (crc & 0x001) == 0x0001 {
                                crc = (crc >> 1) ^ 0x8408
                        } else {
                                crc = crc >> 1
                        }
                }
        }
	// in the original Kermit implementation the two crc bytes are swaped, looks like boost::crc_optimal<16, 0x1021, 0x0000, 0, true, true> doesn't follow that ?
        //lb := (crc & 0xff00) >> 8
        //hb := (crc & 0x00ff) << 8
        //crc = hb | lb
	return crc
} 

func addCRC(packet []byte) []byte {
	crc := crc16_KERMIT(packet)
	packet = append(packet,uint8(crc>>8), uint8(crc&0xff))
	return packet
}

// struct to hold HXB_DTYPE_TIMESTAMP
type Timestamp struct {
	TotalSeconds uint32
}

func (t *Timestamp) Decode(data interface{}) {
        buf := bytes.NewBuffer(data.([]byte))
        err := binary.Read(buf, binary.BigEndian, t)
        if err != nil {
                panic(fmt.Errorf("binary.Write failed:", err))
        }
}

// struct to hold HXB_DTYPE_DATETIME data
type DateTime struct {
        Hours uint8
        Minutes uint8
        Seconds uint8
        Day uint8
        Month uint8
        Year uint16
        DayOfWeek uint8
}

func (d *DateTime) Decode(data interface{}) {
        buf := bytes.NewBuffer(data.([]byte))
        err := binary.Read(buf, binary.BigEndian, d)
        if err != nil {
                panic(fmt.Errorf("binary.Write failed:", err))
        }
}                

type ErrorPacket struct {
	// 4 bytes header
	// 1 byte packet type
	Flags byte	 // 1 byte flags 
	Error byte       // 1 byte error code
}

func (p *ErrorPacket) Encode() []byte {
	packet := make([]byte, 7)
	addHeader(packet)
	packet[4] = HXB_PTYPE_ERROR 
	packet[5] = p.Flags
	packet[6] = p.Error
	packet = addCRC(packet)
	return packet
}

func (p *ErrorPacket) Decode(packet []byte) {
	p.Flags = packet[5]
	p.Error = packet[6]
}

type InfoPacket struct {
	// 4 bytes header
	// 1 byte packet type
	Flags byte	  // 1 byteflags 
	Eid uint32  // 4 bytes endpoint id
	Dtype byte	  // 1 byte data type
	Data interface{}   // ... bytes payload, size depending on datatype
}

func (p *InfoPacket) Encode() []byte {
	packet := make([]byte, 11)
	addHeader(packet)
	packet[4] = HXB_PTYPE_INFO
	packet[5] = p.Flags
	packet[6], packet[7], packet[8], packet[9] = uint8(p.Eid>>24), uint8(p.Eid>>16), uint8(p.Eid>>8), uint8(p.Eid&0xff)
	fmt.Printf("EID bits: %b, %b, %b, %b\n", packet[6], packet[7], packet[8], packet[9])	    
	packet[10] = p.Dtype
	packet = addData(packet, p.Data)
	packet = addCRC(packet)     
	return packet 
}

func (p *InfoPacket) Decode(packet []byte) {
	p.Flags = packet[5]
	p.Eid = uint32(uint8(packet[6])>>24 + uint8(packet[7])>>16 + uint8(packet[8])>>8 + uint8(packet[9])&0xff)
	p.Dtype = packet[10]
	p.Data = packet[11:len(packet)-2]
}

// remove that !!!
func EncodeInfoPacket( flags byte, eid uint32, data interface{}) (p []byte) {
	packet := make([]byte, 11, 141)                                                
	addHeader(packet)
        packet[4] = HXB_PTYPE_INFO
        packet[5] = flags
        packet[6], packet[7], packet[8], packet[9] = uint8(eid>>24), uint8(eid>>16), uint8(eid>>8), uint8(eid&0xff)
	fmt.Printf("EID bits: %b, %b, %b, %b\n", packet[6], packet[7], packet[8], packet[9])
	packet = addData(packet, data)                                                
        packet = addCRC(packet)                                                         
        return packet
}
	
type QueryPacket struct {
	// 4 bytes header
	// 1 byte packet type
	Flags byte	 // flags 
	Eid uint32       // endpoint id
}

func (p *QueryPacket) Encode() []byte {
	packet := make([]byte, 12)
	addHeader(packet)
	packet[4] = HXB_PTYPE_QUERY
	packet[5] = p.Flags
	packet[6], packet[7], packet[8], packet[9] = uint8(p.Eid>>24), uint8(p.Eid>>16), uint8(p.Eid>>8), uint8(p.Eid&0xff)
	packet = addCRC(packet)     
	return packet 
}

func (p  *QueryPacket) Decode(packet []byte) {
	p.Flags = packet[5]
	p.Eid = uint32(uint8(packet[6])>>24 + uint8(packet[7])>>16 + uint8(packet[8])>>8 + uint8(packet[9])&0xff)
}

type WritePacket struct {
	// 4 bytes header
	// 1 byte packet type
	Flags byte       // flags 
	Eid uint32       // endpoint id
	Dtype byte	 // data type
	Data interface{} // payload, size depending on datatype
}     

func (p *WritePacket) Encode() []byte {
	packet := make([]byte, 11)
	addHeader(packet)
	packet[4] = HXB_PTYPE_WRITE
	packet[5] = p.Flags
	packet[6], packet[7], packet[8], packet[9] = uint8(p.Eid>>24), uint8(p.Eid>>16), uint8(p.Eid>>8), uint8(p.Eid&0xff)
	packet = addData(packet, p.Data)                                                
	packet = addCRC(packet)     
	return packet 
}

func (p *WritePacket) Decode(packet []byte) {
	p.Flags = packet[5]
	p.Eid = uint32(uint8(packet[6])>>24 + uint8(packet[7])>>16 + uint8(packet[8])>>8 + uint8(packet[9])&0xff)
	p.Dtype = packet[10]
	p.Data = packet[11:len(packet)-2]

}
