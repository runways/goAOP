package unitTests

import "testing"

func TestFirstStruct_InvokeFirstFunction(t *testing.T) {
	type fields struct {
		name string
		Age  int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Normal test",
			fields: struct {
				name string
				Age  int
			}{name: "", Age: 0},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := FirstStruct{
				name: tt.fields.name,
				Age:  tt.fields.Age,
			}
			if got := fs.InvokeFirstFunction(); got != tt.want {
				t.Errorf("InvokeFirstFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}
