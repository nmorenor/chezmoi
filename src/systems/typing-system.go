package systems

import (
	"image/color"
	"sync"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/f1bonacc1/glippy"
	"github.com/nmorenor/chezmoi/entities"
	"github.com/nmorenor/chezmoi/net"
	"github.com/nmorenor/chezmoi/options"
	"github.com/sacOO7/gowebsocket"

	"github.com/nmorenor/chezmoi-net/client"
)

type TypingSystem struct {
	Session           bool
	label             entities.Label
	runeLock          sync.Mutex
	runesToAdd        []rune
	timeSinceDeletion float32
}

func (t *TypingSystem) New(w *ecs.World) {
	fnt := &common.Font{
		URL:  "fonts/pixelfont.ttf",
		FG:   color.White,
		Size: 12,
	}
	err := fnt.CreatePreloaded()
	if err != nil {
		panic(err)
	}

	t.label = entities.Label{BasicEntity: ecs.NewBasic()}
	if t.Session {
		t.label.SpaceComponent.Position = engo.Point{X: (engo.WindowWidth() / 2) - 150, Y: (engo.WindowHeight() / 2) + 50}
	} else {
		t.label.SpaceComponent.Position = engo.Point{X: (engo.WindowWidth() / 2) - 70, Y: (engo.WindowHeight() / 2) + 50}
	}
	t.label.RenderComponent.Drawable = common.Text{
		Font: fnt,
		Text: "",
	}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&t.label.BasicEntity, &t.label.RenderComponent, &t.label.SpaceComponent)
		}
	}

	engo.Mailbox.Listen("TextMessage", func(msg engo.Message) {
		m, ok := msg.(engo.TextMessage)
		if !ok {
			return
		}
		t.runeLock.Lock()
		t.runesToAdd = append(t.runesToAdd, m.Char)
		t.runeLock.Unlock()
	})
}

func (*TypingSystem) Remove(ecs.BasicEntity) {}

func (t *TypingSystem) StartSession() {
	remoteClient := net.NewRemoteClient(client.NewClient(gowebsocket.New("ws://localhost:8080/ws")), *options.SessionInfo.Username, options.SessionInfo.HostMode)
	if t.Session {
		remoteClient.Session = options.SessionInfo.Session
		remoteClient.Client.Session = options.SessionInfo.Session
	}
	options.SessionInfo.Client = remoteClient
	go func(c *client.Client) {
		c.Connect()
	}(remoteClient.Client)
}

func (t *TypingSystem) Update(dt float32) {
	if options.SessionInfo.Client != nil && options.SessionInfo.Client.Client.Id != nil {
		engo.SetSceneByName("Session", true)
		return
	}
	t.timeSinceDeletion += dt
	str := ""
	txt := t.label.Drawable.(common.Text)
	if engo.Input.Button("backspace").Down() && len(txt.Text) > 0 {
		if t.timeSinceDeletion > 0.2 {
			t.timeSinceDeletion = 0
			txt.Text = txt.Text[:len(txt.Text)-1]
			t.label.Drawable = txt
		}
		return
	}
	if engo.Input.Button("paste").Down() {
		candidate, err := glippy.Get()
		if err == nil {
			txt.Text = candidate
			t.label.Drawable = txt
		}
		return
	}
	t.runeLock.Lock()
	if len(t.runesToAdd) != 0 {
		str = string(t.runesToAdd)
	}
	t.runesToAdd = make([]rune, 0)
	t.runeLock.Unlock()
	txt.Text += str
	if engo.Input.Button("enter").JustPressed() {
		if options.SessionInfo.HostMode {
			username := txt.Text
			options.SessionInfo.Username = &username
			t.StartSession()
			engo.SetSceneByName("Session", true)
		} else {
			candidate := txt.Text
			if t.Session {
				options.SessionInfo.Session = &candidate
				t.StartSession()
			} else {
				options.SessionInfo.Username = &candidate
				engo.SetSceneByName("JoiSession", true)
			}
		}

	}
	t.label.Drawable = txt
}
