package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/bwmarrin/discordgo"
)

var (
	Token                                   string
	Now                                     = time.Now()
	DecID                                   = 184763577729548289
	CalID                                   = 428895427513942019
	thenString, Day, Month, Year, dayString string
	hours, mins, secs                       float64
	then                                    time.Time
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

// func setDate () {

// }

// func getTime(content string) {

// }
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	Now = time.Now()
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If Author is Declan or me
	if m.Author.ID == strconv.Itoa(DecID) || m.Author.ID == strconv.Itoa(CalID) {
		if strings.Contains(m.Content, "!set") {

			// Parse out day, month, year   Strict 01/01/24 format required
			thenString = m.Content[5:]
			// If length isn't 8, format incorrect
			if len(thenString) != 8 {
				fmt.Println("Invaid date format! DD/MM/YY")
				s.ChannelMessageSend(m.ChannelID, "Date does not follow DD/MM/YY")
				return
			}
			//Init rune arrays
			runeArray := []rune(thenString)
			failArray := []rune{}

			for i := 0; i < 7; i++ {
				if unicode.IsLetter(runeArray[i]) {
					failArray = append(failArray, runeArray[i])
				}
			}

			//If fail array contains letters, error.
			if len(failArray) > 0 {
				fmt.Printf("Date %s contains letters %d", thenString, failArray)
				s.ChannelMessageSend(m.ChannelID, "Date contains letters!")
				return
			}
			//Slice date
			Day = thenString[0:2]
			Month = thenString[3:5]
			Year = "20" + thenString[6:8]
			// Convert to Int
			dy, err := strconv.Atoi(Year)
			dm, err := strconv.Atoi(Month)
			dd, err := strconv.Atoi(Day)

			if err != nil {
				fmt.Println("Can't convert this to an int!")
			}

			if dd > 31 || dm > 12 {
				s.ChannelMessageSend(m.ChannelID, "Date is invalid!")
				return
			}

			then = time.Date(dy, time.Month(dm), dd, 0, 0, 0, 0, time.UTC)
			// If date is in the past
			if then.Sub(Now) < 1 {
				s.ChannelMessageSend(m.ChannelID, "Please pick a date in the future")
				return
			}
		}
	}

	// If message contains !dec
	if strings.Contains(m.Content, "!dec") {
		// If then is empty, !set not called yet
		if time.Time.IsZero(then) {
			s.ChannelMessageSend(m.ChannelID, "Declan hasn't set a return date!")
			return
		}

		diff := then.Sub(Now)

		hours = diff.Hours()
		mins = diff.Minutes()
		secs = diff.Seconds()

		days := hours / 24

		dr := strconv.FormatFloat(days, 'f', 0, 64)
		hr := strconv.FormatFloat(hours, 'f', 0, 64)

		if days > 2 {
			dayString = "days"
		} else {
			dayString = "day"
		}

		s.ChannelMessageSend(m.ChannelID, "Dec returns in "+hr+" Hours, or "+dr+" "+dayString+" 🚢")

	}

	// if strings.Contains(m.Content, "@Lethal Company") {
	// 	s.ChannelMessageSend(m.ChannelID,)
	// }

}

func main() {

	// var Token string
	// Token = os.Getenv("TOKEN")

	ds, err := discordgo.New("Bot " + Token)

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	ds.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	ds.Identify.Intents = discordgo.IntentsGuildMessages

	err = ds.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	ds.Close()
}
