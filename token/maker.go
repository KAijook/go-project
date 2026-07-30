package token

//Maker is an interface for creating tokens.
type Maker interface {
	//CreateToken creates a new token for a specific username and duration.
	CreateToken(username string, duration int64) (string, error)

	//VerifyToken checks if the token is valid or not.
	VerifyToken(token string) (*Payload, error)
}
