# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: plugin.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0cplugin.proto\x12\x05proto\x1a\x19google/protobuf/any.proto\"}\n\x05Input\x12$\n\x04\x64\x61ta\x18\x01 \x03(\x0b\x32\x16.proto.Input.DataEntry\x12!\n\x03\x61ny\x18\x02 \x01(\x0b\x32\x14.google.protobuf.Any\x1a+\n\tDataEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\"\x9e\x01\n\x06OutPut\x12!\n\x03\x61ny\x18\x01 \x01(\x0b\x32\x14.google.protobuf.Any\x12\n\n\x02rt\x18\x02 \x01(\x03\x12\x11\n\tstartTime\x18\x03 \x01(\x03\x12%\n\x04\x64\x61ta\x18\x04 \x03(\x0b\x32\x17.proto.OutPut.DataEntry\x1a+\n\tDataEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\x32,\n\x06Plugin\x12\"\n\x03Run\x12\x0c.proto.Input\x1a\r.proto.OutPutB\tZ\x07.;protob\x06proto3')



_INPUT = DESCRIPTOR.message_types_by_name['Input']
_INPUT_DATAENTRY = _INPUT.nested_types_by_name['DataEntry']
_OUTPUT = DESCRIPTOR.message_types_by_name['OutPut']
_OUTPUT_DATAENTRY = _OUTPUT.nested_types_by_name['DataEntry']
Input = _reflection.GeneratedProtocolMessageType('Input', (_message.Message,), {

  'DataEntry' : _reflection.GeneratedProtocolMessageType('DataEntry', (_message.Message,), {
    'DESCRIPTOR' : _INPUT_DATAENTRY,
    '__module__' : 'plugin_pb2'
    # @@protoc_insertion_point(class_scope:proto.Input.DataEntry)
    })
  ,
  'DESCRIPTOR' : _INPUT,
  '__module__' : 'plugin_pb2'
  # @@protoc_insertion_point(class_scope:proto.Input)
  })
_sym_db.RegisterMessage(Input)
_sym_db.RegisterMessage(Input.DataEntry)

OutPut = _reflection.GeneratedProtocolMessageType('OutPut', (_message.Message,), {

  'DataEntry' : _reflection.GeneratedProtocolMessageType('DataEntry', (_message.Message,), {
    'DESCRIPTOR' : _OUTPUT_DATAENTRY,
    '__module__' : 'plugin_pb2'
    # @@protoc_insertion_point(class_scope:proto.OutPut.DataEntry)
    })
  ,
  'DESCRIPTOR' : _OUTPUT,
  '__module__' : 'plugin_pb2'
  # @@protoc_insertion_point(class_scope:proto.OutPut)
  })
_sym_db.RegisterMessage(OutPut)
_sym_db.RegisterMessage(OutPut.DataEntry)

_PLUGIN = DESCRIPTOR.services_by_name['Plugin']
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\007.;proto'
  _INPUT_DATAENTRY._options = None
  _INPUT_DATAENTRY._serialized_options = b'8\001'
  _OUTPUT_DATAENTRY._options = None
  _OUTPUT_DATAENTRY._serialized_options = b'8\001'
  _INPUT._serialized_start=50
  _INPUT._serialized_end=175
  _INPUT_DATAENTRY._serialized_start=132
  _INPUT_DATAENTRY._serialized_end=175
  _OUTPUT._serialized_start=178
  _OUTPUT._serialized_end=336
  _OUTPUT_DATAENTRY._serialized_start=132
  _OUTPUT_DATAENTRY._serialized_end=175
  _PLUGIN._serialized_start=338
  _PLUGIN._serialized_end=382
# @@protoc_insertion_point(module_scope)
