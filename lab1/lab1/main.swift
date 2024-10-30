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
		
		let encryptedFrame = EncryptedFrame(key: key, iv: iv)
		try encryptedFrame.encrypt(clearTextFrame: clearText)

		let (encryptedData, mic) = encryptedFrame.sendFrame()
		print("Sent Encrypted Data (Hex): \(encryptedData.toHexString())")
		print("Sent MIC (Hex): \(mic.toHexString())")

		let receivedData = encryptedData
		let receivedMIC = mic

		let decryptedMessageData = try encryptedFrame.decryptAndVerify(receivedData: receivedData, receivedMIC: receivedMIC, iv: Array(iv))
		let decryptedMessageHex = decryptedMessageData.toHexString()
		let decryptedMessageString = String(decoding: decryptedMessageData, as: UTF8.self)
		
		print("Decrypted Message (Hex): \(decryptedMessageHex)")
		print("Decrypted Message (String): \(decryptedMessageString)")

	} catch {
		print("An error occured")
	}
}


main()
