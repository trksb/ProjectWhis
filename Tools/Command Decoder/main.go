package main

import (
	"bufio"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
)

const (
	_Delta = 0x9e3779b9
)

func main() {
	//Encoded String -> Base64Decode -> XXTeaDecrypt -> String -> Plain Text
	fmt.Println("Project Whis Command Decoder")
MENU:
	fmt.Print("Encoded String => ")
	Encoded := bufio.NewScanner(os.Stdin)
	Encoded.Scan()

	fmt.Print("Passcode => ")
	Passcode := bufio.NewScanner(os.Stdin)
	Passcode.Scan()

	decoded, _ := base64.RawURLEncoding.DecodeString(Encoded.Text())
	Decrypted := XXTeaDecrypt(decoded, []byte(Passcode.Text()))
	fmt.Println("|---------------Decoded---------------|")
	fmt.Println(string(Decrypted))
	fmt.Println("|-------------------------------------|")

	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	goto MENU
}

func Base64Decode(str string) string {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(data)
}

func XXTeaDecrypt(data []byte, key []byte) []byte {
	if data == nil || key == nil || len(data) == 0 || len(key) == 0 {
		return nil
	}
	if len(data)%4 != 0 {
		return nil
	}
	uint32Arr := asUint32Array(data, false)
	decryptedData := btea(uint32Arr, -len(uint32Arr), asKey(key))
	return asByteArray(decryptedData, true)
}

func mx(z, y, sum, p, e uint32, key []uint32) uint32 {
	return ((z>>5 ^ y<<2) + (y>>3 ^ z<<4)) ^ ((sum ^ y) + (key[(p&3)^e] ^ z))
}

func btea(v []uint32, n int, key []uint32) []uint32 {
	var y, z, sum uint32
	var p, rounds, e uint32

	if n > 1 {
		rounds = uint32(6 + 52/n)
		sum = 0
		z = v[n-1]
		for i := uint32(0); i < rounds; i++ {
			sum += _Delta
			e = (sum >> 2) & 3
			for p = 0; p < uint32(n-1); p++ {
				y = v[p+1]
				z = v[p] + mx(z, y, sum, p, e, key)
				v[p] = z
			}
			y = v[0]
			z = v[p] + mx(z, y, sum, p, e, key)
			v[p] = z
		}
	} else if n < -1 {
		n = -n
		rounds = uint32(6 + 52/n)
		sum = rounds * _Delta
		y = v[0]
		for i := uint32(0); i < rounds; i++ {
			e = (sum >> 2) & 3
			for p = uint32(n - 1); p > 0; p-- {
				z = v[p-1]
				y = v[p] - mx(z, y, sum, p, e, key)
				v[p] = y
			}
			z = v[n-1]
			y = v[0] - mx(z, y, sum, p, e, key)
			v[0] = y
			sum -= _Delta
		}
	}

	return v
}

func asKey(key []byte) []uint32 {
	if len(key) > 16 {
		key = key[:16]
	} else if len(key) < 16 {
		padding := make([]byte, 16-len(key))
		key = append(key, padding...)
	}
	return asUint32Array(key, false)
}

func asByteArray(data []uint32, includeLength bool) []byte {
	var result []byte
	dataLen := uint32(len(data))
	size := dataLen << 2
	if includeLength {
		lastByte := data[len(data)-1]
		if lastByte > (size-4) || lastByte < (size-7) {
			return nil
		}
		size = lastByte
		dataLen--
		if size%4 != 0 {
			result = make([]byte, ((size/4)+1)*4)
		} else {
			result = make([]byte, size)
		}
	} else {
		result = make([]byte, size)
	}
	for idx := uint32(0); idx < dataLen; idx++ {
		binary.LittleEndian.PutUint32(result[idx*4:(idx+1)*4], data[idx])
	}
	return result[:size]
}

func toUint32(b []byte) uint32 {
	switch len(b) {
	case 0:
		return uint32(0)
	case 1:
		return uint32(b[0])
	case 2:
		return uint32(b[0]) | uint32(b[1])<<8
	case 3:
		return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
	default:
		return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
	}
}

func asUint32Array(data []byte, includeLength bool) []uint32 {
	var uint32Arr []uint32
	size := uint32(len(data) / 4)
	if len(data)&3 != 0 {
		size++
	}
	if includeLength {
		uint32Arr = make([]uint32, size+1)
		uint32Arr[size] = uint32(len(data))
	} else {
		uint32Arr = make([]uint32, size)
	}
	for idx := uint32(0); idx < size; idx++ {
		uint32Arr[idx] = toUint32(data[idx*4:])
	}

	return uint32Arr
}
