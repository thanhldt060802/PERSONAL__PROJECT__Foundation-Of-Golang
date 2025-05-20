package dto

type SimpleRequest struct {
	Body struct {
		Event string `json:"event" required:"true" enum:"a_off_hook,a_on_hook,dialled_no_busy_or_incorrect,b_side_accepts_call,call_from_other_telephone,b_side_answer,other_side_goes_on_hook,own_side_goes_on_hook" doc:"Event to get next state in FSM."`
	}
}
