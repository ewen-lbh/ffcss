package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsURLClonable(t *testing.T) {
	var actual bool
	var err error

	actual, err = IsURLClonable("https://github.com/ewen-lbh/ffcss/")
	assert.Equal(t, true, actual)
	assert.Nil(t, err)

	actual, err = IsURLClonable("https://github.com/users/schoolsyst")
	assert.Equal(t, false, actual)
	assert.Nil(t, err)

	actual, err = IsURLClonable("https://ewen.works/")
	assert.Equal(t, false, actual)
	assert.Nil(t, err)
}
