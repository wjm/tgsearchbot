package worker

import (
    "fmt"
    "strings"
    "strconv"
    "time"

	"github.com/golang/glog"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/virushuo/tgsearchbot/internal/pkg/service"
	"github.com/virushuo/tgsearchbot/pkg/cypress"
)

// TGBotCommand telegram bot command processor
func TGBotCommand(tgservice *service.Telegram, cypressapi *cypress.API, message *tgbotapi.Message) {
	if strings.HasPrefix(message.Text , "/s ") == true {
        result, err := cypressapi.Search(message.Text[3:], message.Chat.ID)
        if err != nil {
	        glog.Errorf("cypressapi Search error: %v\n", err)
        }
        outputresult := FormatSearchResult(result)
	    replymsg := tgbotapi.NewMessage(message.Chat.ID, outputresult)
	    replymsg.ReplyToMessageID = message.MessageID
        replymsg.ParseMode = "HTML"
	    tgservice.Bot.Send(replymsg)
    }
}

// FormatSearchResult format the search result and build the message text for telegram reply
func FormatSearchResult(result *cypress.Result) string {
    builder := strings.Builder{}
    for idx , item := range result.Items {
        if idx >= 5 {
            break
        }
        humanTimestr := ""
        timestamp , err := strconv.ParseInt(item.Date, 10, 64)
        if err == nil {
            t := time.Unix(timestamp, 0)
            humanTimestr = t.Format("2006-01-02 15:04:05")
        }

        title := strings.Replace(item.Title, "<span class='yx_hl'>", "<b>", -1)
        title = strings.Replace(title, "</span>", "</b>", -1)

        tagline := ""
        if len(item.UserName) > 0 {
            tagline += "["+ item.UserName +"]"
        }
        if len(item.URI) > 0 {
            tagline += fmt.Sprintf(" - <a href=\"%s\">%s</a>", item.URI, humanTimestr)
        } else {
            tagline += " - " + humanTimestr
        }

        itemstr := fmt.Sprintf("%d. %s %s \n", idx, title, tagline)
        builder.WriteString(itemstr)
    }

    return builder.String()
}
