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
	ErrWeakPassword       string = "Password does not meet security requirements. Please use a stronger password with at least 8 characters, including uppercase and lowercase letters, numbers, and special characters."
	ErrAlreadyAnswered    string = "Player has already answered this question."
	ErrGamePlayerNotFound string = "Player not found in the game."
	ErrQuestionNotFound   string = "Question not found."
	ErrNoCorrectAnswer    string = "No correct answer found."
)
