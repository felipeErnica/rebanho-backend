package pregnancytest

import "github.com/google/uuid"

type PregancyTestGroup struct {
	Id       string
	TestDate string
	UserId   string
}

func (p *PregancyTestGroup) New(c *CreatePregancyTestGroup) *PregancyTestGroup {
	p = &PregancyTestGroup{
		Id: uuid.NewString(),
        TestDate: c.TestDate,
        UserId: c.UserId,
	}
    return p
}

type CreatePregancyTestGroup struct {
	TestDate string
	UserId   string
}
