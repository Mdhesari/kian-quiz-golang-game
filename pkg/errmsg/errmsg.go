package errmsg

const (
	ErrInternalServer     string = "Internal Server Error."
	ErrNotFound           string = "Resource not found."
	ErrClaimAssertion     string = "Could not assert claims."
	ErrAuthorization      string = "Authorization required."
	ErrSignKey            string = "SignKey Error."
	ErrMobileUnique       string = "Mobile must be unique."
	ErrEmailUnique        string = "Email must be unique."
	ErrInvalidInput       string = "Invalid input validation."
	ErrCategoryNotFound   string = "Category not found."
	ErrGameNotCreated     string = "Game not created."
	ErrGameIDNotConverted string = "Game ID not converted."
	ErrGameNotFound       string = "Game not found."
	ErrGameNotModified    string = "Game not modified."
)
