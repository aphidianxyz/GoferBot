package message

type TextMessage struct {
	id int64
	text string
	photo any // todo 
	reply Message
}

func (tm *TextMessage) GetID() {

}

func (tm *TextMessage) GetText()  {
	
}

func (tm *TextMessage) GetPhoto()  {
	
}

func (tm *TextMessage) GetReply()  {
	
}
