package api

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const DefaultURL = "http://localhost:8080/api/v1/sink/"

const Encrypt = "encrypt"
const Encrypted = "encrypted"

const Decrypt = "decrypt"
const Decrypted = "decrypted"

type ResponseData struct {
	Data struct {
		Text   string `json: "text"`
		Status string `json: "status"`
		Sink   string `json: "sink"`
	} `json:"data"`
}

type EncryptResponse struct {
	Text   string `json: "text"`
	Status string `json: "status"`
	Sink   string `json: "sink"`
}
type Sink struct {
	Name      string
	Text      string
	Operation string
}

type Login struct {
	Sink    *Sink
	Token   string
	BaseURL string
}

func NewSink(name string) *Sink {
	return &Sink{Name: name}
}

func (s *Sink) SetOperation(operation string) {
	s.Operation = operation
}

func (s *Sink) SetText(text string) {
	s.Text = text
}

func (l *Login) Connect(token string) {
	l.BaseURL = DefaultURL
	l.setToken(token)
}

func (l *Login) SetBaseURL(url string) {
	l.BaseURL = url
}

func (l *Login) AssignSink(s *Sink) {
	l.Sink = s
}

func (l *Login) Encrypt(payload string) (ResponseData, error) {

	l.Sink.SetText(b64.StdEncoding.EncodeToString([]byte(payload)))
	l.Sink.SetOperation(Encrypt)
	b, err := json.Marshal(&l.Sink)
	if err != nil {
		return ResponseData{}, err
	}

	req, err := http.NewRequest("POST", l.BaseURL+l.Sink.Name, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	var bearer = "Bearer " + l.Token
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ResponseData{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ResponseData{}, err
	}
	var er ResponseData
	err = json.Unmarshal(body, &er)
	if err != nil {
		return ResponseData{}, err
	}
	return er, nil
}

func (l *Login) Decrypt(payload string) (ResponseData, error) {
	l.Sink.SetText(payload)
	l.Sink.SetOperation(Decrypt)

	b, err := json.Marshal(&l.Sink)
	if err != nil {
		return ResponseData{}, err
	}

	req, err := http.NewRequest("POST", l.BaseURL+l.Sink.Name, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	var bearer = "Bearer " + l.Token
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ResponseData{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ResponseData{}, err
	}
	var er ResponseData
	err = json.Unmarshal(body, &er)
	if err != nil {
		return ResponseData{}, err
	}
	b, err = b64.StdEncoding.DecodeString(er.Data.Text)
	if err != nil {
		return ResponseData{}, err
	}

	er.Data.Text = string(b)
	return er, nil
}

func (l *Login) setToken(token string) {
	l.Token = token
}
