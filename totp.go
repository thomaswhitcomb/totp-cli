package main

import "os"
import "crypto/hmac"
import "crypto/sha1"
import "encoding/base32"
import "fmt"
import "time"

// convert the interval to base 256 and stick
// each number 0-ff in its own byte.
func interval_to_counter(num int64) []byte {
	bytes := make([]byte, 8)
	for i := 7; num > 0; i-- {
		rem := num % 256
		bytes[i] = byte(rem)
		num = num / 256
	}
	return bytes
}

// Convert the input string to base 32
func to_base32(s string) ([]byte, error) {
	data, err := base32.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Create a SHA1 hash of the counter using the input secret
func generate_hash(secret []byte, counter []byte) []byte {
	mac := hmac.New(sha1.New, secret)
	mac.Write(counter)
	return mac.Sum(nil)
}

func code_from_chunk(chunk []byte) int {
	var i int = 0
	var j int = 0
	shifts := []uint{24, 16, 8, 0}
	for k := 0; k < len(shifts); k++ {
		j = int(chunk[k])
		j = j << shifts[k]
		i = i | j
	}

	i = i % 1000000
	return i
}
func counter_from_time(unixtime int64, period int) ([]byte, int) {
	var period64 int64 = int64(period)
	intervals := unixtime / period64
	remaining_seconds := unixtime - (intervals * period64)
	bytes := interval_to_counter(intervals)
	return bytes, int(remaining_seconds)
}
func chunk_from_hash(hmac []byte) []byte {
	var b byte = hmac[len(hmac)-1]
	offset := int(b & 0x0F)
	chunk := hmac[offset:(offset + 4)]
	chunk[0] = chunk[0] & 0x7f
	return chunk
}

func format_code(code int) string {
	// pad the code out with zeros to 6 chars
	return fmt.Sprintf("%06d", code)
}

func main() {
	var period int = 30
	if len(os.Args) < 2 {
		fmt.Println("Missing secret")
		os.Exit(1)
	}
	if len(os.Args[1]) < 16 {
		fmt.Println("Secret must be equal to or longer than 16 bytes")
		os.Exit(1)
	}

	var secret string = os.Args[1]

	b32, err := to_base32(secret)
	if err != nil {
		fmt.Printf("Error converting secret to base32: %v\n", err)
	}
	now := time.Now().Unix()
	counter, remaining_time := counter_from_time(now, period)
	hash := generate_hash(b32, counter)
	chunk := chunk_from_hash(hash)
	code := code_from_chunk(chunk)
	formatted_code := format_code(code)
	fmt.Println(formatted_code, period-remaining_time)
}
