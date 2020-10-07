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
						roleBrief := "당신의 역할은 `" + wfd.UserRole[item] + "`입니다."
						s.ChannelMessageSend(uChan.ID, roleBrief)
					}
					wfd.CurStage = "Night"
				} else {
					s.ChannelMessageSend(m.ChannelID, "정확한 인원이 모이지 않았습니다. ("+strconv.Itoa(len(wfd.UserIDs))+"/"+strconv.Itoa(MaxUser)+")")
				}
			}
		}
	}

	homeworkMethod(s, m)
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
