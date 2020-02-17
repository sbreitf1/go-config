package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathPrefix(t *testing.T) {
	assert.Equal(t, "Test[7].SubItem[2]", newPathPrefix("Test").Index(7).Field2("SubItem", "Item").Index(2).String())
	assert.Equal(t, "TEST_7_ITEM_2", newPathPrefix("Test").Index(7).Field2("SubItem", "Item").Index(2).Env())
	assert.Equal(t, "[0]", newPathPrefix("").Index(0).String())
	assert.Equal(t, "0", newPathPrefix("").Index(0).Env())
	assert.Equal(t, "Prefix", newPathPrefix("Prefix").String())
	assert.Equal(t, "PREFIX", newPathPrefix("Prefix").Env())
	assert.Equal(t, "Test[0][1][2]", newPathPrefix("Test").Index(0).Index(1).Index(2).String())
	assert.Equal(t, "TEST_0_1_2", newPathPrefix("Test").Index(0).Index(1).Index(2).Env())
	assert.Equal(t, "Test.SubItem[2].Num", newPathPrefix("Test").Field2("SubItem", "Item").Index(2).Field("Num").String())
	assert.Equal(t, "TEST_ITEM_2_NUM", newPathPrefix("Test").Field2("SubItem", "Item").Index(2).Field("Num").Env())
}
