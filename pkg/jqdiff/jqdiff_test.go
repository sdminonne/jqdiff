package jqdiff

import (
	"reflect"
	"testing"
)

func Test_Basic_Types_jqondiff_Compare(t *testing.T) {
	type args struct {
		ref    []byte
		actual []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []Diff
		wantErr bool
	}{
		{
			name: "flat strings: 1 diff",
			args: args{
				ref:    []byte(`"value1"`),
				actual: []byte(`"value2"`),
			},
			want: []Diff{typedDiff[string, string]{
				Selector:  ".",
				Actual:    "value2",
				Reference: "value1",
				Kind:      DifferentValue}},
			wantErr: false,
		},
		{
			name: "flat strings: No diff",
			args: args{
				ref:    []byte(`"value1"`),
				actual: []byte(`"value1"`),
			},
			want:    []Diff{},
			wantErr: false,
		},
		{
			name: "flat bools: No diff",
			args: args{
				ref:    []byte(`true`),
				actual: []byte(`true `),
			},
			want:    []Diff{},
			wantErr: false,
		},

		{
			name: "flat float64s: No diff",
			args: args{
				ref:    []byte(`1`),
				actual: []byte(`1`),
			},
			want:    []Diff{},
			wantErr: false,
		},
		{
			name: "flat nulls: No diff",
			args: args{
				ref:    []byte(`null`),
				actual: []byte(`null`),
			},
			want:    []Diff{},
			wantErr: false,
		},
		{
			name: "flat nulls: 1 diff",
			args: args{
				ref:    []byte(`null`),
				actual: []byte(`true`),
			},
			want: []Diff{typedDiff[any, any]{
				Selector:  ".",
				Reference: nil,
				Actual:    true,
				Kind:      DifferentType,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := NewJqdiff()
			got, err := j.Compare(tt.args.ref, tt.args.actual)
			if (err != nil) != tt.wantErr {
				t.Errorf("jqdiff.Compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jqdiff.Compare() got = %#v\n want = %#v", got, tt.want)
			}
		})
	}
}

func Test_Composed_Types_jqondiff_Compare(t *testing.T) {
	type args struct {
		ref    []byte
		actual []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []Diff
		wantErr bool
	}{
		{
			name: "string in object: 1 diff",
			args: args{
				ref:    []byte(`{ "key": "value1" }`),
				actual: []byte(`{ "key": "value2" }`),
			},
			want: []Diff{typedDiff[string, string]{
				Selector:  ".key",
				Reference: "value1",
				Actual:    "value2",
				Kind:      DifferentValue}},
			wantErr: false,
		},
		{
			name: "string in object No diff",
			args: args{
				ref:    []byte(`{ "key1": "true" }`),
				actual: []byte(`{ "key1": "true" }`),
			},
			want:    []Diff{},
			wantErr: false,
		},
		{
			name: "bool in object No diff",
			args: args{
				ref:    []byte(`{ "key1": true }`),
				actual: []byte(`{ "key1": true }`),
			},
			want:    []Diff{},
			wantErr: false,
		},

		{
			name: "float64 in object No diff",
			args: args{
				ref:    []byte(`{ "key1": 1 }`),
				actual: []byte(`{ "key1": 1 }`),
			},
			want:    []Diff{},
			wantErr: false,
		},
		{
			name: "null in object No diff",
			args: args{
				ref:    []byte(`{ "key1": null }`),
				actual: []byte(`{ "key1": null }`),
			},
			want:    []Diff{},
			wantErr: false,
		},
		{
			name: "null in object 1 diff",
			args: args{
				ref:    []byte(`{ "key1": null }`),
				actual: []byte(`{ "key1": true }`),
			},
			want: []Diff{typedDiff[any, any]{
				Selector:  ".key1",
				Reference: nil,
				Actual:    true,
				Kind:      DifferentType,
			}},
			wantErr: false,
		},

		{
			name: ".key.subkey No diff",
			args: args{
				ref:    []byte(`{ "key": { "subkey" : true } }`),
				actual: []byte(`{ "key": { "subkey" : true } }`),
			},
			want:    []Diff{},
			wantErr: false,
		},

		{
			name: ".key.subkey 1 diff",
			args: args{
				ref:    []byte(`{ "key": { "subkey" : true } }`),
				actual: []byte(`{ "key": { "subkey" : false } }`),
			},
			want: []Diff{typedDiff[bool, bool]{
				Selector:  ".key.subkey",
				Actual:    false,
				Reference: true,
				Kind:      DifferentValue,
			}},
			wantErr: false,
		},
		{
			name: ".key.subkey.sub-subkey 1 diff",
			args: args{
				ref:    []byte(`{ "key": { "subkey" : { "sub-subkey": true } } }`),
				actual: []byte(`{ "key": { "subkey" : { "sub-subkey": false } } }`),
			},
			want: []Diff{typedDiff[bool, bool]{
				Selector:  ".key.subkey.sub-subkey",
				Actual:    false,
				Reference: true,
				Kind:      DifferentValue,
			}},
			wantErr: false,
		},
		{
			name: "array: no diff",
			args: args{
				ref:    []byte(`["one", "two", "three"]`),
				actual: []byte(`["one", "two", "three"]`),
			},
			want:    []Diff{},
			wantErr: false,
		},
		{
			name: "array: 1 diff",
			args: args{
				ref:    []byte(`["one", "two", "three"]`),
				actual: []byte(`["one", "due", "three"]`),
			},
			want: []Diff{
				typedDiff[string, string]{
					Selector:  ".[1]",
					Actual:    "due",
					Reference: "two",
					Kind:      DifferentValue,
				}},
			wantErr: false,
		},
		{
			name: "array: 2 diff",
			args: args{
				ref:    []byte(`["one", "two", "three", "four", "five", "six"]`),
				actual: []byte(`["one", "due", "three", "four", "cinq", "six"]`),
			},
			want: []Diff{
				typedDiff[string, string]{
					Selector:  ".[1]",
					Actual:    "due",
					Reference: "two",
					Kind:      DifferentValue,
				},
				typedDiff[string, string]{
					Selector:  ".[4]",
					Actual:    "cinq",
					Reference: "five",
					Kind:      DifferentValue,
				}},
			wantErr: false,
		},
		{
			name: "array of arrays: no diff",
			args: args{
				ref:    []byte(`[ ["one", "two"] , ["three", "four"], ["five", "six"] ]`),
				actual: []byte(`[ ["one", "two"] , ["three", "four"], ["five", "six"] ]`),
			},
			want:    []Diff{},
			wantErr: false,
		},
		{
			name: "array of arrays: 1 diff",
			args: args{
				ref:    []byte(`[ ["one", "two"] , ["three", "four"], ["five", "six"] ]`),
				actual: []byte(`[ ["one", "two"] , ["trois", "four"], ["five", "six"] ]`),
			},
			want: []Diff{
				typedDiff[string, string]{
					Selector:  ".[1][0]",
					Actual:    "trois",
					Reference: "three",
					Kind:      DifferentValue,
				}},
			wantErr: false,
		},
		{
			name: "array of arrays: 3 diff",
			args: args{
				ref:    []byte(`[ ["one", "two"] , ["three", "four"], ["five", "six"] ]`),
				actual: []byte(`[ ["one", "two"] , ["trois", "four"], ["six", "cinq"] ]`),
			},
			want: []Diff{
				typedDiff[string, string]{
					Selector:  ".[1][0]",
					Actual:    "trois",
					Reference: "three",
					Kind:      DifferentValue,
				},
				typedDiff[string, string]{
					Selector:  ".[2][0]",
					Actual:    "six",
					Reference: "five",
					Kind:      DifferentValue,
				},
				typedDiff[string, string]{
					Selector:  ".[2][1]",
					Actual:    "cinq",
					Reference: "six",
					Kind:      DifferentValue,
				},
			},
			wantErr: false,
		},

		// '[true, "mixed array", 12, { "key": null }, ["not", "mixed", "array"] ]'
		// '{ "key": [true, null, "hello", false, "world" ] }'
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := NewJqdiff()
			got, err := j.Compare(tt.args.ref, tt.args.actual)
			if (err != nil) != tt.wantErr {
				t.Errorf("jqdiff.Compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jqdiff.Compare()\n got -> %#v\n want -> %#v", got, tt.want)
			}
		})
	}
}
