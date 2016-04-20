package ui

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nsf/termbox-go"
)

type Entity interface {
	Children(recursive bool) []Entity
	Destroy()
	RunFunc(f func(e Entity)) // Runs f recorsively
}

type BaseEntity struct {
	Self     Entity // To get around some struct embedding limitations, set to whatever is embedding this, This probably needs some redisigning
	entities []Entity
}

// Run f recursively on all children
func (b *BaseEntity) RunFunc(f func(e Entity)) {
	if b.Self != nil {
		f(b.Self)
	}
	if b.entities != nil && len(b.entities) > 0 {
		for _, v := range b.entities {
			v.RunFunc(f)
		}
	}
}

// Maybe reuse the slice...? probably miniscule performance hit to not...
func (b *BaseEntity) Children(recursive bool) []Entity {
	if b.entities == nil || len(b.entities) < 1 {
		return nil
	}

	ret := make([]Entity, len(b.entities))
	copy(ret, b.entities)
	if recursive {
		for _, entity := range b.entities {
			children := entity.Children(true)
			if children != nil {
				ret = append(ret, children...)
			}
		}
	}

	return ret
}

func (b *BaseEntity) AddChild(children ...Entity) {

	if b.entities == nil {
		b.entities = make([]Entity, len(children))
		copy(b.entities, children)
	} else {
		b.entities = append(b.entities, children...)
	}
}

func (b *BaseEntity) RemoveChild(child Entity, destroy bool) {
	if b.entities == nil || len(b.entities) < 1 {
		return
	}

	if destroy {
		child.Destroy()
	}

	index := -1
	for k, v := range b.entities {
		if v == child {
			index = k
			break
		}

	}

	if index != -1 {
		if index == len(b.entities)-1 {
			b.entities = b.entities[:index]
		} else {
			b.entities = append(b.entities[:index], b.entities[index+1:]...)
		}
	}
}

// Only clears the list, does not call Destroy() on them or anythin
func (b *BaseEntity) ClearChildren() {
	b.entities = make([]Entity, 0)
}

func (b *BaseEntity) DestroyChildren() {
	for _, v := range b.entities {
		if v != nil {
			v.Destroy()
		}
	}
}

type InputHandler interface {
	HandleInput(event termbox.Event)
}

type PreDrawHandler interface {
	PreDraw() // Ran before drawing, if you create any entities here then they wont get predraw called on them so only create simple entities
}

type DrawHandler interface {
	GetDrawLayer() int
	Draw()
}

type MessageCreateHandler interface {
	HandleMessageCreate(session *discordgo.MessageCreate, msg *discordgo.Message)
}

type MessageEditHandler interface {
	HandleMessageEdit(session *discordgo.MessageCreate, msg *discordgo.Message)
}

type MessageRemoveHandler interface {
	HandleMessageRemove(session *discordgo.MessageCreate, msg *discordgo.Message)
}