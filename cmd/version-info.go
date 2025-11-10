package cmd

// These variables are set at build time using ldflags - so embedded in the binary! ✨
var (
	Version = "0.1.0" // Version number - embedded at build time! 💕
	Build   = "0"     // Build number - embedded at build time! 🎀
)
