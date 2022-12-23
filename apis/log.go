package apis

// Log Defines the contract for logging.
type Log struct {
	// +optional
	// Level is the log level. The default is "info". Options are "debug", "info", "warn", "error", "panic", and "fatal".
	Level string `json:"level,omitempty"`
	// +optional
	Format string `json:"format,omitempty"`
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=10
	// +kubebuilder:validation:Enum=0;1;2;3;4;5;6;7;8;9;10
	// +kubebuilder:default=0
	Verbosity int `json:"verbosity,omitempty"`

	// RawMessage defines the log message of the log.
	RawMessage []byte `json:"logMessage,omitempty"`

	// MaxSize defines the maximum size of the log.
	MaxSize int `json:"maxSize,omitempty"`

	// Source defines the source of the log. e.g. PodName, etc.
	Source string `json:"source,omitempty"`

	// TimeFormat defines the time format of the log. e.g. 2006-01-02 15:04:05.000, etc.
	TimeFormat string `json:"timeFormat,omitempty"`

	// TimeStamp defines the time stamp of the log. e.g. 2006-01-02 15:04:05.000, etc.
	TimeStamp string `json:"timeStamp,omitempty"`

	// MaxAge defines the max age of the log. e.g. 7, etc. (7 days)
	MaxAge int `json:"maxAge,omitempty"`

	// Compress defines the compress of the log. e.g. true, etc.
	Compress bool `json:"compress,omitempty"`

	// CompressAlgorithm defines the compress algorithm of the log. e.g. gzip, etc.
	CompressAlgorithm string `json:"compressAlgorithm,omitempty"`

	// CompressedMessage defines the compressed message of the log.
	CompressedMessage []byte `json:"compressedMessage,omitempty"`
}

func (in *Log) DeepCopyInto(out *Log) {
	*out = *in
	if in.RawMessage != nil {
		in, out := &in.RawMessage, &out.RawMessage
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.CompressedMessage != nil {
		in, out := &in.CompressedMessage, &out.CompressedMessage
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
}
