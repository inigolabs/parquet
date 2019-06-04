package parquet_test

// This code is generated by github.com/parsyl/parquet.

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/parsyl/parquet"

	"math"
	"sort"
)

type compression int

const (
	compressionUncompressed compression = 0
	compressionSnappy       compression = 1
	compressionUnknown      compression = -1
)

// ParquetWriter reprents a row group
type ParquetWriter struct {
	fields []Field

	len int

	// child points to the next page
	child *ParquetWriter

	// max is the number of Record items that can get written before
	// a new set of column chunks is written
	max int

	meta        *parquet.Metadata
	w           io.Writer
	compression compression
}

func readHobbyName(x Person) (*string, int64) {
	var def int64
	if x.Hobby == nil {
		return nil, 0
	}
	return &x.Hobby.Name, 2
}

func writeHobbyName(x *Person, v *string, def int64) {
	switch def {
	case 1:
		if x.Hobby == nil {
			x.Hobby = &Hobby{Name: *v}
		} else {
			x.Hobby.Name = *v
		}
	}
}

func readHobbyDifficulty(x Person) (*int32, int64) {
	switch {
	case x.Hobby == nil:
		return nil, 0
	case x.Hobby.Difficulty == nil:
		return nil, 1
	default:
		return x.Hobby.Difficulty, 2
	}
}

func writeHobbyDifficulty(x *Person, v *int32, def int64) {
	switch def {
	case 1:
		if x.Hobby == nil {
			x.Hobby = &Hobby{}
		}
	case 2:
		if x.Hobby == nil {
			x.Hobby = &Hobby{Difficulty: v}
		} else {
			x.Hobby.Difficulty = v
		}
	}
}

func readCode(x Person) (*string, int64) {
	switch {
	case x.Code == nil:
		return nil, 0
	default:
		return x.Code, 1
	}
}

func writeCode(x *Person, v *string, def int64) {
	x.Code = v
}

func Fields(compression compression) []Field {
	return []Field{
		NewInt32Field(func(x Person) int32 { return x.ID }, func(x *Person, v int32) { x.ID = v }, "id", fieldCompression(compression)...),
		NewInt32OptionalField(func(x Person) (*int32, int64) {
			var def int64
			if x.Age != nil {
				def = 1
			}
			return x.Age, def
		}, func(x *Person, v *int32, def int64) { x.Age = v }, "age", optionalFieldCompression(compression)...),
		NewInt64Field(func(x Person) int64 { return x.Happiness }, func(x *Person, v int64) { x.Happiness = v }, "happiness", fieldCompression(compression)...),
		NewInt64OptionalField(func(x Person) *int64 { return x.Sadness }, func(x *Person, v *int64) { x.Sadness = v }, "sadness", optionalFieldCompression(compression)...),
		NewStringOptionalField(readCode, writeCode, "code", optionalFieldCompression(compression)...),
		NewFloat32Field(func(x Person) float32 { return x.Funkiness }, func(x *Person, v float32) { x.Funkiness = v }, "funkiness", fieldCompression(compression)...),
		NewFloat64Field(func(x Person) float64 { return x.Boldness }, func(x *Person, v float64) { x.Boldness = v }, "boldness", fieldCompression(compression)...),
		NewFloat32OptionalField(func(x Person) *float32 { return x.Lameness }, func(x *Person, v *float32) { x.Lameness = v }, "lameness", optionalFieldCompression(compression)...),
		NewBoolOptionalField(func(x Person) *bool { return x.Keen }, func(x *Person, v *bool) { x.Keen = v }, "keen", optionalFieldCompression(compression)...),
		NewUint32Field(func(x Person) uint32 { return x.Birthday }, func(x *Person, v uint32) { x.Birthday = v }, "birthday", fieldCompression(compression)...),
		NewUint64OptionalField(func(x Person) *uint64 { return x.Anniversary }, func(x *Person, v *uint64) { x.Anniversary = v }, "anniversary", optionalFieldCompression(compression)...),
		NewStringField(func(x Person) string { return x.BFF }, func(x *Person, v string) { x.BFF = v }, "bff", fieldCompression(compression)...),
		NewBoolField(func(x Person) bool { return x.Hungry }, func(x *Person, v bool) { x.Hungry = v }, "hungry", fieldCompression(compression)...),
		NewBoolField(func(x Person) bool { return x.Sleepy }, func(x *Person, v bool) { x.Sleepy = v }, "Sleepy", fieldCompression(compression)...),
		NewStringOptionalField(readHobbyName, writeHobbyName, "hobby.name", optionalFieldCompression(compression)...),
	}
}

func fieldCompression(c compression) []func(*parquet.RequiredField) {
	switch c {
	case compressionUncompressed:
		return []func(*parquet.RequiredField){parquet.RequiredFieldUncompressed}
	case compressionSnappy:
		return []func(*parquet.RequiredField){parquet.RequiredFieldSnappy}
	default:
		return []func(*parquet.RequiredField){}
	}
}

func optionalFieldCompression(c compression) []func(*parquet.OptionalField) {
	switch c {
	case compressionUncompressed:
		return []func(*parquet.OptionalField){parquet.OptionalFieldUncompressed}
	case compressionSnappy:
		return []func(*parquet.OptionalField){parquet.OptionalFieldSnappy}
	default:
		return []func(*parquet.OptionalField){}
	}
}

func NewParquetWriter(w io.Writer, opts ...func(*ParquetWriter) error) (*ParquetWriter, error) {
	return newParquetWriter(w, append(opts, begin)...)
}

func newParquetWriter(w io.Writer, opts ...func(*ParquetWriter) error) (*ParquetWriter, error) {
	p := &ParquetWriter{
		max:         1000,
		w:           w,
		compression: compressionSnappy,
	}

	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, err
		}
	}

	p.fields = Fields(p.compression)
	if p.meta == nil {
		ff := Fields(p.compression)
		schema := make([]parquet.Field, len(ff))
		for i, f := range ff {
			schema[i] = f.Schema()
		}
		p.meta = parquet.New(schema...)
	}

	return p, nil
}

// MaxPageSize is the maximum number of rows in each row groups' page.
func MaxPageSize(m int) func(*ParquetWriter) error {
	return func(p *ParquetWriter) error {
		p.max = m
		return nil
	}
}

func begin(p *ParquetWriter) error {
	_, err := p.w.Write([]byte("PAR1"))
	return err
}

func withMeta(m *parquet.Metadata) func(*ParquetWriter) error {
	return func(p *ParquetWriter) error {
		p.meta = m
		return nil
	}
}

func Uncompressed(p *ParquetWriter) error {
	p.compression = compressionUncompressed
	return nil
}

func Snappy(p *ParquetWriter) error {
	p.compression = compressionSnappy
	return nil
}

func withCompression(c compression) func(*ParquetWriter) error {
	return func(p *ParquetWriter) error {
		p.compression = c
		return nil
	}
}

func (p *ParquetWriter) Write() error {
	for i, f := range p.fields {
		if err := f.Write(p.w, p.meta); err != nil {
			return err
		}

		for child := p.child; child != nil; child = child.child {
			if err := child.fields[i].Write(p.w, p.meta); err != nil {
				return err
			}
		}
	}

	p.fields = Fields(p.compression)
	p.child = nil
	p.len = 0

	schema := make([]parquet.Field, len(p.fields))
	for i, f := range p.fields {
		schema[i] = f.Schema()
	}
	p.meta.StartRowGroup(schema...)
	return nil
}

func (p *ParquetWriter) Close() error {
	if err := p.meta.Footer(p.w); err != nil {
		return err
	}

	_, err := p.w.Write([]byte("PAR1"))
	return err
}

func (p *ParquetWriter) Add(rec Person) {
	if p.len == p.max {
		if p.child == nil {
			// an error can't happen here
			p.child, _ = newParquetWriter(p.w, MaxPageSize(p.max), withMeta(p.meta), withCompression(p.compression))
		}

		p.child.Add(rec)
		return
	}

	for _, f := range p.fields {
		f.Add(rec)
	}

	p.len++
}

type Field interface {
	Add(r Person)
	Write(w io.Writer, meta *parquet.Metadata) error
	Schema() parquet.Field
	Scan(r *Person)
	Read(r io.ReadSeeker, pg parquet.Page) error
	Name() string
}

func getFields(ff []Field) map[string]Field {
	m := make(map[string]Field, len(ff))
	for _, f := range ff {
		m[f.Name()] = f
	}
	return m
}

func NewParquetReader(r io.ReadSeeker, opts ...func(*ParquetReader)) (*ParquetReader, error) {
	ff := Fields(compressionUnknown)
	pr := &ParquetReader{
		r: r,
	}

	for _, opt := range opts {
		opt(pr)
	}

	schema := make([]parquet.Field, len(ff))
	for i, f := range ff {
		schema[i] = f.Schema()
	}

	meta := parquet.New(schema...)
	if err := meta.ReadFooter(r); err != nil {
		return nil, err
	}
	pr.rows = meta.Rows()
	var err error
	pr.pages, err = meta.Pages()
	if err != nil {
		return nil, err
	}

	pr.rowGroups = meta.RowGroups()
	_, err = r.Seek(4, io.SeekStart)
	if err != nil {
		return nil, err
	}
	pr.meta = meta

	return pr, pr.readRowGroup()
}

func readerIndex(i int) func(*ParquetReader) {
	return func(p *ParquetReader) {
		p.index = i
	}
}

// ParquetReader reads one page from a row group.
type ParquetReader struct {
	fields         map[string]Field
	index          int
	cursor         int64
	rows           int64
	rowGroupCursor int64
	rowGroupCount  int64
	pages          map[string][]parquet.Page
	meta           *parquet.Metadata
	err            error

	r         io.ReadSeeker
	rowGroups []parquet.RowGroup
}

func (p *ParquetReader) Error() error {
	return p.err
}

func (p *ParquetReader) readRowGroup() error {
	p.rowGroupCursor = 0

	if len(p.rowGroups) == 0 {
		p.rowGroupCount = 0
		return nil
	}

	rg := p.rowGroups[0]
	p.fields = getFields(Fields(compressionUnknown))
	p.rowGroupCount = rg.Rows
	p.rowGroupCursor = 0
	for _, col := range rg.Columns() {
		name := col.MetaData.PathInSchema[len(col.MetaData.PathInSchema)-1]
		f, ok := p.fields[name]
		if !ok {
			return fmt.Errorf("unknown field: %s", name)
		}
		pages := p.pages[f.Name()]
		if len(pages) <= p.index {
			break
		}

		pg := pages[0]
		if err := f.Read(p.r, pg); err != nil {
			return fmt.Errorf("unable to read field %s, err: %s", f.Name(), err)
		}
		p.pages[f.Name()] = p.pages[f.Name()][1:]
	}
	p.rowGroups = p.rowGroups[1:]
	return nil
}

func (p *ParquetReader) Rows() int64 {
	return p.rows
}

func (p *ParquetReader) Next() bool {
	if p.err == nil && p.cursor >= p.rows {
		return false
	}
	if p.rowGroupCursor >= p.rowGroupCount {
		p.err = p.readRowGroup()
		if p.err != nil {
			return false
		}
	}

	p.cursor++
	p.rowGroupCursor++
	return true
}

func (p *ParquetReader) Scan(x *Person) {
	if p.err != nil {
		return
	}

	for _, f := range p.fields {
		f.Scan(x)
	}
}

type Int32Field struct {
	vals []int32
	parquet.RequiredField
	read  func(r Person) (int32, int64)
	write func(r *Person, vals []int32, def int64)
	stats *int32stats
}

func NewInt32Field(read func(r Person) (int32, int64), write func(r *Person, vals []int32, def int64), col string, opts ...func(*parquet.RequiredField)) *Int32Field {
	return &Int32Field{
		read:          read,
		write:         write,
		RequiredField: parquet.NewRequiredField(col, opts...),
		stats:         newInt32stats(),
	}
}

func (f *Int32Field) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Type: parquet.Int32Type, RepetitionType: parquet.RepetitionRequired}
}

func (f *Int32Field) Scan(r *Person) {
	if len(f.vals) == 0 {
		return
	}

	f.write(r, f.vals, 0)
	f.vals = f.vals[1:]
}

func (f *Int32Field) Write(w io.Writer, meta *parquet.Metadata) error {
	var buf bytes.Buffer
	for _, v := range f.vals {
		if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
			return err
		}
	}
	return f.DoWrite(w, meta, buf.Bytes(), len(f.vals), f.stats)
}

func (f *Int32Field) Read(r io.ReadSeeker, pg parquet.Page) error {
	rr, _, err := f.DoRead(r, pg)
	if err != nil {
		return err
	}

	v := make([]int32, int(pg.N))
	err = binary.Read(rr, binary.LittleEndian, &v)
	f.vals = append(f.vals, v...)
	return err
}

func (f *Int32Field) Add(r Person) {
	v, f.defs = f.read(r, f.defs)
	f.stats.add(v)
	f.vals = append(f.vals, v)
}

type Int32OptionalField struct {
	parquet.OptionalField
	vals  []int32
	read  func(r *Person, v *int32, def int64)
	val   func(r Person) (*int32, int64)
	stats *int32optionalStats
}

func NewInt32OptionalField(val func(r Person) (*int32, int64), read func(r *Person, v *int32, def int64), col string, opts ...func(*parquet.OptionalField)) *Int32OptionalField {
	return &Int32OptionalField{
		val:           val,
		read:          read,
		OptionalField: parquet.NewOptionalField(col, opts...),
		stats:         newint32optionalStats(),
	}
}

func (f *Int32OptionalField) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Type: parquet.Int32Type, RepetitionType: parquet.RepetitionOptional}
}

func (f *Int32OptionalField) Write(w io.Writer, meta *parquet.Metadata) error {
	var buf bytes.Buffer
	for _, v := range f.vals {
		if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
			return err
		}
	}
	return f.DoWrite(w, meta, buf.Bytes(), len(f.vals), f.stats)
}

func (f *Int32OptionalField) Read(r io.ReadSeeker, pg parquet.Page) error {
	rr, _, err := f.DoRead(r, pg)
	if err != nil {
		return err
	}

	v := make([]int32, f.Values()-len(f.vals))
	err = binary.Read(rr, binary.LittleEndian, &v)
	f.vals = append(f.vals, v...)
	return err
}

func (f *Int32OptionalField) Add(r Person) {
	v := f.val(r)
	f.stats.add(v)
	if v != nil {
		f.vals = append(f.vals, *v)
		f.Defs = append(f.Defs, 1)
	} else {
		f.Defs = append(f.Defs, 0)
	}
}

func (f *Int32OptionalField) Scan(r *Person) {
	if len(f.Defs) == 0 {
		return
	}

	if f.Defs[0] == 1 {
		var val int32
		v := f.vals[0]
		f.vals = f.vals[1:]
		val = v
		f.read(r, &val)
	}
	f.Defs = f.Defs[1:]
}

type Int64Field struct {
	vals []int64
	parquet.RequiredField
	val   func(r Person) int64
	read  func(r *Person, v int64)
	stats *int64stats
}

func NewInt64Field(val func(r Person) int64, read func(r *Person, v int64), col string, opts ...func(*parquet.RequiredField)) *Int64Field {
	return &Int64Field{
		val:           val,
		read:          read,
		RequiredField: parquet.NewRequiredField(col, opts...),
		stats:         newInt64stats(),
	}
}

func (f *Int64Field) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Type: parquet.Int64Type, RepetitionType: parquet.RepetitionRequired}
}

func (f *Int64Field) Scan(r *Person) {
	if len(f.vals) == 0 {
		return
	}
	v := f.vals[0]
	f.vals = f.vals[1:]
	f.read(r, v)
}

func (f *Int64Field) Write(w io.Writer, meta *parquet.Metadata) error {
	var buf bytes.Buffer
	for _, v := range f.vals {
		if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
			return err
		}
	}
	return f.DoWrite(w, meta, buf.Bytes(), len(f.vals), f.stats)
}

func (f *Int64Field) Read(r io.ReadSeeker, pg parquet.Page) error {
	rr, _, err := f.DoRead(r, pg)
	if err != nil {
		return err
	}

	v := make([]int64, int(pg.N))
	err = binary.Read(rr, binary.LittleEndian, &v)
	f.vals = append(f.vals, v...)
	return err
}

func (f *Int64Field) Add(r Person) {
	v := f.val(r)
	f.stats.add(v)
	f.vals = append(f.vals, v)
}

type Int64OptionalField struct {
	parquet.OptionalField
	vals  []int64
	read  func(r *Person, v *int64)
	val   func(r Person) *int64
	stats *int64optionalStats
}

func NewInt64OptionalField(val func(r Person) *int64, read func(r *Person, v *int64), col string, opts ...func(*parquet.OptionalField)) *Int64OptionalField {
	return &Int64OptionalField{
		val:           val,
		read:          read,
		OptionalField: parquet.NewOptionalField(col, opts...),
		stats:         newint64optionalStats(),
	}
}

func (f *Int64OptionalField) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Type: parquet.Int64Type, RepetitionType: parquet.RepetitionOptional}
}

func (f *Int64OptionalField) Write(w io.Writer, meta *parquet.Metadata) error {
	var buf bytes.Buffer
	for _, v := range f.vals {
		if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
			return err
		}
	}
	return f.DoWrite(w, meta, buf.Bytes(), len(f.vals), f.stats)
}

func (f *Int64OptionalField) Read(r io.ReadSeeker, pg parquet.Page) error {
	rr, _, err := f.DoRead(r, pg)
	if err != nil {
		return err
	}

	v := make([]int64, f.Values()-len(f.vals))
	err = binary.Read(rr, binary.LittleEndian, &v)
	f.vals = append(f.vals, v...)
	return err
}

func (f *Int64OptionalField) Add(r Person) {
	v := f.val(r)
	f.stats.add(v)
	if v != nil {
		f.vals = append(f.vals, *v)
		f.Defs = append(f.Defs, 1)
	} else {
		f.Defs = append(f.Defs, 0)
	}
}

func (f *Int64OptionalField) Scan(r *Person) {
	if len(f.Defs) == 0 {
		return
	}

	if f.Defs[0] == 1 {
		var val int64
		v := f.vals[0]
		f.vals = f.vals[1:]
		val = v
		f.read(r, &val)
	}
	f.Defs = f.Defs[1:]
}

type StringOptionalField struct {
	parquet.OptionalField
	vals  []string
	read  func(r Person) (*string, int64)
	write func(r *Person, v *string, def int64)
	stats *stringOptionalStats
}

func NewStringOptionalField(read func(r Person) (*string, int64), write func(r *Person, v *string, def int64), col string, opts ...func(*parquet.OptionalField)) *StringOptionalField {
	return &StringOptionalField{
		read:          read,
		write:         write,
		OptionalField: parquet.NewOptionalField(col, opts...),
		stats:         newStringOptionalStats(),
	}
}

func (f *StringOptionalField) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Type: parquet.StringType, RepetitionType: parquet.RepetitionOptional}
}

func (f *StringOptionalField) Scan(r *Person) {
	if len(f.Defs) == 0 {
		return
	}

	if f.Defs[0] == int64(f.OptionalField.Depth) {
		var val *string
		v := f.vals[0]
		f.vals = f.vals[1:]
		val = &v
		f.write(r, val, f.Defs[0])
	} else {
		f.write(r, nil, f.Defs[0])
	}
	f.Defs = f.Defs[1:]
}

func (f *StringOptionalField) Add(r Person) {
	v, depth := f.read(r)
	f.stats.add(v)
	if v != nil {
		f.vals = append(f.vals, *v)
		f.Defs = append(f.Defs, depth)
	} else {
		f.Defs = append(f.Defs, depth)
	}
}

func (f *StringOptionalField) Write(w io.Writer, meta *parquet.Metadata) error {
	buf := bytes.Buffer{}

	for _, s := range f.vals {
		if err := binary.Write(&buf, binary.LittleEndian, int32(len(s))); err != nil {
			return err
		}
		buf.Write([]byte(s))
	}

	return f.DoWrite(w, meta, buf.Bytes(), len(f.vals), f.stats)
}

func (f *StringOptionalField) Read(r io.ReadSeeker, pg parquet.Page) error {
	start := len(f.Defs)
	rr, _, err := f.DoRead(r, pg)
	if err != nil {
		return err
	}

	for j := 0; j < pg.N; j++ {
		if f.Defs[start+j] != int64(f.Depth) {
			continue
		}

		var x int32
		if err := binary.Read(rr, binary.LittleEndian, &x); err != nil {
			return err
		}
		s := make([]byte, x)
		if _, err := rr.Read(s); err != nil {
			return err
		}

		f.vals = append(f.vals, string(s))
	}
	return nil
}

type Float32Field struct {
	vals []float32
	parquet.RequiredField
	val   func(r Person) float32
	read  func(r *Person, v float32)
	stats *float32stats
}

func NewFloat32Field(val func(r Person) float32, read func(r *Person, v float32), col string, opts ...func(*parquet.RequiredField)) *Float32Field {
	return &Float32Field{
		val:           val,
		read:          read,
		RequiredField: parquet.NewRequiredField(col, opts...),
		stats:         newFloat32stats(),
	}
}

func (f *Float32Field) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Type: parquet.Float32Type, RepetitionType: parquet.RepetitionRequired}
}

func (f *Float32Field) Scan(r *Person) {
	if len(f.vals) == 0 {
		return
	}
	v := f.vals[0]
	f.vals = f.vals[1:]
	f.read(r, v)
}

func (f *Float32Field) Write(w io.Writer, meta *parquet.Metadata) error {
	var buf bytes.Buffer
	for _, v := range f.vals {
		if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
			return err
		}
	}
	return f.DoWrite(w, meta, buf.Bytes(), len(f.vals), f.stats)
}

func (f *Float32Field) Read(r io.ReadSeeker, pg parquet.Page) error {
	rr, _, err := f.DoRead(r, pg)
	if err != nil {
		return err
	}

	v := make([]float32, int(pg.N))
	err = binary.Read(rr, binary.LittleEndian, &v)
	f.vals = append(f.vals, v...)
	return err
}

func (f *Float32Field) Add(r Person) {
	v := f.val(r)
	f.stats.add(v)
	f.vals = append(f.vals, v)
}

type Float64Field struct {
	vals []float64
	parquet.RequiredField
	val   func(r Person) float64
	read  func(r *Person, v float64)
	stats *float64stats
}

func NewFloat64Field(val func(r Person) float64, read func(r *Person, v float64), col string, opts ...func(*parquet.RequiredField)) *Float64Field {
	return &Float64Field{
		val:           val,
		read:          read,
		RequiredField: parquet.NewRequiredField(col, opts...),
		stats:         newFloat64stats(),
	}
}

func (f *Float64Field) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Type: parquet.Float64Type, RepetitionType: parquet.RepetitionRequired}
}

func (f *Float64Field) Scan(r *Person) {
	if len(f.vals) == 0 {
		return
	}
	v := f.vals[0]
	f.vals = f.vals[1:]
	f.read(r, v)
}

func (f *Float64Field) Write(w io.Writer, meta *parquet.Metadata) error {
	var buf bytes.Buffer
	for _, v := range f.vals {
		if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
			return err
		}
	}
	return f.DoWrite(w, meta, buf.Bytes(), len(f.vals), f.stats)
}

func (f *Float64Field) Read(r io.ReadSeeker, pg parquet.Page) error {
	rr, _, err := f.DoRead(r, pg)
	if err != nil {
		return err
	}

	v := make([]float64, int(pg.N))
	err = binary.Read(rr, binary.LittleEndian, &v)
	f.vals = append(f.vals, v...)
	return err
}

func (f *Float64Field) Add(r Person) {
	v := f.val(r)
	f.stats.add(v)
	f.vals = append(f.vals, v)
}

type Float32OptionalField struct {
	parquet.OptionalField
	vals  []float32
	read  func(r *Person, v *float32)
	val   func(r Person) *float32
	stats *float32optionalStats
}

func NewFloat32OptionalField(val func(r Person) *float32, read func(r *Person, v *float32), col string, opts ...func(*parquet.OptionalField)) *Float32OptionalField {
	return &Float32OptionalField{
		val:           val,
		read:          read,
		OptionalField: parquet.NewOptionalField(col, opts...),
		stats:         newfloat32optionalStats(),
	}
}

func (f *Float32OptionalField) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Type: parquet.Float32Type, RepetitionType: parquet.RepetitionOptional}
}

func (f *Float32OptionalField) Write(w io.Writer, meta *parquet.Metadata) error {
	var buf bytes.Buffer
	for _, v := range f.vals {
		if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
			return err
		}
	}
	return f.DoWrite(w, meta, buf.Bytes(), len(f.vals), f.stats)
}

func (f *Float32OptionalField) Read(r io.ReadSeeker, pg parquet.Page) error {
	rr, _, err := f.DoRead(r, pg)
	if err != nil {
		return err
	}

	v := make([]float32, f.Values()-len(f.vals))
	err = binary.Read(rr, binary.LittleEndian, &v)
	f.vals = append(f.vals, v...)
	return err
}

func (f *Float32OptionalField) Add(r Person) {
	v := f.val(r)
	f.stats.add(v)
	if v != nil {
		f.vals = append(f.vals, *v)
		f.Defs = append(f.Defs, 1)
	} else {
		f.Defs = append(f.Defs, 0)
	}
}

func (f *Float32OptionalField) Scan(r *Person) {
	if len(f.Defs) == 0 {
		return
	}

	if f.Defs[0] == 1 {
		var val float32
		v := f.vals[0]
		f.vals = f.vals[1:]
		val = v
		f.read(r, &val)
	}
	f.Defs = f.Defs[1:]
}

type BoolOptionalField struct {
	parquet.OptionalField
	vals  []bool
	val   func(r Person) *bool
	read  func(r *Person, v *bool)
	stats *boolOptionalStats
}

func NewBoolOptionalField(val func(r Person) *bool, read func(r *Person, v *bool), col string, opts ...func(*parquet.OptionalField)) *BoolOptionalField {
	return &BoolOptionalField{
		val:           val,
		read:          read,
		OptionalField: parquet.NewOptionalField(col, opts...),
		stats:         newBoolOptionalStats(),
	}
}

func (f *BoolOptionalField) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Type: parquet.BoolType, RepetitionType: parquet.RepetitionOptional}
}

func (f *BoolOptionalField) Read(r io.ReadSeeker, pg parquet.Page) error {
	rr, sizes, err := f.DoRead(r, pg)
	if err != nil {
		return err
	}

	v, err := parquet.GetBools(rr, f.Values()-len(f.vals), sizes)
	f.vals = append(f.vals, v...)
	return err
}

func (f *BoolOptionalField) Scan(r *Person) {
	if len(f.Defs) == 0 {
		return
	}

	var val *bool
	if f.Defs[0] == 1 {
		v := f.vals[0]
		f.vals = f.vals[1:]
		val = &v
		f.read(r, val)
	}
	f.Defs = f.Defs[1:]
}

func (f *BoolOptionalField) Add(r Person) {
	v := f.val(r)
	f.stats.add(v)
	if v != nil {
		f.vals = append(f.vals, *v)
		f.Defs = append(f.Defs, 1)
	} else {
		f.Defs = append(f.Defs, 0)
	}
}

func (f *BoolOptionalField) Write(w io.Writer, meta *parquet.Metadata) error {
	ln := len(f.vals)
	byteNum := (ln + 7) / 8
	rawBuf := make([]byte, byteNum)

	for i := 0; i < ln; i++ {
		if f.vals[i] {
			rawBuf[i/8] = rawBuf[i/8] | (1 << uint32(i%8))
		}
	}

	return f.DoWrite(w, meta, rawBuf, len(f.vals), f.stats)
}

type Uint32Field struct {
	vals []uint32
	parquet.RequiredField
	val   func(r Person) uint32
	read  func(r *Person, v uint32)
	stats *uint32stats
}

func NewUint32Field(val func(r Person) uint32, read func(r *Person, v uint32), col string, opts ...func(*parquet.RequiredField)) *Uint32Field {
	return &Uint32Field{
		val:           val,
		read:          read,
		RequiredField: parquet.NewRequiredField(col, opts...),
		stats:         newUint32stats(),
	}
}

func (f *Uint32Field) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Type: parquet.Uint32Type, RepetitionType: parquet.RepetitionRequired}
}

func (f *Uint32Field) Scan(r *Person) {
	if len(f.vals) == 0 {
		return
	}
	v := f.vals[0]
	f.vals = f.vals[1:]
	f.read(r, v)
}

func (f *Uint32Field) Write(w io.Writer, meta *parquet.Metadata) error {
	var buf bytes.Buffer
	for _, v := range f.vals {
		if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
			return err
		}
	}
	return f.DoWrite(w, meta, buf.Bytes(), len(f.vals), f.stats)
}

func (f *Uint32Field) Read(r io.ReadSeeker, pg parquet.Page) error {
	rr, _, err := f.DoRead(r, pg)
	if err != nil {
		return err
	}

	v := make([]uint32, int(pg.N))
	err = binary.Read(rr, binary.LittleEndian, &v)
	f.vals = append(f.vals, v...)
	return err
}

func (f *Uint32Field) Add(r Person) {
	v := f.val(r)
	f.stats.add(v)
	f.vals = append(f.vals, v)
}

type Uint64OptionalField struct {
	parquet.OptionalField
	vals  []uint64
	read  func(r *Person, v *uint64)
	val   func(r Person) *uint64
	stats *uint64optionalStats
}

func NewUint64OptionalField(val func(r Person) *uint64, read func(r *Person, v *uint64), col string, opts ...func(*parquet.OptionalField)) *Uint64OptionalField {
	return &Uint64OptionalField{
		val:           val,
		read:          read,
		OptionalField: parquet.NewOptionalField(col, opts...),
		stats:         newuint64optionalStats(),
	}
}

func (f *Uint64OptionalField) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Type: parquet.Uint64Type, RepetitionType: parquet.RepetitionOptional}
}

func (f *Uint64OptionalField) Write(w io.Writer, meta *parquet.Metadata) error {
	var buf bytes.Buffer
	for _, v := range f.vals {
		if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
			return err
		}
	}
	return f.DoWrite(w, meta, buf.Bytes(), len(f.vals), f.stats)
}

func (f *Uint64OptionalField) Read(r io.ReadSeeker, pg parquet.Page) error {
	rr, _, err := f.DoRead(r, pg)
	if err != nil {
		return err
	}

	v := make([]uint64, f.Values()-len(f.vals))
	err = binary.Read(rr, binary.LittleEndian, &v)
	f.vals = append(f.vals, v...)
	return err
}

func (f *Uint64OptionalField) Add(r Person) {
	v := f.val(r)
	f.stats.add(v)
	if v != nil {
		f.vals = append(f.vals, *v)
		f.Defs = append(f.Defs, 1)
	} else {
		f.Defs = append(f.Defs, 0)
	}
}

func (f *Uint64OptionalField) Scan(r *Person) {
	if len(f.Defs) == 0 {
		return
	}

	if f.Defs[0] == 1 {
		var val uint64
		v := f.vals[0]
		f.vals = f.vals[1:]
		val = v
		f.read(r, &val)
	}
	f.Defs = f.Defs[1:]
}

type StringField struct {
	parquet.RequiredField
	vals  []string
	val   func(r Person) string
	read  func(r *Person, v string)
	stats *stringStats
}

func NewStringField(val func(r Person) string, read func(r *Person, v string), col string, opts ...func(*parquet.RequiredField)) *StringField {
	return &StringField{
		val:           val,
		read:          read,
		RequiredField: parquet.NewRequiredField(col, opts...),
		stats:         newStringStats(),
	}
}

func (f *StringField) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Type: parquet.StringType, RepetitionType: parquet.RepetitionRequired}
}

func (f *StringField) Scan(r *Person) {
	if len(f.vals) == 0 {
		return
	}

	v := f.vals[0]
	f.vals = f.vals[1:]
	f.read(r, v)
}

func (f *StringField) Add(r Person) {
	v := f.val(r)
	f.stats.add(v)
	f.vals = append(f.vals, v)
}

func (f *StringField) Write(w io.Writer, meta *parquet.Metadata) error {
	buf := bytes.Buffer{}

	for _, s := range f.vals {
		if err := binary.Write(&buf, binary.LittleEndian, int32(len(s))); err != nil {
			return err
		}
		buf.Write([]byte(s))
	}

	return f.DoWrite(w, meta, buf.Bytes(), len(f.vals), f.stats)
}

func (f *StringField) Read(r io.ReadSeeker, pg parquet.Page) error {
	rr, _, err := f.DoRead(r, pg)
	if err != nil {
		return err
	}

	for j := 0; j < pg.N; j++ {
		var x int32
		if err := binary.Read(rr, binary.LittleEndian, &x); err != nil {
			return err
		}
		s := make([]byte, x)
		if _, err := rr.Read(s); err != nil {
			return err
		}

		f.vals = append(f.vals, string(s))
	}
	return nil
}

type BoolField struct {
	parquet.RequiredField
	vals []bool
	val  func(r Person) bool
	read func(r *Person, v bool)
}

func NewBoolField(val func(r Person) bool, read func(r *Person, v bool), col string, opts ...func(*parquet.RequiredField)) *BoolField {
	return &BoolField{
		val:           val,
		read:          read,
		RequiredField: parquet.NewRequiredField(col, opts...),
	}
}

func (f *BoolField) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Type: parquet.BoolType, RepetitionType: parquet.RepetitionRequired}
}

func (f *BoolField) Scan(r *Person) {
	if len(f.vals) == 0 {
		return
	}

	v := f.vals[0]
	f.vals = f.vals[1:]
	f.read(r, v)
}

func (f *BoolField) Add(r Person) {
	f.vals = append(f.vals, f.val(r))
}

func (f *BoolField) Write(w io.Writer, meta *parquet.Metadata) error {
	ln := len(f.vals)
	byteNum := (ln + 7) / 8
	rawBuf := make([]byte, byteNum)

	for i := 0; i < ln; i++ {
		if f.vals[i] {
			rawBuf[i/8] = rawBuf[i/8] | (1 << uint32(i%8))
		}
	}

	return f.DoWrite(w, meta, rawBuf, len(f.vals), newBoolStats())
}

func (f *BoolField) Read(r io.ReadSeeker, pg parquet.Page) error {
	rr, sizes, err := f.DoRead(r, pg)
	if err != nil {
		return err
	}

	f.vals, err = parquet.GetBools(rr, int(pg.N), sizes)
	return err
}

type int32stats struct {
	min int32
	max int32
}

func newInt32stats() *int32stats {
	return &int32stats{
		min: int32(math.MaxInt32),
	}
}

func (i *int32stats) add(val int32) {
	if val < i.min {
		i.min = val
	}
	if val > i.max {
		i.max = val
	}
}

func (f *int32stats) bytes(val int32) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, val)
	return buf.Bytes()
}

func (f *int32stats) NullCount() *int64 {
	return nil
}

func (f *int32stats) DistinctCount() *int64 {
	return nil
}

func (f *int32stats) Min() []byte {
	return f.bytes(f.min)
}

func (f *int32stats) Max() []byte {
	return f.bytes(f.max)
}

type int32optionalStats struct {
	min     int32
	max     int32
	nils    int64
	nonNils int64
}

func newint32optionalStats() *int32optionalStats {
	return &int32optionalStats{
		min: int32(math.MaxInt32),
	}
}

func (f *int32optionalStats) add(val *int32) {
	if val == nil {
		f.nils++
		return
	}

	f.nonNils++
	if *val < f.min {
		f.min = *val
	}
	if *val > f.max {
		f.max = *val
	}
}

func (f *int32optionalStats) bytes(val int32) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, val)
	return buf.Bytes()
}

func (f *int32optionalStats) NullCount() *int64 {
	return &f.nils
}

func (f *int32optionalStats) DistinctCount() *int64 {
	return nil
}

func (f *int32optionalStats) Min() []byte {
	if f.nonNils == 0 {
		return nil
	}
	return f.bytes(f.min)
}

func (f *int32optionalStats) Max() []byte {
	if f.nonNils == 0 {
		return nil
	}
	return f.bytes(f.max)
}

type int64stats struct {
	min int64
	max int64
}

func newInt64stats() *int64stats {
	return &int64stats{
		min: int64(math.MaxInt64),
	}
}

func (i *int64stats) add(val int64) {
	if val < i.min {
		i.min = val
	}
	if val > i.max {
		i.max = val
	}
}

func (f *int64stats) bytes(val int64) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, val)
	return buf.Bytes()
}

func (f *int64stats) NullCount() *int64 {
	return nil
}

func (f *int64stats) DistinctCount() *int64 {
	return nil
}

func (f *int64stats) Min() []byte {
	return f.bytes(f.min)
}

func (f *int64stats) Max() []byte {
	return f.bytes(f.max)
}

type int64optionalStats struct {
	min     int64
	max     int64
	nils    int64
	nonNils int64
}

func newint64optionalStats() *int64optionalStats {
	return &int64optionalStats{
		min: int64(math.MaxInt64),
	}
}

func (f *int64optionalStats) add(val *int64) {
	if val == nil {
		f.nils++
		return
	}

	f.nonNils++
	if *val < f.min {
		f.min = *val
	}
	if *val > f.max {
		f.max = *val
	}
}

func (f *int64optionalStats) bytes(val int64) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, val)
	return buf.Bytes()
}

func (f *int64optionalStats) NullCount() *int64 {
	return &f.nils
}

func (f *int64optionalStats) DistinctCount() *int64 {
	return nil
}

func (f *int64optionalStats) Min() []byte {
	if f.nonNils == 0 {
		return nil
	}
	return f.bytes(f.min)
}

func (f *int64optionalStats) Max() []byte {
	if f.nonNils == 0 {
		return nil
	}
	return f.bytes(f.max)
}

type stringOptionalStats struct {
	vals []string
	min  []byte
	max  []byte
	nils int64
}

func newStringOptionalStats() *stringOptionalStats {
	return &stringOptionalStats{}
}

func (s *stringOptionalStats) add(val *string) {
	if val == nil {
		s.nils++
		return
	}
	s.vals = append(s.vals, *val)
}

func (s *stringOptionalStats) NullCount() *int64 {
	return &s.nils
}

func (s *stringOptionalStats) DistinctCount() *int64 {
	return nil
}

func (s *stringOptionalStats) Min() []byte {
	if s.min == nil {
		s.minMax()
	}
	return s.min
}

func (s *stringOptionalStats) Max() []byte {
	if s.max == nil {
		s.minMax()
	}
	return s.max
}

func (s *stringOptionalStats) minMax() {
	if len(s.vals) == 0 {
		return
	}

	tmp := make([]string, len(s.vals))
	copy(tmp, s.vals)
	sort.Strings(tmp)
	s.min = []byte(tmp[0])
	s.max = []byte(tmp[len(tmp)-1])
}

type float32stats struct {
	min float32
	max float32
}

func newFloat32stats() *float32stats {
	return &float32stats{
		min: float32(math.MaxFloat32),
	}
}

func (i *float32stats) add(val float32) {
	if val < i.min {
		i.min = val
	}
	if val > i.max {
		i.max = val
	}
}

func (f *float32stats) bytes(val float32) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, val)
	return buf.Bytes()
}

func (f *float32stats) NullCount() *int64 {
	return nil
}

func (f *float32stats) DistinctCount() *int64 {
	return nil
}

func (f *float32stats) Min() []byte {
	return f.bytes(f.min)
}

func (f *float32stats) Max() []byte {
	return f.bytes(f.max)
}

type float64stats struct {
	min float64
	max float64
}

func newFloat64stats() *float64stats {
	return &float64stats{
		min: float64(math.MaxFloat64),
	}
}

func (i *float64stats) add(val float64) {
	if val < i.min {
		i.min = val
	}
	if val > i.max {
		i.max = val
	}
}

func (f *float64stats) bytes(val float64) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, val)
	return buf.Bytes()
}

func (f *float64stats) NullCount() *int64 {
	return nil
}

func (f *float64stats) DistinctCount() *int64 {
	return nil
}

func (f *float64stats) Min() []byte {
	return f.bytes(f.min)
}

func (f *float64stats) Max() []byte {
	return f.bytes(f.max)
}

type float32optionalStats struct {
	min     float32
	max     float32
	nils    int64
	nonNils int64
}

func newfloat32optionalStats() *float32optionalStats {
	return &float32optionalStats{
		min: float32(math.MaxFloat32),
	}
}

func (f *float32optionalStats) add(val *float32) {
	if val == nil {
		f.nils++
		return
	}

	f.nonNils++
	if *val < f.min {
		f.min = *val
	}
	if *val > f.max {
		f.max = *val
	}
}

func (f *float32optionalStats) bytes(val float32) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, val)
	return buf.Bytes()
}

func (f *float32optionalStats) NullCount() *int64 {
	return &f.nils
}

func (f *float32optionalStats) DistinctCount() *int64 {
	return nil
}

func (f *float32optionalStats) Min() []byte {
	if f.nonNils == 0 {
		return nil
	}
	return f.bytes(f.min)
}

func (f *float32optionalStats) Max() []byte {
	if f.nonNils == 0 {
		return nil
	}
	return f.bytes(f.max)
}

type boolOptionalStats struct {
	nils int64
}

func newBoolOptionalStats() *boolOptionalStats {
	return &boolOptionalStats{}
}

func (b *boolOptionalStats) add(val *bool) {
	if val == nil {
		b.nils++
	}
}

func (b *boolOptionalStats) NullCount() *int64 {
	return &b.nils
}

func (b *boolOptionalStats) DistinctCount() *int64 {
	return nil
}

func (b *boolOptionalStats) Min() []byte {
	return nil
}

func (b *boolOptionalStats) Max() []byte {
	return nil
}

type uint32stats struct {
	min uint32
	max uint32
}

func newUint32stats() *uint32stats {
	return &uint32stats{
		min: uint32(math.MaxUint32),
	}
}

func (i *uint32stats) add(val uint32) {
	if val < i.min {
		i.min = val
	}
	if val > i.max {
		i.max = val
	}
}

func (f *uint32stats) bytes(val uint32) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, val)
	return buf.Bytes()
}

func (f *uint32stats) NullCount() *int64 {
	return nil
}

func (f *uint32stats) DistinctCount() *int64 {
	return nil
}

func (f *uint32stats) Min() []byte {
	return f.bytes(f.min)
}

func (f *uint32stats) Max() []byte {
	return f.bytes(f.max)
}

type uint64optionalStats struct {
	min     uint64
	max     uint64
	nils    int64
	nonNils int64
}

func newuint64optionalStats() *uint64optionalStats {
	return &uint64optionalStats{
		min: uint64(math.MaxUint64),
	}
}

func (f *uint64optionalStats) add(val *uint64) {
	if val == nil {
		f.nils++
		return
	}

	f.nonNils++
	if *val < f.min {
		f.min = *val
	}
	if *val > f.max {
		f.max = *val
	}
}

func (f *uint64optionalStats) bytes(val uint64) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, val)
	return buf.Bytes()
}

func (f *uint64optionalStats) NullCount() *int64 {
	return &f.nils
}

func (f *uint64optionalStats) DistinctCount() *int64 {
	return nil
}

func (f *uint64optionalStats) Min() []byte {
	if f.nonNils == 0 {
		return nil
	}
	return f.bytes(f.min)
}

func (f *uint64optionalStats) Max() []byte {
	if f.nonNils == 0 {
		return nil
	}
	return f.bytes(f.max)
}

type stringStats struct {
	vals []string
	min  []byte
	max  []byte
}

func newStringStats() *stringStats {
	return &stringStats{}
}

func (s *stringStats) add(val string) {
	s.vals = append(s.vals, val)
}

func (s *stringStats) NullCount() *int64 {
	return nil
}

func (s *stringStats) DistinctCount() *int64 {
	return nil
}

func (s *stringStats) Min() []byte {
	if s.min == nil {
		s.minMax()
	}
	return s.min
}

func (s *stringStats) Max() []byte {
	if s.max == nil {
		s.minMax()
	}
	return s.max
}

func (s *stringStats) minMax() {
	if len(s.vals) == 0 {
		return
	}

	tmp := make([]string, len(s.vals))
	copy(tmp, s.vals)
	sort.Strings(tmp)
	s.min = []byte(tmp[0])
	s.max = []byte(tmp[len(tmp)-1])
}

type boolStats struct{}

func newBoolStats() *boolStats             { return &boolStats{} }
func (b *boolStats) NullCount() *int64     { return nil }
func (b *boolStats) DistinctCount() *int64 { return nil }
func (b *boolStats) Min() []byte           { return nil }
func (b *boolStats) Max() []byte           { return nil }
