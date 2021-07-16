package hlclock

import (
	"testing"
	"time"
)

var (
	current int64 = time.Now().Unix()
	before        = current - 1000
	after         = current + 1000
)

func TestHTimestamp_Increment(t *testing.T) {
	type fields struct {
		timestamp int64
		counter   uint16
	}

	tests := []struct {
		name   string
		fields fields
		args   int64
		want   HTimestamp
	}{
		{"PTime newer", fields{current, 0}, after, NewHTimestamp(after, 0)},
		{"PTime drift before", fields{current, 0}, before, NewHTimestamp(current, 1)},
		{"PTime equal", fields{current, 0}, current, NewHTimestamp(current, 1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ht := &HTimestamp{
				timestamp: tt.fields.timestamp,
				counter:   tt.fields.counter,
			}
			ht.Increment(tt.args)
			if got := ht.Copy(); !got.Equal(&tt.want) {
				t.Errorf("HTimestamp.Increment() = tc: %s, got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestHTimestamp_Merge(t *testing.T) {
	type fields struct {
		timestamp int64
		counter   uint16
	}
	type args struct {
		pt  int64
		msg HTimestamp
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   HTimestamp
	}{
		{"PTime newer", fields{current, 0}, args{after, NewHTimestamp(current, 1)}, NewHTimestamp(after, 0)},
		{"PTime drift before, recv msg timestamp counter newer", fields{current, 0}, args{before, NewHTimestamp(current, 1)}, NewHTimestamp(current, 2)},
		{"PTime equal, recv msg timestamp newer", fields{current, 0}, args{after, NewHTimestamp(after, 1)}, NewHTimestamp(after, 2)},
		{"PTime newer, current timestamp newer", fields{after, 0}, args{after, NewHTimestamp(current, 3)}, NewHTimestamp(after, 1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ht := &HTimestamp{
				timestamp: tt.fields.timestamp,
				counter:   tt.fields.counter,
			}
			ht.Merge(tt.args.pt, &tt.args.msg)
			if got := ht.Copy(); !got.Equal(&tt.want) {
				t.Errorf("HTimestamp.Merge() = tc: %s, got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestHTimestamp_Compare(t *testing.T) {
	type fields struct {
		timestamp int64
		counter   uint16
	}
	type args struct {
		other HTimestamp
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{"Equal", fields{current, 0}, args{NewHTimestamp(current, 0)}, 0},
		{"Gt counter", fields{current, 1}, args{NewHTimestamp(current, 0)}, 1},
		{"Gt", fields{after, 0}, args{NewHTimestamp(current, 0)}, 1},
		{"Lt", fields{before, 0}, args{NewHTimestamp(current, 0)}, -1},
		{"Lt counter", fields{current, 0}, args{NewHTimestamp(current, 1)}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ht := &HTimestamp{
				timestamp: tt.fields.timestamp,
				counter:   tt.fields.counter,
			}
			if got := ht.Compare(&tt.args.other); got != tt.want {
				t.Errorf("HTimestamp.Compare() = tc: %s, got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestHCLock_Increment(t *testing.T) {
	type fields struct {
		sysClock SysClock
		latest   *HTimestamp
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hcl := &HCLock{
				sysClock: tt.fields.sysClock,
				latest:   tt.fields.latest,
			}
			hcl.Increment()
		})
	}
}

func TestHCLock_Merge(t *testing.T) {
	type fields struct {
		sysClock SysClock
		latest   *HTimestamp
	}
	type args struct {
		e *HTimestamp
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hcl := &HCLock{
				sysClock: tt.fields.sysClock,
				latest:   tt.fields.latest,
			}
			hcl.Merge(tt.args.e)
		})
	}
}
