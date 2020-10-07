/*
Copyright 2018 Blindside Networks

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/mattermost"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	bbbAPI "github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/api"
	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/dataStructs"
	"github.com/mattermost/mattermost-server/model"
)

type RequestProfilesJSON struct {
	UserId    string `json:"user_id"`
	ChannelId string `json:"channel_id"`
	Topic     string `json:"title"`
	Desc      string `json:"description"`
}
type RequestCreateMeetingJSON struct {
	UserId    string `json:"user_id"`
	ChannelId string `json:"channel_id"`
	Topic     string `json:"title"`
	Desc      string `json:"description"`
}

type RequestUserProfileJSON struct {
	UserId string `json:"user_id"`
}
type RequestMultipleUserProfileJSON struct {
	UserIds []string `json:"user_ids"`
}
type ResponseUserProfileJSON struct {
	ProfileId   string `json:"profile_id"`
	UserId      string `json:"user_id"`
	Twitter     string `json:"twitter"`
	Age         string `json:"age"`
	LivingPlace string `json:"livingPlace"`
	Pronoun     string `json:"pronoun"`
	MagicPower  string `json:"magicPower"`
	Food        string `json:"food"`
	Misc        string `json:"misc"`
}

type UpdateUserProfileJSON struct {
	UserId string `json:"user_id"`
	Field  string `json:"field"`
	Value  string `json:"val"`
}

type SpeedDatingCreateJSON struct {
	CreatorId       string   `json:"creator_id"`
	RoomDisplayName string   `json:"room_display_name"`
	UserIds         []string `json:"user_ids"`
	TeamId          string   `json:"team_id"`
	ExcludeUserIds  []string `json:"excluded_user_ids"`
	UsersPerRoom    int      `json:"users_per_room"`
	Duration        int      `json:duration`
}

func (p *Plugin) handleProfiles(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	var request RequestProfilesJSON
	json.Unmarshal(body, &request)

	p.createProfilesPost(request.UserId, request.ChannelId)
}

//Create meeting doesn't call the BBB api to start a meeting
//Only populates the meeting with details. Meeting is started when first person joins
func (p *Plugin) handleCreateMeeting(w http.ResponseWriter, r *http.Request) {

	// reads in information to create a meeting from client inside
	// whats being read in is the stuff in RequestCreateMeetingJSON
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	var request RequestCreateMeetingJSON
	json.Unmarshal(body, &request)

	meetingpointer := new(dataStructs.MeetingRoom)
	var err error
	if request.Topic == "" {
		err = p.PopulateMeeting(meetingpointer, nil, request.Desc)
	} else {
		err = p.PopulateMeeting(meetingpointer, []string{"create", request.Topic}, request.Desc)
	}

	if err != nil {
		http.Error(w, "Please provide a 'Site URL' in Settings > General > Configuration.", http.StatusUnprocessableEntity)
		return
	}

	//creates the start meeting post
	p.createStartMeetingPost(request.UserId, request.ChannelId, meetingpointer)

	// add our newly created meeting to our array of meetings
	p.Meetings = append(p.Meetings, *meetingpointer)

	w.WriteHeader(http.StatusOK)
}

type ButtonRequestJSON struct {
	UserId    string `json:"user_id"`
	MeetingId string `json:"meeting_id"`
	IsMod     string `json:"is_mod"`
}

type ButtonResponseJSON struct {
	Url string `json:"url"`
}

func (p *Plugin) handleJoinMeeting(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request ButtonRequestJSON
	json.Unmarshal(body, &request)
	meetingID := request.MeetingId
	meetingpointer := p.FindMeeting(meetingID)

	if meetingpointer == nil {
		myresp := ButtonResponseJSON{
			Url: "error",
		}
		userJson, _ := json.Marshal(myresp)
		w.Write(userJson)
		return
	} else {
		//check if meeting has actually been created and can be joined
		if !meetingpointer.Created {
			bbbAPI.CreateMeeting(meetingpointer)
			meetingpointer.Created = true
			var fullMeetingInfo dataStructs.GetMeetingInfoResponse
			bbbAPI.GetMeetingInfo(meetingID, meetingpointer.ModeratorPW_, &fullMeetingInfo) // this is used to get the InternalMeetingID
			meetingpointer.InternalMeetingId = fullMeetingInfo.InternalMeetingID
			meetingpointer.CreatedAt = time.Now().Unix()
		}

		user, _ := p.API.GetUser(request.UserId)
		username := user.Username

		//golang doesnt have sets so have to iterate through array to check if meeting participant is already in meeeting
		if !IsItemInArray(username, meetingpointer.AttendeeNames) {
			meetingpointer.AttendeeNames = append(meetingpointer.AttendeeNames, username)
		}

		var participant = dataStructs.Participants{} //set participant as an empty struct of type Participants
		participant.FullName_ = username
		participant.MeetingID_ = meetingID

		post, appErr := p.API.GetPost(meetingpointer.PostId)
		if appErr != nil {
			http.Error(w, appErr.Error(), appErr.StatusCode)
			return
		}
		config := p.config()
		if config.AdminOnly {
			participant.Password_ = meetingpointer.AttendeePW_
			if request.UserId == post.UserId {
				mattermost.API.LogInfo("userID same as post userID")
				participant.Password_ = meetingpointer.ModeratorPW_ // the creator of a room is always moderator
			} else {
				for _, role := range user.GetRoles() {
					if role == "SYSTEM_ADMIN" || role == "TEAM_ADMIN" {
						mattermost.API.LogInfo("user is system or team administrator")
						participant.Password_ = meetingpointer.ModeratorPW_
						break
					}
				}
			}
		} else {
			mattermost.API.LogInfo("everyone should be a moderator")
			participant.Password_ = meetingpointer.ModeratorPW_ //make everyone in channel a mod
		}
		joinURL, err := bbbAPI.GetJoinURL(&participant)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		myresp := ButtonResponseJSON{
			Url: joinURL,
		}
		userJson, _ := json.Marshal(myresp)

		Length, attendantsarray := GetAttendees(meetingID, meetingpointer.ModeratorPW_)
		// we immediately add our current attendee thats trying to join the meeting
		// to avoid the delay
		attendantsarray = append(attendantsarray, username)
		post.Props["user_count"] = Length + 1
		post.Props["attendees"] = strings.Join(attendantsarray, ",")

		if _, err := p.API.UpdatePost(post); err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(userJson)
	}
}

//this method is responsible for updating meeting has ended inside mattermost when
// we end our meeting from inside BigBlueButton
func (p *Plugin) handleImmediateEndMeetingCallback(w http.ResponseWriter, r *http.Request, path string) {

	startpoint := len("/meetingendedcallback?")
	endpoint := strings.Index(path, "&")
	meetingid := path[startpoint:endpoint]
	validation := path[endpoint+1:]
	meetingpointer := p.FindMeeting(meetingid)
	if meetingpointer == nil || meetingpointer.ValidToken != validation {
		w.WriteHeader(http.StatusOK)
		return
	}
	post, err := p.API.GetPost(meetingpointer.PostId)
	if err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}
	if meetingpointer.EndedAt == 0 {
		meetingpointer.EndedAt = time.Now().Unix()
	}
	p.MeetingsWaitingforRecordings = append(p.MeetingsWaitingforRecordings, *meetingpointer)
	post.Props["meeting_status"] = "ENDED"
	post.Props["attendents"] = strings.Join(meetingpointer.AttendeeNames, ",")
	timediff := meetingpointer.EndedAt - meetingpointer.CreatedAt
	durationstring := FormatSeconds(timediff)
	post.Props["duration"] = durationstring

	p.API.UpdatePost(post)

	w.WriteHeader(http.StatusOK)
}

//when user clicks endmeeting button inside Mattermost
func (p *Plugin) handleEndMeeting(w http.ResponseWriter, r *http.Request) {

	//for debugging
	mattermost.API.LogInfo("Processing End Meeting Request")

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request ButtonRequestJSON
	json.Unmarshal(body, &request)
	meetingID := request.MeetingId
	meetingpointer := p.FindMeeting(meetingID)

	user, _ := p.API.GetUser(request.UserId)
	username := user.Username

	if meetingpointer == nil {
		myresp := model.PostActionIntegrationResponse{
			EphemeralText: "meeting has already ended",
		}
		userJson, _ := json.Marshal(myresp)
		w.Write(userJson)
		return
	} else {
		bbbAPI.EndMeeting(meetingpointer.MeetingID_, meetingpointer.ModeratorPW_)
		//for debugging
		mattermost.API.LogInfo("Meeting Ended")

		if meetingpointer.EndedAt == 0 {
			meetingpointer.EndedAt = time.Now().Unix()
		}
		p.MeetingsWaitingforRecordings = append(p.MeetingsWaitingforRecordings, *meetingpointer)

		post, err := p.API.GetPost(meetingpointer.PostId)
		if err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}

		post.Props["meeting_status"] = "ENDED"
		post.Props["attendents"] = strings.Join(meetingpointer.AttendeeNames, ",")
		post.Props["ended_by"] = username
		timediff := meetingpointer.EndedAt - meetingpointer.CreatedAt
		if meetingpointer.CreatedAt == 0 {
			timediff = 0
		}
		durationstring := FormatSeconds(timediff)
		post.Props["duration"] = durationstring

		if _, err := p.API.UpdatePost(post); err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

type isRunningRequestJSON struct {
	MeetingId string `json:"meeting_id"`
}

type isRunningResponseJSON struct {
	IsRunning bool `json:"running"`
}

func (p *Plugin) handleIsMeetingRunning(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request ButtonRequestJSON
	json.Unmarshal(body, &request)
	meetingID := request.MeetingId

	resp, _ := bbbAPI.IsMeetingRunning(meetingID)

	myresp := isRunningResponseJSON{
		IsRunning: resp,
	}
	userJson, _ := json.Marshal(myresp)

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJson)

}

// type WebHookRequestJSON struct {
// 	Header struct {
// 		Timestamp   string `json:"timestamp"`
// 		Name        string `json:"name"`
// 		CurrentTime string `json:"current_time"`
// 		Version     string `json:"version"`
// 	} `json:"header"`
// 	Payload struct {
// 		MeetingId string `json:"meeting_id"`
// 	} `json:"payload"`
// }
//
// type WebHookResponseEncoded struct {
// 	Payload map[string]interface{} `form:"payload"`
// }

//webhook to send additional information about meeting that had ended
//has a 4-5 minute delay which is why handleImmediateEndMeetingCallback() is
//used instead for updating end meeting post on Mattermost. Keeping it here in
//case bbbserver is not up to date with the immediate meeting ended callback feature
// func (p *Plugin) handleWebhookMeetingEnded(w http.ResponseWriter, r *http.Request) {
//
// 	out := ""
// 	r.ParseForm()
// 	for key, value := range r.Form {
// 		out += fmt.Sprintf("%s = %s\n", key, value)
// 	}
// 	events := (r.FormValue("event"))
//
// 	internal_meetingid := events[strings.Index(events, "\""+"meeting_id"+"\"")+14:]
// 	internal_meetingid = internal_meetingid[:strings.IndexByte(internal_meetingid, '"')]
//
// 	meetingpointer := p.FindMeetingfromInternal(internal_meetingid)
//
// 	if meetingpointer == nil {
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}
//
// 	postid := meetingpointer.PostId
// 	if postid == "" {
// 		panic("no post id found")
// 	}
// 	post, err := p.API.GetPost(postid)
// 	if err != nil {
// 		http.Error(w, err.Error(), err.StatusCode)
// 		return
// 	}
// 	if meetingpointer.EndedAt == 0 {
// 		meetingpointer.EndedAt = time.Now().Unix()
// 	}
// 	p.MeetingsWaitingforRecordings = append(p.MeetingsWaitingforRecordings, *meetingpointer)
// 	post.Props["meeting_status"] = "ENDED"
// 	post.Props["attendents"] = strings.Join(meetingpointer.AttendeeNames, ",")
// 	timediff := meetingpointer.EndedAt - meetingpointer.CreatedAt
// 	durationstring := FormatSeconds(timediff)
// 	post.Props["duration"] = durationstring
//
// 	if _, err := p.API.UpdatePost(post); err != nil {
// 		http.Error(w, err.Error(), err.StatusCode)
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusOK)
// }

// type MyCustomClaims struct {
// 	MeetingID string `json:"meeting_id"`
// 	RecordID  string `json:"record_id"`
// 	jwt.StandardClaims
// }

func (p *Plugin) handleRecordingReady(w http.ResponseWriter, r *http.Request) {
	// p.API.LogDebug("handleRecordingReady reached")
	// r.ParseForm()
	// parameters := (r.FormValue("signed_parameters"))
	// token, _ := jwt.ParseWithClaims(parameters, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte("AllYourBase"), nil
	// })
	// claims, _ := token.Claims.(*MyCustomClaims)
	// meetingid := claims.MeetingID
	// recordid := claims.RecordID
	// p.API.LogDebug(meetingid + " " + recordid)
	// recordingsresponse, _ := bbbAPI.GetRecordings(meetingid, recordid, "")
	// if recordingsresponse.ReturnCode != "SUCCESS" {
	// 	w.WriteHeader(http.StatusOK)
	// 	return
	// }
	//
	// meetingpointer := p.FindMeeting(meetingid)
	//
	// if meetingpointer == nil {
	// 	w.WriteHeader(http.StatusOK)
	// 	return
	// }
	//
	// postid := meetingpointer.PostId
	// if postid == "" {
	// 	panic("no post id found")
	// }
	// post, err := p.API.GetPost(postid)
	// if err != nil {
	// 	http.Error(w, err.Error(), err.StatusCode)
	// 	return
	// }
	//
	// post.Message = "#BigBlueButton #" + meetingpointer.Name_ + " #" + recordid + " #recording" + " recordings"
	// post.Props["recording_status"] = "COMPLETE"
	// post.Props["is_published"] = "true"
	// post.Props["record_id"] = recordid
	// post.Props["recording_url"] = recordingsresponse.Recordings.Recording[0].Playback.Format[0].Url
	//
	// post.Props["images"] = strings.Join(recordingsresponse.Recordings.Recording[0].Playback.Format[0].Images, ",")
	// if _, err := p.API.UpdatePost(post); err != nil {
	// 	http.Error(w, err.Error(), err.StatusCode)
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
	return
}

type AttendeesRequestJSON struct {
	MeetingId string `json:"meeting_id"`
}

type AttendeesResponseJSON struct {
	Num       int      `json:"num"`
	Attendees []string `json:"attendees"`
}

func (p *Plugin) handleGetAttendeesInfo(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request AttendeesRequestJSON
	json.Unmarshal(body, &request)
	meetingID := request.MeetingId
	meetingpointer := p.FindMeeting(meetingID)
	if meetingpointer == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	postid := meetingpointer.PostId
	if postid == "" {
		w.WriteHeader(http.StatusOK)
		return
	}
	post, err := p.API.GetPost(postid)
	if err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}
	Length, Array := GetAttendees(meetingID, meetingpointer.ModeratorPW_)
	post.Props["user_count"] = Length
	post.Props["attendees"] = strings.Join(Array, ",")

	if _, err := p.API.UpdatePost(post); err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}
	myresp := AttendeesResponseJSON{
		Num:       Length,
		Attendees: Array,
	}
	userJson, _ := json.Marshal(myresp)

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJson)
}

type RecordingsRequestJSON struct {
	ChannelId string `json:"channel_id"`
}

type RecordingsResponseJSON struct {
	Recordings []SingleRecording `json:"recordings"`
}

type SingleRecording struct {
	RecordingUrl string `json:"recordingurl"`
	Title        string `json:"title"`
}

type PublishRecordingsRequestJSON struct {
	RecordId  string `json:"record_id"`
	Publish   string `json:"publish"` //string  true or false
	MeetingId string `json:"meeting_id"`
}

func (p *Plugin) handlePublishRecordings(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request PublishRecordingsRequestJSON
	json.Unmarshal(body, &request)
	recordid := request.RecordId
	publish := request.Publish

	meetingpointer := p.FindMeeting(request.MeetingId)
	if meetingpointer == nil {
		http.Error(w, "Error: Cannot find the meeting_id for the recording, MeetingID#"+request.MeetingId, http.StatusForbidden)
		return
	}

	if _, err := bbbAPI.PublishRecordings(recordid, publish); err != nil {
		http.Error(w, "Error: Recording not found", http.StatusForbidden)
		return
	}

	post, appErr := p.API.GetPost(meetingpointer.PostId)
	if appErr != nil {
		http.Error(w, "Error: cannot find the post message for this recording \n"+appErr.Error(), appErr.StatusCode)
		return
	}

	post.Props["is_published"] = publish

	if _, err := p.API.UpdatePost(post); err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}
	//update post props with new recording  status
	w.WriteHeader(http.StatusOK)
}

type DeleteRecordingsRequestJSON struct {
	RecordId  string `json:"record_id"`
	MeetingId string `json:"meeting_id"`
}

func (p *Plugin) handleDeleteRecordings(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request DeleteRecordingsRequestJSON
	json.Unmarshal(body, &request)
	recordid := request.RecordId

	if _, err := bbbAPI.DeleteRecordings(recordid); err != nil {
		http.Error(w, "Error: Recording not found", http.StatusForbidden)
		return
	}

	meetingID := request.MeetingId
	meetingpointer := p.FindMeeting(meetingID)
	if meetingpointer == nil {
		http.Error(w, "Error: Cannot find the meeting_id for the recording", http.StatusForbidden)
		return
	}

	post, appErr := p.API.GetPost(meetingpointer.PostId)
	if appErr != nil {
		http.Error(w, "Error: cannot find the post message for this recording \n"+appErr.Error(), appErr.StatusCode)
		return
	}

	post.Props["is_deleted"] = "true"
	post.Props["record_status"] = "Recording Deleted"
	if _, err := p.API.UpdatePost(post); err != nil {
		http.Error(w, "Error: could not update post \n"+err.Error(), err.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *Plugin) Loopthroughrecordings() {

	for i := 0; i < len(p.MeetingsWaitingforRecordings); i++ {
		Meeting := p.MeetingsWaitingforRecordings[i]
		// TODO Harshil Sharma: explore better alternative of waiting for specific count of re-tries
		// instead of duration of re-tries.
		if Meeting.LoopCount > 144 {
			p.MeetingsWaitingforRecordings = append(p.MeetingsWaitingforRecordings[:i], p.MeetingsWaitingforRecordings[i+1:]...)
			i--
			continue
		}

		recordingsresponse, _, _ := bbbAPI.GetRecordings(Meeting.MeetingID_, "", "")
		if recordingsresponse.ReturnCode == "SUCCESS" {
			if len(recordingsresponse.Recordings.Recording) > 0 {
				postid := Meeting.PostId
				if postid != "" {
					post, _ := p.API.GetPost(postid)
					post.Message = "#BigBlueButton #" + Meeting.Name_ + " #" + Meeting.MeetingID_ + " #recording" + " recordings"
					post.Props["recording_status"] = "COMPLETE"
					post.Props["is_published"] = "true"
					post.Props["record_id"] = recordingsresponse.Recordings.Recording[0].RecordID
					post.Props["recording_url"] = recordingsresponse.Recordings.Recording[0].Playback.Format[0].Url
					post.Props["images"] = strings.Join(recordingsresponse.Recordings.Recording[0].Playback.Format[0].Images, ",")

					if _, err := p.API.UpdatePost(post); err == nil {
						p.MeetingsWaitingforRecordings = append(p.MeetingsWaitingforRecordings[:i], p.MeetingsWaitingforRecordings[i+1:]...)
						i--
					}
				}
			}
		}
	}
}

const profilesPrefix = "bbb_profiles_"

func NewProfile(profileId string, userId string) *dataStructs.Profile {
	profile := dataStructs.Profile{
		ID:          profileId,
		User:        userId,
		CreatedAt:   model.GetMillis(),
		Twitter:     "",
		Age:         "",
		LivingPlace: "",
		Pronoun:     "",
		MagicPower:  "",
		Food:        "",
		Misc:        "",
	}
	return &profile
}

func (p *Plugin) handleGetMultipleProfilesInfo(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request RequestMultipleUserProfileJSON
	json.Unmarshal(body, &request)

	var responseArray = make([]ResponseUserProfileJSON, 0, 0)

	for _, userId := range request.UserIds {

		profileId := profilesPrefix + userId

		var userProfile dataStructs.Profile
		_userProfile, _ := p.Store.Profiles().Get(profileId)
		if _userProfile == nil {
			userProfile = *NewProfile(profileId, userId)
		} else {
			userProfile = *_userProfile
		}

		resp := ResponseUserProfileJSON{
			ProfileId:   profileId,
			UserId:      userProfile.User,
			Twitter:     userProfile.Twitter,
			Age:         userProfile.Age,
			LivingPlace: userProfile.LivingPlace,
			Pronoun:     userProfile.Pronoun,
			MagicPower:  userProfile.MagicPower,
			Food:        userProfile.Food,
			Misc:        userProfile.Misc,
		}

		responseArray = append(responseArray, resp)
	}

	responseJsonArray, _ := json.Marshal(responseArray)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJsonArray)
}

func (p *Plugin) handleGetProfileInfo(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request RequestUserProfileJSON
	json.Unmarshal(body, &request)

	userId := request.UserId
	profileId := profilesPrefix + userId

	var userProfile dataStructs.Profile
	_userProfile, _ := p.Store.Profiles().Get(profileId)
	if _userProfile == nil {
		userProfile = *NewProfile(profileId, userId)
	} else {
		userProfile = *_userProfile
	}

	resp := ResponseUserProfileJSON{
		ProfileId:   profileId,
		UserId:      userProfile.User,
		Twitter:     userProfile.Twitter,
		Age:         userProfile.Age,
		LivingPlace: userProfile.LivingPlace,
		Pronoun:     userProfile.Pronoun,
		MagicPower:  userProfile.MagicPower,
		Food:        userProfile.Food,
		Misc:        userProfile.Misc,
	}

	profileJSON, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(profileJSON)
}

func (p *Plugin) createRoom(roomDisplayName string, teamId string, creatorId string, userIds []string, duration int) {

	p.API.LogInfo("Create room and meeting by " + creatorId + " in Team " + teamId + " for users " + strings.Join(userIds, ","))

	cname := "kennenlernen_" + model.NewId() //string(time.Now().Unix()) + "_" +  string(rand.Intn(1000))
	prepChannel := &model.Channel{
		DisplayName: roomDisplayName,
		Name:        cname,
		Type:        model.CHANNEL_PRIVATE,
		TeamId:      teamId,
		CreatorId:   creatorId,
	}
	//channel, err := p.API.GetGroupChannel(userIds)
	channel, err := p.API.CreateChannel(prepChannel)
	if err != nil {
		p.API.LogError("Cannot create group channel - " + err.Error())
		return
	}
	p.API.LogInfo("Created private room " + cname + " will now populate with users")
	for _, userId := range userIds {
		//_, err1 := p.API.AddUserToChannel(channel.Id, userId, creatorId)
		_, err1 := p.API.AddChannelMember(channel.Id, userId)
		if err1 != nil {
			p.API.LogError("Cannot add " + userId + " to the channel " + channel.Name)
		} else {
			p.API.LogInfo("Added " + userId + " to channel " + channel.Name)
		}
	}

	meetingpointer := new(dataStructs.MeetingRoom)
	err2 := p.PopulateMeeting(meetingpointer, nil, "Kennlern-Runde")

	if err2 != nil {
		//http.Error(w, "Please provide a 'Site URL' in Settings > General > Configuration.", http.StatusUnprocessableEntity)
		p.API.LogError("Cannot PopulateMeeting! - " + err.Error())
		return
	}
	meetingpointer.Duration = duration

	//creates the start meeting post
	p.createStartMeetingPost(creatorId, channel.Id, meetingpointer)

	// add our newly created meeting to our array of meetings
	p.Meetings = append(p.Meetings, *meetingpointer)
	//channel.Id
	var payload = map[string]interface{}{"event": "SpeeddatingChannelCreated", "channelId": channel.Id, "meetingId": meetingpointer.MeetingID_, "userIds": strings.Join(userIds, ",")}
	p.API.PublishWebSocketEvent("SpeeddatingChannelCreated", payload, &model.WebsocketBroadcast{
		ChannelId: channel.Id,
	})

	//add the creator
	_, err3 := p.API.AddChannelMember(channel.Id, creatorId)
	if err3 != nil {
		p.API.LogError("Cannot add Creator " + creatorId + " to the channel " + channel.Name)
	} else {
		p.API.LogInfo("Added Creator " + creatorId + " to channel " + channel.Name)
	}

}

func getNthKeyOf(m *map[string]bool, n int) (string, bool) {
	res := ""
	ok := false
	index := 0
	for key, _ := range *m {
		if index == n {
			res = key
			ok = true
			break
		}
		index++
	}
	return res, ok
}

func (p *Plugin) handleCreateSpeeddatingRooms(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	p.API.LogInfo("handleUpdateProfileInfo")

	var request SpeedDatingCreateJSON
	json.Unmarshal(body, &request)

	p.API.LogInfo("handleCreateSpeeddatingRooms request " + strings.Join(request.UserIds, ","))

	usersPerRoom := request.UsersPerRoom
	var usersLeftMap map[string]bool
	usersLeftMap = make(map[string]bool)

	for _, m := range request.UserIds {
		usersLeftMap[m] = true
	}

	/*
		p.API.LogInfo("Will exclude users of Channel " + request.ExcludeFromChannel)
		orgaChannel, err := p.API.GetChannelByName(request.TeamId, request.ExcludeFromChannel, false)
		if err == nil {
			members, err1 := p.API.GetChannelMembers(orgaChannel.Id, 0, 100)
			p.API.LogInfo("Will remove " + string(len(*members)) + " members from usersLeft" )
			if err1 == nil {
				for _, m := range *members {
					delete(usersLeftMap, m.UserId)
				}
			}
		}*/
	p.API.LogInfo("Will exclude users of Channel " + strings.Join(request.ExcludeUserIds, ","))
	for _, userId := range request.ExcludeUserIds {
		delete(usersLeftMap, userId)
	}

	//usersLeft := request.UserIds
	userCount := len(usersLeftMap)
	leftOverCount := userCount % usersPerRoom
	p.API.LogInfo(" there would be " + strconv.Itoa(leftOverCount) + " useres in a smaller room, will distribute")
	//roomsCount := int(math.Floor(float64(userCount) / float64(usersPerRoom)))
	//usersLeftCount := userCount % roomsCount
	roomIndex := 1
	for len(usersLeftMap) > 0 {
		uCount := usersPerRoom
		if len(usersLeftMap) <= usersPerRoom {
			uCount = len(usersLeftMap)
		} else {
			if leftOverCount > 0 {
				p.API.LogInfo("This time one more user for the lefties")
				uCount++
				leftOverCount--
			}
		}

		p.API.LogInfo("Will select " + strconv.Itoa(uCount) + " amount of users randomly")

		roomUserIds := []string{}
		for a := 0; a < uCount; a++ {
			user_index := 0
			if len(usersLeftMap) > 1 {
				user_index = rand.Intn(len(usersLeftMap) - 1)
			}
			userId, ok := getNthKeyOf(&usersLeftMap, user_index)
			if ok {
				roomUserIds = append(roomUserIds, userId)
				//usersLeft = append(usersLeft[:user_index], usersLeft[user_index + 1:]...)
				delete(usersLeftMap, userId)
			} else {
				p.API.LogError("Oops, no user forund at position " + strconv.Itoa(user_index) + " this should not happen")
			}
		}
		p.createRoom(request.RoomDisplayName+" "+strconv.Itoa(roomIndex), request.TeamId, request.CreatorId, roomUserIds, request.Duration)
		roomIndex++
	}

}

func (p *Plugin) handleUpdateProfileInfo(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	mattermost.API.LogInfo("handleUpdateProfileInfo")

	var request UpdateUserProfileJSON
	json.Unmarshal(body, &request)

	mattermost.API.LogInfo("handleUpdateProfileInfo request " + request.UserId)

	userId := request.UserId
	profileId := profilesPrefix + userId

	var userProfile dataStructs.Profile
	mattermost.API.LogInfo("p.Store.Profiles().Get(profileId)")
	_userProfile, _ := p.Store.Profiles().Get(profileId)
	mattermost.API.LogInfo("got UserProfile")
	if _userProfile == nil {
		userProfile = *NewProfile(profileId, userId)
		mattermost.API.LogInfo("got nil as profile")
		err := p.Store.Profiles().Insert(&userProfile)
		if err != nil {
			mattermost.API.LogError(err.Error() + "got error from set profile")
			//log.Error(err.Error() + " - Cannot store new Profile")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		userProfile = *_userProfile
	}

	userProfileNew := new(dataStructs.Profile)
	*userProfileNew = userProfile
	var changed = true
	if request.Field == "age" {
		userProfileNew.Age = request.Value
	} else if request.Field == "livingPlace" {
		userProfileNew.LivingPlace = request.Value
	} else if request.Field == "pronoun" {
		userProfileNew.Pronoun = request.Value
	} else if request.Field == "magicPower" {
		userProfileNew.MagicPower = request.Value
	} else if request.Field == "food" {
		userProfileNew.Food = request.Value
	} else if request.Field == "misc" {
		userProfileNew.Misc = request.Value
	} else {
		changed = false
	}

	if changed {
		userProfileNew.ID = profileId
		p.API.LogInfo("Something has changed - will update")
		out, errJ := json.Marshal(userProfileNew)
		if errJ == nil {
			p.API.LogInfo(string(out))
		}
		err := p.Store.Profiles().Update(&userProfile, userProfileNew)
		if err != nil {
			p.API.LogError(err.Error() + " - Cannot store new Profile")
			//log.Error(err.Error() + " - Cannot store new Profile")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		p.API.LogInfo("nothing has changed")
	}

	w.WriteHeader(http.StatusOK)
}
