package routes

import wisdomType "github.com/wisdom-oss/commonTypes/v2"

var ErrEmptyWaterRightID = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.1",
	Status: 400,
	Title:  "Empty Water Right ID",
	Detail: "The water right id set in the query is empty. Please check your query",
}

var ErrNoWaterRightAvailable = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.5",
	Status: 404,
	Title:  "No Water Right Available",
	Detail: "The water right id set in the query does not point to a water right available in the database",
}
