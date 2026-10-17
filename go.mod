module github.com/AD-Tech-Projects/Ledger

go 1.25.7

// ANCHOR - direct dependencies
require github.com/AD-Tech-Projects/Go-Libs v0.0.1

// require (
// 	adtech.com/concurrency v0.0.0-unpublished
// 	rsc.io/quote v1.5.2
// )

// ANCHOR - dependencies replaces
// replace adtech.com/concurrency v0.0.0-unpublished => ./modules/concurrency

// ANCHOR - indirect dependencies
// require (
// 	golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c // indirect
// 	rsc.io/sampler v1.3.0 // indirect
// )
