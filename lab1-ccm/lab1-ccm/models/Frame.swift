//
//  Frame.swift
//  lab1-ccm
//
//  Created by Sebastijan Zindl on 4.11.24.
//
import Foundation

class Frame {
	public let PN: UInt64 //48-bit packet Number
	public let Header: FrameHeader
	public let MIC: Array<UInt8>
	public let IV: Array<UInt8>
	
	init(PN: UInt64, Header: FrameHeader, MIC: Array<UInt8>, IV: Array<UInt8>) {
		self.PN = PN
		self.Header = Header
		self.MIC = MIC
		self.IV = IV
	}
}
