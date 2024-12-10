package service

// pingDatabaseService пингование БД.
func (sv *pingDatabaseService) Ping() error {
	return sv.st.Ping()
}
