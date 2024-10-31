//
//  util.swift
//  lab1
//
//  Created by Sebastijan Zindl on 30.10.24.
//

import Foundation
import CryptoKit

enum NotEnoughDataError: Error {
	case noMessage
}

func getMessage() throws -> String {
	print("Please enter the message you want to encrypt/decrypt: ")
	
	if let message = readLine() {
		return message
	} else {
		throw NotEnoughDataError.noMessage
	}
}

func generateKeyAndIV() -> (key: Data, iv: Data) {
	let key = SymmetricKey(size: .bits128)
	let keyData = key.withUnsafeBytes { Data(Array($0)) }
	
	print("Generated Key (Hex): \(keyData.map { String(format: "%02hhx", $0) }.joined())")
	
	var ivData = Data(count: 16)
	let _ = ivData.withUnsafeMutableBytes { SecRandomCopyBytes(kSecRandomDefault, 16, $0.baseAddress!) }
	print("Generated IV (Hex): \(ivData.map { String(format: "%02hhx", $0) }.joined())")
	print (separator: "\n")
	
	return (keyData, ivData)

}

func splitFrame(frame: Data) -> (receivedData: Data, receivedMIC: Data?) {
	let micLength = 32
	
	let encryptedData = frame.prefix(frame.count - micLength)
	let receivedMIC = frame.suffix(micLength)
	
	return (encryptedData, receivedMIC)
}

