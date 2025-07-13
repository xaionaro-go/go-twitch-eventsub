package twitch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/coder/websocket"
)

const (
	twitchWebsocketUrl = "wss://eventsub.wss.twitch.tv/ws"
)

var (
	ErrConnClosed   = fmt.Errorf("connection closed")
	ErrNilOnWelcome = fmt.Errorf("OnWelcome function was not set")

	messageTypeMap = map[string]func() any{
		"session_welcome":   zeroPtrGen[WelcomeMessage](),
		"session_keepalive": zeroPtrGen[KeepAliveMessage](),
		"notification":      zeroPtrGen[NotificationMessage](),
		"session_reconnect": zeroPtrGen[ReconnectMessage](),
		"revocation":        zeroPtrGen[RevokeMessage](),
	}
)

func zeroPtrGen[T any]() func() any {
	return func() any {
		return new(T)
	}
}

func callFunc[T any](f func(T), v T) {
	if f != nil {
		go f(v)
	}
}

func callFuncWithMsg[T any](f func(T, NotificationMessage), v T, msg NotificationMessage) {
	if f != nil {
		go f(v, msg)
	}
}

type Client struct {
	Address   string
	ws        *websocket.Conn
	connected bool
	ctx       context.Context

	reconnecting bool
	reconnected  chan struct{}

	// Responses
	onError        func(err error)
	onWelcome      func(message WelcomeMessage)
	onKeepAlive    func(message KeepAliveMessage)
	onNotification func(message NotificationMessage)
	onReconnect    func(message ReconnectMessage)
	onRevoke       func(message RevokeMessage)

	// Events
	onRawEvent                                              func(event string, metadata MessageMetadata, subscription PayloadSubscription)
	onEventChannelUpdate                                    func(event EventChannelUpdate, msg NotificationMessage)
	onEventChannelFollow                                    func(event EventChannelFollow, msg NotificationMessage)
	onEventChannelSubscribe                                 func(event EventChannelSubscribe, msg NotificationMessage)
	onEventChannelSubscriptionEnd                           func(event EventChannelSubscriptionEnd, msg NotificationMessage)
	onEventChannelSubscriptionGift                          func(event EventChannelSubscriptionGift, msg NotificationMessage)
	onEventChannelSubscriptionMessage                       func(event EventChannelSubscriptionMessage, msg NotificationMessage)
	onEventChannelCheer                                     func(event EventChannelCheer, msg NotificationMessage)
	onEventChannelRaid                                      func(event EventChannelRaid, msg NotificationMessage)
	onEventChannelBan                                       func(event EventChannelBan, msg NotificationMessage)
	onEventChannelUnban                                     func(event EventChannelUnban, msg NotificationMessage)
	onEventChannelModeratorAdd                              func(event EventChannelModeratorAdd, msg NotificationMessage)
	onEventChannelModeratorRemove                           func(event EventChannelModeratorRemove, msg NotificationMessage)
	onEventChannelVIPAdd                                    func(event EventChannelVIPAdd, msg NotificationMessage)
	onEventChannelVIPRemove                                 func(event EventChannelVIPRemove, msg NotificationMessage)
	onEventChannelChannelPointsCustomRewardAdd              func(event EventChannelChannelPointsCustomRewardAdd, msg NotificationMessage)
	onEventChannelChannelPointsCustomRewardUpdate           func(event EventChannelChannelPointsCustomRewardUpdate, msg NotificationMessage)
	onEventChannelChannelPointsCustomRewardRemove           func(event EventChannelChannelPointsCustomRewardRemove, msg NotificationMessage)
	onEventChannelChannelPointsCustomRewardRedemptionAdd    func(event EventChannelChannelPointsCustomRewardRedemptionAdd, msg NotificationMessage)
	onEventChannelChannelPointsCustomRewardRedemptionUpdate func(event EventChannelChannelPointsCustomRewardRedemptionUpdate, msg NotificationMessage)
	onEventChannelChannelPointsAutomaticRewardRedemptionAdd func(event EventChannelChannelPointsAutomaticRewardRedemptionAdd, msg NotificationMessage)
	onEventChannelPollBegin                                 func(event EventChannelPollBegin, msg NotificationMessage)
	onEventChannelPollProgress                              func(event EventChannelPollProgress, msg NotificationMessage)
	onEventChannelPollEnd                                   func(event EventChannelPollEnd, msg NotificationMessage)
	onEventChannelPredictionBegin                           func(event EventChannelPredictionBegin, msg NotificationMessage)
	onEventChannelPredictionProgress                        func(event EventChannelPredictionProgress, msg NotificationMessage)
	onEventChannelPredictionLock                            func(event EventChannelPredictionLock, msg NotificationMessage)
	onEventChannelPredictionEnd                             func(event EventChannelPredictionEnd, msg NotificationMessage)
	onEventDropEntitlementGrant                             func(event []EventDropEntitlementGrant, msg NotificationMessage)
	onEventExtensionBitsTransactionCreate                   func(event EventExtensionBitsTransactionCreate, msg NotificationMessage)
	onEventChannelGoalBegin                                 func(event EventChannelGoalBegin, msg NotificationMessage)
	onEventChannelGoalProgress                              func(event EventChannelGoalProgress, msg NotificationMessage)
	onEventChannelGoalEnd                                   func(event EventChannelGoalEnd, msg NotificationMessage)
	onEventChannelHypeTrainBegin                            func(event EventChannelHypeTrainBegin, msg NotificationMessage)
	onEventChannelHypeTrainProgress                         func(event EventChannelHypeTrainProgress, msg NotificationMessage)
	onEventChannelHypeTrainEnd                              func(event EventChannelHypeTrainEnd, msg NotificationMessage)
	onEventStreamOnline                                     func(event EventStreamOnline, msg NotificationMessage)
	onEventStreamOffline                                    func(event EventStreamOffline, msg NotificationMessage)
	onEventUserAuthorizationGrant                           func(event EventUserAuthorizationGrant, msg NotificationMessage)
	onEventUserAuthorizationRevoke                          func(event EventUserAuthorizationRevoke, msg NotificationMessage)
	onEventUserUpdate                                       func(event EventUserUpdate, msg NotificationMessage)
	onEventChannelCharityCampaignDonate                     func(event EventChannelCharityCampaignDonate, msg NotificationMessage)
	onEventChannelCharityCampaignProgress                   func(event EventChannelCharityCampaignProgress, msg NotificationMessage)
	onEventChannelCharityCampaignStart                      func(event EventChannelCharityCampaignStart, msg NotificationMessage)
	onEventChannelCharityCampaignStop                       func(event EventChannelCharityCampaignStop, msg NotificationMessage)
	onEventChannelShieldModeBegin                           func(event EventChannelShieldModeBegin, msg NotificationMessage)
	onEventChannelShieldModeEnd                             func(event EventChannelShieldModeEnd, msg NotificationMessage)
	onEventChannelShoutoutCreate                            func(event EventChannelShoutoutCreate, msg NotificationMessage)
	onEventChannelShoutoutReceive                           func(event EventChannelShoutoutReceive, msg NotificationMessage)
	onEventChannelModerate                                  func(event EventChannelModerate, msg NotificationMessage)
	onEventAutomodMessageHold                               func(event EventAutomodMessageHold, msg NotificationMessage)
	onEventAutomodMessageUpdate                             func(event EventAutomodMessageUpdate, msg NotificationMessage)
	onEventAutomodSettingsUpdate                            func(event EventAutomodSettingsUpdate, msg NotificationMessage)
	onEventAutomodTermsUpdate                               func(event EventAutomodTermsUpdate, msg NotificationMessage)
	onEventChannelChatUserMessageHold                       func(event EventChannelChatUserMessageHold, msg NotificationMessage)
	onEventChannelChatUserMessageUpdate                     func(event EventChannelChatUserMessageUpdate, msg NotificationMessage)
	onEventChannelChatClear                                 func(event EventChannelChatClear, msg NotificationMessage)
	onEventChannelChatClearUserMessages                     func(event EventChannelChatClearUserMessages, msg NotificationMessage)
	onEventChannelChatMessage                               func(event EventChannelChatMessage, msg NotificationMessage)
	onEventChannelChatMessageDelete                         func(event EventChannelChatMessageDelete, msg NotificationMessage)
	onEventChannelChatNotification                          func(event EventChannelChatNotification, msg NotificationMessage)
	onEventChannelChatSettingsUpdate                        func(event EventChannelChatSettingsUpdate, msg NotificationMessage)
	onEventChannelSuspiciousUserMessage                     func(event EventChannelSuspiciousUserMessage, msg NotificationMessage)
	onEventChannelSuspiciousUserUpdate                      func(event EventChannelSuspiciousUserUpdate, msg NotificationMessage)
	onEventChannelSharedChatBegin                           func(event EventChannelSharedChatBegin, msg NotificationMessage)
	onEventChannelSharedChatUpdate                          func(event EventChannelSharedChatUpdate, msg NotificationMessage)
	onEventChannelSharedChatEnd                             func(event EventChannelSharedChatEnd, msg NotificationMessage)
	onEventUserWhisperMessage                               func(event EventUserWhisperMessage, msg NotificationMessage)
	onEventChannelAdBreakBegin                              func(event EventChannelAdBreakBegin, msg NotificationMessage)
	onEventChannelWarningAcknowledge                        func(event EventChannelWarningAcknowledge, msg NotificationMessage)
	onEventChannelWarningSend                               func(event EventChannelWarningSend, msg NotificationMessage)
	onEventChannelUnbanRequestCreate                        func(event EventChannelUnbanRequestCreate, msg NotificationMessage)
	onEventChannelUnbanRequestResolve                       func(event EventChannelUnbanRequestResolve, msg NotificationMessage)
	onEventConduitShardDisabled                             func(event EventConduitShardDisabled, msg NotificationMessage)
}

func NewClient() *Client {
	return NewClientWithUrl(twitchWebsocketUrl)
}

func NewClientWithUrl(url string) *Client {
	return &Client{
		Address:     url,
		reconnected: make(chan struct{}),
		onError:     func(err error) { fmt.Printf("ERROR: %v\n", err) },
	}
}

func (c *Client) Connect(
	onReadError func(context.Context, error),
) error {
	return c.ConnectWithContext(context.Background(), onReadError)
}

func (c *Client) ConnectWithContext(
	ctx context.Context,
	onReadError func(context.Context, error),
) error {
	if c.onWelcome == nil {
		return ErrNilOnWelcome
	}

	c.ctx = ctx
	ws, err := c.dial()
	if err != nil {
		return err
	}
	c.ws = ws
	c.connected = true

	go c.readLoop(ctx, onReadError)
	return nil
}

func (c *Client) readLoop(
	ctx context.Context,
	onReadError func(context.Context, error),
) {
	for {
		_, data, err := c.ws.Read(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}

			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				if c.reconnecting {
					c.reconnecting = false
					<-c.reconnected
					continue
				}
				return
			}

			if onReadError == nil {
				c.onError(err)
				return
			}
			onReadError(ctx, err)
			return
		}

		err = c.handleMessage(data)
		if err != nil {
			c.onError(err)
		}
	}
}

func (c *Client) Close() error {
	defer func() { c.ws = nil }()
	if !c.connected {
		return nil
	}
	c.connected = false

	err := c.ws.Close(websocket.StatusNormalClosure, "Stopping Connection")

	var closeError websocket.CloseError
	if err != nil && !errors.As(err, &closeError) {
		return fmt.Errorf("could not close websocket connection: %w", err)
	}
	return nil
}

func (c *Client) handleMessage(data []byte) error {
	metadata, err := parseBaseMessage(data)
	if err != nil {
		return err
	}

	messageType := metadata.MessageType
	genMessage, ok := messageTypeMap[messageType]
	if !ok {
		return fmt.Errorf("unknown message type %s: %s", messageType, string(data))
	}

	message := genMessage()
	err = json.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("could not unmarshal message into %s: %w", messageType, err)
	}

	switch msg := message.(type) {
	case *WelcomeMessage:
		callFunc(c.onWelcome, *msg)
	case *KeepAliveMessage:
		callFunc(c.onKeepAlive, *msg)
	case *NotificationMessage:
		callFunc(c.onNotification, *msg)

		err = c.handleNotification(*msg)
		if err != nil {
			return fmt.Errorf("could not handle notification: %w", err)
		}
	case *ReconnectMessage:
		callFunc(c.onReconnect, *msg)

		err = c.reconnect(*msg)
		if err != nil {
			return fmt.Errorf("could not handle reconnect: %w", err)
		}
	case *RevokeMessage:
		callFunc(c.onRevoke, *msg)
	default:
		return fmt.Errorf("unhandled %T message: %v", msg, msg)
	}

	return nil
}

func (c *Client) reconnect(message ReconnectMessage) error {
	c.Address = message.Payload.Session.ReconnectUrl
	ws, err := c.dial()
	if err != nil {
		return fmt.Errorf("could not dial to reconnect")
	}

	go func() {
		_, data, err := ws.Read(c.ctx)
		if err != nil {
			c.onError(fmt.Errorf("reconnect failed: could not read reconnect websocket for welcome: %w", err))
		}

		metadata, err := parseBaseMessage(data)
		if err != nil {
			c.onError(fmt.Errorf("reconnect failed: could parse base message: %w", err))
		}

		if metadata.MessageType != "session_welcome" {
			c.onError(fmt.Errorf("reconnect failed: did not get a session_welcome message first: got message %s", metadata.MessageType))
			return
		}

		c.reconnecting = true
		c.ws.Close(websocket.StatusNormalClosure, "Stopping Connection")
		c.ws = ws
		c.reconnected <- struct{}{}
	}()

	return nil
}

func (c *Client) handleNotification(message NotificationMessage) error {
	data, err := message.Payload.Event.MarshalJSON()
	if err != nil {
		return fmt.Errorf("could not get event json: %w", err)
	}

	subscription := message.Payload.Subscription
	metadata, ok := subMetadata[subscription.Type]
	if !ok {
		return fmt.Errorf("unknown subscription type %s", subscription.Type)
	}

	if c.onRawEvent != nil {
		c.onRawEvent(string(data), message.Metadata, subscription)
	}

	var newEvent any
	if metadata.EventGen != nil {
		newEvent = metadata.EventGen()
		err = json.Unmarshal(data, newEvent)
		if err != nil {
			return fmt.Errorf("could not unmarshal %s into %T: %w", subscription.Type, newEvent, err)
		}
	}

	switch event := newEvent.(type) {
	case *EventChannelUpdate:
		callFuncWithMsg(c.onEventChannelUpdate, *event, message)
	case *EventChannelFollow:
		callFuncWithMsg(c.onEventChannelFollow, *event, message)
	case *EventChannelSubscribe:
		callFuncWithMsg(c.onEventChannelSubscribe, *event, message)
	case *EventChannelSubscriptionEnd:
		callFuncWithMsg(c.onEventChannelSubscriptionEnd, *event, message)
	case *EventChannelSubscriptionGift:
		callFuncWithMsg(c.onEventChannelSubscriptionGift, *event, message)
	case *EventChannelSubscriptionMessage:
		callFuncWithMsg(c.onEventChannelSubscriptionMessage, *event, message)
	case *EventChannelCheer:
		callFuncWithMsg(c.onEventChannelCheer, *event, message)
	case *EventChannelRaid:
		callFuncWithMsg(c.onEventChannelRaid, *event, message)
	case *EventChannelBan:
		callFuncWithMsg(c.onEventChannelBan, *event, message)
	case *EventChannelUnban:
		callFuncWithMsg(c.onEventChannelUnban, *event, message)
	case *EventChannelModeratorAdd:
		callFuncWithMsg(c.onEventChannelModeratorAdd, *event, message)
	case *EventChannelModeratorRemove:
		callFuncWithMsg(c.onEventChannelModeratorRemove, *event, message)
	case *EventChannelVIPAdd:
		callFuncWithMsg(c.onEventChannelVIPAdd, *event, message)
	case *EventChannelVIPRemove:
		callFuncWithMsg(c.onEventChannelVIPRemove, *event, message)
	case *EventChannelChannelPointsCustomRewardAdd:
		callFuncWithMsg(c.onEventChannelChannelPointsCustomRewardAdd, *event, message)
	case *EventChannelChannelPointsCustomRewardUpdate:
		callFuncWithMsg(c.onEventChannelChannelPointsCustomRewardUpdate, *event, message)
	case *EventChannelChannelPointsCustomRewardRemove:
		callFuncWithMsg(c.onEventChannelChannelPointsCustomRewardRemove, *event, message)
	case *EventChannelChannelPointsCustomRewardRedemptionAdd:
		callFuncWithMsg(c.onEventChannelChannelPointsCustomRewardRedemptionAdd, *event, message)
	case *EventChannelChannelPointsCustomRewardRedemptionUpdate:
		callFuncWithMsg(c.onEventChannelChannelPointsCustomRewardRedemptionUpdate, *event, message)
	case *EventChannelChannelPointsAutomaticRewardRedemptionAdd:
		callFuncWithMsg(c.onEventChannelChannelPointsAutomaticRewardRedemptionAdd, *event, message)
	case *EventChannelPollBegin:
		callFuncWithMsg(c.onEventChannelPollBegin, *event, message)
	case *EventChannelPollProgress:
		callFuncWithMsg(c.onEventChannelPollProgress, *event, message)
	case *EventChannelPollEnd:
		callFuncWithMsg(c.onEventChannelPollEnd, *event, message)
	case *EventChannelPredictionBegin:
		callFuncWithMsg(c.onEventChannelPredictionBegin, *event, message)
	case *EventChannelPredictionProgress:
		callFuncWithMsg(c.onEventChannelPredictionProgress, *event, message)
	case *EventChannelPredictionLock:
		callFuncWithMsg(c.onEventChannelPredictionLock, *event, message)
	case *EventChannelPredictionEnd:
		callFuncWithMsg(c.onEventChannelPredictionEnd, *event, message)
	case *[]EventDropEntitlementGrant:
		callFuncWithMsg(c.onEventDropEntitlementGrant, *event, message)
	case *EventExtensionBitsTransactionCreate:
		callFuncWithMsg(c.onEventExtensionBitsTransactionCreate, *event, message)
	case *EventChannelGoalBegin:
		callFuncWithMsg(c.onEventChannelGoalBegin, *event, message)
	case *EventChannelGoalProgress:
		callFuncWithMsg(c.onEventChannelGoalProgress, *event, message)
	case *EventChannelGoalEnd:
		callFuncWithMsg(c.onEventChannelGoalEnd, *event, message)
	case *EventChannelHypeTrainBegin:
		callFuncWithMsg(c.onEventChannelHypeTrainBegin, *event, message)
	case *EventChannelHypeTrainProgress:
		callFuncWithMsg(c.onEventChannelHypeTrainProgress, *event, message)
	case *EventChannelHypeTrainEnd:
		callFuncWithMsg(c.onEventChannelHypeTrainEnd, *event, message)
	case *EventStreamOnline:
		callFuncWithMsg(c.onEventStreamOnline, *event, message)
	case *EventStreamOffline:
		callFuncWithMsg(c.onEventStreamOffline, *event, message)
	case *EventUserAuthorizationGrant:
		callFuncWithMsg(c.onEventUserAuthorizationGrant, *event, message)
	case *EventUserAuthorizationRevoke:
		callFuncWithMsg(c.onEventUserAuthorizationRevoke, *event, message)
	case *EventUserUpdate:
		callFuncWithMsg(c.onEventUserUpdate, *event, message)
	case *EventChannelCharityCampaignDonate:
		callFuncWithMsg(c.onEventChannelCharityCampaignDonate, *event, message)
	case *EventChannelCharityCampaignProgress:
		callFuncWithMsg(c.onEventChannelCharityCampaignProgress, *event, message)
	case *EventChannelCharityCampaignStart:
		callFuncWithMsg(c.onEventChannelCharityCampaignStart, *event, message)
	case *EventChannelCharityCampaignStop:
		callFuncWithMsg(c.onEventChannelCharityCampaignStop, *event, message)
	case *EventChannelShieldModeBegin:
		callFuncWithMsg(c.onEventChannelShieldModeBegin, *event, message)
	case *EventChannelShieldModeEnd:
		callFuncWithMsg(c.onEventChannelShieldModeEnd, *event, message)
	case *EventChannelShoutoutCreate:
		callFuncWithMsg(c.onEventChannelShoutoutCreate, *event, message)
	case *EventChannelShoutoutReceive:
		callFuncWithMsg(c.onEventChannelShoutoutReceive, *event, message)
	case *EventChannelModerate:
		callFuncWithMsg(c.onEventChannelModerate, *event, message)
	case *EventAutomodMessageHold:
		callFuncWithMsg(c.onEventAutomodMessageHold, *event, message)
	case *EventAutomodMessageUpdate:
		callFuncWithMsg(c.onEventAutomodMessageUpdate, *event, message)
	case *EventAutomodSettingsUpdate:
		callFuncWithMsg(c.onEventAutomodSettingsUpdate, *event, message)
	case *EventAutomodTermsUpdate:
		callFuncWithMsg(c.onEventAutomodTermsUpdate, *event, message)
	case *EventChannelChatUserMessageHold:
		callFuncWithMsg(c.onEventChannelChatUserMessageHold, *event, message)
	case *EventChannelChatUserMessageUpdate:
		callFuncWithMsg(c.onEventChannelChatUserMessageUpdate, *event, message)
	case *EventChannelChatClear:
		callFuncWithMsg(c.onEventChannelChatClear, *event, message)
	case *EventChannelChatClearUserMessages:
		callFuncWithMsg(c.onEventChannelChatClearUserMessages, *event, message)
	case *EventChannelChatMessage:
		callFuncWithMsg(c.onEventChannelChatMessage, *event, message)
	case *EventChannelChatMessageDelete:
		callFuncWithMsg(c.onEventChannelChatMessageDelete, *event, message)
	case *EventChannelChatNotification:
		callFuncWithMsg(c.onEventChannelChatNotification, *event, message)
	case *EventChannelChatSettingsUpdate:
		callFuncWithMsg(c.onEventChannelChatSettingsUpdate, *event, message)
	case *EventChannelSuspiciousUserMessage:
		callFuncWithMsg(c.onEventChannelSuspiciousUserMessage, *event, message)
	case *EventChannelSuspiciousUserUpdate:
		callFuncWithMsg(c.onEventChannelSuspiciousUserUpdate, *event, message)
	case *EventChannelSharedChatBegin:
		callFuncWithMsg(c.onEventChannelSharedChatBegin, *event, message)
	case *EventChannelSharedChatUpdate:
		callFuncWithMsg(c.onEventChannelSharedChatUpdate, *event, message)
	case *EventChannelSharedChatEnd:
		callFuncWithMsg(c.onEventChannelSharedChatEnd, *event, message)
	case *EventUserWhisperMessage:
		callFuncWithMsg(c.onEventUserWhisperMessage, *event, message)
	case *EventChannelAdBreakBegin:
		callFuncWithMsg(c.onEventChannelAdBreakBegin, *event, message)
	case *EventChannelWarningAcknowledge:
		callFuncWithMsg(c.onEventChannelWarningAcknowledge, *event, message)
	case *EventChannelWarningSend:
		callFuncWithMsg(c.onEventChannelWarningSend, *event, message)
	case *EventChannelUnbanRequestCreate:
		callFuncWithMsg(c.onEventChannelUnbanRequestCreate, *event, message)
	case *EventChannelUnbanRequestResolve:
		callFuncWithMsg(c.onEventChannelUnbanRequestResolve, *event, message)
	case *EventConduitShardDisabled:
		callFuncWithMsg(c.onEventConduitShardDisabled, *event, message)
	default:
		c.onError(fmt.Errorf("unknown event type %s", subscription.Type))
	}

	return nil
}

func (c *Client) dial() (*websocket.Conn, error) {
	ws, _, err := websocket.Dial(c.ctx, c.Address, nil)
	if err != nil {
		return nil, fmt.Errorf("could not dial %s: %w", c.Address, err)
	}
	return ws, nil
}

func parseBaseMessage(data []byte) (MessageMetadata, error) {
	type BaseMessage struct {
		Metadata MessageMetadata `json:"metadata"`
	}

	var baseMessage BaseMessage
	err := json.Unmarshal(data, &baseMessage)
	if err != nil {
		return MessageMetadata{}, fmt.Errorf("could not unmarshal basemessage to get message type: %w", err)
	}

	return baseMessage.Metadata, nil
}

func (c *Client) OnError(callback func(err error)) {
	c.onError = callback
}

func (c *Client) OnWelcome(callback func(message WelcomeMessage)) {
	c.onWelcome = callback
}

func (c *Client) OnKeepAlive(callback func(message KeepAliveMessage)) {
	c.onKeepAlive = callback
}

func (c *Client) OnNotification(callback func(message NotificationMessage)) {
	c.onNotification = callback
}

func (c *Client) OnReconnect(callback func(message ReconnectMessage)) {
	c.onReconnect = callback
}

func (c *Client) OnRevoke(callback func(message RevokeMessage)) {
	c.onRevoke = callback
}

func (c *Client) OnRawEvent(callback func(event string, metadata MessageMetadata, subscription PayloadSubscription)) {
	c.onRawEvent = callback
}

func (c *Client) OnEventChannelUpdate(callback func(event EventChannelUpdate, msg NotificationMessage)) {
	c.onEventChannelUpdate = callback
}

func (c *Client) OnEventChannelFollow(callback func(event EventChannelFollow, msg NotificationMessage)) {
	c.onEventChannelFollow = callback
}

func (c *Client) OnEventChannelSubscribe(callback func(event EventChannelSubscribe, msg NotificationMessage)) {
	c.onEventChannelSubscribe = callback
}

func (c *Client) OnEventChannelSubscriptionEnd(callback func(event EventChannelSubscriptionEnd, msg NotificationMessage)) {
	c.onEventChannelSubscriptionEnd = callback
}

func (c *Client) OnEventChannelSubscriptionGift(callback func(event EventChannelSubscriptionGift, msg NotificationMessage)) {
	c.onEventChannelSubscriptionGift = callback
}

func (c *Client) OnEventChannelSubscriptionMessage(callback func(event EventChannelSubscriptionMessage, msg NotificationMessage)) {
	c.onEventChannelSubscriptionMessage = callback
}

func (c *Client) OnEventChannelCheer(callback func(event EventChannelCheer, msg NotificationMessage)) {
	c.onEventChannelCheer = callback
}

func (c *Client) OnEventChannelRaid(callback func(event EventChannelRaid, msg NotificationMessage)) {
	c.onEventChannelRaid = callback
}

func (c *Client) OnEventChannelBan(callback func(event EventChannelBan, msg NotificationMessage)) {
	c.onEventChannelBan = callback
}

func (c *Client) OnEventChannelUnban(callback func(event EventChannelUnban, msg NotificationMessage)) {
	c.onEventChannelUnban = callback
}

func (c *Client) OnEventChannelModeratorAdd(callback func(event EventChannelModeratorAdd, msg NotificationMessage)) {
	c.onEventChannelModeratorAdd = callback
}

func (c *Client) OnEventChannelModeratorRemove(callback func(event EventChannelModeratorRemove, msg NotificationMessage)) {
	c.onEventChannelModeratorRemove = callback
}

func (c *Client) OnEventChannelVIPAdd(callback func(event EventChannelVIPAdd, msg NotificationMessage)) {
	c.onEventChannelVIPAdd = callback
}

func (c *Client) OnEventChannelVIPRemove(callback func(event EventChannelVIPRemove, msg NotificationMessage)) {
	c.onEventChannelVIPRemove = callback
}

func (c *Client) OnEventChannelChannelPointsCustomRewardAdd(callback func(event EventChannelChannelPointsCustomRewardAdd, msg NotificationMessage)) {
	c.onEventChannelChannelPointsCustomRewardAdd = callback
}

func (c *Client) OnEventChannelChannelPointsCustomRewardUpdate(callback func(event EventChannelChannelPointsCustomRewardUpdate, msg NotificationMessage)) {
	c.onEventChannelChannelPointsCustomRewardUpdate = callback
}

func (c *Client) OnEventChannelChannelPointsCustomRewardRemove(callback func(event EventChannelChannelPointsCustomRewardRemove, msg NotificationMessage)) {
	c.onEventChannelChannelPointsCustomRewardRemove = callback
}

func (c *Client) OnEventChannelChannelPointsCustomRewardRedemptionAdd(callback func(event EventChannelChannelPointsCustomRewardRedemptionAdd, msg NotificationMessage)) {
	c.onEventChannelChannelPointsCustomRewardRedemptionAdd = callback
}

func (c *Client) OnEventChannelChannelPointsCustomRewardRedemptionUpdate(callback func(event EventChannelChannelPointsCustomRewardRedemptionUpdate, msg NotificationMessage)) {
	c.onEventChannelChannelPointsCustomRewardRedemptionUpdate = callback
}

func (c *Client) OnEventChannelChannelPointsAutomaticRewardRedemptionAdd(callback func(event EventChannelChannelPointsAutomaticRewardRedemptionAdd, msg NotificationMessage)) {
	c.onEventChannelChannelPointsAutomaticRewardRedemptionAdd = callback
}

func (c *Client) OnEventChannelPollBegin(callback func(event EventChannelPollBegin, msg NotificationMessage)) {
	c.onEventChannelPollBegin = callback
}

func (c *Client) OnEventChannelPollProgress(callback func(event EventChannelPollProgress, msg NotificationMessage)) {
	c.onEventChannelPollProgress = callback
}

func (c *Client) OnEventChannelPollEnd(callback func(event EventChannelPollEnd, msg NotificationMessage)) {
	c.onEventChannelPollEnd = callback
}

func (c *Client) OnEventChannelPredictionBegin(callback func(event EventChannelPredictionBegin, msg NotificationMessage)) {
	c.onEventChannelPredictionBegin = callback
}

func (c *Client) OnEventChannelPredictionProgress(callback func(event EventChannelPredictionProgress, msg NotificationMessage)) {
	c.onEventChannelPredictionProgress = callback
}

func (c *Client) OnEventChannelPredictionLock(callback func(event EventChannelPredictionLock, msg NotificationMessage)) {
	c.onEventChannelPredictionLock = callback
}

func (c *Client) OnEventChannelPredictionEnd(callback func(event EventChannelPredictionEnd, msg NotificationMessage)) {
	c.onEventChannelPredictionEnd = callback
}

func (c *Client) OnEventDropEntitlementGrant(callback func(event []EventDropEntitlementGrant, msg NotificationMessage)) {
	c.onEventDropEntitlementGrant = callback
}

func (c *Client) OnEventExtensionBitsTransactionCreate(callback func(event EventExtensionBitsTransactionCreate, msg NotificationMessage)) {
	c.onEventExtensionBitsTransactionCreate = callback
}

func (c *Client) OnEventChannelGoalBegin(callback func(event EventChannelGoalBegin, msg NotificationMessage)) {
	c.onEventChannelGoalBegin = callback
}

func (c *Client) OnEventChannelGoalProgress(callback func(event EventChannelGoalProgress, msg NotificationMessage)) {
	c.onEventChannelGoalProgress = callback
}

func (c *Client) OnEventChannelGoalEnd(callback func(event EventChannelGoalEnd, msg NotificationMessage)) {
	c.onEventChannelGoalEnd = callback
}

func (c *Client) OnEventChannelHypeTrainBegin(callback func(event EventChannelHypeTrainBegin, msg NotificationMessage)) {
	c.onEventChannelHypeTrainBegin = callback
}

func (c *Client) OnEventChannelHypeTrainProgress(callback func(event EventChannelHypeTrainProgress, msg NotificationMessage)) {
	c.onEventChannelHypeTrainProgress = callback
}

func (c *Client) OnEventChannelHypeTrainEnd(callback func(event EventChannelHypeTrainEnd, msg NotificationMessage)) {
	c.onEventChannelHypeTrainEnd = callback
}

func (c *Client) OnEventStreamOnline(callback func(event EventStreamOnline, msg NotificationMessage)) {
	c.onEventStreamOnline = callback
}

func (c *Client) OnEventStreamOffline(callback func(event EventStreamOffline, msg NotificationMessage)) {
	c.onEventStreamOffline = callback
}

func (c *Client) OnEventUserAuthorizationGrant(callback func(event EventUserAuthorizationGrant, msg NotificationMessage)) {
	c.onEventUserAuthorizationGrant = callback
}

func (c *Client) OnEventUserAuthorizationRevoke(callback func(event EventUserAuthorizationRevoke, msg NotificationMessage)) {
	c.onEventUserAuthorizationRevoke = callback
}

func (c *Client) OnEventUserUpdate(callback func(event EventUserUpdate, msg NotificationMessage)) {
	c.onEventUserUpdate = callback
}

func (c *Client) OnEventChannelCharityCampaignDonate(callback func(event EventChannelCharityCampaignDonate, msg NotificationMessage)) {
	c.onEventChannelCharityCampaignDonate = callback
}

func (c *Client) OnEventChannelCharityCampaignProgress(callback func(event EventChannelCharityCampaignProgress, msg NotificationMessage)) {
	c.onEventChannelCharityCampaignProgress = callback
}

func (c *Client) OnEventChannelCharityCampaignStart(callback func(event EventChannelCharityCampaignStart, msg NotificationMessage)) {
	c.onEventChannelCharityCampaignStart = callback
}

func (c *Client) OnEventChannelCharityCampaignStop(callback func(event EventChannelCharityCampaignStop, msg NotificationMessage)) {
	c.onEventChannelCharityCampaignStop = callback
}

func (c *Client) OnEventChannelShieldModeBegin(callback func(event EventChannelShieldModeBegin, msg NotificationMessage)) {
	c.onEventChannelShieldModeBegin = callback
}

func (c *Client) OnEventChannelShieldModeEnd(callback func(event EventChannelShieldModeEnd, msg NotificationMessage)) {
	c.onEventChannelShieldModeEnd = callback
}

func (c *Client) OnEventChannelShoutoutCreate(callback func(event EventChannelShoutoutCreate, msg NotificationMessage)) {
	c.onEventChannelShoutoutCreate = callback
}

func (c *Client) OnEventChannelShoutoutReceive(callback func(event EventChannelShoutoutReceive, msg NotificationMessage)) {
	c.onEventChannelShoutoutReceive = callback
}

func (c *Client) OnEventChannelModerate(callback func(event EventChannelModerate, msg NotificationMessage)) {
	c.onEventChannelModerate = callback
}

func (c *Client) OnEventAutomodMessageHold(callback func(event EventAutomodMessageHold, msg NotificationMessage)) {
	c.onEventAutomodMessageHold = callback
}

func (c *Client) OnEventAutomodMessageUpdate(callback func(event EventAutomodMessageUpdate, msg NotificationMessage)) {
	c.onEventAutomodMessageUpdate = callback
}

func (c *Client) OnEventAutomodSettingsUpdate(callback func(event EventAutomodSettingsUpdate, msg NotificationMessage)) {
	c.onEventAutomodSettingsUpdate = callback
}

func (c *Client) OnEventAutomodTermsUpdate(callback func(event EventAutomodTermsUpdate, msg NotificationMessage)) {
	c.onEventAutomodTermsUpdate = callback
}

func (c *Client) OnEventChannelChatUserMessageHold(callback func(event EventChannelChatUserMessageHold, msg NotificationMessage)) {
	c.onEventChannelChatUserMessageHold = callback
}

func (c *Client) OnEventChannelChatUserMessageUpdate(callback func(event EventChannelChatUserMessageUpdate, msg NotificationMessage)) {
	c.onEventChannelChatUserMessageUpdate = callback
}

func (c *Client) OnEventChannelChatClear(callback func(event EventChannelChatClear, msg NotificationMessage)) {
	c.onEventChannelChatClear = callback
}

func (c *Client) OnEventChannelChatClearUserMessages(callback func(event EventChannelChatClearUserMessages, msg NotificationMessage)) {
	c.onEventChannelChatClearUserMessages = callback
}

func (c *Client) OnEventChannelChatMessage(callback func(event EventChannelChatMessage, msg NotificationMessage)) {
	c.onEventChannelChatMessage = callback
}

func (c *Client) OnEventChannelChatMessageDelete(callback func(event EventChannelChatMessageDelete, msg NotificationMessage)) {
	c.onEventChannelChatMessageDelete = callback
}

func (c *Client) OnEventChannelChatNotification(callback func(event EventChannelChatNotification, msg NotificationMessage)) {
	c.onEventChannelChatNotification = callback
}

func (c *Client) OnEventChannelChatSettingsUpdate(callback func(event EventChannelChatSettingsUpdate, msg NotificationMessage)) {
	c.onEventChannelChatSettingsUpdate = callback
}

func (c *Client) OnEventChannelSuspiciousUserMessage(callback func(event EventChannelSuspiciousUserMessage, msg NotificationMessage)) {
	c.onEventChannelSuspiciousUserMessage = callback
}

func (c *Client) OnEventChannelSuspiciousUserUpdate(callback func(event EventChannelSuspiciousUserUpdate, msg NotificationMessage)) {
	c.onEventChannelSuspiciousUserUpdate = callback
}

func (c *Client) OnEventChannelSharedChatBegin(callback func(event EventChannelSharedChatBegin, msg NotificationMessage)) {
	c.onEventChannelSharedChatBegin = callback
}

func (c *Client) OnEventChannelSharedChatUpdate(callback func(event EventChannelSharedChatUpdate, msg NotificationMessage)) {
	c.onEventChannelSharedChatUpdate = callback
}

func (c *Client) OnEventChannelSharedChatEnd(callback func(event EventChannelSharedChatEnd, msg NotificationMessage)) {
	c.onEventChannelSharedChatEnd = callback
}

func (c *Client) OnEventUserWhisperMessage(callback func(event EventUserWhisperMessage, msg NotificationMessage)) {
	c.onEventUserWhisperMessage = callback
}

func (c *Client) OnEventChannelAdBreakBegin(callback func(event EventChannelAdBreakBegin, msg NotificationMessage)) {
	c.onEventChannelAdBreakBegin = callback
}

func (c *Client) OnEventChannelWarningAcknowledge(callback func(event EventChannelWarningAcknowledge, msg NotificationMessage)) {
	c.onEventChannelWarningAcknowledge = callback
}

func (c *Client) OnEventChannelWarningSend(callback func(event EventChannelWarningSend, msg NotificationMessage)) {
	c.onEventChannelWarningSend = callback
}

func (c *Client) OnEventChannelUnbanRequestCreate(callback func(event EventChannelUnbanRequestCreate, msg NotificationMessage)) {
	c.onEventChannelUnbanRequestCreate = callback
}

func (c *Client) OnEventChannelUnbanRequestResolve(callback func(event EventChannelUnbanRequestResolve, msg NotificationMessage)) {
	c.onEventChannelUnbanRequestResolve = callback
}

func (c *Client) OnEventConduitShardDisabled(callback func(event EventConduitShardDisabled, msg NotificationMessage)) {
	c.onEventConduitShardDisabled = callback
}
