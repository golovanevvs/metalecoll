package service

func (sv *PingDatabaseService) Ping() error {
	return sv.st.Ping()
}
