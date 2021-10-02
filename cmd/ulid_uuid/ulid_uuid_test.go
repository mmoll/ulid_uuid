package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_checkType(t *testing.T) {
	type args struct {
		ulidUUID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"convert ulid",
			args{"08A1YW3WAH8SNTQVYGDB2EP69T"},
			"08507dc1-f151-466b-abef-d06ac4eb193a",
			false,
		},
		{
			"convert uuid",
			args{"cfa45f5d-9c38-4772-b39a-036a0b9f8d30"},
			"6FMHFNV71R8XSB76G3D85SZ39G",
			false,
		},
		{
			"ulid out of range",
			args{"FFMHFNV71R8XSB76G3D85SZ39G"},
			"",
			true,
		},
		{
			"invalid uuid",
			args{"cfa45f5k-9c38-4772-!39a-036a0b9f8d30"},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkType(tt.args.ulidUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mainControl(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name       string
		args       args
		want       int
		wantOutput string
	}{
		{
			"call program without any parameter",
			args{[]string{"ulid_uuid"}},
			1,
			"Please give one Parameter to convert!",
		},
		{
			"convert failed",
			args{[]string{"ulid_uuid", "F8A1YW3WAH8SNTQVYGDB2EP69T"}},
			1,
			"not valid ULID|UUID|GUID",
		},
		{
			"display help",
			args{[]string{"ulid_uuid", "-h"}},
			0,
			"ulid_uuid [-hn] [UUID|GUID|ULID]",
		},
		{
			"get uuid without newline",
			args{[]string{"ulid_uuid", "-n", "08A1YW3WAH8SNTQVYGDB2EP69T"}},
			0,
			"08507dc1-f151-466b-abef-d06ac4eb193a",
		},
		{
			"get uuid with newline",
			args{[]string{"ulid_uuid", "08A1YW3WAH8SNTQVYGDB2EP69T"}},
			0,
			"08507dc1-f151-466b-abef-d06ac4eb193a\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			reader, writer, err := os.Pipe()
			if err != nil {
				panic(err)
			}
			origStdout := os.Stdout
			origStderr := os.Stderr
			os.Stdout = writer
			os.Stderr = writer
			if got := mainControl(tt.args.args); got != tt.want {
				t.Errorf("mainControl() = %v, want %v", got, tt.want)
			}
			_ = writer.Close()
			_, _ = io.Copy(buf, reader)
			_ = reader.Close()
			os.Stdout = origStdout
			os.Stderr = origStderr
			output := buf.String()

			if !strings.Contains(output, tt.wantOutput) {
				t.Errorf("mainControl() got output = %v, want part output %v", output, tt.wantOutput)
			}
		})
	}
}

func Test_toUUID(t *testing.T) {
	type args struct {
		ulidString string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"to uuid failed",
			args{"ADSJFKEWIFJWFEW"},
			"",
			true,
		},
		{
			"to uuid sucess",
			args{"08A1YW3WAH8SNTQVYGDB2EP69T"},
			"08507dc1-f151-466b-abef-d06ac4eb193a",
			false,
		},
		{
			"",
			args{"FFMHFNV71R8XSB76G3D85SZ39G"},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toUUID(tt.args.ulidString)
			if (err != nil) != tt.wantErr {
				t.Errorf("toUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("toUUID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toUlid(t *testing.T) {
	type args struct {
		uuidString string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"to ulid failed",
			args{"jfkldsfaj"},
			"",
			true,
		},
		{
			"to ulid sucess",
			args{"cfa45f5d-9c38-4772-b39a-036a0b9f8d30"},
			"6FMHFNV71R8XSB76G3D85SZ39G",
			false,
		},
		{
			"to ulid failed unmarshal wrong uuid",
			args{"cfa45f5k-9c38-4772-!39a-036a0b9f8d30"},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toUlid(tt.args.uuidString)
			if (err != nil) != tt.wantErr {
				t.Errorf("toUlid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("toUlid() got = %v, want %v", got, tt.want)
			}
		})
	}
}
