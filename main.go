package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const token = "Bot OTUzNzkyNjk3MTQ2MTU1MDg4.YjJuqw.ghRikYLa898p1gD6JlkO27uWpv8"
const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data
)

//Messages Log
const (
	online              = "Bot Online"
	failLogin           = "Erro ao se conectar com Discord API."
	failVoiceConnection = "Falha ao se conectar ao canal de voz."
)

func main() {
	fmt.Println("Iniciando bot.")
	bot, err := discordgo.New(token)
	if err != nil {
		panic(err)
	}
	//new connection websocket
	err = bot.Open()
	defer bot.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println(online)
	setHandlers(bot)
	for {
		time.Sleep(time.Second)
	}
}

func setHandlers(bot *discordgo.Session) {
	bot.AddHandler(reciveMessage)
}

func reciveMessage(session *discordgo.Session, message *discordgo.MessageCreate) {

	if !message.Author.Bot {
		if !strings.HasPrefix(message.Content, "!") {
			return
		}
		if strings.Count(message.Content, "!") > 1 {
			session.ChannelMessageSend(message.ChannelID, "Desculpe, ainda apenas um comando é validado por mensagem.\nExperimente quebrar seu comando em várias mensagens.")
			return
		}
		switch strings.ToLower(message.Content) {
		case "!pruai":
			session.ChannelMessageSend(message.ChannelID, "tocando....")
			for _, g := range session.State.Guilds {
				for _, v := range g.VoiceStates {
					if v.UserID == message.Author.ID {
						channelVoice, err := session.ChannelVoiceJoin(message.GuildID, v.ChannelID, false, false)
						if err != nil {
							session.ChannelMessageSend(message.ChannelID, failVoiceConnection)
							return
						}
						if !channelVoice.Ready {
							session.ChannelMessageSend(message.ChannelID, "canal não esta pronto, tente novamente.")
							return
						}
						session.ChannelMessageSend(message.ChannelID, "entrou no canal. Iniciando stream...")
						time.Sleep(time.Second * 5)
						var bufferSize = 40
						var buffer = make([]byte, bufferSize)
						var file *os.File
						var readBytes int
						var streamDone bool
						channelVoice.Speaking(true)
						file, err = os.Open("./X2Download.com-Joji-Yeah-Right-_LEGENDADO_TRADUCAO_.opus")
						defer file.Close()
						if err != nil {
							panic("unable to read the file")
						}
						for !streamDone {
							readBytes, err = file.Read(buffer)
							if err != nil {
								if err == io.EOF {
									println("stream done.")
									streamDone = true
									channelVoice.Speaking(false)
									channelVoice.Close()
									channelVoice.Disconnect()
									session.ChannelMessageSend(message.ChannelID, "stream node.")
									continue
								}
							}
							channelVoice.OpusSend <- buffer[:readBytes]
							//println(buffer[:readBytes	])
						}
					}
				}
			}

		case "!goroutines":
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Total de Goroutines: **%v**\nTotal de núcloes de processamento: **%v**", runtime.NumGoroutine(), runtime.NumCPU()))
		case "!ping":
			var times []time.Duration
			var times_ms []int = []int{10, 20, 30}
			var ping time.Duration
			var average_ms int
			msgPing, err := session.ChannelMessageSend(message.ChannelID, "Iniciando teste de latência..")
			defer session.ChannelMessageDelete(msgPing.ChannelID, msgPing.Reference().MessageID)
			for i := 1; i < 30; i++ {
				time.Sleep(time.Second)
				newPing := time.Now()
				if err != nil {
					session.ChannelMessageSend(message.ChannelID, "Ocorreu um erro ao efetuar a consulta de latência.")
					break
				}
				session.ChannelMessageEdit(msgPing.ChannelID, msgPing.Reference().MessageID, fmt.Sprintf("Enviando pacotes: %v/30", i))
				ping = time.Since(newPing)
				times = append(times, ping)
			}
			for _, p := range times {
				average_ms += int(p.Milliseconds())
				times_ms = append(times_ms, int(p.Milliseconds()))
			}
			average_ms = average_ms / len(times)
			sort.Ints(times_ms)

			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Velocidade de lantência:\nLatência Minima: **%vms**, Latência Média: **%vms**, Latência Máxima: **%vms**", times_ms[0], average_ms, times_ms[len(times_ms)-1]))
		default:
			session.ChannelMessageSend(message.ChannelID, "Por favor, digite um comando válido.")
		}
	}
}
