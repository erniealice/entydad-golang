package auth

import (
	login02mod "github.com/erniealice/entydad-golang/service/auth/views/login02"
	resetpassword02mod "github.com/erniealice/entydad-golang/service/auth/views/reset-password02"
	signup02mod "github.com/erniealice/entydad-golang/service/auth/views/signup02"
)

// CarouselSlide holds data for a single carousel slide used across the
// auth pages (login, signup, reset-password).
type CarouselSlide struct {
	Title       string
	Description string
}

// DefaultCarouselSlides returns the default carousel slides for login,
// signup, and reset-password pages.
func DefaultCarouselSlides() []CarouselSlide {
	return []CarouselSlide{
		{
			Title:       "Manage your business",
			Description: "Everything you need to run your operations in one place.",
		},
		{
			Title:       "Track your finances",
			Description: "Stay on top of revenue, expenses, and reporting.",
		},
		{
			Title:       "Serve your clients",
			Description: "Deliver great experiences with streamlined workflows.",
		},
	}
}

// toLogin02Slides converts auth CarouselSlides to login02 module slides.
func toLogin02Slides(slides []CarouselSlide) []login02mod.CarouselSlide {
	out := make([]login02mod.CarouselSlide, len(slides))
	for i, s := range slides {
		out[i] = login02mod.CarouselSlide{Title: s.Title, Description: s.Description}
	}
	return out
}

// toSignup02Slides converts auth CarouselSlides to signup02 module slides.
func toSignup02Slides(slides []CarouselSlide) []signup02mod.CarouselSlide {
	out := make([]signup02mod.CarouselSlide, len(slides))
	for i, s := range slides {
		out[i] = signup02mod.CarouselSlide{Title: s.Title, Description: s.Description}
	}
	return out
}

// toResetPassword02Slides converts auth CarouselSlides to resetpassword02 module slides.
func toResetPassword02Slides(slides []CarouselSlide) []resetpassword02mod.CarouselSlide {
	out := make([]resetpassword02mod.CarouselSlide, len(slides))
	for i, s := range slides {
		out[i] = resetpassword02mod.CarouselSlide{Title: s.Title, Description: s.Description}
	}
	return out
}
