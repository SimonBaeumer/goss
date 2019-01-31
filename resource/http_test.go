package resource

import (
	"github.com/SimonBaeumer/goss/system"
	"github.com/SimonBaeumer/goss/system/mock_system"
	"github.com/SimonBaeumer/goss/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

const DoNotSkip = false
const ID = ""

func TestNewHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockSystemHTTP := mock_system.NewMockHTTP(ctrl)
	mockSystemHTTP.EXPECT().HTTP().Return("http://goss.rocks")
	mockSystemHTTP.EXPECT().Status().Return(200, nil)
	mockSystemHTTP.EXPECT().Headers().Return(system.Header{}, nil)

	got, _ := NewHTTP(mockSystemHTTP, util.Config{})

	assert.IsType(t, new(HTTP), got)
	assert.Equal(t, "http://goss.rocks", got.HTTP)
}

func Test_validateHeader_ShouldFindKey(t *testing.T) {
	fakeRes := &FakeResource{ID}
	expected := make(map[string][]string)
	expected["key"] = []string{"value"}

	actual := system.Header{}
	actual["key"] = []string{"value"}

	got := validateHeader(fakeRes, "Headers", expected, actual, DoNotSkip)

	assert.IsType(t, new(TestResult), &got)
	assert.Equal(t, SUCCESS, got.Result)
	assert.True(t, got.Successful)
}

func Test_validateHeader_ShouldNotFindKey(t *testing.T) {
	fakeRes := &FakeResource{ID}
	expected := make(map[string][]string)
	expected["key"] = []string{"value"}

	actual := system.Header{}
	actual["another"] = []string{"value"}

	got := validateHeader(fakeRes, "Headers", expected, actual, DoNotSkip)

	assert.IsType(t, new(TestResult), &got)
	assert.Equal(t, FAIL, got.Result)
	assert.False(t, got.Successful)
	assert.Equal(t, "Did not find header key got \nanother: value\n", got.Human)
}

func Test_validateHeader_ShouldFindValue(t *testing.T) {
	fakeRes := &FakeResource{ID}
	expected := make(map[string][]string)
	expected["key"] = []string{"value", "value1"}
	expected["another"] = []string{"value"}

	actual := make(system.Header)
	actual["key"] = []string{"value", "value1"}
	actual["another"] = []string{"value"}

	got := validateHeader(fakeRes, "Header", expected, actual, DoNotSkip)

	assert.Equal(t, SUCCESS, got.Result)
	assert.True(t, got.Successful)
}

func Test_validateHeader_ShouldNotFindValue(t *testing.T) {
	fakeRes := &FakeResource{ID}
	expected := make(map[string][]string)
	expected["key"] = []string{"does not match"}

	actual := make(system.Header)
	actual["key"] = []string{"value"}

	got := validateHeader(fakeRes, "Header", expected, actual, DoNotSkip)

	assert.Equal(t, FAIL, got.Result)
	assert.False(t, got.Successful)
	assert.Equal(t, "Did not find header key: does not match", got.Human)
}

func Test_validateHeader_ShouldNotFindValue_GivenMultiple(t *testing.T) {
	fakeRes := &FakeResource{ID}
	expected := make(map[string][]string)
	expected["key"] = []string{"value", "does not match"}

	actual := make(system.Header)
	actual["key"] = []string{"value"}

	got := validateHeader(fakeRes, "Header", expected, actual, DoNotSkip)

	assert.Equal(t, FAIL, got.Result)
	assert.False(t, got.Successful)
	assert.Equal(t, "Did not find header key: does not match", got.Human)
}

func Test_validateHeader_ShouldFindValuesInComplex(t *testing.T) {
	fakeRes := &FakeResource{ID}
	expected := make(map[string][]string)
	expected["key"] = []string{"value0"}
	expected["another"] = []string{"value1", "value2", "value3"}
	expected["yet"] = []string{"my", "my1", "my2"}

	actual := make(system.Header)
	actual["key"] = []string{"value0"}
	actual["another"] = []string{"value1", "value2", "value3"}
	actual["yet"] = []string{"my", "my1", "my2"}

	got := validateHeader(fakeRes, "Header", expected, actual, DoNotSkip)

	assert.Equal(t, SUCCESS, got.Result)
	assert.True(t, got.Successful)
}

func Test_isInSlice(t *testing.T) {
	haystack := []string{"apple", "banana"}

	got := isInStringSlice(haystack, "banana")
	assert.True(t, got)
}

func Test_isNotInSlice(t *testing.T) {
	haystack := []string{"apple", "banana"}

	got := isInStringSlice(haystack, "pear")
	assert.False(t, got)
}

func Test_ParseYAML(t *testing.T) {
	configString :=  []byte(`
status: 200
allow-insecure: true
no-follow-redirects: true
request-headers:
    key:
    - value
    - value2
`)

	got := new(HTTP)
	err := yaml.Unmarshal(configString, got)

	assert.Nil(t, err)
	assert.Equal(t, 200, got.Status)
	assert.Equal(t, true, got.AllowInsecure)
	assert.Equal(t, true, got.NoFollowRedirects)
	assert.IsType(t, make(map[string][]string), got.RequestHeaders)
	assert.Equal(t, []string{"value", "value2"}, got.RequestHeaders["key"])
}
