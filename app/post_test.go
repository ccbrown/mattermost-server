// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/platform/model"
)

func TestUpdatePostEditAt(t *testing.T) {
	th := Setup().InitBasic()

	post := &model.Post{}
	*post = *th.BasicPost

	post.IsPinned = true
	if saved, err := UpdatePost(post, true); err != nil {
		t.Fatal(err)
	} else if saved.EditAt != post.EditAt {
		t.Fatal("shouldn't have updated post.EditAt when pinning post")

		*post = *saved
	}

	time.Sleep(time.Millisecond * 100)

	post.Message = model.NewId()
	if saved, err := UpdatePost(post, true); err != nil {
		t.Fatal(err)
	} else if saved.EditAt == post.EditAt {
		t.Fatal("should have updated post.EditAt when updating post message")
	}
}

func TestPostReplyToPostWhereRootPosterLeftChannel(t *testing.T) {
	// This test ensures that when replying to a root post made by a user who has since left the channel, the reply
	// post completes successfully. This is a regression test for PLT-6523.
	th := Setup().InitBasic()

	channel := th.BasicChannel
	userInChannel := th.BasicUser2
	userNotInChannel := th.BasicUser
	rootPost := th.BasicPost

	if _, err := AddUserToChannel(userInChannel, channel); err != nil {
		t.Fatal(err)
	}

	if err := RemoveUserFromChannel(userNotInChannel.Id, "", channel); err != nil {
		t.Fatal(err)
	}

	replyPost := model.Post{
		Message:       "asd",
		ChannelId:     channel.Id,
		RootId:        rootPost.Id,
		ParentId:      rootPost.Id,
		PendingPostId: model.NewId() + ":" + fmt.Sprint(model.GetMillis()),
		UserId:        userInChannel.Id,
		CreateAt:      0,
	}

	if _, err := CreatePostAsUser(&replyPost); err != nil {
		t.Fatal(err)
	}
}

func TestPostChannelMentions(t *testing.T) {
	th := Setup().InitBasic()

	channel := th.BasicChannel
	user := th.BasicUser

	channelToMention, err := CreateChannel(&model.Channel{
		DisplayName: "Mention Test",
		Name:        "mention-test",
		Type:        model.CHANNEL_OPEN,
		TeamId:      th.BasicTeam.Id,
	}, false)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer func() {
		PermanentDeleteChannel(channelToMention)
	}()

	_, err = AddUserToChannel(user, channel)
	require.Nil(t, err)

	post := &model.Post{
		Message:       fmt.Sprintf("hello, ~%v!", channelToMention.Name),
		ChannelId:     channel.Id,
		PendingPostId: model.NewId() + ":" + fmt.Sprint(model.GetMillis()),
		UserId:        user.Id,
		CreateAt:      0,
	}

	result, err := CreatePostAsUser(post)
	require.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"mention-test": map[string]interface{}{
			"display_name": "Mention Test",
		},
	}, result.Props["channel_mentions"])
}
