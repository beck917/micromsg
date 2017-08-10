package entities

import (
	"time"
)

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(j).Format("2006-01-02 15:04:05") + `"`), nil
}

func (j *JsonTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+"2006-01-02 15:04:05"+`"`, string(data), time.Local)
	*j = JsonTime(now)
	return
}

func (j JsonTime) MarshalMsgpack() ([]byte, error) {
	return []byte(`"` + time.Time(j).Format("2006-01-02 15:04:05") + `"`), nil
}

func (j *JsonTime) UnmarshalMsgpack(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+"2006-01-02 15:04:05"+`"`, string(data), time.Local)
	*j = JsonTime(now)
	return
}

func (j JsonTime) String() string {
	return time.Time(j).Format("2006-01-02 15:04:05")
}
