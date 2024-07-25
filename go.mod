module github.com/iamdlfl/sunday-worship-playlist

go 1.22.5

require gopkg.in/ini.v1 v1.67.0 // indirect

replace "github.com/iamdlfl/spotify" => ./spotify
require "github.com/iamdlfl/spotify" v0.0.0

replace "github.com/iamdlfl/pco" => ./pco
require "github.com/iamdlfl/pco" v0.0.0