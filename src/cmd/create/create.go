package create

import (
	"fmt"
	"github.com/iocage/libiocage/copyFile"
	"github.com/iocage/libiocage/set/jailProp"
	"github.com/nu7hatch/gouuid"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// create accepts the zpool, the path and any number of properties.
// It then creates a jail with those properties set.
func Args(Pool, Iocroot *string, props []string) {
	var hname, huuid, t bool
	var c, ntag string
	var cn, i int

	if len(props) < 0 {
		if !strings.Contains(props[0], "=") {
			fmt.Printf("Property %s has invalid syntax!\n", props[0])
			os.Exit(1)
		}
	}

	// A quick check to see if they supplied us with base.
	m := make(map[string]string)
	for _, v := range props {
		p := strings.Split(v, "=")
		prop, val := p[0], p[1]
		m[prop] = val
	}

	// Check if any of these properties exist, as the users choice takes
	// preference.
	switch {
	case m["base"] == "":
		fmt.Println("You must supply the 'base' property.")
		os.Exit(1)
	case m["host_hostname"] != "":
		hname = true
		fallthrough
	case m["host_hostuuid"] != "":
		huuid = true
		fallthrough
	case m["tag"] != "":
		t = true
		ntag = m["tag"]
		fallthrough
	case m["count"] != "":
		c = m["count"]
		cn, _ = strconv.Atoi(c)
	}

	// This is so a user can supply 'count=N' and they get that many jails with
	// unique tags and host_hostnames. i needs some massaging if N > 1.
	if cn != 0 {
		i = 1
	}

	for i <= cn {
		uuid4, _ := uuid.NewV4()
		uuidstr := uuid4.String()

		fs := *Pool + "/ioc/bases/" + m["base"]

		// Snapshot the root dataset we need for the jail creation.
		_, err := exec.Command("/sbin/zfs", "snapshot", fs+"/root"+"@"+
			"jail_"+uuidstr).Output()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Clone the root dataset, using '-p' will create an independent parent
		// dataset above it.
		_, err = exec.Command("/sbin/zfs", "clone", "-p",
			fs+"/root"+"@"+"jail_"+uuidstr,
			*Pool+"/ioc/jails/"+uuidstr+"/root").Output()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Use the defaults file as a template, then we set whatever props the
		// user wanted on it.
		copyFile.Args(*Iocroot+"/.default", *Iocroot+"/jails/"+uuidstr+"/config")
		path := *Iocroot + "/jails/" + uuidstr

		// Set this to a sane default before the user properties.
		jailProp.Args(path, "jail_zfs_dataset", path+"/data", false)

		// Check to see if the user supplied these critical properties that
		// we would otherwise overwrite.
		switch {
		case !hname:
			jailProp.Args(path, "host_hostname", uuidstr, false)
			fallthrough
		case !huuid:
			jailProp.Args(path, "host_hostuuid", uuidstr, false)
			fallthrough
		case !t:
			jailProp.Args(path, "tag", "$(date +%F@%T)", false)
		}

		// Now we set the properties the user wants.
		for name, value := range m {
			// Adding the _ to the user supplied tags and host_hostname.
			// Otherwise we set their properties.
			switch {
			case cn >= 2 && name == "tag":
				value = value + "_" + strconv.Itoa(i)
				ntag = value
				jailProp.Args(path, name, value, false)
			case cn >= 2 && name == "host_hostname":
				value = value + "_" + strconv.Itoa(i)
				jailProp.Args(path, name, value, false)
			default:
				// We finally set the remaining props.
				jailProp.Args(path, name, value, false)
			}
		}

		fmt.Printf("Successfully created: %s (%s)\n", uuidstr, ntag)
		i++
	}
}
