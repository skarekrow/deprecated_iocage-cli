package main

import (
	"cmd/activate"
	"cmd/clean"
	"cmd/create"
	"cmd/deactivate"
	"cmd/destroy"
	"cmd/fetch"
	"cmd/get"
	"cmd/list"
	"cmd/set"
	"fmt"
	"github.com/iocage/libiocage/checkDatasets"
	"github.com/iocage/libiocage/rootCmds"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

type DEFAULTLOC struct {
	pool, iocroot string
}

var defloc DEFAULTLOC

func init() {
	var match int
	var d string

	u, _ := user.Current()

	if u.Uid != "0" && len(os.Args) >= 2 {
		cmd := rootCmds.Args(os.Args[1])

		if cmd {
			fmt.Printf("The %s command needs root credentials!\n", os.Args[1])
			os.Exit(1)
		}
	}

	pools, _ := exec.Command("/sbin/zpool", "list", "-H",
		"-o", "name").Output()
	str := strings.Split(strings.TrimSpace(string(pools)), "\n")

	for _, z := range str {
		dataset, _ := exec.Command("/sbin/zfs", "get", "-H", "-o", "value",
			"org.freebsd.ioc:active", z).Output()
		d = strings.TrimSpace(string(dataset))

		if d == "yes" {
			d = z // z is the zpool
			match++
		}
	}

	if match == 1 {
		defloc = DEFAULTLOC{pool: d, iocroot: "/ioc"}
	} else if match >= 2 {
		if len(os.Args) >= 2 {
			if os.Args[1] == "deactivate" {
				main()
				os.Exit(0)
			}
		}
		fmt.Printf("You have %d pools marked active for iocage usage.\n"+
			"Run 'ioc deactivate ZPOOL' on %d of the pools.\n",
			match, match-1)
		os.Exit(1)
	}
}

func main() {
	Pool := &defloc.pool
	Iocroot := &defloc.iocroot

	// User needs to make sure they have a valid zpool before we can do anything
	// useful
	if *Pool == "" {
		if len(os.Args) <= 2 ||
			os.Args[1] != "activate" && os.Args[1] != "deactivate" {
			fmt.Println("No valid zpool available. " +
				"Please run 'ioc activate ZPOOL'.")
			os.Exit(1)
		}
	}

	// Check if all of our datasets are here.
	checkDatasets.Args(Pool, Iocroot)

	// This check goes here so we ensure that there is a valid pool and iocroot.
	exists := Exists(*Iocroot + "/.default")
	if !exists {
		defaultprops(Pool, Iocroot)
	}

	app := cli.NewApp()
	app.Name = "iocage"
	app.Usage = "a jail manager"
	app.Version = "0.1"
	app.Commands = []cli.Command{
		{
			Name:  "activate",
			Usage: "Mark a zpool for iocage usage",
			Action: func(c *cli.Context) {
				if c.NArg() > 0 {
					fmt.Printf(activate.Args(c.Args()[0]))
				} else {
					fmt.Println("Please supply a zpool!")
					os.Exit(1)
				}
			},
		},
		{
			Name:  "clean",
			Usage: "Clean a particular type's dataset.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "f, force",
					Usage: "No promts given and stops any jails or templates.",
				},
				cli.BoolFlag{
					Name:  "j, jails",
					Usage: "List all releases fetched.",
				},
				cli.BoolFlag{
					Name:  "r, releases",
					Usage: "List all releases fetched.",
				},
				cli.BoolFlag{
					Name:  "t, templates",
					Usage: "List all templates.",
				},
			},
			Action: func(c *cli.Context) {
				clean.Args(Pool, c.Bool("f"), c.Bool("j"),
					c.Bool("r"), c.Bool("t"))
			},
		},
		{
			Name:  "create",
			Usage: "Create a jail",
			Action: func(c *cli.Context) {
				create.Args(Pool, Iocroot, c.Args())
			},
		},
		{
			Name:  "deactivate",
			Usage: "Remove the mark denoting a zpool for iocage usage",
			Action: func(c *cli.Context) {
				if c.NArg() > 0 {
					fmt.Printf(deactivate.Args(c.Args()[0]))
				} else {
					fmt.Println("Please supply a zpool!")
					os.Exit(1)
				}
			},
		},
		{
			Name:  "destroy",
			Usage: "Destroy a jail.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "f, force",
					Usage: "Force a jail do be destroyed. Skips prompts" +
						" and stops the jail.",
				},
			},
			Action: func(c *cli.Context) {
				if c.NArg() > 0 {
					destroy.Args(Pool, c.Bool("f"), false, c.Args())
				} else {
					fmt.Println("Please supply a jail!")
					os.Exit(1)
				}
			},
		},
		{
			Name:  "fetch",
			Usage: "Fetch a release or plugin",
			Flags: []cli.Flag{
				// TODO: Make a real bool
				cli.BoolFlag{
					Name: "f, force",
					Usage: "Force a jail do be destroyed. Skips prompts" +
						" and stops the jail.",
				},
			},
			Action: func(c *cli.Context) {
				fetch.Args(Pool, Iocroot, c.Args())
			},
		},
		{
			Name:  "get",
			Usage: "Get a jails property",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "r, recursive",
					Usage: "Recursively get a property for all jails.",
				},
				cli.BoolFlag{
					Name:  "H, no-header",
					Usage: "Remove the header and table format.",
				},
			},
			Action: func(c *cli.Context) {
				if len(c.Args()) > 1 {
					fmt.Printf("%s\n", get.Args(c.Args()[0], c.Args()[1],
						Pool, Iocroot, c.Bool("r"), c.Bool("H")))
				} else {
					if c.Bool("r") {
						fmt.Printf(get.Args(c.Args()[0], "jail",
							Pool, Iocroot, c.Bool("r"), c.Bool("H")))
					} else {
						fmt.Println("Please supply a jail!")
						os.Exit(1)
					}
				}
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List all running jails or templates.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "H, no-header",
					Usage: "Remove the header and table format.",
				},
				cli.BoolFlag{
					Name:  "r, releases",
					Usage: "List all releases fetched.",
				},
				cli.BoolFlag{
					Name:  "t, templates",
					Usage: "List all templates.",
				},
			},
			Action: func(c *cli.Context) {
				if c.NArg() > 0 {
					list.Args(c.Args()[0], Pool, c.Bool("H"))
				} else {
					list.Args("all", Pool, c.Bool("H"), c.Bool("r"),
						c.Bool("t"))
				}
			},
		},
		{
			Name:  "set",
			Usage: "Set a jail property.",
			Action: func(c *cli.Context) {
				if c.NArg() > 0 {
					fmt.Println(set.Args(c.Args()[0], c.Args()[1], Pool,
						Iocroot))
				} else {
					fmt.Println("Please supply a jail!")
					os.Exit(1)
				}
			},
		},
	}
	app.Run(os.Args)
}

func defaultprops(Pool, Iocroot *string) {
	m := map[string]string{
		// Network properties
		"ipv6":            "off",
		"interfaces":      "vnet0:bridge0,vnet1:bridge1",
		"host_domainname": "none",
		"host_hostname":   "$uuid",
		"exec_fib":        "0",
		"ip4_addr":        "none",
		"ip4_autostart":   "none",
		"ip4_autoend":     "none",
		"ip4_autosubnet":  "none",
		"ip4_saddrsel":    "1",
		"ip4":             "new",
		"ip6_addr":        "none",
		"ip6_saddrsel":    "1",
		"ip6":             "new",
		"defaultrouter":   "none",
		"defaultrouter6":  "none",
		"resolver":        "none",
		"mac_prefix":      "02ff60",
		"vnet0_mac":       "none",
		"vnet1_mac":       "none",
		"vnet2_mac":       "none",
		"vnet3_mac":       "none",
		// Jail Properties
		"devfs_ruleset":         "4",
		"exec_start":            "/bin/sh /etc/rc",
		"exec_stop":             "/bin/sh /etc/rc.shutdown",
		"exec_prestart":         "/usr/bin/true",
		"exec_poststart":        "/usr/bin/true",
		"exec_prestop":          "/usr/bin/true",
		"exec_poststop":         "/usr/bin/true",
		"exec_clean":            "1",
		"exec_timeout":          "60",
		"stop_timeout":          "30",
		"exec_jail_user":        "root",
		"exec_system_jail_user": "0",
		"exec_system_user":      "root",
		"mount_devfs":           "1",
		"mount_fdescfs":         "1",
		"enforce_statfs":        "2",
		"children_max":          "0",
		"login_flags":           "-f root",
		"securelevel":           "2",
		"host_hostuuid":         "$uuid",
		"allow_set_hostname":    "1",
		"allow_sysvipc":         "0",
		"allow_raw_sockets":     "0",
		"allow_chflags":         "0",
		"allow_mount":           "0",
		"allow_mount_devfs":     "0",
		"allow_mount_nullfs":    "0",
		"allow_mount_procfs":    "0",
		"allow_mount_tmpfs":     "0",
		"allow_mount_zfs":       "0",
		"allow_quotas":          "0",
		"allow_socket_af":       "0",
		// RCTL limits
		"cpuset":          "off",
		"rlimits":         "off",
		"memoryuse":       "8G:log",
		"memorylocked":    "off",
		"vmemoryuse":      "off",
		"maxproc":         "off",
		"cputime":         "off",
		"pcpu":            "off",
		"datasize":        "off",
		"stacksize":       "off",
		"coredumpsize":    "off",
		"openfiles":       "off",
		"pseudoterminals": "off",
		"swapuse":         "off",
		"nthr":            "off",
		"msgqqueued":      "off",
		"msgqsize":        "off",
		"nmsgq":           "off",
		"nsemop":          "off",
		"nshm":            "off",
		"shmsize":         "off",
		"wallclock":       "off",
		// Custom properties
		"tag":                 "$(date \"+%F@%T\")",
		"istemplate":          "no",
		"bpf":                 "off",
		"dhcp":                "off",
		"boot":                "off",
		"notes":               "none",
		"owner":               "root",
		"priority":            "99",
		"last_started":        "none",
		"base":                "$(uname -r|cut -f 1,2 -d'-')",
		"template":            "none",
		"hostid":              "$(cat /etc/hostid)",
		"jail_zfs":            "off",
		"jail_zfs_dataset":    "iocage/jails/${uuid}/data",
		"jail_zfs_mountpoint": "none",
		"mount_procfs":        "0",
		"mount_linprocfs":     "0",
		"hack88":              "0",
		"count":               "1",
		// Sync properties
		"sync_state":     "none",
		"sync_target":    "none",
		"sync_tgt_zpool": "none",
		// FTP variables
		"ftphost":  "ftp.freebsd.org",
		"ftpdir":   "/pub/FreeBSD/releases/amd64",
		"ftpfiles": "base.txz doc.txz lib32.txz",
		// Git properties
		"gitlocation": "https://github.com"}
	f, err := os.Create(*Iocroot + "/.default")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer f.Close()

	for n, v := range m {
		str := n + "=" + v
		set.Args(str, "default", Pool, Iocroot)
	}
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
