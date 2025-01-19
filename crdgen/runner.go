package crdgen

import (
	"flag"
	"fmt"
	"github.com/rancher/wrangler/v3/pkg/crd"
	"github.com/sirupsen/logrus"
)

var mode string
var chartName string

func init() {
	const (
		defaultMode      = "clean"
		defaultChartName = "default"
	)
	flag.StringVar(&mode, "mode", defaultMode, "mode for use")
	flag.StringVar(&chartName, "chart-name", defaultChartName, "chart name")
}

func Run(crds []crd.CRD) {

	flag.Parse()

	var p *printer

	switch mode {
	case "clean":
		p = &printer{
			nodeVisitor:     cleanNodeVisitorInstance,
			documentVisitor: &cleanDocumentVisitor{},
		}
	case "helm":
		p = &printer{
			nodeVisitor: helmNodeVisitor,
			documentVisitor: &helmDocumentVisitor{
				chartName: chartName,
			},
		}
	default:
		panic("invalid mode")
	}

	fmt.Println("# START CRD {{- .Values.crds.enabled }}")
	err := p.run(crds)
	fmt.Println("# END CRD {{- end }}")

	if err != nil {
		logrus.Fatal(err.Error())
	}

}
