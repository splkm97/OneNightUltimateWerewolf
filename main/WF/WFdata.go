package WF

const (
	MaxUser = 6
)

type ElectInfo struct {
	VoterName string
	CandiID   string
}

func NewElectInfo(voterName, candiID string) *ElectInfo {
	return &ElectInfo{VoterName: voterName, CandiID: candiID}
}

type SettingData struct {
	CardDeck *CardDeck
	MaxUser  int
}

func NewSettingData(cardDeck CardDeck, maxUser int) *SettingData {
	newCardDeck := NewCardDeck()
	for _, item := range cardDeck.Cards {
		newCardDeck.Cards = append(newCardDeck.Cards, item)
	}
	return &SettingData{CardDeck: newCardDeck, MaxUser: maxUser}
}

type Data struct {
	CardDeck     CardDeck
	MaxUser      int
	DIFlag       bool
	GameLog      string
	IndexChan    chan int
	TimingChan   chan bool
	DoppelChan   chan bool
	ElectChan    chan *ElectInfo
	CurStage     string
	AdminUserID  string
	UseChannelID string
	UserIDs      []string
	UserRole     map[string]string
	FinalRole    map[string]string
	DeadUserMap  map[string]bool
}

func NewWFData(uid, cid string) *Data {
	nc := NewCardDeck()
	nc.ShuffleCards()
	return &Data{
		CardDeck:     *nc,
		GameLog:      "",
		DIFlag:       false,
		TimingChan:   make(chan bool, 2),
		DoppelChan:   make(chan bool, 2),
		IndexChan:    make(chan int, 10),
		ElectChan:    make(chan *ElectInfo, 10),
		CurStage:     "Prepare",
		AdminUserID:  uid,
		UseChannelID: cid,
		UserRole:     make(map[string]string),
		FinalRole:    make(map[string]string),
		UserIDs:      make([]string, 0, MaxUser),
		DeadUserMap:  make(map[string]bool),
	}
}

func (wfd *Data) AppendUser(uid string) {
	wfd.UserIDs = append(wfd.UserIDs, uid)
	wfd.UserRole[uid] = wfd.CardDeck.PopCard()
	wfd.FinalRole[uid] = wfd.UserRole[uid]
}
