package gen

var readTpl = `{{define "readFunc"}}
{{end}}`

var optionalNumericTpl = `{{define "optionalField"}}
type {{.FieldType}} struct {
	parquet.OptionalField
	vals  []{{removeStar .TypeName}}
	read   func(r {{.StructType}}, vals []{{removeStar .TypeName}}, defs, reps []uint8) ([]{{removeStar .TypeName}}, []uint8, []uint8)
	write  func(r *{{.StructType}}, vals []{{removeStar .TypeName}}, defs, reps []uint8) (int, int)
	stats *{{removeStar .TypeName}}optionalStats
}

func New{{.FieldType}}(read func(r {{.StructType}}, vals []{{removeStar .TypeName}}, defs, reps []uint8) ([]{{removeStar .TypeName}}, []uint8, []uint8), write func(r *{{.StructType}}, vals []{{removeStar .TypeName}}, defs, reps []uint8) (int, int), path []string, types []int, opts ...func(*parquet.OptionalField)) *{{.FieldType}} {
	return &{{.FieldType}}{
		read:          read,
		write:         write,
		OptionalField: parquet.NewOptionalField(path, types, opts...),
		stats:         new{{removeStar .TypeName}}optionalStats(maxDef(types)),
	}
}

func (f *{{.FieldType}}) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Path: f.Path(), Type: {{.ParquetType}}, RepetitionType: f.RepetitionType, Types: f.Types}
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
	return f.DoWrite(w, meta, buf.Bytes(), len(f.Defs), f.stats)
}

func (f *{{.FieldType}}) Read(ctx context.Context, r io.ReadSeeker, pg parquet.Page) error {
	rr, _, err := f.DoRead(ctx, r, pg)
	if err != nil {
		return err
	}

	v := make([]{{removeStar .TypeName}}, f.Values()-len(f.vals))
	err = binary.Read(rr, binary.LittleEndian, &v)
	f.vals = append(f.vals, v...)
	return err
}

func (f *{{.FieldType}}) Add(r {{.StructType}}) {
	vals, defs, reps := f.read(r, f.vals, f.Defs, f.Reps)
	f.stats.add(vals[len(f.vals):], defs[len(f.Defs):])
	f.vals = vals
	f.Defs = defs
	f.Reps = reps
}

func (f *{{.FieldType}}) Scan(r *{{.StructType}}) {
	if len(f.Defs) == 0 {
		return
	}

	v, l := f.write(r, f.vals, f.Defs, f.Reps)
	f.vals = f.vals[v:]
	f.Defs = f.Defs[l:]
	if len(f.Reps) > 0 {
		f.Reps = f.Reps[l:]
	}
}

func (f *{{.FieldType}}) Levels() ([]uint8, []uint8) {
	return f.Defs, f.Reps
}
{{end}}`

var optionalStatsTpl = `{{define "optionalStats"}}
type {{removeStar .TypeName}}optionalStats struct {
	min {{removeStar .TypeName}}
	max {{removeStar .TypeName}}
	nils int64
	nonNils int64
	maxDef uint8
}

func new{{removeStar .TypeName}}optionalStats(d uint8) *{{removeStar .TypeName}}optionalStats {
	return &{{removeStar .TypeName}}optionalStats{
		min: {{removeStar .TypeName}}(math.Max{{camelCaseRemoveStar .TypeName}}),
		maxDef: d,
	}
}

func (f *{{removeStar .TypeName}}optionalStats) add(vals []{{removeStar .TypeName}}, defs []uint8) {
	var i int
	for _, def := range defs {
		if def < f.maxDef {
			f.nils++
		} else {
			val := vals[i]
			i++

			f.nonNils++
			if val < f.min {
				f.min = val
			}
			if val > f.max {
				f.max = val
			}
		}
	}
}

func (f *{{removeStar .TypeName}}optionalStats) bytes(v {{removeStar .TypeName}}) []byte {
	bs := make([]byte, {{byteSize .}})
	binary.LittleEndian.{{ putFunc . }}(bs, {{ uintFunc . }})
	return bs
}

func (f *{{removeStar .TypeName}}optionalStats) NullCount() *int64 {
	return &f.nils
}

func (f *{{removeStar .TypeName}}optionalStats) DistinctCount() *int64 {
	return nil
}

func (f *{{removeStar .TypeName}}optionalStats) Min() []byte {
	if f.nonNils == 0  {
		return nil
	}
	return f.bytes(f.min)
}

func (f *{{removeStar .TypeName}}optionalStats) Max() []byte {
	if f.nonNils == 0  {
		return nil
	}
	return f.bytes(f.max)
}
{{end}}`
