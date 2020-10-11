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
	prefix = "ㅁ"
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
	classList = []string{
		"늑대인간",
		"하수인",
		"예언자",
		"강도",
		"말썽쟁이",
		"주정뱅이",
		"무두장이",
		"마을주민",
		"불면증환자",
	}
	gameSeq = []string{
		"늑대인간",
		"하수인",
		"예언자",
		"강도",
		"말썽쟁이",
		"주정뱅이",
		"불면증환자",
	}
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

	cardMap = make(map[string][]string)
	prevSettingMap = make(map[string]*WF.SettingData)
	uidToGid = make(map[string]string)
	isGuildIn = make(map[string]bool)
	isUserIn = make(map[string]bool)
	wfDataMap = make(map[string]*WF.Data)
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

func nextStage(wfd *WF.Data) {
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
				s.ChannelMessageSend(uChan.ID, "<1번: `"+wfd.CardDeck.Cards[0]+"` >")
			}
			if r.Emoji.Name == eTwo {
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				uChan, _ := s.UserChannelCreate(r.UserID)
				go func() {
					time.Sleep(10 * time.Second)
					wfd.TimingChan <- true
				}()
				wfd.GameLog += "\n유일한 늑대인간은 버려진 `" + wfd.CardDeck.Cards[1] + "` 를 확인하였습니다."
				s.ChannelMessageSend(uChan.ID, "<2번: `"+wfd.CardDeck.Cards[1]+"` >")
			}
			if r.Emoji.Name == eThree {
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				uChan, _ := s.UserChannelCreate(r.UserID)
				go func() {
					time.Sleep(10 * time.Second)
					wfd.TimingChan <- true
				}()
				wfd.GameLog += "\n유일한 늑대인간은 버려진 `" + wfd.CardDeck.Cards[2] + "` 를 확인하였습니다."
				s.ChannelMessageSend(uChan.ID, "<3번: `"+wfd.CardDeck.Cards[2]+"` >")
			}
		}
	}
	if wfd.CurStage == "예언자" {
		if wfd.UserRole[r.UserID] == "예언자" {
			if r.Emoji.Name == eBin {
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				wfd.CurStage = "예언자_trash"
				msg, _ := s.ChannelMessageSend(r.ChannelID, "보지 않을 직업을 고르세요.\n< 1 > < 2 > < 3 >")
				s.MessageReactionAdd(r.ChannelID, msg.ID, eOne)
				s.MessageReactionAdd(r.ChannelID, msg.ID, eTwo)
				s.MessageReactionAdd(r.ChannelID, msg.ID, eThree)
			}
			for i := 0; i < len(wfd.UserIDs); i++ {
				if r.Emoji.Name == eNum[i] {
					if wfd.UserIDs[i] != r.UserID {
						seer, _ := s.User(r.UserID)
						s.ChannelMessageDelete(r.ChannelID, r.MessageID)
						wfd.CurStage = "예언자_used_power"
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
	}
	if wfd.CurStage == "예언자_trash" {
		if wfd.UserRole[r.UserID] == "예언자" {
			seer, _ := s.User(r.UserID)
			trashMsg := ""
			wfd.GameLog += "\n`예언자` `" + seer.Username + "` (은)는 버려진 직업 "
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
	}
	if wfd.CurStage == "강도" {
		if wfd.UserRole[r.UserID] == "강도" {
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
	}
	if wfd.CurStage == "말썽쟁이_oneMoreChoice" {
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
	}
	if wfd.CurStage == "말썽쟁이" {
		if wfd.UserRole[r.UserID] == "말썽쟁이" {
			tm, _ := s.User(r.UserID)
			tmMsg := ""
			for i := 0; i < len(wfd.UserIDs); i++ {
				if r.Emoji.Name == eNum[i] {
					if wfd.UserIDs[i] != r.UserID {
						wfd.CurStage = "말썽쟁이_choiceWaiting"
						wfd.IndexChan <- i

						s.ChannelMessageSend(r.ChannelID, "다음 사람을 고르세요")

						wfd.CurStage = "말썽쟁이_oneMoreChoice"
						s.ChannelMessageDelete(r.ChannelID, r.MessageID)
						user, _ := s.User(wfd.UserIDs[i])
						selectMsg := "`" + user.Username + "`님을 선택하였습니다."
						wfd.GameLog += "\n말썽쟁이 `" + tm.Username +
							"` (은)는 `" + wfd.FinalRole[wfd.UserIDs[i]] + "` 인 `" +
							user.Username + "` 의 직업과, "
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
	if wfd.CurStage == "주정뱅이" {
		if wfd.UserRole[r.UserID] == "주정뱅이" {
			dr, _ := s.User(r.UserID)
			for i := 0; i < 3; i++ {
				if r.Emoji.Name == eNum[i] {
					s.ChannelMessageDelete(r.ChannelID, r.MessageID)
					wfd.GameLog += "\n주정뱅이 `" + dr.Username +
						"` 는 버려진 직업 중 `" + wfd.CardDeck.Cards[i] + "` (와)과" +
						" 자신의 직업 `주정뱅이` 를 맞바꾸었습니다."
					temp := wfd.CardDeck.Cards[i]
					wfd.CardDeck.Cards[i] = "주정뱅이"
					wfd.FinalRole[dr.ID] = temp
					s.ChannelMessageSend(r.ChannelID, "술에 취한 당신은, "+
						strconv.Itoa(i+1)+"번 직업와 맞바꾸었습니다."+
						"\n이런... 술에 취해 무슨 직업이었는지도 잊어버렸군요..")
					wfd.TimingChan <- true
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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, prefix) { // 프리픽스로 시작하는 메시지일 경우
		var wfd *WF.Data

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
			wfDataMap[m.GuildID].CurStage = "Prepare_card"
			wfd = wfDataMap[m.GuildID]
			cardSetting(s, m.GuildID, wfd)
			<-wfd.TimingChan
			prevSettingMap[m.GuildID] = WF.NewSettingData(wfd.CardDeck, wfd.MaxUser)
			cardMsg := "> 직업 설정이 완료되었습니다. 설정된 직업들은 다음과 같습니다."
			for _, item := range wfd.CardDeck.Cards {
				cardMsg += "\n" + item
			}
			cardMsg += "\n**총 " + strconv.Itoa(len(wfd.CardDeck.Cards)) +
				"개의 직업이 선정되었습니다. 총 플레이 인원은 " +
				strconv.Itoa(len(wfd.CardDeck.Cards)-3) + "명 입니다.**"
			s.ChannelMessageSend(wfd.UseChannelID, cardMsg)
			wfd.CardDeck.ShuffleCards()
			newUserTask(m)
			wfDataMap[m.GuildID].CurStage = "Prepare"
			s.ChannelMessageSend(m.ChannelID, "> 게임 시작!\n> `"+prefix+"입장` 으로 입장하세요")
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
				if m.Content == prefix+"늑대인간" {
					s.ChannelMessageSend(wfd.UseChannelID, "늑대인간은 2명이 최대입니다.")
				}
				if m.Content == prefix+"하수인" {
					for _, item := range cardMap[m.GuildID] {
						if item == "하수인" {
							s.ChannelMessageSend(wfd.UseChannelID, "하수인은 최대 1장입니다.")
							return
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "하수인")
					s.ChannelMessageSend(wfd.UseChannelID, "하수인을 넣었습니다.")
				}
				if m.Content == prefix+"예언자" {
					for _, item := range cardMap[m.GuildID] {
						if item == "예언자" {
							s.ChannelMessageSend(wfd.UseChannelID, "예언자은 최대 1장입니다.")
							return
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "예언자")
					s.ChannelMessageSend(wfd.UseChannelID, "에언자를 넣었습니다.")
				}
				if m.Content == prefix+"말썽쟁이" {
					for _, item := range cardMap[m.GuildID] {
						if item == "말썽쟁이" {
							s.ChannelMessageSend(wfd.UseChannelID, "말썽쟁이은 최대 1장입니다.")
							return
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "말썽쟁이")
					s.ChannelMessageSend(wfd.UseChannelID, "말썽쟁이를 넣었습니다.")
				}
				if m.Content == prefix+"무두장이" {
					for _, item := range cardMap[m.GuildID] {
						if item == "무두장이" {
							s.ChannelMessageSend(wfd.UseChannelID, "무두장이은 최대 1장입니다.")
							return
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "무두장이")
					s.ChannelMessageSend(wfd.UseChannelID, "무두장이를 넣었습니다.")
				}
				if m.Content == prefix+"불면증환자" {
					for _, item := range cardMap[m.GuildID] {
						if item == "불면증환자" {
							s.ChannelMessageSend(wfd.UseChannelID, "불면증환자은 최대 1장입니다.")
							return
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "불면증환자")
					s.ChannelMessageSend(wfd.UseChannelID, "불면증환자를 넣었습니다.")
				}
				if m.Content == prefix+"강도" {
					for _, item := range cardMap[m.GuildID] {
						if item == "강도" {
							s.ChannelMessageSend(wfd.UseChannelID, "강도은 최대 1장입니다.")
							return
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "강도")
					s.ChannelMessageSend(wfd.UseChannelID, "강도를 넣었습니다.")
				}
				if m.Content == prefix+"마을주민" {
					count := 0
					for _, item := range cardMap[m.GuildID] {
						if item == "마을주민" {
							count++
							if count == 3 {
								s.ChannelMessageSend(wfd.UseChannelID, "마을주민은 최대 3장입니다.")
								return
							}
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "마을주민")
					s.ChannelMessageSend(wfd.UseChannelID, "마을주민을 넣었습니다.")
				}
				if m.Content == prefix+"주정뱅이" {
					count := 0
					for _, item := range cardMap[m.GuildID] {
						if item == "주정뱅이" {
							count++
							if count == 1 {
								s.ChannelMessageSend(wfd.UseChannelID, "주정뱅이는 최대 3장입니다.")
								return
							}
						}
					}
					cardMap[m.GuildID] = append(cardMap[m.GuildID], "주정뱅이")
					s.ChannelMessageSend(wfd.UseChannelID, "주정뱅이를 넣었습니다.")
				}
				if m.Content == prefix+"직업설정 완료" {

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
							wfd.GameLog += "`" + user.Username + "` 님의 역할이 `" + wfd.UserRole[uid] + "` (으)로 배정되었습니다.\n"
							roleBrief := "> **당신의 역할은 **`" + wfd.UserRole[uid] + "`**입니다.**\n"
							roleBrief += getRoleInfo(wfd.UserRole[uid])
							s.ChannelMessageSend(uChan.ID, roleBrief)
						}(item)
					}

					nextStage(wfd)
					werewolfTask(s, wfd)
					minionTask(s, wfd)
					seerTask(s, wfd)
					robberTask(s, wfd)
					tmTask(s, wfd)
					drunkTask(s, wfd)
					insomniacTask(s, wfd)
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

func cardSetting(s *discordgo.Session, gid string, wfd *WF.Data) {
	cardMap[gid] = make([]string, 0, 10)

	s.ChannelMessageSend(wfd.UseChannelID, "> 직업 설정을 시작합니다.")
	s.ChannelMessageSend(wfd.UseChannelID, "이전 설정과 동일한 직업 설정을 사용할까요?\n("+prefix+"ㅇㅇ/"+prefix+"ㄴㄴ)")

	choice := <-wfd.CardDeck.ChoiceChan
	if choice == 0 {
		if prevSettingMap[gid] != nil {
			wfd.CardDeck.Cards = prevSettingMap[gid].CardDeck.Cards
			wfd.CardDeck.ShuffleCards()
			wfd.MaxUser = prevSettingMap[gid].MaxUser
			wfd.TimingChan <- true
			return
		} else {
			s.ChannelMessageSend(wfd.UseChannelID, "> 이전 게임 기록이 남아있지 않습니다.."+
				"\n> 게임을 한 적이 없거나, 서버가 재부팅되었을 수 있습니다.")
		}
	}
	cardMap[gid] = append(cardMap[gid], "늑대인간")
	cardMap[gid] = append(cardMap[gid], "늑대인간")
	s.ChannelMessageSend(wfd.UseChannelID, "늑대인간 2장은 필수입니다. 직업 덱에 넣었습니다.")
	for true {
		wfd.CurStage = "Prepare_class"
		s.ChannelMessageSend(wfd.UseChannelID, "추가할 직업들을 입력하세요. (ex: ㅁ마을주민)"+
			"\n모두 입력한 후 `"+prefix+"직업설정 완료` 로 다음 단계로 넘어가세요.")
		classMsg := "구현된 직업 목록:"
		for _, item := range classList {
			classMsg += " " + item
		}
		s.ChannelMessageSend(wfd.UseChannelID, classMsg)
		<-wfd.CardDeck.ChoiceChan
		if len(cardMap[gid]) < 6 {
			s.ChannelMessageSend(wfd.UseChannelID, "6장 이상을 골라야 합니다..\n("+strconv.Itoa(len(cardMap[gid]))+"/6)")
		} else {
			wfd.CardDeck.Cards = cardMap[gid]
			wfd.MaxUser = len(cardMap[gid]) - 3
			break
		}
	}

	wfd.TimingChan <- true
}

func dayBriefTask(s *discordgo.Session, wfd *WF.Data) {
	briefMsg := ""

	briefMsg += "> 모든 특수 능력 사용이 끝났습니다." +
		"\n3초 후 여러분들에게 각자의 투표 용지가 전송됩니다." +
		"\n한번 투표한 내용은 바꿀 수 없기에, 신중하게 투표하세요" +
		"\n"
	go func() {
		time.Sleep(time.Second * 3)
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
		msg, _ := s.ChannelMessageSend(cid, briefMsg)
		for i := 0; i < len(wfd.UserIDs); i++ {
			s.MessageReactionAdd(cid, msg.ID, eNum[i])
		}
	}
	electFinishTask(s, wfd)
}

func electFinishTask(s *discordgo.Session, wfd *WF.Data) {
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
		s.ChannelMessageSend(wfd.UseChannelID, electResultMsg)
		time.Sleep(3 * time.Second)
		s.ChannelMessageSend(wfd.UseChannelID, "> 게임 로그:")
		s.ChannelMessageSend(wfd.UseChannelID, wfd.GameLog)
		return
	}
	electResultMsg += "> "

	for i, item := range electResult {
		if item == max {
			user, _ := s.User(wfd.UserIDs[i])
			electResultMsg += " `" + user.Username + "`"
		}
	}
	electResultMsg += " 님이 총 " + strconv.Itoa(electResult[maxi]) + " 표로 처형되었습니다."
	s.ChannelMessageSend(wfd.UseChannelID, electResultMsg)
	s.ChannelMessageSend(wfd.UseChannelID, "처형된 사람의 직업은..?")
	for i := 0; i < 3; i++ {
		s.ChannelMessageSend(wfd.UseChannelID, "...")
		time.Sleep(time.Second)
	}
	electResultMsg = "> 처형된 사람의 직업은"
	for i, item := range electResult {
		if item == max {
			user, _ := s.User(wfd.UserIDs[i])
			electResultMsg += "\n<`" + user.Username + "` : `" +
				wfd.UserRole[wfd.UserIDs[i]] + "`-> `" +
				wfd.FinalRole[wfd.UserIDs[i]] + "`>"
		}
	}
	electResultMsg += " 입니다."
	s.ChannelMessageSend(wfd.UseChannelID, electResultMsg)
	electResultMsg = "> 모두의 최종 직업:"
	for _, item := range wfd.UserIDs {
		user, _ := s.User(item)
		electResultMsg += "\n<`" + user.Username + "` : `" + wfd.FinalRole[item] + "`>"
	}
	s.ChannelMessageSend(wfd.UseChannelID, electResultMsg)
	time.Sleep(3 * time.Second)
	s.ChannelMessageSend(wfd.UseChannelID, "> 게임 로그:")
	s.ChannelMessageSend(wfd.UseChannelID, wfd.GameLog)
}

func werewolfTask(s *discordgo.Session, wfd *WF.Data) {
	s.ChannelMessageSend(wfd.UseChannelID, "늑대인간의 차례입니다.")
	wolvesID := make([]string, 0, 10)
	for _, item := range wfd.UserIDs {
		if wfd.UserRole[item] == "늑대인간" {
			wolvesID = append(wolvesID, item)
		}
	}
	if len(wolvesID) != 1 {
		go func() {
			wfd.GameLog += "`늑대인간`이 " + strconv.Itoa(len(wolvesID)) + " 명이라 서로를 확인만 합니다.\n"
			time.Sleep(10 * time.Second)
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
	nextStage(wfd)
	s.ChannelMessageSend(wfd.UseChannelID, "늑대인간의 차례 종료.")
}

func minionTask(s *discordgo.Session, wfd *WF.Data) {
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

	if minionID != "" {
		wfd.GameLog += "하수인의 차례: 하수인은 늑대인간이 " + strconv.Itoa(len(wolvesID)) + " 명인 것을 확인했습니다.\n"
	}
	user, err := s.User(minionID)
	if err != nil {
		<-wfd.TimingChan
		nextStage(wfd)
		s.ChannelMessageSend(wfd.UseChannelID, "하수인의 차례 종료.")
		return
	}

	uChan, _ := s.UserChannelCreate(user.ID)
	s.ChannelMessageSend(uChan.ID, minionMsg)
	<-wfd.TimingChan
	nextStage(wfd)
	s.ChannelMessageSend(wfd.UseChannelID, "하수인의 차례 종료.")
}

func seerTask(s *discordgo.Session, wfd *WF.Data) {
	s.ChannelMessageSend(wfd.UseChannelID, "예언자의 차례입니다.")
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
	nextStage(wfd)
	s.ChannelMessageSend(wfd.UseChannelID, "예언자의 차례 종료.")
}

func robberTask(s *discordgo.Session, wfd *WF.Data) {
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
		wfd.GameLog += "\n강도는 없었습니다."
		go func() {
			time.Sleep(20 * time.Second)
			wfd.TimingChan <- true
		}()
	}
	if robberID == "" {
		<-wfd.TimingChan
		nextStage(wfd)
		s.ChannelMessageSend(wfd.UseChannelID, "강도의 차례 종료.")
		return
	}
	uChan, _ := s.UserChannelCreate(robberID)
	msg, _ := s.ChannelMessageSend(uChan.ID, robberMsg)
	for i := 0; i < len(wfd.UserIDs); i++ {
		s.MessageReactionAdd(uChan.ID, msg.ID, eNum[i])
	}

	<-wfd.TimingChan
	nextStage(wfd)
	s.ChannelMessageSend(wfd.UseChannelID, "강도의 차례 종료.")

}

func tmTask(s *discordgo.Session, wfd *WF.Data) {
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
		wfd.GameLog += "\n말썽쟁이는 없었습니다."
		go func() {
			time.Sleep(20 * time.Second)
			wfd.TimingChan <- true
		}()
	}
	if tmID == "" {
		<-wfd.TimingChan
		nextStage(wfd)
		s.ChannelMessageSend(wfd.UseChannelID, "말썽쟁이의 차례 종료.")
		return
	}
	uChan, _ := s.UserChannelCreate(tmID)
	msg, _ := s.ChannelMessageSend(uChan.ID, tmMsg)
	for i := 0; i < len(wfd.UserIDs); i++ {
		s.MessageReactionAdd(uChan.ID, msg.ID, eNum[i])
	}

	<-wfd.TimingChan
	nextStage(wfd)
	s.ChannelMessageSend(wfd.UseChannelID, "말썽쟁이의 차례 종료.")

}

func drunkTask(s *discordgo.Session, wfd *WF.Data) {
	s.ChannelMessageSend(wfd.UseChannelID, "주정뱅이의 차례입니다.")

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
		s.ChannelMessageSend(wfd.UseChannelID, "주정뱅이의 차례 종료.")
		return
	}

	drunkChan, _ := s.UserChannelCreate(drunkID)
	msg, _ := s.ChannelMessageSend(drunkChan.ID, "세 장의 비공개 직업 중 한 개를 선택하세요."+
		"\n< 1 > < 2 > < 3 >")
	s.MessageReactionAdd(drunkChan.ID, msg.ID, eOne)
	s.MessageReactionAdd(drunkChan.ID, msg.ID, eTwo)
	s.MessageReactionAdd(drunkChan.ID, msg.ID, eThree)

	<-wfd.TimingChan
	nextStage(wfd)
	s.ChannelMessageSend(wfd.UseChannelID, "주정뱅이의 차례 종료")
}

func insomniacTask(s *discordgo.Session, wfd *WF.Data) {
	s.ChannelMessageSend(wfd.UseChannelID, "불면증환자의 차례입니다.")

	go func() {
		time.Sleep(10 * time.Second)
		wfd.TimingChan <- true
	}()

	inID := ""
	for _, item := range wfd.UserIDs {
		if wfd.UserRole[item] == "불면증환자" {
			inID = item
		}
	}

	if inID == "" {
		wfd.GameLog += "\n불면증환자는 없었습니다."
		<-wfd.TimingChan
		nextStage(wfd)
		s.ChannelMessageSend(wfd.UseChannelID, "불면증환자의 차례 종료.")
		return
	}
	inUser, _ := s.User(inID)
	uChan, _ := s.UserChannelCreate(inID)
	s.ChannelMessageSend(uChan.ID, "다른 모든 사람이 능력을 쓴 후, 당신의 역할은 다음과 같습니다."+
		"\n역할: "+wfd.FinalRole[inID])

	wfd.GameLog += "\n불면증환자 `" + inUser.Username + "` (은)는 자신의 최종 직업 `" +
		wfd.FinalRole[inID] + "` (을)를 확인하였습니다."
	<-wfd.TimingChan
	nextStage(wfd)
	s.ChannelMessageSend(wfd.UseChannelID, "불면증환자의 차례 종료.")

}

func getRoleInfo(role string) string {
	info := ""

	if role == "늑대인간" {
		info = "당신은 게임이 시작된 후에 동료 늑대인간을 확인할 수 있습니다." +
			"\n만약 동료 늑대인간이 없다면, 버려진 직업들 3개 중 1개를 확인할 수 있습니다." +
			"\n마을 사람들을 혼란에 빠뜨리세요. 능력이 사라지지 않았다면 말이죠."
	}
	if role == "무두장이" {
		info = "당신은 죽기로 결심했죠." +
			"\n당신이 늑대인간인 것 처럼 연기하세요." +
			"\n처형된다면, 당신의 승리입니다."
	}
	if role == "마을주민" {
		info = "당신은 아무런 능력도 가지지 못했습니다." +
			"\n불안과 공포속에서 늑대인간을 찾아서 처형하세요"
	}
	if role == "하수인" {
		info = "당신은 누가 늑대인간인지 알고 있어요." +
			"\n그들을 도와 모든 늑대인간이 처형당하지 않도록 하세요."
	}
	if role == "예언자" {
		info = "당신은 버려진 3개의 직업들 중 2개를 보거나," +
			"\n다른 사람 하나의 능력을 볼 수 있습니다." +
			"\n예언이 밝혀준 곳을 따라 늑대인간을 찾아서 처형하세요."
	}
	if role == "강도" {
		info = "당신은 무언갈 훔쳤습니다. 그것은 물건이 아닌 능력이죠." +
			"\n늑대인간의 능력을 훔쳤다면, 당신은 늑대인간이 될 것이고," +
			"\n그 늑대인간은 자신이 아직도 늑대인간인 줄 알 겁니다." +
			"\n훔친 능력을 확인하고, 훔친 능력에 맞게 늑대인간을 처형할지 말지 판단하세요"
	}
	if role == "말썽쟁이" {
		info = "당신도 모르는 새에 두 사람의 능력을 바꾸어버리다니.. 말도 안돼죠." +
			"\n당신은 하지만 그런 능력을 갖고 있어요. 그래도 두 사람이 무슨 능력이 있었는지는 알 수 없습니다." +
			"\n혼란스럽겠지만, 늑대인간을 찾아 처형하세요."
	}
	if role == "불면증환자" {
		info = "당신은 겨우 잠이 들었지만 얼마 지나지 않아 깨어났습니다." +
			"\n덕분에 당신은 어떤 새로운 힘이 생겼는지 알 수 있었죠." +
			"\n늑대인간을 처형하세요, 당신이 늑대인간이 되지 않았다면요."
	}
	if role == "주정뱅이" {
		info = "당신은 술에 잔뜩 취해 어젯밤 일을 기억하지 못합니다.." +
			"\n어쩌면 당신은 늑대인간이었을지도?"
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
