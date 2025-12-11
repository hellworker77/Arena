use std::io::{self, Read};
use byteorder::{LittleEndian, ReadBytesExt};

#[derive(Debug, Clone, Copy)]
pub enum PacketType {
    PTInput = 1,
    PTSnapshot = 2,
    PTReliableCmd = 3,
    PTAuth = 10,
    PTAuthResp = 11,
}

impl TryFrom<u8> for PacketType {
    type Error = ();
    
    fn try_from(value: u8) -> Result<Self, Self::Error> {
        match value {
            1 => Ok(PacketType::PTInput),
            2 => Ok(PacketType::PTSnapshot),
            3 => Ok(PacketType::PTReliableCmd),
            10 => Ok(PacketType::PTAuth),
            11 => Ok(PacketType::PTAuthResp),
            _ => Err(()),
        }
    }
}

#[derive(Debug, Clone)]
pub struct PacketHeader {
    pub ver: u8,
    pub ptype: u8,
    pub connection: u32,
    pub seq: u32,
    pub ack_latest: u32,
    pub ack_bitmap: u64,
}

impl PacketHeader {
    pub fn new(ptype: PacketType) -> Self {
        PacketHeader {
            ver: 1,
            ptype: ptype as u8,
            connection: 0,
            seq: 0,
            ack_latest: 0,
            ack_bitmap: 0,
        }
    }
}

pub fn write_packet(buf: &mut Vec<u8>, header: &PacketHeader, payload: &[u8]) {
    buf.push(header.ver);
    buf.push(header.ptype);
    buf.extend_from_slice(&header.connection.to_le_bytes());
    buf.extend_from_slice(&header.seq.to_le_bytes());
    buf.extend_from_slice(&header.ack_latest.to_le_bytes());
    buf.extend_from_slice(&header.ack_bitmap.to_le_bytes());
    buf.extend_from_slice(payload);
}

pub fn parse_packet(data: &[u8]) -> io::Result<(PacketHeader, Vec<u8>)> {
    let mut rdr = io::Cursor::new(data);

    let ver = rdr.read_u8()?;
    let ptype = rdr.read_u8()?;
    let connection = rdr.read_u32::<LittleEndian>()?;
    let seq = rdr.read_u32::<LittleEndian>()?;
    let ack_latest = rdr.read_u32::<LittleEndian>()?;
    let ack_bitmap = rdr.read_u64::<LittleEndian>()?;

    let mut body = Vec::new();
    rdr.read_to_end(&mut body)?;

    Ok((
        PacketHeader {
            ver,
            ptype,
            connection,
            seq,
            ack_latest,
            ack_bitmap,
        },
        body,
    ))
}