package net

import (
	"fmt"
	"sync"

	"github.com/EngoEngine/engo"
	"github.com/f1bonacc1/glippy"
	"github.com/nmorenor/chezmoi-net/client"
	"github.com/nmorenor/chezmoi-net/utils"
)

const (
	broadcast = "-1"
)

func NewRemoteClient(currentClient *client.Client, userName string, hostMode bool) *RemoteClient {
	remoteClient := &RemoteClient{Client: currentClient, Participants: nil, mutex: &sync.Mutex{}, queueMutex: &sync.Mutex{}, locationMutex: &sync.Mutex{}, Host: hostMode, outUueue: utils.NewQueue[string](), Username: userName}
	remoteClient.Client.OnConnect = remoteClient.onReady
	remoteClient.Client.OnSessionChange = remoteClient.onSessionChange
	return remoteClient
}

type Point struct {
	X float32
	Y float32
}

type Message struct {
	Source    string
	Position  Point
	Point     Point
	Animation string
}

type PositionMessage struct {
	Id string
}

type PositionResponseMessage struct {
	Position Point
	Anim     string
}

type RemoteClient struct {
	Host           bool
	initialized    bool
	Client         *client.Client
	Participants   map[string]*string
	outUueue       *utils.Queue[string]
	mutex          *sync.Mutex
	locationMutex  *sync.Mutex
	queueMutex     *sync.Mutex
	Username       string
	Session        *string
	LocalPosition  *engo.Point
	LocalAnimation *string
	OnRemoteUpdate func(client *RemoteClient, from *string, msg Message)
	OnSessionJoin  func(client *RemoteClient, target *string, position *Point, anim *string)
	OnSessionLeave func(client *RemoteClient, target *string)
	OnSessionEnd   func()
}

func (remoteClient *RemoteClient) target() *string {
	remoteClient.queueMutex.Lock()
	defer remoteClient.queueMutex.Unlock()

	if remoteClient.outUueue.IsEmpty() {
		return nil
	}
	target := remoteClient.outUueue.Remove()
	if *target == broadcast {
		return nil
	}
	return target
}

// This will be called when web socket is connected
func (remoteClient *RemoteClient) onReady() {
	// Register this (RemoteClient) instance to receive rcp calls
	client.RegisterService(remoteClient, remoteClient.Client, remoteClient.target)

	if remoteClient.Host {
		remoteClient.Client.StartHosting(remoteClient.Username)
		glippy.Set(*remoteClient.Client.Session)
		fmt.Println("Session on clipboard: " + *remoteClient.Client.Session)
	} else {
		remoteClient.Client.JoinSession(remoteClient.Username, *remoteClient.Session)
	}

	response := remoteClient.Client.SessionMembers()
	remoteClient.Participants = response.Members

}

func (remoteClient *RemoteClient) Initialize() {
	if remoteClient.initialized {
		return
	}
	if !remoteClient.Host {
		for id := range remoteClient.Participants {
			if id != *remoteClient.Client.Id {
				remoteClient.mutex.Lock()
				defer remoteClient.mutex.Unlock()
				rpcClient := remoteClient.Client.GetRpcClientForService(*remoteClient)
				sname := remoteClient.Client.GetServiceName(*remoteClient)
				var position PositionResponseMessage
				if rpcClient != nil {
					remoteClient.queueMutex.Lock()
					remoteClient.outUueue.Add(&id)
					remoteClient.queueMutex.Unlock()
					msg := PositionMessage{Id: id}
					rpcClient.Call(sname+".GetPosition", msg, &position)
				}
				remoteClient.OnSessionJoin(remoteClient, &id, &position.Position, &position.Anim)
				break
			}
		}
	}
	remoteClient.initialized = true
}

func (remoteClient *RemoteClient) SetLocalPosition(position *engo.Point, anim *string) {
	remoteClient.locationMutex.Lock()
	defer remoteClient.locationMutex.Unlock()
	remoteClient.LocalPosition = position
	if anim != nil {
		remoteClient.LocalAnimation = anim
	}
}

func (remoteClient *RemoteClient) SendMessage(vector Point, position Point, animation string, target *string) {
	remoteClient.mutex.Lock()
	defer remoteClient.mutex.Unlock()
	rpcClient := remoteClient.Client.GetRpcClientForService(*remoteClient)
	sname := remoteClient.Client.GetServiceName(*remoteClient)
	msg := &Message{
		Source:    *remoteClient.Client.Id,
		Point:     vector,
		Position:  position,
		Animation: animation,
	}
	// if message starts with [memberName] try to lookup as target
	if target != nil {
		candidate := remoteClient.findParticipantFromName(*target)
		if candidate != nil {
			remoteClient.queueMutex.Lock()
			remoteClient.outUueue.Add(candidate)
			remoteClient.queueMutex.Unlock()
		} else {
			remoteClient.queueMutex.Lock()
			remoteClient.outUueue.Add(ptr(broadcast))
			remoteClient.queueMutex.Unlock()
		}
	}

	if rpcClient != nil {
		var reply string
		rpcClient.Call(sname+".OnMessage", msg, &reply)
	}
}

func ptr[T any](t T) *T {
	return &t
}

func (remoteClient *RemoteClient) findParticipantFromName(target string) *string {
	for id, name := range remoteClient.Participants {
		if *name == target {
			return &id
		}
	}
	return nil
}

/**
 * Message received from rcp call, RPC methods must follow the signature
 */
func (remoteClient *RemoteClient) OnMessage(message *Message, reply *string) error {
	remoteClient.mutex.Lock()
	defer remoteClient.mutex.Unlock()
	if remoteClient.Participants[message.Source] != nil {
		if remoteClient.OnRemoteUpdate != nil {
			fmt.Println("hola")
			remoteClient.OnRemoteUpdate(remoteClient, &message.Source, *message)
		}
	}
	*reply = "OK"
	return nil
}

func (remoteClient *RemoteClient) GetPosition(message *PositionMessage, reply *PositionResponseMessage) error {
	remoteClient.locationMutex.Lock()
	defer remoteClient.locationMutex.Unlock()
	if remoteClient.LocalPosition != nil && remoteClient.LocalAnimation != nil {
		*reply = PositionResponseMessage{Position: Point{X: remoteClient.LocalPosition.X, Y: remoteClient.LocalPosition.Y}, Anim: *remoteClient.LocalAnimation}
	}
	return nil
}

func (remoteClient *RemoteClient) onSessionChange(event client.SessionChangeEvent) {
	remoteClient.mutex.Lock()
	defer remoteClient.mutex.Unlock()
	response := remoteClient.Client.SessionMembers()
	oldParticipants := remoteClient.Participants
	remoteClient.Participants = response.Members
	if event.EventType == client.SESSION_JOIN && remoteClient.Participants[event.EventSource] != nil {
		if remoteClient.OnSessionJoin != nil {
			remoteClient.OnSessionJoin(remoteClient, &event.EventSource, nil, nil)
		}
	}
	if event.EventType == client.SESSION_LEAVE && oldParticipants[event.EventSource] != nil {
		if remoteClient.OnSessionLeave != nil {
			remoteClient.OnSessionLeave(remoteClient, &event.EventSource)
		}
	}
	if event.EventType == client.SESSION_END {
		if remoteClient.OnSessionEnd != nil {
			remoteClient.OnSessionEnd()
		}
	}
}
