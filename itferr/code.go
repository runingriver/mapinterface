package itferr

//go:generate stringer -type=MapItfErrorCode
type MapItfErrorCode int

const (
	UnknownErr       MapItfErrorCode = -1
	InitParseFailed  MapItfErrorCode = 1001
	InitParamTypeErr MapItfErrorCode = 2001
	ExceptObject     MapItfErrorCode = 2002

	KeyTypeErr              MapItfErrorCode = 3001
	ValueTypeErr            MapItfErrorCode = 3002
	ValueConvertFailed      MapItfErrorCode = 3003
	BaseTypeConvertFailed   MapItfErrorCode = 3004
	KeyNotFound             MapItfErrorCode = 3005
	GetFuncTypeInconsistent MapItfErrorCode = 3006
	IllegalMapObject        MapItfErrorCode = 3007
	EmptyMapObject          MapItfErrorCode = 3008

	ListIndexIllegal MapItfErrorCode = 4001

	UnSupportInterfaceFunc MapItfErrorCode = 5001
	CurrentCannotUseIndex  MapItfErrorCode = 5002
	TypeMismatchErr        MapItfErrorCode = 5003
	FuncUsedErr            MapItfErrorCode = 5004
	UnrecoverablePanicErr  MapItfErrorCode = 5005

	SetValueErr              MapItfErrorCode = 6001
	UnSupportSetValTypeErr   MapItfErrorCode = 6002
	IterChainIsEmpty         MapItfErrorCode = 6003
	IterChainPreElementIsNil MapItfErrorCode = 6004
)
