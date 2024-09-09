package pkg

import (
	"encoding/json"
	"testing"
)

func TestToInt64E(t *testing.T) {
	type args struct {
		v interface{}
	}
	n := int64(6911300862917002766)
	o := struct{ A *int64 }{}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{args: args{"6911300862917002766"}, want: 6911300862917002766},
		{args: args{int64(6911300862917002766)}, want: 6911300862917002766},
		{args: args{&n}, want: 6911300862917002766},
		{args: args{120}, want: 120},
		{args: args{uint64(121)}, want: 121},
		{args: args{122.}, want: 122},
		{args: args{o.A}, want: 0, wantErr: true},
		{args: args{"120.0"}, want: 120},
		{args: args{"-120."}, want: -120},
		{args: args{""}, want: 0, wantErr: true},
		{args: args{[]byte("125.")}, want: 125},
		{args: args{[]byte("125")}, want: 125},
		{args: args{true}, want: 1},
		{args: args{false}, want: 0},
		{args: args{nil}, want: 0, wantErr: true},
		{args: args{json.Number("123")}, want: 123},
		{args: args{json.Number("123e")}, want: 0, wantErr: true},
		{args: args{json.Number("123e.")}, want: 0, wantErr: true},
		{args: args{json.Number("123.1")}, want: 123},
		{args: args{json.Number("")}, want: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt64(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt64E() error = %v, wantErr %v,v %v", err, tt.wantErr, tt.args.v)
				return
			}
			if got != tt.want {
				t.Errorf("ToInt64E() got = %v, want %v,v %v", got, tt.want, tt.args.v)
			}
		})
	}
}

func TestToFloat64E(t *testing.T) {
	type args struct {
		v interface{}
	}
	n := float64(6911300862917002766)
	o := struct{ A *float64 }{}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{args: args{"6911300862917002766"}, want: 6911300862917002766},
		{args: args{int64(6911300862917002766)}, want: 6911300862917002766},
		{args: args{&n}, want: 6911300862917002766},
		{args: args{120}, want: 120},
		{args: args{uint64(121)}, want: 121},
		{args: args{122.}, want: 122},
		{args: args{o.A}, want: 0, wantErr: true},
		{args: args{"120.0"}, want: 120},
		{args: args{"-120."}, want: -120},
		{args: args{""}, want: 0, wantErr: true},
		{args: args{[]byte("125.")}, want: 125},
		{args: args{[]byte("125")}, want: 125},
		{args: args{true}, want: 1},
		{args: args{false}, want: 0},
		{args: args{nil}, want: 0, wantErr: true},
		{args: args{json.Number("123")}, want: 123},
		{args: args{json.Number("123e")}, want: 0, wantErr: true},
		{args: args{json.Number("123e.")}, want: 0, wantErr: true},
		{args: args{json.Number("123.1")}, want: 123.1},
		{args: args{json.Number("")}, want: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToFloat64(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt64E() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToInt64E() got = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestToFloat64(t *testing.T) {
	a := struct {
		A *int64
	}{}
	n1 := int64(6911300862917002766)
	n2 := "6911300862917002766"
	cse := []interface{}{
		int64(6911300862917002766), // ok
		"6911300862917002766",      // ok
		"6911300862917002.766",     // ok
		nil,                        // error
		"",                         // error
		&n1,                        // error
		&n2,                        // error
		a.A,                        // error
		json.Number("123"),         // ok
		json.Number("123e"),        // error
		json.Number("123e."),       // error
		json.Number("123."),        // ok
		json.Number("123.1"),       // ok
		json.Number(""),            // error
	}
	for i, n := range cse {
		f, err := ToFloat64(n)
		t.Logf("i:%d,n:%#v=>%v:%v", i, n, f, err)
	}
}

func TestToInt64(t *testing.T) {
	a := struct {
		A *int64
	}{}
	n1 := int64(6911300862917002766)
	n2 := "6911300862917002766"
	cse := []interface{}{
		int64(6911300862917002766), // ok
		"6911300862917002766",      // ok
		"6911300862917002.766",     // ok
		nil,                        // error
		"",                         // error
		&n1,                        // error
		&n2,                        // error
		a.A,                        // error
		json.Number("123"),         // ok
		json.Number("123e"),        // error
		json.Number("123e."),       // error
		json.Number("123."),        // ok
		json.Number("123.1"),       // ok
		json.Number(""),            // error
	}
	for i, n := range cse {
		in, err := ToInt64(n)
		t.Logf("i:%d,n:%#v=>%v,%v", i, n, in, err)
	}
}

func TestByteToStr(t *testing.T) {
	cse := []string{
		"", " ", "a", "hello world",
	}
	for _, s := range cse {
		strToByte := StrToByte(s)
		byteToStr := ByteToStr(strToByte)
		t.Logf("StrToByte:%v,ByteToStr:%v", strToByte, byteToStr)
	}
}
