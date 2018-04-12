package version

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"text/template"

	"github.com/prometheus/client_golang/prometheus"
)

// Build information. Populated at build-time.
var (
	Version   string
	Revision  string
	Branch    string
	BuildDate string
	GoVersion = runtime.Version()
)

// NewCollector returns a collector which exports metrics about current version information.
func NewCollector(program string) *prometheus.GaugeVec {
	buildInfo := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: program,
			Name:      "build_info",
			Help: fmt.Sprintf(
				"A metric with a constant '1' value labeled by version, revision, branch, and goversion from which %s was built.",
				program,
			),
		},
		[]string{"version", "revision", "branch", "goversion"},
	)
	buildInfo.WithLabelValues(Version, Revision, Branch, GoVersion).Set(1)
	return buildInfo
}

// versionInfoTmpl contains the template used by Info.
var versionInfoTmpl = `
{{.program}}, version {{.version}} (branch: {{.branch}}, revision: {{.revision}})
  build date:       {{.buildDate}}
  go version:       {{.goVersion}}
`

// Print returns version information.
func Print(program string) string {
	m := map[string]string{
		"program":   program,
		"version":   Version,
		"revision":  Revision,
		"branch":    Branch,
		"buildDate": BuildDate,
		"goVersion": GoVersion,
	}
	t := template.Must(template.New("version").Parse(versionInfoTmpl))

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "version", m); err != nil {
		panic(err)
	}
	return strings.TrimSpace(buf.String())
}

// Info returns version, branch and revision information.
func Info() string {
	return fmt.Sprintf("(version=%s, branch=%s, revision=%s)", Version, Branch, Revision)
}

// BuildContext returns goVersion, buildUser and buildDate information.
func BuildContext() string {
	return fmt.Sprintf("(go=%s, date=%s)", GoVersion, BuildDate)
}
