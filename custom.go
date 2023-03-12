package jsonschema

type CustomView struct {
	SizeXs       uint8       `json:"sizeXs,omitempty" bson:"sizeXs,omitempty"`
	SizeSm       uint8       `json:"sizeSm,omitempty" bson:"sizeSm,omitempty"`
	SizeMd       uint8       `json:"sizeMd,omitempty" bson:"sizeMd,omitempty"`
	SizeLg       uint8       `json:"sizeLg,omitempty" bson:"sizeLg,omitempty"`
	SizeXl       uint8       `json:"sizeXl,omitempty" bson:"sizeXl,omitempty"`
	NoGrid       bool        `json:"noGrid,omitempty" bson:"noGrid,omitempty"`
	Spacing      int         `json:"spacing,omitempty" bson:"spacing,omitempty"`
	Rows         uint        `json:"rows,omitempty" bson:"rows,omitempty"`
	RowsMax      uint        `json:"rowsMax,omitempty" bson:"rowsMax,omitempty"`
	Variant      string      `json:"variant,omitempty" bson:"variant,omitempty"`
	Margin       string      `json:"margin,omitempty" bson:"margin,omitempty"`
	Dense        bool        `json:"dense,omitempty" bson:"dense,omitempty"`
	DenseOptions bool        `json:"denseOptions,omitempty" bson:"denseOptions,omitempty"`
	Bg           bool        `json:"bg,omitempty" bson:"bg,omitempty"`
	Shrink       bool        `json:"shrink,omitempty" bson:"shrink,omitempty"`
	Formats      []string    `json:"formats,omitempty" bson:"formats,omitempty"`
	Justify      string      `json:"justify,omitempty" bson:"justify,omitempty"`
	Marks        interface{} `json:"marks,omitempty" bson:"marks,omitempty"`
	MarksLabel   string      `json:"marksLabel,omitempty" bson:"marksLabel,omitempty"`
	Tooltip      string      `json:"tooltip,omitempty" bson:"tooltip,omitempty"`
	TopControls  bool        `json:"topControls,omitempty" bson:"topControls,omitempty"`
	Alpha        bool        `json:"alpha,omitempty" bson:"alpha,omitempty"`
	IconOn       bool        `json:"iconOn,omitempty" bson:"iconOn,omitempty"`
	Colors       []string    `json:"colors,omitempty" bson:"colors,omitempty"`
	BtnSize      string      `json:"btnSize,omitempty" bson:"btnSize,omitempty"`
	Track        interface{} `json:"track,omitempty" bson:"track,omitempty"`
	Mt           int         `json:"mt,omitempty" bson:"mt,omitempty"`
	Mb           int         `json:"mb,omitempty" bson:"mb,omitempty"`
}

type CustomDate struct {
	Format        string   `json:"format,omitempty" bson:"format,omitempty"`
	FormatDate    string   `json:"formatDate,omitempty" bson:"formatDate,omitempty"`
	Keyboard      bool     `json:"keyboard,omitempty" bson:"keyboard,omitempty"`
	Views         []string `json:"views,omitempty" bson:"views,omitempty"`
	Variant       string   `json:"variant,omitempty" bson:"variant,omitempty"`
	AutoOk        bool     `json:"autoOk,omitempty" bson:"autoOk,omitempty"`
	DisableFuture bool     `json:"disableFuture,omitempty" bson:"disableFuture,omitempty"`
	DisablePast   bool     `json:"disablePast,omitempty" bson:"disablePast,omitempty"`
	Toolbar       bool     `json:"toolbar,omitempty" bson:"toolbar,omitempty"`
	Clearable     bool     `json:"clearable,omitempty" bson:"clearable,omitempty"`
	MinDate       string   `json:"minDate,omitempty" bson:"minDate,omitempty"`
	MaxDate       string   `json:"maxDate,omitempty" bson:"maxDate,omitempty"`
	OpenTo        string   `json:"openTo,omitempty" bson:"openTo,omitempty"`
	Orientation   string   `json:"orientation,omitempty" bson:"orientation,omitempty"`
	Tabs          bool     `json:"tabs,omitempty" bson:"tabs,omitempty"`
	MinutesStep   int      `json:"minutesStep,omitempty" bson:"minutesStep,omitempty"`
	Ampm          bool     `json:"ampm,omitempty" bson:"ampm,omitempty"`
}
