module geminibb.com/m

go 1.20

require github.com/pitr/gig v0.9.8
require "github.com/userName/otherModule" v0.0.0
replace "github.com/userName/otherModule" v0.0.0 => "./db"

require (
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.1.0 // indirect
)
