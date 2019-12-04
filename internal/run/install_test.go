package run

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestInstall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cmd := NewMockcmd(ctrl)
	cmd.EXPECT().
		Args(gomock.Eq([]string{"install", "arg1", "arg2"}))
	cmd.EXPECT().
		Run().
		Times(1)

	install(cmd, []string{"arg1", "arg2"})
}
