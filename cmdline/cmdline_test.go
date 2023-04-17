// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmdline

import (
	"io"
	"strings"
	"testing"
)

func TestCmdline(t *testing.T) {
	exampleCmdLine := `BOOT_IMAGE=/vmlinuz-4.11.2 ro ` +
		`test-flag test2-flag=8 ` +
		`uroot.initflags="systemd test-flag=3  test2-flag runlevel=2" ` +
		`root=LABEL=/ biosdevname=0 net.ifnames=0 fsck.repair=yes ` +
		`ipv6.autoconf=0 erst_disable nox2apic crashkernel=128M ` +
		`systemd.unified_cgroup_hierarchy=1 cgroup_no_v1=all console=tty0 ` +
		`console=ttyS0,115200 security=selinux selinux=1 enforcing=0`

	c := parse(strings.NewReader(exampleCmdLine))
	wantLen := len(exampleCmdLine)
	if len(c.Raw) != wantLen {
		t.Errorf("c.Raw wrong length: %v != %d", len(c.Raw), wantLen)
	}

	if len(c.AsMap) != 21 {
		t.Errorf("c.AsMap wrong length: %v != 21", len(c.AsMap))
	}

	if c.ContainsFlag("biosdevname") == false {
		t.Errorf("couldn't find biosdevname in kernel flags: map is %v", c.AsMap)
	}

	if c.ContainsFlag("fsck.repair") == false {
		t.Error("could find fsck.repair in kernel flags, but should")
	}

	if c.ContainsFlag("biosname") == true {
		t.Error("could find biosname in kernel flags, but shouldn't")
	}

	if security, present := c.Flag("security"); !present || security != "selinux" {
		t.Errorf("Flag 'security' is %v instead of 'selinux'", security)
	}

	if c.AsBool("enforcing") {
		t.Errorf("enforcing should be casted to false")
	}
	if !c.AsBool("fsck.repair") {
		t.Errorf("fsck.repair should be casted to true")
	}
}

type badreader struct{}

// Read implements io.Reader, always returning io.ErrClosedPipe
func (*badreader) Read([]byte) (int, error) {
	// Interesting. If you return a -1 for the length,
	// it tickles a bug in io.ReadAll. It uses the returned
	// length BEFORE seeing if there was an error.
	// Note to self: file an issue on Go.
	return 0, io.ErrClosedPipe
}

func TestBadRead(t *testing.T) {
	if err := parse(&badreader{}); err == nil {
		t.Errorf("parse(&badreader{}): got nil, want %v", io.ErrClosedPipe)
	}
}
