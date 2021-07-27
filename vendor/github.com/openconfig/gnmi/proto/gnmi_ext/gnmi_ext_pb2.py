# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: proto/gnmi_ext/gnmi_ext.proto

from google.protobuf.internal import enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='proto/gnmi_ext/gnmi_ext.proto',
  package='gnmi_ext',
  syntax='proto3',
  serialized_options=b'Z)github.com/openconfig/gnmi/proto/gnmi_ext',
  serialized_pb=b'\n\x1dproto/gnmi_ext/gnmi_ext.proto\x12\x08gnmi_ext\"\x86\x01\n\tExtension\x12\x37\n\x0eregistered_ext\x18\x01 \x01(\x0b\x32\x1d.gnmi_ext.RegisteredExtensionH\x00\x12\x39\n\x12master_arbitration\x18\x02 \x01(\x0b\x32\x1b.gnmi_ext.MasterArbitrationH\x00\x42\x05\n\x03\x65xt\"E\n\x13RegisteredExtension\x12!\n\x02id\x18\x01 \x01(\x0e\x32\x15.gnmi_ext.ExtensionID\x12\x0b\n\x03msg\x18\x02 \x01(\x0c\"Y\n\x11MasterArbitration\x12\x1c\n\x04role\x18\x01 \x01(\x0b\x32\x0e.gnmi_ext.Role\x12&\n\x0b\x65lection_id\x18\x02 \x01(\x0b\x32\x11.gnmi_ext.Uint128\"$\n\x07Uint128\x12\x0c\n\x04high\x18\x01 \x01(\x04\x12\x0b\n\x03low\x18\x02 \x01(\x04\"\x12\n\x04Role\x12\n\n\x02id\x18\x01 \x01(\t*3\n\x0b\x45xtensionID\x12\r\n\tEID_UNSET\x10\x00\x12\x15\n\x10\x45ID_EXPERIMENTAL\x10\xe7\x07\x42+Z)github.com/openconfig/gnmi/proto/gnmi_extb\x06proto3'
)

_EXTENSIONID = _descriptor.EnumDescriptor(
  name='ExtensionID',
  full_name='gnmi_ext.ExtensionID',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='EID_UNSET', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='EID_EXPERIMENTAL', index=1, number=999,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=400,
  serialized_end=451,
)
_sym_db.RegisterEnumDescriptor(_EXTENSIONID)

ExtensionID = enum_type_wrapper.EnumTypeWrapper(_EXTENSIONID)
EID_UNSET = 0
EID_EXPERIMENTAL = 999



_EXTENSION = _descriptor.Descriptor(
  name='Extension',
  full_name='gnmi_ext.Extension',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='registered_ext', full_name='gnmi_ext.Extension.registered_ext', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='master_arbitration', full_name='gnmi_ext.Extension.master_arbitration', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
    _descriptor.OneofDescriptor(
      name='ext', full_name='gnmi_ext.Extension.ext',
      index=0, containing_type=None, fields=[]),
  ],
  serialized_start=44,
  serialized_end=178,
)


_REGISTEREDEXTENSION = _descriptor.Descriptor(
  name='RegisteredExtension',
  full_name='gnmi_ext.RegisteredExtension',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='gnmi_ext.RegisteredExtension.id', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='msg', full_name='gnmi_ext.RegisteredExtension.msg', index=1,
      number=2, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=b"",
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=180,
  serialized_end=249,
)


_MASTERARBITRATION = _descriptor.Descriptor(
  name='MasterArbitration',
  full_name='gnmi_ext.MasterArbitration',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='role', full_name='gnmi_ext.MasterArbitration.role', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='election_id', full_name='gnmi_ext.MasterArbitration.election_id', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=251,
  serialized_end=340,
)


_UINT128 = _descriptor.Descriptor(
  name='Uint128',
  full_name='gnmi_ext.Uint128',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='high', full_name='gnmi_ext.Uint128.high', index=0,
      number=1, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='low', full_name='gnmi_ext.Uint128.low', index=1,
      number=2, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=342,
  serialized_end=378,
)


_ROLE = _descriptor.Descriptor(
  name='Role',
  full_name='gnmi_ext.Role',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='gnmi_ext.Role.id', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=380,
  serialized_end=398,
)

_EXTENSION.fields_by_name['registered_ext'].message_type = _REGISTEREDEXTENSION
_EXTENSION.fields_by_name['master_arbitration'].message_type = _MASTERARBITRATION
_EXTENSION.oneofs_by_name['ext'].fields.append(
  _EXTENSION.fields_by_name['registered_ext'])
_EXTENSION.fields_by_name['registered_ext'].containing_oneof = _EXTENSION.oneofs_by_name['ext']
_EXTENSION.oneofs_by_name['ext'].fields.append(
  _EXTENSION.fields_by_name['master_arbitration'])
_EXTENSION.fields_by_name['master_arbitration'].containing_oneof = _EXTENSION.oneofs_by_name['ext']
_REGISTEREDEXTENSION.fields_by_name['id'].enum_type = _EXTENSIONID
_MASTERARBITRATION.fields_by_name['role'].message_type = _ROLE
_MASTERARBITRATION.fields_by_name['election_id'].message_type = _UINT128
DESCRIPTOR.message_types_by_name['Extension'] = _EXTENSION
DESCRIPTOR.message_types_by_name['RegisteredExtension'] = _REGISTEREDEXTENSION
DESCRIPTOR.message_types_by_name['MasterArbitration'] = _MASTERARBITRATION
DESCRIPTOR.message_types_by_name['Uint128'] = _UINT128
DESCRIPTOR.message_types_by_name['Role'] = _ROLE
DESCRIPTOR.enum_types_by_name['ExtensionID'] = _EXTENSIONID
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Extension = _reflection.GeneratedProtocolMessageType('Extension', (_message.Message,), {
  'DESCRIPTOR' : _EXTENSION,
  '__module__' : 'proto.gnmi_ext.gnmi_ext_pb2'
  # @@protoc_insertion_point(class_scope:gnmi_ext.Extension)
  })
_sym_db.RegisterMessage(Extension)

RegisteredExtension = _reflection.GeneratedProtocolMessageType('RegisteredExtension', (_message.Message,), {
  'DESCRIPTOR' : _REGISTEREDEXTENSION,
  '__module__' : 'proto.gnmi_ext.gnmi_ext_pb2'
  # @@protoc_insertion_point(class_scope:gnmi_ext.RegisteredExtension)
  })
_sym_db.RegisterMessage(RegisteredExtension)

MasterArbitration = _reflection.GeneratedProtocolMessageType('MasterArbitration', (_message.Message,), {
  'DESCRIPTOR' : _MASTERARBITRATION,
  '__module__' : 'proto.gnmi_ext.gnmi_ext_pb2'
  # @@protoc_insertion_point(class_scope:gnmi_ext.MasterArbitration)
  })
_sym_db.RegisterMessage(MasterArbitration)

Uint128 = _reflection.GeneratedProtocolMessageType('Uint128', (_message.Message,), {
  'DESCRIPTOR' : _UINT128,
  '__module__' : 'proto.gnmi_ext.gnmi_ext_pb2'
  # @@protoc_insertion_point(class_scope:gnmi_ext.Uint128)
  })
_sym_db.RegisterMessage(Uint128)

Role = _reflection.GeneratedProtocolMessageType('Role', (_message.Message,), {
  'DESCRIPTOR' : _ROLE,
  '__module__' : 'proto.gnmi_ext.gnmi_ext_pb2'
  # @@protoc_insertion_point(class_scope:gnmi_ext.Role)
  })
_sym_db.RegisterMessage(Role)


DESCRIPTOR._options = None
# @@protoc_insertion_point(module_scope)
