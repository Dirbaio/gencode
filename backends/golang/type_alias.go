package golang

import (
	"text/template"

	"github.com/andyleap/gencode/schema"
)

var (
	AliasTemps *template.Template
)

func init() {
	AliasTemps = template.New("AliasTemps")
	template.Must(AliasTemps.New("unmarshal").Parse(`
	{
        var aTemp {{.SubType}}
		{{.SubTypeCode}}
        {{.Target}} = {{.Alias}}(aTemp)
	}`))
}

type AliasTemp struct {
	Target      string
	SubType     string
	SubTypeCode string
	Alias       string
}

func (w *Walker) WalkAliasDef(at *schema.AliasType) (parts *StringBuilder, err error) {
	parts = &StringBuilder{}
	parts.Append(at.Alias)
	return
}

func (w *Walker) WalkAliasSize(at *schema.AliasType, target string) (parts *StringBuilder, err error) {
	return w.WalkTypeSize(at.SubType, target)
}

func (w *Walker) WalkAliasMarshal(at *schema.AliasType, target string) (parts *StringBuilder, err error) {
	return w.WalkTypeMarshal(at.SubType, target)
}

func (w *Walker) WalkAliasUnmarshal(at *schema.AliasType, target string) (parts *StringBuilder, err error) {
	parts = &StringBuilder{}

	subtypecode, err := w.WalkTypeUnmarshal(at.SubType, "aTemp")
	if err != nil {
		return nil, err
	}

	subtype, err := w.WalkTypeDef(at.SubType)
	if err != nil {
		return nil, err
	}

	err = parts.AddTemplate(AliasTemps, "unmarshal", AliasTemp{target, subtype.String(), subtypecode.String(), at.Alias})
	return
}
