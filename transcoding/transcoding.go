package transcoding

import (
	"fmt"
	"os/exec"
	"reflect"
)

const (
	OriginalFileName          string = "original"
	HLSMasterPlaylistFilename string = "master.m3u8"
	HLSPlaylistFilename       string = "index.m3u8"
)

type OriginalFile struct {
	Filepath        string  `json:"filepath"`
	DurationSeconds float64 `json:"duration_seconds"`
	Format          string  `json:"format"`
	FrameRate       uint8   `json:"frame_rate"`
}

type ITranscoder interface {
	Process() IProcess
	Run(p IProcess) error
}

type IProcess interface {
	Input(string) IProcess
	Output(string) IProcess
	WithOptions(ProcessOptions) IProcess
	GetStrArguments() []string
	WithAdditionalOptions(map[string]string) IProcess
}

// ProcessOptions defines allowed FFmpeg arguments
type ProcessOptions struct {
	Aspect                *string           `flag:"-aspect"`
	Resolution            *string           `flag:"-s"`
	VideoBitRate          *int              `flag:"-b:v"`
	VideoBitRateTolerance *int              `flag:"-bt"`
	VideoMaxBitRate       *int              `flag:"-maxrate"`
	VideoMinBitrate       *int              `flag:"-minrate"`
	VideoCodec            *string           `flag:"-c:v"`
	Vframes               *int              `flag:"-vframes"`
	FrameRate             *int              `flag:"-r"`
	AudioRate             *int              `flag:"-ar"`
	KeyframeInterval      *int              `flag:"-g"`
	AudioCodec            *string           `flag:"-c:a"`
	AudioBitrate          *int              `flag:"-ab"`
	AudioChannels         *int              `flag:"-ac"`
	AudioVariableBitrate  *bool             `flag:"-q:a"`
	BufferSize            *int              `flag:"-bufsize"`
	Threadset             *bool             `flag:"-threads"`
	Threads               *int              `flag:"-threads"`
	Preset                *string           `flag:"-preset"`
	Tune                  *string           `flag:"-tune"`
	AudioProfile          *string           `flag:"-profile:a"`
	VideoProfile          *string           `flag:"-profile:v"`
	Target                *string           `flag:"-target"`
	Duration              *string           `flag:"-t"`
	Qscale                *uint32           `flag:"-qscale"`
	Crf                   *uint32           `flag:"-crf"`
	Strict                *int              `flag:"-strict"`
	MuxDelay              *string           `flag:"-muxdelay"`
	SeekTime              *string           `flag:"-ss"`
	SeekUsingTimestamp    *bool             `flag:"-seek_timestamp"`
	MovFlags              *string           `flag:"-movflags"`
	HideBanner            *bool             `flag:"-hide_banner"`
	OutputFormat          *string           `flag:"-f"`
	CopyTs                *bool             `flag:"-copyts"`
	NativeFramerateInput  *bool             `flag:"-re"`
	InputInitialOffset    *string           `flag:"-itsoffset"`
	RtmpLive              *string           `flag:"-rtmp_live"`
	HlsPlaylistType       *string           `flag:"-hls_playlist_type"`
	HlsListSize           *int              `flag:"-hls_list_size"`
	HlsSegmentDuration    *int              `flag:"-hls_time"`
	HlsMasterPlaylistName *string           `flag:"-master_pl_name"`
	HlsSegmentFilename    *string           `flag:"-hls_segment_filename"`
	HTTPMethod            *string           `flag:"-method"`
	HTTPKeepAlive         *bool             `flag:"-multiple_requests"`
	Hwaccel               *string           `flag:"-hwaccel"`
	StreamIds             map[string]string `flag:"-streamid"`
	VideoFilter           *string           `flag:"-vf"`
	AudioFilter           *string           `flag:"-af"`
	SkipVideo             *bool             `flag:"-vn"`
	SkipAudio             *bool             `flag:"-an"`
	CompressionLevel      *int              `flag:"-compression_level"`
	MapMetadata           *string           `flag:"-map_metadata"`
	Metadata              map[string]string `flag:"-metadata"`
	EncryptionKey         *string           `flag:"-hls_key_info_file"`
	Bframe                *int              `flag:"-bf"`
	PixFmt                *string           `flag:"-pix_fmt"`
	WhiteListProtocols    []string          `flag:"-protocol_whitelist"`
	Overwrite             *bool             `flag:"-y"`
}

type Process struct {
	output    string
	extraArgs map[string]string
	options   ProcessOptions
}

func (p *Process) WithOptions(opts ProcessOptions) IProcess {
	p.options = opts
	return p
}

func (p *Process) WithAdditionalOptions(opts map[string]string) IProcess {
	p.extraArgs = map[string]string{}
	for k, v := range opts {
		p.extraArgs[k] = v
	}
	return p
}

func (p *Process) Input(i string) IProcess {
	p.extraArgs["-i"] = i
	return p
}

func (p *Process) Output(o string) IProcess {
	p.output = o
	return p
}

func (p *Process) GetStrArguments() []string {
	f := reflect.TypeOf(p.options)
	v := reflect.ValueOf(p.options)

	values := []string{}

	for i := 0; i < f.NumField(); i++ {
		flag := f.Field(i).Tag.Get("flag")
		value := v.Field(i).Interface()

		if !v.Field(i).IsNil() {
			if vs, ok := value.(*bool); ok && *vs {
				values = append(values, flag)
			}

			if vs, ok := value.(*string); ok {
				values = append(values, flag, *vs)
			}

			if va, ok := value.([]string); ok {
				for _, item := range va {
					values = append(values, flag, item)
				}
			}

			if vm, ok := value.(map[string]string); ok {
				for k, v := range vm {
					values = append(values, flag, fmt.Sprintf("%v:%v", k, v))
				}
			}

			if vi, ok := value.(*int); ok {
				values = append(values, flag, fmt.Sprintf("%d", *vi))
			}
		}
	}

	for k, v := range p.extraArgs {
		values = append(values, k, v)
	}

	values = append(values, p.output)

	return values
}

type Options struct {
	FFmpegBinPath  string
	FFprobeBinPath string
}

type Transcoder struct {
	ffmpegBinPath  string
	ffprobeBinPath string
}

func (t *Transcoder) Process() IProcess {
	return &Process{
		options:   ProcessOptions{},
		extraArgs: map[string]string{},
	}
}

func (t *Transcoder) Run(p IProcess) error {
	return exec.Command(t.ffmpegBinPath, p.GetStrArguments()...).Run()
}

func NewTranscoder(opts Options) ITranscoder {
	return &Transcoder{
		ffmpegBinPath:  opts.FFmpegBinPath,
		ffprobeBinPath: opts.FFprobeBinPath,
	}
}
