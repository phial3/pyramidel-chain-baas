package configtxgen

import (
	"reflect"
	"testing"
)

func TestReadTemplateFile(t *testing.T) {
	tests := []struct {
		name    string
		want    *TopLevel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadTemplateFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadTemplateFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadTemplateFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
