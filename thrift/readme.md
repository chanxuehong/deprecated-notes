```Go
func NewThriftHandlerFunc(processor TProcessor, inPfactory, outPfactory TProtocolFactory) func(w http.ResponseWriter, r *http.Request)
func NewStoredMessageProtocol(protocol TProtocol, name string, typeId TMessageType, seqid int32) *storedMessageProtocol
func Skip(self TProtocol, fieldType TType, maxDepth int) (err error)
func SkipDefaultDepth(prot TProtocol, typeId TType) (err error)
func PrependError(prepend string, err error) error

type TSerializer
    func NewTSerializer() *TSerializer
    func (t *TSerializer) Write(msg TStruct) (b []byte, err error)
    func (t *TSerializer) WriteString(msg TStruct) (s string, err error)
	
type TDeserializer
    func NewTDeserializer() *TDeserializer
    func (t *TDeserializer) Read(msg TStruct, b []byte) (err error)
    func (t *TDeserializer) ReadString(msg TStruct, s string) (err error)
	
type THttpClientOptions

type TType
    func (p TType) String() string
type TMessageType

// helper --------------------------------------------------------------------------------------------------------------
func BoolPtr(v bool) *bool
func ByteSlicePtr(v []byte) *[]byte
func Float32Ptr(v float32) *float32
func Float64Ptr(v float64) *float64
func Int32Ptr(v int32) *int32
func Int64Ptr(v int64) *int64
func IntPtr(v int) *int
func StringPtr(v string) *string
func Uint32Ptr(v uint32) *uint32
func Uint64Ptr(v uint64) *uint64

// interface -----------------------------------------------------------------------------------------------------------
type Numeric interface {
	Int64() int64
	Int32() int32
	Int16() int16
	Byte() byte
	Int() int
	Float64() float64
	Float32() float32
	String() string
	isNull() bool
}
    func NewNullNumeric() Numeric
    func NewNumericFromDouble(dValue float64) Numeric
    func NewNumericFromI32(iValue int32) Numeric
    func NewNumericFromI64(iValue int64) Numeric
    func NewNumericFromJSONString(sValue string, isNull bool) Numeric
    func NewNumericFromString(sValue string) Numeric
type TStruct interface {
	Write(p TProtocol) error
	Read(p TProtocol) error
}

type Flusher interface {
	Flush() (err error)
}
type ReadSizeProvider interface {
	RemainingBytes() (num_bytes uint64)
}

type TException interface {
	error
}
type TApplicationException interface {
	TException
	TypeId() int32
	Read(iprot TProtocol) (TApplicationException, error)
	Write(oprot TProtocol) error
}
    func NewTApplicationException(type_ int32, message string) TApplicationException
type TProtocolException interface {
	TException
	TypeId() int
}
    func NewTProtocolException(err error) TProtocolException
    func NewTProtocolExceptionWithType(errType int, err error) TProtocolException
type TTransportException interface {
	TException
	TypeId() int
	Err() error
}
    func NewTTransportException(t int, e string) TTransportException
    func NewTTransportExceptionFromError(e error) TTransportException

type TServer interface {
	ProcessorFactory() TProcessorFactory
	ServerTransport() TServerTransport
	InputTransportFactory() TTransportFactory
	OutputTransportFactory() TTransportFactory
	InputProtocolFactory() TProtocolFactory
	OutputProtocolFactory() TProtocolFactory

	// Starts the server
	Serve() error
	// Stops the server. This is optional on a per-implementation basis. Not
	// all servers are required to be cleanly stoppable.
	Stop() error
}

type TProcessor interface {
	Process(in, out TProtocol) (bool, TException)
}
type TProcessorFactory interface {
	GetProcessor(trans TTransport) TProcessor
}
    func NewTProcessorFactory(p TProcessor) TProcessorFactory

type TProcessorFunction interface {
	Process(seqId int32, in, out TProtocol) (bool, TException)
}
type TProcessorFunctionFactory interface {
	GetProcessorFunction(trans TTransport) TProcessorFunction
}
    func NewTProcessorFunctionFactory(p TProcessorFunction) TProcessorFunctionFactory


type TProtocol interface {
	WriteMessageBegin(name string, typeId TMessageType, seqid int32) error
	WriteMessageEnd() error
	WriteStructBegin(name string) error
	WriteStructEnd() error
	WriteFieldBegin(name string, typeId TType, id int16) error
	WriteFieldEnd() error
	WriteFieldStop() error
	WriteMapBegin(keyType TType, valueType TType, size int) error
	WriteMapEnd() error
	WriteListBegin(elemType TType, size int) error
	WriteListEnd() error
	WriteSetBegin(elemType TType, size int) error
	WriteSetEnd() error
	WriteBool(value bool) error
	WriteByte(value int8) error
	WriteI16(value int16) error
	WriteI32(value int32) error
	WriteI64(value int64) error
	WriteDouble(value float64) error
	WriteString(value string) error
	WriteBinary(value []byte) error

	ReadMessageBegin() (name string, typeId TMessageType, seqid int32, err error)
	ReadMessageEnd() error
	ReadStructBegin() (name string, err error)
	ReadStructEnd() error
	ReadFieldBegin() (name string, typeId TType, id int16, err error)
	ReadFieldEnd() error
	ReadMapBegin() (keyType TType, valueType TType, size int, err error)
	ReadMapEnd() error
	ReadListBegin() (elemType TType, size int, err error)
	ReadListEnd() error
	ReadSetBegin() (elemType TType, size int, err error)
	ReadSetEnd() error
	ReadBool() (value bool, err error)
	ReadByte() (value int8, err error)
	ReadI16() (value int16, err error)
	ReadI32() (value int32, err error)
	ReadI64() (value int64, err error)
	ReadDouble() (value float64, err error)
	ReadString() (value string, err error)
	ReadBinary() (value []byte, err error)

	Skip(fieldType TType) (err error)
	Flush() (err error)

	Transport() TTransport
}
type TProtocolFactory interface {
	GetProtocol(trans TTransport) TProtocol
}
type ProtocolFactory interface {
	GetProtocol(t TTransport) TProtocol
}

type TTransport interface {
	io.ReadWriteCloser
	Flusher
	ReadSizeProvider

	// Opens the transport for communication
	Open() error

	// Returns true if the transport is open
	IsOpen() bool
}
type TTransportFactory interface {
	GetTransport(trans TTransport) TTransport
}
type TRichTransport interface {
	io.ReadWriter
	io.ByteReader
	io.ByteWriter
	stringWriter
	Flusher
	ReadSizeProvider
}
type TServerTransport interface {
	Listen() error
	Accept() (TTransport, error)
	Close() error

	// Optional method implementation. This signals to the server transport
	// that it should break out of any accept() or listen() that it is currently
	// blocked on. This method, if implemented, MUST be thread safe, as it may
	// be called from a different thread context than the other TServerTransport
	// methods.
	Interrupt() error
}

// Server
type TSimpleServer
    func NewTSimpleServer2(processor TProcessor, serverTransport TServerTransport) *TSimpleServer
    func NewTSimpleServer4(processor TProcessor, serverTransport TServerTransport, transportFactory TTransportFactory, protocolFactory TProtocolFactory) *TSimpleServer
    func NewTSimpleServer6(processor TProcessor, serverTransport TServerTransport, inputTransportFactory TTransportFactory, outputTransportFactory TTransportFactory, inputProtocolFactory TProtocolFactory, outputProtocolFactory TProtocolFactory) *TSimpleServer
    func NewTSimpleServerFactory2(processorFactory TProcessorFactory, serverTransport TServerTransport) *TSimpleServer
    func NewTSimpleServerFactory4(processorFactory TProcessorFactory, serverTransport TServerTransport, transportFactory TTransportFactory, protocolFactory TProtocolFactory) *TSimpleServer
    func NewTSimpleServerFactory6(processorFactory TProcessorFactory, serverTransport TServerTransport, inputTransportFactory TTransportFactory, outputTransportFactory TTransportFactory, inputProtocolFactory TProtocolFactory, outputProtocolFactory TProtocolFactory) *TSimpleServer
    func (p *TSimpleServer) AcceptLoop() error
    func (p *TSimpleServer) InputProtocolFactory() TProtocolFactory
    func (p *TSimpleServer) InputTransportFactory() TTransportFactory
    func (p *TSimpleServer) Listen() error
    func (p *TSimpleServer) OutputProtocolFactory() TProtocolFactory
    func (p *TSimpleServer) OutputTransportFactory() TTransportFactory
    func (p *TSimpleServer) ProcessorFactory() TProcessorFactory
    func (p *TSimpleServer) Serve() error
    func (p *TSimpleServer) ServerTransport() TServerTransport
    func (p *TSimpleServer) Stop() error

// Processor -----------------------------------------------------------------------------------------------------------
type TMultiplexedProcessor
    func NewTMultiplexedProcessor() *TMultiplexedProcessor
    func (t *TMultiplexedProcessor) Process(in, out TProtocol) (bool, TException)
    func (t *TMultiplexedProcessor) RegisterDefault(processor TProcessor)
    func (t *TMultiplexedProcessor) RegisterProcessor(name string, processor TProcessor)
	
// ProtocolFactory -----------------------------------------------------------------------------------------------------
type TCompactProtocolFactory
    func NewTCompactProtocolFactory() *TCompactProtocolFactory
    func (p *TCompactProtocolFactory) GetProtocol(trans TTransport) TProtocol
type TDebugProtocolFactory
    func NewTDebugProtocolFactory(underlying TProtocolFactory, logPrefix string) *TDebugProtocolFactory
    func (t *TDebugProtocolFactory) GetProtocol(trans TTransport) TProtocol
type TJSONProtocolFactory
    func NewTJSONProtocolFactory() *TJSONProtocolFactory
    func (p *TJSONProtocolFactory) GetProtocol(trans TTransport) TProtocol
type TSimpleJSONProtocolFactory
    func NewTSimpleJSONProtocolFactory() *TSimpleJSONProtocolFactory
    func (p *TSimpleJSONProtocolFactory) GetProtocol(trans TTransport) TProtocol
	
// Protocol ------------------------------------------------------------------------------------------------------------
type TBinaryProtocol
    func NewTBinaryProtocol(t TTransport, strictRead, strictWrite bool) *TBinaryProtocol
    func NewTBinaryProtocolTransport(t TTransport) *TBinaryProtocol
    func (p *TBinaryProtocol) Flush() (err error)
    func (p *TBinaryProtocol) ReadBinary() ([]byte, error)
    func (p *TBinaryProtocol) ReadBool() (bool, error)
    func (p *TBinaryProtocol) ReadByte() (int8, error)
    func (p *TBinaryProtocol) ReadDouble() (value float64, err error)
    func (p *TBinaryProtocol) ReadFieldBegin() (name string, typeId TType, seqId int16, err error)
    func (p *TBinaryProtocol) ReadFieldEnd() error
    func (p *TBinaryProtocol) ReadI16() (value int16, err error)
    func (p *TBinaryProtocol) ReadI32() (value int32, err error)
    func (p *TBinaryProtocol) ReadI64() (value int64, err error)
    func (p *TBinaryProtocol) ReadListBegin() (elemType TType, size int, err error)
    func (p *TBinaryProtocol) ReadListEnd() error
    func (p *TBinaryProtocol) ReadMapBegin() (kType, vType TType, size int, err error)
    func (p *TBinaryProtocol) ReadMapEnd() error
    func (p *TBinaryProtocol) ReadMessageBegin() (name string, typeId TMessageType, seqId int32, err error)
    func (p *TBinaryProtocol) ReadMessageEnd() error
    func (p *TBinaryProtocol) ReadSetBegin() (elemType TType, size int, err error)
    func (p *TBinaryProtocol) ReadSetEnd() error
    func (p *TBinaryProtocol) ReadString() (value string, err error)
    func (p *TBinaryProtocol) ReadStructBegin() (name string, err error)
    func (p *TBinaryProtocol) ReadStructEnd() error
    func (p *TBinaryProtocol) Skip(fieldType TType) (err error)
    func (p *TBinaryProtocol) Transport() TTransport
    func (p *TBinaryProtocol) WriteBinary(value []byte) error
    func (p *TBinaryProtocol) WriteBool(value bool) error
    func (p *TBinaryProtocol) WriteByte(value int8) error
    func (p *TBinaryProtocol) WriteDouble(value float64) error
    func (p *TBinaryProtocol) WriteFieldBegin(name string, typeId TType, id int16) error
    func (p *TBinaryProtocol) WriteFieldEnd() error
    func (p *TBinaryProtocol) WriteFieldStop() error
    func (p *TBinaryProtocol) WriteI16(value int16) error
    func (p *TBinaryProtocol) WriteI32(value int32) error
    func (p *TBinaryProtocol) WriteI64(value int64) error
    func (p *TBinaryProtocol) WriteListBegin(elemType TType, size int) error
    func (p *TBinaryProtocol) WriteListEnd() error
    func (p *TBinaryProtocol) WriteMapBegin(keyType TType, valueType TType, size int) error
    func (p *TBinaryProtocol) WriteMapEnd() error
    func (p *TBinaryProtocol) WriteMessageBegin(name string, typeId TMessageType, seqId int32) error
    func (p *TBinaryProtocol) WriteMessageEnd() error
    func (p *TBinaryProtocol) WriteSetBegin(elemType TType, size int) error
    func (p *TBinaryProtocol) WriteSetEnd() error
    func (p *TBinaryProtocol) WriteString(value string) error
    func (p *TBinaryProtocol) WriteStructBegin(name string) error
    func (p *TBinaryProtocol) WriteStructEnd() error
type TCompactProtocol
    func NewTCompactProtocol(trans TTransport) *TCompactProtocol
    func (p *TCompactProtocol) Flush() (err error)
    func (p *TCompactProtocol) ReadBinary() (value []byte, err error)
    func (p *TCompactProtocol) ReadBool() (value bool, err error)
    func (p *TCompactProtocol) ReadByte() (int8, error)
    func (p *TCompactProtocol) ReadDouble() (value float64, err error)
    func (p *TCompactProtocol) ReadFieldBegin() (name string, typeId TType, id int16, err error)
    func (p *TCompactProtocol) ReadFieldEnd() error
    func (p *TCompactProtocol) ReadI16() (value int16, err error)
    func (p *TCompactProtocol) ReadI32() (value int32, err error)
    func (p *TCompactProtocol) ReadI64() (value int64, err error)
    func (p *TCompactProtocol) ReadListBegin() (elemType TType, size int, err error)
    func (p *TCompactProtocol) ReadListEnd() error
    func (p *TCompactProtocol) ReadMapBegin() (keyType TType, valueType TType, size int, err error)
    func (p *TCompactProtocol) ReadMapEnd() error
    func (p *TCompactProtocol) ReadMessageBegin() (name string, typeId TMessageType, seqId int32, err error)
    func (p *TCompactProtocol) ReadMessageEnd() error
    func (p *TCompactProtocol) ReadSetBegin() (elemType TType, size int, err error)
    func (p *TCompactProtocol) ReadSetEnd() error
    func (p *TCompactProtocol) ReadString() (value string, err error)
    func (p *TCompactProtocol) ReadStructBegin() (name string, err error)
    func (p *TCompactProtocol) ReadStructEnd() error
    func (p *TCompactProtocol) Skip(fieldType TType) (err error)
    func (p *TCompactProtocol) Transport() TTransport
    func (p *TCompactProtocol) WriteBinary(bin []byte) error
    func (p *TCompactProtocol) WriteBool(value bool) error
    func (p *TCompactProtocol) WriteByte(value int8) error
    func (p *TCompactProtocol) WriteDouble(value float64) error
    func (p *TCompactProtocol) WriteFieldBegin(name string, typeId TType, id int16) error
    func (p *TCompactProtocol) WriteFieldEnd() error
    func (p *TCompactProtocol) WriteFieldStop() error
    func (p *TCompactProtocol) WriteI16(value int16) error
    func (p *TCompactProtocol) WriteI32(value int32) error
    func (p *TCompactProtocol) WriteI64(value int64) error
    func (p *TCompactProtocol) WriteListBegin(elemType TType, size int) error
    func (p *TCompactProtocol) WriteListEnd() error
    func (p *TCompactProtocol) WriteMapBegin(keyType TType, valueType TType, size int) error
    func (p *TCompactProtocol) WriteMapEnd() error
    func (p *TCompactProtocol) WriteMessageBegin(name string, typeId TMessageType, seqid int32) error
    func (p *TCompactProtocol) WriteMessageEnd() error
    func (p *TCompactProtocol) WriteSetBegin(elemType TType, size int) error
    func (p *TCompactProtocol) WriteSetEnd() error
    func (p *TCompactProtocol) WriteString(value string) error
    func (p *TCompactProtocol) WriteStructBegin(name string) error
    func (p *TCompactProtocol) WriteStructEnd() error
type TDebugProtocol
    func (tdp *TDebugProtocol) Flush() (err error)
    func (tdp *TDebugProtocol) ReadBinary() (value []byte, err error)
    func (tdp *TDebugProtocol) ReadBool() (value bool, err error)
    func (tdp *TDebugProtocol) ReadByte() (value int8, err error)
    func (tdp *TDebugProtocol) ReadDouble() (value float64, err error)
    func (tdp *TDebugProtocol) ReadFieldBegin() (name string, typeId TType, id int16, err error)
    func (tdp *TDebugProtocol) ReadFieldEnd() (err error)
    func (tdp *TDebugProtocol) ReadI16() (value int16, err error)
    func (tdp *TDebugProtocol) ReadI32() (value int32, err error)
    func (tdp *TDebugProtocol) ReadI64() (value int64, err error)
    func (tdp *TDebugProtocol) ReadListBegin() (elemType TType, size int, err error)
    func (tdp *TDebugProtocol) ReadListEnd() (err error)
    func (tdp *TDebugProtocol) ReadMapBegin() (keyType TType, valueType TType, size int, err error)
    func (tdp *TDebugProtocol) ReadMapEnd() (err error)
    func (tdp *TDebugProtocol) ReadMessageBegin() (name string, typeId TMessageType, seqid int32, err error)
    func (tdp *TDebugProtocol) ReadMessageEnd() (err error)
    func (tdp *TDebugProtocol) ReadSetBegin() (elemType TType, size int, err error)
    func (tdp *TDebugProtocol) ReadSetEnd() (err error)
    func (tdp *TDebugProtocol) ReadString() (value string, err error)
    func (tdp *TDebugProtocol) ReadStructBegin() (name string, err error)
    func (tdp *TDebugProtocol) ReadStructEnd() (err error)
    func (tdp *TDebugProtocol) Skip(fieldType TType) (err error)
    func (tdp *TDebugProtocol) Transport() TTransport
    func (tdp *TDebugProtocol) WriteBinary(value []byte) error
    func (tdp *TDebugProtocol) WriteBool(value bool) error
    func (tdp *TDebugProtocol) WriteByte(value int8) error
    func (tdp *TDebugProtocol) WriteDouble(value float64) error
    func (tdp *TDebugProtocol) WriteFieldBegin(name string, typeId TType, id int16) error
    func (tdp *TDebugProtocol) WriteFieldEnd() error
    func (tdp *TDebugProtocol) WriteFieldStop() error
    func (tdp *TDebugProtocol) WriteI16(value int16) error
    func (tdp *TDebugProtocol) WriteI32(value int32) error
    func (tdp *TDebugProtocol) WriteI64(value int64) error
    func (tdp *TDebugProtocol) WriteListBegin(elemType TType, size int) error
    func (tdp *TDebugProtocol) WriteListEnd() error
    func (tdp *TDebugProtocol) WriteMapBegin(keyType TType, valueType TType, size int) error
    func (tdp *TDebugProtocol) WriteMapEnd() error
    func (tdp *TDebugProtocol) WriteMessageBegin(name string, typeId TMessageType, seqid int32) error
    func (tdp *TDebugProtocol) WriteMessageEnd() error
    func (tdp *TDebugProtocol) WriteSetBegin(elemType TType, size int) error
    func (tdp *TDebugProtocol) WriteSetEnd() error
    func (tdp *TDebugProtocol) WriteString(value string) error
    func (tdp *TDebugProtocol) WriteStructBegin(name string) error
    func (tdp *TDebugProtocol) WriteStructEnd() error
type TJSONProtocol
    func NewTJSONProtocol(t TTransport) *TJSONProtocol
    func (p *TJSONProtocol) Flush() (err error)
    func (p *TJSONProtocol) OutputElemListBegin(elemType TType, size int) error
    func (p *TJSONProtocol) ParseElemListBegin() (elemType TType, size int, e error)
    func (p *TJSONProtocol) ReadBinary() ([]byte, error)
    func (p *TJSONProtocol) ReadBool() (bool, error)
    func (p *TJSONProtocol) ReadByte() (int8, error)
    func (p *TJSONProtocol) ReadDouble() (float64, error)
    func (p *TJSONProtocol) ReadFieldBegin() (string, TType, int16, error)
    func (p *TJSONProtocol) ReadFieldEnd() error
    func (p *TJSONProtocol) ReadI16() (int16, error)
    func (p *TJSONProtocol) ReadI32() (int32, error)
    func (p *TJSONProtocol) ReadI64() (int64, error)
    func (p *TJSONProtocol) ReadListBegin() (elemType TType, size int, e error)
    func (p *TJSONProtocol) ReadListEnd() error
    func (p *TJSONProtocol) ReadMapBegin() (keyType TType, valueType TType, size int, e error)
    func (p *TJSONProtocol) ReadMapEnd() error
    func (p *TJSONProtocol) ReadMessageBegin() (name string, typeId TMessageType, seqId int32, err error)
    func (p *TJSONProtocol) ReadMessageEnd() error
    func (p *TJSONProtocol) ReadSetBegin() (elemType TType, size int, e error)
    func (p *TJSONProtocol) ReadSetEnd() error
    func (p *TJSONProtocol) ReadString() (string, error)
    func (p *TJSONProtocol) ReadStructBegin() (name string, err error)
    func (p *TJSONProtocol) ReadStructEnd() error
    func (p *TJSONProtocol) Skip(fieldType TType) (err error)
    func (p *TJSONProtocol) StringToTypeId(fieldType string) (TType, error)
    func (p *TJSONProtocol) Transport() TTransport
    func (p *TJSONProtocol) TypeIdToString(fieldType TType) (string, error)
    func (p *TJSONProtocol) WriteBinary(v []byte) error
    func (p *TJSONProtocol) WriteBool(b bool) error
    func (p *TJSONProtocol) WriteByte(b int8) error
    func (p *TJSONProtocol) WriteDouble(v float64) error
    func (p *TJSONProtocol) WriteFieldBegin(name string, typeId TType, id int16) error
    func (p *TJSONProtocol) WriteFieldEnd() error
    func (p *TJSONProtocol) WriteFieldStop() error
    func (p *TJSONProtocol) WriteI16(v int16) error
    func (p *TJSONProtocol) WriteI32(v int32) error
    func (p *TJSONProtocol) WriteI64(v int64) error
    func (p *TJSONProtocol) WriteListBegin(elemType TType, size int) error
    func (p *TJSONProtocol) WriteListEnd() error
    func (p *TJSONProtocol) WriteMapBegin(keyType TType, valueType TType, size int) error
    func (p *TJSONProtocol) WriteMapEnd() error
    func (p *TJSONProtocol) WriteMessageBegin(name string, typeId TMessageType, seqId int32) error
    func (p *TJSONProtocol) WriteMessageEnd() error
    func (p *TJSONProtocol) WriteSetBegin(elemType TType, size int) error
    func (p *TJSONProtocol) WriteSetEnd() error
    func (p *TJSONProtocol) WriteString(v string) error
    func (p *TJSONProtocol) WriteStructBegin(name string) error
    func (p *TJSONProtocol) WriteStructEnd() error
type TMultiplexedProtocol
    func NewTMultiplexedProtocol(protocol TProtocol, serviceName string) *TMultiplexedProtocol
    func (t *TMultiplexedProtocol) WriteMessageBegin(name string, typeId TMessageType, seqid int32) error
type TSimpleJSONProtocol
    func NewTSimpleJSONProtocol(t TTransport) *TSimpleJSONProtocol
    func (p *TSimpleJSONProtocol) Flush() (err error)
    func (p *TSimpleJSONProtocol) OutputBool(value bool) error
    func (p *TSimpleJSONProtocol) OutputElemListBegin(elemType TType, size int) error
    func (p *TSimpleJSONProtocol) OutputF64(value float64) error
    func (p *TSimpleJSONProtocol) OutputI64(value int64) error
    func (p *TSimpleJSONProtocol) OutputListBegin() error
    func (p *TSimpleJSONProtocol) OutputListEnd() error
    func (p *TSimpleJSONProtocol) OutputNull() error
    func (p *TSimpleJSONProtocol) OutputObjectBegin() error
    func (p *TSimpleJSONProtocol) OutputObjectEnd() error
    func (p *TSimpleJSONProtocol) OutputPostValue() error
    func (p *TSimpleJSONProtocol) OutputPreValue() error
    func (p *TSimpleJSONProtocol) OutputString(s string) error
    func (p *TSimpleJSONProtocol) OutputStringData(s string) error
    func (p *TSimpleJSONProtocol) ParseBase64EncodedBody() ([]byte, error)
    func (p *TSimpleJSONProtocol) ParseElemListBegin() (elemType TType, size int, e error)
    func (p *TSimpleJSONProtocol) ParseF64() (float64, bool, error)
    func (p *TSimpleJSONProtocol) ParseI64() (int64, bool, error)
    func (p *TSimpleJSONProtocol) ParseListBegin() (isNull bool, err error)
    func (p *TSimpleJSONProtocol) ParseListEnd() error
    func (p *TSimpleJSONProtocol) ParseObjectEnd() error
    func (p *TSimpleJSONProtocol) ParseObjectStart() (bool, error)
    func (p *TSimpleJSONProtocol) ParsePostValue() error
    func (p *TSimpleJSONProtocol) ParsePreValue() error
    func (p *TSimpleJSONProtocol) ParseQuotedStringBody() (string, error)
    func (p *TSimpleJSONProtocol) ParseStringBody() (string, error)
    func (p *TSimpleJSONProtocol) ReadBinary() ([]byte, error)
    func (p *TSimpleJSONProtocol) ReadBool() (bool, error)
    func (p *TSimpleJSONProtocol) ReadByte() (int8, error)
    func (p *TSimpleJSONProtocol) ReadDouble() (float64, error)
    func (p *TSimpleJSONProtocol) ReadFieldBegin() (string, TType, int16, error)
    func (p *TSimpleJSONProtocol) ReadFieldEnd() error
    func (p *TSimpleJSONProtocol) ReadI16() (int16, error)
    func (p *TSimpleJSONProtocol) ReadI32() (int32, error)
    func (p *TSimpleJSONProtocol) ReadI64() (int64, error)
    func (p *TSimpleJSONProtocol) ReadListBegin() (elemType TType, size int, e error)
    func (p *TSimpleJSONProtocol) ReadListEnd() error
    func (p *TSimpleJSONProtocol) ReadMapBegin() (keyType TType, valueType TType, size int, e error)
    func (p *TSimpleJSONProtocol) ReadMapEnd() error
    func (p *TSimpleJSONProtocol) ReadMessageBegin() (name string, typeId TMessageType, seqId int32, err error)
    func (p *TSimpleJSONProtocol) ReadMessageEnd() error
    func (p *TSimpleJSONProtocol) ReadSetBegin() (elemType TType, size int, e error)
    func (p *TSimpleJSONProtocol) ReadSetEnd() error
    func (p *TSimpleJSONProtocol) ReadString() (string, error)
    func (p *TSimpleJSONProtocol) ReadStructBegin() (name string, err error)
    func (p *TSimpleJSONProtocol) ReadStructEnd() error
    func (p *TSimpleJSONProtocol) Skip(fieldType TType) (err error)
    func (p *TSimpleJSONProtocol) Transport() TTransport
    func (p *TSimpleJSONProtocol) WriteBinary(v []byte) error
    func (p *TSimpleJSONProtocol) WriteBool(b bool) error
    func (p *TSimpleJSONProtocol) WriteByte(b int8) error
    func (p *TSimpleJSONProtocol) WriteDouble(v float64) error
    func (p *TSimpleJSONProtocol) WriteFieldBegin(name string, typeId TType, id int16) error
    func (p *TSimpleJSONProtocol) WriteFieldEnd() error
    func (p *TSimpleJSONProtocol) WriteFieldStop() error
    func (p *TSimpleJSONProtocol) WriteI16(v int16) error
    func (p *TSimpleJSONProtocol) WriteI32(v int32) error
    func (p *TSimpleJSONProtocol) WriteI64(v int64) error
    func (p *TSimpleJSONProtocol) WriteListBegin(elemType TType, size int) error
    func (p *TSimpleJSONProtocol) WriteListEnd() error
    func (p *TSimpleJSONProtocol) WriteMapBegin(keyType TType, valueType TType, size int) error
    func (p *TSimpleJSONProtocol) WriteMapEnd() error
    func (p *TSimpleJSONProtocol) WriteMessageBegin(name string, typeId TMessageType, seqId int32) error
    func (p *TSimpleJSONProtocol) WriteMessageEnd() error
    func (p *TSimpleJSONProtocol) WriteSetBegin(elemType TType, size int) error
    func (p *TSimpleJSONProtocol) WriteSetEnd() error
    func (p *TSimpleJSONProtocol) WriteString(v string) error
    func (p *TSimpleJSONProtocol) WriteStructBegin(name string) error
    func (p *TSimpleJSONProtocol) WriteStructEnd() error
	
// TransportFactory ----------------------------------------------------------------------------------------------------
    func NewTFramedTransportFactory(factory TTransportFactory) TTransportFactory
    func NewTFramedTransportFactoryMaxLength(factory TTransportFactory, maxLength uint32) TTransportFactory
    func NewTTransportFactory() TTransportFactory
type StreamTransportFactory
    func NewStreamTransportFactory(reader io.Reader, writer io.Writer, isReadWriter bool) *StreamTransportFactory
    func (p *StreamTransportFactory) GetTransport(trans TTransport) TTransport
type TBinaryProtocolFactory
    func NewTBinaryProtocolFactory(strictRead, strictWrite bool) *TBinaryProtocolFactory
    func NewTBinaryProtocolFactoryDefault() *TBinaryProtocolFactory
    func (p *TBinaryProtocolFactory) GetProtocol(t TTransport) TProtocol
type TBufferedTransportFactory
    func NewTBufferedTransportFactory(bufferSize int) *TBufferedTransportFactory
    func (p *TBufferedTransportFactory) GetTransport(trans TTransport) TTransport
type THttpClientTransportFactory
    func NewTHttpClientTransportFactory(url string) *THttpClientTransportFactory
    func NewTHttpClientTransportFactoryWithOptions(url string, options THttpClientOptions) *THttpClientTransportFactory
    func NewTHttpPostClientTransportFactory(url string) *THttpClientTransportFactory
    func NewTHttpPostClientTransportFactoryWithOptions(url string, options THttpClientOptions) *THttpClientTransportFactory
    func (p *THttpClientTransportFactory) GetTransport(trans TTransport) TTransport
type TMemoryBufferTransportFactory
    func NewTMemoryBufferTransportFactory(size int) *TMemoryBufferTransportFactory
    func (p *TMemoryBufferTransportFactory) GetTransport(trans TTransport) TTransport
type TZlibTransportFactory
    func NewTZlibTransportFactory(level int) *TZlibTransportFactory
    func (p *TZlibTransportFactory) GetTransport(trans TTransport) TTransport

// Transport -----------------------------------------------------------------------------------------------------------
    func NewTHttpClient(urlstr string) (TTransport, error)
    func NewTHttpClientWithOptions(urlstr string, options THttpClientOptions) (TTransport, error)
    func NewTHttpPostClient(urlstr string) (TTransport, error)
    func NewTHttpPostClientWithOptions(urlstr string, options THttpClientOptions) (TTransport, error)
type RichTransport // TRichTransport
    func NewTRichTransport(trans TTransport) *RichTransport
    func (r *RichTransport) ReadByte() (c byte, err error)
    func (r *RichTransport) RemainingBytes() (num_bytes uint64)
    func (r *RichTransport) WriteByte(c byte) error
    func (r *RichTransport) WriteString(s string) (n int, err error)
type StreamTransport // TRichTransport
    func NewStreamTransport(r io.Reader, w io.Writer) *StreamTransport
    func NewStreamTransportR(r io.Reader) *StreamTransport
    func NewStreamTransportRW(rw io.ReadWriter) *StreamTransport
    func NewStreamTransportW(w io.Writer) *StreamTransport
    func (p *StreamTransport) Close() error
    func (p *StreamTransport) Flush() error
    func (p *StreamTransport) IsOpen() bool
    func (p *StreamTransport) Open() error
    func (p *StreamTransport) Read(c []byte) (n int, err error)
    func (p *StreamTransport) ReadByte() (c byte, err error)
    func (p *StreamTransport) RemainingBytes() (num_bytes uint64)
    func (p *StreamTransport) Write(c []byte) (n int, err error)
    func (p *StreamTransport) WriteByte(c byte) (err error)
    func (p *StreamTransport) WriteString(s string) (n int, err error)
type TBufferedTransport
    func NewTBufferedTransport(trans TTransport, bufferSize int) *TBufferedTransport
    func (p *TBufferedTransport) Close() (err error)
    func (p *TBufferedTransport) Flush() error
    func (p *TBufferedTransport) IsOpen() bool
    func (p *TBufferedTransport) Open() (err error)
    func (p *TBufferedTransport) Read(b []byte) (int, error)
    func (p *TBufferedTransport) RemainingBytes() (num_bytes uint64)
    func (p *TBufferedTransport) Write(b []byte) (int, error)
type TFramedTransport // TRichTransport
    func NewTFramedTransport(transport TTransport) *TFramedTransport
    func NewTFramedTransportMaxLength(transport TTransport, maxLength uint32) *TFramedTransport
    func (p *TFramedTransport) Close() error
    func (p *TFramedTransport) Flush() error
    func (p *TFramedTransport) IsOpen() bool
    func (p *TFramedTransport) Open() error
    func (p *TFramedTransport) Read(buf []byte) (l int, err error)
    func (p *TFramedTransport) ReadByte() (c byte, err error)
    func (p *TFramedTransport) RemainingBytes() (num_bytes uint64)
    func (p *TFramedTransport) Write(buf []byte) (int, error)
    func (p *TFramedTransport) WriteByte(c byte) error
    func (p *TFramedTransport) WriteString(s string) (n int, err error)
type THttpClient // TRichTransport
    func (p *THttpClient) Close() error
    func (p *THttpClient) DelHeader(key string)
    func (p *THttpClient) Flush() error
    func (p *THttpClient) GetHeader(key string) string
    func (p *THttpClient) IsOpen() bool
    func (p *THttpClient) Open() error
    func (p *THttpClient) Read(buf []byte) (int, error)
    func (p *THttpClient) ReadByte() (c byte, err error)
    func (p *THttpClient) RemainingBytes() (num_bytes uint64)
    func (p *THttpClient) SetHeader(key string, value string)
    func (p *THttpClient) Write(buf []byte) (int, error)
    func (p *THttpClient) WriteByte(c byte) error
    func (p *THttpClient) WriteString(s string) (n int, err error)
type TMemoryBuffer
    func NewTMemoryBuffer() *TMemoryBuffer
    func NewTMemoryBufferLen(size int) *TMemoryBuffer
    func (p *TMemoryBuffer) Close() error
    func (p *TMemoryBuffer) Flush() error
    func (p *TMemoryBuffer) IsOpen() bool
    func (p *TMemoryBuffer) Open() error
    func (p *TMemoryBuffer) RemainingBytes() (num_bytes uint64)
type TSSLSocket
    func NewTSSLSocket(hostPort string, cfg *tls.Config) (*TSSLSocket, error)
    func NewTSSLSocketFromAddrTimeout(addr net.Addr, cfg *tls.Config, timeout time.Duration) *TSSLSocket
    func NewTSSLSocketFromConnTimeout(conn net.Conn, cfg *tls.Config, timeout time.Duration) *TSSLSocket
    func NewTSSLSocketTimeout(hostPort string, cfg *tls.Config, timeout time.Duration) (*TSSLSocket, error)
    func (p *TSSLSocket) Close() error
    func (p *TSSLSocket) Conn() net.Conn
    func (p *TSSLSocket) Flush() error
    func (p *TSSLSocket) Interrupt() error
    func (p *TSSLSocket) IsOpen() bool
    func (p *TSSLSocket) Open() error
    func (p *TSSLSocket) Read(buf []byte) (int, error)
    func (p *TSSLSocket) RemainingBytes() (num_bytes uint64)
    func (p *TSSLSocket) SetTimeout(timeout time.Duration) error
    func (p *TSSLSocket) Write(buf []byte) (int, error)
type TSocket
    func NewTSocket(hostPort string) (*TSocket, error)
    func NewTSocketFromAddrTimeout(addr net.Addr, timeout time.Duration) *TSocket
    func NewTSocketFromConnTimeout(conn net.Conn, timeout time.Duration) *TSocket
    func NewTSocketTimeout(hostPort string, timeout time.Duration) (*TSocket, error)
    func (p *TSocket) Addr() net.Addr
    func (p *TSocket) Close() error
    func (p *TSocket) Conn() net.Conn
    func (p *TSocket) Flush() error
    func (p *TSocket) Interrupt() error
    func (p *TSocket) IsOpen() bool
    func (p *TSocket) Open() error
    func (p *TSocket) Read(buf []byte) (int, error)
    func (p *TSocket) RemainingBytes() (num_bytes uint64)
    func (p *TSocket) SetTimeout(timeout time.Duration) error
    func (p *TSocket) Write(buf []byte) (int, error)
type TZlibTransport
    func NewTZlibTransport(trans TTransport, level int) (*TZlibTransport, error)
    func (z *TZlibTransport) Close() error
    func (z *TZlibTransport) Flush() error
    func (z *TZlibTransport) IsOpen() bool
    func (z *TZlibTransport) Open() error
    func (z *TZlibTransport) Read(p []byte) (int, error)
    func (z *TZlibTransport) RemainingBytes() uint64
    func (z *TZlibTransport) Write(p []byte) (int, error)

// TServerTransport ----------------------------------------------------------------------------------------------------
type TSSLServerSocket
    func NewTSSLServerSocket(listenAddr string, cfg *tls.Config) (*TSSLServerSocket, error)
    func NewTSSLServerSocketTimeout(listenAddr string, cfg *tls.Config, clientTimeout time.Duration) (*TSSLServerSocket, error)
    func (p *TSSLServerSocket) Accept() (TTransport, error)
    func (p *TSSLServerSocket) Addr() net.Addr
    func (p *TSSLServerSocket) Close() error
    func (p *TSSLServerSocket) Interrupt() error
    func (p *TSSLServerSocket) IsListening() bool
    func (p *TSSLServerSocket) Listen() error
    func (p *TSSLServerSocket) Open() error
type TServerSocket
    func NewTServerSocket(listenAddr string) (*TServerSocket, error)
    func NewTServerSocketTimeout(listenAddr string, clientTimeout time.Duration) (*TServerSocket, error)
    func (p *TServerSocket) Accept() (TTransport, error)
    func (p *TServerSocket) Addr() net.Addr
    func (p *TServerSocket) Close() error
    func (p *TServerSocket) Interrupt() error
    func (p *TServerSocket) IsListening() bool
    func (p *TServerSocket) Listen() error
    func (p *TServerSocket) Open() error
```