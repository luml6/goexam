package services

import (
	"notepad-api/model"
	"notepad-api/utils"
	"time"

	"github.com/astaxie/beego"
)

//VoteGet get one voted note's content
func VoteGet(Ob map[string]interface{}) *utils.Response {

	var path, noteID string
	if _, ok := Ob["Path"].(string); ok {
		path = Ob["Path"].(string)
		share := model.NewShare(path)
		shareNote, err := share.Get()
		if shareNote == nil || err != nil {
			return utils.SHARE_ID_NOT_EXIST
		}
		noteID = shareNote.NoteID
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	voteresp := GetVote(noteID)
	return utils.NewResponse(0, "", voteresp)
}

//VoteGet get one voted note's content
func VoteGetByNote(Ob map[string]interface{}, userId string) *utils.Response {

	var noteID string
	if _, ok := Ob["Path"].(string); ok {
		noteID = Ob["Path"].(string)
		note := model.NewNote(noteID)
		note.UID = userId
		if !note.IsExist() {
			return utils.NOTE_NOT_EXIST
		}
	} else {
		return utils.WITHOUT_PARAMETERS
	}

	voteresp := GetVote(noteID)
	return utils.NewResponse(0, "", voteresp)
}

//VoteAdd Add one vote
func VoteAdd(Ob map[string]interface{}, userId string) *utils.Response {
	var (
		vote     *model.Vote
		answer   []interface{}
		note     *model.Note
		ID, path string
	)
	vote = model.NewVote("")
	if _, ok := Ob["Path"].(string); ok {
		path = Ob["Path"].(string)
		note = model.NewNote(path)
		note.UID = userId
		notebook, err := note.Get()
		if notebook == nil || err != nil {
			return utils.NOTE_NOT_EXIST
		}
		vote.NoteID = path
		if votenote, err := vote.GetByNoteID(); votenote != nil && err == nil {
			return utils.NOTE_ALREADY_VOTE
		}
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	ID = utils.CreateId()
	vote = model.NewVote(ID)
	vote.NoteID = path
	if _, ok := Ob["VoteType"].(float64); ok {
		vote.Type = uint8(Ob["VoteType"].(float64))
		if vote.Type != 0 && vote.Type != 1 {
			return utils.VOTE_TYPE_NOT_EXIST
		}
	} else {
		//		vote.Type = 0
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["InvalidTime"]; ok {
		vote.InvalidTime = utils.GetTime(Ob["InvalidTime"])
		//		vote.InvalidTime = int(Ob["InvalidTime"].(float64))
	}
	if _, ok := Ob["VoteQuestion"].(string); ok {
		question := Ob["VoteQuestion"].(string)
		if question == "" {
			return utils.VOTE_QUESTION_NOT_NULL
		}
		vote.Question = question
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["VoteSelect"].([]interface{}); ok {
		answer = Ob["VoteSelect"].([]interface{})
		if len(answer) <= 1 || len(answer) > 20 {
			return utils.VOTE_OPTION_LIMIT
		}
		if utils.CheckSlice(&answer) {
			return utils.VOTE_OPTION_NOT_SAME
		}
		if utils.CheckStingsNotNil(answer) {
			return utils.VOTE_OPTION_NOT_NULL
		}
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if err := vote.Add(); err != nil {

		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	err := note.SetVote(model.NOTE_SHARE)
	if err != nil {
		beego.Error(err)
	}

	if err = addvoteoptions(answer, ID); err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	t := struct {
		VoteID string
	}{ID}
	return utils.NewResponse(0, "", t)

}

func addvoteoptions(options []interface{}, voteID string) error {
	var err error
	for i := 0; i < len(options); i++ {
		ID := utils.CreateId()
		voteOption := model.NewVoteOption(ID)
		voteOption.Option, _ = options[i].(string)
		voteOption.VoteID = voteID
		err = voteOption.Add()
	}
	return err
}
func VoteUpdate(Ob map[string]interface{}, userId string) *utils.Response {
	var (
		vote   *model.Vote
		path   string
		file   []string
		answer []interface{}
	)
	if _, ok := Ob["Path"].(string); ok {
		path = Ob["Path"].(string)
		note := model.NewNote(path)
		note.UID = userId
		notebook, err := note.Get()
		if notebook == nil || err != nil {
			return utils.NOTE_NOT_EXIST
		}
		if notebook.IsShare == 1 {
			return utils.NOTE_ALREADY_SHARE
		}
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	vote = model.NewVote("")
	vote.NoteID = path
	votenote, err := vote.GetByNoteID()
	if votenote == nil || err != nil {
		return utils.VOTE_ID_NOT_EXIST
	}
	if _, ok := Ob["InvalidTime"].(float64); ok {
		vote.InvalidTime = utils.GetTime(Ob["InvalidTime"])
		file = append(file, "InvalidTime")
	}
	if _, ok := Ob["VoteType"].(float64); ok {
		vote.Type = uint8(Ob["VoteType"].(float64))
	}
	if _, ok := Ob["VoteQuestion"].(string); ok {
		question := Ob["VoteQuestion"].(string)
		vote.Question = question
		file = append(file, "Question")
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["VoteSelect"].([]interface{}); ok {
		answer = Ob["VoteSelect"].([]interface{})
		if len(answer) <= 1 || len(answer) > 20 {
			return utils.VOTE_OPTION_LIMIT
		}
		if utils.CheckSlice(&answer) {
			return utils.VOTE_OPTION_NOT_SAME
		}
		if utils.CheckStingsNotNil(answer) {
			return utils.VOTE_OPTION_NOT_NULL
		}
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	file = append(file, "CreateTime")
	vote.CreateTime = utils.NowSecond()
	if err := vote.Update(file); err != nil {
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	} else {
		options := model.NewOptionList(votenote.VoteID)
		if err = options.RemoveAll(); err == nil {
			if err = addvoteoptions(answer, votenote.VoteID); err != nil {
				return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
			}
		}
		beego.Debug(err)
	}
	t := struct {
		VoteID string
	}{votenote.VoteID}
	return utils.NewResponse(0, "", t)
}

//ClickVote get vote note url
func ClickVote(Ob map[string]interface{}) *utils.Response {
	var (
		vote       *model.Vote
		option     *model.VoteOption
		path       string
		clickCount int
		optionID   []interface{}
	)
	filed := []string{"OptionCount"}
	if _, ok := Ob["VoteID"].(string); ok {
		path = Ob["VoteID"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	if _, ok := Ob["SelectID"].([]interface{}); ok {
		optionID = Ob["SelectID"].([]interface{})
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	vote = model.NewVote(path)
	if _, err := vote.Get(); err != nil {
		return utils.VOTE_ID_NOT_EXIST
	}
	clickNum := vote.ClickNum
	vote.ClickNum = clickNum + 1

	if time.Now().Sub(vote.InvalidTime) > 0 {
		return utils.VOTE_OPTION_TIMEOUT
	}
	if vote.Type == 0 && len(optionID) > 1 {
		return utils.VOTE_OPTION_NOT_ONEMORE
	}
	for i := 0; i < len(optionID); i++ {
		option = model.NewVoteOption(optionID[i].(string))
		if _, err := option.Get(); err == nil {
			clickCount = option.OptionCount
		} else {
			beego.Error(err)
			return utils.OPTION_ID_NOT_EXIST
		}
		clickCount += 1
		option.OptionCount = clickCount
		if err := option.Update(filed); err != nil {
			beego.Error(err)
			return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
		}
		filed1 := []string{"ClickNum"}
		if err := vote.Update(filed1); err != nil {
			beego.Error(err)
		}
	}

	t := struct {
		IsSuccess bool
	}{true}
	return utils.NewResponse(0, "", t)
}

//CancleVote cancle one note vote
func CancleVote(Ob map[string]interface{}, userId string) *utils.Response {
	var vote *model.Vote
	var noteID string
	if _, ok := Ob["Path"].(string); ok {
		noteID = Ob["Path"].(string)
	} else {
		return utils.WITHOUT_PARAMETERS
	}
	vote = model.NewVote("")
	vote.NoteID = noteID
	votenote, err := vote.GetByNoteID()
	if votenote == nil || err != nil {
		return utils.NOTE_NOT_VOTE
	}
	voteID := votenote.VoteID
	err = votenote.Remove()
	if err != nil {
		beego.Error(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	options := model.NewOptionList(voteID)
	if err = options.RemoveAll(); err != nil {
		beego.Error(err)
	}
	note := model.NewNote(noteID)
	note.UID = userId
	err = note.SetVote(model.NOTE_NOT_SHARE)
	if err != nil {
		beego.Debug(err)
		return utils.NewResponse(utils.SYSTEM_CODE, err.Error(), nil)
	}
	t := struct {
		IsSuccess bool
	}{true}
	return utils.NewResponse(0, "", t)
}

func GetVote(notePath string) *utils.VoteResp {
	var voteresp utils.VoteResp
	vote := model.NewVote("")
	vote.NoteID = notePath
	votenote, err := vote.GetByNoteID()
	if votenote == nil || err != nil {
		beego.Error(err)
		return nil
	} else {
		voteresp.VoteID = votenote.VoteID
		voteresp.VoteQuestion = votenote.Question
		voteresp.VoteType = votenote.Type
		voteresp.VoteCount = votenote.ClickNum
	}
	voteoption := model.NewOptionList(votenote.VoteID)
	options, err := voteoption.GetList()
	if err != nil {
		beego.Error(err)
	} else {
		voteresp.VoteSelect = options
	}
	return &voteresp
}
