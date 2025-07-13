package twitch_test

import (
	"testing"

	"github.com/joeyak/go-twitch-eventsub/v3"
)

func assertSpecificEventOccurred(t *testing.T, register func(client *twitch.Client, ch chan struct{}), event twitch.EventSubscription, suffixes ...string) {
	assertEventOccurred(t, func(ch chan struct{}) {
		client := newClientWithWelcome(t, "", event, getTestEventData(event, suffixes...))
		register(client, ch)
		go connect(t, client)
	})
}

func TestNotification(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnNotification(func(message twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubStreamOnline)
}

func TestUnknownSubscription(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnError(func(err error) {
			close(ch)
		})
	}, "unknown")
}

func TestEventChannelUpdate(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelUpdate(func(event twitch.EventChannelUpdate, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelUpdate)
}

func TestEventChannelFollow(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelFollow(func(event twitch.EventChannelFollow, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelFollow)
}

func TestEventChannelSubscribe(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelSubscribe(func(event twitch.EventChannelSubscribe, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelSubscribe)
}

func TestEventChannelSubscriptionEnd(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelSubscriptionEnd(func(event twitch.EventChannelSubscriptionEnd, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelSubscriptionEnd)
}

func TestEventChannelSubscriptionGift(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelSubscriptionGift(func(event twitch.EventChannelSubscriptionGift, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelSubscriptionGift)
}

func TestEventChannelSubscriptionGiftAnon(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelSubscriptionGift(func(event twitch.EventChannelSubscriptionGift, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelSubscriptionGift, "anon")
}

func TestEventChannelSubscriptionMessage(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelSubscriptionMessage(func(event twitch.EventChannelSubscriptionMessage, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelSubscriptionMessage)
}

func TestEventChannelSubscriptionMessageNoStreak(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelSubscriptionMessage(func(event twitch.EventChannelSubscriptionMessage, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelSubscriptionMessage, "nostreak")
}

func TestEventChannelCheer(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelCheer(func(event twitch.EventChannelCheer, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelCheer)
}

func TestEventChannelCheerAnon(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelCheer(func(event twitch.EventChannelCheer, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelCheer, "anon")
}

func TestEventChannelRaid(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelRaid(func(event twitch.EventChannelRaid, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelRaid)
}

func TestEventChannelBan(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelBan(func(event twitch.EventChannelBan, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelBan)
}

func TestEventChannelUnban(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelUnban(func(event twitch.EventChannelUnban, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelUnban)
}

func TestEventChannelModeratorAdd(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelModeratorAdd(func(event twitch.EventChannelModeratorAdd, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelModeratorAdd)
}

func TestEventChannelModeratorRemove(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelModeratorRemove(func(event twitch.EventChannelModeratorRemove, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelModeratorRemove)
}

func TestEventChannelVIPAdd(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelVIPAdd(func(event twitch.EventChannelVIPAdd, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelVIPAdd)
}

func TestEventChannelVIPRemove(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelVIPRemove(func(event twitch.EventChannelVIPRemove, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelVIPRemove)
}

func TestEventChannelChannelPointsCustomRewardAdd(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelChannelPointsCustomRewardAdd(func(event twitch.EventChannelChannelPointsCustomRewardAdd, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelChannelPointsCustomRewardAdd)
}

func TestEventChannelChannelPointsCustomRewardUpdate(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelChannelPointsCustomRewardUpdate(func(event twitch.EventChannelChannelPointsCustomRewardUpdate, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelChannelPointsCustomRewardUpdate)
}

func TestEventChannelChannelPointsCustomRewardRemove(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelChannelPointsCustomRewardRemove(func(event twitch.EventChannelChannelPointsCustomRewardRemove, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelChannelPointsCustomRewardRemove)
}

func TestEventChannelChannelPointsCustomRewardRedemptionAdd(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelChannelPointsCustomRewardRedemptionAdd(func(event twitch.EventChannelChannelPointsCustomRewardRedemptionAdd, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelChannelPointsCustomRewardRedemptionAdd)
}

func TestEventChannelChannelPointsCustomRewardRedemptionUpdate(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelChannelPointsCustomRewardRedemptionUpdate(func(event twitch.EventChannelChannelPointsCustomRewardRedemptionUpdate, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelChannelPointsCustomRewardRedemptionUpdate)
}

func TestEventChannelChannelPointsAutomaticRewardRedemptionAdd(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelChannelPointsAutomaticRewardRedemptionAdd(func(event twitch.EventChannelChannelPointsAutomaticRewardRedemptionAdd, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelChannelPointsAutomaticRewardRedemptionAdd)
}

func TestEventChannelPollBegin(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelPollBegin(func(event twitch.EventChannelPollBegin, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelPollBegin)
}

func TestEventChannelPollProgress(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelPollProgress(func(event twitch.EventChannelPollProgress, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelPollProgress)
}

func TestEventChannelPollEnd(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelPollEnd(func(event twitch.EventChannelPollEnd, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelPollEnd)
}

func TestEventChannelPredictionBegin(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelPredictionBegin(func(event twitch.EventChannelPredictionBegin, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelPredictionBegin)
}

func TestEventChannelPredictionProgress(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelPredictionProgress(func(event twitch.EventChannelPredictionProgress, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelPredictionProgress)
}

func TestEventChannelPredictionLock(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelPredictionLock(func(event twitch.EventChannelPredictionLock, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelPredictionLock)
}

func TestEventChannelPredictionEnd(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelPredictionEnd(func(event twitch.EventChannelPredictionEnd, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelPredictionEnd)
}

func TestEventDropEntitlementGrant(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventDropEntitlementGrant(func(event []twitch.EventDropEntitlementGrant, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubDropEntitlementGrant)
}

func TestEventExtensionBitsTransactionCreate(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventExtensionBitsTransactionCreate(func(event twitch.EventExtensionBitsTransactionCreate, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubExtensionBitsTransactionCreate)
}

func TestEventChannelGoalBegin(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelGoalBegin(func(event twitch.EventChannelGoalBegin, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelGoalBegin)
}

func TestEventChannelGoalProgress(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelGoalProgress(func(event twitch.EventChannelGoalProgress, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelGoalProgress)
}

func TestEventChannelGoalEnd(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelGoalEnd(func(event twitch.EventChannelGoalEnd, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelGoalEnd)
}

func TestEventChannelHypeTrainBegin(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelHypeTrainBegin(func(event twitch.EventChannelHypeTrainBegin, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelHypeTrainBegin)
}

func TestEventChannelHypeTrainProgress(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelHypeTrainProgress(func(event twitch.EventChannelHypeTrainProgress, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelHypeTrainProgress)
}

func TestEventChannelHypeTrainEnd(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelHypeTrainEnd(func(event twitch.EventChannelHypeTrainEnd, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelHypeTrainEnd)
}

func TestEventStreamOnline(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventStreamOnline(func(event twitch.EventStreamOnline, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubStreamOnline)
}

func TestEventStreamOffline(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventStreamOffline(func(event twitch.EventStreamOffline, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubStreamOffline)
}

func TestEventUserAuthorizationGrant(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventUserAuthorizationGrant(func(event twitch.EventUserAuthorizationGrant, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubUserAuthorizationGrant)
}

func TestEventUserAuthorizationRevoke(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventUserAuthorizationRevoke(func(event twitch.EventUserAuthorizationRevoke, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubUserAuthorizationRevoke)
}

func TestEventUserAuthorizationRevokeNoUser(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventUserAuthorizationRevoke(func(event twitch.EventUserAuthorizationRevoke, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubUserAuthorizationRevoke, "nouser")
}

func TestEventUserUpdate(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventUserUpdate(func(event twitch.EventUserUpdate, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubUserUpdate)
}

func TestEventUserUpdateNoEmail(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventUserUpdate(func(event twitch.EventUserUpdate, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubUserUpdate, "noemail")
}

func TestEventChannelCharityCampaignDonate(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelCharityCampaignDonate(func(event twitch.EventChannelCharityCampaignDonate, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelCharityCampaignDonate)
}

func TestEventChannelCharityCampaignProgress(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelCharityCampaignProgress(func(event twitch.EventChannelCharityCampaignProgress, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelCharityCampaignProgress)
}

func TestEventChannelCharityCampaignStart(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelCharityCampaignStart(func(event twitch.EventChannelCharityCampaignStart, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelCharityCampaignStart)
}

func TestEventChannelCharityCampaignStop(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelCharityCampaignStop(func(event twitch.EventChannelCharityCampaignStop, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelCharityCampaignStop)
}

func TestEventChannelShieldModeBegin(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelShieldModeBegin(func(event twitch.EventChannelShieldModeBegin, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelShieldModeBegin)
}

func TestEventChannelShieldModeEnd(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelShieldModeEnd(func(event twitch.EventChannelShieldModeEnd, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelShieldModeEnd)
}

func TestEventChannelShoutoutCreate(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelShoutoutCreate(func(event twitch.EventChannelShoutoutCreate, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelShoutoutCreate)
}

func TestEventChannelShoutoutReceive(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelShoutoutReceive(func(event twitch.EventChannelShoutoutReceive, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelShoutoutReceive)
}

func TestEventChannelModerate(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelModerate(func(event twitch.EventChannelModerate, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelModerate)
}

func TestEventAutomodMessageHold(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventAutomodMessageHold(func(event twitch.EventAutomodMessageHold, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubAutomodMessageHold)
}

func TestEventAutomodMessageUpdate(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventAutomodMessageUpdate(func(event twitch.EventAutomodMessageUpdate, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubAutomodMessageUpdate)
}

func TestEventAutomodSettingsUpdate(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventAutomodSettingsUpdate(func(event twitch.EventAutomodSettingsUpdate, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubAutomodSettingsUpdate)
}

func TestEventAutomodTermsUpdate(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventAutomodTermsUpdate(func(event twitch.EventAutomodTermsUpdate, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubAutomodTermsUpdate)
}

func TestEventChannelChatUserMessageHold(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelChatUserMessageHold(func(event twitch.EventChannelChatUserMessageHold, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelChatUserMessageHold)
}

func TestEventChannelChatUserMessageUpdate(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelChatUserMessageUpdate(func(event twitch.EventChannelChatUserMessageUpdate, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelChatUserMessageUpdate)
}

func TestEventChannelChatClear(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelChatClear(func(event twitch.EventChannelChatClear, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelChatClear)
}

func TestEventChannelChatClearUserMessages(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelChatClearUserMessages(func(event twitch.EventChannelChatClearUserMessages, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelChatClearUserMessages)
}

func TestEventChannelChatMessage(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelChatMessage(func(event twitch.EventChannelChatMessage, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelChatMessage)
}

func TestEventChannelChatMessageDelete(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelChatMessageDelete(func(event twitch.EventChannelChatMessageDelete, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelChatMessageDelete)
}

func TestEventChannelChatNotification(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelChatNotification(func(event twitch.EventChannelChatNotification, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelChatNotification)
}

func TestEventChannelChatSettingsUpdate(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelChatSettingsUpdate(func(event twitch.EventChannelChatSettingsUpdate, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelChatSettingsUpdate)
}

func TestEventChannelSuspiciousUserMessage(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelSuspiciousUserMessage(func(event twitch.EventChannelSuspiciousUserMessage, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelSuspiciousUserMessage)
}

func TestEventChannelSuspiciousUserUpdate(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelSuspiciousUserUpdate(func(event twitch.EventChannelSuspiciousUserUpdate, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelSuspiciousUserUpdate)
}

func TestEventChannelSharedChatBegin(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelSharedChatBegin(func(event twitch.EventChannelSharedChatBegin, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelSharedChatBegin)
}

func TestEventChannelSharedChatUpdate(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelSharedChatUpdate(func(event twitch.EventChannelSharedChatUpdate, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelSharedChatUpdate)
}

func TestEventChannelSharedChatEnd(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelSharedChatEnd(func(event twitch.EventChannelSharedChatEnd, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelSharedChatEnd)
}

func TestEventUserWhisperMessage(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventUserWhisperMessage(func(event twitch.EventUserWhisperMessage, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubUserWhisperMessage)
}

func TestEventChannelAdBreakBegin(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelAdBreakBegin(func(event twitch.EventChannelAdBreakBegin, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelAdBreakBegin)
}

func TestEventChannelWarningAcknowledge(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelWarningAcknowledge(func(event twitch.EventChannelWarningAcknowledge, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelWarningAcknowledge)
}

func TestEventChannelWarningSend(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelWarningSend(func(event twitch.EventChannelWarningSend, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelWarningSend)
}

func TestEventChannelUnbanRequestCreate(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelUnbanRequestCreate(func(event twitch.EventChannelUnbanRequestCreate, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelUnbanRequestCreate)
}

func TestEventChannelUnbanRequestResolve(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventChannelUnbanRequestResolve(func(event twitch.EventChannelUnbanRequestResolve, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubChannelUnbanRequestResolve)
}

func TestEventConduitShardDisabled(t *testing.T) {
	t.Parallel()

	assertSpecificEventOccurred(t, func(client *twitch.Client, ch chan struct{}) {
		client.OnEventConduitShardDisabled(func(event twitch.EventConduitShardDisabled, msg twitch.NotificationMessage) {
			close(ch)
		})
	}, twitch.SubConduitShardDisabled)
}
