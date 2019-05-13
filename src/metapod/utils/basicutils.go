package utils

import "time"

//ParseUnixTimeOrDie
func ParseUnixTimeOrDie(unixTime string) time.Time {
	t, err := time.Parse(time.UnixDate, unixTime)
	if err != nil {
		panic(err)
	}
	return t
}