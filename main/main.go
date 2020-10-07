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

	"github.com/bwmarrin/discordgo"
)

const (
	MaxUser = 6
	prefix  = "="
)

var (
	hwUserIDs []string
)

// Variables used for command line parameters
var (
	Token     string
	isGuildIn map[string]bool
	isUserIn  map[string]bool
	uidToGid  map[string]string
	wfDataMap map[string]*WF.WFData
)

func init() {
	hwUserIDs = make([]string, 0, 10)

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
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

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

		if m.Content == prefix+"시작" {
			isGuildIn[m.GuildID] = true
			wfDataMap[m.GuildID] = WF.NewWFData(m.Author.ID, m.ChannelID)
			newUserTask(m)
			_, _ = s.ChannelMessageSend(m.ChannelID, "게임 시작!\n`"+prefix+"입장` 으로 입장하세요")
		}
		if isGuildIn[m.GuildID] {
			wfd = wfDataMap[m.GuildID]
			if m.Content == prefix+"입장" {
				if isUserIn[m.Author.ID] {
					_, _ = s.ChannelMessageSend(m.ChannelID, "이미 입장한 유저입니다.")
					return
				}
				newUserTask(m)
				s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+"님이 입장하셨습니다.")
			}
			if m.Author.ID == wfd.AdminUserID && m.Content == prefix+"취소" {
				cancelGameTask(m)
				s.ChannelMessageSend(m.ChannelID, "게임이 취소되었습니다.")
			}
			if m.Author.ID == wfd.AdminUserID && strings.HasPrefix(m.Content, prefix+"더미추가") {
				sepMsg := strings.Split(m.Content, " ")
				if len(sepMsg) == 1 {
					s.ChannelMessageSend(m.ChannelID, "추가할 인원 숫자를 입력하세요.\t<usage: =더미추가 num>")
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
			if strings.HasPrefix(m.Content, prefix+"마감") {
				if len(wfd.UserIDs) == 6 {
					for _, item := range wfd.UserIDs {
						uChan, _ := s.UserChannelCreate(item)
						roleBrief := "> **당신의 역할은 **`" + wfd.UserRole[item] + "`**입니다.**\n"
						roleBrief += getRoleInfo(wfd.UserRole[item])
						s.ChannelMessageSend(uChan.ID, roleBrief)
					}

				} else {
					s.ChannelMessageSend(m.ChannelID, "정확한 인원이 모이지 않았습니다. ("+strconv.Itoa(len(wfd.UserIDs))+"/"+strconv.Itoa(MaxUser)+")")
				}
			}
		}
	}

	homeworkMethod(s, m)
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
	if role == "마을시민" {
		info = "당신은 잔뜩 겁에 질린 평범한 마을 시민입니다." +
			"\n당신은 아무런 능력도 가지지 못했습니다." +
			"\n단지 당신은 늑대인간을 찾아내어 처형하세요, 불안과 공포에서 해방될 수 있을거에요." +
			"\n늑대인간을 한명이라도 찾아내어 처형한다면 승리입니다." +
			"\n하수인이 죽더라도 늑대인간이 남는다면 오늘 밤 당신은 아마.... 끔찍하죠." +
			"\n하지만 늑대인간과 그들의 하수인이 한명도 없을 수도 있어요." +
			"\n그 경우에는 아무도 죽어서는 안됩니다." +
			"\n**늑대인간들은 당신을 혼란스럽게 만들겁니다. 행운을 빕니다.**"
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
