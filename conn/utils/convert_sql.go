package utils

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	"mycfgrest/types"
)

func ChangeSqlToNumBindSupportSql(sql string, p *types.ParsingValue) (query string, param []any, err *types.AppError) {
	fetch, err := p.Fetch()
	if err != nil {
		return "", nil, types.NewAppError(err, "parsing value is fetch error")
	}

	keys := make([]string, fetch.FullSize())
	vals := make([]any, fetch.FullSize())

	fetch.Keys(keys)
	fetch.Values(vals)

	fetch.Close()

	var buffer bytes.Buffer
	lastIndex := 0

	for i := 0; i < len(sql); i++ {
		if sql[i] == '#' {
			if i+1 < len(sql) && sql[i+1] == '#' {
				// '##' 처리
				buffer.WriteString(sql[lastIndex:i])
				buffer.WriteByte('#')
				i++
				lastIndex = i + 1
			} else if i+1 < len(sql) && sql[i+1] == '{' {
				// '#{...}' 처리
				end := strings.IndexByte(sql[i:], '}')
				if end != -1 {
					end += i
					key := sql[i+2 : end]

					if idx := slices.Index(keys, key); idx != -1 {
						buffer.WriteString(sql[lastIndex:i])
						buffer.WriteString(fmt.Sprintf("$%d", idx+1))

						i = end
						lastIndex = end + 1
					}
				}
			}
		}
	}
	buffer.WriteString(sql[lastIndex:])

	return buffer.String(), vals, nil
}

func ChangeSqlToQuestionMarkBindSupportSql(sql string, p *types.ParsingValue) (query string, param []any, err *types.AppError) {
	fetch, err := p.Fetch()
	if err != nil {
		return "", nil, types.NewAppError(err, "parsing value is fetch error")
	}

	keys := make([]string, fetch.FullSize())
	vals := make([]any, fetch.FullSize())

	fetch.Values(vals)
	fetch.Keys(keys)

	fetch.Close()

	var buffer bytes.Buffer
	lastIndex := 0

	for i := 0; i < len(sql); i++ {
		if sql[i] == '#' {
			if i+1 < len(sql) && sql[i+1] == '#' {
				// '##' 처리
				buffer.WriteString(sql[lastIndex:i])
				buffer.WriteByte('#')
				i++
				lastIndex = i + 1
			} else if i+1 < len(sql) && sql[i+1] == '{' {
				// '#{...}' 처리
				end := strings.IndexByte(sql[i:], '}')
				if end != -1 {
					end += i
					key := sql[i+2 : end]

					if exists := slices.Index(keys, key) != -1; exists {
						buffer.WriteString(sql[lastIndex:i])
						buffer.WriteByte('?')
						i = end
						lastIndex = end + 1
					}
				}
			}
		}
	}
	buffer.WriteString(sql[lastIndex:])

	return buffer.String(), vals, nil
}
