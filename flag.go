package gg

import (
	"encoding"
	r "reflect"
)

/*
Parses CLI flags into an instance of the given type, which must be a struct.
For parsing rules, see `FlagParser`.
*/
func FlagParseTo[A any](src []string) (out A) {
	FlagParse(src, &out)
	return
}

/*
Parses CLI flags into the given value, which must be a struct.
Panics on error. For parsing rules, see `FlagParser`.
*/
func FlagParse[A any](src []string, out *A) {
	if out != nil {
		FlagParseReflect(src, r.ValueOf(AnyNoEscUnsafe(out)).Elem())
	}
}

/*
Parses CLI flags into the given value, which must be a struct.
For parsing rules, see `FlagParser`.
*/
func FlagParseCatch[A any](src []string, out *A) (err error) {
	defer Rec(&err)
	FlagParse(src, out)
	return
}

/*
Parses CLI flags into the given output, which must be a settable struct value.
For parsing rules, see `FlagParser`.
*/
func FlagParseReflect(src []string, out r.Value) {
	if !out.IsValid() {
		return
	}

	var par FlagParser
	par.Init(out)
	par.Args(src)
	par.Default()
}

/*
Tool for parsing lists of CLI flags into structs. Partial replacement for the
standard library package "flag". Example:

	type Opt struct {
		Args []string `flag:""`
		Help bool     `flag:"-h"          desc:"Print help and exit."`
		Verb bool     `flag:"-v"          desc:"Verbose logging."`
		Src  string   `flag:"-s" init:"." desc:"Source path."`
		Out  string   `flag:"-o"          desc:"Destination path."`
	}

	func (self *Opt) Init() {
		gg.FlagParse(os.Args[1:], self)

		if self.Help {
			log.Println(gg.FlagHelp[Opt]())
			os.Exit(0)
		}

		if gg.IsZero(self.Out) {
			log.Println(`missing output path: "-o"`)
			os.Exit(1)
		}
	}

Supported struct tags:

	* `flag`: must be "" or a valid flag like "-v" or "--verbose".
	  Fields without the `flag` tag are ignored. Flags must be unique.
	* Field with `flag:""` is used for remaining non-flag args.
	  It must have a type convertible to `[]string`.
	* `init`: initial value. Used if the flag was not provided.
	* `desc`: description. Used for help printing.

Parsing rules:

	* Supports all primitive types.
	* Supports slices of arbitrary types.
	* Supports `gg.Parser`.
	* Supports `encoding.TextUnmarshaler`.
	* Supports `flag.Value`.
	* Each flag may be listed multiple times.
		* If the target is a parser, invoke its parsing method.
		* If the target is a scalar, replace the old value with the new value.
		* If the target is a slice, append the new value.
*/
type FlagParser struct {
	Tar r.Value
	Def FlagDef
	Got Set[string]
}

/*
Initializes the parser for the given destination, which must be a settable
struct value.
*/
func (self *FlagParser) Init(tar r.Value) {
	self.Tar = tar
	self.Def = FlagDefCache.Get(tar.Type())
	self.Got = make(Set[string], len(self.Def.Flags))
}

/*
Parses the given CLI args into the destination. May be called multiple times.
Must be called after `(*FlagParser).Init`, and before `FlagParser.Default`.
*/
func (self FlagParser) Args(src []string) {
	for HasLen(src) {
		if !isCliFlag(Head(src)) {
			self.SetArgs(src)
			return
		}

		head := PopHead(&src)
		key, val, split := cliFlagSplit(head)
		if split {
			self.Got.Add(key)
			self.Flag(key, val)
			continue
		}

		self.Got.Add(head)

		if IsEmpty(src) || isCliFlag(Head(src)) {
			self.TrailingFlag(head)
			continue
		}
		if self.TrailingBool(head) {
			continue
		}

		self.Flag(head, PopHead(&src))
	}
}

// For internal use.
func (self FlagParser) SetArgs(src []string) {
	field := self.Def.Args

	if field.IsValid() {
		self.Tar.
			FieldByIndex(field.Index).
			Set(r.ValueOf(src).Convert(field.Type))
		return
	}

	if IsEmpty(src) {
		return
	}

	panic(Errf(`unexpected non-flag args: %q`, src))
}

// For internal use.
func (self FlagParser) FlagField(key string) r.Value {
	return self.Tar.FieldByIndex(self.Def.Get(key).Index)
}

// For internal use.
func (self FlagParser) Flag(key, src string) {
	self.FieldParse(src, self.FlagField(key))
}

// For internal use.
func (self FlagParser) FieldParse(src string, out r.Value) {
	var nested bool

interfaces:
	ptr := out.Addr().Interface()

	// Part of the `flag.Value` interface.
	setter, _ := ptr.(interface{ Set(string) error })
	if setter != nil {
		Try(setter.Set(src))
		return
	}

	parser, _ := ptr.(Parser)
	if parser != nil {
		Try(parser.Parse(src))
		return
	}

	unmarshaler, _ := ptr.(encoding.TextUnmarshaler)
	if unmarshaler != nil {
		Try(unmarshaler.UnmarshalText(ToBytes(src)))
		return
	}

	if out.Kind() == r.Slice {
		growLenReflect(out)
		out = out.Index(out.Len() - 1)

		if !nested {
			nested = true
			goto interfaces
		}
	}

	if out.Kind() == r.Bool && src == `` {
		src = `true`
	}
	Try(ParseReflectCatch(src, out.Addr()))
}

// For internal use.
func (self FlagParser) TrailingFlag(key string) {
	// TODO: consider supporting various parser interfaces here.
	if self.TrailingBool(key) {
		return
	}
	panic(Errf(`missing value for trailing flag %q`, key))
}

// For internal use.
func (self FlagParser) TrailingBool(key string) bool {
	/**
	Following the established conventions, bool flags don't support
	"-flag value", only "-flag=value". A boolean flag always terminates
	immediately, without looking for a following space-separated value.
	*/

	tar := self.FlagField(key)

	if tar.Kind() == r.Bool {
		tar.SetBool(true)
		return true
	}

	if tar.Kind() == r.Slice && tar.Type().Elem().Kind() == r.Bool {
		growLenReflect(tar)
		tar.Index(tar.Len() - 1).SetBool(true)
		return true
	}

	return false
}

/*
Applies defaults to all flags which have not been found during parsing.
Explicitly providing an empty value suppresses a default, although
an empty string may not be a viable input to some types.
*/
func (self FlagParser) Default() {
	for _, field := range self.Def.Flags {
		if !self.Got.Has(field.Flag) {
			if field.InitHas {
				self.FieldParse(field.Init, self.Tar.FieldByIndex(field.Index))
			}
		}
	}
}

// Returns a help string for the given struct type, using `FlagFmtDefault`.
func FlagHelp[A any]() string {
	return FlagDefCache.Get(Type[A]()).Help()
}

// Stores cached `FlagDef` definitions for struct types.
var FlagDefCache = TypeCacheOf[FlagDef]()

/*
Struct type definition suitable for flag parsing. Used internally by
`FlagParser`. User code shouldn't have to use this type, but it's exported for
customization purposes.
*/
type FlagDef struct {
	Type  r.Type
	Args  FlagDefField
	Flags []FlagDefField
	Index map[string]int
}

// For internal use.
func (self *FlagDef) Init(src r.Type) {
	self.Type = src
	Each(StructDeepPublicFieldCache.Get(src), self.AddField)
}

// For internal use.
func (self *FlagDef) AddField(src r.StructField) {
	var field FlagDefField
	field.Set(src)
	if !field.FlagHas {
		return
	}

	if MapHas(self.Index, field.Flag) ||
		(field.Flag == `` && self.Args.IsValid()) {
		panic(Errf(`redundant flag %q in type %v`, field.Flag, self.Type))
	}

	if field.Flag == `` {
		if !field.Type.ConvertibleTo(Type[[]string]()) {
			panic(Errf(
				`invalid type %v in field %q of type %v: args field must be convertible to []string`,
				field.Type, field.Name, self.Type,
			))
		}
		self.Args = field
		return
	}

	if !isCliFlagValid(field.Flag) {
		panic(Errf(
			`invalid flag %q in field %q of type %v`,
			field.Flag, field.Name, self.Type,
		))
	}

	MapInit(&self.Index)[field.Flag] = len(self.Flags)
	AppendVals(&self.Flags, field)
}

// For internal use.
func (self FlagDef) Got(key string) (FlagDefField, bool) {
	ind, ok := self.Index[key]
	if !ok {
		return Zero[FlagDefField](), false
	}
	return Got(self.Flags, ind)
}

// For internal use.
func (self FlagDef) Get(key string) FlagDefField {
	val, ok := self.Got(key)
	if !ok {
		panic(Errf(`unable to find flag %q in type %v`, key, self.Type))
	}
	return val
}

// Creates a help string listing the available flags, using `FlagFmtDefault`.
func (self FlagDef) Help() string { return FlagFmtDefault.String(self) }

// Used internally by `FlagDef`.
type FlagDefField struct {
	r.StructField
	Flag    string
	FlagHas bool
	FlagLen int
	Init    string
	InitHas bool
	InitLen int
	Desc    string
	DescHas bool
	DescLen int
}

func (self FlagDefField) IsValid() bool { return IsNonZero(self) }

func (self *FlagDefField) Set(src r.StructField) {
	self.StructField = src

	self.Flag, self.FlagHas = self.Tag.Lookup(`flag`)
	self.Init, self.InitHas = self.Tag.Lookup(`init`)
	self.Desc = self.Tag.Get(`desc`)

	self.FlagLen = CharCount(self.Flag)
	self.InitLen = CharCount(self.Init)
	self.DescLen = CharCount(self.Desc)

	self.DescHas = self.DescLen > 0
}

func (self FlagDefField) GetFlagHas() bool { return self.FlagHas }
func (self FlagDefField) GetInitHas() bool { return self.InitHas }
func (self FlagDefField) GetDescHas() bool { return self.DescHas }

func (self FlagDefField) GetFlagLen() int { return self.FlagLen }
func (self FlagDefField) GetInitLen() int { return self.InitLen }
func (self FlagDefField) GetDescLen() int { return self.DescLen }

// Default help formatter, used by `FlagHelp` and `FlagDef.Help`.
var FlagFmtDefault = With((*FlagFmt).Default)

/*
Table-like formatter for listing available flags, initial values, and
descriptions. Used via `FlagFmtDefault`, `FlagHelp`, `FlagDef.Help`.
To customize printing, mutate `FlagFmtDefault`.
*/
type FlagFmt struct {
	Prefix    string // Prepended before each line.
	Infix     string // Inserted between columns.
	Head      bool   // If true, print table header.
	FlagHead  string // Title for header cell "flag".
	InitHead  string // Title for header cell "init".
	DescHead  string // Title for header cell "desc".
	HeadUnder string // Separator between table header and body.
}

// Sets default values.
func (self *FlagFmt) Default() {
	self.Infix = `    `
	self.Head = true
	self.FlagHead = `flag`
	self.InitHead = `init`
	self.DescHead = `desc`
	self.HeadUnder = `-`
}

// Returns a table-like help string for the given definition.
func (self FlagFmt) String(def FlagDef) string {
	return ToString(self.Append(nil, def))
}

/*
Appends table-like help for the given definition. Known limitation: assumes
monospace, doesn't support wider characters such as kanji or emoji.
*/
func (self FlagFmt) Append(src []byte, def FlagDef) []byte {
	flags := def.Flags
	if IsEmpty(flags) {
		return src
	}

	prefixLen := CharCount(self.Prefix)
	sepLen := CharCount(self.Infix)
	newlineLen := CharCount(Newline)
	flagLen := MaxPrimBy(flags, FlagDefField.GetFlagLen)

	var flagHeadLen int
	if self.Head {
		flagHeadLen = CharCount(self.FlagHead)
		flagLen = MaxPrim2(flagHeadLen, flagLen)
	}

	var initHeadLen int
	var initLen int
	if Some(flags, FlagDefField.GetInitHas) {
		if self.Head {
			initHeadLen = CharCount(self.InitHead)
		}
		initLen = MaxPrim2(
			initHeadLen,
			MaxPrimBy(flags, FlagDefField.GetInitLen),
		)
	}

	var initLenOuter int
	if initLen > 0 {
		initLenOuter = sepLen + initLen
	}

	var descHeadLen int
	var descLen int
	if Some(flags, FlagDefField.GetDescHas) {
		if self.Head {
			descHeadLen = CharCount(self.DescHead)
		}
		descLen = MaxPrim2(
			descHeadLen,
			MaxPrimBy(flags, FlagDefField.GetDescLen),
		)
	}

	var descLenOuter int
	if descLen > 0 {
		descLenOuter = sepLen + descLen
	}

	rowLenInner := prefixLen + flagLen + initLenOuter + descLenOuter
	rowLen := rowLenInner + newlineLen
	headUnderLen := CharCount(self.HeadUnder)

	buf := Buf(src)
	buf.GrowCap(((2 + len(flags)) * rowLen))

	if self.Head {
		buf.AppendString(self.Prefix)

		buf.AppendString(self.FlagHead)
		buf.AppendSpaces(flagLen - flagHeadLen)
		if initLen > 0 || descLen > 0 {
			buf.AppendSpaces(flagLen - flagHeadLen)
		}

		if initLen > 0 {
			buf.AppendString(self.Infix)
			buf.AppendString(self.InitHead)
			if descLen > 0 {
				buf.AppendSpaces(initLen - initHeadLen)
			}
		}

		if descLen > 0 {
			buf.AppendString(self.Infix)
			buf.AppendString(self.DescHead)
		}

		buf.AppendString(Newline)

		if rowLenInner > 0 && headUnderLen > 0 {
			buf.AppendString(self.Prefix)
			buf.AppendStringN(self.HeadUnder, rowLenInner/headUnderLen)
			buf.AppendString(Newline)
		}
	}

	for _, field := range flags {
		buf.AppendString(self.Prefix)

		buf.AppendString(field.Flag)
		if initLen > 0 || descLen > 0 {
			buf.AppendSpaces(flagLen - field.FlagLen)
		}

		if initLen > 0 {
			buf.AppendString(self.Infix)
			buf.AppendString(field.Init)
			if descLen > 0 {
				buf.AppendSpaces(initLen - field.InitLen)
			}
		}

		if descLen > 0 {
			buf.AppendString(self.Infix)
			buf.AppendString(field.Desc)
		}
		buf.AppendString(Newline)
	}

	return buf
}
