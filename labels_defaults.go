package entydad

// labels_defaults.go — default-constructor functions for auth label structs.
//
// These constructors centralise the English fallback strings that previously
// lived as hardcoded literals in the app's composition/translations.go.
// The app calls e.g. entydad.DefaultLogin02Labels() and then overlays it
// with translations.LoadPathIfExists("en", bt, "auth.json", "login02", &l.Login02).
// Identical values; identical JSON overlay behaviour.

// DefaultLogin02Labels returns Login02Labels populated with English defaults.
func DefaultLogin02Labels() Login02Labels {
	return Login02Labels{
		Title:               "Sign In",
		Heading:             "Welcome back",
		Subheading:          "Sign in to your account",
		EmailLabel:          "Email",
		EmailPlaceholder:    "Enter your email",
		PasswordLabel:       "Password",
		PasswordPlaceholder: "Enter your password",
		RememberMe:          "Remember me",
		ForgotPassword:      "Forgot password?",
		SignInButton:        "Sign In",
		NoAccount:           "Don't have an account?",
		SignUpLink:          "Sign up",
		SocialDivider:       "or continue with",
		Error:               "Invalid email or password. Please try again.",
		PreviousSlide:       "Previous slide",
		NextSlide:           "Next slide",
		ContinueWith:        "Continue with",
	}
}

// DefaultSignup02Labels returns Signup02Labels populated with English defaults.
func DefaultSignup02Labels() Signup02Labels {
	return Signup02Labels{
		Title:                      "Create Account",
		Heading:                    "Get started",
		Subheading:                 "Create your account",
		FirstNameLabel:             "First Name",
		FirstNamePlaceholder:       "Enter your first name",
		LastNameLabel:              "Last Name",
		LastNamePlaceholder:        "Enter your last name",
		EmailLabel:                 "Email",
		EmailPlaceholder:           "Enter your email",
		PasswordLabel:              "Password",
		PasswordPlaceholder:        "Create a password",
		ConfirmPasswordLabel:       "Confirm Password",
		ConfirmPasswordPlaceholder: "Confirm your password",
		SignUpButton:               "Create Account",
		HasAccount:                 "Already have an account?",
		SignInLink:                 "Sign in",
		SocialDivider:              "or sign up with",
		TermsText:                  "By signing up, you agree to our Terms and Privacy Policy.",
		Error:                      "Registration failed. Please try again.",
		PreviousSlide:              "Previous slide",
		NextSlide:                  "Next slide",
		ContinueWith:               "Continue with",
		PasswordStrength:           "Password strength",
		TermsLink:                  "Terms",
	}
}

// DefaultResetPassword02Labels returns ResetPassword02Labels populated with English defaults.
func DefaultResetPassword02Labels() ResetPassword02Labels {
	return ResetPassword02Labels{
		Title:                      "Reset Password",
		Heading:                    "Forgot your password?",
		Subheading:                 "Enter your email and we'll send you a reset link.",
		EmailLabel:                 "Email",
		EmailPlaceholder:           "Enter your email",
		SendResetButton:            "Send Reset Link",
		BackToLogin:                "Back to sign in",
		ConfirmHeading:             "Set new password",
		ConfirmSubheading:          "Enter your new password below.",
		NewPasswordLabel:           "New Password",
		NewPasswordPlaceholder:     "Enter your new password",
		ConfirmPasswordLabel:       "Confirm Password",
		ConfirmPasswordPlaceholder: "Confirm your new password",
		ResetButton:                "Reset Password",
		SuccessHeading:             "Password reset sent",
		SuccessMessage:             "If that email exists in our system, you'll receive a reset link shortly.",
		Error:                      "Password reset failed. Please try again.",
		ErrorMismatch:              "The two passwords don't match. Please enter the same value in both fields.",
		ErrorInvalidToken:          "This reset link is no longer valid. Request a new one to continue.",
		ErrorExpiredToken:          "This reset link has expired. Request a new one to continue.",
		ErrorWeakPassword:          "Your new password is too short. Choose at least 8 characters.",
		PreviousSlide:              "Previous slide",
		NextSlide:                  "Next slide",
	}
}

// DefaultChangePasswordLabels returns ChangePasswordLabels populated with English defaults.
func DefaultChangePasswordLabels() ChangePasswordLabels {
	return ChangePasswordLabels{
		Title:                      "Change Password",
		Heading:                    "Change your password",
		Subheading:                 "Enter your current password and choose a new one.",
		OldPasswordLabel:           "Current Password",
		OldPasswordPlaceholder:     "Enter your current password",
		NewPasswordLabel:           "New Password",
		NewPasswordPlaceholder:     "Enter your new password",
		ConfirmPasswordLabel:       "Confirm New Password",
		ConfirmPasswordPlaceholder: "Confirm your new password",
		SubmitButton:               "Change Password",
		SuccessMessage:             "Password changed successfully",
		Error:                      "We couldn't change your password. Please try again.",
		ErrorMismatch:              "The two new-password fields don't match. Please retype to confirm.",
		ErrorCurrentIncorrect:      "Current password is incorrect",
		ErrorTooShort:              "New password must be at least 8 characters",
		BackToApp:                  "Back to dashboard",
	}
}
