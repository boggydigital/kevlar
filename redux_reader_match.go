package kevlar

import (
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"strings"
)

func (rdx *redux) MatchAsset(asset string, terms []string, scope []string, options ...MatchOption) []string {
	if scope == nil {
		scope = rdx.Keys(asset)
	}

	matches := make(map[string]interface{})
	for _, term := range terms {
		if !slices.Contains(options, CaseSensitive) {
			term = strings.ToLower(term)
		}
		for _, key := range scope {
			if values, ok := rdx.GetAllValues(asset, key); !ok {
				continue
			} else if anyValueMatchesTerm(term, values, options...) {
				matches[key] = nil
			}
		}
	}

	return maps.Keys(matches)
}

func (rdx *redux) Match(query map[string][]string, options ...MatchOption) []string {
	var matches []string
	for asset, terms := range query {
		if !rdx.HasAsset(asset) {
			continue
		}
		matches = rdx.MatchAsset(asset, terms, matches, options...)
	}
	return matches
}

func anyValueMatchesTerm(term string, values []string, options ...MatchOption) bool {
	anyCase := true
	contains := true

	if options != nil {
		anyCase = !slices.Contains(options, CaseSensitive)
		contains = !slices.Contains(options, FullMatch)
	}

	for _, val := range values {
		if anyCase {
			val = strings.ToLower(val)
		}
		if contains {
			if strings.Contains(val, term) {
				return true
			}
		} else {
			if val == term {
				return true
			}
		}
	}
	return false
}
