package run

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestInstall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mCmd := NewMockcmd(ctrl)
	originalCommand := Command
	Command = func() cmd { return mCmd }
	defer func() { Command = originalCommand }()

	mCmd.EXPECT().
		Path(HELM_BIN)
	mCmd.EXPECT().
		Args(gomock.Eq([]string{"install", "arg1", "arg2"}))
	mCmd.EXPECT().
		Run().
		Times(1)

	Install("arg1", "arg2")
}
