package main

import "testing"
import "bytes"

func TestCounterCreation1(t *testing.T) {
	b := interval_to_counter(30)
	if b[7] != 0x1e {
		t.Error("Expected 0x1e but got:", b[7])
	}
}
func TestCounterCreation2(t *testing.T) {
	b := interval_to_counter(255)
	if b[7] != 0xff {
		t.Error("Expected 0xff but got:", b[7])
	}
}
func TestCounterCreation3(t *testing.T) {
	b := interval_to_counter(1025)
	if b[7] != 0x01 || b[6] != 0x04 {
		t.Error("Expected 0x04 and 0x01 but got:", b[7], b[6])
	}
}

func TestToBase32_1(t *testing.T) {
	b32, err := to_base32("ABCDEFGHIJKLMNOP")
	if err != nil || !bytes.Equal(b32, []byte{0, 68, 50, 20, 199, 66, 84, 182, 53, 207}) {
		t.Errorf("Expected good conversion. got %v", err)
	}
}
func TestGenerateHash1(t *testing.T) {
	secret, _ := to_base32("ABCDEFGHIJKLMNOP")
	b := interval_to_counter(1025)
	hash := generate_hash(secret, b)
	if !bytes.Equal(hash, []byte{207, 6, 67, 115, 30, 56, 191, 55, 60, 97, 113, 240, 48, 124, 85, 120, 193, 99, 170, 82}) {
		t.Errorf("Expected good hash. got %v", hash)
	}
}
func TestCodeFromChunk1(t *testing.T) {
	code := code_from_chunk([]byte{9, 181, 88, 199})
	if code != 879687 {
		t.Errorf("1) Bad code from chunk. got %v", code)
	}
}

// This test is from RFC 4226
func TestCodeFromChunk2(t *testing.T) {
	code := code_from_chunk([]byte{0x50, 0xef, 0x7f, 0x19})
	if code != 872921 {
		t.Errorf("2) Bad code from chunk. got %v", code)
	}
}

func TestChunkFromHash1(t *testing.T) {
	chunk := chunk_from_hash([]byte{237, 121, 86, 241, 173, 223, 106, 81, 252, 233, 111, 140, 26, 124, 77, 117, 209, 122, 142, 221})
	if !bytes.Equal(chunk, []byte{124, 77, 117, 209}) {
		t.Errorf("1) Bad chunk from hash. got %v", chunk)
	}
}

// This test is from RFC4226
func TestChunkFromHash2(t *testing.T) {
	chunk := chunk_from_hash([]byte{0x1f, 0x86, 0x98, 0x69, 0x0e, 0x02, 0xca, 0x16, 0x61, 0x85, 0x50, 0xef, 0x7f, 0x19, 0xda, 0x8e, 0x94, 0x5b, 0x55, 0x5a})
	if !bytes.Equal(chunk, []byte{0x50, 0xef, 0x7f, 0x19}) {
		t.Errorf("2) Bad chunk from hash. got %v", chunk)
	}
}

func TestFormatCode(t *testing.T) {
	fs := format_code(67676)
	if fs != "067676" {
		t.Errorf("Formatting failed. got %v", fs)
	}
}
func TestCounterFromTime1(t *testing.T) {
	counter, remaining := counter_from_time(0, 30)
	if !bytes.Equal(counter, []byte{0, 0, 0, 0, 0, 0, 0, 0}) {
		t.Errorf("1) Bad counter from time. got %v", counter)
	}
	if remaining != 0 {
		t.Errorf("1) Bad remaining from time. got %v", remaining)
	}
}
func TestCounterFromTime2(t *testing.T) {
	counter, remaining := counter_from_time(1494368473, 30)
	if !bytes.Equal(counter, []byte{0, 0, 0, 0, 2, 248, 19, 58}) {
		t.Errorf("2) Bad counter from time. got %v", counter)
	}
	if remaining != 13 {
		t.Errorf("2) Bad remaining from time. got %v", remaining)
	}
}
