package crdgen

import (
	"gopkg.in/yaml.v3"
)

var (
	cleanNodeVisitorInstance = &cleanNodeVisitor{}
)

type cleanNodeVisitor struct {
}

func (d *cleanNodeVisitor) Visit(n *yaml.Node, ctx any) (any, NodeVisitor, error) {
	return ctx, d, nil
}

type cleanDocumentVisitor struct{}

func (c cleanDocumentVisitor) Visit(n map[string]any) error {
	return nil
}
