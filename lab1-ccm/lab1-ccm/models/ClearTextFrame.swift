//
//  ClearTextFrame.swift
//  lab1-ccm
//
//  Created by Sebastijan Zindl on 4.11.24.
//
import Foundation

class ClearTextFrame: Frame {
	let clearText: String
	var blocks: Array<[UInt8]> = []
	
	init(PN: UInt64, Header: FrameHeader, MIC: Array<UInt8>, IV: Array<UInt8>, clearText: String) {
		self.clearText = clearText
		super.init(PN: PN, Header: Header, MIC: MIC, IV: IV)
		self.blocks = padding()
	}
	
	private func padding() -> Array<[UInt8]> {
		let messageBytes = Data(clearText.utf8)
		let blockSize = 16
		let paddingSize = blockSize - (messageBytes.count % blockSize)
		
		var paddedBytes: Array<UInt8> = Array(repeating: UInt8.max, count: paddingSize)
		paddedBytes.replaceSubrange(0..<messageBytes.count, with: messageBytes)
		
		for index in messageBytes.count..<paddedBytes.count {
			paddedBytes[index] = UInt8(paddingSize)
		}
		
		var blocks: Array<[UInt8]> = []
		for i in stride(from: 0, to: paddedBytes.count, by: blockSize) {
			let end = min(i + blockSize, paddedBytes.count)
			let block = Array(paddedBytes[i..<end])
			blocks.append(block)
		}
		return blocks
	}
	
	private func calculateMIC() -> Void{
		let packetNumber = self.PN
		let sourceMac = self.Header.sourceMAC
		let sourceMacBytes = Data(sourceMac.utf8)
		let destinationMac = self.Header.destinationMAC
		let desintationMacBytes = Data(destinationMac.utf8)
		let qosPriority = UInt64(4)
		
		let nonce: UInt64 = (packetNumber << 56) | (sourceMacBytes.uint64 << 8) | qosPriority;
	}
	
}
