from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Question(_message.Message):
    __slots__ = ("board", "solverType")
    BOARD_FIELD_NUMBER: _ClassVar[int]
    SOLVERTYPE_FIELD_NUMBER: _ClassVar[int]
    board: _containers.RepeatedScalarFieldContainer[int]
    solverType: str
    def __init__(self, board: _Optional[_Iterable[int]] = ..., solverType: _Optional[str] = ...) -> None: ...

class Solution(_message.Message):
    __slots__ = ("solved", "board", "duration")
    SOLVED_FIELD_NUMBER: _ClassVar[int]
    BOARD_FIELD_NUMBER: _ClassVar[int]
    DURATION_FIELD_NUMBER: _ClassVar[int]
    solved: bool
    board: _containers.RepeatedScalarFieldContainer[int]
    duration: float
    def __init__(self, solved: bool = ..., board: _Optional[_Iterable[int]] = ..., duration: _Optional[float] = ...) -> None: ...
