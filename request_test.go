package nap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildQueryParamUrl1(t *testing.T) {
	type args struct {
		fullPath    string
		queryStruct interface{}
	}
	tests := []struct {
		name     string
		args     args
		expected string
		wantErr  bool
	}{
		{
			name: "Test case 1 - Parse URL with struct request for GET method Successfully",
			args: args{
				fullPath: "https://github.com/v2/create",
				queryStruct: struct {
					AppID      int32  `url:"app_id"`
					AppTransID string `url:"app_trans_id"`
					Mac        string `url:"mac"`
				}{
					AppID:      1,
					AppTransID: "mmf_transid_210822001",
					Mac:        "MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAPl9eHJltu48w1P",
				},
			},
			expected: "https://github.com/v2/create?app_id=1&app_trans_id=mmf_transid_210822001&mac=MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAPl9eHJltu48w1P",
			wantErr:  false,
		},
		{
			name: "Test case 2 - Request is not a struct, should return error",
			args: args{
				fullPath:    "https://sb-openapi.zalopay.vn/v2/create",
				queryStruct: 6,
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "Test case 3 - Full path is invalid, should return error",
			args: args{
				fullPath:    "github.com/%zzzzz",
				queryStruct: 6,
			},
			expected: "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlStr, err := buildQueryParamUrl(tt.args.fullPath, tt.args.queryStruct)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, urlStr)
				return
			}
			assert.Equal(t, tt.expected, urlStr)
		})
	}
}
