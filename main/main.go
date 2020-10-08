package main

import (
	"OneNightUltimateWerewolf/main/WF"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	MaxUser = 6
	prefix  = "ㄴ"
	eBin    = "🚮"
	eOne    = "1️⃣"
	eTwo    = "2️⃣"
	eThree  = "3️⃣"
	eFour   = "4️⃣"
	eFive   = "5️⃣"
	eSix    = "6️⃣"
	eSeven  = "7️⃣"
	eEight  = "8️⃣"
	eNine   = "9️⃣"
	eTen    = "🔟"
)

var (
	hwUserIDs []string
)

// Variables used for command line parameters
var (
	Token     string
	eNum      []string
	isGuildIn map[string]bool
	isUserIn  map[string]bool
	uidToGid  map[string]string
	wfDataMap map[string]*WF.WFData
)

func init() {
	hwUserIDs = make([]string, 0, 10)

	eNum = []string{eOne, eTwo, eThree, eFour, eFive, eSix, eSeven, eEight, eNine, eTen}

	uidToGid = make(map[string]string)
	isGuildIn = make(map[string]bool)
	isUserIn = make(map[string]bool)
	wfDataMap = make(map[string]*WF.WFData)
	flag.StringVar(&Token, "t", "NzYyNjUzOTczNjgwODgxNjg1.X3sS3A.Goy20AhNusZK4kGbLYJe1r8w1UA", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageReactionAdd)
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	_ = dg.Close()
}

func messageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}
	if !isUserIn[r.UserID] {
		return
	}

	gid := uidToGid[r.UserID]
	wfd := wfDataMap[gid]
	if wfd.CurStage == "Werewolf_only" {
		if wfd.UserRole[r.UserID] == "늑대인간" {
			if r.Emoji.Name == eOne {
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				uChan, _ := s.UserChannelCreate(r.UserID)
				go func() {
					time.Sleep(10 * time.Second)
					wfd.TimingChan <- true
				}()
				s.ChannelMessageSend(uChan.ID, "<1번: `"+wfd.CardDeck.Cards[0]+"` >")
			}
			if r.Emoji.Name == eTwo {
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				uChan, _ := s.UserChannelCreate(r.UserID)
				go func() {
					time.Sleep(10 * time.Second)
					wfd.TimingChan <- true
				}()
				s.ChannelMessageSend(uChan.ID, "<2번: `"+wfd.CardDeck.Cards[1]+"` >")
			}
			if r.Emoji.Name == eThree {
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				uChan, _ := s.UserChannelCreate(r.UserID)
				go func() {
					time.Sleep(10 * time.Second)
					wfd.TimingChan <- true
				}()
				s.ChannelMessageSend(uChan.ID, "<3번: `"+wfd.CardDeck.Cards[2]+"` >")
			}
		}
	}
	if wfd.CurStage == "Seer" {
		if wfd.UserRole[r.UserID] == "예언자" {
			if r.Emoji.Name == eBin {
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				wfd.CurStage = "Seer_trash"
				msg, _ := s.ChannelMessageSend(r.ChannelID, "보지 않고 덮어둘 카드를 고르시오\n< 1 > < 2 > < 3 >")
				s.MessageReactionAdd(r.ChannelID, msg.ID, eOne)
				s.MessageReactionAdd(r.ChannelID, msg.ID, eTwo)
				s.MessageReactionAdd(r.ChannelID, msg.ID, eThree)
			}
			for i := 0; i < len(wfd.UserIDs); i++ {
				if r.Emoji.Name == eNum[i] {
					if wfd.UserIDs[i] != r.UserID {
						s.ChannelMessageDelete(r.ChannelID, r.MessageID)
						wfd.CurStage = "Seer_used_power"
						powerMsg := strconv.Itoa(i+1) + "번째 유저는 `" + wfd.UserRole[wfd.UserIDs[i]] + "` 입니다."
						go func() {
							time.Sleep(5 * time.Second)
							wfd.TimingChan <- true
						}()
						s.ChannelMessageSend(r.ChannelID, powerMsg)
					}
				}
			}
		}
	}
	if wfd.CurStage == "Seer_trash" {
		if wfd.UserRole[r.UserID] == "예언자" {
			trashMsg := ""
			if r.Emoji.Name == eOne {
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				trashMsg += "<2번: " + wfd.CardDeck.Cards[1] + "> <3번: " + wfd.CardDeck.Cards[2] + ">"
				go func() {
					time.Sleep(5 * time.Second)
					wfd.TimingChan <- true
				}()
				s.ChannelMessageSend(r.ChannelID, trashMsg)
			}
			if r.Emoji.Name == eTwo {
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				trashMsg += "<1번: " + wfd.CardDeck.Cards[0] + "> <3번: " + wfd.CardDeck.Cards[2] + ">"
				go func() {
					time.Sleep(5 * time.Second)
					wfd.TimingChan <- true
				}()
				s.ChannelMessageSend(r.ChannelID, trashMsg)
			}
			if r.Emoji.Name == eThree {
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				trashMsg += "<1번: " + wfd.CardDeck.Cards[0] + "> <2번: " + wfd.CardDeck.Cards[1] + ">"
				go func() {
					time.Sleep(5 * time.Second)
					wfd.TimingChan <- true
				}()
				s.ChannelMessageSend(r.ChannelID, trashMsg)
			}
		}
	}
	if wfd.CurStage == "Robber" {
		if wfd.UserRole[r.UserID] == "강도" {
			robberMsg := ""
			for i := 0; i < len(wfd.UserIDs); i++ {
				if r.Emoji.Name == eNum[i] {
					if wfd.UserIDs[i] != r.UserID {
						s.ChannelMessageDelete(r.ChannelID, r.MessageID)
						user, _ := s.User(wfd.UserIDs[i])
						robberMsg = user.Username + "은 `" + wfd.FinalRole[wfd.UserIDs[i]] + "` 이었습니다. 하지만 이젠 아니죠."
						wfd.FinalRole[r.UserID] = wfd.FinalRole[wfd.UserIDs[i]]
						wfd.FinalRole[wfd.UserIDs[i]] = "강도"
						go func() {
							time.Sleep(5 * time.Second)
							wfd.TimingChan <- true
						}()
						s.ChannelMessageSend(r.ChannelID, robberMsg)
					}
				}
			}
		}
	}
	if wfd.CurStage == "TroubleMaker" {
		if wfd.UserRole[r.UserID] == "말썽쟁이" {
			tmMsg := ""
			for i := 0; i < len(wfd.UserIDs); i++ {
				if r.Emoji.Name == eNum[i] {
					if wfd.UserIDs[i] != r.UserID {
						wfd.CurStage = "TroubleMaker_oneMoreChoice"
						wfd.IndexChan <- i
						s.ChannelMessageDelete(r.ChannelID, r.MessageID)
						user, _ := s.User(wfd.UserIDs[i])
						selectMsg := "`" + user.Username + "`님을 선택하였습니다."
						s.ChannelMessageSend(r.ChannelID, selectMsg)
						index := len(wfd.UserIDs)
						for j := 0; j < len(wfd.UserIDs); j++ {
							if i == j {
								index = j
								break
							}
							if wfd.UserRole[wfd.UserIDs[j]] == "말썽쟁이" {
								tmMsg += "~~"
							}
							user, _ := s.User(wfd.UserIDs[j])
							tmMsg += "<" + strconv.Itoa(j+1) + "번 사용자: " + user.Username + "> "
							if wfd.UserRole[wfd.UserIDs[j]] == "말썽쟁이" {
								tmMsg += "~~"
							}
						}
						for j := index + 1; j < len(wfd.UserIDs); j++ {
							if wfd.UserRole[wfd.UserIDs[j]] == "말썽쟁이" {
								tmMsg += "~~"
							}
							user, _ := s.User(wfd.UserIDs[j])
							tmMsg += "<" + strconv.Itoa(j) + "번 사용자: " + user.Username + "> "
							if wfd.UserRole[wfd.UserIDs[j]] == "말썽쟁이" {
								tmMsg += "~~"
							}
						}
						msg, _ := s.ChannelMessageSend(r.ChannelID, tmMsg)
						for i := 0; i < len(wfd.UserIDs)-1; i++ {
							s.MessageReactionAdd(r.ChannelID, msg.ID, eNum[i])
						}
					}
				}
			}
		}
	}
	if wfd.CurStage == "TroubleMaker_oneMoreChoice" {
		if wfd.UserRole[r.UserID] == "말썽쟁이" {
			prev := <-wfd.IndexChan
			index := len(wfd.UserIDs)
			for i := 0; i < len(wfd.UserIDs); i++ {
				if i == prev {
					index = i
					break
				}
				if wfd.UserIDs[i] != r.UserID {
					if r.Emoji.Name == eNum[i] {
						s.ChannelMessageDelete(r.ChannelID, r.MessageID)
						user, _ := s.User(wfd.UserIDs[i])
						tmMsg := "`" + user.Username + "` 님을 선택하였습니다."
						s.ChannelMessageSend(r.ChannelID, tmMsg)
						temp := wfd.FinalRole[wfd.UserIDs[i]]
						wfd.FinalRole[wfd.UserIDs[i]] = wfd.FinalRole[wfd.UserIDs[prev]]
						wfd.FinalRole[wfd.UserIDs[prev]] = temp
						go func() {
							time.Sleep(5 * time.Second)
							wfd.TimingChan <- true
						}()
						s.ChannelMessageSend(r.ChannelID, "성공적으로 교환되었습니다.")
					}
				}
			}
			for i := index + 1; i < len(wfd.UserIDs); i++ {
				if wfd.UserIDs[i] != r.UserID {
					if r.Emoji.Name == eNum[i-1] {
						s.ChannelMessageDelete(r.ChannelID, r.MessageID)
						user, _ := s.User(wfd.UserIDs[i])
						tmMsg := "`" + user.Username + "` 님을 선택하였습니다."
						s.ChannelMessageSend(r.ChannelID, tmMsg)
						temp := wfd.FinalRole[wfd.UserIDs[i]]
						wfd.FinalRole[wfd.UserIDs[i]] = wfd.FinalRole[wfd.UserIDs[prev]]
						wfd.FinalRole[wfd.UserIDs[prev]] = temp
						go func() {
							time.Sleep(5 * time.Second)
							wfd.TimingChan <- true
						}()
						s.ChannelMessageSend(r.ChannelID, "성공적으로 교환되었습니다.")
					}
				}
			}
		}
	}
	if wfd.CurStage == "Election" {
		for _, item := range wfd.UserIDs {
			go func(uid string) {
				if uid == r.UserID {
					for i := 0; i < len(wfd.UserIDs); i++ {
						if r.Emoji.Name == eNum[i] {
							s.ChannelMessageDelete(r.ChannelID, r.MessageID)
							selCandi := wfd.UserIDs[i]
							voteUser, _ := s.User(uid)
							selUser, _ := s.User(selCandi)
							s.ChannelMessageSend(r.ChannelID, "`"+
								selUser.Username+"` 님에게 투표하였습니다.")
							wfd.ElectChan <- WF.NewElectInfo(voteUser.Username, selCandi)
						}
					}
				}
			}(item)
		}
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, prefix) { // 프리픽스로 시작하는 메시지일 경우
		var wfd *WF.WFData

		if m.Content == prefix+"강제종료" {
			wfd = wfDataMap[m.GuildID]
			s.ChannelMessageSend(wfd.UseChannelID, "안전하게 강제종료 수행중..")
			time.Sleep(time.Second * 5)
			cancelGameTask(m)
			s.ChannelMessageSend(wfd.UseChannelID, "사용 종료가 정상적으로 완료되었습니다.")
		}

		if m.Content == prefix+"시작" && !isGuildIn[m.GuildID] {
			isGuildIn[m.GuildID] = true
			wfDataMap[m.GuildID] = WF.NewWFData(m.Author.ID, m.ChannelID)
			wfDataMap[m.GuildID].CurStage = "Prepare"
			newUserTask(m)
			_, _ = s.ChannelMessageSend(m.ChannelID, "게임 시작!\n`"+prefix+"입장` 으로 입장하세요")
		}
		if isGuildIn[m.GuildID] {
			wfd = wfDataMap[m.GuildID]
			if m.Content == prefix+"입장" && wfd.CurStage == "Prepare" {
				if isUserIn[m.Author.ID] {
					s.ChannelMessageSend(m.ChannelID, "이미 입장한 유저입니다.")
					return
				}
				newUserTask(m)
				s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+"님이 입장하셨습니다.")
			}
			if m.Author.ID == wfd.AdminUserID && m.Content == prefix+"취소" && wfd.CurStage == "Prepare" {
				cancelGameTask(m)
				s.ChannelMessageSend(m.ChannelID, "게임이 취소되었습니다.")
			}
			if m.Author.ID == wfd.AdminUserID && strings.HasPrefix(m.Content, prefix+"더미추가") && wfd.CurStage == "Prepare" {
				sepMsg := strings.Split(m.Content, " ")
				if len(sepMsg) == 1 {
					s.ChannelMessageSend(m.ChannelID, "추가할 인원 숫자를 입력하세요.")
					return
				}
				num, err := strconv.Atoi(sepMsg[1])
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "숫자가 아닌걸 입력했네요....")
				}
				for i := 0; i < num; i++ {
					newUserTask(m)
				}
				s.ChannelMessageSend(m.ChannelID, "현재인원: ("+strconv.Itoa(len(wfd.UserIDs))+"/6)")
			}
			if strings.HasPrefix(m.Content, prefix+"마감") && wfd.CurStage == "Prepare" {
				if len(wfd.UserIDs) == 6 {
					wfd.CurStage = "Prepare_finishing"
					for _, item := range wfd.UserIDs {
						go func(uid string) {
							uChan, _ := s.UserChannelCreate(uid)
							user, _ := s.User(uid)
							wfd.GameLog += "> `" + user.Username + "` 님의 역할이 `" + wfd.UserRole[uid] + "` 로 배정되었습니다.\n"
							roleBrief := "> **당신의 역할은 **`" + wfd.UserRole[uid] + "`**입니다.**\n"
							roleBrief += getRoleInfo(wfd.UserRole[uid])
							s.ChannelMessageSend(uChan.ID, roleBrief)
						}(item)
					}
					wfd.CurStage = "Werewolf"
					werewolfTask(s, wfd)
					minionTask(s, wfd)
					seerTask(s, wfd)
					robberTask(s, wfd)
					tmTask(s, wfd)
					insomniacTask(s, wfd)
					dayBriefTask(s, wfd)
				} else {
					s.ChannelMessageSend(m.ChannelID, "정확한 인원이 모이지 않았습니다. ("+strconv.Itoa(len(wfd.UserIDs))+"/"+strconv.Itoa(MaxUser)+")")
				}
			}
		}
	}

	homeworkMethod(s, m)
}

func dayBriefTask(s *discordgo.Session, wfd *WF.WFData) {
	briefMsg := ""

	briefMsg += "> 모든 특수 능력 사용이 끝났습니다." +
		"\n> 3분 후 여러분들에게 각자의 투표 용지가 전송됩니다." +
		"\n> 한번 투표한 내용은 바꿀 수 없기에, 신중하게 투표하세요" +
		"\n"
	go func() {
		time.Sleep(time.Minute * 3)
		wfd.TimingChan <- true
	}()
	s.ChannelMessageSend(wfd.UseChannelID, briefMsg)
	userChans := make([]string, 0, 10)
	userNames := make([]string, 0, 10)
	for _, item := range wfd.UserIDs {
		user, _ := s.User(item)
		uChan, _ := s.UserChannelCreate(item)
		userChans = append(userChans, uChan.ID)
		userNames = append(userNames, user.Username)
	}
	briefMsg = ""
	for i, item := range userNames {
		briefMsg += "<" + strconv.Itoa(i+1) + "번: " + item + "> "
	}
	<-wfd.TimingChan
	wfd.CurStage = "Election"
	for _, cid := range userChans {
		go func(item string) {
			msg, _ := s.ChannelMessageSend(item, briefMsg)
			for i := 0; i < len(wfd.UserIDs); i++ {
				s.MessageReactionAdd(item, msg.ID, eNum[i])
			}
		}(cid)
	}
	electFinishTask(s, wfd)
}

func electFinishTask(s *discordgo.Session, wfd *WF.WFData) {
	electData := make([]*WF.ElectInfo, 0, 10)
	electResult := make([]int, len(wfd.UserIDs))
	s.ChannelMessageSend(wfd.UseChannelID, "> 투표를 시작합니다!")
	for i := 0; i < len(wfd.UserIDs); i++ {
		electData = append(electData, <-wfd.ElectChan)
		electAlarmMsg := "`" + electData[i].VoterName + "`님이 투표하셨습니다."
		s.ChannelMessageSend(wfd.UseChannelID, electAlarmMsg)
	}
	s.ChannelMessageSend(wfd.UseChannelID, "> 투표가 끝났습니다.")
	s.ChannelMessageSend(wfd.UseChannelID, "결과 계산중...")
	for i, uid := range wfd.UserIDs {
		for _, elc := range electData {
			if uid == elc.CandiID {
				electResult[i]++
			}
		}
	}
	electResultMsg := "> 투표 결과 :\n"
	max := 0
	maxi := -1
	maxName := ""
	for i, item := range electResult {
		if max < item {
			max = item
			maxi = i
		}
		if item != 0 {
			user, _ := s.User(wfd.UserIDs[i])
			if maxi == i {
				maxName = user.Username
			}
			electResultMsg += "<`" + user.Username + "` : " + strconv.Itoa(item) + "표>\n"
		}
	}
	electResultMsg += "> `" + maxName + "` 님이 총 " + strconv.Itoa(electResult[maxi]) + " 표로 처형당하였습니다."
	s.ChannelMessageSend(wfd.UseChannelID, electResultMsg)
	s.ChannelMessageSend(wfd.UseChannelID, "`"+maxName+"` 님의 직업은..?")
	for i := 0; i < 3; i++ {
		s.ChannelMessageSend(wfd.UseChannelID, "...")
		time.Sleep(time.Second)
	}
	s.ChannelMessageSend(wfd.UseChannelID, "`"+wfd.UserRole[wfd.UserIDs[maxi]]+
		"` -> `"+wfd.FinalRole[wfd.UserIDs[maxi]]+"` 입니다.")
}

func werewolfTask(s *discordgo.Session, wfd *WF.WFData) {
	s.ChannelMessageSend(wfd.UseChannelID, "늑대인간의 차례입니다.")
	wfd.GameLog += "> 늑대인간의 차례"
	wolvesID := make([]string, 0, 10)
	for _, item := range wfd.UserIDs {
		if wfd.UserRole[item] == "늑대인간" {
			wolvesID = append(wolvesID, item)
		}
	}
	if len(wolvesID) != 1 {
		go func() {
			time.Sleep(40 * time.Second)
			wfd.CurStage = "Minion"
			wfd.TimingChan <- true
		}()
	}

	if len(wolvesID) == 1 {
		wfd.CurStage = "Werewolf_only"

		wolvesMsg := "20초 안에 세 장의 비공개 카드 중 한 장을 선택하세요"
		wolvesMsg += "< 1 > < 2 > < 3 >"
		uChan, _ := s.UserChannelCreate(wolvesID[0])
		msg, _ := s.ChannelMessageSend(uChan.ID, wolvesMsg)
		s.MessageReactionAdd(uChan.ID, msg.ID, eOne)
		s.MessageReactionAdd(uChan.ID, msg.ID, eTwo)
		s.MessageReactionAdd(uChan.ID, msg.ID, eThree)
		<-wfd.TimingChan
		s.ChannelMessageSend(uChan.ID, "당신의 차례가 끝났습니다.")
		s.ChannelMessageSend(wfd.UseChannelID, "늑대인간의 차례 종료")
		return
	}
	for _, item := range wolvesID {
		uChan, _ := s.UserChannelCreate(item)
		wolvesMsg := "늑대인간: "
		for _, item := range wolvesID {
			user, _ := s.User(item)
			wolvesMsg += "<" + user.Username + "> "
		}
		s.ChannelMessageSend(uChan.ID, wolvesMsg)

	}
	<-wfd.TimingChan
	s.ChannelMessageSend(wfd.UseChannelID, "늑대인간의 차례 종료.")
}

func minionTask(s *discordgo.Session, wfd *WF.WFData) {
	s.ChannelMessageSend(wfd.UseChannelID, "하수인의 차례입니다.")

	wolvesID := make([]string, 0, 10)
	minionID := ""
	for _, item := range wfd.UserIDs {
		if wfd.UserRole[item] == "늑대인간" {
			wolvesID = append(wolvesID, item)
		}
		if wfd.UserRole[item] == "하수인" {
			minionID = item
		}
	}
	go func() {
		time.Sleep(time.Second * 10)
		wfd.CurStage = "Seer"
		wfd.TimingChan <- true
	}()
	minionMsg := "늑대인간은, "
	for _, item := range wolvesID {
		user, _ := s.User(item)
		minionMsg += "<" + user.Username + "> "
	}
	minionMsg += "입니다."

	if len(wolvesID) == 0 {
		minionMsg = "늑대인간이 존재하지 않습니다."
	}

	user, err := s.User(minionID)
	if err != nil {
		<-wfd.TimingChan
		s.ChannelMessageSend(wfd.UseChannelID, "하수인의 차례 종료.")
		return
	}

	uChan, _ := s.UserChannelCreate(user.ID)
	s.ChannelMessageSend(uChan.ID, minionMsg)
	<-wfd.TimingChan
	s.ChannelMessageSend(wfd.UseChannelID, "하수인의 차례 종료.")
}

func seerTask(s *discordgo.Session, wfd *WF.WFData) {
	s.ChannelMessageSend(wfd.UseChannelID, "예언자의 차례입니다.")
	seerID := ""
	seerMsg := "30초 안에 버려진 카드중 2장 또는, 확인하고싶은 사람 한 명을 선택하세요\n자신은 선택할 수 없어요\t(" + eBin + "): 버려진 카드에서 고르기\n"
	for i, item := range wfd.UserIDs {
		if wfd.UserRole[item] == "예언자" {
			seerID = item
			seerMsg += "~~"
		}
		user, _ := s.User(item)
		seerMsg += "<" + strconv.Itoa(i+1) + "번 사용자: " + user.Username + ">\t"
		if wfd.UserRole[item] == "예언자" {
			seerMsg += "~~"
		}
	}
	if seerID == "" {
		go func() {
			time.Sleep(40 * time.Second)
			wfd.CurStage = "Robber"
			wfd.TimingChan <- true
		}()
	}
	if seerID == "" {
		<-wfd.TimingChan
		s.ChannelMessageSend(wfd.UseChannelID, "예언자의 차례 종료.")
		return
	}
	uChan, _ := s.UserChannelCreate(seerID)
	msg, _ := s.ChannelMessageSend(uChan.ID, seerMsg)
	s.MessageReactionAdd(uChan.ID, msg.ID, eBin)
	for i := 0; i < len(wfd.UserIDs); i++ {
		s.MessageReactionAdd(uChan.ID, msg.ID, eNum[i])
	}

	<-wfd.TimingChan
	s.ChannelMessageSend(wfd.UseChannelID, "예언자의 차례 종료.")
}

func robberTask(s *discordgo.Session, wfd *WF.WFData) {
	s.ChannelMessageSend(wfd.UseChannelID, "강도의 차례입니다.")

	robberID := ""
	robberMsg := "자신은 선택할 수 없어요\n"
	for i, item := range wfd.UserIDs {
		if wfd.UserRole[item] == "강도" {
			robberID = item
			robberMsg += "~~"
		}
		user, _ := s.User(item)
		robberMsg += "<" + strconv.Itoa(i+1) + "번 사용자: " + user.Username + ">\t"
		if wfd.UserRole[item] == "강도" {
			robberMsg += "~~"
		}
	}
	if robberID == "" {
		go func() {
			time.Sleep(30 * time.Second)
			wfd.CurStage = "TroubleMaker"
			wfd.TimingChan <- true
		}()
	}
	if robberID == "" {
		<-wfd.TimingChan
		wfd.CurStage = "TroubleMaker"
		s.ChannelMessageSend(wfd.UseChannelID, "강도의 차례 종료.")
		return
	}
	uChan, _ := s.UserChannelCreate(robberID)
	msg, _ := s.ChannelMessageSend(uChan.ID, robberMsg)
	for i := 0; i < len(wfd.UserIDs); i++ {
		s.MessageReactionAdd(uChan.ID, msg.ID, eNum[i])
	}

	<-wfd.TimingChan
	wfd.CurStage = "TroubleMaker"
	s.ChannelMessageSend(wfd.UseChannelID, "강도의 차례 종료.")

}

func tmTask(s *discordgo.Session, wfd *WF.WFData) {
	s.ChannelMessageSend(wfd.UseChannelID, "말썽쟁이의 차례입니다.")

	tmID := ""
	tmMsg := "자신은 선택할 수 없어요\n"
	for i, item := range wfd.UserIDs {
		if wfd.UserRole[item] == "말썽쟁이" {
			tmID = item
			tmMsg += "~~"
		}
		user, _ := s.User(item)
		tmMsg += "<" + strconv.Itoa(i+1) + "번 사용자: " + user.Username + ">\t"
		if wfd.UserRole[item] == "말썽쟁이" {
			tmMsg += "~~"
		}
	}
	if tmID == "" {
		go func() {
			time.Sleep(60 * time.Second)
			wfd.CurStage = "Insomniac"
			wfd.TimingChan <- true
		}()
	}
	if tmID == "" {
		<-wfd.TimingChan
		s.ChannelMessageSend(wfd.UseChannelID, "말썽쟁이의 차례 종료.")
		return
	}
	uChan, _ := s.UserChannelCreate(tmID)
	msg, _ := s.ChannelMessageSend(uChan.ID, tmMsg)
	for i := 0; i < len(wfd.UserIDs); i++ {
		s.MessageReactionAdd(uChan.ID, msg.ID, eNum[i])
	}

	<-wfd.TimingChan
	s.ChannelMessageSend(wfd.UseChannelID, "말썽쟁이의 차례 종료.")

}

func insomniacTask(s *discordgo.Session, wfd *WF.WFData) {
	s.ChannelMessageSend(wfd.UseChannelID, "불면증환자의 차례입니다.")

	go func() {
		time.Sleep(10 * time.Second)
		wfd.CurStage = "Day"
		wfd.TimingChan <- true
	}()

	inID := ""
	for _, item := range wfd.UserIDs {
		if wfd.UserRole[item] == "불면증환자" {
			inID = item
		}
	}

	if inID == "" {
		<-wfd.TimingChan
		s.ChannelMessageSend(wfd.UseChannelID, "불면증환자의 차례 종료.")
		return
	}

	uChan, _ := s.UserChannelCreate(inID)
	s.ChannelMessageSend(uChan.ID, "다른 모든 사람이 능력을 쓴 후, 당신의 역할은 다음과 같습니다."+
		"\n역할: "+wfd.FinalRole[inID])

	<-wfd.TimingChan
	s.ChannelMessageSend(wfd.UseChannelID, "불면증환자의 차례 종료.")

}

func getRoleInfo(role string) string {
	info := ""

	if role == "늑대인간" {
		info = "당신은 피에 굶주린 늑대인간입니다." +
			"\n당신은 게임이 시작된 후에 동료 늑대인간을 확인할 수 있습니다." +
			"\n만약 동료 늑대인간이 없다면, 아무에게도 배정되지 않은 역할 카드들 3장 중 1장을" +
			"\n무작위로 확인할 수 있습니다." +
			"\n당신을 도와줄 하수인 한명이 있을수도, 없을수도 있습니다." +
			"\n승리를 위해서는 당신과 당신의 동료 늑대인간 중 아무도 죽어서는 안됩니다." +
			"\n또한 무두장이가 자살에 성공한다면, 당신은 패배하게 됩니다." +
			"\n**보름달이 뜨기 전 마을 사람들을 혼란에 빠뜨리세요. 행운을 빕니다!**"
	}
	if role == "무두장이" {
		info = "당신은 일에 지쳐 극단적 선택을 꿈꾸는 무두장이입니다." +
			"\n당신은 마을 시민들이 늑대인간에게 죽든 말든 아무런 관심이 없습니다." +
			"\n왜냐하면 당신은 더 이상 희망이 보이지 않은 삶을 이어가고 있으니까요.." +
			"\n당신은 투표로 당신이 처형당하길 바라고 있습니다." +
			"\n사람들이 당신을 늑대인간이라고 믿게 하세요. 당신이 죽는다면, 당신의 승리입니다." +
			"\n당신이 죽으면 마을 사람들과 늑대인간들은 패배합니다." +
			"\n**당신의 불운을 여기서 멈추고 싶다면, 이제 영원한 잠에 빠질 때입니다.. 안타깝지만, 행운을 빕니다...**"
	}
	if role == "마을주민" {
		info = "당신은 아무런 능력도 가지지 못했습니다." +
			"\n불안과 공포속에서 늑대인간을 찾아서 처형하세요"
	}
	if role == "하수인" {
		info = "당신은 늑대인간들을 위해 목숨바칠 각오가 되어 있는 하수인입니다." +
			"\n당신은 누가 늑대인간인지 잘 알고 있습니다." +
			"\n하지만 늑대인간들은 당신이 존재하는지조차도 알지 못해요.." +
			"\n당신은 당신이 처형당하는 한이 있더라도 늑대인간을 지켜야 합니다." +
			"\n아무도 늑대인간이 아닌 척 하는 것도 좋은 방법일 겁니다." +
			"\n무두장이가 자살하는것을 막는 것 또한 당신의 임무입니다." +
			"\n**내일은 보름달이네요.. 행운을 빕니다!**"
	}
	if role == "예언자" {
		info = "당신은 장막을 들추고 미래를 엿보았지만 그곳엔 오직... 보름달 뿐이었습니다." +
			"\n당신은 아무에게도 배정되지 않은 역할 3장 중 2장을 보거나," +
			"\n원하는 사람 한명의 역할을 간파할 수 있습니다." +
			"\n늑대인간들은 예언자행세를 하며 마을시민들을 혼란스럽게 할 수도 있습니다." +
			"\n승리를 위해선 진짜 늑대인간을 찾아 처형시키세요" +
			"\n늑대인간들과 하수인이 아무도 없다면, 아무도 처형시키면 안됩니다." +
			"\n**많은걸 알고 있기에, 더 의심받을 수 있어요.. 행운을 빕니다..**"
	}
	if role == "강도" {
		info = "당신은 당신에게 특별한 힘이 생긴걸 알 수 있었습니다." +
			"\n다른 사람의 능력까지 훔쳐올 수 있는 능력이죠" +
			"\n원하는 사람의 능력을 알아낼 수 있습니다." +
			"\n그리고 그 사람과 능력을 바꿔치기합니다." +
			"\n늑대인간의 능력을 훔쳤다면, 늑대인간을 지켜야합니다." +
			"\n훔친 능력이 마땅치 않다면, 늑대인간을 처형시켜야 해요." +
			"\n**당신의 운을 시험해 보세요, 행운을 빕니다.**"
	}
	if role == "말썽쟁이" {
		info = "당신 또 무슨짓을 한 거죠?" +
			"\n이런! 두 사람의 능력을 바꿔버리다니!" +
			"\n어떤 능력인진 최소한 알고 바꿔야죠... 정말 이름값 하는군요" +
			"\n어휴... 그래도 우린 늑대인간을 처형시켜야 해요." +
			"\n**다음엔 말썽 안피우기로 약속합시다. 행운을 빌어요!**"
	}
	if role == "불면증환자" {
		info = "오늘 밤은 잠이 올까요..." +
			"\n새벽에 깨어나 잠을 설치진 않을까요..." +
			"\n당신은 다른 사람들이 밤사이 했던 모든 일들이 당신의 역할을 어떻게 바꿔놓았는지" +
			"\n아니면 그저 아무것도 변한게 없는지 알 수 있어요." +
			"\n당신이 늑대인간이 되었다면, 늑대인간을 지켜야죠" +
			"\n그게 아니라면.. 늑대인간을 꼭 찾아내세요!" +
			"\n**오늘 밤은 두 다리 쭉 뻗고 단잠에 빠질 수 있기를.. 행운을 빌어요!**"
	}

	return info
}

// 새로운 유저 등록시 수행
func newUserTask(m *discordgo.MessageCreate) {
	wfd := wfDataMap[m.GuildID]
	if len(wfd.UserIDs) >= MaxUser {
		return
	}
	isUserIn[m.Author.ID] = true
	uidToGid[m.Author.ID] = m.GuildID
	wfd.AppendUser(m.Author.ID)
}

// 게임 취소시 데이터 삭제 수행
func cancelGameTask(m *discordgo.MessageCreate) {
	wfd := wfDataMap[m.GuildID]
	for _, item := range wfd.UserIDs {
		uidToGid[item] = ""
		isUserIn[item] = false
	}
	wfd.UserIDs = make([]string, 0, MaxUser)
	isGuildIn[m.GuildID] = false
}

func homeworkMethod(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Content, "=info") {
		if len(m.Mentions) == 0 {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Your message doesn't have any mentions!")
			return
		}
		msg := ""
		for i, item := range m.Mentions {
			msg = strconv.Itoa(i+1) + "번째 멘션 정보\n" +
				"UID:\t" + item.ID + "\n" +
				"Username:\t" + item.Username + "\n" +
				"Mention:\t" + item.Mention() + "\n"
			_, _ = s.ChannelMessageSend(m.ChannelID, msg)
		}
	}
	if strings.HasPrefix(m.Content, "=save") {
		if len(m.Mentions) == 0 {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Your message doesn't have any mentions!")
			return
		}
		if len(m.Mentions) > 3 {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Your message contains too many mentions!")
			return
		}
		msg := ""
		for i, item := range m.Mentions {
			hwUserIDs = append(hwUserIDs, item.ID)
			msg = "총 " + strconv.Itoa(i+1) + "명의 사용자를 저장하였습니다."
		}
		_, _ = s.ChannelMessageSend(m.ChannelID, msg)
	}
	if m.Content == "=load" {
		if len(hwUserIDs) == 0 {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Your server doesn't have any user ID")
			return
		}
		msg := ""
		for i, item := range hwUserIDs {
			user, _ := s.User(item)
			msg = strconv.Itoa(i+1) + "번째 저장된 정보\n" +
				"UID:\t" + user.ID + "\n" +
				"Username:\t" + user.Username + "\n" +
				"Mention:\t" + user.Mention() + "\n"
			_, _ = s.ChannelMessageSend(m.ChannelID, msg)
		}
	}
	if m.Content == "=delete" {
		if len(hwUserIDs) == 0 {
			_, _ = s.ChannelMessageSend(m.ChannelID, "저장된 정보가 없었습니다.")
			return
		}
		hwUserIDs = make([]string, 0, 10)
		msg := "모든 입력된 데이터들을 삭제하였습니다."
		_, _ = s.ChannelMessageSend(m.ChannelID, msg)
	}
}
