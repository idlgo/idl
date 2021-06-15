struct PenGetPenResponseDataItem{
	0: required string Name
	1: required bool IsAsync
	2: required i32 ConfigType
}
struct PenGetPenResponseData{
	0: required i64 Total
	1: required list<PenGetPenResponseDataItem>  Items
}
enum PenEnumStatusType{
	Draft
	Validated
	Pending
	Published
}
struct PenGetPenResponse{
	0: required i64 Code
	1: required string Message
	2: required PenGetPenResponseData Data
}
struct PenGetPenRequest{
	0: required string EventCode
	1: required PenEnumStatusType Status
	2: required string Name
	3: required i64 Page
	4: required i64 Limit
}
struct StoreResponse{
	0: required list<PenGetPenResponse>  PenL
	1: optional map<string,PenGetPenResponse>  PenM
}
service appleBan{
	StoreResponse BuyPen(0: PenGetPenRequest arg0,)(path='/store/buy/pen',typ='json',)
}