package vmrc

import (
	"context"
	"fmt"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"
)

// VMRC ..
func VMRC(datacenter, s, user, password string) ([]string, []string, error) {
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
	vms, err := f.VirtualMachineList(ctx, "*")
	if err != nil {
	}
	var refs []*object.VirtualMachine
	var refsName []string
	for _, vm := range vms {
		refs = append(refs, vm)
		refsName = append(refsName, vm.Name())
	}
	var links []string
	__ := []string{""}
	for i := 0; i < len(refs); i++ {
		vm := refs[i]
		state, err := vm.PowerState(ctx)
		if err != nil {
			return __, __, err
		}

		if state != types.VirtualMachinePowerStatePoweredOn {
			//fmt.Errorf("vm is not powered on (%s)", state)
		}

		c := vm.Client()

		u := c.URL()
		req := types.AcquireCloneTicket{
			This: *c.ServiceContent.SessionManager,
		}

		res, err := methods.AcquireCloneTicket(ctx, c, &req)
		if err != nil {
			return __, __, err
		}

		ticket := res.Returnval

		var link string
		h5 := false

		if h5 {
			m := object.NewOptionManager(c, *c.ServiceContent.Setting)

			opt, err := m.Query(ctx, "VirtualCenter.FQDN")
			if err != nil {
				return __, __, err
			}

			fqdn := opt[0].GetOptionValue().Value.(string)

			var info object.HostCertificateInfo
			_ = info.FromURL(u, nil)

			u.Path = "/ui/webconsole.html"

			u.RawQuery = url.Values{
				"vmId":          []string{vm.Reference().Value},
				"vmName":        []string{vm.Name()},
				"serverGuid":    []string{c.ServiceContent.About.InstanceUuid},
				"host":          []string{fqdn},
				"sessionTicket": []string{ticket},
				"thumbprint":    []string{info.ThumbprintSHA1},
			}.Encode()

			link = u.String()
		} else {

			link = fmt.Sprintf("vmrc://clone:%s@%s/?moid=%s", ticket, s, vm.Reference().Value)
		}

		//fmt.Fprintln(os.Stdout, link)
		//fmt.Printf("%+v\n%+v\n", refs, refsName)
		links = append(links, link)
	}
	return refsName, links, err
}
