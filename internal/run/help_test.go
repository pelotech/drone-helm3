package run

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHelp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mCmd := NewMockcmd(ctrl)
	originalCommand := command

	command = func(path string, args ...string) cmd {
		assert.Equal(t, helmBin, path)
		assert.Equal(t, []string{"help"}, args)
		return mCmd
	}
	defer func() { command = originalCommand }()

	mCmd.EXPECT().
		Stdout(gomock.Any())
	mCmd.EXPECT().
		Stderr(gomock.Any())
	mCmd.EXPECT().
		Run().
		Times(1)

	h := NewHelp()
	h.Run()
}
