/*
	Copyright © 2023 Gabriel Pozo

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/
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
