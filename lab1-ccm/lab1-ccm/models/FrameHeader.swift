//
//  FrameHeader.swift
//  lab1-ccm
//
//  Created by Sebastijan Zindl on 4.11.24.
//

enum FrameType {
	case ip
	case arp
}

public class FrameHeader {
	let sourceMAC: String
	let destinationMAC: String
	let frameType: FrameType
	
	init(sourceMAC: String, destinationMAC: String, frameType: FrameType) {
		self.sourceMAC = sourceMAC
		self.destinationMAC = destinationMAC
		self.frameType = frameType
	}
	
	func toString() -> String {
		return "Frame Type: \(frameType), Source MAC: \(sourceMAC), Destination MAC: \(destinationMAC)"
	}
}
