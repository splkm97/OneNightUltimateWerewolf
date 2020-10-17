package main

import (
	"OneNightUltimateWerewolf/main/WF"
	"flag"
	"fmt"
	embed "github.com/clinet/discordgo-embed"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	//	prefix = "ㅌ"	// 테스트 봇 프리픽스
	prefix = "ㅁ" // 본계 봇 프리픽스
	eBin   = "🚮"
	eOne   = "1️⃣"
	eTwo   = "2️⃣"
	eThree = "3️⃣"
	eFour  = "4️⃣"
	eFive  = "5️⃣"
	eSix   = "6️⃣"
	eSeven = "7️⃣"
	eEight = "8️⃣"
	eNine  = "9️⃣"
	eTen   = "🔟"
)

var (
	hwUserIDs []string
)

// Variables used for command line parameters
var (
	commandMsg     string
	helpMsg        string
	wereWinTitle   string
	wereWinMsg     string
	vilWinTitle    string
	vilWinMsg      string
	tannerWinTitle string
	psMsg          string
	namuDeck       = []string{
		"늑대인간",
		"늑대인간",
		"예언자",
		"강도",
		"말썽쟁이",
		"마을주민",
		"마을주민",
		"마을주민",
	}
	doppelBugDeck = []string{
		"늑대인간",
		"늑대인간",
		"도플갱어",
		"예언자",
		"강도",
		"말썽쟁이",
		"주정뱅이",
		"마을주민",
		"무두장이",
	}
	anywayDeck = []string{
		"늑대인간",
		"늑대인간",
		"도플갱어",
		"하수인",
		"예언자",
		"강도",
		"말썽쟁이",
		"사냥꾼",
		"무두장이",
		"불면증환자",
	}
	comeOnSeerDeck = []string{
		"늑대인간",
		"늑대인간",
		"하수인",
		"예언자",
		"강도",
		"말썽쟁이",
		"사냥꾼",
		"마을주민",
		"마을주민",
		"무두장이",
		"불면증환자",
	}
	chaosDeck = []string{
		"늑대인간",
		"늑대인간",
		"도플갱어",
		"프리메이슨",
		"프리메이슨",
		"하수인",
		"예언자",
		"강도",
		"말썽쟁이",
		"주정뱅이",
		"무두장이",
		"마을주민",
		"사냥꾼",
		"불면증환자",
	}
	classList = []string{
		"도플갱어",
		"늑대인간",
		"하수인",
		"프리메이슨",
		"예언자",
		"강도",
		"말썽쟁이",
		"주정뱅이",
		"무두장이",
		"마을주민",
		"사냥꾼",
		"불면증환자",
	}
	gameSeq = []string{
		"도플갱어",
		"늑대인간",
		"하수인",
		"프리메이슨",
		"예언자",
		"강도",
		"말썽쟁이",
		"주정뱅이",
		"불면증환자",
	}
	isTest         bool
	Token          string
	cardMap        map[string][]string
	eNum           []string
	isGuildIn      map[string]bool
	isUserIn       map[string]bool
	uidToGid       map[string]string
	prevSettingMap map[string]*WF.SettingData
	wfDataMap      map[string]*WF.Data
)

func init() {
	hwUserIDs = make([]string, 0, 10)

	eNum = []string{eOne, eTwo, eThree, eFour, eFive, eSix, eSeven, eEight, eNine, eTen}

	dat, err := ioutil.ReadFile("./Asset/help.txt")
	if err != nil {
		panic(err)
	}
	helpMsg = string(dat)
	dat, err = ioutil.ReadFile("./Asset/Commands.txt")
	if err != nil {
		panic(err)
	}
	commandMsg = string(dat)
	wereWinTitle = "**늑대인간 팀 승리조건**"
	vilWinTitle = "**마을주민 팀 승리조건**"
	dat, err = ioutil.ReadFile("./Asset/WereWin.txt")
	if err != nil {
		panic(err)
	}
	wereWinMsg = string(dat)
	dat, err = ioutil.ReadFile("./Asset/VilWin.txt")
	if err != nil {
		panic(err)
	}
	vilWinMsg = string(dat)
	dat, err = ioutil.ReadFile("./Asset/TannerWin.txt")
	if err != nil {
		panic(err)
	}
	tannerWinTitle = string(dat)
	dat, err = ioutil.ReadFile("./Asset/ps.txt")
	if err != nil {
		panic(err)
	}
	psMsg = string(dat)
	for _, item := range classList {
		psMsg += item + " "
	}
	isTest = false
	cardMap = make(map[string][]string)
	prevSettingMap = make(map[string]*WF.SettingData)
	uidToGid = make(map[string]string)
	isGuildIn = make(map[string]bool)
	isUserIn = make(map[string]bool)
	wfDataMap = make(map[string]*WF.Data)
	//	flag.StringVar(&Token, "t", "NzY1NDUxMDA3MDE4MjcwNzMw.X4U_zQ._U1RlF8BtOvQzYnDrv7RpInDr44", "Bot Token")	// 테스트 봇 토큰
	flag.StringVar(&Token, "t", "NzYyNjUzOTczNjgwODgxNjg1.X3sS3A.Goy20AhNusZK4kGbLYJe1r8w1UA", "Bot Token") // 본계 봇 토큰

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

func exit(stage string) bool {
	if stage == "Exit" {
		return true
	} else {
		return false
	}
}

func nextStage(wfd *WF.Data) {
	if exit(wfd.CurStage) {
		return
	}
	for i, item := range gameSeq {
		if strings.HasPrefix(wfd.CurStage, item) {
			if len(gameSeq) == i+1 {
				wfd.CurStage = "Day"
				break
			} else {
				wfd.CurStage = gameSeq[i+1]
				break
			}
		}
	}
	if wfd.CurStage == "Prepare_finishing" {
		wfd.CurStage = gameSeq[0]
	}
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

	if r.UserID == wfd.AdminUserID && wfd.CurStage == "recomChoice" {
		for i := 0; i < 5; i++ {
			if r.Emoji.Name == eNum[i] { //ddddd
				wfd.CardDeck.ChoiceChan <- i
			}
		}
	}

	if wfd.CurStage == "도플갱어_강도" {
		if wfd.UserRole[r.UserID] == "도플갱어" {
			robberReactionTask(s, r, wfd)
		}
	}
	if wfd.CurStage == "도플갱어_주정뱅이" {
		if wfd.UserRole[r.UserID] == "도플갱어" {
			drunkReactionTask(s, r, wfd)
		}
	}
	if wfd.CurStage == "도플갱어_예언자_trash" {
		if wfd.UserRole[r.UserID] == "도플갱어" {
			seerTrashReactionTask(s, r, wfd)
		}
	}
	if wfd.CurStage == "도플갱어_예언자" {
		if wfd.UserRole[r.UserID] == "도플갱어" {
			if r.Emoji.Name == eBin {
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				wfd.CurStage = "도플갱어_예언자_trash"
				msg, _ := s.ChannelMessageSendEmbed(r.ChannelID, embed.NewGenericEmbed("보지 않을 직업을 고르세요.\n< 1 > < 2 > < 3 >", ""))
				s.MessageReactionAdd(r.ChannelID, msg.ID, eOne)
				s.MessageReactionAdd(r.ChannelID, msg.ID, eTwo)
				s.MessageReactionAdd(r.ChannelID, msg.ID, eThree)
			}
			seerUserReactionTask(s, r, wfd)
		}
	}
	if wfd.CurStage == "도플갱어_말썽쟁이_oneMoreChoice" {
		if wfd.UserRole[r.UserID] == "도플갱어" {
			tmOneMoreReactionTask(s, r, wfd)
		}
	}
	if wfd.CurStage == "도플갱어_말썽쟁이" {
		if wfd.UserRole[r.UserID] == "도플갱어" {
			tmReactionTask(s, r, wfd)
		}
	}
	if wfd.CurStage == "도플갱어" {
		if wfd.UserRole[r.UserID] == "도플갱어" {
			s.ChannelMessageDelete(r.ChannelID, r.MessageID)
			uChan, _ := s.UserChannelCreate(r.UserID)
			dUser, _ := s.User(r.UserID)
			copyRole := ""
			copyUserID := ""
			dMsg := ""
			for i, item := range eNum {
				if item == r.Emoji.Name {
					copyUserID = wfd.UserIDs[i]
					copyRole = wfd.UserRole[copyUserID]
					if copyRole == "도플갱어" {
						return
					}
					copyUser, _ := s.User(copyUserID)
					wfd.GameLog += "\n도플갱어 `" + dUser.Username + "` 는 " +
						"`" + copyUser.Username + "` 의 직업\n`" + copyRole +
						"` (을)를 복사하였습니다."
					dMsg += "당신은 `" + copyUser.Username + "` 의 직업 `" +
						copyRole + "` (을)를 복사하였습니다."
					wfd.FinalRole[r.UserID] = copyRole
				}
			}
			s.ChannelMessageSendEmbed(uChan.ID, embed.NewGenericEmbed(dMsg, ""))
			sendClassInfo(s, copyRole, uChan.ID)
			if copyRole == "강도" {
				wfd.DoppelRobberFlag = true
				sendAllUserAddReaction(s, wfd, "강도", uChan)
				wfd.CurStage = "도플갱어_강도"
			} else if copyRole == "주정뱅이" {
				wfd.DoppelDrunkFlag = true
				sendDiscardsAddReaction(s, uChan)
				wfd.CurStage = "도플갱어_주정뱅이"
			} else if copyRole == "예언자" {
				doppelMsg := "버려진 직업들 중 2개 또는, 직업을 확인하고싶은 사람 한 명을 선택하세요" +
					"\n자신은 선택할 수 없어요\t(" + eBin + "): 버려진 직업들에서 고르기\n"
				for i, item := range wfd.UserIDs {
					user, _ := s.User(item)
					if item == r.UserID {
						doppelMsg += "~~"
					}
					doppelMsg += "<" + strconv.Itoa(i+1) + "번 사용자 : `" +
						user.Username + "`>\t"
					if item == r.UserID {
						doppelMsg += "~~"
					}
				}
				msg, _ := s.ChannelMessageSendEmbed(uChan.ID, embed.NewGenericEmbed(doppelMsg, ""))
				s.MessageReactionAdd(uChan.ID, msg.ID, eBin)
				for i := 0; i < len(wfd.UserIDs); i++ {
					s.MessageReactionAdd(uChan.ID, msg.ID, eNum[i])
				}
				wfd.CurStage = "도플갱어_예언자"
			} else if copyRole == "말썽쟁이" {
				sendAllUserAddReaction(s, wfd, "말썽쟁이", uChan)
				wfd.CurStage = "도플갱어_말썽쟁이"
			} else if copyRole == "불면증환자" {
				wfd.DIFlag = true
				wfd.TimingChan <- true
			} else {
				wfd.TimingChan <- true
			}
		}
	}
	if wfd.CurStage == "늑대인간_only" {
		if wfd.UserRole[r.UserID] == "늑대인간" {
			if r.Emoji.Name == eOne {
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				uChan, _ := s.UserChannelCreate(r.UserID)
				go func() {
					time.Sleep(10 * time.Second)
					wfd.TimingChan <- true
				}()
				wfd.GameLog += "\n유일한 늑대인간은 버려진 `" + wfd.CardDeck.Cards[0] + "` 를 확인하였습니다."
				s.ChannelMessageSendEmbed(uChan.ID, embed.NewGenericEmbed("<1번: `"+wfd.CardDeck.Cards[0]+"` >", ""))
			}
			if r.Emoji.Name == eTwo {
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				uChan, _ := s.UserChannelCreate(r.UserID)
				go func() {
					time.Sleep(10 * time.Second)
					wfd.TimingChan <- true
				}()
				wfd.GameLog += "\n유일한 늑대인간은 버려진 `" + wfd.CardDeck.Cards[1] + "` 를 확인하였습니다."
				s.ChannelMessageSendEmbed(uChan.ID, embed.NewGenericEmbed("<2번: `"+wfd.CardDeck.Cards[1]+"` >", ""))
			}
			if r.Emoji.Name == eThree {
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				uChan, _ := s.UserChannelCreate(r.UserID)
				go func() {
					time.Sleep(10 * time.Second)
					wfd.TimingChan <- true
				}()
				wfd.GameLog += "\n유일한 늑대인간은 버려진 `" + wfd.CardDeck.Cards[2] + "` 를 확인하였습니다."
				s.ChannelMessageSendEmbed(uChan.ID, embed.NewGenericEmbed("<3번: `"+wfd.CardDeck.Cards[2]+"` >", ""))
			}
		}
	}
	if wfd.CurStage == "예언자" {
		if wfd.UserRole[r.UserID] == "예언자" {
			if r.Emoji.Name == eBin {
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				wfd.CurStage = "예언자_trash"
				msg, _ := s.ChannelMessageSendEmbed(r.ChannelID, embed.NewGenericEmbed("보지 않을 직업을 고르세요.\n< 1 > < 2 > < 3 >", ""))
				s.MessageReactionAdd(r.ChannelID, msg.ID, eOne)
				s.MessageReactionAdd(r.ChannelID, msg.ID, eTwo)
				s.MessageReactionAdd(r.ChannelID, msg.ID, eThree)
			}
			seerUserReactionTask(s, r, wfd)
		}
	}
	if wfd.CurStage == "예언자_trash" {
		if wfd.UserRole[r.UserID] == "예언자" {
			seerTrashReactionTask(s, r, wfd)
		}
	}
	if wfd.CurStage == "강도" {
		if wfd.UserRole[r.UserID] == "강도" {
			robberReactionTask(s, r, wfd)
		}
	}
	if wfd.CurStage == "말썽쟁이_oneMoreChoice" {
		if wfd.UserRole[r.UserID] == "말썽쟁이" {
			tmOneMoreReactionTask(s, r, wfd)
		}
	}
	if wfd.CurStage == "말썽쟁이" {
		if wfd.UserRole[r.UserID] == "말썽쟁이" {
			tmReactionTask(s, r, wfd)
		}
	}
	if wfd.CurStage == "주정뱅이" {
		if wfd.UserRole[r.UserID] == "주정뱅이" {
			drunkReactionTask(s, r, wfd)
		}
	}

	if wfd.CurStage == "Election" {
		for _, item := range wfd.UserIDs {
			go func(uid string) {
				if uid == r.UserID {
					for i := 0; i < len(wfd.UserIDs); i++ {
						if r.Emoji.Name == eNum[i] {
							if r.UserID == wfd.UserIDs[i] && !isTest {
								return
							}
							s.ChannelMessageDelete(r.ChannelID, r.MessageID)
							selCandi := wfd.UserIDs[i]
							voteUser, _ := s.User(uid)
							selUser, _ := s.User(selCandi)
							s.ChannelMessageSend(r.ChannelID, "`"+
								selUser.Username+"` 님에게 투표하였습니다.")
							wfd.ElectChan <- WF.NewElectInfo(uid, voteUser.Username, selCandi)
						}
					}
				}
			}(item)
		}
	}
}

func sendClassInfo(s *discordgo.Session, role string, id string) {
	s.ChannelMessageSendEmbed(id, embed.NewGenericEmbed("**"+role+"**", getRoleInfo(role)))
}

func tmReactionTask(s *discordgo.Session, r *discordgo.MessageReactionAdd, wfd *WF.Data) {
	if exit(wfd.CurStage) {
		return
	}
	tm, _ := s.User(r.UserID)
	tmMsg := ""
	for i := 0; i < len(wfd.UserIDs); i++ {
		if r.Emoji.Name == eNum[i] {
			if wfd.UserIDs[i] != r.UserID {
				if strings.HasPrefix(wfd.CurStage, "도플갱어") {
					wfd.CurStage = "도플갱어_말썽쟁이_choiceWating"
				} else {
					wfd.CurStage = "말썽쟁이_choiceWaiting"
				}
				wfd.IndexChan <- i

				s.ChannelMessageSend(r.ChannelID, "다음 사람을 고르세요")
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				user, _ := s.User(wfd.UserIDs[i])
				selectMsg := "`" + user.Username + "`님을 선택하였습니다."
				wfd.GameLog += "\n말썽쟁이 `" + tm.Username +
					"` (은)는 `" + wfd.FinalRole[wfd.UserIDs[i]] + "` 인 `" +
					user.Username + "` 의 직업과,\n"
				s.ChannelMessageSend(r.ChannelID, selectMsg)
				index := len(wfd.UserIDs)
				for j := 0; j < len(wfd.UserIDs); j++ {
					if i == j {
						index = j
						break
					}
					if (wfd.UserRole[wfd.UserIDs[j]] == "도플갱어" && strings.HasPrefix(wfd.CurStage, "도플갱어")) ||
						(wfd.UserRole[wfd.UserIDs[j]] == "말썽쟁이" && strings.HasPrefix(wfd.CurStage, "말썽쟁이")) {
						tmMsg += "~~"
					}
					user, _ := s.User(wfd.UserIDs[j])
					tmMsg += "<" + strconv.Itoa(j+1) + "번 사용자: " + user.Username + "> "
					if (wfd.UserRole[wfd.UserIDs[j]] == "도플갱어" && strings.HasPrefix(wfd.CurStage, "도플갱어")) ||
						(wfd.UserRole[wfd.UserIDs[j]] == "말썽쟁이" && strings.HasPrefix(wfd.CurStage, "말썽쟁이")) {
						tmMsg += "~~"
					}
				}
				for j := index + 1; j < len(wfd.UserIDs); j++ {
					if (wfd.UserRole[wfd.UserIDs[j]] == "도플갱어" && strings.HasPrefix(wfd.CurStage, "도플갱어")) ||
						(wfd.UserRole[wfd.UserIDs[j]] == "말썽쟁이" && strings.HasPrefix(wfd.CurStage, "말썽쟁이")) {
						tmMsg += "~~"
					}
					user, _ := s.User(wfd.UserIDs[j])
					tmMsg += "<" + strconv.Itoa(j) + "번 사용자: " + user.Username + "> "
					if (wfd.UserRole[wfd.UserIDs[j]] == "도플갱어" && strings.HasPrefix(wfd.CurStage, "도플갱어")) ||
						(wfd.UserRole[wfd.UserIDs[j]] == "말썽쟁이" && strings.HasPrefix(wfd.CurStage, "말썽쟁이")) {
						tmMsg += "~~"
					}
				}
				msg, _ := s.ChannelMessageSendEmbed(r.ChannelID, embed.NewGenericEmbed(tmMsg, ""))
				for i := 0; i < len(wfd.UserIDs)-1; i++ {
					s.MessageReactionAdd(r.ChannelID, msg.ID, eNum[i])
				}
				if strings.HasPrefix(wfd.CurStage, "도플갱어") {
					wfd.CurStage = "도플갱어_말썽쟁이_oneMoreChoice"
				} else {
					wfd.CurStage = "말썽쟁이_oneMoreChoice"
				}
			}
		}
	}
}

func tmOneMoreReactionTask(s *discordgo.Session, r *discordgo.MessageReactionAdd, wfd *WF.Data) {
	if exit(wfd.CurStage) {
		return
	}
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
				wfd.GameLog += "`" + wfd.FinalRole[wfd.UserIDs[i]] + "`" +
					" 인 `" + user.Username + "` 의 직업을 맞바꾸었습니다."
				s.ChannelMessageSend(r.ChannelID, tmMsg)
				temp := wfd.FinalRole[wfd.UserIDs[i]]
				wfd.FinalRole[wfd.UserIDs[i]] = wfd.FinalRole[wfd.UserIDs[prev]]
				wfd.FinalRole[wfd.UserIDs[prev]] = temp
				go func() {
					time.Sleep(3 * time.Second)
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
				wfd.GameLog += "`" + wfd.FinalRole[wfd.UserIDs[i]] + "`" +
					" 인 `" + user.Username + "` 의 직업을 맞바꾸었습니다."
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

func seerTrashReactionTask(s *discordgo.Session, r *discordgo.MessageReactionAdd, wfd *WF.Data) {
	if exit(wfd.CurStage) {
		return
	}
	seer, _ := s.User(r.UserID)
	trashMsg := ""
	wfd.GameLog += "\n예언자 `" + seer.Username + "` (은)는 버려진 직업 "
	if r.Emoji.Name == eOne {
		s.ChannelMessageDelete(r.ChannelID, r.MessageID)
		trashMsg += "<2번: `" + wfd.CardDeck.Cards[1] + "`>" +
			" <3번: `" + wfd.CardDeck.Cards[2] + "`>"
		go func() {
			time.Sleep(5 * time.Second)
			wfd.TimingChan <- true
		}()
		wfd.GameLog += trashMsg + " (을)를 확인하였습니다."
		s.ChannelMessageSend(r.ChannelID, trashMsg)
	}
	if r.Emoji.Name == eTwo {
		s.ChannelMessageDelete(r.ChannelID, r.MessageID)
		trashMsg += "<1번: `" + wfd.CardDeck.Cards[0] + "`>" +
			" <3번: `" + wfd.CardDeck.Cards[2] + "`>"
		go func() {
			time.Sleep(5 * time.Second)
			wfd.TimingChan <- true
		}()
		wfd.GameLog += trashMsg + " (을)를 확인하였습니다."
		s.ChannelMessageSend(r.ChannelID, trashMsg)
	}
	if r.Emoji.Name == eThree {
		s.ChannelMessageDelete(r.ChannelID, r.MessageID)
		trashMsg += "<1번: `" + wfd.CardDeck.Cards[0] + "`>" +
			" <2번: `" + wfd.CardDeck.Cards[1] + "`>"
		go func() {
			time.Sleep(5 * time.Second)
			wfd.TimingChan <- true
		}()
		wfd.GameLog += trashMsg + " (을)를 확인하였습니다."
		s.ChannelMessageSend(r.ChannelID, trashMsg)
	}
}

func seerUserReactionTask(s *discordgo.Session, r *discordgo.MessageReactionAdd, wfd *WF.Data) {
	if exit(wfd.CurStage) {
		return
	}
	for i := 0; i < len(wfd.UserIDs); i++ {
		if r.Emoji.Name == eNum[i] {
			if wfd.UserIDs[i] != r.UserID {
				seer, _ := s.User(r.UserID)
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				if strings.HasPrefix(wfd.CurStage, "도플갱어") {
					wfd.CurStage = "도플갱어_예언자_used_power"
				} else {
					wfd.CurStage = "예언자_used_power"
				}
				user, _ := s.User(wfd.UserIDs[i])
				powerMsg := "`" + user.Username + "` (은)는 `" + wfd.FinalRole[wfd.UserIDs[i]] + "` 입니다."
				wfd.GameLog += "\n`예언자` `" + seer.Username + "` (은)는 `" +
					user.Username + "` 의 직업 `" + wfd.FinalRole[wfd.UserIDs[i]] + "` (을)를 확인하였습니다."
				go func() {
					time.Sleep(5 * time.Second)
					wfd.TimingChan <- true
				}()
				s.ChannelMessageSend(r.ChannelID, powerMsg)
			}
		}
	}
}

func drunkReactionTask(s *discordgo.Session, r *discordgo.MessageReactionAdd, wfd *WF.Data) {
	if exit(wfd.CurStage) {
		return
	}
	dr, _ := s.User(r.UserID)
	for i := 0; i < 3; i++ {
		if r.Emoji.Name == eNum[i] {
			s.ChannelMessageDelete(r.ChannelID, r.MessageID)
			wfd.GameLog += "\n주정뱅이 `" + dr.Username +
				"` 는 버려진 직업 중 `" + wfd.CardDeck.Cards[i] + "` (와)과\n" +
				" 자신의 직업 `주정뱅이` 를 맞바꾸었습니다."
			temp := wfd.CardDeck.Cards[i]
			wfd.CardDeck.Cards[i] = "주정뱅이"
			wfd.FinalRole[dr.ID] = temp
			s.ChannelMessageSendEmbed(r.ChannelID, embed.NewGenericEmbed("", "술에 취한 당신은, "+
				strconv.Itoa(i+1)+"번 직업와 맞바꾸었습니다."+
				"\n이런... 술에 취해 무슨 직업이었는지도 잊어버렸군요.."))
			wfd.TimingChan <- true
		}
	}
}

func robberReactionTask(s *discordgo.Session, r *discordgo.MessageReactionAdd, wfd *WF.Data) {
	if exit(wfd.CurStage) {
		return
	}
	robber, _ := s.User(r.UserID)
	robberMsg := ""
	for i := 0; i < len(wfd.UserIDs); i++ {
		if r.Emoji.Name == eNum[i] {
			if wfd.UserIDs[i] != r.UserID {
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				user, _ := s.User(wfd.UserIDs[i])
				robberMsg = user.Username + "은 `" + wfd.FinalRole[wfd.UserIDs[i]] + "` 이었습니다. 하지만 이젠 아니죠."
				wfd.FinalRole[r.UserID] = wfd.FinalRole[wfd.UserIDs[i]]
				wfd.FinalRole[wfd.UserIDs[i]] = "강도"
				wfd.TimingChan <- true
				wfd.GameLog += "\n강도 `" + robber.Username + "` (은)는 `" +
					user.Username + "` 의 직업 `" + wfd.FinalRole[r.UserID] +
					"` (을)를 확인하고 훔쳤습니다."
				s.ChannelMessageSend(r.ChannelID, robberMsg)
			}
		}
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.ID == "318743234601811969" && strings.HasPrefix(m.Content, "상태설정") {
		if m.Content == "상태설정" {
			s.UpdateStatus(0, prefix+"게임방법 / "+prefix+"도움말")
		}
		splits := strings.Split(m.Content, " ")
		if len(splits) == 2 {
			s.UpdateStatus(0, splits[1])
		}
	}

	if strings.HasPrefix(m.Content, prefix) { // 프리픽스로 시작하는 메시지일 경우
		if m.Content == prefix+"직업목록" {
			classMsg := ""
			for i, item := range classList {
				classMsg += item + " "
				if i%5 == 4 {
					classMsg += "\n"
				}
			}
			classMsg += "\n\n`" +
				prefix + "직업소개 <직업명>` 으로 직업소개 불러오기" +
				"\n`" + prefix + "직업소개 모두` 로 모든 직업소개 DM으로 받기"
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("**구현된 직업 목록**", classMsg))
		}
		if m.Content == prefix+"게임방법" {
			contextEmbedSend(s, m)
		}
		if m.Content == prefix+"명령어" {
			commandEmbedSend(s, m)
		}
		if m.Content == prefix+"참고" {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("**참고**", psMsg))
		}
		if m.Content == prefix+"게임배경" {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("**게임배경**", helpMsg))
		}
		if m.Content == prefix+"능력순서" {
			gameSeqMsg := ""
			for i, item := range gameSeq {
				gameSeqMsg += item + " -> "
				if i%3 == 2 {
					gameSeqMsg += "\n"
				}
			}
			gameSeqMsg += "투표 시작"
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("**특수능력 사용 순서**", gameSeqMsg))
		}
		if m.Content == prefix+"help" || m.Content == prefix+"도움말" {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("**게임배경**", helpMsg))
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("**참고**", psMsg))
			contextEmbedSend(s, m)
			commandEmbedSend(s, m)
		}
		if m.Content == prefix+"승리조건" {
			winCheckEmbed := embed.NewEmbed()
			winCheckEmbed.SetTitle("**각 진영의 승리조건**")
			winCheckEmbed.AddField(vilWinTitle, vilWinMsg)
			winCheckEmbed.AddField(wereWinTitle, wereWinMsg+"\n\n"+tannerWinTitle)
			s.ChannelMessageSendEmbed(m.ChannelID, winCheckEmbed.MessageEmbed)
		}
		if m.Content == prefix+"테스트" && m.Author.ID == "318743234601811969" {
			if isTest {
				s.ChannelMessageSend(m.ChannelID, "테스트 모드를 끕니다.")
			} else {
				s.ChannelMessageSend(m.ChannelID, "테스트 모드를 켭니다.")
			}
			isTest = !isTest
		}

		if strings.HasPrefix(m.Content, prefix+"직업소개") {
			classStr := strings.Split(m.Content, " ")
			if len(classStr) != 1 {
				classFlag := false
				for _, item := range classList {
					if classStr[1] == item {
						classFlag = true
						s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("**"+item+" 소개**", getRoleInfo(item)))
						break
					}
				}
				if classStr[1] == "모두" {
					classFlag = true
					uChan, _ := s.UserChannelCreate(m.Author.ID)
					for _, item := range classList {
						s.ChannelMessageSendEmbed(uChan.ID, embed.NewGenericEmbed("**"+item+" 소개**", getRoleInfo(item)))
					}
					s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("모든 직업 소개가 DM으로 전송되었습니다.", ""))
				}
				if !classFlag {
					s.ChannelMessageSend(m.ChannelID, "직업 이름이 잘못되었습니다.")
				}
			} else {
				s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("직업소개", prefix+"직업소개 <직업명> 으로 요청하세요."))
			}
		}

		var wfd *WF.Data

		if m.Content == prefix+"강제종료" {
			if !isGuildIn[m.GuildID] {
				s.ChannelMessageSend(m.ChannelID, "현재 서버에서 진행중인 게임이 없습니다.")
				return
			}
			wfd = wfDataMap[m.GuildID]
			wfd.CurStage = "Exit"
			s.ChannelMessageSend(wfd.UseChannelID, "안전하게 강제종료 수행중..")
			time.Sleep(time.Second * 5)
			wfd.CurStage = "Exit"
			cancelGameTask(m)
			s.ChannelMessageSend(wfd.UseChannelID, "사용 종료가 정상적으로 완료되었습니다.")
		}

		if m.Content == prefix+"시작" && !isGuildIn[m.GuildID] {
			if isUserIn[m.Author.ID] {
				s.ChannelMessageSend(m.ChannelID, "이미 다른 서버에서 게임중인 유저입니다.")
				return
			}
			isGuildIn[m.GuildID] = true
			isUserIn[m.Author.ID] = true
			wfDataMap[m.GuildID] = WF.NewWFData(m.Author.ID, m.ChannelID)
			uidToGid[m.Author.ID] = m.GuildID
			wfDataMap[m.GuildID].CurStage = "Prepare_card"
			wfd = wfDataMap[m.GuildID]
			cardSetting(s, m.GuildID, wfd)
			<-wfd.TimingChan
			prevSettingMap[m.GuildID] = WF.NewSettingData(wfd.CardDeck, wfd.MaxUser)
			wfd.OriCardDeck = wfd.CardDeck
			prevSettingMap[m.GuildID].CardDeck.SortCards()
			cardMsg := ""
			for _, item := range wfd.CardDeck.Cards {
				cardMsg += "\n" + item
			}
			cardMsg += "\n**총 " + strconv.Itoa(len(wfd.CardDeck.Cards)) + "개의 직업이 선정되었습니다." +
				"\n총 플레이 인원은 " + strconv.Itoa(len(wfd.CardDeck.Cards)-3) + "명 입니다.**"
			s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed("**설정된 직업**", cardMsg))
			wfd.CardDeck.ShuffleCards()
			newUserTask(m)
			wfDataMap[m.GuildID].CurStage = "Prepare"

			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(
				"** 시작 준비...**", "`"+prefix+"입장` 으로 입장하세요"+
					"\n`"+prefix+"마감` 으로 입장을 마감하세요."))
		}
		if isGuildIn[m.GuildID] {
			wfd = wfDataMap[m.GuildID]
			if wfd.CurStage == "Prepare_card" {
				if m.Content == prefix+"ㅇㅇ" {
					wfd.CardDeck.ChoiceChan <- 0
				} else if m.Content == prefix+"ㄴㄴ" {
					wfd.CardDeck.ChoiceChan <- 1
				}
			}
			if wfd.CurStage == "Prepare_class" {
				deckLen := len(cardMap[m.GuildID])
				userLen := deckLen - 3
				if m.Content == prefix+"늑대인간" {
					s.ChannelMessageSend(wfd.UseChannelID, "늑대인간은 2명이 최대입니다."+
						"\n현재 직업 수: "+strconv.Itoa(deckLen)+""+
						"\n현재 플레이어 수: "+strconv.Itoa(userLen))
				}
				if m.Content == prefix+"사냥꾼" {
					for _, item := range cardMap[m.GuildID] {
						if item == "사냥꾼" {
							s.ChannelMessageSend(wfd.UseChannelID, "사냥꾼은 최대 1명입니다."+
								"\n현재 직업 수: "+strconv.Itoa(deckLen)+""+
								"\n현재 플레이어 수: "+strconv.Itoa(userLen))
							return
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "사냥꾼")
					s.ChannelMessageSend(wfd.UseChannelID, "사냥꾼을 넣었습니다."+
						"\n현재 직업 수: "+strconv.Itoa(deckLen+1)+""+
						"\n현재 플레이어 수: "+strconv.Itoa(userLen+1))
				}
				if m.Content == prefix+"도플갱어" {
					for _, item := range cardMap[m.GuildID] {
						if item == "도플갱어" && !isTest {
							s.ChannelMessageSend(wfd.UseChannelID, "도플갱어는 최대 1명입니다."+
								"\n현재 직업 수: "+strconv.Itoa(deckLen)+""+
								"\n현재 플레이어 수: "+strconv.Itoa(userLen))
							return
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "도플갱어")
					s.ChannelMessageSend(wfd.UseChannelID, "도플갱어를 넣었습니다."+
						"\n현재 직업 수: "+strconv.Itoa(deckLen+1)+""+
						"\n현재 플레이어 수: "+strconv.Itoa(userLen+1))
				}
				if m.Content == prefix+"하수인" {
					for _, item := range cardMap[m.GuildID] {
						if item == "하수인" {
							s.ChannelMessageSend(wfd.UseChannelID, "하수인은 최대 1명입니다."+
								"\n현재 직업 수: "+strconv.Itoa(deckLen)+""+
								"\n현재 플레이어 수: "+strconv.Itoa(userLen))
							return
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "하수인")
					s.ChannelMessageSend(wfd.UseChannelID, "하수인을 넣었습니다."+
						"\n현재 직업 수: "+strconv.Itoa(deckLen+1)+""+
						"\n현재 플레이어 수: "+strconv.Itoa(userLen+1))
				}
				if m.Content == prefix+"프리메이슨" {
					for _, item := range cardMap[m.GuildID] {
						if item == "프리메이슨" {
							s.ChannelMessageSend(wfd.UseChannelID, "프리메이슨은 최대 2명입니다."+
								"\n현재 직업 수: "+strconv.Itoa(deckLen)+""+
								"\n현재 플레이어 수: "+strconv.Itoa(userLen))
							return
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "프리메이슨")
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "프리메이슨")
					s.ChannelMessageSend(wfd.UseChannelID, "프리메이슨을 2번 추가하였습니다."+
						"\n프리메이슨은 한쌍씩만 넣을 수 없습니다."+
						"\n현재 직업 수: "+strconv.Itoa(deckLen+2)+""+
						"\n현재 플레이어 수: "+strconv.Itoa(userLen+2))
				}
				if m.Content == prefix+"예언자" {
					for _, item := range cardMap[m.GuildID] {
						if item == "예언자" {
							s.ChannelMessageSend(wfd.UseChannelID, "예언자은 최대 1명입니다."+
								"\n현재 직업 수: "+strconv.Itoa(deckLen)+""+
								"\n현재 플레이어 수: "+strconv.Itoa(userLen))
							return
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "예언자")
					s.ChannelMessageSend(wfd.UseChannelID, "에언자를 넣었습니다."+
						"\n현재 직업 수: "+strconv.Itoa(deckLen+1)+""+
						"\n현재 플레이어 수: "+strconv.Itoa(userLen+1))
				}
				if m.Content == prefix+"말썽쟁이" {
					for _, item := range cardMap[m.GuildID] {
						if item == "말썽쟁이" && !isTest {
							s.ChannelMessageSend(wfd.UseChannelID, "말썽쟁이은 최대 1명입니다."+
								"\n현재 직업 수: "+strconv.Itoa(deckLen)+""+
								"\n현재 플레이어 수: "+strconv.Itoa(userLen))
							return
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "말썽쟁이")
					s.ChannelMessageSend(wfd.UseChannelID, "말썽쟁이를 넣었습니다."+
						"\n현재 직업 수: "+strconv.Itoa(deckLen+1)+""+
						"\n현재 플레이어 수: "+strconv.Itoa(userLen+1))
				}
				if m.Content == prefix+"무두장이" {
					for _, item := range cardMap[m.GuildID] {
						if item == "무두장이" {
							s.ChannelMessageSend(wfd.UseChannelID, "무두장이은 최대 1명입니다."+
								"\n현재 직업 수: "+strconv.Itoa(deckLen)+""+
								"\n현재 플레이어 수: "+strconv.Itoa(userLen))
							return
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "무두장이")
					s.ChannelMessageSend(wfd.UseChannelID, "무두장이를 넣었습니다."+
						"\n현재 직업 수: "+strconv.Itoa(deckLen+1)+""+
						"\n현재 플레이어 수: "+strconv.Itoa(userLen+1))
				}
				if m.Content == prefix+"불면증환자" {
					for _, item := range cardMap[m.GuildID] {
						if item == "불면증환자" {
							s.ChannelMessageSend(wfd.UseChannelID, "불면증환자은 최대 1명입니다."+
								"\n현재 직업 수: "+strconv.Itoa(deckLen)+""+
								"\n현재 플레이어 수: "+strconv.Itoa(userLen))
							return
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "불면증환자")
					s.ChannelMessageSend(wfd.UseChannelID, "불면증환자를 넣었습니다."+
						"\n현재 직업 수: "+strconv.Itoa(deckLen+1)+""+
						"\n현재 플레이어 수: "+strconv.Itoa(userLen+1))
				}
				if m.Content == prefix+"강도" {
					for _, item := range cardMap[m.GuildID] {
						if item == "강도" {
							s.ChannelMessageSend(wfd.UseChannelID, "강도은 최대 1명입니다."+
								"\n현재 직업 수: "+strconv.Itoa(deckLen)+""+
								"\n현재 플레이어 수: "+strconv.Itoa(userLen))
							return
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "강도")
					s.ChannelMessageSend(wfd.UseChannelID, "강도를 넣었습니다."+
						"\n현재 직업 수: "+strconv.Itoa(deckLen+1)+""+
						"\n현재 플레이어 수: "+strconv.Itoa(userLen+1))
				}
				if m.Content == prefix+"마을주민" {
					count := 0
					for _, item := range cardMap[m.GuildID] {
						if item == "마을주민" {
							count++
							if count == 3 {
								s.ChannelMessageSend(wfd.UseChannelID, "마을주민은 최대 3명입니다."+
									"\n현재 직업 수: "+strconv.Itoa(deckLen)+""+
									"\n현재 플레이어 수: "+strconv.Itoa(userLen))
								return
							}
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "마을주민")
					s.ChannelMessageSend(wfd.UseChannelID, "마을주민을 넣었습니다."+
						"\n현재 직업 수: "+strconv.Itoa(deckLen+1)+""+
						"\n현재 플레이어 수: "+strconv.Itoa(userLen+1))
				}
				if m.Content == prefix+"주정뱅이" {
					count := 0
					for _, item := range cardMap[m.GuildID] {
						if item == "주정뱅이" {
							count++
							if count == 1 && !isTest {
								s.ChannelMessageSend(wfd.UseChannelID, "주정뱅이는 최대 1명입니다."+
									"\n현재 직업 수: "+strconv.Itoa(deckLen)+""+
									"\n현재 플레이어 수: "+strconv.Itoa(userLen))
								return
							}
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "주정뱅이")
					s.ChannelMessageSend(wfd.UseChannelID, "주정뱅이를 넣었습니다."+
						"\n현재 직업 수: "+strconv.Itoa(deckLen+1)+""+
						"\n현재 플레이어 수: "+strconv.Itoa(userLen+1))
				}
				if m.Content == prefix+"완료" {

					wfd.CardDeck.ChoiceChan <- 0
				}
			}
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
				s.ChannelMessageSend(m.ChannelID, "현재인원: ("+strconv.Itoa(len(wfd.UserIDs))+"/"+strconv.Itoa(wfd.MaxUser)+")")
			}
			if strings.HasPrefix(m.Content, prefix+"마감") && wfd.CurStage == "Prepare" {
				if len(wfd.UserIDs) == wfd.MaxUser {
					wfd.CurStage = "Prepare_finishing"
					for _, item := range wfd.UserIDs {
						go func(uid string) {
							uChan, _ := s.UserChannelCreate(uid)
							user, _ := s.User(uid)
							wfd.GameLog += "`" + user.Username + "` 님의 직업이 `" + wfd.UserRole[uid] + "` (으)로 배정되었습니다.\n"
							roleTitle := "**당신의 직업은 **`" + wfd.UserRole[uid] + "`**입니다.**"
							roleBrief := getRoleInfo(wfd.UserRole[uid])

							s.ChannelMessageSend(uChan.ID, "")
							s.ChannelMessageSendEmbed(uChan.ID, embed.NewGenericEmbed(roleTitle, roleBrief))
						}(item)
					}

					s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed("", "게임을 시작합니다."+
						"\n게임의 상태는 아래 메시지에 표시됩니다."))
					stageMsg, _ := s.ChannelMessageSend(wfd.UseChannelID, "게임 시작 준비중...")
					guild, _ := s.Guild(m.GuildID)
					filestr := "서버 " + guild.Name + " 에서 게임 시작됨\n"
					file, _ := os.OpenFile("GuildLog.txt", os.O_WRONLY|os.O_APPEND, os.FileMode(644))
					defer file.Close()
					file.WriteString(filestr)

					nextStage(wfd)
					doppelTask(s, wfd, stageMsg)
					werewolfTask(s, wfd, stageMsg)
					minionTask(s, wfd, stageMsg)
					masonTask(s, wfd, stageMsg)
					seerTask(s, wfd, stageMsg)
					robberTask(s, wfd, stageMsg)
					tmTask(s, wfd, stageMsg)
					drunkTask(s, wfd, stageMsg)
					insomniacTask(s, wfd, stageMsg)
					s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "모두의 차례 종료.")
					dayBriefTask(s, wfd)
					cancelGameTask(m)
				} else {
					s.ChannelMessageSend(m.ChannelID, "정확한 인원이 모이지 않았습니다. ("+strconv.Itoa(len(wfd.UserIDs))+"/"+strconv.Itoa(wfd.MaxUser)+")")
				}
			}
		}
	}

	homeworkMethod(s, m)
}

func commandEmbedSend(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("**명령어 목록**", commandMsg))
}

func contextEmbedSend(s *discordgo.Session, m *discordgo.MessageCreate) {
	contextMsg := "게임을 시작하려면, `" + prefix + "시작` 을 입력하세요." +
		"\n이후 추가할 직업을 `" + prefix + "<직업명>` 으로 덱에 넣은 뒤" +
		"\n`" + prefix + "완료` 로 직업 설정을 완료하세요." +
		"\n게임에 참가할 플레이어들은 `" + prefix + "입장` 을 입력하고," +
		"\n모두 참가했다면 `" + prefix + "마감` 으로 게임을 시작하세요." +
		"\n" +
		"\n`" + prefix + "명령어` 로 사용 가능한 명령어를 볼 수 있습니다." +
		"\n`" + prefix + "게임배경` 로 게임배경를 볼 수 있습니다." +
		"\n`" + prefix + "참고` 로 참고 메시지를 볼 수 있습니다." +
		"\n승리 조건은 `" + prefix + "승리조건` 으로," +
		"\n직업 설명은 `" + prefix + "직업소개 <직업명>` 으로 확인하세요." +
		"\n`" + prefix + "직업소개 모두` 를 입력하면 모든 직업 소개를" +
		"\nDM으로 받을 수 있습니다." +
		"\n" +
		"\n만약 문제가 발생했다면 `" + prefix + "강제종료` 로" +
		"\n게임을 강제 종료할 수 있습니다."
	s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("**게임방법**", contextMsg))
}

func cardSetting(s *discordgo.Session, gid string, wfd *WF.Data) {
	if exit(wfd.CurStage) {
		return
	}
	cardMap[gid] = make([]string, 0, 10)

	classSetTitle := "**직업 설정을 시작합니다.**"
	classSetMsg := "이전 설정과 동일한 직업 설정을 사용할까요?\n`" + prefix + "ㅇㅇ` / `" + prefix + "ㄴㄴ`"
	s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed(classSetTitle, classSetMsg))

	choice := <-wfd.CardDeck.ChoiceChan
	if choice == 0 {
		if prevSettingMap[gid] != nil {
			wfd.CardDeck.Cards = prevSettingMap[gid].CardDeck.Cards
			wfd.MaxUser = prevSettingMap[gid].MaxUser
			wfd.TimingChan <- true
			return
		} else {
			s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed("", "\n이전 게임 기록이 남아있지 않습니다.."+
				"\n게임을 한 적이 없거나, 서버가 재부팅되었을 수 있습니다.\n"))
		}
	}
	classSetMsg = "추천 직업 설정을 사용할까요?\n`" + prefix + "ㅇㅇ` / `" + prefix + "ㄴㄴ`"
	s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed(classSetTitle, classSetMsg))

	choice = <-wfd.CardDeck.ChoiceChan
	if choice == 0 {
		classSetTitle = "**추천 직업 목록 (방장만 선택가능)**"
		recomMsg := "1. 나무위키 추천 (5인)\n"
		for i, item := range namuDeck {
			recomMsg += item + " "
			if i%5 == 4 {
				recomMsg += "\n"
			}
		}
		recomMsg += "\n"
		recomMsg += "\n2. 도플갱어 또 버그남? (6인)\n"
		for i, item := range doppelBugDeck {
			recomMsg += item + " "
			if i%5 == 4 {
				recomMsg += "\n"
			}
		}
		recomMsg += "\n"
		recomMsg += "\n3. 아무튼 나 마을주민임 (7인)\n"
		for i, item := range anywayDeck {
			recomMsg += item + " "
			if i%5 == 4 && i != len(anywayDeck)-1 {
				recomMsg += "\n"
			}
		}
		recomMsg += "\n"
		recomMsg += "\n4. 예언자 나오세요 (8인)\n"
		for i, item := range comeOnSeerDeck {
			recomMsg += item + " "
			if i%5 == 4 && i != len(comeOnSeerDeck)-1 {
				recomMsg += "\n"
			}
		}
		recomMsg += "\n"
		recomMsg += "\n5. 뭘 좋아할지 몰라 다 준비해봤어 (11인)\n"
		for i, item := range chaosDeck {
			recomMsg += item + " "
			if i%5 == 4 {
				recomMsg += "\n"
			}
		}
		wfd.CurStage = "recomChoice"
		recomReactMsg, _ := s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed(classSetTitle, recomMsg))
		for i := 0; i < 5; i++ {
			s.MessageReactionAdd(wfd.UseChannelID, recomReactMsg.ID, eNum[i])
		}
		recChoice := <-wfd.CardDeck.ChoiceChan
		s.ChannelMessageDelete(wfd.UseChannelID, recomReactMsg.ID)
		if recChoice == 0 {
			for _, item := range namuDeck {
				cardMap[gid] = append(cardMap[gid], item)
			}
		} else if recChoice == 1 {
			for _, item := range doppelBugDeck {
				cardMap[gid] = append(cardMap[gid], item)
			}
		} else if recChoice == 2 {
			for _, item := range anywayDeck {
				cardMap[gid] = append(cardMap[gid], item)
			}
		} else if recChoice == 3 {
			for _, item := range comeOnSeerDeck {
				cardMap[gid] = append(cardMap[gid], item)
			}
		} else if recChoice == 4 {
			for _, item := range chaosDeck {
				cardMap[gid] = append(cardMap[gid], item)
			}
		}

		wfd.CardDeck.Cards = cardMap[gid]
		wfd.CardDeck.ShuffleCards()
		wfd.MaxUser = len(cardMap[gid]) - 3

		wfd.TimingChan <- true
		return
	}

	cardMap[gid] = append(cardMap[gid], "늑대인간")
	cardMap[gid] = append(cardMap[gid], "늑대인간")
	s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed("", "늑대인간 2장은 필수입니다. 직업 덱에 넣었습니다."))
	for true {
		wfd.CurStage = "Prepare_class"
		classTitle := "추가할 직업들을 입력하세요. (ex: " + prefix + "마을주민)"
		classMsg := "모두 입력한 후 `" + prefix + "완료` 로 다음 단계로 넘어가세요."
		classFieldTitle := "**구현된 직업 목록:**"
		classFieldMsg := ""
		for _, item := range classList {
			classFieldMsg += item + "\n"
		}
		s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed(classTitle, classMsg))
		s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed(classFieldTitle, classFieldMsg))
		<-wfd.CardDeck.ChoiceChan
		if len(cardMap[gid]) < 6 {
			s.ChannelMessageSendEmbed(wfd.UseChannelID,
				embed.NewGenericEmbed("직업을 6개 이상을 골라야 합니다..",
					"현재 직업 수: ("+strconv.Itoa(len(cardMap[gid]))+"/6)"))
		} else {
			wfd.CardDeck.Cards = cardMap[gid]
			wfd.MaxUser = len(cardMap[gid]) - 3
			break
		}
	}

	wfd.TimingChan <- true
}

func dayBriefTask(s *discordgo.Session, wfd *WF.Data) {
	if exit(wfd.CurStage) {
		return
	}
	briefMsg := ""
	briefTitle := "**모든 특수 능력 사용이 끝났습니다.**"
	briefMsg += "3초 후 여러분들에게 투표 용지가 전송됩니다." +
		"\n한번 투표한 내용은 바꿀 수 없기에," +
		"\n**신중하게 투표하세요.**"
	go func() {
		time.Sleep(time.Second * 3)
		wfd.TimingChan <- true
	}()

	s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed(briefTitle, briefMsg))
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
		msg, _ := s.ChannelMessageSend(cid, briefMsg)
		for i := 0; i < len(wfd.UserIDs); i++ {
			s.MessageReactionAdd(cid, msg.ID, eNum[i])
		}
	}
	electFinishTask(s, wfd)
}

func electFinishTask(s *discordgo.Session, wfd *WF.Data) {
	if exit(wfd.CurStage) {
		return
	}
	electData := make([]*WF.ElectInfo, 0, 10)
	electResult := make([]int, len(wfd.UserIDs))

	s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed("투표를 시작합니다!", ""))
	for i := 0; i < len(wfd.UserIDs); i++ {
		electData = append(electData, <-wfd.ElectChan)
		electAlarmMsg := "`" + electData[i].VoterName + "`님이 투표하셨습니다."
		s.ChannelMessageSend(wfd.UseChannelID, electAlarmMsg)
	}

	s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed("투표가 끝났습니다.", ""))
	s.ChannelMessageSend(wfd.UseChannelID, "결과 계산중...")
	huntTargets := make([]string, 0, 3)
	for i, uid := range wfd.UserIDs {
		for _, elc := range electData {
			if uid == elc.CandiID {
				electResult[i]++
			}
			if wfd.FinalRole[elc.VoterID] == "사냥꾼" {
				huntTargets = append(huntTargets, elc.CandiID)
			}
		}
	}
	electResultMsg := ""
	max := 0
	maxi := -1
	for i, item := range electResult {
		if max < item {
			max = item
			maxi = i
		}
		if item != 0 {
			user, _ := s.User(wfd.UserIDs[i])
			electResultMsg += "<`" + user.Username + "` : " +
				strconv.Itoa(item) + "표>\n"
		}
	}
	if max == 1 {
		electResultMsg += "모두가 한표씩을 받게 되었습니다. 아무도 처형되지 않았습니다.\n"

		s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed("**투표 결과**", electResultMsg))
		finalRoleTitle := "**모두의 최종 직업**"
		finalRoleMsg := ""
		for _, item := range wfd.UserIDs {
			user, _ := s.User(item)
			finalRoleMsg += "\n<`" + user.Username + "` : `" + wfd.FinalRole[item] + "`>"
		}

		s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed(finalRoleTitle, finalRoleMsg))
		time.Sleep(3 * time.Second)
		s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed("**게임 로그**", wfd.GameLog))
		return
	}

	for i, item := range electResult {
		if item == max {
			user, _ := s.User(wfd.UserIDs[i])
			electResultMsg += " `" + user.Username + "`"
		}
	}

	electResultMsg += " 님이 총 " + strconv.Itoa(electResult[maxi]) + " 표로 처형되었습니다."

	s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed("**투표 결과**", electResultMsg))
	s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed("", "처형된 사람의 직업은...?"))
	for i := 0; i < 3; i++ {
		s.ChannelMessageSend(wfd.UseChannelID, "...")
		time.Sleep(time.Second)
	}
	executeTitle := "**처형된 사람의 직업은**"
	executeMsg := ""
	for i, item := range electResult {
		if item == max {
			wfd.DeadUserMap[wfd.UserIDs[i]] = true
			user, _ := s.User(wfd.UserIDs[i])
			executeMsg += "\n<`" + user.Username + "` : `" +
				wfd.UserRole[wfd.UserIDs[i]] + "`-> `" +
				wfd.FinalRole[wfd.UserIDs[i]] + "`>"
		}
	}
	executeMsg += " 입니다."

	s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed(executeTitle, executeMsg))
	if len(huntTargets) != 0 {
		hunterUser := make([]*discordgo.User, 0, 3)
		for _, item := range wfd.UserIDs {
			if wfd.FinalRole[item] == "사냥꾼" && wfd.DeadUserMap[item] {
				tempUser, _ := s.User(item)
				hunterUser = append(hunterUser, tempUser)
			}
		}
		huntMsg := ""
		for i, item := range hunterUser {
			huntTargetUser, _ := s.User(huntTargets[i])
			huntMsg += "사냥꾼 `" + item.Username + "` (은)는`" +
				huntTargetUser.Username + "` 를 투표하여\n처형당하기 전, `" +
				huntTargetUser.Username + "` (을)를 사냥하였습니다." +
				"\n`" + huntTargetUser.Username + "` (은)는 `" + wfd.FinalRole[huntTargets[i]] + "` 이었습니다.\n"
		}
		if len(hunterUser) != 0 {
			s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed("**사냥꾼이 처형되었습니다!**", huntMsg))
		}
	}
	finalRoleTitle := "**모두의 최종 직업**"
	finalRoleMsg := ""
	for _, item := range wfd.UserIDs {
		user, _ := s.User(item)
		finalRoleMsg += "\n<`" + user.Username + "` : `" + wfd.FinalRole[item] + "`>"
	}

	s.ChannelMessageSendEmbed(wfd.UseChannelID, embed.NewGenericEmbed(finalRoleTitle, finalRoleMsg))
	time.Sleep(3 * time.Second)

	logEmbed := embed.NewEmbed()
	logEmbed.AddField("게임 로그", wfd.GameLog)
	s.ChannelMessageSendEmbed(wfd.UseChannelID, logEmbed.MessageEmbed)
}

func masonTask(s *discordgo.Session, wfd *WF.Data, stageMsg *discordgo.Message) {
	if exit(wfd.CurStage) {
		return
	}
	role := "프리메이슨"
	isIn := isInCheck(wfd, role)

	if !isIn {
		s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "프리메이슨은 넣지 않았습니다.")
		nextStage(wfd)
		return
	}
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "프리메이슨의 차례 진행중....")

	masonID := make([]string, 0, 3)
	for _, item := range wfd.UserIDs {
		if wfd.UserRole[item] == "프리메이슨" ||
			(wfd.UserRole[item] == "도플갱어" && wfd.FinalRole[item] == "프리메이슨") &&
				!wfd.DoppelRobberFlag && !wfd.DoppelDrunkFlag {
			masonID = append(masonID, item)
			break
		}
	}
	if len(masonID) == 0 {
		go func() {
			time.Sleep(10 * time.Second)
			wfd.GameLog += "\n프리메이슨은 없었습니다."
			wfd.TimingChan <- true
		}()
		<-wfd.TimingChan
		nextStage(wfd)
		s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "프리메이슨의 차례 종료.")
		return
	}

	masonChanID := make([]string, 0, 3)
	temp := make([]*discordgo.Channel, 0, 3)
	for i, item := range masonID {
		temp[i], _ = s.UserChannelCreate(item)
		masonChanID = append(masonChanID, temp[i].ID)
	}
	for _, item := range masonChanID {
		go func(cid string) {
			masonMsg := ""
			for _, item := range masonID {
				user, _ := s.User(item)
				masonMsg += "<`" + user.Username + "`> "
			}

			if len(masonChanID) == 1 {
				wfd.GameLog += "\n유일한 프리메이슨 " + masonMsg + "(은)는" +
					"\n자신이 유일한 프리메이슨임을 확인하였습니다."
				s.ChannelMessageSendEmbed(cid, embed.NewGenericEmbed("당신은 유일한 프리메이슨입니다.", masonMsg))
			} else {
				wfd.GameLog += "\n프리메이슨 " + masonMsg + "(은)는" +
					"\n서로를 확인하였습니다."
				s.ChannelMessageSendEmbed(cid, embed.NewGenericEmbed("프리메이슨 동료들이 모습을 드러냅니다.", masonMsg))
			}
		}(item)
	}

	go func() {
		time.Sleep(10 * time.Second)
		wfd.TimingChan <- true
	}()

	<-wfd.TimingChan
	nextStage(wfd)
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "프리메이슨의 차례 종료.")
}

func doppelTask(s *discordgo.Session, wfd *WF.Data, stageMsg *discordgo.Message) {
	if exit(wfd.CurStage) {
		return
	}
	role := "도플갱어"
	isIn := isInCheck(wfd, role)
	if !isIn {
		s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "도플갱어는 넣지 않았습니다.")
		nextStage(wfd)
		return
	}
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "도플갱어의 차례 진행중....")

	doppelID := ""
	for _, item := range wfd.UserIDs {
		if wfd.UserRole[item] == "도플갱어" {
			doppelID = item
			break
		}
	}

	if doppelID == "" {
		go func() {
			time.Sleep(20 * time.Second)
			wfd.GameLog += "\n도플갱어는 없었습니다."
			wfd.TimingChan <- true
		}()
		<-wfd.TimingChan
		nextStage(wfd)
		s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "도플갱어의 차례 종료.")
		return
	}
	uChan, _ := s.UserChannelCreate(doppelID)
	doppelRole := "도플갱어"
	sendAllUserAddReaction(s, wfd, doppelRole, uChan)

	<-wfd.TimingChan
	nextStage(wfd)
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "도플갱어의 차례 종료.")
}

func isInCheck(wfd *WF.Data, role string) bool {
	if exit(wfd.CurStage) {
		return false
	}
	isIn := false
	for _, item := range wfd.OriCardDeck.Cards {
		if item == role {
			isIn = true
			break
		}
	}
	return isIn
}

func sendAllUserAddReaction(s *discordgo.Session, wfd *WF.Data, role string, uChan *discordgo.Channel) {
	if exit(wfd.CurStage) {
		return
	}
	userListMsg := "자신은 선택할 수 없습니다.\n"
	for i, item := range wfd.UserIDs {
		if wfd.UserRole[item] == role {
			userListMsg += "~~"
		}
		user, _ := s.User(item)
		userListMsg += "<" + strconv.Itoa(i+1) + "번 사용자: " + user.Username + ">\t"
		if wfd.UserRole[item] == role {
			userListMsg += "~~"
		}
	}

	msg, _ := s.ChannelMessageSend(uChan.ID, userListMsg)
	for i := 0; i < len(wfd.UserIDs); i++ {
		s.MessageReactionAdd(uChan.ID, msg.ID, eNum[i])
	}
}

func werewolfTask(s *discordgo.Session, wfd *WF.Data, stageMsg *discordgo.Message) {
	if exit(wfd.CurStage) {
		return
	}
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "늑대인간의 차례 진행중....")
	wolvesID := make([]string, 0, 10)
	for _, item := range wfd.UserIDs {
		if wfd.UserRole[item] == "늑대인간" ||
			(wfd.UserRole[item] == "도플갱어" && wfd.FinalRole[item] == "늑대인간" &&
				!wfd.DoppelDrunkFlag && !wfd.DoppelRobberFlag) {
			wolvesID = append(wolvesID, item)
		}
	}
	if len(wolvesID) != 1 {
		go func() {
			time.Sleep(10 * time.Second)
			wfd.GameLog += "\n늑대인간이 " + strconv.Itoa(len(wolvesID)) + " 명이라 서로를 확인만 합니다."
			wfd.TimingChan <- true
		}()
	}

	if len(wolvesID) == 1 {
		wfd.CurStage = "늑대인간_only"

		wolvesMsg := "세 장의 비공개 직업 중 한 개를 선택하세요\n"
		wolvesMsg += "< 1 > < 2 > < 3 >"
		uChan, _ := s.UserChannelCreate(wolvesID[0])
		msg, _ := s.ChannelMessageSend(uChan.ID, wolvesMsg)
		s.MessageReactionAdd(uChan.ID, msg.ID, eOne)
		s.MessageReactionAdd(uChan.ID, msg.ID, eTwo)
		s.MessageReactionAdd(uChan.ID, msg.ID, eThree)
		<-wfd.TimingChan
		nextStage(wfd)
		s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "늑대인간의 차례 종료")
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
	nextStage(wfd)
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "늑대인간의 차례 종료.")
}

func minionTask(s *discordgo.Session, wfd *WF.Data, stageMsg *discordgo.Message) {
	if exit(wfd.CurStage) {
		return
	}
	role := "하수인"
	isIn := isInCheck(wfd, role)
	if !isIn {
		s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "하수인은 넣지 않았습니다.")
		nextStage(wfd)
		return
	}
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "하수인의 차례 진행중....")

	wolvesID := make([]string, 0, 10)
	minionID := make([]string, 0, 10)
	for _, item := range wfd.UserIDs {
		if wfd.UserRole[item] == "늑대인간" ||
			(wfd.UserRole[item] == "도플갱어" && wfd.FinalRole[item] == "늑대인간" &&
				!wfd.DoppelDrunkFlag) {
			wolvesID = append(wolvesID, item)
		}
		if wfd.UserRole[item] == "하수인" ||
			(wfd.UserRole[item] == "도플갱어" && wfd.FinalRole[item] == "하수인" &&
				!wfd.DoppelDrunkFlag) {
			minionID = append(minionID, item)
		}
	}
	minionMsg := "늑대인간은, "
	for _, item := range wolvesID {
		user, _ := s.User(item)
		minionMsg += "<" + user.Username + "> "
	}
	minionMsg += "입니다."

	if len(wolvesID) == 0 {
		minionMsg = "늑대인간이 존재하지 않습니다."
	}

	go func() {
		time.Sleep(10 * time.Second)
		wfd.TimingChan <- true
	}()

	if len(minionID) == 0 {

		<-wfd.TimingChan
		nextStage(wfd)
		s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "하수인의 차례 종료.")
		return
	}
	for _, item := range minionID {
		user, _ := s.User(item)
		uChan, _ := s.UserChannelCreate(user.ID)
		s.ChannelMessageSend(uChan.ID, minionMsg)
	}
	<-wfd.TimingChan
	nextStage(wfd)
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "하수인의 차례 종료.")
}

func seerTask(s *discordgo.Session, wfd *WF.Data, stageMsg *discordgo.Message) {
	if exit(wfd.CurStage) {
		return
	}
	role := "예언자"
	isIn := isInCheck(wfd, role)
	if !isIn {
		s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "예언자는 넣지 않았습니다.")
		nextStage(wfd)
		return
	}
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "예언자의 차례 진행중....")
	seerID := ""
	seerMsg := "버려진 직업들 중 2개 또는, 직업을 확인하고싶은 사람 한 명을 선택하세요" +
		"\n자신은 선택할 수 없어요\t(" + eBin + "): 버려진 직업들에서 고르기\n"
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
		wfd.GameLog += "\n예언자는 없었습니다."
		go func() {
			time.Sleep(20 * time.Second)
			wfd.TimingChan <- true
		}()
	}
	if seerID == "" {
		<-wfd.TimingChan
		nextStage(wfd)
		s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "예언자의 차례 종료.")
		return
	}
	uChan, _ := s.UserChannelCreate(seerID)
	msg, _ := s.ChannelMessageSend(uChan.ID, seerMsg)
	s.MessageReactionAdd(uChan.ID, msg.ID, eBin)
	for i := 0; i < len(wfd.UserIDs); i++ {
		s.MessageReactionAdd(uChan.ID, msg.ID, eNum[i])
	}

	<-wfd.TimingChan
	nextStage(wfd)
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "예언자의 차례 종료.")
}

func robberTask(s *discordgo.Session, wfd *WF.Data, stageMsg *discordgo.Message) {
	if exit(wfd.CurStage) {
		return
	}
	role := "강도"
	isIn := isInCheck(wfd, role)
	if !isIn {
		s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "강도는 넣지 않았습니다.")
		nextStage(wfd)
		return
	}
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "강도의 차례 진행중....")

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
		wfd.GameLog += "\n강도는 없었습니다."
		go func() {
			time.Sleep(20 * time.Second)
			wfd.TimingChan <- true
		}()
	}
	if robberID == "" {
		<-wfd.TimingChan
		nextStage(wfd)
		s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "강도의 차례 종료.")
		return
	}
	uChan, _ := s.UserChannelCreate(robberID)
	msg, _ := s.ChannelMessageSend(uChan.ID, robberMsg)
	for i := 0; i < len(wfd.UserIDs); i++ {
		s.MessageReactionAdd(uChan.ID, msg.ID, eNum[i])
	}

	<-wfd.TimingChan
	nextStage(wfd)
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "강도의 차례 종료.")

}

func tmTask(s *discordgo.Session, wfd *WF.Data, stageMsg *discordgo.Message) {
	if exit(wfd.CurStage) {
		return
	}
	role := "말썽쟁이"
	isIn := isInCheck(wfd, role)
	if !isIn {
		s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "말썽쟁이는 넣지 않았습니다.")
		nextStage(wfd)
		return
	}
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "말썽쟁이의 차례 진행중....")

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
		wfd.GameLog += "\n말썽쟁이는 없었습니다."
		go func() {
			time.Sleep(20 * time.Second)
			wfd.TimingChan <- true
		}()
	}
	if tmID == "" {
		<-wfd.TimingChan
		nextStage(wfd)
		s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "말썽쟁이의 차례 종료.")
		return
	}
	uChan, _ := s.UserChannelCreate(tmID)
	msg, _ := s.ChannelMessageSend(uChan.ID, tmMsg)
	for i := 0; i < len(wfd.UserIDs); i++ {
		s.MessageReactionAdd(uChan.ID, msg.ID, eNum[i])
	}

	<-wfd.TimingChan
	nextStage(wfd)
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "말썽쟁이의 차례 종료.")

}

func drunkTask(s *discordgo.Session, wfd *WF.Data, stageMsg *discordgo.Message) {
	if exit(wfd.CurStage) {
		return
	}
	role := "주정뱅이"
	isIn := isInCheck(wfd, role)
	if !isIn {
		s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "주정뱅이는 넣지 않았습니다.")
		nextStage(wfd)
		return
	}
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "주정뱅이의 차례 진행중....")

	drunkID := ""
	for _, item := range wfd.UserIDs {
		if wfd.UserRole[item] == "주정뱅이" {
			drunkID = item
		}
	}

	if drunkID == "" {
		wfd.GameLog += "\n주정뱅이는 없었습니다."
		go func() {
			time.Sleep(10 * time.Second)
			wfd.TimingChan <- true
		}()
		<-wfd.TimingChan
		nextStage(wfd)
		s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "주정뱅이의 차례 종료.")
		return
	}

	drunkChan, _ := s.UserChannelCreate(drunkID)
	sendDiscardsAddReaction(s, drunkChan)

	<-wfd.TimingChan
	nextStage(wfd)
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "주정뱅이의 차례 종료")
}

func sendDiscardsAddReaction(s *discordgo.Session, uChan *discordgo.Channel) {

	msg, _ := s.ChannelMessageSend(uChan.ID, "세 장의 비공개 직업 중 한 개를 선택하세요."+
		"\n< 1 > < 2 > < 3 >")
	s.MessageReactionAdd(uChan.ID, msg.ID, eOne)
	s.MessageReactionAdd(uChan.ID, msg.ID, eTwo)
	s.MessageReactionAdd(uChan.ID, msg.ID, eThree)
}

func insomniacTask(s *discordgo.Session, wfd *WF.Data, stageMsg *discordgo.Message) {
	if exit(wfd.CurStage) {
		return
	}
	role := "불면증환자"
	isIn := isInCheck(wfd, role)
	if !isIn {
		s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "불면증환자는 넣지 않았습니다.")
		nextStage(wfd)
		return
	}
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "불면증환자의 차례 진행중....")

	go func() {
		time.Sleep(10 * time.Second)
		wfd.TimingChan <- true
	}()

	inID := ""
	for _, item := range wfd.UserIDs {
		if wfd.UserRole[item] == "불면증환자" ||
			(wfd.DIFlag && wfd.UserRole[item] == "도플갱어") {
			inID = item
		}
	}

	if inID == "" {
		wfd.GameLog += "\n불면증환자는 없었습니다."
		<-wfd.TimingChan
		nextStage(wfd)
		s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "불면증환자의 차례 종료.")
		return
	}
	inUser, _ := s.User(inID)
	uChan, _ := s.UserChannelCreate(inID)
	s.ChannelMessageSend(uChan.ID, "다른 모든 사람이 능력을 쓴 후, 당신의 직업은 다음과 같습니다."+
		"\n직업: "+wfd.FinalRole[inID])

	wfd.GameLog += "\n불면증환자 `" + inUser.Username + "` (은)는 자신의 최종 직업\n`" +
		wfd.FinalRole[inID] + "` (을)를 확인하였습니다."
	<-wfd.TimingChan
	nextStage(wfd)
	s.ChannelMessageEdit(wfd.UseChannelID, stageMsg.ID, "불면증환자의 차례 종료.")

}

func getRoleInfo(role string) string {

	info := ""

	if role == "사냥꾼" {
		info = "당신은 노련한 사냥꾼입니다." +
			"\n당신이 늑대인간이라고 생각하는 한 사람에게 투표하세요." +
			"\n당신이 처형된다면, 길동무로 데려갈 수 있습니다."
	}
	if role == "도플갱어" {
		info = "처음으로 만난 사람의 직업을 복제합니다." +
			"\n복제한 능력을 사용할 수 있습니다." +
			"\n당신은 복제한 사람 편에 섭니다."
	}
	if role == "늑대인간" {
		info = "당신의 차례에 동료 늑대인간을 확인할 수 있습니다." +
			"\n만약 동료 늑대인간이 없다면," +
			"\n버려진 직업 3개 중 1개를 확인할 수 있습니다." +
			"\n하수인이 있다면, 당신을 몰래 도와줄 수도 있어요" +
			"\n마을 사람들을 혼란에 빠뜨리고 살아남으세요."
	}
	if role == "무두장이" {
		info = "당신은 죽기로 결심했죠." +
			"\n당신이 늑대인간인 것 처럼 연기하세요." +
			"\n처형된다면, 당신의 승리입니다."
	}
	if role == "마을주민" {
		info = "당신은 아무런 능력도 가지지 못했습니다." +
			"\n불안과 공포속에서 늑대인간을 찾아 처형하세요"
	}
	if role == "하수인" {
		info = "당신은 누가 늑대인간인지 알고 있어요." +
			"\n하지만 늑대인간들은 누가 하수인인지 모르죠." +
			"\n당신이 죽어도 늑대인간만 죽지 않는다면 승리입니다." +
			"\n모든 늑대인간이 죽지 않도록 도우세요."
	}
	if role == "프리메이슨" {
		info = "당신은 누가 동료 프리메이슨인지 확인합니다." +
			"\n만약 프리메이슨이 버려졌다면," +
			"\n다른 프리메이슨이 없음을 확인합니다." +
			"\n동료와 함께 늑대인간을 처형하세요."
	}
	if role == "예언자" {
		info = "당신은 버려진 3개의 직업들 중 2개를 보거나," +
			"\n다른 사람 하나의 직업을 볼 수 있습니다." +
			"\n예언이 밝혀준 곳을 따라 늑대인간을 찾아 처형하세요."
	}
	if role == "강도" {
		info = "당신은 누군가의 직업을 훔칠 수 있습니다." +
			"\n능력을 도둑맞은 사람은 강도가 되고," +
			"\n자신이 아직 원래 직업인줄 알 겁니다." +
			"\n훔친 능력에 맞게 누군가 처형하세요."
	}
	if role == "말썽쟁이" {
		info = "당신의 차례에 두 사람을 고릅니다." +
			"\n그 두 사람의 직업을 맞바꿉니다." +
			"\n말썽쟁이는 두 사람의 직업을 확인하지는 못합니다." +
			"\n혼란스럽겠지만, 늑대인간을 찾아 처형하세요."
	}
	if role == "불면증환자" {
		info = "당신은 잠이 든 지 얼마 지나지 않아 깨어났습니다." +
			"\n덕분에 당신은 당신이 무엇인지 알 수 있었죠." +
			"\n늑대인간을 처형하세요.." +
			"\n당신이 늑대인간이 되지 않았다면요."
	}
	if role == "주정뱅이" {
		info = "당신은 술에 잔뜩 취해 직업을 하나 주웠습니다." +
			"\n그치만 그 직업이 어떤 직업인지 기억이 안나요..." +
			"\n어쩌면 당신은 늑대인간일지도?"
	}

	return info
}

// 새로운 유저 등록시 수행
func newUserTask(m *discordgo.MessageCreate) {
	wfd := wfDataMap[m.GuildID]
	if len(wfd.UserIDs) >= wfd.MaxUser {
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
	wfd.UserIDs = make([]string, 0, wfd.MaxUser)
	wfd = WF.NewWFData("", "")
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
			s.ChannelMessageSend(m.ChannelID, "저장된 정보가 없었습니다.")
			return
		}
		hwUserIDs = make([]string, 0, 10)
		msg := "모든 입력된 데이터들을 삭제하였습니다."
		_, _ = s.ChannelMessageSend(m.ChannelID, msg)
	}
}
