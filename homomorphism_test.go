package paillier

import (
	"math/big"
	"testing"
)

func TestPublicKey_Add(t *testing.T) {
	pk, sk, err := GenerateKeyPair(1024)
	if err != nil {
		t.Errorf("Error generating key pair")
		return
	}
	ct2, _ := pk.Encrypt(2)
	ct245, _ := pk.Encrypt(245)

	type args struct {
		ct1 *big.Int
		ct2 *big.Int
	}

	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			"invalid inputs, must return error",
			args{zero, zero},
			0,
			true,
		},
		{
			"valid inputs, must return a valid ciphertext",
			args{ct2, ct245},
			247,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pk.Add(tt.args.ct1, tt.args.ct2)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublicKey.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			// Test the homomorphic property
			sum, err := sk.Decrypt(got)
			if err != nil {
				t.Errorf("PublicKey.Add() error = invalid ciphertext generated by addition")
				return
			}
			if sum != tt.want {
				t.Errorf("PublicKey.Add() = %v, want %v", sum, tt.want)
			}
		})
	}
}

func TestPublicKey_MultPlaintext(t *testing.T) {
	pk, sk, _ := GenerateKeyPair(1024)

	ct2, _ := pk.Encrypt(2)
	ct36, _ := pk.Encrypt(36)

	type args struct {
		ct *big.Int
		pt int64
	}

	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			"invalid inputs, must return error",
			args{zero, 0},
			0,
			true,
		},
		{
			"valid inputs, must return a valid ciphertext",
			args{ct2, 2},
			4,
			false,
		},
		{
			"valid inputs, must return a valid ciphertext",
			args{ct36, 36},
			1296,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pk.MultPlaintext(tt.args.ct, tt.args.pt)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublicKey.MultPlaintext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			// Test the homomorphic property
			sum, err := sk.Decrypt(got)
			if err != nil {
				t.Errorf("PublicKey.MultPlaintext() error = invalid ciphertext generated by addition")
				return
			}
			if sum != tt.want {
				t.Errorf("PublicKey.MultPlaintext() = %v, want %v", sum, tt.want)
			}
		})
	}
}

func TestPublicKey_AddPlaintext(t *testing.T) {
	pk, sk, _ := GenerateKeyPair(1024)

	ct2, _ := pk.Encrypt(2)
	ct36, _ := pk.Encrypt(36)

	type args struct {
		ct *big.Int
		pt int64
	}

	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			"invalid inputs, must return error",
			args{zero, 0},
			0,
			true,
		},
		{
			"valid inputs, must return a valid ciphertext",
			args{ct2, 272},
			274,
			false,
		},
		{
			"valid inputs, must return a valid ciphertext",
			args{ct36, 36},
			72,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pk.AddPlaintext(tt.args.ct, tt.args.pt)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublicKey.AddPlaintext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			// Test the homomorphic property
			sum, err := sk.Decrypt(got)
			if err != nil {
				t.Errorf("PublicKey.AddPlaintext() error = invalid ciphertext generated by addition")
				return
			}
			if sum != tt.want {
				t.Errorf("PublicKey.AddPlaintext() = %v, want %v", sum, tt.want)
			}
		})
	}
}

func TestPublicKey_BatchAdd(t *testing.T) {
	pk, sk, _ := GenerateKeyPair(1024)

	ct2, _ := pk.Encrypt(2)
	ct3, _ := pk.Encrypt(3)
	ct4, _ := pk.Encrypt(4)
	ct5, _ := pk.Encrypt(5)

	cap := 1000
	ctCAPx5 := make([]*big.Int, cap)
	for i := 0; i < cap; i++ {
		ctCAPx5[i], _ = pk.Encrypt(5)
	}
	tests := []struct {
		name string
		args []*big.Int
		want int64
	}{
		{
			"sum 2..4",
			[]*big.Int{ct2, ct3, ct4},
			9,
		},
		{
			"sum 4..5",
			[]*big.Int{ct4, ct5},
			9,
		},
		{
			"sum 2..5",
			[]*big.Int{ct2, ct3, ct4, ct5},
			14,
		},
		{
			"sum trying very large number, around 2000x5 or as much as the memory allows",
			ctCAPx5,
			int64(cap * 5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pk.BatchAdd(tt.args...)

			// Test the homomorphic property
			sum, err := sk.Decrypt(got)
			if err != nil {
				t.Errorf("PublicKey.BatchAdd() error = invalid ciphertext generated by addition")
				return
			}
			if sum != tt.want {
				t.Errorf("PublicKey.BatchAdd() = %v, want %v", sum, tt.want)
			}
		})
	}
}

func TestPublicKey_Sub(t *testing.T) {
	pk, sk, _ := GenerateKeyPair(1024)

	c2, _ := pk.Encrypt(2)
	c3, _ := pk.Encrypt(3)
	c4, _ := pk.Encrypt(4)
	c5, _ := pk.Encrypt(5)
	c23578, _ := pk.Encrypt(23578)
	c115, _ := pk.Encrypt(115)

	tests := []struct {
		name string
		ct1  *big.Int
		ct2  *big.Int
		want int64
	}{
		{"3-2", c3, c2, 3 - 2},
		{"5-2", c5, c2, 5 - 2},
		{"5-3", c5, c3, 5 - 3},
		{"5-4", c5, c4, 5 - 4},
		{"5-4", c23578, c115, 23578 - 115},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pk.Sub(tt.ct1, tt.ct2)

			// Test the homomorphic property
			sum, err := sk.Decrypt(got)
			if err != nil {
				t.Errorf("PublicKey.Sub() error = invalid ciphertext generated by addition")
				return
			}
			if sum != tt.want {
				t.Errorf("PublicKey.Sub() = %v, want %v", sum, tt.want)
			}
		})
	}
}

func TestPublicKey_DivPlaintext(t *testing.T) {
	pk, sk, _ := GenerateKeyPair(1024)

	ct2, _ := pk.Encrypt(2)
	ct36, _ := pk.Encrypt(36)

	type args struct {
		ct *big.Int
		pt int64
	}

	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			"invalid inputs, must return error",
			args{zero, 0},
			0,
			true,
		},
		{
			"valid inputs, must return a valid ciphertext",
			args{ct2, 2},
			1,
			false,
		},
		{
			"valid inputs, must return a valid ciphertext",
			args{ct36, 2},
			18,
			false,
		},
		{
			"valid inputs, must return a valid ciphertext",
			args{ct36, 6},
			6,
			false,
		},
		{
			"valid inputs, must return a valid ciphertext",
			args{ct36, 5},
			7,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pk.DivPlaintext(tt.args.ct, tt.args.pt)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublicKey.DivPlaintext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			// Test the homomorphic property
			sum, err := sk.Decrypt(got)
			if err != nil {
				t.Errorf("PublicKey.DivPlaintext() error = invalid ciphertext generated by addition")
				return
			}
			if sum != tt.want {
				t.Errorf("PublicKey.MultPlaintext() = %v, want %v", sum, tt.want)
			}
		})
	}
}
