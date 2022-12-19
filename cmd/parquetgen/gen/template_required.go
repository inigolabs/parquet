package gen

var requiredNumericTpl = `{{define "numericField"}}
type {{.FieldType}} struct {
	vals []{{.TypeName}}
	parquet.RequiredField
	read  func(r {{.StructType}}) {{.TypeName}}
	write func(r *{{.StructType}}, vals []{{removeStar .TypeName}})
	stats *{{.TypeName}}stats
}

func New{{.FieldType}}(read func(r {{.StructType}}) {{.TypeName}}, write func(r *{{.StructType}}, vals []{{removeStar .TypeName}}), path []string, opts ...func(*parquet.RequiredField)) *{{.FieldType}} {
	return &{{.FieldType}}{
		read:           read,
		write:          write,
		RequiredField: parquet.NewRequiredField(path, opts...),
		stats:         new{{camelCase .TypeName}}stats(),
	}
}

func (f *{{.FieldType}}) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Path: f.Path(), Type: {{.ParquetType}}, RepetitionType: parquet.RepetitionRequired, Types: []int{0}}
}

func (f *{{.FieldType}}) Read(ctx context.Context, r io.ReadSeeker, pg parquet.Page) error {
	rr, _, err := f.DoRead(ctx, r, pg)
	if err != nil {
		return err
	}

	v := make([]{{.TypeName}}, int(pg.N))
	err = binary.Read(rr, binary.LittleEndian, &v)
	f.vals = append(f.vals, v...)
	return err
}

func (f *{{.FieldType}}) Write(w io.Writer, meta *parquet.Metadata) error {
	buf := buffpool.Get()
	defer buffpool.Put(buf)

	bs := make([]byte, {{byteSize .}})
	for _, v := range f.vals {
		binary.LittleEndian.{{ putFunc . }}(bs, {{ uintFunc . }})
		if _, err := buf.Write(bs); err != nil {
			return err
		}
	}
	return f.DoWrite(w, meta, buf.Bytes(), len(f.vals), f.stats)
}

func (f *{{.FieldType}}) Scan(r *{{.StructType}}) {
	if len(f.vals) == 0 {
		return
	}

	f.write(r, f.vals)
	f.vals = f.vals[1:]
}

func (f *{{.FieldType}}) Add(r {{.Parent.StructType}}) {
	v := f.read(r)
	f.stats.add(v)
	f.vals = append(f.vals, v)
}

func (f *{{.FieldType}}) Levels() ([]uint8, []uint8) {
	return nil, nil
}
{{end}}`

var requiredStatsTpl = `{{define "requiredStats"}}
type {{.TypeName}}stats struct {
	min {{.TypeName}}
	max {{.TypeName}}
}

func new{{camelCase .TypeName}}stats() *{{.TypeName}}stats {
	return &{{.TypeName}}stats{
		min: {{.TypeName}}(math.Max{{camelCase .TypeName}}),
	}
}

func (i *{{.TypeName}}stats) add(val {{.TypeName}}) {
	if val < i.min {
		i.min = val
	}
	if val > i.max {
		i.max = val
	}
}

func (f *{{.TypeName}}stats) bytes(v {{.TypeName}}) []byte {
	bs := make([]byte, {{byteSize .}})
	binary.LittleEndian.{{ putFunc . }}(bs, {{ uintFunc . }})
	return bs
}

func (f *{{.TypeName}}stats) NullCount() *int64 {
	return nil
}

func (f *{{.TypeName}}stats) DistinctCount() *int64 {
	return nil
}

func (f *{{.TypeName}}stats) Min() []byte {
	return f.bytes(f.min)
}

func (f *{{.TypeName}}stats) Max() []byte {
	return f.bytes(f.max)
}
{{end}}`
