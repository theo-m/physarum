// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.11.0
// source: physarum.proto

package pb

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Config_InitDistribution int32

const (
	Config_UNK       Config_InitDistribution = 0
	Config_UNIFORM   Config_InitDistribution = 1
	Config_CENTROIDS Config_InitDistribution = 2
	Config_CENTRE    Config_InitDistribution = 3
	Config_GRID      Config_InitDistribution = 4
)

// Enum value maps for Config_InitDistribution.
var (
	Config_InitDistribution_name = map[int32]string{
		0: "UNK",
		1: "UNIFORM",
		2: "CENTROIDS",
		3: "CENTRE",
		4: "GRID",
	}
	Config_InitDistribution_value = map[string]int32{
		"UNK":       0,
		"UNIFORM":   1,
		"CENTROIDS": 2,
		"CENTRE":    3,
		"GRID":      4,
	}
)

func (x Config_InitDistribution) Enum() *Config_InitDistribution {
	p := new(Config_InitDistribution)
	*p = x
	return p
}

func (x Config_InitDistribution) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Config_InitDistribution) Descriptor() protoreflect.EnumDescriptor {
	return file_physarum_proto_enumTypes[0].Descriptor()
}

func (Config_InitDistribution) Type() protoreflect.EnumType {
	return &file_physarum_proto_enumTypes[0]
}

func (x Config_InitDistribution) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Config_InitDistribution.Descriptor instead.
func (Config_InitDistribution) EnumDescriptor() ([]byte, []int) {
	return file_physarum_proto_rawDescGZIP(), []int{0, 0}
}

type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Width             int32                   `protobuf:"varint,1,opt,name=width,proto3" json:"width,omitempty"`
	Height            int32                   `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
	Particles         int64                   `protobuf:"varint,3,opt,name=particles,proto3" json:"particles,omitempty"`
	Iterations        int32                   `protobuf:"varint,4,opt,name=iterations,proto3" json:"iterations,omitempty"`
	BlurRadius        int32                   `protobuf:"varint,5,opt,name=blur_radius,json=blurRadius,proto3" json:"blur_radius,omitempty"`
	BlurPasses        int32                   `protobuf:"varint,6,opt,name=blur_passes,json=blurPasses,proto3" json:"blur_passes,omitempty"`
	ZoomFactor        float32                 `protobuf:"fixed32,7,opt,name=zoom_factor,json=zoomFactor,proto3" json:"zoom_factor,omitempty"`
	Agents            []*AgentConfig          `protobuf:"bytes,9,rep,name=agents,proto3" json:"agents,omitempty"`
	InteractionMatrix []float32               `protobuf:"fixed32,10,rep,packed,name=interaction_matrix,json=interactionMatrix,proto3" json:"interaction_matrix,omitempty"`
	Idist             Config_InitDistribution `protobuf:"varint,11,opt,name=idist,proto3,enum=physarium.Config_InitDistribution" json:"idist,omitempty"`
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_physarum_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_physarum_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_physarum_proto_rawDescGZIP(), []int{0}
}

func (x *Config) GetWidth() int32 {
	if x != nil {
		return x.Width
	}
	return 0
}

func (x *Config) GetHeight() int32 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *Config) GetParticles() int64 {
	if x != nil {
		return x.Particles
	}
	return 0
}

func (x *Config) GetIterations() int32 {
	if x != nil {
		return x.Iterations
	}
	return 0
}

func (x *Config) GetBlurRadius() int32 {
	if x != nil {
		return x.BlurRadius
	}
	return 0
}

func (x *Config) GetBlurPasses() int32 {
	if x != nil {
		return x.BlurPasses
	}
	return 0
}

func (x *Config) GetZoomFactor() float32 {
	if x != nil {
		return x.ZoomFactor
	}
	return 0
}

func (x *Config) GetAgents() []*AgentConfig {
	if x != nil {
		return x.Agents
	}
	return nil
}

func (x *Config) GetInteractionMatrix() []float32 {
	if x != nil {
		return x.InteractionMatrix
	}
	return nil
}

func (x *Config) GetIdist() Config_InitDistribution {
	if x != nil {
		return x.Idist
	}
	return Config_UNK
}

type AgentConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SensorAngle      float32 `protobuf:"fixed32,1,opt,name=sensor_angle,json=sensorAngle,proto3" json:"sensor_angle,omitempty"`
	SensorDistance   float32 `protobuf:"fixed32,2,opt,name=sensor_distance,json=sensorDistance,proto3" json:"sensor_distance,omitempty"`
	RotationAngle    float32 `protobuf:"fixed32,3,opt,name=rotation_angle,json=rotationAngle,proto3" json:"rotation_angle,omitempty"`
	StepDistance     float32 `protobuf:"fixed32,4,opt,name=step_distance,json=stepDistance,proto3" json:"step_distance,omitempty"`
	DepositionAmount float32 `protobuf:"fixed32,5,opt,name=deposition_amount,json=depositionAmount,proto3" json:"deposition_amount,omitempty"`
	DecayFactor      float32 `protobuf:"fixed32,6,opt,name=decay_factor,json=decayFactor,proto3" json:"decay_factor,omitempty"`
	Color            string  `protobuf:"bytes,7,opt,name=color,proto3" json:"color,omitempty"`
}

func (x *AgentConfig) Reset() {
	*x = AgentConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_physarum_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgentConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentConfig) ProtoMessage() {}

func (x *AgentConfig) ProtoReflect() protoreflect.Message {
	mi := &file_physarum_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgentConfig.ProtoReflect.Descriptor instead.
func (*AgentConfig) Descriptor() ([]byte, []int) {
	return file_physarum_proto_rawDescGZIP(), []int{1}
}

func (x *AgentConfig) GetSensorAngle() float32 {
	if x != nil {
		return x.SensorAngle
	}
	return 0
}

func (x *AgentConfig) GetSensorDistance() float32 {
	if x != nil {
		return x.SensorDistance
	}
	return 0
}

func (x *AgentConfig) GetRotationAngle() float32 {
	if x != nil {
		return x.RotationAngle
	}
	return 0
}

func (x *AgentConfig) GetStepDistance() float32 {
	if x != nil {
		return x.StepDistance
	}
	return 0
}

func (x *AgentConfig) GetDepositionAmount() float32 {
	if x != nil {
		return x.DepositionAmount
	}
	return 0
}

func (x *AgentConfig) GetDecayFactor() float32 {
	if x != nil {
		return x.DecayFactor
	}
	return 0
}

func (x *AgentConfig) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Content:
	//	*Event_Picture
	//	*Event_Video
	//	*Event_Step
	Content isEvent_Content `protobuf_oneof:"content"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_physarum_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_physarum_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_physarum_proto_rawDescGZIP(), []int{2}
}

func (m *Event) GetContent() isEvent_Content {
	if m != nil {
		return m.Content
	}
	return nil
}

func (x *Event) GetPicture() []byte {
	if x, ok := x.GetContent().(*Event_Picture); ok {
		return x.Picture
	}
	return nil
}

func (x *Event) GetVideo() []byte {
	if x, ok := x.GetContent().(*Event_Video); ok {
		return x.Video
	}
	return nil
}

func (x *Event) GetStep() string {
	if x, ok := x.GetContent().(*Event_Step); ok {
		return x.Step
	}
	return ""
}

type isEvent_Content interface {
	isEvent_Content()
}

type Event_Picture struct {
	Picture []byte `protobuf:"bytes,1,opt,name=picture,proto3,oneof"`
}

type Event_Video struct {
	Video []byte `protobuf:"bytes,2,opt,name=video,proto3,oneof"`
}

type Event_Step struct {
	Step string `protobuf:"bytes,3,opt,name=step,proto3,oneof"`
}

func (*Event_Picture) isEvent_Content() {}

func (*Event_Video) isEvent_Content() {}

func (*Event_Step) isEvent_Content() {}

var File_physarum_proto protoreflect.FileDescriptor

var file_physarum_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x70, 0x68, 0x79, 0x73, 0x61, 0x72, 0x75, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x09, 0x70, 0x68, 0x79, 0x73, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x22, 0xbf, 0x03, 0x0a, 0x06,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68, 0x12, 0x16, 0x0a, 0x06,
	0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x68, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c,
	0x65, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x74, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x69, 0x74, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x62, 0x6c, 0x75, 0x72, 0x5f, 0x72, 0x61, 0x64, 0x69, 0x75,
	0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x62, 0x6c, 0x75, 0x72, 0x52, 0x61, 0x64,
	0x69, 0x75, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x62, 0x6c, 0x75, 0x72, 0x5f, 0x70, 0x61, 0x73, 0x73,
	0x65, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x62, 0x6c, 0x75, 0x72, 0x50, 0x61,
	0x73, 0x73, 0x65, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x7a, 0x6f, 0x6f, 0x6d, 0x5f, 0x66, 0x61, 0x63,
	0x74, 0x6f, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0a, 0x7a, 0x6f, 0x6f, 0x6d, 0x46,
	0x61, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x2e, 0x0a, 0x06, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x18,
	0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x68, 0x79, 0x73, 0x61, 0x72, 0x69, 0x75,
	0x6d, 0x2e, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x06, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x2d, 0x0a, 0x12, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x61, 0x74, 0x72, 0x69, 0x78, 0x18, 0x0a, 0x20, 0x03, 0x28,
	0x02, 0x52, 0x11, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x61,
	0x74, 0x72, 0x69, 0x78, 0x12, 0x38, 0x0a, 0x05, 0x69, 0x64, 0x69, 0x73, 0x74, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x70, 0x68, 0x79, 0x73, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2e,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x44, 0x69, 0x73, 0x74, 0x72,
	0x69, 0x62, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x69, 0x64, 0x69, 0x73, 0x74, 0x22, 0x4d,
	0x0a, 0x10, 0x49, 0x6e, 0x69, 0x74, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x07, 0x0a, 0x03, 0x55, 0x4e, 0x4b, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x55,
	0x4e, 0x49, 0x46, 0x4f, 0x52, 0x4d, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x43, 0x45, 0x4e, 0x54,
	0x52, 0x4f, 0x49, 0x44, 0x53, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x45, 0x4e, 0x54, 0x52,
	0x45, 0x10, 0x03, 0x12, 0x08, 0x0a, 0x04, 0x47, 0x52, 0x49, 0x44, 0x10, 0x04, 0x22, 0x8b, 0x02,
	0x0a, 0x0b, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x21, 0x0a,
	0x0c, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x5f, 0x61, 0x6e, 0x67, 0x6c, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x02, 0x52, 0x0b, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x41, 0x6e, 0x67, 0x6c, 0x65,
	0x12, 0x27, 0x0a, 0x0f, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x5f, 0x64, 0x69, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0e, 0x73, 0x65, 0x6e, 0x73, 0x6f,
	0x72, 0x44, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x72, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x61, 0x6e, 0x67, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x02, 0x52, 0x0d, 0x72, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x6e, 0x67, 0x6c, 0x65,
	0x12, 0x23, 0x0a, 0x0d, 0x73, 0x74, 0x65, 0x70, 0x5f, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0c, 0x73, 0x74, 0x65, 0x70, 0x44, 0x69, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x2b, 0x0a, 0x11, 0x64, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x02,
	0x52, 0x10, 0x64, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x6d, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x65, 0x63, 0x61, 0x79, 0x5f, 0x66, 0x61, 0x63, 0x74,
	0x6f, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0b, 0x64, 0x65, 0x63, 0x61, 0x79, 0x46,
	0x61, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x22, 0x5c, 0x0a, 0x05, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x07, 0x70, 0x69, 0x63, 0x74, 0x75, 0x72, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x07, 0x70, 0x69, 0x63, 0x74, 0x75, 0x72, 0x65,
	0x12, 0x16, 0x0a, 0x05, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x48,
	0x00, 0x52, 0x05, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x12, 0x14, 0x0a, 0x04, 0x73, 0x74, 0x65, 0x70,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x73, 0x74, 0x65, 0x70, 0x42, 0x09,
	0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x42, 0x20, 0x5a, 0x1e, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x68, 0x65, 0x6f, 0x2d, 0x6d, 0x2f, 0x70,
	0x68, 0x79, 0x73, 0x61, 0x72, 0x69, 0x75, 0x6d, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_physarum_proto_rawDescOnce sync.Once
	file_physarum_proto_rawDescData = file_physarum_proto_rawDesc
)

func file_physarum_proto_rawDescGZIP() []byte {
	file_physarum_proto_rawDescOnce.Do(func() {
		file_physarum_proto_rawDescData = protoimpl.X.CompressGZIP(file_physarum_proto_rawDescData)
	})
	return file_physarum_proto_rawDescData
}

var file_physarum_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_physarum_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_physarum_proto_goTypes = []interface{}{
	(Config_InitDistribution)(0), // 0: physarium.Config.InitDistribution
	(*Config)(nil),               // 1: physarium.Config
	(*AgentConfig)(nil),          // 2: physarium.AgentConfig
	(*Event)(nil),                // 3: physarium.Event
}
var file_physarum_proto_depIdxs = []int32{
	2, // 0: physarium.Config.agents:type_name -> physarium.AgentConfig
	0, // 1: physarium.Config.idist:type_name -> physarium.Config.InitDistribution
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_physarum_proto_init() }
func file_physarum_proto_init() {
	if File_physarum_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_physarum_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Config); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_physarum_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AgentConfig); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_physarum_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_physarum_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Event_Picture)(nil),
		(*Event_Video)(nil),
		(*Event_Step)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_physarum_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_physarum_proto_goTypes,
		DependencyIndexes: file_physarum_proto_depIdxs,
		EnumInfos:         file_physarum_proto_enumTypes,
		MessageInfos:      file_physarum_proto_msgTypes,
	}.Build()
	File_physarum_proto = out.File
	file_physarum_proto_rawDesc = nil
	file_physarum_proto_goTypes = nil
	file_physarum_proto_depIdxs = nil
}
