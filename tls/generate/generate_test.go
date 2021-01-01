// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/31

package generate

import "testing"

func TestDefaultCertificate(t *testing.T) {
	tlsConfig, err := DefaultCertificate("cert.pem", "key.pem")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", tlsConfig)
}
