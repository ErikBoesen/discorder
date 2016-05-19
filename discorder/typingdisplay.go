package discorder

import (
	"github.com/jonas747/discorder/common"
	"github.com/jonas747/discorder/ui"
	"github.com/nsf/termbox-go"
)

// Shows whos typing
type TypingDisplay struct {
	*ui.BaseEntity
	Transform *ui.Transform

	App *App

	text *ui.Text
}

func NewTypingDisplay(app *App) *TypingDisplay {
	td := &TypingDisplay{
		BaseEntity: &ui.BaseEntity{},
		Transform:  &ui.Transform{},
		App:        app,
	}

	t := ui.NewText()

	t.Transform.Parent = td.Transform
	t.Transform.AnchorMax = common.NewVector2I(1, 1)
	t.FG = termbox.ColorCyan
	td.AddChild(t)
	td.text = t
	return td
}

func (t *TypingDisplay) Update() {
	channels := make([]string, len(t.App.ViewManager.mv.Channels))
	copy(channels, t.App.ViewManager.mv.Channels)
	if t.App.ViewManager.mv.ShowAllPrivate {
		t.App.session.State.RLock()
		for _, pChan := range t.App.session.State.PrivateChannels {
			found := false
			for _, subChan := range t.App.ViewManager.mv.Channels {
				if subChan == pChan.ID {
					found = true
					break
				}
			}
			if !found {
				channels = append(channels, pChan.ID)
			}
		}
		t.App.session.State.RUnlock()
	}

	typing := t.App.typingRoutine.GetTyping(channels)

	if len(typing) > 0 {
		t.text.Disabled = false

		typingStr := "Typing: asd"

		for _, v := range typing {
			channel, err := t.App.session.State.Channel(v.ChannelID)
			if err != nil {
				continue
			}
			if !channel.IsPrivate {
				member, err := t.App.session.State.Member(channel.GuildID, v.UserID)
				if err != nil {
					continue
				}
				typingStr += "#" + GetChannelNameOrRecipient(channel) + ":" + member.User.Username + ", "
			} else {
				typingStr += "#DM:" + channel.Recipient.Username + ", "
			}
		}
		// Remove trailing ","
		typingStr = typingStr[:len(typingStr)-2]
		t.text.Text = typingStr
	} else {
		t.text.Disabled = true
	}
}

func (t *TypingDisplay) Destroy() { t.DestroyChildren() }

func (t *TypingDisplay) GetRequiredSize() common.Vector2F {
	rect := t.text.Transform.GetRect()
	//log.Println(float32(t.text.HeightRequired()), t.text.Text)
	return common.Vector2F{rect.W, float32(t.text.HeightRequired())}
}

func (t *TypingDisplay) GetTransform() *ui.Transform {
	return t.Transform
}

func (t *TypingDisplay) IsLayoutDynamic() bool {
	return false
}
