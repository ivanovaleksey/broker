package topics

import (
	"errors"
	"github.com/ivanovaleksey/broker/pkg/types"
	"regexp"
	"strings"
)

var (
	ErrTopicEmpty       = errors.New("empty topic")
	ErrTopicTooLong     = errors.New("topic too long")
	ErrTopicInvalidChar = errors.New("topic contains invalid character")
	ErrTopicEmptyPart   = errors.New("topic contains empty part")
)

var (
	topicRegexp = regexp.MustCompile(`^(?:[.A-Za-z0-9])+$`)
	patternRegexp = regexp.MustCompile(`^(?:[.*#A-Za-z0-9])+$`)
)

type Parser struct {
}

func NewParser() Parser {
	return Parser{}
}

func (p Parser) ParseTopic(topic types.Topic) ([]string, error) {
	return p.parse(topic, topicRegexp)
}

func (p Parser) ParsePattern(pattern types.Topic) ([]string, error) {
	return p.parse(pattern, patternRegexp)
}

func (p Parser) parse(str string, exp *regexp.Regexp) ([]string, error) {
	topicLen := len(str)
	if topicLen == 0 {
		return nil, ErrTopicEmpty
	}
	if topicLen > 64 {
		return nil, ErrTopicTooLong
	}
	if !exp.MatchString(str) {
		return nil, ErrTopicInvalidChar
	}
	parts := strings.Split(str, ".")
	for _, part := range parts {
		if len(part) == 0 {
			return nil, ErrTopicEmptyPart
		}
	}
	return parts, nil
}
