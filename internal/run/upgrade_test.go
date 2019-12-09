package run

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUpgrade(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mCmd := NewMockcmd(ctrl)
	originalCommand := Command

	Command = func(path string, args ...string) cmd {
		assert.Equal(t, HELM_BIN, path)
		assert.Equal(t, []string{"upgrade", "--install", "jonas_brothers_only_human", "at40"}, args)

		return mCmd
	}
	defer func() { Command = originalCommand }()

	mCmd.EXPECT().
		Stdout(gomock.Any())
	mCmd.EXPECT().
		Stderr(gomock.Any())
	mCmd.EXPECT().
		Run().
		Times(1)

	u := NewUpgrade("jonas_brothers_only_human", "at40")
	u.Run()
}
