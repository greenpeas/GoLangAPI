package modemLogRaw

import "time"

type Db struct {
	RegTime        time.Time `json:"reg_time" db:"reg_time"`
	Imei           uint64
	Src            int
	Hex            string
	Payload        string
	CmdName        string `json:"cmd_name" db:"cmd_name"`
	CmdDescription string `json:"cmd_description" db:"cmd_description"`
}

type Repo interface {
	List(params ListParams) ([]ModemLogRaw, error)
	ListTelemetry(params ListParams) ([]Telemetry, error)
}

type Usecase interface {
	List(params ListParams) ([]ModemLogRaw, error)
	ListTelemetry(params ListParams) ([]Telemetry, error)
}
