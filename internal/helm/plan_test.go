package helm

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/pelotech/drone-helm3/internal/run"
)

type PlanTestSuite struct {
	suite.Suite
}

func TestPlanTestSuite(t *testing.T) {
	suite.Run(t, new(PlanTestSuite))
}

func (suite *PlanTestSuite) TestNewPlanUpgradeCommand() {
	cfg := Config{
		Command: "upgrade",
		Chart:   "billboard_top_100",
		Release: "post_malone_circles",
	}

	plan, err := NewPlan(cfg)
	suite.Require().Nil(err)
	suite.Require().Equal(1, len(plan.steps))

	suite.Require().IsType(&run.Upgrade{}, plan.steps[0])
	step, _ := plan.steps[0].(*run.Upgrade)

	suite.Equal("billboard_top_100", step.Chart)
	suite.Equal("post_malone_circles", step.Release)
}

func (suite *PlanTestSuite) TestNewPlanUpgradeFromDroneEvent() {
	cfg := Config{
		Chart:   "billboard_top_100",
		Release: "lizzo_good_as_hell",
	}

	upgradeEvents := []string{"push", "tag", "deployment", "pull_request", "promote", "rollback"}
	for _, event := range upgradeEvents {
		cfg.DroneEvent = event
		plan, err := NewPlan(cfg)
		suite.Require().Nil(err)
		suite.Require().Equal(1, len(plan.steps), fmt.Sprintf("for event type '%s'", event))
		suite.IsType(&run.Upgrade{}, plan.steps[0], fmt.Sprintf("for event type '%s'", event))
	}
}

func (suite *PlanTestSuite) TestNewPlanHelpCommand() {
	cfg := Config{
		Command: "help",
	}

	plan, err := NewPlan(cfg)
	suite.Require().Nil(err)
	suite.Equal(1, len(plan.steps))

	suite.Require().IsType(&run.Help{}, plan.steps[0])
}
