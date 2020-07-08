package helpers

import "testing"

func Test_convertSQLDayToDOW(t *testing.T) {
	type args struct {
		daynum int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Monday",
			args: args{
				daynum: 1,
			},
			want: "Monday",
		},
		{
			name: "OOB",
			args: args{
				daynum: 7,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertSQLDayToDOW(tt.args.daynum); got != tt.want {
				t.Errorf("convertSQLDayToDOW() = %v, want %v", got, tt.want)
			}
		})
	}
}
