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
				{Name: "available", Label: "Available", Stacked: false},
				{Name: "used", Label: "Used", Stacked: false},
				{Name: "kernelstack", Label: "KernelStack", Stacked: true},
				{Name: "vmallocused", Label: "VmallocUsed", Stacked: true},
				{Name: "pagetables", Label: "PageTables", Stacked: true},
				{Name: "mapped", Label: "Mapped", Stacked: true},
				{Name: "anonpages", Label: "AnonPages", Stacked: true},
				{Name: "slab", Label: "Slab", Stacked: true},
				{Name: "buffers", Label: "Buffers", Stacked: true},
				{Name: "cached", Label: "Cached", Stacked: true},
				{Name: "free", Label: "Free", Stacked: true},
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

	result := map[string]float64{
		"total":       float64(*m.MemTotal * 1024),
		"kernelstack": float64(*m.KernelStack * 1024),
		"vmallocused": float64(*m.VmallocUsed * 1024),
		"pagetables":  float64(*m.PageTables * 1024),
		"mapped":      float64(*m.Mapped * 1024),
		"anonpages":   float64(*m.AnonPages * 1024),
		"slab":        float64(*m.Slab * 1024),
		"buffers":     float64(*m.Buffers * 1024),
		"cached":      float64(*m.Cached * 1024),
		"free":        float64(*m.MemFree * 1024),
	}

	if m.MemAvailable != nil {
		result["used"] = float64((*m.MemTotal - *m.MemAvailable) * 1024)
		result["available"] = float64(*m.MemAvailable * 1024)
	} else {
		result["used"] = float64(*m.MemTotal - *m.MemFree - *m.Buffers - *m.Cached)
	}

	return result, nil
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
