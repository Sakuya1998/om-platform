// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.27.3
// source: pkg/utils/pagination/v1/pagination.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 分页请求
// 用于列表查询操作的通用分页请求结构
type PagingRequest struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	PageNum        uint32                 `protobuf:"varint,1,opt,name=page_num,json=pageNum,proto3" json:"page_num,omitempty"`                                                                                       // 页码，从1开始
	PageSize       uint32                 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`                                                                                    // 每页大小，默认为20
	SortBy         string                 `protobuf:"bytes,3,opt,name=sort_by,json=sortBy,proto3" json:"sort_by,omitempty"`                                                                                           // 排序字段，格式：field_name:asc/desc，例如：created_at:desc
	FilterBy       []string               `protobuf:"bytes,4,rep,name=filter_by,json=filterBy,proto3" json:"filter_by,omitempty"`                                                                                     // 过滤条件，格式：field_name:operator:value，例如：status:eq:active
	CountTotal     bool                   `protobuf:"varint,5,opt,name=count_total,json=countTotal,proto3" json:"count_total,omitempty"`                                                                              // 是否返回总记录数，默认为true
	TimeRangeStart *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=time_range_start,json=timeRangeStart,proto3" json:"time_range_start,omitempty"`                                                                 // 时间范围开始
	TimeRangeEnd   *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=time_range_end,json=timeRangeEnd,proto3" json:"time_range_end,omitempty"`                                                                       // 时间范围结束
	TimeField      string                 `protobuf:"bytes,8,opt,name=time_field,json=timeField,proto3" json:"time_field,omitempty"`                                                                                  // 时间字段名称，默认为created_at
	SearchKeyword  string                 `protobuf:"bytes,9,opt,name=search_keyword,json=searchKeyword,proto3" json:"search_keyword,omitempty"`                                                                      // 搜索关键字
	ExtraParams    map[string]string      `protobuf:"bytes,10,rep,name=extra_params,json=extraParams,proto3" json:"extra_params,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"` // 额外参数
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *PagingRequest) Reset() {
	*x = PagingRequest{}
	mi := &file_pkg_utils_pagination_v1_pagination_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PagingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PagingRequest) ProtoMessage() {}

func (x *PagingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_utils_pagination_v1_pagination_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PagingRequest.ProtoReflect.Descriptor instead.
func (*PagingRequest) Descriptor() ([]byte, []int) {
	return file_pkg_utils_pagination_v1_pagination_proto_rawDescGZIP(), []int{0}
}

func (x *PagingRequest) GetPageNum() uint32 {
	if x != nil {
		return x.PageNum
	}
	return 0
}

func (x *PagingRequest) GetPageSize() uint32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *PagingRequest) GetSortBy() string {
	if x != nil {
		return x.SortBy
	}
	return ""
}

func (x *PagingRequest) GetFilterBy() []string {
	if x != nil {
		return x.FilterBy
	}
	return nil
}

func (x *PagingRequest) GetCountTotal() bool {
	if x != nil {
		return x.CountTotal
	}
	return false
}

func (x *PagingRequest) GetTimeRangeStart() *timestamppb.Timestamp {
	if x != nil {
		return x.TimeRangeStart
	}
	return nil
}

func (x *PagingRequest) GetTimeRangeEnd() *timestamppb.Timestamp {
	if x != nil {
		return x.TimeRangeEnd
	}
	return nil
}

func (x *PagingRequest) GetTimeField() string {
	if x != nil {
		return x.TimeField
	}
	return ""
}

func (x *PagingRequest) GetSearchKeyword() string {
	if x != nil {
		return x.SearchKeyword
	}
	return ""
}

func (x *PagingRequest) GetExtraParams() map[string]string {
	if x != nil {
		return x.ExtraParams
	}
	return nil
}

// 分页响应
// 用于列表查询操作的通用分页响应结构
type PagingResponse struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	PageNum        uint32                 `protobuf:"varint,1,opt,name=page_num,json=pageNum,proto3" json:"page_num,omitempty"`                                                                                       // 当前页码
	PageSize       uint32                 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`                                                                                    // 每页大小
	Total          uint64                 `protobuf:"varint,3,opt,name=total,proto3" json:"total,omitempty"`                                                                                                          // 总记录数
	TotalPages     uint32                 `protobuf:"varint,4,opt,name=total_pages,json=totalPages,proto3" json:"total_pages,omitempty"`                                                                              // 总页数
	HasNext        bool                   `protobuf:"varint,5,opt,name=has_next,json=hasNext,proto3" json:"has_next,omitempty"`                                                                                       // 是否有下一页
	HasPrev        bool                   `protobuf:"varint,6,opt,name=has_prev,json=hasPrev,proto3" json:"has_prev,omitempty"`                                                                                       // 是否有上一页
	SortBy         string                 `protobuf:"bytes,7,opt,name=sort_by,json=sortBy,proto3" json:"sort_by,omitempty"`                                                                                           // 排序字段
	FilterBy       []string               `protobuf:"bytes,8,rep,name=filter_by,json=filterBy,proto3" json:"filter_by,omitempty"`                                                                                     // 过滤条件
	TimeRangeStart *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=time_range_start,json=timeRangeStart,proto3" json:"time_range_start,omitempty"`                                                                 // 时间范围开始
	TimeRangeEnd   *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=time_range_end,json=timeRangeEnd,proto3" json:"time_range_end,omitempty"`                                                                      // 时间范围结束
	TimeField      string                 `protobuf:"bytes,11,opt,name=time_field,json=timeField,proto3" json:"time_field,omitempty"`                                                                                 // 时间字段名称
	SearchKeyword  string                 `protobuf:"bytes,12,opt,name=search_keyword,json=searchKeyword,proto3" json:"search_keyword,omitempty"`                                                                     // 搜索关键字
	ExtraParams    map[string]string      `protobuf:"bytes,13,rep,name=extra_params,json=extraParams,proto3" json:"extra_params,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"` // 额外参数
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *PagingResponse) Reset() {
	*x = PagingResponse{}
	mi := &file_pkg_utils_pagination_v1_pagination_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PagingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PagingResponse) ProtoMessage() {}

func (x *PagingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_utils_pagination_v1_pagination_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PagingResponse.ProtoReflect.Descriptor instead.
func (*PagingResponse) Descriptor() ([]byte, []int) {
	return file_pkg_utils_pagination_v1_pagination_proto_rawDescGZIP(), []int{1}
}

func (x *PagingResponse) GetPageNum() uint32 {
	if x != nil {
		return x.PageNum
	}
	return 0
}

func (x *PagingResponse) GetPageSize() uint32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *PagingResponse) GetTotal() uint64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *PagingResponse) GetTotalPages() uint32 {
	if x != nil {
		return x.TotalPages
	}
	return 0
}

func (x *PagingResponse) GetHasNext() bool {
	if x != nil {
		return x.HasNext
	}
	return false
}

func (x *PagingResponse) GetHasPrev() bool {
	if x != nil {
		return x.HasPrev
	}
	return false
}

func (x *PagingResponse) GetSortBy() string {
	if x != nil {
		return x.SortBy
	}
	return ""
}

func (x *PagingResponse) GetFilterBy() []string {
	if x != nil {
		return x.FilterBy
	}
	return nil
}

func (x *PagingResponse) GetTimeRangeStart() *timestamppb.Timestamp {
	if x != nil {
		return x.TimeRangeStart
	}
	return nil
}

func (x *PagingResponse) GetTimeRangeEnd() *timestamppb.Timestamp {
	if x != nil {
		return x.TimeRangeEnd
	}
	return nil
}

func (x *PagingResponse) GetTimeField() string {
	if x != nil {
		return x.TimeField
	}
	return ""
}

func (x *PagingResponse) GetSearchKeyword() string {
	if x != nil {
		return x.SearchKeyword
	}
	return ""
}

func (x *PagingResponse) GetExtraParams() map[string]string {
	if x != nil {
		return x.ExtraParams
	}
	return nil
}

var File_pkg_utils_pagination_v1_pagination_proto protoreflect.FileDescriptor

const file_pkg_utils_pagination_v1_pagination_proto_rawDesc = "" +
	"\n" +
	"(pkg/utils/pagination/v1/pagination.proto\x12\x17pkg.utils.pagination.v1\x1a\x1fgoogle/protobuf/timestamp.proto\"\x88\x04\n" +
	"\rPagingRequest\x12\x19\n" +
	"\bpage_num\x18\x01 \x01(\rR\apageNum\x12\x1b\n" +
	"\tpage_size\x18\x02 \x01(\rR\bpageSize\x12\x17\n" +
	"\asort_by\x18\x03 \x01(\tR\x06sortBy\x12\x1b\n" +
	"\tfilter_by\x18\x04 \x03(\tR\bfilterBy\x12\x1f\n" +
	"\vcount_total\x18\x05 \x01(\bR\n" +
	"countTotal\x12D\n" +
	"\x10time_range_start\x18\x06 \x01(\v2\x1a.google.protobuf.TimestampR\x0etimeRangeStart\x12@\n" +
	"\x0etime_range_end\x18\a \x01(\v2\x1a.google.protobuf.TimestampR\ftimeRangeEnd\x12\x1d\n" +
	"\n" +
	"time_field\x18\b \x01(\tR\ttimeField\x12%\n" +
	"\x0esearch_keyword\x18\t \x01(\tR\rsearchKeyword\x12Z\n" +
	"\fextra_params\x18\n" +
	" \x03(\v27.pkg.utils.pagination.v1.PagingRequest.ExtraParamsEntryR\vextraParams\x1a>\n" +
	"\x10ExtraParamsEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01\"\xd6\x04\n" +
	"\x0ePagingResponse\x12\x19\n" +
	"\bpage_num\x18\x01 \x01(\rR\apageNum\x12\x1b\n" +
	"\tpage_size\x18\x02 \x01(\rR\bpageSize\x12\x14\n" +
	"\x05total\x18\x03 \x01(\x04R\x05total\x12\x1f\n" +
	"\vtotal_pages\x18\x04 \x01(\rR\n" +
	"totalPages\x12\x19\n" +
	"\bhas_next\x18\x05 \x01(\bR\ahasNext\x12\x19\n" +
	"\bhas_prev\x18\x06 \x01(\bR\ahasPrev\x12\x17\n" +
	"\asort_by\x18\a \x01(\tR\x06sortBy\x12\x1b\n" +
	"\tfilter_by\x18\b \x03(\tR\bfilterBy\x12D\n" +
	"\x10time_range_start\x18\t \x01(\v2\x1a.google.protobuf.TimestampR\x0etimeRangeStart\x12@\n" +
	"\x0etime_range_end\x18\n" +
	" \x01(\v2\x1a.google.protobuf.TimestampR\ftimeRangeEnd\x12\x1d\n" +
	"\n" +
	"time_field\x18\v \x01(\tR\ttimeField\x12%\n" +
	"\x0esearch_keyword\x18\f \x01(\tR\rsearchKeyword\x12[\n" +
	"\fextra_params\x18\r \x03(\v28.pkg.utils.pagination.v1.PagingResponse.ExtraParamsEntryR\vextraParams\x1a>\n" +
	"\x10ExtraParamsEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01BC\n" +
	"\x17pkg.utils.pagination.v1P\x01Z&om-platform/pkg/utils/pagination/v1;v1b\x06proto3"

var (
	file_pkg_utils_pagination_v1_pagination_proto_rawDescOnce sync.Once
	file_pkg_utils_pagination_v1_pagination_proto_rawDescData []byte
)

func file_pkg_utils_pagination_v1_pagination_proto_rawDescGZIP() []byte {
	file_pkg_utils_pagination_v1_pagination_proto_rawDescOnce.Do(func() {
		file_pkg_utils_pagination_v1_pagination_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_pkg_utils_pagination_v1_pagination_proto_rawDesc), len(file_pkg_utils_pagination_v1_pagination_proto_rawDesc)))
	})
	return file_pkg_utils_pagination_v1_pagination_proto_rawDescData
}

var file_pkg_utils_pagination_v1_pagination_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pkg_utils_pagination_v1_pagination_proto_goTypes = []any{
	(*PagingRequest)(nil),         // 0: pkg.utils.pagination.v1.PagingRequest
	(*PagingResponse)(nil),        // 1: pkg.utils.pagination.v1.PagingResponse
	nil,                           // 2: pkg.utils.pagination.v1.PagingRequest.ExtraParamsEntry
	nil,                           // 3: pkg.utils.pagination.v1.PagingResponse.ExtraParamsEntry
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_pkg_utils_pagination_v1_pagination_proto_depIdxs = []int32{
	4, // 0: pkg.utils.pagination.v1.PagingRequest.time_range_start:type_name -> google.protobuf.Timestamp
	4, // 1: pkg.utils.pagination.v1.PagingRequest.time_range_end:type_name -> google.protobuf.Timestamp
	2, // 2: pkg.utils.pagination.v1.PagingRequest.extra_params:type_name -> pkg.utils.pagination.v1.PagingRequest.ExtraParamsEntry
	4, // 3: pkg.utils.pagination.v1.PagingResponse.time_range_start:type_name -> google.protobuf.Timestamp
	4, // 4: pkg.utils.pagination.v1.PagingResponse.time_range_end:type_name -> google.protobuf.Timestamp
	3, // 5: pkg.utils.pagination.v1.PagingResponse.extra_params:type_name -> pkg.utils.pagination.v1.PagingResponse.ExtraParamsEntry
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_pkg_utils_pagination_v1_pagination_proto_init() }
func file_pkg_utils_pagination_v1_pagination_proto_init() {
	if File_pkg_utils_pagination_v1_pagination_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_pkg_utils_pagination_v1_pagination_proto_rawDesc), len(file_pkg_utils_pagination_v1_pagination_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_utils_pagination_v1_pagination_proto_goTypes,
		DependencyIndexes: file_pkg_utils_pagination_v1_pagination_proto_depIdxs,
		MessageInfos:      file_pkg_utils_pagination_v1_pagination_proto_msgTypes,
	}.Build()
	File_pkg_utils_pagination_v1_pagination_proto = out.File
	file_pkg_utils_pagination_v1_pagination_proto_goTypes = nil
	file_pkg_utils_pagination_v1_pagination_proto_depIdxs = nil
}
