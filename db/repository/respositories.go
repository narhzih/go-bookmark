package repository

type Repositories struct {
	User                UserRepository
	Pipe                PipeRepository
	PipeShare           PipeShareRepository
	Bookmark            BookmarkRepository
	AccountVerification AccountVerificationRepository
	Notification        NotificationRepository
	PasswordReset       PasswordResetRepository
	Tag                 TagRepository
	Search              SearchRepository
}
