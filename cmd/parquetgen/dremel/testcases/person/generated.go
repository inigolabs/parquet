package person

// Code generated by github.com/inigolabs/parquet.  DO NOT EDIT.

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/inigolabs/parquet"
	sch "github.com/inigolabs/parquet/schema"
	"github.com/valyala/bytebufferpool"
)

var _ = math.MaxInt32 // to avoid unused import

type compression int

const (
	compressionUncompressed compression = 0
	compressionSnappy       compression = 1
	compressionGzip         compression = 2
	compressionUnknown      compression = -1
)

var buffpool = bytebufferpool.Pool{}

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

func Fields(compression compression) []Field {
	return []Field{
		NewStringField(readName, writeName, []string{"name"}, fieldCompression(compression)),
		NewStringOptionalField(readHobbyName, writeHobbyName, []string{"hobby", "name"}, []int{1, 0}, optionalFieldCompression(compression)),
		NewInt32OptionalField(readHobbyDifficulty, writeHobbyDifficulty, []string{"hobby", "difficulty"}, []int{1, 1}, optionalFieldCompression(compression)),
		NewStringOptionalField(readHobbySkillsName, writeHobbySkillsName, []string{"hobby", "skills", "name"}, []int{1, 2, 0}, optionalFieldCompression(compression)),
		NewStringOptionalField(readHobbySkillsDifficulty, writeHobbySkillsDifficulty, []string{"hobby", "skills", "difficulty"}, []int{1, 2, 0}, optionalFieldCompression(compression)),
	}
}

func readName(x Person) string {
	return x.Name
}

func writeName(x *Person, vals []string) {
	x.Name = vals[0]
}

func readHobbyName(x Person, vals []string, defs, reps []uint8) ([]string, []uint8, []uint8) {
	switch {
	case x.Hobby == nil:
		defs = append(defs, 0)
		return vals, defs, reps
	default:
		vals = append(vals, x.Hobby.Name)
		defs = append(defs, 1)
		return vals, defs, reps
	}
}

func writeHobbyName(x *Person, vals []string, defs, reps []uint8) (int, int) {
	def := defs[0]
	switch def {
	case 1:
		x.Hobby = &Hobby{Name: vals[0]}
		return 1, 1
	}

	return 0, 1
}

func readHobbyDifficulty(x Person, vals []int32, defs, reps []uint8) ([]int32, []uint8, []uint8) {
	switch {
	case x.Hobby == nil:
		defs = append(defs, 0)
		return vals, defs, reps
	case x.Hobby.Difficulty == nil:
		defs = append(defs, 1)
		return vals, defs, reps
	default:
		vals = append(vals, *x.Hobby.Difficulty)
		defs = append(defs, 2)
		return vals, defs, reps
	}
}

func writeHobbyDifficulty(x *Person, vals []int32, defs, reps []uint8) (int, int) {
	def := defs[0]
	switch def {
	case 2:
		x.Hobby.Difficulty = pint32(vals[0])
		return 1, 1
	}

	return 0, 1
}

func readHobbySkillsName(x Person, vals []string, defs, reps []uint8) ([]string, []uint8, []uint8) {
	var lastRep uint8

	if x.Hobby == nil {
		defs = append(defs, 0)
		reps = append(reps, lastRep)
	} else {
		if len(x.Hobby.Skills) == 0 {
			defs = append(defs, 1)
			reps = append(reps, lastRep)
		} else {
			for i0, x0 := range x.Hobby.Skills {
				if i0 >= 1 {
					lastRep = 1
				}
				defs = append(defs, 2)
				reps = append(reps, lastRep)
				vals = append(vals, x0.Name)
			}
		}
	}

	return vals, defs, reps
}

func writeHobbySkillsName(x *Person, vals []string, defs, reps []uint8) (int, int) {
	var nVals, nLevels int
	ind := make(indices, 1)

	for i := range defs {
		def := defs[i]
		rep := reps[i]
		if i > 0 && rep == 0 {
			break
		}

		nLevels++
		ind.rep(rep)

		switch def {
		case 2:
			x.Hobby.Skills = append(x.Hobby.Skills, Skill{Name: vals[nVals]})
			nVals++
		}
	}

	return nVals, nLevels
}

func readHobbySkillsDifficulty(x Person, vals []string, defs, reps []uint8) ([]string, []uint8, []uint8) {
	var lastRep uint8

	if x.Hobby == nil {
		defs = append(defs, 0)
		reps = append(reps, lastRep)
	} else {
		if len(x.Hobby.Skills) == 0 {
			defs = append(defs, 1)
			reps = append(reps, lastRep)
		} else {
			for i0, x0 := range x.Hobby.Skills {
				if i0 >= 1 {
					lastRep = 1
				}
				defs = append(defs, 2)
				reps = append(reps, lastRep)
				vals = append(vals, x0.Difficulty)
			}
		}
	}

	return vals, defs, reps
}

func writeHobbySkillsDifficulty(x *Person, vals []string, defs, reps []uint8) (int, int) {
	var nVals, nLevels int
	ind := make(indices, 1)

	for i := range defs {
		def := defs[i]
		rep := reps[i]
		if i > 0 && rep == 0 {
			break
		}

		nLevels++
		ind.rep(rep)

		switch def {
		case 2:
			x.Hobby.Skills[ind[0]].Difficulty = vals[nVals]
			nVals++
		}
	}

	return nVals, nLevels
}

func fieldCompression(c compression) func(*parquet.RequiredField) {
	switch c {
	case compressionUncompressed:
		return parquet.RequiredFieldUncompressed
	case compressionSnappy:
		return parquet.RequiredFieldSnappy
	case compressionGzip:
		return parquet.RequiredFieldGzip
	default:
		return parquet.RequiredFieldUncompressed
	}
}

func optionalFieldCompression(c compression) func(*parquet.OptionalField) {
	switch c {
	case compressionUncompressed:
		return parquet.OptionalFieldUncompressed
	case compressionSnappy:
		return parquet.OptionalFieldSnappy
	case compressionGzip:
		return parquet.OptionalFieldGzip
	default:
		return parquet.OptionalFieldUncompressed
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

var par1 = []byte("PAR1")

func begin(p *ParquetWriter) error {
	_, err := p.w.Write(par1)
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

func Gzip(p *ParquetWriter) error {
	p.compression = compressionGzip
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

	_, err := p.w.Write(par1)
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

	p.meta.NextDoc()
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
	Levels() ([]uint8, []uint8)
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
		pr.fieldNames = append(pr.fieldNames, f.Name())
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
	fieldNames     []string
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

type Levels struct {
	Name string
	Defs []uint8
	Reps []uint8
}

func (p *ParquetReader) Levels() []Levels {
	var out []Levels
	//for {
	for _, name := range p.fieldNames {
		f := p.fields[name]
		d, r := f.Levels()
		out = append(out, Levels{Name: f.Name(), Defs: d, Reps: r})
	}
	//	if err := p.readRowGroup(); err != nil {
	//		break
	//	}
	//}
	return out
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
		name := strings.Join(col.MetaData.PathInSchema, ".")
		f, ok := p.fields[name]
		if !ok {
			return fmt.Errorf("unknown field: %s", name)
		}
		pages := p.pages[name]
		if len(pages) <= p.index {
			break
		}

		pg := pages[0]
		if err := f.Read(p.r, pg); err != nil {
			return fmt.Errorf("unable to read field %s, err: %s", f.Name(), err)
		}
		p.pages[name] = p.pages[name][1:]
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

	for _, name := range p.fieldNames {
		f := p.fields[name]
		f.Scan(x)
	}
}

type StringField struct {
	parquet.RequiredField
	vals  []string
	read  func(r Person) string
	write func(r *Person, vals []string)
	stats *stringStats
}

func NewStringField(read func(r Person) string, write func(r *Person, vals []string), path []string, opts ...func(*parquet.RequiredField)) *StringField {
	return &StringField{
		read:          read,
		write:         write,
		RequiredField: parquet.NewRequiredField(path, opts...),
		stats:         newStringStats(),
	}
}

func (f *StringField) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Path: f.Path(), Type: StringType, RepetitionType: parquet.RepetitionRequired, Types: []int{0}}
}

func (f *StringField) Write(w io.Writer, meta *parquet.Metadata) error {
	buf := buffpool.Get()
	defer buffpool.Put(buf)

	bs := make([]byte, 4)
	for _, s := range f.vals {
		binary.LittleEndian.PutUint32(bs, uint32(len(s)))
		if _, err := buf.Write(bs); err != nil {
			return err
		}
		buf.WriteString(s)
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

func (f *StringField) Scan(r *Person) {
	if len(f.vals) == 0 {
		return
	}

	f.write(r, f.vals)
	f.vals = f.vals[1:]
}

func (f *StringField) Add(r Person) {
	v := f.read(r)
	f.stats.add(v)
	f.vals = append(f.vals, v)
}

func (f *StringField) Levels() ([]uint8, []uint8) {
	return nil, nil
}

type StringOptionalField struct {
	parquet.OptionalField
	vals  []string
	read  func(r Person, vals []string, def, rep []uint8) ([]string, []uint8, []uint8)
	write func(r *Person, vals []string, def, rep []uint8) (int, int)
	stats *stringOptionalStats
}

func NewStringOptionalField(read func(r Person, vals []string, def, rep []uint8) ([]string, []uint8, []uint8), write func(r *Person, vals []string, defs, reps []uint8) (int, int), path []string, types []int, opts ...func(*parquet.OptionalField)) *StringOptionalField {
	return &StringOptionalField{
		read:          read,
		write:         write,
		OptionalField: parquet.NewOptionalField(path, types, opts...),
		stats:         newStringOptionalStats(maxDef(types)),
	}
}

func (f *StringOptionalField) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Path: f.Path(), Type: StringType, RepetitionType: f.RepetitionType, Types: f.Types}
}

func (f *StringOptionalField) Add(r Person) {
	vals, defs, reps := f.read(r, f.vals, f.Defs, f.Reps)
	f.stats.add(vals[len(f.vals):], defs[len(f.Defs):])
	f.vals = vals
	f.Defs = defs
	f.Reps = reps
}

func (f *StringOptionalField) Scan(r *Person) {
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

func (f *StringOptionalField) Write(w io.Writer, meta *parquet.Metadata) error {
	buf := buffpool.Get()
	defer buffpool.Put(buf)

	bs := make([]byte, 4)
	for _, s := range f.vals {
		binary.LittleEndian.PutUint32(bs, uint32(len(s)))
		if _, err := buf.Write(bs); err != nil {
			return err
		}
		buf.WriteString(s)
	}

	return f.DoWrite(w, meta, buf.Bytes(), len(f.Defs), f.stats)
}

func (f *StringOptionalField) Read(r io.ReadSeeker, pg parquet.Page) error {
	rr, _, err := f.DoRead(r, pg)
	if err != nil {
		return err
	}

	for j := 0; j < f.Values(); j++ {
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

func (f *StringOptionalField) Levels() ([]uint8, []uint8) {
	return f.Defs, f.Reps
}

type Int32OptionalField struct {
	parquet.OptionalField
	vals  []int32
	read  func(r Person, vals []int32, defs, reps []uint8) ([]int32, []uint8, []uint8)
	write func(r *Person, vals []int32, defs, reps []uint8) (int, int)
	stats *int32optionalStats
}

func NewInt32OptionalField(read func(r Person, vals []int32, defs, reps []uint8) ([]int32, []uint8, []uint8), write func(r *Person, vals []int32, defs, reps []uint8) (int, int), path []string, types []int, opts ...func(*parquet.OptionalField)) *Int32OptionalField {
	return &Int32OptionalField{
		read:          read,
		write:         write,
		OptionalField: parquet.NewOptionalField(path, types, opts...),
		stats:         newint32optionalStats(maxDef(types)),
	}
}

func (f *Int32OptionalField) Schema() parquet.Field {
	return parquet.Field{Name: f.Name(), Path: f.Path(), Type: Int32Type, RepetitionType: f.RepetitionType, Types: f.Types}
}

func (f *Int32OptionalField) Write(w io.Writer, meta *parquet.Metadata) error {
	buf := buffpool.Get()
	defer buffpool.Put(buf)

	bs := make([]byte, 4)
	for _, v := range f.vals {
		binary.LittleEndian.PutUint32(bs, uint32(v))
		if _, err := buf.Write(bs); err != nil {
			return err
		}
	}
	return f.DoWrite(w, meta, buf.Bytes(), len(f.Defs), f.stats)
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
	vals, defs, reps := f.read(r, f.vals, f.Defs, f.Reps)
	f.stats.add(vals[len(f.vals):], defs[len(f.Defs):])
	f.vals = vals
	f.Defs = defs
	f.Reps = reps
}

func (f *Int32OptionalField) Scan(r *Person) {
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

func (f *Int32OptionalField) Levels() ([]uint8, []uint8) {
	return f.Defs, f.Reps
}

const nilString = "__#NIL#__"

type stringStats struct {
	min string
	max string
}

func newStringStats() *stringStats {
	return &stringStats{
		min: nilString,
		max: nilString,
	}
}

func (s *stringStats) add(val string) {
	if s.min == nilString {
		s.min = val
	} else {
		if val < s.min {
			s.min = val
		}
	}
	if s.max == nilString {
		s.max = val
	} else {
		if val > s.max {
			s.max = val
		}
	}
}

func (s *stringStats) NullCount() *int64 {
	return nil
}

func (s *stringStats) DistinctCount() *int64 {
	return nil
}

func (s *stringStats) Min() []byte {
	if s.min == nilString {
		return nil
	}
	return []byte(s.min)
}

func (s *stringStats) Max() []byte {
	if s.max == nilString {
		return nil
	}
	return []byte(s.max)
}

const nilOptString = "__#NIL#__"

type stringOptionalStats struct {
	min    string
	max    string
	nils   int64
	maxDef uint8
}

func newStringOptionalStats(d uint8) *stringOptionalStats {
	return &stringOptionalStats{
		min:    nilOptString,
		max:    nilOptString,
		maxDef: d,
	}
}

func (s *stringOptionalStats) add(vals []string, defs []uint8) {
	var i int
	for _, def := range defs {
		if def < s.maxDef {
			s.nils++
		} else {
			val := vals[i]
			if s.min == nilOptString {
				s.min = val
			} else {
				if val < s.min {
					s.min = val
				}
			}
			if s.max == nilOptString {
				s.max = val
			} else {
				if val > s.max {
					s.max = val
				}
			}
			i++
		}
	}
}

func (s *stringOptionalStats) NullCount() *int64 {
	return &s.nils
}

func (s *stringOptionalStats) DistinctCount() *int64 {
	return nil
}

func (s *stringOptionalStats) Min() []byte {
	if s.min == nilOptString {
		return nil
	}
	return []byte(s.min)
}

func (s *stringOptionalStats) Max() []byte {
	if s.max == nilOptString {
		return nil
	}
	return []byte(s.max)
}

type int32optionalStats struct {
	min     int32
	max     int32
	nils    int64
	nonNils int64
	maxDef  uint8
}

func newint32optionalStats(d uint8) *int32optionalStats {
	return &int32optionalStats{
		min:    int32(math.MaxInt32),
		maxDef: d,
	}
}

func (f *int32optionalStats) add(vals []int32, defs []uint8) {
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

func (f *int32optionalStats) bytes(v int32) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(v))
	return bs
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

func pint32(i int32) *int32       { return &i }
func puint32(i uint32) *uint32    { return &i }
func pint64(i int64) *int64       { return &i }
func puint64(i uint64) *uint64    { return &i }
func pbool(b bool) *bool          { return &b }
func pstring(s string) *string    { return &s }
func pfloat32(f float32) *float32 { return &f }
func pfloat64(f float64) *float64 { return &f }

// keeps track of the indices of repeated fields
// that have already been handled by a previous field
type indices []int

func (i indices) rep(rep uint8) {
	if rep > 0 {
		r := int(rep) - 1
		i[r] = i[r] + 1
		for j := int(rep); j < len(i); j++ {
			i[j] = 0
		}
	}
}

func maxDef(types []int) uint8 {
	var out uint8
	for _, typ := range types {
		if typ > 0 {
			out++
		}
	}
	return out
}

func Int32Type(se *sch.SchemaElement) {
	t := sch.Type_INT32
	se.Type = &t
}

func Uint32Type(se *sch.SchemaElement) {
	t := sch.Type_INT32
	se.Type = &t
	ct := sch.ConvertedType_UINT_32
	se.ConvertedType = &ct
}

func Int64Type(se *sch.SchemaElement) {
	t := sch.Type_INT64
	se.Type = &t
}

func Uint64Type(se *sch.SchemaElement) {
	t := sch.Type_INT64
	se.Type = &t
	ct := sch.ConvertedType_UINT_64
	se.ConvertedType = &ct
}

func Float32Type(se *sch.SchemaElement) {
	t := sch.Type_FLOAT
	se.Type = &t
}

func Float64Type(se *sch.SchemaElement) {
	t := sch.Type_DOUBLE
	se.Type = &t
}

func BoolType(se *sch.SchemaElement) {
	t := sch.Type_BOOLEAN
	se.Type = &t
}

func StringType(se *sch.SchemaElement) {
	t := sch.Type_BYTE_ARRAY
	se.Type = &t
}
