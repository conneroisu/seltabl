package http

import "testing"

// TestParseAddress tests the ParseAddress function by asserting that it returns
// the correct protocol and path for the given address.
func TestParseAddress(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name         string
		args         args
		wantProtocol string
		wantPath     string
		wantErr      bool
	}{
		{
			name:         "parse ipv4 hello world",
			args:         args{addr: "ipv4://127.0.0.1:80/hello"},
			wantProtocol: "",
			wantPath:     "",
			wantErr:      true,
		},
		{
			name:         "parse tcp hello world uppercase",
			args:         args{addr: "TCP://127.0.0.1:80/hello"},
			wantProtocol: "TCP",
			wantPath:     "127.0.0.1:80/hello",
			wantErr:      false,
		},
		{
			name:         "parse tcp hello world lowercase",
			args:         args{addr: "tcp://127.0.0.1:80/hello"},
			wantProtocol: "tcp",
			wantPath:     "127.0.0.1:80/hello",
			wantErr:      false,
		},
		{
			name:         "parse udp hello world",
			args:         args{addr: "udp://127.0.0.1:443/secure"},
			wantProtocol: "udp",
			wantPath:     "127.0.0.1:443/secure",
			wantErr:      false,
		},
		{
			name:         "parse unix hello world",
			args:         args{addr: "unix:///var/run/rexray/rexray.sock"},
			wantProtocol: "unix",
			wantPath:     "/var/run/rexray/rexray.sock",
			wantErr:      false,
		},
		{
			name:         "parse unixgram hello world",
			args:         args{addr: "unixgram:///var/run/rexray/rexray.sock"},
			wantProtocol: "unixgram",
			wantPath:     "/var/run/rexray/rexray.sock",
			wantErr:      false,
		},
		{
			name: "unixpacket hello world",
			args: args{
				addr: "unixpacket:///var/run/rexray/rexray.sock",
			},
			wantProtocol: "unixpacket",
			wantPath:     "/var/run/rexray/rexray.sock",
			wantErr:      false,
		},
		{
			name:         "ip hello world",
			args:         args{addr: "ip://127.0.0.1:443/hello"},
			wantProtocol: "ip",
			wantPath:     "127.0.0.1:443/hello",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			gotProto, gotPath, err := ParseAddress(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"ParseAddress() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if gotProto != tt.wantProtocol {
				t.Errorf(
					"ParseAddress() gotProto = %v, want %v",
					gotProto,
					tt.wantProtocol,
				)
			}
			if gotPath != tt.wantPath {
				t.Errorf(
					"ParseAddress() gotPath = %v, want %v",
					gotPath,
					tt.wantPath,
				)
			}
		})
	}
}
