package providers

type CryptoProvider interface {
	Encrypt(password string) (string, error)
	Compare(password, hashedPassword string) (bool, error)
}
