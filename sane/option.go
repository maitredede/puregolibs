package sane

// Option represents a scanning option.
type Option struct {
	Name         string        // option name
	Group        string        // option group
	Title        string        // option title
	Desc         string        // option description
	Type         Type          // option type
	Unit         Unit          // units
	Length       int           // vector length for vector-valued options
	ConstrSet    []interface{} // constraint set
	ConstrRange  *Range        // constraint range
	IsActive     bool          // whether option is active
	IsSettable   bool          // whether option can be set
	IsDetectable bool          // whether option value can be detected
	IsAutomatic  bool          // whether option has an auto value
	IsEmulated   bool          // whether option is emulated
	IsAdvanced   bool          // whether option is advanced
	index        int           // internal option index
	size         int           // internal option size in bytes
}
