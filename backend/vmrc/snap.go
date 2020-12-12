package vmrc

import (
	"context"
	"fmt"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
)

// Snap ..
func Snap(s, user, password, datacenter, device string) error {
	ur, err := url.Parse(s + "/sdk")
	ur.Scheme = "https"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ct, err := govmomi.NewClient(ctx, ur, true)

	//i := &url.Userinfo{}
	i := url.UserPassword(user, password)
	err = ct.Login(ctx, i)

	f := find.NewFinder(ct.Client, true)

	// Find one and only datacenter
	dc, err := f.Datacenter(ctx, datacenter)
	if err != nil {
		fmt.Println(err)
	}

	// Make future calls local to this datacenter
	f.SetDatacenter(dc)

	// Find virtual machines in datacenter
	vms, err := f.VirtualMachineList(ctx, device)
	if err != nil {
	}
	var refs []*object.VirtualMachine
	var refsName []string
	for _, vm := range vms {
		refs = append(refs, vm)
		refsName = append(refsName, vm.Name())
	}
	task, err := vms[0].RevertToSnapshot(ctx, "CLEAR", false)
	fmt.Println(user, task)
	return err
}
