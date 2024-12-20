//
//  main.swift
//  lab1
//
//  Created by Sebastijan Zindl on 30.10.24.
//

import Foundation
import CryptoKit

func main() {
	do {
		let message = try getMessage()
		
		let (key, iv) = generateKeyAndIV()
		
		let clearText = ClearTextFrame(message: message.data(using: .utf8)!)
		
		print("Original Message (Hex): \(clearText.message.toHexString())")
		print("Original Message (String): \(message)")
		print (separator: "\n")
		
		//Encrypt frame
		let encryptedFrame = try clearText.encrypt(key: Array(key), iv: iv)
		print("Sent Encrypted Data (Hex): \(encryptedFrame.encryptedData.toHexString())")
		print("Sent MIC (Hex): \(encryptedFrame.mic!.toHexString())")
		print (separator: "\n")

		
		//'Send' frame in JSON format
		let receivedFrameJSON = try encryptedFrame.sendFrame()
		print("Received Frame: \(String(decoding: receivedFrameJSON, as: UTF8.self))")
		print (separator: "\n")
		
		//'Receive' frame and decode it
		let receivedFrame = try JSONDecoder().decode(EncryptedFrame.self, from: receivedFrameJSON)
	
		//Decrypt the data and verify mic
		let decryptedFrame = try receivedFrame.decryptAndVerify(key: Array(key))
		
		let decryptedMessageHex = decryptedFrame.message.toHexString()
		let decryptedMessageString = String(decoding: decryptedFrame.message, as: UTF8.self)
		
		print("Decrypted Message (Hex): \(decryptedMessageHex)")
		print("Decrypted Message (String): \(decryptedMessageString)")

	} catch {
		print("An error occured")
	}
}


main()
