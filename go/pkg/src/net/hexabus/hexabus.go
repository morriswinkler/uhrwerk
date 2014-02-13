package hexabus

// Constants
const (
      /* Header */

      // The UDP Data of a Hexabus Packet starts with the Bytes 0x48 0x58 0x30 0x43 
      // (HX0C) to identify it as a Hexabus Packet
      HXB_HEADER = [0x48, 0x58, 0x30, 0x43]

      /* Boolean values */      

      // boolean false
      HXB_FALSE = 0
      
      // boolean true
      HXB_TRUE = 1

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
      HXB_DTYPE_NONE = 0x00
      
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

type HXB_ErrorPacket_ST stuct {
     header [4]byte   // header "HX0C"
     type byte	      // type 
     flags byte	      // flags 
     error_code byte  // error code
}

func HXB_ErrorPacket(err byte) *HXB_ErrorPacket_ST {
     return &HXB_ErrorPacket_ST{
				header: [0x48, 0x58, 0x30, 0x43]
				type: 0x00
				flags: 0x00
				error_code: err 
}

type HXB_InfoPacket stuct {
     header byte[4]   // header "HX0C"
     type byte	      // type 
     flags byte	      // flags 
     eid byte[3]      // endpoint id
     dtype byte	      // data type
     data  byte[]     // payload, size depending on datatype
}

type HXB_QueryPacket stuct {
     header byte[4]   // header "HX0C"
     type byte	      // type 
     flags byte	      // flags 
     eid byte[3]      // endpoint id
}

type HXB_WritePacket stuct {
     header byte[4]   // header "HX0C"
     type byte	      // type 
     flags byte	      // flags 
     eid byte[3]      // endpoint id
     dtype byte	      // data type
     data  byte[]     // payload, size depending on datatype
}     