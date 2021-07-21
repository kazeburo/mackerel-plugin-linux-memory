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
				{Name: "total", Label: "Total", Stacked: false},
				{Name: "used", Label: "Used", Stacked: false},
				{Name: "kernelstack", Label: "KernelStack", Stacked: true},
				{Name: "pagetables", Label: "PageTables", Stacked: true},
				{Name: "anonpages", Label: "AnonPages", Stacked: true},
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
		"total":       float64(*m.MemTotal * 1024),
		"used":        float64((*m.MemTotal - *m.MemAvailable) * 1024),
		"kernelstack": float64(*m.KernelStack * 1024),
		"pagetables":  float64(*m.PageTables * 1024),
		"anonpages":   float64(*m.AnonPages * 1024),
		"buffers":     float64(*m.Buffers * 1024),
		"cached":      float64(*m.Cached * 1024),
		"slab":        float64(*m.Slab * 1024),
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
