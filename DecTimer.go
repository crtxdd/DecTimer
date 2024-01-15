package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	Now = time.Now()
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If Author is Declan
	if m.Author.ID == strconv.Itoa(DecID) || m.Author.ID == strconv.Itoa(CalID) {
		if strings.Contains(m.Content, "!set") {
			// Parse out day, month, year   Strict 01/01/24 format required

			thenString = m.Content[5:]

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
				s.ChannelMessageSend(m.ChannelID, "Invalid date")
				return
			}

			then = time.Date(dy, time.Month(dm), dd, 0, 0, 0, 0, time.UTC)
			// If date is in the past
			if then.Sub(Now) < 1 {
				s.ChannelMessageSend(m.ChannelID, "Please pick a date in the future")
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

		s.ChannelMessageSend(m.ChannelID, "Dec returns in "+hr+" Hours, or "+dr+" "+dayString)

	}

}

func main() {

	Token = "MTE5NjUyNDg5NDUxNTY0MjQ3OA.GmsvEd.IsH7FKO8UtXn4Ih3tSObbINKeGMlsm6LuWJXrw"

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
