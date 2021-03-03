// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dreamvo/gilfoyle/ent/media"
	"github.com/dreamvo/gilfoyle/ent/mediafile"
	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/facebook/ent/schema/field"
	"github.com/google/uuid"
)

// MediaFileCreate is the builder for creating a MediaFile entity.
type MediaFileCreate struct {
	config
	mutation *MediaFileMutation
	hooks    []Hook
}

// SetRenditionName sets the rendition_name field.
func (mfc *MediaFileCreate) SetRenditionName(s string) *MediaFileCreate {
	mfc.mutation.SetRenditionName(s)
	return mfc
}

// SetFormat sets the format field.
func (mfc *MediaFileCreate) SetFormat(s string) *MediaFileCreate {
	mfc.mutation.SetFormat(s)
	return mfc
}

// SetTargetBandwidth sets the target_bandwidth field.
func (mfc *MediaFileCreate) SetTargetBandwidth(u uint64) *MediaFileCreate {
	mfc.mutation.SetTargetBandwidth(u)
	return mfc
}

// SetNillableTargetBandwidth sets the target_bandwidth field if the given value is not nil.
func (mfc *MediaFileCreate) SetNillableTargetBandwidth(u *uint64) *MediaFileCreate {
	if u != nil {
		mfc.SetTargetBandwidth(*u)
	}
	return mfc
}

// SetVideoBitrate sets the video_bitrate field.
func (mfc *MediaFileCreate) SetVideoBitrate(i int64) *MediaFileCreate {
	mfc.mutation.SetVideoBitrate(i)
	return mfc
}

// SetAudioBitrate sets the audio_bitrate field.
func (mfc *MediaFileCreate) SetAudioBitrate(i int64) *MediaFileCreate {
	mfc.mutation.SetAudioBitrate(i)
	return mfc
}

// SetVideoCodec sets the video_codec field.
func (mfc *MediaFileCreate) SetVideoCodec(s string) *MediaFileCreate {
	mfc.mutation.SetVideoCodec(s)
	return mfc
}

// SetNillableVideoCodec sets the video_codec field if the given value is not nil.
func (mfc *MediaFileCreate) SetNillableVideoCodec(s *string) *MediaFileCreate {
	if s != nil {
		mfc.SetVideoCodec(*s)
	}
	return mfc
}

// SetAudioCodec sets the audio_codec field.
func (mfc *MediaFileCreate) SetAudioCodec(s string) *MediaFileCreate {
	mfc.mutation.SetAudioCodec(s)
	return mfc
}

// SetNillableAudioCodec sets the audio_codec field if the given value is not nil.
func (mfc *MediaFileCreate) SetNillableAudioCodec(s *string) *MediaFileCreate {
	if s != nil {
		mfc.SetAudioCodec(*s)
	}
	return mfc
}

// SetResolutionWidth sets the resolution_width field.
func (mfc *MediaFileCreate) SetResolutionWidth(u uint16) *MediaFileCreate {
	mfc.mutation.SetResolutionWidth(u)
	return mfc
}

// SetResolutionHeight sets the resolution_height field.
func (mfc *MediaFileCreate) SetResolutionHeight(u uint16) *MediaFileCreate {
	mfc.mutation.SetResolutionHeight(u)
	return mfc
}

// SetFramerate sets the framerate field.
func (mfc *MediaFileCreate) SetFramerate(u uint8) *MediaFileCreate {
	mfc.mutation.SetFramerate(u)
	return mfc
}

// SetDurationSeconds sets the duration_seconds field.
func (mfc *MediaFileCreate) SetDurationSeconds(f float64) *MediaFileCreate {
	mfc.mutation.SetDurationSeconds(f)
	return mfc
}

// SetMediaType sets the media_type field.
func (mfc *MediaFileCreate) SetMediaType(mt mediafile.MediaType) *MediaFileCreate {
	mfc.mutation.SetMediaType(mt)
	return mfc
}

// SetStatus sets the status field.
func (mfc *MediaFileCreate) SetStatus(m mediafile.Status) *MediaFileCreate {
	mfc.mutation.SetStatus(m)
	return mfc
}

// SetMessage sets the message field.
func (mfc *MediaFileCreate) SetMessage(s string) *MediaFileCreate {
	mfc.mutation.SetMessage(s)
	return mfc
}

// SetNillableMessage sets the message field if the given value is not nil.
func (mfc *MediaFileCreate) SetNillableMessage(s *string) *MediaFileCreate {
	if s != nil {
		mfc.SetMessage(*s)
	}
	return mfc
}

// SetEntryFile sets the entry_file field.
func (mfc *MediaFileCreate) SetEntryFile(s string) *MediaFileCreate {
	mfc.mutation.SetEntryFile(s)
	return mfc
}

// SetNillableEntryFile sets the entry_file field if the given value is not nil.
func (mfc *MediaFileCreate) SetNillableEntryFile(s *string) *MediaFileCreate {
	if s != nil {
		mfc.SetEntryFile(*s)
	}
	return mfc
}

// SetMimetype sets the mimetype field.
func (mfc *MediaFileCreate) SetMimetype(s string) *MediaFileCreate {
	mfc.mutation.SetMimetype(s)
	return mfc
}

// SetNillableMimetype sets the mimetype field if the given value is not nil.
func (mfc *MediaFileCreate) SetNillableMimetype(s *string) *MediaFileCreate {
	if s != nil {
		mfc.SetMimetype(*s)
	}
	return mfc
}

// SetCreatedAt sets the created_at field.
func (mfc *MediaFileCreate) SetCreatedAt(t time.Time) *MediaFileCreate {
	mfc.mutation.SetCreatedAt(t)
	return mfc
}

// SetNillableCreatedAt sets the created_at field if the given value is not nil.
func (mfc *MediaFileCreate) SetNillableCreatedAt(t *time.Time) *MediaFileCreate {
	if t != nil {
		mfc.SetCreatedAt(*t)
	}
	return mfc
}

// SetUpdatedAt sets the updated_at field.
func (mfc *MediaFileCreate) SetUpdatedAt(t time.Time) *MediaFileCreate {
	mfc.mutation.SetUpdatedAt(t)
	return mfc
}

// SetNillableUpdatedAt sets the updated_at field if the given value is not nil.
func (mfc *MediaFileCreate) SetNillableUpdatedAt(t *time.Time) *MediaFileCreate {
	if t != nil {
		mfc.SetUpdatedAt(*t)
	}
	return mfc
}

// SetID sets the id field.
func (mfc *MediaFileCreate) SetID(u uuid.UUID) *MediaFileCreate {
	mfc.mutation.SetID(u)
	return mfc
}

// SetMediaID sets the media edge to Media by id.
func (mfc *MediaFileCreate) SetMediaID(id uuid.UUID) *MediaFileCreate {
	mfc.mutation.SetMediaID(id)
	return mfc
}

// SetMedia sets the media edge to Media.
func (mfc *MediaFileCreate) SetMedia(m *Media) *MediaFileCreate {
	return mfc.SetMediaID(m.ID)
}

// Mutation returns the MediaFileMutation object of the builder.
func (mfc *MediaFileCreate) Mutation() *MediaFileMutation {
	return mfc.mutation
}

// Save creates the MediaFile in the database.
func (mfc *MediaFileCreate) Save(ctx context.Context) (*MediaFile, error) {
	var (
		err  error
		node *MediaFile
	)
	mfc.defaults()
	if len(mfc.hooks) == 0 {
		if err = mfc.check(); err != nil {
			return nil, err
		}
		node, err = mfc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*MediaFileMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = mfc.check(); err != nil {
				return nil, err
			}
			mfc.mutation = mutation
			node, err = mfc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(mfc.hooks) - 1; i >= 0; i-- {
			mut = mfc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, mfc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (mfc *MediaFileCreate) SaveX(ctx context.Context) *MediaFile {
	v, err := mfc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// defaults sets the default values of the builder before save.
func (mfc *MediaFileCreate) defaults() {
	if _, ok := mfc.mutation.TargetBandwidth(); !ok {
		v := mediafile.DefaultTargetBandwidth
		mfc.mutation.SetTargetBandwidth(v)
	}
	if _, ok := mfc.mutation.VideoCodec(); !ok {
		v := mediafile.DefaultVideoCodec
		mfc.mutation.SetVideoCodec(v)
	}
	if _, ok := mfc.mutation.AudioCodec(); !ok {
		v := mediafile.DefaultAudioCodec
		mfc.mutation.SetAudioCodec(v)
	}
	if _, ok := mfc.mutation.Message(); !ok {
		v := mediafile.DefaultMessage
		mfc.mutation.SetMessage(v)
	}
	if _, ok := mfc.mutation.EntryFile(); !ok {
		v := mediafile.DefaultEntryFile
		mfc.mutation.SetEntryFile(v)
	}
	if _, ok := mfc.mutation.Mimetype(); !ok {
		v := mediafile.DefaultMimetype
		mfc.mutation.SetMimetype(v)
	}
	if _, ok := mfc.mutation.CreatedAt(); !ok {
		v := mediafile.DefaultCreatedAt()
		mfc.mutation.SetCreatedAt(v)
	}
	if _, ok := mfc.mutation.UpdatedAt(); !ok {
		v := mediafile.DefaultUpdatedAt()
		mfc.mutation.SetUpdatedAt(v)
	}
	if _, ok := mfc.mutation.ID(); !ok {
		v := mediafile.DefaultID()
		mfc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mfc *MediaFileCreate) check() error {
	if _, ok := mfc.mutation.RenditionName(); !ok {
		return &ValidationError{Name: "rendition_name", err: errors.New("ent: missing required field \"rendition_name\"")}
	}
	if v, ok := mfc.mutation.RenditionName(); ok {
		if err := mediafile.RenditionNameValidator(v); err != nil {
			return &ValidationError{Name: "rendition_name", err: fmt.Errorf("ent: validator failed for field \"rendition_name\": %w", err)}
		}
	}
	if _, ok := mfc.mutation.Format(); !ok {
		return &ValidationError{Name: "format", err: errors.New("ent: missing required field \"format\"")}
	}
	if v, ok := mfc.mutation.Format(); ok {
		if err := mediafile.FormatValidator(v); err != nil {
			return &ValidationError{Name: "format", err: fmt.Errorf("ent: validator failed for field \"format\": %w", err)}
		}
	}
	if _, ok := mfc.mutation.TargetBandwidth(); !ok {
		return &ValidationError{Name: "target_bandwidth", err: errors.New("ent: missing required field \"target_bandwidth\"")}
	}
	if _, ok := mfc.mutation.VideoBitrate(); !ok {
		return &ValidationError{Name: "video_bitrate", err: errors.New("ent: missing required field \"video_bitrate\"")}
	}
	if v, ok := mfc.mutation.VideoBitrate(); ok {
		if err := mediafile.VideoBitrateValidator(v); err != nil {
			return &ValidationError{Name: "video_bitrate", err: fmt.Errorf("ent: validator failed for field \"video_bitrate\": %w", err)}
		}
	}
	if _, ok := mfc.mutation.AudioBitrate(); !ok {
		return &ValidationError{Name: "audio_bitrate", err: errors.New("ent: missing required field \"audio_bitrate\"")}
	}
	if v, ok := mfc.mutation.AudioBitrate(); ok {
		if err := mediafile.AudioBitrateValidator(v); err != nil {
			return &ValidationError{Name: "audio_bitrate", err: fmt.Errorf("ent: validator failed for field \"audio_bitrate\": %w", err)}
		}
	}
	if _, ok := mfc.mutation.VideoCodec(); !ok {
		return &ValidationError{Name: "video_codec", err: errors.New("ent: missing required field \"video_codec\"")}
	}
	if v, ok := mfc.mutation.VideoCodec(); ok {
		if err := mediafile.VideoCodecValidator(v); err != nil {
			return &ValidationError{Name: "video_codec", err: fmt.Errorf("ent: validator failed for field \"video_codec\": %w", err)}
		}
	}
	if _, ok := mfc.mutation.AudioCodec(); !ok {
		return &ValidationError{Name: "audio_codec", err: errors.New("ent: missing required field \"audio_codec\"")}
	}
	if v, ok := mfc.mutation.AudioCodec(); ok {
		if err := mediafile.AudioCodecValidator(v); err != nil {
			return &ValidationError{Name: "audio_codec", err: fmt.Errorf("ent: validator failed for field \"audio_codec\": %w", err)}
		}
	}
	if _, ok := mfc.mutation.ResolutionWidth(); !ok {
		return &ValidationError{Name: "resolution_width", err: errors.New("ent: missing required field \"resolution_width\"")}
	}
	if v, ok := mfc.mutation.ResolutionWidth(); ok {
		if err := mediafile.ResolutionWidthValidator(v); err != nil {
			return &ValidationError{Name: "resolution_width", err: fmt.Errorf("ent: validator failed for field \"resolution_width\": %w", err)}
		}
	}
	if _, ok := mfc.mutation.ResolutionHeight(); !ok {
		return &ValidationError{Name: "resolution_height", err: errors.New("ent: missing required field \"resolution_height\"")}
	}
	if v, ok := mfc.mutation.ResolutionHeight(); ok {
		if err := mediafile.ResolutionHeightValidator(v); err != nil {
			return &ValidationError{Name: "resolution_height", err: fmt.Errorf("ent: validator failed for field \"resolution_height\": %w", err)}
		}
	}
	if _, ok := mfc.mutation.Framerate(); !ok {
		return &ValidationError{Name: "framerate", err: errors.New("ent: missing required field \"framerate\"")}
	}
	if v, ok := mfc.mutation.Framerate(); ok {
		if err := mediafile.FramerateValidator(v); err != nil {
			return &ValidationError{Name: "framerate", err: fmt.Errorf("ent: validator failed for field \"framerate\": %w", err)}
		}
	}
	if _, ok := mfc.mutation.DurationSeconds(); !ok {
		return &ValidationError{Name: "duration_seconds", err: errors.New("ent: missing required field \"duration_seconds\"")}
	}
	if v, ok := mfc.mutation.DurationSeconds(); ok {
		if err := mediafile.DurationSecondsValidator(v); err != nil {
			return &ValidationError{Name: "duration_seconds", err: fmt.Errorf("ent: validator failed for field \"duration_seconds\": %w", err)}
		}
	}
	if _, ok := mfc.mutation.MediaType(); !ok {
		return &ValidationError{Name: "media_type", err: errors.New("ent: missing required field \"media_type\"")}
	}
	if v, ok := mfc.mutation.MediaType(); ok {
		if err := mediafile.MediaTypeValidator(v); err != nil {
			return &ValidationError{Name: "media_type", err: fmt.Errorf("ent: validator failed for field \"media_type\": %w", err)}
		}
	}
	if _, ok := mfc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New("ent: missing required field \"status\"")}
	}
	if v, ok := mfc.mutation.Status(); ok {
		if err := mediafile.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf("ent: validator failed for field \"status\": %w", err)}
		}
	}
	if v, ok := mfc.mutation.Message(); ok {
		if err := mediafile.MessageValidator(v); err != nil {
			return &ValidationError{Name: "message", err: fmt.Errorf("ent: validator failed for field \"message\": %w", err)}
		}
	}
	if _, ok := mfc.mutation.EntryFile(); !ok {
		return &ValidationError{Name: "entry_file", err: errors.New("ent: missing required field \"entry_file\"")}
	}
	if v, ok := mfc.mutation.EntryFile(); ok {
		if err := mediafile.EntryFileValidator(v); err != nil {
			return &ValidationError{Name: "entry_file", err: fmt.Errorf("ent: validator failed for field \"entry_file\": %w", err)}
		}
	}
	if _, ok := mfc.mutation.Mimetype(); !ok {
		return &ValidationError{Name: "mimetype", err: errors.New("ent: missing required field \"mimetype\"")}
	}
	if v, ok := mfc.mutation.Mimetype(); ok {
		if err := mediafile.MimetypeValidator(v); err != nil {
			return &ValidationError{Name: "mimetype", err: fmt.Errorf("ent: validator failed for field \"mimetype\": %w", err)}
		}
	}
	if _, ok := mfc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New("ent: missing required field \"created_at\"")}
	}
	if _, ok := mfc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New("ent: missing required field \"updated_at\"")}
	}
	if _, ok := mfc.mutation.MediaID(); !ok {
		return &ValidationError{Name: "media", err: errors.New("ent: missing required edge \"media\"")}
	}
	return nil
}

func (mfc *MediaFileCreate) sqlSave(ctx context.Context) (*MediaFile, error) {
	_node, _spec := mfc.createSpec()
	if err := sqlgraph.CreateNode(ctx, mfc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return _node, nil
}

func (mfc *MediaFileCreate) createSpec() (*MediaFile, *sqlgraph.CreateSpec) {
	var (
		_node = &MediaFile{config: mfc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: mediafile.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: mediafile.FieldID,
			},
		}
	)
	if id, ok := mfc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := mfc.mutation.RenditionName(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: mediafile.FieldRenditionName,
		})
		_node.RenditionName = value
	}
	if value, ok := mfc.mutation.Format(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: mediafile.FieldFormat,
		})
		_node.Format = value
	}
	if value, ok := mfc.mutation.TargetBandwidth(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeUint64,
			Value:  value,
			Column: mediafile.FieldTargetBandwidth,
		})
		_node.TargetBandwidth = value
	}
	if value, ok := mfc.mutation.VideoBitrate(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: mediafile.FieldVideoBitrate,
		})
		_node.VideoBitrate = value
	}
	if value, ok := mfc.mutation.AudioBitrate(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: mediafile.FieldAudioBitrate,
		})
		_node.AudioBitrate = value
	}
	if value, ok := mfc.mutation.VideoCodec(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: mediafile.FieldVideoCodec,
		})
		_node.VideoCodec = value
	}
	if value, ok := mfc.mutation.AudioCodec(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: mediafile.FieldAudioCodec,
		})
		_node.AudioCodec = value
	}
	if value, ok := mfc.mutation.ResolutionWidth(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeUint16,
			Value:  value,
			Column: mediafile.FieldResolutionWidth,
		})
		_node.ResolutionWidth = value
	}
	if value, ok := mfc.mutation.ResolutionHeight(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeUint16,
			Value:  value,
			Column: mediafile.FieldResolutionHeight,
		})
		_node.ResolutionHeight = value
	}
	if value, ok := mfc.mutation.Framerate(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeUint8,
			Value:  value,
			Column: mediafile.FieldFramerate,
		})
		_node.Framerate = value
	}
	if value, ok := mfc.mutation.DurationSeconds(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: mediafile.FieldDurationSeconds,
		})
		_node.DurationSeconds = value
	}
	if value, ok := mfc.mutation.MediaType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: mediafile.FieldMediaType,
		})
		_node.MediaType = value
	}
	if value, ok := mfc.mutation.Status(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: mediafile.FieldStatus,
		})
		_node.Status = value
	}
	if value, ok := mfc.mutation.Message(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: mediafile.FieldMessage,
		})
		_node.Message = value
	}
	if value, ok := mfc.mutation.EntryFile(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: mediafile.FieldEntryFile,
		})
		_node.EntryFile = value
	}
	if value, ok := mfc.mutation.Mimetype(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: mediafile.FieldMimetype,
		})
		_node.Mimetype = value
	}
	if value, ok := mfc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: mediafile.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := mfc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: mediafile.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	if nodes := mfc.mutation.MediaIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   mediafile.MediaTable,
			Columns: []string{mediafile.MediaColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: media.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// MediaFileCreateBulk is the builder for creating a bulk of MediaFile entities.
type MediaFileCreateBulk struct {
	config
	builders []*MediaFileCreate
}

// Save creates the MediaFile entities in the database.
func (mfcb *MediaFileCreateBulk) Save(ctx context.Context) ([]*MediaFile, error) {
	specs := make([]*sqlgraph.CreateSpec, len(mfcb.builders))
	nodes := make([]*MediaFile, len(mfcb.builders))
	mutators := make([]Mutator, len(mfcb.builders))
	for i := range mfcb.builders {
		func(i int, root context.Context) {
			builder := mfcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*MediaFileMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, mfcb.builders[i+1].mutation)
				} else {
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, mfcb.driver, &sqlgraph.BatchCreateSpec{Nodes: specs}); err != nil {
						if cerr, ok := isSQLConstraintError(err); ok {
							err = cerr
						}
					}
				}
				mutation.done = true
				if err != nil {
					return nil, err
				}
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, mfcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX calls Save and panics if Save returns an error.
func (mfcb *MediaFileCreateBulk) SaveX(ctx context.Context) []*MediaFile {
	v, err := mfcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}
