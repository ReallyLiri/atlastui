package inspect

// adapted from ariga.io/atlas/cmdlog SchemaInspect.MarshalJSON()

type Attrs struct {
	Comment string `json:"comment,omitempty"`
	Charset string `json:"charset,omitempty"`
	Collate string `json:"collate,omitempty"`
}

type Column struct {
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
	Null bool   `json:"null,omitempty"`
	Attrs
}

type IndexPart struct {
	Desc   bool   `json:"desc,omitempty"`
	Column string `json:"column,omitempty"`
	Expr   string `json:"expr,omitempty"`
}

type Index struct {
	Name   string      `json:"name,omitempty"`
	Unique bool        `json:"unique,omitempty"`
	Parts  []IndexPart `json:"parts,omitempty"`
}

type ForeignKey struct {
	Name       string   `json:"name"`
	Columns    []string `json:"columns,omitempty"`
	References struct {
		Table   string   `json:"table"`
		Columns []string `json:"columns,omitempty"`
	} `json:"references"`
}

type Table struct {
	Name        string       `json:"name"`
	Columns     []Column     `json:"columns,omitempty"`
	Indexes     []Index      `json:"indexes,omitempty"`
	PrimaryKey  *Index       `json:"primary_key,omitempty"`
	ForeignKeys []ForeignKey `json:"foreign_keys,omitempty"`
	Attrs
}

type Schema struct {
	Name   string  `json:"name"`
	Tables []Table `json:"tables,omitempty"`
	Attrs
}

type Schemas struct {
	Schemas []Schema `json:"schemas"`
}
