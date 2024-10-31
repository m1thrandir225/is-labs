//
//  EncryptedFrame.swift
//  lab1
//
//  Created by Sebastijan Zindl on 30.10.24.
//
import Foundation
import CryptoSwift

public class EncryptedFrame : Codable {
	public let iv: Array<UInt8>
	public let encryptedData: Data
	public var mic: Data?
	
	init(iv: Data, encryptedData: Data, key: Array<UInt8>) {
		self.iv = Array(iv)
		self.encryptedData = encryptedData
		calculateMIC(key: key)
	}
	
	private func calculateMIC(key: Array<UInt8>) {
		let hmac = HMAC(key: key, variant: .sha2(.sha256))
		
		let macBytes = try! hmac.authenticate(Array(encryptedData))
		
		mic = Data(macBytes)
	}

	func sendFrame()  throws -> Data {
		let encoder = JSONEncoder()
		var frame = encryptedData
		frame.append(mic ?? Data())
		
		let json = try encoder.encode(frame)
		
		return json
	}
	

	func decryptAndVerify(
		key: Array<UInt8>,
		receivedData: Data,
		receivedMIC: Data
	) throws -> ClearTextFrame {
		let receivedBytes = Array(receivedData)
		
		let aes = try AES(key: key, blockMode: CTR(iv: iv), padding: .pkcs7)
		let decryptedBytes = try aes.decrypt(receivedBytes)
		
		let decryptedData = Data(decryptedBytes)
		
		let calculatedMIC = try HMAC(key: key, variant: .sha2(.sha256)).authenticate(receivedBytes)
		
		guard Data(calculatedMIC) == receivedMIC else {
			throw CCMError.integrityCheckFailed
		}
		
		return ClearTextFrame(message: decryptedData)
	}
}

enum CCMError: Error {
	case integrityCheckFailed
}
