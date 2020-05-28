package topics

import (
	"errors"
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
	topicRegexp   = regexp.MustCompile(`^(?:[.A-Za-z0-9])+$`)
	patternRegexp = regexp.MustCompile(`^(?:[.*#A-Za-z0-9])+$`)
)

type Parser struct {
}

func NewParser() Parser {
	return Parser{}
}

func (p Parser) IsStatic(topic string) (bool, error) {
	if err := p.validate(topic); err != nil {
		return false, err
	}
	return topicRegexp.MatchString(topic), nil
}

func (p Parser) ParseTopic(topic string) ([]string, error) {
	return p.parse(topic, topicRegexp)
}

func (p Parser) ParsePattern(pattern string) ([]string, error) {
	return p.parse(pattern, patternRegexp)
}

func (p Parser) parse(str string, exp *regexp.Regexp) ([]string, error) {
	if err := p.validate(str); err != nil {
		return nil, err
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

func (p Parser) validate(topic string) error {
	topicLen := len(topic)
	if topicLen == 0 {
		return ErrTopicEmpty
	}
	if topicLen > 64 {
		return ErrTopicTooLong
	}
	return nil
}
