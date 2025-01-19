package crdgen

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/mapstructure"
	crd2 "github.com/rancher/wrangler/v3/pkg/crd"
	"gopkg.in/yaml.v3"
)

type printer struct {
	documentVisitor DocumentVisitor
	nodeVisitor     NodeVisitor
}

func walk(v NodeVisitor, n *yaml.Node, ctx any) (any, error) {

	ctx, next, err := v.Visit(n, ctx)

	if err != nil {
		return nil, err
	}

	if next == nil {
		next = v
	}

	var childCtx any

	for _, nn := range n.Content {
		childCtx, err = walk(next, nn, childCtx)
		if err != nil {
			return nil, err
		}
	}
	return ctx, nil

}

func (p *printer) run(crds []crd2.CRD) error {

	for _, c := range crds {

		o, err := c.ToCustomResourceDefinition()
		if err != nil {
			return err
		}
		bs, err := p.cleanMarshall(o)
		if err != nil {
			return err
		}

		fmt.Println(string(bs))

	}

	return nil
}

func (p *printer) cleanMarshall(in any) ([]byte, error) {

	out := map[string]any{}

	err := mapstructure.Decode(in, &out)

	if err != nil {
		return nil, err
	}

	out = out["Object"].(map[string]any)
	delete(out, "status")
	delete(out["metadata"].(map[string]any), "creationTimestamp")

	err = p.documentVisitor.Visit(out)

	if err != nil {
		return nil, err
	}

	var bs []byte
	bs, err = yaml.Marshal(out)

	if err != nil {
		return nil, err
	}

	n := &yaml.Node{}

	err = yaml.Unmarshal(bs, n)

	if err != nil {
		return nil, err
	}

	_, err = walk(p.nodeVisitor, n, nil)

	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	buf.WriteString("---\n")

	bs, err = yaml.Marshal(n)

	if err != nil {
		return nil, err
	}
	buf.Write(bs)
	return buf.Bytes(), nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
