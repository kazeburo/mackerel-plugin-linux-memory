package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/jessevdk/go-flags"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/prometheus/procfs"
)

// version by Makefile
var version string

type cmdOpts struct {
	Version bool `short:"v" long:"version" description:"Show version"`
}

type LinuxMemoryPlugin struct{}

func (u LinuxMemoryPlugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"": {
			Label: "Linux Memory",
			Unit:  mp.UnitBytes,
			Metrics: []mp.Metrics{
				{Name: "buffers", Label: "Buffers", Stacked: true},
				{Name: "cached", Label: "Cached", Stacked: true},
				{Name: "slab", Label: "Slab", Stacked: true},
			},
		},
	}
}

func (u LinuxMemoryPlugin) MetricKeyPrefix() string {
	return "linux-memory"
}

func (u LinuxMemoryPlugin) FetchMetrics() (map[string]float64, error) {
	fs, err := procfs.NewFS("/proc")
	if err != nil {
		return nil, err
	}
	m, err := fs.Meminfo()
	if err != nil {
		return nil, err
	}

	return map[string]float64{
		"buffers": float64(*m.Buffers),
		"cached":  float64(*m.Cached),
		"slab":    float64(*m.Slab),
	}, nil
}

func main() {
	os.Exit(_main())
}

func _main() int {
	opts := cmdOpts{}
	psr := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	_, err := psr.Parse()
	if opts.Version {
		fmt.Printf(`%s %s
Compiler: %s %s
`,
			os.Args[0],
			version,
			runtime.Compiler,
			runtime.Version())
		return 0
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	u := LinuxMemoryPlugin{}
	plugin := mp.NewMackerelPlugin(u)
	plugin.Run()
	return 0
}
