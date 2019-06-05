package errors

const (
	DocumentNotFoundError = "mongo: no documents in result"
)

func IsDocumentNotFound(e error) bool {
	return e.Error() == DocumentNotFoundError
}
