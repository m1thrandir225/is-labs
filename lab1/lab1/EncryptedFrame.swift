//
//  EncryptedFrame.swift
//  lab1
//
//  Created by Sebastijan Zindl on 30.10.24.
//
import Foundation
import CryptoSwift

class EncryptedFrame {
	private let key: Array<UInt8>
	private let iv: Array<UInt8>
	private(set) var encryptedData: Data = Data()
	private var mic: Data?
	
	init(key: Data, iv: Data) {
		self.key = Array(key)
		self.iv = Array(iv)
	}
	
	func encrypt(clearTextFrame: ClearTextFrame) throws {
		let clearTextBytes = Array(clearTextFrame.message)
		
		let aes = try AES(key: key, blockMode: CTR(iv: iv),padding: .noPadding)
		
		let encryptedBytes = try aes.encrypt(clearTextBytes)
		encryptedData = Data(encryptedBytes)
		
		calculateMIC()
		
	}
	
	private func calculateMIC() {
		let hmac = HMAC(key: key, variant: .sha256)
		
		let macBytes = try! hmac.authenticate(Array(encryptedData))
		
		mic = Data(macBytes)
	}
	
	func sendFrame() -> (Data, Data) {
		return (encryptedData, mic ?? Data())
	}
	
	func decryptAndVerify(receivedData: Data, receivedMIC: Data) throws -> Data {
		let receivedBytes = Array(receivedData)
		
		let aes = try AES(key: key, blockMode: CTR(iv: iv), padding: .noPadding)
		let decryptedBytes = try aes.decrypt(receivedBytes)
		
		let decryptedData = Data(decryptedBytes)
		
		let calculatedMIC = try HMAC(key: key, variant: .sha256).authenticate(receivedBytes)
		
		guard Data(calculatedMIC) == receivedMIC else {
			throw CCMError.integrityCheckFailed
		}
		
		return decryptedData
	}
}

enum CCMError: Error {
	case integrityCheckFailed
}
