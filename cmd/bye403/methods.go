package main

// skip Delete
func (b *bye403) methods() []string {
	return []string{
		"ACL", "BASELINE-CONTROL", "BIND", "CHECKIN",
		"CHECKOUT", "CONNECT", "COPY", "GET", "HEAD",
		"LABEL", "LINK", "LOCK", "MERGE", "MKACTIVITY",
		"MKCALENDAR", "MKCOL", "MKREDIRECTREF", "MKWORKSPACE",
		"MOVE", "OPTIONS", "ORDERPATCH", "PATCH", "POST", "PRI",
		"PROPFIND", "PROPPATCH", "PUT", "QUERY", "REBIND",
		"REPORT", "SEARCH", "TRACE", "UNBIND", "UNCHECKOUT",
		"UNLINK", "UNLOCK", "UPDATE", "UPDATEREDIRECTREF",
		"VERSION-CONTROL", "FOOBAR",
	}
}
