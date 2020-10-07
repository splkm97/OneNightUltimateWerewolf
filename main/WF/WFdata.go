package WF

const (
	MaxUser = 6
)

type WFData struct {
	CardDeck     CardDeck
	IndexChan    chan int
	CurStage     string
	AdminUserID  string
	UseChannelID string
	UserIDs      []string
	UserRole     map[string]string
	FinalRole    map[string]string
	DeadUserMap  map[string]bool
}

func NewWFData(uid, cid string) *WFData {
	nc := NewCardDeck()
	nc.ShuffleCards()
	return &WFData{
		CardDeck:     *nc,
		IndexChan:    make(chan int, 10),
		CurStage:     "Prepare",
		AdminUserID:  uid,
		UseChannelID: cid,
		UserRole:     make(map[string]string),
		FinalRole:    make(map[string]string),
		UserIDs:      make([]string, 0, MaxUser),
		DeadUserMap:  make(map[string]bool),
	}
}

func (wfd *WFData) AppendUser(uid string) {
	wfd.UserIDs = append(wfd.UserIDs, uid)
	wfd.UserRole[uid] = wfd.CardDeck.PopCard()
	wfd.FinalRole[uid] = wfd.UserRole[uid]
}
