package wkdb

type DB interface {
	Open() error
	Close() error
	// AppendMessages appends messages to the db.
	AppendMessages(channelId string, channelType uint8, msgs []Message) error
	// LoadPrevRangeMsgs 向上加载指定范围的消息 end=0表示不做限制 比如 start=100 end=0 limit=10 则返回的消息seq为91-100的消息, 比如 start=100 end=95 limit=10 则返回的消息seq为96-100的消息
	// 结果包含start,不包含end
	LoadPrevRangeMsgs(channelId string, channelType uint8, startMessageSeq, endMessageSeq uint64, limit int) ([]Message, error)

	// LoadNextRangeMsgs 向下加载指定范围的消息 end=0表示不做限制 比如 start=100 end=200 limit=10 则返回的消息seq为101-111的消息，
	// 比如start=100 end=105 limit=10 则返回的消息seq为101-104的消息
	// 结果包含start,不包含end
	LoadNextRangeMsgs(channelId string, channelType uint8, start, end uint64, limit int) ([]Message, error)

	LoadNextRangeMsgsForSize(channelId string, channelType uint8, startMessageSeq, endMessageSeq uint64, limitSize uint64) ([]Message, error)
	// LoadMsg 加载指定seq的消息
	LoadMsg(channelId string, channelType uint8, seq uint64) (Message, error)
	// // TruncateLogTo 截断消息, 从messageSeq开始截断,messageSeq=0 表示清空所有日志 （保留下来的内容包含messageSeq）
	TruncateLogTo(channelId string, channelType uint8, messageSeq uint64) error

	// LoadLastMsgsWithEnd 加载最新的消息 endMessageSeq表示加载到endMessageSeq的位置结束加载 endMessageSeq=0表示不做限制 结果不包含endMessageSeq
	LoadLastMsgsWithEnd(channelID string, channelType uint8, endMessageSeq uint64, limit int) ([]Message, error)
	// LoadLastMsgs 加载最后的消息
	LoadLastMsgs(channelID string, channelType uint8, limit int) ([]Message, error)

	// GetChannelLastMessageSeq 获取最后一条消息的seq
	GetChannelLastMessageSeq(channelId string, channelType uint8) (uint64, error)

	SetChannelLastMessageSeq(channelId string, channelType uint8, seq uint64) error
}
