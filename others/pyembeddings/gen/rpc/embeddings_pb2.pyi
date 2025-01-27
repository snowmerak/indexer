from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class GetEmbeddingsRequest(_message.Message):
    __slots__ = ("model", "contents")
    MODEL_FIELD_NUMBER: _ClassVar[int]
    CONTENTS_FIELD_NUMBER: _ClassVar[int]
    model: str
    contents: str
    def __init__(self, model: _Optional[str] = ..., contents: _Optional[str] = ...) -> None: ...

class GetEmbeddingsResponse(_message.Message):
    __slots__ = ("embeddings",)
    EMBEDDINGS_FIELD_NUMBER: _ClassVar[int]
    embeddings: _containers.RepeatedScalarFieldContainer[float]
    def __init__(self, embeddings: _Optional[_Iterable[float]] = ...) -> None: ...
