package model

import (
	"encoding/json"
	"fmt"
	"github.com/akmalfairuz/finance/server/model/errortype"
	"net/http"
	"strings"
)

const (
	ProductDestinationFieldTextType   = "text"
	ProductDestinationFieldNumberType = "number"
	ProductDestinationFieldPhoneType  = "phone"
	ProductDestinationFieldOptionType = "option"
)

type ProductDestination struct {
	ID          int64  `db:"id"`
	Name        string `db:"name"`
	Fields_     string `db:"fields"`
	Format_     string `db:"format"`
	CheckerID   string `db:"checkerId"`
	Description string `db:"description"`

	fields map[string]DestinationField `db:"-"`
}

type DestinationField struct {
	Priority    int                      `json:"priority"`
	Label       string                   `json:"label"`
	Type        string                   `json:"type"`
	Description string                   `json:"description"`
	Options     []DestinationFieldOption `json:"options"`
	Required    bool                     `json:"required"`
}

type DestinationFieldOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func (d ProductDestination) Fields() map[string]DestinationField {
	if d.fields != nil {
		return d.fields
	}

	if d.Fields_ == "" {
		return nil
	}

	var fields map[string]DestinationField
	if err := json.Unmarshal([]byte(d.Fields_), &fields); err != nil {
		return nil
	}

	d.fields = fields
	return fields
}

func (d ProductDestination) Readable(input map[string]string) map[string]string {
	fields := d.Fields()
	out := map[string]string{}
	for fieldName, field := range fields {
		v, ok := input[fieldName]
		if !ok {
			continue
		}
		if field.Type == ProductDestinationFieldOptionType {
			for _, opt := range field.Options {
				if opt.Value == v {
					v = opt.Label
				}
			}
		}
		out[field.Label] = v
	}
	return out
}

func (d ProductDestination) Format(input map[string]string) (string, error) {
	fields := d.Fields()
	output := d.Format_
	for fieldName, field := range fields {
		v, ok := input[fieldName]
		if !ok || v == "" {
			if field.Required {
				return "", errortype.NewErrorVisible(errortype.CodeInvalidDestination, http.StatusUnauthorized, "Kolom "+field.Label+" wajib diisi")
			}
		}
		if field.Type == ProductDestinationFieldOptionType {
			optFound := false
			for _, opt := range field.Options {
				if opt.Value == v {
					optFound = true
				}
			}
			if !optFound {
				return "", fmt.Errorf("invalid option %s in %s field", v, fieldName)
			}
		}
		output = strings.ReplaceAll(output, "{{"+fieldName+"}}", v)
	}
	return output, nil
}

type CreateOrUpdateProductDestination struct {
	Name        string
	Description *string
	Format      *string
	CheckerID   *string
	Fields      map[string]DestinationField
}
