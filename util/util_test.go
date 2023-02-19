package util

import "testing"

func Test_encrypt(t *testing.T) {
	type args struct {
		plaintext string
		key       string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "happy test",
			args: args{
				plaintext: "kingsman",
				key:       "04076d64bdb6fcf31706eea85ec98431"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// encrypt the plaintext
			ciphertext, err := Encrypt(tt.args.plaintext, tt.args.key)
			if err != nil {
				t.Errorf("encrypt() error = %v", err)
				return
			}
			t.Logf("ciphertext = %s", ciphertext)
			//// decrypt the ciphertext from previous encrypt function
			plaintext, err := Decrypt(ciphertext, tt.args.key)
			if err != nil {
				t.Errorf("encrypt() error = %v", err)
				return
			}
			t.Logf("plaintext = %s", plaintext)
			//// compare the initial plaintext with output of previous decrypt function
			if plaintext != tt.args.plaintext {
				t.Errorf("plaintext = %v, want %v", plaintext, tt.args.plaintext)
			}
			//
		})
	}
}

func Test_encryptPassword(t *testing.T) {
	type args struct {
		plaintext string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "password happy test",
			args: args{
				plaintext: "kingsman"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			hash, err := HashPassword(tt.args.plaintext)
			if err != nil {
				t.Errorf("HashPassword() error = %v", err)
				return
			}

			t.Logf("hash= %s", hash)

			err = CheckPasswordHash(tt.args.plaintext, hash)
			if err != nil {
				t.Errorf("CheckPasswordHash() error = %v", err)
				return
			}

			err = CheckPasswordHash(tt.args.plaintext+"a", hash)
			if err == nil {
				t.Errorf("CheckPasswordHash() error = %v", err)
				return
			}
		})
	}
}
