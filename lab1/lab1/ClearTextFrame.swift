//
//  ClearTextFrame.swift
//  lab1
//
//  Created by Sebastijan Zindl on 30.10.24.
//


import Foundation
import CryptoSwift

class ClearTextFrame {
	
	let message: Data
	
	init(message: Data) {
		self.message = message
	}
	
	func encrypt(key: Array<UInt8>, iv: Data) throws -> EncryptedFrame {
		let clearTextBytes = Array(message)
		
		let aes = try AES(key: key, blockMode: CTR(iv: Array(iv)), padding: .pkcs7)
		
		let encryptedBytes = try aes.encrypt(clearTextBytes)
		
		let encryptedData = Data(encryptedBytes)
		
		
		return EncryptedFrame(iv: iv, encryptedData: encryptedData, key:key)
	
	}
}
