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
	var err error

	switch mode {
	case "clean":
		p = &printer{
			nodeVisitor:     cleanNodeVisitorInstance,
			documentVisitor: &cleanDocumentVisitor{},
		}
		err = p.run(crds)
	case "helm":
		p = &printer{
			nodeVisitor: helmNodeVisitor,
			documentVisitor: &helmDocumentVisitor{
				chartName: chartName,
			},
		}
		fmt.Println("# START CRD {{- if .Values.crds.enabled }}")
		err = p.run(crds)
		fmt.Println("# END CRD {{- end }}")
	default:
		panic("invalid mode")
	}

	if err != nil {
		logrus.Fatal(err.Error())
	}

}
