package WF

const (
	MaxUser = 6
)

type WFData struct {
	CardDeck     CardDeck
	CurStage     string
	AdminUserID  string
	UseChannelID string
	UserIDs      []string
	UserRole     map[string]string
	DeadUserMap  map[string]bool
}

func NewWFData(uid, cid string) *WFData {
	nc := NewCardDeck()
	nc.ShuffleCards()
	return &WFData{
		CardDeck:     *nc,
		CurStage:     "Prepare",
		AdminUserID:  uid,
		UseChannelID: cid,
		UserRole:     make(map[string]string),
		UserIDs:      make([]string, 0, MaxUser),
		DeadUserMap:  make(map[string]bool),
	}
}

func (wfd *WFData) AppendUser(uid string) {
	wfd.UserIDs = append(wfd.UserIDs, uid)
	wfd.UserRole[uid] = wfd.CardDeck.PopCard()
}
