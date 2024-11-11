//
//  Data+UInt64.swift
//  lab1-ccm
//
//  Created by Sebastijan Zindl on 4.11.24.
//
import Foundation

extension Data {
	var uint64: UInt64 {
		get {
			if count >= 8 {
				return self.withUnsafeBytes {
					$0.load(as: UInt64.self)
				}
			} else {
				return (self + Data(repeating: 0, count:  8 - count)).uint64
			}
		}
	}
}
