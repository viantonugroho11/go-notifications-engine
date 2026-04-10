package email


// Message parameter untuk kirim email.
type Message struct {
	From    string   // pengirim (opsional, bisa pakai default dari dialer)
	To      []string // penerima
	Subject string
	Body    string // body plain text
}