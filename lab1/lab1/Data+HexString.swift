//
//  Data+String.swift
//  lab1
//
//  Created by Sebastijan Zindl on 30.10.24.
//
import Foundation

extension Data {
	func toHexString() -> String {
		return map {
			String(format: "%02hhx", $0)
		}.joined()
	}
}
