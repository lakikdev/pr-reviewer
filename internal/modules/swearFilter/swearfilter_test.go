package swearFilter_test

import (
	"flag"
	"testing"

	"pr-reviewer/internal/modules/swearFilter"

	"github.com/stretchr/testify/suite"
)

type SwearFilterTestSuite struct {
	suite.Suite
}

func (s *SwearFilterTestSuite) SetupSuite() {
	flag.Parse()
}

func (s *SwearFilterTestSuite) TestIsClean() {
	valid := swearFilter.Check("clean")
	s.Require().Equal(valid, true)
}

func (s *SwearFilterTestSuite) TestIsClean_Dirty() {
	valid := swearFilter.Check("ass")
	s.Require().Equal(valid, false)
}

func TestSwearFilterTestSuite(t *testing.T) {
	suite.Run(t, new(SwearFilterTestSuite))
}
